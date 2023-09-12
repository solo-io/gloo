package local_ratelimit

import (
	"errors"
	"fmt"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	gloo_local_ratelimit "github.com/solo-io/gloo/projects/gloo/pkg/plugins/local_ratelimit"
)

var (
	_ plugins.Plugin            = new(plugin)
	_ plugins.VirtualHostPlugin = new(plugin)
	_ plugins.RoutePlugin       = new(plugin)
)

const (
	HTTPFilterStatPrefix  = gloo_local_ratelimit.HTTPFilterStatPrefix
	HTTPFilterName        = gloo_local_ratelimit.HTTPFilterName
	CustomStageBeforeAuth = gloo_local_ratelimit.CustomStageBeforeAuth
)

var (
	ExtensionName = fmt.Sprintf("%s_ee", gloo_local_ratelimit.ExtensionName)

	ErrFilterDefinedInOSS     = errors.New("configuration for local_ratelimit already defined in options.ratelimit. Ignoring the config defined in options.ratelimitEarly")
	ErrFilterDefinedInRegular = errors.New("configuration for local_ratelimit can only be configured in options.ratelimitEarly or options.ratelimit. Ignoring the config defined in options.ratelimitRegular")
)

type plugin struct {
	filterRequiredForListener map[*v1.HttpListener]struct{}
}

// The OSS plugin adds the typed config on the vHost and route when the config is specified in options.ratelimit
// It also adds the http filter, so adding it is not needed in the enterprise plugin.
// The enterprise plugin adds the typed config on the vHost and route when specified in options.ratelimitEarly
// and errors if it is already configured in OSS.
// It also errors if specified in options.ratelimitRegular as this plugin is intended to be used pre-auth
func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
	p.filterRequiredForListener = make(map[*v1.HttpListener]struct{})
}

func (p *plugin) ProcessVirtualHost(
	params plugins.VirtualHostParams,
	in *v1.VirtualHost,
	out *envoy_config_route_v3.VirtualHost,
) error {
	settings := params.HttpListener.GetOptions().GetHttpLocalRatelimit()

	if beforeAuthLimits := in.GetOptions().GetRatelimitEarly().GetLocalRatelimit(); beforeAuthLimits != nil {
		err := gloo_local_ratelimit.ConfigureVirtualHostFilter(settings, beforeAuthLimits, CustomStageBeforeAuth, out)
		// It could error out if the plugin was already configured in OSS or any other reason. So return the appropriate error
		if err != nil {
			if err != gloo_local_ratelimit.ErrConfigurationExists {
				return err
			}
			return ErrFilterDefinedInOSS
		}
		p.filterRequiredForListener[params.HttpListener] = struct{}{}
	}
	// Error if the plugin is configured in options.ratelimitRegular as this is intended to executed be pre-auth
	if afterAuthLimits := in.GetOptions().GetRatelimitRegular().GetLocalRatelimit(); afterAuthLimits != nil {
		return ErrFilterDefinedInRegular
	}
	return nil
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	settings := params.HttpListener.GetOptions().GetHttpLocalRatelimit()

	if beforeAuthLimits := in.GetOptions().GetRatelimitEarly().GetLocalRatelimit(); beforeAuthLimits != nil {
		err := gloo_local_ratelimit.ConfigureRouteFilter(settings, beforeAuthLimits, CustomStageBeforeAuth, out)
		// It could error out if the plugin was already configured in OSS or any other reason. So return the appropriate error
		if err != nil {
			if err != gloo_local_ratelimit.ErrConfigurationExists {
				return err
			}
			return ErrFilterDefinedInOSS
		}
		p.filterRequiredForListener[params.HttpListener] = struct{}{}
	}
	// Error if the plugin is configured in options.ratelimitRegular as this is intended to executed be pre-auth
	if afterAuthLimits := in.GetOptions().GetRatelimitRegular().GetLocalRatelimit(); afterAuthLimits != nil {
		return ErrFilterDefinedInRegular
	}
	return nil
}
