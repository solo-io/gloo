package plugin

import (
	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	apiroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/translator"
)

type Plugin interface {
	translator.UserResource

	EnvoyFilters(fi *PluginInputs) []FilterWrapper

	UpdateEnvoyRoute(fi *PluginInputs, in *v1.Route, out *apiroute.Route) error

	UpdateEnvoyCluster(fi *PluginInputs, in *v1.Upstream, out *api.Cluster) error

	UpdateFunctionToEnvoyCluster(fi *PluginInputs, in *v1.Upstream, infunc *v1.Function, out *api.Cluster) error
}

type FunctionalPlugin interface {
	IsMyUpstream(upstream *v1.Upstream) bool
	GetFunctionSpec(in *v1.Function) (*types.Struct, error)
}
