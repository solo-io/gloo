package istio_integration

import (
	"context"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/go-utils/contextutils"
)

var (
	_ plugins.RoutePlugin = new(plugin)
)

const (
	ExtensionName = "istio_integration"
)

// Handles transformations required to integrate with Istio
type plugin struct {
	appendXForwardedHost bool
}

// Deprecated
func NewPlugin(ctx context.Context) *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
	if xfh := params.Settings.GetGloo().GetIstioOptions().GetAppendXForwardedHost(); xfh != nil {
		p.appendXForwardedHost = xfh.GetValue()
	} else {
		p.appendXForwardedHost = false
	}
}

// Deprecated
// When istio integration is enabled, we need to access k8s services using a host that istio will recognize (servicename.namespace)
// We do this by adding a hostRewrite for kube destinations kube upstreams. In case the upstream also wants the original host,
// we also set x-forwarded-host
// We ignore other destinations and routes that already have a rewrite applied and return an error if we can't look up an Upstream.
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	if p.appendXForwardedHost {
		contextutils.LoggerFrom(params.Ctx).Warn("append_x_forwarded_host is deprecated")
	}

	return nil
}
