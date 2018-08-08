package plugins

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoylistener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

type Plugin interface {
	Init() error
}

type Params struct {
	Snapshot *v1.Snapshot
}

/*
	Upstream Plugins
 */

type UpstreamPlugin interface {
	Plugin
	ProcessUpstream(params Params, in *v1.Upstream, out *envoyapi.Cluster) error
}

/*
	Routing Plugins
 */

type RoutePlugin interface {
	Plugin
	ProcessRoute(params Params, in *v1.Route, out *envoyroute.Route) error
}

/*
	Listener Plugins
 */

type ListenerPlugin interface {
	Plugin
	ProcessListener(params Params, in *v1.Listener, out *envoyapi.Listener) error
}

type ListenerFilterPlugin interface {
	Plugin
	ProcessListenerFilter(params Params, in *v1.Listener) ([]StagedListenerFilter, error)
}

type StagedListenerFilter struct {
	ListenerFilter envoylistener.Filter
	Stage          FilterStage
}

type HttpFilterPlugin interface {
	Plugin
	HttpFilters(params Params, listener *v1.HttpListener) ([]StagedHttpFilter, error)
}

type StagedHttpFilter struct {
	HttpFilter *envoyhttp.HttpFilter
	Stage      FilterStage
}

type FilterStage int

const (
	PreInAuth  FilterStage = iota
	InAuth
	PostInAuth
	PreOutAuth
	OutAuth
)

/*
	Non-translator plugins
	TODO(ilackarms): consider combining eds plugin and uds
*/

type EdsPlugin interface {
	Plugin
	RunEds(opts clients.WatchOpts) error
	SubscribeUpstream(upstream *v1.Upstream) (<-chan []*v1.Endpoint, error)
}

type UdsPlugin interface {

}
