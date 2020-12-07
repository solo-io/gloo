package extauth

import (
	"fmt"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	errors "github.com/rotisserie/eris"
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/extauth"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
)

const (
	ExtensionName      = "extauth"
	SanitizeFilterName = "io.solo.filters.http.sanitize"
	DefaultAuthHeader  = "x-user-id"

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

	sanitizeFilterStage = plugins.BeforeStage(plugins.AuthNStage)

	NoMatchesForGroupError = func(labelSelector map[string]string) error {
		return errors.Errorf("no matching apikey secrets for the provided label selector %v", labelSelector)
	}
)

const (
	SourceTypeVirtualHost         = "virtual_host"
	SourceTypeRoute               = "route"
	SourceTypeWeightedDestination = "weighted_destination"
)

type Plugin struct {
	userIdHeader    string
	extAuthSettings *extauthapi.Settings
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func BuildVirtualHostName(proxy *v1.Proxy, listener *v1.Listener, virtualHost *v1.VirtualHost) string {
	return fmt.Sprintf("%s-%s-%s", proxy.Metadata.Ref().Key(), listener.Name, virtualHost.Name)
}

func GetAuthHeader(e *extauthapi.Settings) string {
	if e != nil {
		if e.UserIdHeader != "" {
			return e.UserIdHeader
		}
	}
	return DefaultAuthHeader
}

func (p *Plugin) Init(params plugins.InitParams) error {
	p.userIdHeader = ""
	p.extAuthSettings = nil

	settings := params.Settings.GetExtauth()
	p.extAuthSettings = settings
	p.userIdHeader = GetAuthHeader(settings)
	return nil
}

// This function just needs to add the sanitize filter. If extauth has been configured in the settings,
// the ext_authz will already have been created by the extauth plugin in OS Gloo.
func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter

	userIdHeader := listener.GetOptions().GetExtauth().GetUserIdHeader()
	if userIdHeader == "" {
		userIdHeader = p.userIdHeader
	}

	// Add sanitize filter if a user ID header is defined in the settings
	if userIdHeader != "" {
		sanitizeConf := &Sanitize{
			HeadersToRemove: []string{userIdHeader},
		}
		stagedFilter, err := plugins.NewStagedFilterWithConfig(SanitizeFilterName, sanitizeConf, sanitizeFilterStage)
		if err != nil {
			return nil, err
		}
		filters = append(filters, stagedFilter)
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

	extAuthConfig := in.GetOptions().GetExtauth()

	// No config was defined or explicitly disabled, disable the filter for this virtual host
	if extAuthConfig == nil || extAuthConfig.GetDisable() {
		return markVirtualHostNoAuth(out)
	}

	// No auth config ref provided, must be using custom auth (which has already been configured by open-source plugin)
	if extAuthConfig.GetConfigRef() == nil {
		return nil
	}

	config, err := buildFilterConfig(
		SourceTypeVirtualHost,
		BuildVirtualHostName(params.Proxy, params.Listener, in),
		extAuthConfig.GetConfigRef().Key(),
	)
	if err != nil {
		return err
	}

	return pluginutils.SetVhostPerFilterConfig(out, wellknown.HTTPExternalAuthorization, config)
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

	extAuthConfig := in.GetOptions().GetExtauth()

	// No config was defined, just return
	if extAuthConfig == nil {
		return nil
	}

	// Explicitly disable the filter for this route
	if extAuthConfig.GetDisable() {
		return markRouteNoAuth(out)
	}

	// No auth config ref provided, must be using custom auth (which has already been configured by open-source plugin)
	if extAuthConfig.GetConfigRef() == nil {
		return nil
	}

	config, err := buildFilterConfig(SourceTypeRoute, "", extAuthConfig.GetConfigRef().Key())
	if err != nil {
		return err
	}

	return pluginutils.SetRoutePerFilterConfig(out, wellknown.HTTPExternalAuthorization, config)
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

	extAuthConfig := in.GetOptions().GetExtauth()

	// No config was defined, just return
	if extAuthConfig == nil {
		return nil
	}

	// Explicitly disable the filter for this route
	if extAuthConfig.GetDisable() {
		return markWeightedClusterNoAuth(out)
	}

	// No auth config ref provided, must be using custom auth (which has already been configured by open-source plugin)
	if extAuthConfig.GetConfigRef() == nil {
		return nil
	}

	config, err := buildFilterConfig(SourceTypeWeightedDestination, "", extAuthConfig.GetConfigRef().Key())
	if err != nil {
		return err
	}

	return pluginutils.SetWeightedClusterPerFilterConfig(out, wellknown.HTTPExternalAuthorization, config)
}

func buildFilterConfig(sourceType, sourceName, authConfigRef string) (*envoyauth.ExtAuthzPerRoute, error) {
	requestContext, err := extauthservice.NewRequestContext(authConfigRef, sourceType, sourceName)
	if err != nil {
		return nil, err
	}

	return &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_CheckSettings{
			CheckSettings: &envoyauth.CheckSettings{
				ContextExtensions: requestContext.ToContextExtensions(),
			},
		},
	}, nil
}

func markVirtualHostNoAuth(out *envoy_config_route_v3.VirtualHost) error {
	return pluginutils.SetVhostPerFilterConfig(out, wellknown.HTTPExternalAuthorization, getNoAuthConfig())
}

func markWeightedClusterNoAuth(out *envoy_config_route_v3.WeightedCluster_ClusterWeight) error {
	return pluginutils.SetWeightedClusterPerFilterConfig(out, wellknown.HTTPExternalAuthorization, getNoAuthConfig())
}

func markRouteNoAuth(out *envoy_config_route_v3.Route) error {
	return pluginutils.SetRoutePerFilterConfig(out, wellknown.HTTPExternalAuthorization, getNoAuthConfig())
}

func getNoAuthConfig() *envoyauth.ExtAuthzPerRoute {
	return &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}
}

func (p *Plugin) isExtAuthzFilterConfigured(listener *v1.HttpListener, upstreams v1.UpstreamList) bool {

	// Call the same function called by HttpFilters to verify whether the filter was created
	filters, err := extauth.BuildHttpFilters(p.extAuthSettings, listener, upstreams)
	if err != nil {
		// If it returned an error, the filter was not configured
		return false
	}

	// Check for a filter called "envoy.filters.http.ext_authz"
	for _, filter := range filters {
		if filter.HttpFilter.GetName() == wellknown.HTTPExternalAuthorization {
			return true
		}
	}

	return false
}
