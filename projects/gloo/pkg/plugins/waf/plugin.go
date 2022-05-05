package waf

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/waf"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
)

var (
	_ plugins.Plugin            = new(plugin)
	_ plugins.VirtualHostPlugin = new(plugin)
	_ plugins.RoutePlugin       = new(plugin)
	_ plugins.HttpFilterPlugin  = new(plugin)
)

const (
	ExtensionName = "waf"
	FilterName    = "io.solo.filters.http.modsecurity"
)

type plugin struct {
	removeUnused              bool
	filterRequiredForListener map[*v1.HttpListener]struct{}
}

var (
	// waf should happen before any code is run
	filterStage = plugins.DuringStage(plugins.WafStage)
)

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) error {
	p.removeUnused = params.Settings.GetGloo().GetRemoveUnusedFilters().GetValue()
	p.filterRequiredForListener = make(map[*v1.HttpListener]struct{})
	return nil
}

// Process virtual host plugin
func (p *plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoy_config_route_v3.VirtualHost) error {
	wafConfig := in.Options.GetWaf()
	if wafConfig == nil {
		// no config found, nothing to do here
		return nil
	}

	// should never be nil
	p.filterRequiredForListener[params.HttpListener] = struct{}{}

	perVhostCfg := &ModSecurityPerRoute{
		Disabled:                  wafConfig.GetDisabled(),
		AuditLogging:              wafConfig.GetAuditLogging(),
		CustomInterventionMessage: wafConfig.GetCustomInterventionMessage(),
		RequestHeadersOnly:        wafConfig.GetRequestHeadersOnly(),
		ResponseHeadersOnly:       wafConfig.GetResponseHeadersOnly(),
	}

	perVhostCfg.RuleSets = wafConfig.GetRuleSets()
	if coreRuleSet := getCoreRuleSet(wafConfig.GetCoreRuleSet()); coreRuleSet != nil {
		perVhostCfg.RuleSets = append(perVhostCfg.RuleSets, coreRuleSet...)
	}

	pluginutils.SetVhostPerFilterConfig(out, FilterName, perVhostCfg)

	return nil
}

// Process route plugin
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	wafConfig := in.GetOptions().GetWaf()
	if wafConfig == nil {
		// no config found, nothing to do here
		return nil
	}

	p.filterRequiredForListener[params.HttpListener] = struct{}{}

	perRouteCfg := &ModSecurityPerRoute{
		Disabled:                  wafConfig.GetDisabled(),
		AuditLogging:              wafConfig.GetAuditLogging(),
		CustomInterventionMessage: wafConfig.GetCustomInterventionMessage(),
		RequestHeadersOnly:        wafConfig.GetRequestHeadersOnly(),
		ResponseHeadersOnly:       wafConfig.GetResponseHeadersOnly(),
	}

	perRouteCfg.RuleSets = wafConfig.GetRuleSets()
	if coreRuleSet := getCoreRuleSet(wafConfig.GetCoreRuleSet()); coreRuleSet != nil {
		perRouteCfg.RuleSets = append(perRouteCfg.RuleSets, coreRuleSet...)
	}

	pluginutils.SetRoutePerFilterConfig(out, FilterName, perRouteCfg)
	return nil
}

// Http Filter to return the waf filter
func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	// If the list does not already have the listener than it is necessary to check for nil

	wafSettings := listener.GetOptions().GetWaf()

	_, ok := p.filterRequiredForListener[listener]
	if !ok && p.removeUnused && wafSettings == nil {
		return filters, nil
	}

	var settings waf.Settings
	if wafSettings != nil {
		settings = *wafSettings
	}

	modSecurityConfig := &ModSecurity{}

	if settings.GetCoreRuleSet() == nil && settings.GetRuleSets() == nil {
		modSecurityConfig.Disabled = true
	} else {
		modSecurityConfig.RuleSets = settings.GetRuleSets()
		modSecurityConfig.AuditLogging = settings.GetAuditLogging()
		modSecurityConfig.Disabled = settings.GetDisabled()
		modSecurityConfig.CustomInterventionMessage = settings.GetCustomInterventionMessage()
		modSecurityConfig.RequestHeadersOnly = settings.GetRequestHeadersOnly()
		modSecurityConfig.ResponseHeadersOnly = settings.GetResponseHeadersOnly()

		if coreRuleSet := getCoreRuleSet(settings.GetCoreRuleSet()); coreRuleSet != nil {
			modSecurityConfig.RuleSets = append(modSecurityConfig.RuleSets, coreRuleSet...)
		}
	}

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, modSecurityConfig, filterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

func getCoreRuleSet(crs *waf.CoreRuleSet) []*RuleSet {
	if crs == nil {
		return nil
	}
	coreRuleSet := &RuleSet{
		Directory: crsPathPrefix,
	}
	coreRuleSetSettings := &RuleSet{}
	switch additionalSettings := crs.GetCustomSettingsType().(type) {
	case *waf.CoreRuleSet_CustomSettingsString:
		coreRuleSetSettings.RuleStr = additionalSettings.CustomSettingsString
	case *waf.CoreRuleSet_CustomSettingsFile:
		coreRuleSetSettings.Files = append([]string{additionalSettings.CustomSettingsFile}, coreRuleSet.Files...)
	}
	return []*RuleSet{coreRuleSetSettings, coreRuleSet}
}
