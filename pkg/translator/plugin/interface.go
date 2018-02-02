package plugin

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/module"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/hashicorp/go-multierror"
)

type ConfigError struct {
	Cfg v1.ConfigObject
	Err *multierror.Error
}

func NewConfigError(Cfg v1.ConfigObject, Err *multierror.Error) ConfigError {
	return ConfigError{Cfg: Cfg, Err: Err}
}

type Plugin interface {
	GetDependencies(cfg *v1.Config) DependenciesDescription
	Validate(fi *PluginInputs) []ConfigError

	EnvoyFilters(fi *PluginInputs) []FilterWrapper
	UpdateEnvoyRoute(fi *PluginInputs, in *v1.Route, out *api.Route)
	UpdateEnvoyCluster(fi *PluginInputs, in *v1.Upstream, out *api.Cluster)
}

type DependenciesDescription interface {
	SecretRefs() []string
}

type Dependencies interface {
	Secrets() module.SecretMap
}

type NameTranslator interface {
	UpstreamToClusterName(string) string
}
