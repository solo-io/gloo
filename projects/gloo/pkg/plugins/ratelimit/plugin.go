package ratelimit

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	rlplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
)

var (
	_ plugins.VirtualHostPlugin = new(plugin)
	_ plugins.RoutePlugin       = new(plugin)
	_ plugins.HttpFilterPlugin  = new(plugin)
)

const (
	IngressDomain   = "ingress"
	ConfigCrdDomain = "crd"
	SetActionDomain = rlplugin.CustomDomain

	IngressRateLimitStage             = uint32(0)                      // 0
	SetActionRateLimitStage           = rlplugin.CustomStage           // 1
	CrdRateLimitStage                 = uint32(2)                      // 2
	SetActionRateLimitStageBeforeAuth = rlplugin.CustomStageBeforeAuth // 3
	CrdRateLimitStageBeforeAuth       = uint32(4)                      // 4
)

var (
	rateLimitStageToDomain = map[uint32]string{
		IngressRateLimitStage: IngressDomain,

		CrdRateLimitStage:           ConfigCrdDomain,
		CrdRateLimitStageBeforeAuth: ConfigCrdDomain,

		SetActionRateLimitStage:           SetActionDomain,
		SetActionRateLimitStageBeforeAuth: SetActionDomain,
	}

	rateLimitStageToFilterStage = map[uint32]plugins.FilterStage{
		IngressRateLimitStage: rlplugin.DefaultFilterStage,

		CrdRateLimitStage:           rlplugin.DefaultFilterStage,
		CrdRateLimitStageBeforeAuth: rlplugin.BeforeAuthStage,

		SetActionRateLimitStage:           rlplugin.DefaultFilterStage,
		SetActionRateLimitStageBeforeAuth: rlplugin.BeforeAuthStage,
	}
)

type plugin struct {
	serverSettings *ratelimit.Settings

	filterNeeded     bool // is set to indicate if the ratelimit configs exist
	configuredStages map[uint32]struct{}
	stagedTranslator StagedTranslator
}

func NewPlugin() *plugin {
	return NewPluginWithTranslators(
		translation.NewBasicRateLimitTranslator(),
		shims.NewGlobalRateLimitTranslator(),
		shims.NewRateLimitConfigTranslator(),
	)
}

func NewPluginWithTranslators(
	basic translation.BasicRateLimitTranslator,
	global shims.GlobalRateLimitTranslator,
	crd shims.RateLimitConfigTranslator,
) *plugin {
	return &plugin{
		stagedTranslator: getStagedTranslatorForRateLimitPlugin(basic, global, crd),
	}
}

func (p *plugin) Name() string {
	return fmt.Sprintf("%s_ee", rlplugin.ExtensionName)
}

func (p *plugin) Init(params plugins.InitParams) error {
	if rlServer := params.Settings.GetRatelimitServer(); rlServer != nil {
		p.serverSettings = rlServer
	}

	p.filterNeeded = !params.Settings.GetGloo().GetRemoveUnusedFilters().GetValue()
	p.configuredStages = make(map[uint32]struct{})

	return p.stagedTranslator.Init(params)
}

func (p *plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoy_config_route_v3.VirtualHost) error {
	var limits []*envoy_config_route_v3.RateLimit

	stagedRateLimits, err := p.stagedTranslator.GetVirtualHostRateLimitsByStage(params, in)
	for stage, rateLimits := range stagedRateLimits {
		if len(rateLimits) > 0 {
			// Mark this stage as having user-defined configuration
			p.configuredStages[stage] = struct{}{}
			limits = append(limits, rateLimits...)
		}
	}

	if len(limits) > 0 {
		out.RateLimits = append(out.RateLimits, limits...)
	}
	return err
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	routeAction := in.GetRouteAction()
	if routeAction == nil {
		// Only route actions can have rate limits
		return nil
	}

	outRouteAction := out.GetRoute()
	if outRouteAction == nil {
		return RouteTypeMismatchErr // should never happen
	}

	var limits []*envoy_config_route_v3.RateLimit

	stagedRateLimits, err := p.stagedTranslator.GetRouteRateLimitsByStage(params, in)
	for stage, rateLimits := range stagedRateLimits {
		if len(rateLimits) > 0 {
			// Mark this stage as having user-defined configuration
			p.configuredStages[stage] = struct{}{}
			limits = append(limits, rateLimits...)
		}
	}

	if len(limits) > 0 {
		outRouteAction.RateLimits = append(outRouteAction.RateLimits, limits...)
	}

	return err
}

func (p *plugin) getServerSettingsForListener(listener *v1.HttpListener) *ratelimit.Settings {
	if rlServer := listener.GetOptions().GetRatelimitServer(); rlServer != nil {
		return rlServer
	}

	return p.serverSettings
}

// HttpFilters returns Rate Limit Http Filters for the Gloo Enterprise API
// There are 3 type of rate limit filters that can be configured:
//  - Ingress
//  - RateLimitConfig
//  - SetAction
//
// There are 2 filter stages that these filters can be placed in:
//  - BeforeExtAuth
//  - AfterExtAuth
//
// To guarantee isolation between the configuration types, we generate separate
// filters for each combination of {config type, filter stage}
func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	serverSettings := p.getServerSettingsForListener(listener)
	upstreamRef := serverSettings.GetRatelimitServerRef()
	if upstreamRef == nil {
		return nil, nil
	}

	// Make sure the server exists
	_, err := params.Snapshot.Upstreams.Find(upstreamRef.GetNamespace(), upstreamRef.GetName())
	if err != nil {
		return nil, rlplugin.ServerNotFound(upstreamRef)
	}

	// The open source plugin creates a single http filter. Ensure that we do not
	// create a duplicate http filter referencing the same rate limit stage
	alreadyConfiguredRateLimitStage := getRateLimitStageConfiguredByOpenSourcePlugin(serverSettings)

	var stagedFilters []plugins.StagedHttpFilter
	for rateLimitStage, filterDomain := range rateLimitStageToDomain {
		if !p.filterNeeded {
			// If the filter is not required, check if there is configuration for the particular stage
			if _, ok := p.configuredStages[rateLimitStage]; !ok {
				// There is no configuration for the stage, skip it
				continue
			}
		}

		if rateLimitStage == alreadyConfiguredRateLimitStage {
			// If the stage matches the stage of the filter configured by the open source plugin, skip it
			continue
		}

		filterStage := rateLimitStageToFilterStage[rateLimitStage]
		rateLimitFilter := rlplugin.GenerateEnvoyConfigForFilterWith(
			upstreamRef,
			filterDomain,
			rateLimitStage,
			serverSettings.GetRequestTimeout(),
			serverSettings.GetDenyOnFail())

		stagedFilter, err := plugins.NewStagedFilterWithConfig(wellknown.HTTPRateLimit, rateLimitFilter, filterStage)
		if err != nil {
			return nil, err
		}
		stagedFilters = append(stagedFilters, stagedFilter)
	}

	return stagedFilters, nil
}

func getRateLimitStageConfiguredByOpenSourcePlugin(serverSettings *ratelimit.Settings) uint32 {
	return rlplugin.GetRateLimitStageForServerSettings(serverSettings)
}
