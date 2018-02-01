package aws

import (
	"errors"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/glue/implemented_modules/aws/pkg/upstream"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator/plugin"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/api/filter/network"
	"github.com/hashicorp/go-multierror"
)

const AwsPluginType = "aws"
const AwsSecretAccessKey = "access_key"
const AwsSecretSecretKey = "secret_key"

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

		if _, ok := secret[AwsSecretAccessKey]; !ok {
			errorchan <- errors.New("secret doesn't contain " + AwsSecretAccessKey)
		}
		if _, ok := secret[AwsSecretSecretKey]; !ok {
			errorchan <- errors.New("secret doesn't contain " + AwsSecretSecretKey)
		}
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

func (a *AwsPlugin) UpdateEnvoyRoute(pi *plugin.PluginInputs, routein *v1.Route, routeout *api.Route) {
	a.helper.UpdateRoute(pi, a, routein, routeout)
}

func (a *AwsPlugin) UpdateEnvoyCluster(pi *plugin.PluginInputs, routein *v1.Upstream, routeout *api.Cluster) {
	panic("panic")
}

func (a *AwsPlugin) IsMyUpstream(upstream *v1.Upstream) bool {
	if upstream == nil {
		// TODO: log warning
		return false
	}
	return upstream.Type == AwsPluginType
}

func (a *AwsPlugin) GetFunctionSpec(functioname string) *types.Struct {
	panic("panic")
}
