package extauth

import (
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"

	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/ext_authz/v2"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/extauth"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/go-utils/errors"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
)

const (
	ExtensionName      = "extauth"
	SanitizeFilterName = "io.solo.filters.http.sanitize"
	FilterName         = "envoy.ext_authz"
	DefaultAuthHeader  = "x-user-id"

	// when extauth is deployed into a sidecar in the envoy pod, an upstream should
	// be created that points to that sidecar and has this name. A regression test
	// attempts to talk to that upstream
	SidecarUpstreamName = "extauth-sidecar"
)

var (
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

	// Add sanitize filter if a user ID header is defined in the settings
	if p.userIdHeader != "" {
		sanitizeConf := &Sanitize{
			HeadersToRemove: []string{p.userIdHeader},
		}
		stagedFilter, err := plugins.NewStagedFilterWithConfig(SanitizeFilterName, sanitizeConf, sanitizeFilterStage)
		if err != nil {
			return nil, err
		}
		filters = append(filters, stagedFilter)
	}

	return filters, nil
}

// This function generates the ext_authz PerFilterConfig for this virtual host.
// If the virtual host does not explicitly define an extauth configuration, we disable the ext_authz filter.
// Since the ext_authz filter is always enabled on the listener, we need this to disable authentication by default on
// a virtual host and its child resources (routes, weighted destinations). Extauth should be opt-in.
func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.Snapshot.Upstreams) {
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

	return pluginutils.SetVhostPerFilterConfig(out, FilterName, config)
}

// This function generates the ext_authz PerFilterConfig for this route:
// - if the route defines auth configuration, set the filter correspondingly;
// - if auth is explicitly disabled, disable the filter (will apply by default also to WeightedDestinations);
// - if not auth config is defined, do nothing (will inherit config from parent virtual host).
func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.Snapshot.Upstreams) {
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

	return pluginutils.SetRoutePerFilterConfig(out, FilterName, config)
}

// This function generates the ext_authz PerFilterConfig for this weightedDestination:
// - if the weightedDestination defines auth configuration, set the filter correspondingly;
// - if auth is explicitly disabled, disable the filter;
// - if not auth config is defined, do nothing (will inherit config from parent virtual host and/or route).
func (p *Plugin) ProcessWeightedDestination(params plugins.RouteParams, in *v1.WeightedDestination, out *envoyroute.WeightedCluster_ClusterWeight) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.Snapshot.Upstreams) {
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

	return pluginutils.SetWeightedClusterPerFilterConfig(out, FilterName, config)
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

func markVirtualHostNoAuth(out *envoyroute.VirtualHost) error {
	return pluginutils.SetVhostPerFilterConfig(out, FilterName, getNoAuthConfig())
}

func markWeightedClusterNoAuth(out *envoyroute.WeightedCluster_ClusterWeight) error {
	return pluginutils.SetWeightedClusterPerFilterConfig(out, FilterName, getNoAuthConfig())
}

func markRouteNoAuth(out *envoyroute.Route) error {
	return pluginutils.SetRoutePerFilterConfig(out, FilterName, getNoAuthConfig())
}

func getNoAuthConfig() *envoyauth.ExtAuthzPerRoute {
	return &envoyauth.ExtAuthzPerRoute{
		Override: &envoyauth.ExtAuthzPerRoute_Disabled{
			Disabled: true,
		},
	}
}

func (p *Plugin) isExtAuthzFilterConfigured(upstreams v1.UpstreamList) bool {

	// Call the same function called by HttpFilters to verify whether the filter was created
	filters, err := extauth.BuildHttpFilters(p.extAuthSettings, upstreams)
	if err != nil {
		// If it returned an error, the filter was not configured
		return false
	}

	// Check for a filter called "envoy.ext_authz"
	for _, filter := range filters {
		if filter.HttpFilter.GetName() == FilterName {
			return true
		}
	}

	return false
}
