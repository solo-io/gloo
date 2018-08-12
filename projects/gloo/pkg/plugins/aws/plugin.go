package aws

import (
	"fmt"
	"unicode/utf8"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/pluginutils"
)

const (
	// filter info
	filterName  = "io.solo.lambda"
	pluginStage = plugins.OutAuth

	// cluster info
	AccessKey = "access_key"
	SecretKey = "secret_key"
	awsRegion = "region"
	awsHost   = "host"
)

func init() {
	plugins.Register(&plugin{recordedUpstreams: make(map[string]*v1.Upstream)})
}

type plugin struct {
	recordedUpstreams map[string]*UpstreamSpec
}

func (s *UpstreamSpec) getLambdaHostname() string {
	return fmt.Sprintf("lambda.%s.amazonaws.com", s.Region)
}

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *envoyapi.Cluster) error {
	upstreamSpec, ok := in.UpstreamType.UpstreamType.(*v1.UpstreamSpec_Aws)
	if !ok {
		// not ours
		return nil
	}

	lambdaHostname := upstreamSpec.Aws.getLambdaHostname()

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

	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	common.InitFilterMetadata(filterName, out.Metadata)
	out.Metadata.FilterMetadata[filterName] = &types.Struct{
		Fields: map[string]*types.Value{
			AccessKey: {Kind: &types.Value_StringValue{StringValue: accessKey}},
			SecretKey: {Kind: &types.Value_StringValue{StringValue: secretKey}},
			awsRegion: {Kind: &types.Value_StringValue{StringValue: upstreamSpec.Aws.Region}},
			awsHost:   {Kind: &types.Value_StringValue{StringValue: lambdaHostname}},
		},
	}

	p.recordedUpstreams[in.Metadata.Name] = upstreamSpec.Aws

	return secretErrs

	return nil
}

func (p *plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	return pluginutils.MarkPerFilterConfig(in, out, filterName, func(spec *v1.Destination) (proto.Message, error) {
		// check if it's aws destination
		awsDestinationSpec, ok := spec.DestinationSpec.DestinationType.(*v1.DestinationSpec_Aws)
		if !ok {
			return nil, nil
		}
		// get upstream
		lambdaSpec, ok := p.recordedUpstreams[spec.UpstreamName]
		if !ok {
			return nil, errors.Errorf("%v is not an AWS upstream", spec.UpstreamName)
		}
		// should be aws upstream

		// get function
		logicalName := awsDestinationSpec.Aws.LogicalName
		for _, lambdaFunc := range lambdaSpec.LambdaFunctions {
			if lambdaFunc.LogicalName == logicalName {
				// TODO(ilackarms, yuval-k): return the expected message to the Filter
				return lambdaFunc, nil
			}
		}
		return nil, errors.Errorf("unknown function %v", logicalName)
	})
}

func (p *plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	// flush cache
	defer func() { p.recordedUpstreams = make(map[string]*UpstreamSpec) }()
}

