package aws

import (
	"fmt"
	"unicode/utf8"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoyendpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/pluginutils"
)

//go:generate protoc -I$GOPATH/src/github.com/lyft/protoc-gen-validate -I. -I$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types:${GOPATH}/src/ filter.proto

const (
	// filter info
	FilterName  = "io.solo.aws_lambda"
	pluginStage = plugins.OutAuth

	// cluster info
	AccessKey = "access_key"
	SecretKey = "secret_key"
	awsRegion = "region"
	awsHost   = "host"
)

func getLambdaHostname(s *aws.UpstreamSpec) string {
	return fmt.Sprintf("lambda.%s.amazonaws.com", s.Region)
}

func init() {
	plugins.RegisterFunc(NewAwsPlugin)
}

func NewAwsPlugin() plugins.Plugin {
	return &plugin{recordedUpstreams: make(map[string]*aws.UpstreamSpec)}
}

type plugin struct {
	recordedUpstreams map[string]*aws.UpstreamSpec
}

func (p *plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *envoyapi.Cluster) error {
	upstreamSpec, ok := in.UpstreamSpec.UpstreamType.(*v1.UpstreamSpec_Aws)
	if !ok {
		// not ours
		return nil
	}

	lambdaHostname := getLambdaHostname(upstreamSpec.Aws)

	// configure Envoy cluster routing info
	out.Type = envoyapi.Cluster_LOGICAL_DNS
	// TODO(yuval-k): why do we need to make sure we use ipv4 only dns?
	out.DnsLookupFamily = envoyapi.Cluster_V4_ONLY
	out.LoadAssignment = &envoyapi.ClusterLoadAssignment{
		ClusterName: out.Name,
		Endpoints: []envoyendpoint.LocalityLbEndpoints{
			{
				LbEndpoints: []envoyendpoint.LbEndpoint{
					{
						Endpoint: pluginutils.EnvoyEndpoint(lambdaHostname, 443),
					},
				},
			},
		},
	}
	out.TlsContext = &envoyauth.UpstreamTlsContext{
		// TODO(yuval-k): Add verification context
		Sni: lambdaHostname,
	}

	// TODO(ilacakrms): consider if secretRef should be namespace+name
	awsSecrets, err := params.Snapshot.SecretList.Find("", upstreamSpec.Aws.SecretRef)
	if err != nil {
		return errors.Wrapf(err, "retrieving aws secret")
	}
	var secretErrs error

	accessKey, ok := awsSecrets.Data[AccessKey]
	if !ok {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("key %v missing from provided secret", AccessKey))
	}
	if accessKey == "" || !utf8.Valid([]byte(accessKey)) {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("%s not a valid string", AccessKey))
	}
	secretKey, ok := awsSecrets.Data[SecretKey]
	if !ok {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("key %v missing from provided secret", SecretKey))
	}
	if secretKey == "" || !utf8.Valid([]byte(secretKey)) {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("%s not a valid string", SecretKey))
	}

	if secretErrs != nil {
		return secretErrs
	}

	lpe := &LambdaProtocolExtension{
		Host:      lambdaHostname,
		Region:    upstreamSpec.Aws.Region,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	if out.ExtensionProtocolOptions == nil {
		out.ExtensionProtocolOptions = make(map[string]*types.Struct)
	}

	lpeStruct, err := util.MessageToStruct(lpe)
	if err != nil {
		return errors.Wrapf(err, "converting aws protocol options to struct")
	}
	out.ExtensionProtocolOptions[FilterName] = lpeStruct

	// TODO(yuval-k): What about namespace?!
	p.recordedUpstreams[in.Metadata.Name] = upstreamSpec.Aws

	return nil
}

func (p *plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	return pluginutils.MarkPerFilterConfig(in, out, FilterName, func(spec *v1.Destination) (proto.Message, error) {
		// check if it's aws destination
		if spec.DestinationSpec == nil {
			return nil, nil
		}
		awsDestinationSpec, ok := spec.DestinationSpec.DestinationType.(*v1.DestinationSpec_Aws)
		if !ok {
			return nil, nil
		}
		// get upstream
		lambdaSpec, ok := p.recordedUpstreams[spec.UpstreamName]
		if !ok {
			// TODO(yuval-k): panic in debug
			return nil, errors.Errorf("%v is not an AWS upstream", spec.UpstreamName)
		}
		// should be aws upstream

		// get function
		logicalName := awsDestinationSpec.Aws.LogicalName
		for _, lambdaFunc := range lambdaSpec.LambdaFunctions {
			if lambdaFunc.LogicalName == logicalName {

				lambdaRouteFunc := &LambdaPerRoute{
					Async:     awsDestinationSpec.Aws.InvocationStyle == aws.DestinationSpec_ASYNC,
					Qualifier: lambdaFunc.Qualifier,
					Name:      lambdaFunc.LambdaFunctionName,
				}

				return lambdaRouteFunc, nil
			}
		}
		return nil, errors.Errorf("unknown function %v", logicalName)
	})
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if len(p.recordedUpstreams) == 0 {
		// no upstreams no filter
		return nil, nil
	}
	return []plugins.StagedHttpFilter{
		{
			HttpFilter: &envoyhttp.HttpFilter{Name: FilterName},
			Stage:      pluginStage,
		},
	}, nil
}
