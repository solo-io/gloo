package plugin

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gogo/protobuf/types"
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

type UserResource interface {
	GetDependencies(cfg *v1.Config) DependenciesDescription
	Validate(fi *PluginInputs) []ConfigError
}

type Plugin interface {
	UserResource

	EnvoyFilters(fi *PluginInputs) []FilterWrapper
	UpdateEnvoyRoute(fi *PluginInputs, in *v1.Route, out *api.Route)
	UpdateEnvoyCluster(fi *PluginInputs, in *v1.Upstream, out *api.Cluster)
}

type FunctionalPlugin interface {
	IsMyUpstream(upstream *v1.Upstream) bool
	GetFunctionSpec(in *v1.Function) (*types.Struct, error)
}

type DependenciesDescription interface {
	SecretRefs() []string
}

type Dependencies interface {
	Secrets() module.SecretMap
}

type NameTranslator interface {
	UpstreamToClusterName(string) string
	ToEnvoyVhostName(*v1.VirtualHost) string
}

type Translator interface {
	UserResource

	Translate(cfg *v1.Config, secretMap module.SecretMap, endpoints module.EndpointGroups) (*envoycache.Snapshot, error)
}
