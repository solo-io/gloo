package plugin

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/module"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/hashicorp/go-multierror"
)

type ConfigStatus struct {
	Cfg v1.ConfigObject
	Err *multierror.Error
}

func NewConfigMultiError(Cfg v1.ConfigObject, err *multierror.Error) ConfigStatus {
	return ConfigStatus{Cfg: Cfg, Err: err}
}

func NewConfigError(Cfg v1.ConfigObject, err error) ConfigStatus {
	return ConfigStatus{Cfg: Cfg, Err: multierror.Append(nil, err)}
}

func NewConfigOk(Cfg v1.ConfigObject) ConfigStatus {
	return ConfigStatus{Cfg: Cfg, Err: nil}
}

type UserResource interface {
	GetDependencies(cfg *v1.Config) DependenciesDescription
}

type Plugin interface {
	UserResource

	EnvoyFilters(fi *PluginInputs) []FilterWrapper

	UpdateEnvoyRoute(fi *PluginInputs, in *v1.Route, out *api.Route) error

	UpdateEnvoyCluster(fi *PluginInputs, in *v1.Upstream, out *api.Cluster) error

	UpdateFunctionToEnvoyCluster(fi *PluginInputs, in *v1.Upstream, infunc *v1.Function, out *api.Cluster) error
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

	Translate(cfg *v1.Config, secretMap module.SecretMap, endpoints module.EndpointGroups) (*envoycache.Snapshot, []ConfigStatus)
}
