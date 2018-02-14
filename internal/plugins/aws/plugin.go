package aws

import (
	"unicode/utf8"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/gogo/protobuf/types"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/solo-io/glue/internal/plugins"
	"github.com/solo-io/glue/pkg/secretwatcher"

	"github.com/solo-io/glue/internal/plugins/common"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/plugin"
)

func init() {
	plugins.Register(&Plugin{}, nil)
}

type Plugin struct{}

const (
	// define Upstream type name
	UpstreamTypeAws = "aws"

	// generic plugin info
	filterName  = "io.solo.aws"
	pluginStage = plugin.OutAuth

	// filter-specific metadata
	filterMetadataKeyAsync = "async"

	// upstream-specific metadata
	awsAccessKey = "access_key"
	awsSecretKey = "secret_key"
	awsRegion    = "region"
	awsHost      = "host"

	// function-specific metadata
	functionNameKey      = "name"
	functionQualifierKey = "qualifier"
)

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	deps := new(plugin.Dependencies)
	for _, upstream := range cfg.Upstreams {
		if upstream.Type != UpstreamTypeAws {
			continue
		}
		awsUpstream, err := DecodeUpstreamSpec(upstream.Spec)
		if err != nil {
			// errors will be handled during validation
			// TODO: consider logging error here
			continue
		}
		deps.SecretRefs = append(deps.SecretRefs, awsUpstream.SecretRef)
	}
	return deps
}

func (p *Plugin) HttpFilter() (*envoyhttp.HttpFilter, plugin.Stage) {
	return &envoyhttp.HttpFilter{Name: filterName}, pluginStage
}

func (p *Plugin) ProcessRoute(in *v1.Route, out *envoyroute.Route) error {
	executionStyle, err := GetExecutionStyle(in.Extensions)
	if err != nil {
		return err
	}
	setRouteAsync(executionStyle == ExecutionStyleAsync, out)
	return nil
}

func setRouteAsync(async bool, out *envoyroute.Route) {
	common.InitFilterMetadataField(filterName, filterMetadataKeyAsync, out.Metadata)
	out.Metadata.FilterMetadata[filterName].Fields[filterMetadataKeyAsync].Kind = &types.Value_BoolValue{
		BoolValue: async,
	}
}

func (p *Plugin) ProcessUpstream(in *v1.Upstream, secrets secretwatcher.SecretMap, out *envoyapi.Cluster) error {
	if in.Type != UpstreamTypeAws {
		return nil
	}

	awsUpstream, err := DecodeUpstreamSpec(in.Spec)
	if err != nil {
		return errors.Wrap(err, "invalid AWS upstream spec")
	}

	awsSecrets, ok := secrets[awsUpstream.SecretRef]
	if !ok {
		return errors.Errorf("aws secrets for ref %v not found", awsUpstream.SecretRef)
	}
	var secretErrs error

	accessKey, ok := awsSecrets[awsAccessKey]
	if !ok {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("key %v missing from provided secret", awsAccessKey))
	}
	if accessKey != "" && !utf8.Valid([]byte(accessKey)) {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("%s not a valid string", awsAccessKey))
	}
	secretKey, ok := awsSecrets[awsSecretKey]
	if !ok {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("key %v missing from provided secret", awsSecretKey))
	}
	if secretKey != "" && !utf8.Valid([]byte(secretKey)) {
		secretErrs = multierror.Append(secretErrs, errors.Errorf("%s not a valid string", awsSecretKey))
	}
	if secretErrs != nil {
		return secretErrs
	}

	common.InitFilterMetadata(filterName, out.Metadata)
	out.Metadata.FilterMetadata[filterName] = &types.Struct{
		Fields: map[string]*types.Value{
			awsAccessKey: {Kind: &types.Value_StringValue{StringValue: accessKey}},
			awsSecretKey: {Kind: &types.Value_StringValue{StringValue: secretKey}},
			awsRegion:    {Kind: &types.Value_StringValue{StringValue: awsUpstream.Region}},
			awsHost:      {Kind: &types.Value_StringValue{StringValue: awsUpstream.GetLambdaHostname()}},
		},
	}

	out.Type = envoyapi.Cluster_LOGICAL_DNS
	// need to make sure we use ipv4 only dns
	out.DnsLookupFamily = envoyapi.Cluster_V4_ONLY

	return nil
}

func (p *Plugin) ParseFunctionSpec(upstreamType string, in v1.FunctionSpec) (*types.Struct, error) {
	if upstreamType != UpstreamTypeAws {
		return nil, nil
	}
	functionSpec, err := DecodeFunctionSpec(in)
	if err != nil {
		return nil, errors.Wrap(err, "invalid lambda function spec")
	}
	return &types.Struct{
		Fields: map[string]*types.Value{
			functionNameKey:      {Kind: &types.Value_StringValue{StringValue: functionSpec.FunctionName}},
			functionQualifierKey: {Kind: &types.Value_StringValue{StringValue: functionSpec.Qualifier}},
		},
	}, nil
}
