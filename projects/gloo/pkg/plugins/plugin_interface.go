package plugins

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type Plugin interface {
	Init() error
}

type PluginParams struct {
	Snapshot *v1.Snapshot
}

type UpstreamPlugin interface {
	Plugin
	ProcessUpstream(params PluginParams, in *v1.Upstream, out *envoyapi.Cluster) error
}

type RoutePlugin interface {
	Plugin
	ProcessRoute(params PluginParams, in *v1.Route, out *envoyroute.Route) error
}
