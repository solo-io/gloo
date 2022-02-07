package extauth

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	errors "github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
)

const (
	DefaultAuthHeader = "x-user-id"

	// when extauth is deployed into a sidecar in the envoy pod, an upstream should
	// be created that points to that sidecar and has this name. A regression test
	// attempts to talk to that upstream
	SidecarUpstreamName = "extauth-sidecar"
)

var (
	_ plugins.Plugin                    = new(Plugin)
	_ plugins.VirtualHostPlugin         = new(Plugin)
	_ plugins.RoutePlugin               = new(Plugin)
	_ plugins.HttpFilterPlugin          = new(Plugin)
	_ plugins.WeightedDestinationPlugin = new(Plugin)
	_ plugins.Upgradable                = new(Plugin)

	sanitizeFilterStage = plugins.BeforeStage(plugins.AuthNStage)

	NoMatchesForGroupError = func(labelSelector map[string]string) error {
		return errors.Errorf("no matching apikey secrets for the provided label selector %v", labelSelector)
	}
)

type Plugin struct {
	userIdHeader         string
	namedExtAuthSettings map[string]*extauthapi.Settings

	extAuthzConfigGenerator extauth.ExtAuthzConfigGenerator
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func GetAuthHeader(e *extauthapi.Settings) string {
	if e != nil {
		if e.UserIdHeader != "" {
			return e.UserIdHeader
		}
	}
	return DefaultAuthHeader
}

func (p *Plugin) PluginName() string {
	return extauth.ExtensionName
}

func (p *Plugin) IsUpgrade() bool {
	return true
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.userIdHeader = ""

	settings := params.Settings.GetExtauth()
	p.userIdHeader = GetAuthHeader(settings)

	p.namedExtAuthSettings = params.Settings.GetNamedExtauth()
	p.extAuthzConfigGenerator = getEnterpriseConfigGenerator(settings, p.namedExtAuthSettings)
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter

	// Configure ext_authz http filters
	stagedFilters, err := extauth.BuildStagedHttpFilters(func() ([]*envoyauth.ExtAuthz, error) {
		return p.extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(listener, params.Snapshot.Upstreams)
	}, extauth.FilterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilters...)

	// ExtAuth relies on the sanitize filter to achieve some of its functionality
	// Add sanitize filter if a user ID header is defined in the settings
	// or if multiple ext_authz filters are configured
	userIdHeader := listener.GetOptions().GetExtauth().GetUserIdHeader()
	if userIdHeader == "" {
		userIdHeader = p.userIdHeader
	}

	includeCustomAuthServiceName := len(stagedFilters) >= 2
	if userIdHeader != "" || includeCustomAuthServiceName {
		// In the case where multiple ext_authz filters are configured, we want to be sure that at least
		// the default filter is enabled. To ensure this, we configure the sanitize filter
		sanitizeFilter, err := buildSanitizeFilter(userIdHeader, includeCustomAuthServiceName)
		if err != nil {
			return nil, err
		}
		filters = append(filters, sanitizeFilter)
	}

	return filters, nil
}

// This function generates the ext_authz TypedPerFilterConfig for this virtual host.
// If the virtual host does not explicitly define an extauth configuration, we disable the ext_authz filter.
// Since the ext_authz filter is always enabled on the listener, we need this to disable authentication by default on
// a virtual host and its child resources (routes, weighted destinations). Extauth should be opt-in.
func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoy_config_route_v3.VirtualHost) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.Listener.GetHttpListener(), params.Snapshot.Upstreams) {
		return nil
	}

	// Configure the sanitize filter in the case of multiple ext_authz filters
	if p.extAuthzConfigGenerator.IsMulti() {
		err := setVirtualHostCustomAuth(out, in.GetOptions().GetExtauth(), p.namedExtAuthSettings)
		if err != nil {
			return err
		}
	}

	// Configure the ext_authz filter per route config
	extAuthPerRouteConfig, err := p.extAuthzConfigGenerator.GenerateVirtualHostExtAuthzConfig(in, params)
	if err != nil {
		return err
	}
	if extAuthPerRouteConfig == nil {
		return nil
	}
	return pluginutils.SetVhostPerFilterConfig(out, wellknown.HTTPExternalAuthorization, extAuthPerRouteConfig)
}

// This function generates the ext_authz TypedPerFilterConfig for this route:
// - if the route defines auth configuration, set the filter correspondingly;
// - if auth is explicitly disabled, disable the filter (will apply by default also to WeightedDestinations);
// - if not auth config is defined, do nothing (will inherit config from parent virtual host).
func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.Listener.GetHttpListener(), params.Snapshot.Upstreams) {
		return nil
	}

	// Configure the sanitize filter in the case of multiple ext_authz filters
	if p.extAuthzConfigGenerator.IsMulti() {
		err := setRouteCustomAuth(out, in.GetOptions().GetExtauth(), p.namedExtAuthSettings)
		if err != nil {
			return err
		}
	}

	// Configure the ext_authz filter per route config
	extAuthPerRouteConfig, err := p.extAuthzConfigGenerator.GenerateRouteExtAuthzConfig(in)
	if err != nil {
		return err
	}
	if extAuthPerRouteConfig == nil {
		return nil
	}
	return pluginutils.SetRoutePerFilterConfig(out, wellknown.HTTPExternalAuthorization, extAuthPerRouteConfig)
}

// This function generates the ext_authz TypedPerFilterConfig for this weightedDestination:
// - if the weightedDestination defines auth configuration, set the filter correspondingly;
// - if auth is explicitly disabled, disable the filter;
// - if not auth config is defined, do nothing (will inherit config from parent virtual host and/or route).
func (p *Plugin) ProcessWeightedDestination(params plugins.RouteParams, in *v1.WeightedDestination, out *envoy_config_route_v3.WeightedCluster_ClusterWeight) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.Listener.GetHttpListener(), params.Snapshot.Upstreams) {
		return nil
	}

	// Configure the sanitize filter in the case of multiple ext_authz filters
	if p.extAuthzConfigGenerator.IsMulti() {
		err := setWeightedClusterCustomAuth(out, in.GetOptions().GetExtauth(), p.namedExtAuthSettings)
		if err != nil {
			return err
		}
	}

	// Configure the ext_authz filter per route config
	extAuthPerRouteConfig, err := p.extAuthzConfigGenerator.GenerateWeightedDestinationExtAuthzConfig(in)
	if err != nil {
		return err
	}
	if extAuthPerRouteConfig == nil {
		return nil
	}
	return pluginutils.SetWeightedClusterPerFilterConfig(out, wellknown.HTTPExternalAuthorization, extAuthPerRouteConfig)
}

func (p *Plugin) isExtAuthzFilterConfigured(listener *v1.HttpListener, upstreams v1.UpstreamList) bool {
	// Call the same function called by HttpFilters to verify whether the filter was created
	stagedFilters, err := extauth.BuildStagedHttpFilters(func() ([]*envoyauth.ExtAuthz, error) {
		return p.extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(listener, upstreams)
	}, extauth.FilterStage)

	if err != nil {
		// If it returned an error, the filter was not configured
		return false
	}

	// Check for a filter called "envoy.filters.http.ext_authz"
	return plugins.StagedFilterListContainsName(stagedFilters, wellknown.HTTPExternalAuthorization)
}
