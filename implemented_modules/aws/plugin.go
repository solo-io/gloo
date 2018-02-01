package aws

import (
	"errors"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/glue/implemented_modules/aws/pkg/function"
	"github.com/solo-io/glue/implemented_modules/aws/pkg/upstream"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator/plugin"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/api/filter/network"
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

const AwsFilterName = "io.solo.aws"
const AwsFilterStage = plugin.OutAuth

type AwsPlugin struct {
	helper plugin.FuncitonalFilterHelper
}

func New() plugin.Plugin {
	return &AwsPlugin{}
}

func (a *AwsPlugin) GetDependencies(cfg *v1.Config) plugin.DependenciesDescription {
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

func (a *AwsPlugin) Validate(pi *plugin.PluginInputs) []plugin.ConfigError {

	var cfgErrors []plugin.ConfigError

	for _, u := range pi.State.Config.Upstreams {
		if !a.IsMyUpstream(&u) {
			continue
		}
		var errors *multierror.Error
		for err := range a.validateUpstream(pi, &u) {
			errors = multierror.Append(errors, err)
		}
		cfgErrors = append(cfgErrors, plugin.NewConfigError(&u, errors))
	}
	return cfgErrors
}

func (a *AwsPlugin) validateUpstream(pi *plugin.PluginInputs, u *v1.Upstream) <-chan error {
	errorchan := make(chan error)
	go func() {
		defer close(errorchan)
		spec, err := upstream.FromMap(u.Spec)
		if err != nil {
			errorchan <- err
			return
		}

		// verify that the secrets look right
		secret := pi.State.Dependencies.Secrets()[spec.Secret]
		if secret == nil {
			errorchan <- err
			return
		}

		if accesskey, ok := secret[AwsSecretAccessKey]; !ok {
			errorchan <- errors.New("secret doesn't contain " + AwsSecretAccessKey)
		} else if !utf8.Valid(accesskey) {
			errorchan <- errors.New("access key is not a valid string ")
		}

		if secretkey, ok := secret[AwsSecretSecretKey]; !ok {
			errorchan <- errors.New("secret doesn't contain " + AwsSecretSecretKey)
		} else if !utf8.Valid(secretkey) {
			errorchan <- errors.New("secret key is not a valid string ")
		}

		// TODO: validate all the functions of the upstream as well

	}()
	return errorchan
}

func (a *AwsPlugin) EnvoyFilters(pi *plugin.PluginInputs) []plugin.FilterWrapper {
	filter := plugin.FilterWrapper{
		Filter: network.HttpFilter{
			Name: AwsFilterName,
		},
		Stage: AwsFilterStage,
	}
	return []plugin.FilterWrapper{filter}
}

func (a *AwsPlugin) UpdateEnvoyRoute(pi *plugin.PluginInputs, in *v1.Route, out *api.Route) {
	a.helper.UpdateRoute(pi, a, in, out)
}

func (a *AwsPlugin) UpdateEnvoyCluster(pi *plugin.PluginInputs, in *v1.Upstream, out *api.Cluster) {

	if !a.IsMyUpstream(in) {
		return
	}

	spec, err := upstream.FromMap(in.Spec)
	if err != nil {
		// TODO: log error
		return
	}

	secret := pi.State.Dependencies.Secrets()[spec.Secret]

	if out.Metadata == nil {
		out.Metadata = &api.Metadata{
			FilterMetadata: make(map[string]*types.Struct),
		}
	}
	awsstruct := &types.Struct{Fields: make(map[string]*types.Value)}
	out.Metadata.FilterMetadata[AwsFilterName] = awsstruct

	awsstruct.Fields[AwsSecretAccessKey].Kind = &types.Value_StringValue{StringValue: string(secret[AwsSecretAccessKey])}
	awsstruct.Fields[AwsSecretSecretKey].Kind = &types.Value_StringValue{StringValue: string(secret[AwsSecretSecretKey])}
	awsstruct.Fields[AwsRegionKey].Kind = &types.Value_StringValue{StringValue: spec.Region}
	awsstruct.Fields[AwsHostKey].Kind = &types.Value_StringValue{StringValue: spec.GetLambdaHostname()}
}

func (a *AwsPlugin) IsMyUpstream(upstream *v1.Upstream) bool {
	if upstream == nil {
		// TODO: log warning
		return false
	}
	return upstream.Type == AwsPluginType
}

func (a *AwsPlugin) GetFunctionSpec(in *v1.Function) *types.Struct {
	spec, err := function.FromMap(in.Spec)

	funcstruct := &types.Struct{Fields: make(map[string]*types.Value)}
	funcstruct.Fields[AwsFunctionNameKey].Kind = &types.Value_StringValue{StringValue: string(spec.FunctionName)}
	// TODO: add qualifier.

	return funcstruct
}
