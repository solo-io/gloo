package prefixrewrite

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

type Plugin struct{}

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *Plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	if in.RoutePlugins == nil {
		return nil
	}
	if in.RoutePlugins.PrefixRewrite == nil {
		return nil
	}
	routeAction, ok := out.Action.(*envoyroute.Route_Route)
	if !ok {
		return errors.Errorf("prefix rewrite is only available for Route Actions")
	}
	if routeAction.Route == nil {
		return errors.Errorf("internal error: route %v specified a prefix, but output Envoy object "+
			"had nil route", in.Action)
	}
	routeAction.Route.PrefixRewrite = in.RoutePlugins.PrefixRewrite.PrefixRewrite
	return nil
}
