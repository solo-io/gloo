package aws

import (
	"errors"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/glue/implemented_modules/aws/pkg/function"
	awspluginspec "github.com/solo-io/glue/implemented_modules/aws/pkg/plugin"
	"github.com/solo-io/glue/implemented_modules/aws/pkg/upstream"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator"
	"github.com/solo-io/glue/pkg/translator/plugin"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"github.com/hashicorp/go-multierror"
)

const (
	AwsPluginType      = "aws"
	AwsSecretAccessKey = "access_key"
	AwsSecretSecretKey = "secret_key"
	AwsRegionKey       = "region"
	AwsHostKey         = "host"
)

const AwsFunctionNameKey = "name"
const AwsFunctionQualifierKey = "qualifier"

const AwsFilterName = "io.solo.aws"
const AwsAsyncKey = "async"
const AwsFilterStage = plugin.OutAuth

type AwsPlugin struct {
}

func New() plugin.Plugin {
	return &AwsPlugin{}
}

func (a *AwsPlugin) GetDependencies(cfg *v1.Config) translator.DependenciesDescription {
	var deps plugin.SimpleDependenciesDescription
	// go over all upstream and get secrets.
	for _, u := range cfg.Upstreams {
		if !a.IsMyUpstream(&u) {
			continue
		}

		spec, err := upstream.FromMap(u.Spec)
		if err != nil {
			// no need to return error now, as this will happen during validation
			continue
		}

		deps.AddSecretRef(spec.Secret)
	}
	return &deps
}

func (a *AwsPlugin) EnvoyFilters(pi *plugin.PluginInputs) []plugin.FilterWrapper {
	filter := plugin.FilterWrapper{
		Filter: hcm.HttpFilter{
			Name: AwsFilterName,
		},
		Stage: AwsFilterStage,
	}
	return []plugin.FilterWrapper{filter}
}

func (a *AwsPlugin) UpdateEnvoyRoute(pi *plugin.PluginInputs, in *v1.Route, out *apiroute.Route) error {
	// nothing here

	plugins := in.Plugins
	if plugins != nil {
		awsplugin := plugins[AwsPluginType]
		spec, err := awspluginspec.FromMap(awsplugin)
		if err != nil {
			return err
		}
		if spec.Async {

			if out.Metadata == nil {
				out.Metadata = &envoy_api_v2_core.Metadata{
					FilterMetadata: make(map[string]*types.Struct),
				}
			}

			if out.Metadata.FilterMetadata == nil {
				out.Metadata.FilterMetadata = make(map[string]*types.Struct)
			}

			if out.Metadata.FilterMetadata[AwsFilterName] == nil {
				out.Metadata.FilterMetadata[AwsFilterName] = &types.Struct{
					Fields: make(map[string]*types.Value),
				}
			}
			out.Metadata.FilterMetadata[AwsFilterName].Fields[AwsAsyncKey] = &types.Value{
				Kind: &types.Value_BoolValue{BoolValue: spec.Async},
			}
		}

	}

	return nil
}

func (a *AwsPlugin) UpdateFunctionToEnvoyCluster(pi *plugin.PluginInputs, in *v1.Upstream, infunc *v1.Function, out *api.Cluster) error {
	// validate the spec
	_, err := function.FromMap(in.Spec)
	// no need to udpate the cluster, as we implement functional filter.
	return err
}

func (a *AwsPlugin) UpdateEnvoyCluster(pi *plugin.PluginInputs, in *v1.Upstream, out *api.Cluster) error {

	if !a.IsMyUpstream(in) {
		return nil
	}

	spec, err := upstream.FromMap(in.Spec)
	if err != nil {
		return err
	}

	out.Type = api.Cluster_LOGICAL_DNS
	out.Hosts = append(out.Hosts, &envoy_api_v2_core.Address{Address: &envoy_api_v2_core.Address_SocketAddress{SocketAddress: &envoy_api_v2_core.SocketAddress{
		Address: spec.GetLambdaHostname(),
	}}})

	secret := pi.State.Dependencies.Secrets()[spec.Secret]
	if secret == nil {
		return errors.New("secret not found")
	}

	var specerror error

	var accesskey []byte
	var secretkey []byte
	var ok bool
	if accesskey, ok = secret[AwsSecretAccessKey]; !ok {
		specerror = multierror.Append(specerror, errors.New("secret doesn't contain "+AwsSecretAccessKey))
	}
	if !utf8.Valid(accesskey) {
		specerror = multierror.Append(specerror, errors.New(AwsSecretAccessKey+" is not a valid string"))
	}

	if secretkey, ok = secret[AwsSecretSecretKey]; !ok {
		specerror = multierror.Append(specerror, errors.New("secret doesn't contain "+AwsSecretSecretKey))
	}
	if !utf8.Valid(secretkey) {
		specerror = multierror.Append(specerror, errors.New(AwsSecretSecretKey+" is not a valid string"))
	}

	if specerror != nil {
		return specerror
	}

	if out.Metadata == nil {
		out.Metadata = &envoy_api_v2_core.Metadata{
			FilterMetadata: make(map[string]*types.Struct),
		}
	}
	awsstruct := &types.Struct{Fields: make(map[string]*types.Value)}
	out.Metadata.FilterMetadata[AwsFilterName] = awsstruct

	awsstruct.Fields[AwsSecretAccessKey].Kind = &types.Value_StringValue{StringValue: string(accesskey)}
	awsstruct.Fields[AwsSecretSecretKey].Kind = &types.Value_StringValue{StringValue: string(secretkey)}
	awsstruct.Fields[AwsRegionKey].Kind = &types.Value_StringValue{StringValue: spec.Region}
	awsstruct.Fields[AwsHostKey].Kind = &types.Value_StringValue{StringValue: spec.GetLambdaHostname()}
	return nil
}

func (a *AwsPlugin) IsMyUpstream(upstream *v1.Upstream) bool {
	if upstream == nil {
		// TODO: log warning
		return false
	}
	return upstream.Type == AwsPluginType
}

func (a *AwsPlugin) GetFunctionSpec(in *v1.Function) (*types.Struct, error) {
	spec, err := function.FromMap(in.Spec)
	if err != nil {
		return nil, err
	}

	funcstruct := &types.Struct{Fields: make(map[string]*types.Value)}
	funcstruct.Fields[AwsFunctionNameKey].Kind = &types.Value_StringValue{StringValue: string(spec.FunctionName)}
	funcstruct.Fields[AwsFunctionQualifierKey].Kind = &types.Value_StringValue{StringValue: string(spec.Qualifier)}

	return funcstruct, nil
}