// const (
// 	// define Upstream type name
// 	UpstreamTypeAws = "aws"
//
// 	// generic plugin info
// 	filterName  = "io.solo.lambda"
// 	pluginStage = plugins.OutAuth
//
// 	// filter-specific metadata
// 	filterMetadataKeyAsync = "async"
//
// 	// upstream-specific metadata
// 	AccessKey = "access_key"
// 	SecretKey = "secret_key"
// 	awsRegion    = "region"
// 	awsHost      = "host"
//
// 	// function-specific metadata
// 	functionNameKey      = "name"
// 	functionQualifierKey = "qualifier"
// )
//
// func (p *plugin) aHttpFilters(params
// plugins.Params) ([]plugins.StagedHttpFilter, error) {
// defer func () { p.isNeeded = false }()
//
// if p.isNeeded {
// return []plugins.StagedHttpFilter{{HttpFilter: &envoyhttp.HttpFilter{Name: filterName}, Stage: pluginStage}}, nil
// }
// return nil, nil
// }
//
// func (p *plugin) ProcessRoute(_ *plugins.RoutePluginParams,
// 	in *v1.Route, out *envoyroute.Route) error {
// 	executionStyle, err := GetExecutionStyle(in.Extensions)
// 	if err != nil {
// 		return err
// 	}
// 	if executionStyle == ExecutionStyleNone {
// 		return nil
// 	}
// 	setRouteAsync(executionStyle == ExecutionStyleAsync, out)
// 	return nil
// }
//
// func
// setRouteAsync(async
// bool, out * envoyroute.Route) {
// if out.Metadata == nil {
// out.Metadata = &envoycore.Metadata{}
// }
// common.InitFilterMetadataField(filterName, filterMetadataKeyAsync, out.Metadata).Kind = &types.Value_BoolValue{
// BoolValue: async,
// }
// }
//
// func (p *plugin) ProcessUpstream(params *plugins.UpstreamPluginParams,
// 	in *v1.Upstream, out *envoyapi.Cluster) error {
// 	if in.Type != UpstreamTypeAws {
// 		return nil
// 	}
// 	p.isNeeded = true
//
// 	out.Type = envoyapi.Cluster_LOGICAL_DNS
// 	// need to make sure we use ipv4 only dns
// 	out.DnsLookupFamily = envoyapi.Cluster_V4_ONLY
//
// 	awsUpstream, err := DecodeUpstreamSpec(in.Spec)
// 	if err != nil {
// 		return errors.Wrap(err, "invalid AWS upstream spec")
// 	}
//
// 	out.Hosts = append(out.Hosts, &envoycore.Address{Address: &envoycore.Address_SocketAddress{SocketAddress: &envoycore.SocketAddress{
// 		Address:       awsUpstream.GetLambdaHostname(),
// 		PortSpecifier: &envoycore.SocketAddress_PortValue{PortValue: 443},
// 	}}})
// 	out.TlsContext = &envoyauth.UpstreamTlsContext{
// 		Sni: awsUpstream.GetLambdaHostname(),
// 	}
//
// 	awsSecrets, ok := params.Secrets[awsUpstream.SecretRef]
// 	if !ok {
// 		return errors.Errorf("aws secrets for ref %v not found", awsUpstream.SecretRef)
// 	}
//
// 	var secretErrs error
//
// 	accessKey, ok := awsSecrets.Data[AccessKey]
// 	if !ok {
// 		secretErrs = multierror.Append(secretErrs, errors.Errorf("key %v missing from provided secret", AccessKey))
// 	}
// 	if accessKey == "" || !utf8.Valid([]byte(accessKey)) {
// 		secretErrs = multierror.Append(secretErrs, errors.Errorf("%s not a valid string", AccessKey))
// 	}
// 	secretKey, ok := awsSecrets.Data[SecretKey]
// 	if !ok {
// 		secretErrs = multierror.Append(secretErrs, errors.Errorf("key %v missing from provided secret", SecretKey))
// 	}
// 	if secretKey == "" || !utf8.Valid([]byte(secretKey)) {
// 		secretErrs = multierror.Append(secretErrs, errors.Errorf("%s not a valid string", SecretKey))
// 	}
//
// 	if out.Metadata == nil {
// 		out.Metadata = &envoycore.Metadata{}
// 	}
// 	common.InitFilterMetadata(filterName, out.Metadata)
// 	out.Metadata.FilterMetadata[filterName] = &types.Struct{
// 		Fields: map[string]*types.Value{
// 			AccessKey: {Kind: &types.Value_StringValue{StringValue: accessKey}},
// 			SecretKey: {Kind: &types.Value_StringValue{StringValue: secretKey}},
// 			awsRegion:    {Kind: &types.Value_StringValue{StringValue: awsUpstream.Region}},
// 			awsHost:      {Kind: &types.Value_StringValue{StringValue: awsUpstream.GetLambdaHostname()}},
// 		},
// 	}
//
// 	return secretErrs
// }
//
// func (p *plugin) ParseFunctionSpec(params *plugins.FunctionPluginParams,
// 	in
// v1.FunctionSpec) (*types.Struct, error) {
// if params.UpstreamType != UpstreamTypeAws {
// return nil, nil
// }
// functionSpec, err := DecodeFunctionSpec(in)
// if err != nil {
// return nil, errors.Wrap(err, "invalid lambda function spec")
// }
// return &types.Struct{
// Fields: map[string]*types.Value{
// functionNameKey:      {Kind: &types.Value_StringValue{StringValue: functionSpec.FunctionName}},
// functionQualifierKey: {Kind: &types.Value_StringValue{StringValue: functionSpec.Qualifier}},
// },
// }, nil
// }
