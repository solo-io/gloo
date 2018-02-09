package plugin

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

type Dependencies struct {
	SecretRefs []string
}

type TranslatorPlugin interface {
	GetDependencies(cfg v1.Config) *Dependencies
}

type UpstreamPlugin interface {
	TranslatorPlugin
	ProcessUpstream(in v1.Upstream, secrets secretwatcher.SecretMap, out *envoyapi.Cluster) error
}

type FunctionPlugin interface {
	UpstreamPlugin
	// if the FunctionSpec does not belong to this plugin, return nil, nil
	// if the FunctionSpec belongs to this plugin but is not valid, return nil, err
	// if the FunctionSpec belogns to this plugin and is valid, return *Struct, nil
	ParseFunctionSpec(upstreamType v1.UpstreamType, in v1.FunctionSpec) (*types.Struct, error)
}

type RoutePlugin interface {
	TranslatorPlugin
	ProcessRoute(in v1.Route, out *envoyroute.Route) error
}
