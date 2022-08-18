package extauth

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	errors "github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	_ plugins.Plugin                    = new(plugin)
	_ plugins.VirtualHostPlugin         = new(plugin)
	_ plugins.RoutePlugin               = new(plugin)
	_ plugins.HttpFilterPlugin          = new(plugin)
	_ plugins.WeightedDestinationPlugin = new(plugin)
)

const (
	DefaultAuthHeader = "x-user-id"

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

type plugin struct {
	userIdHeader         string
	namedExtAuthSettings map[string]*extauthapi.Settings

	extAuthzConfigGenerator extauth.ExtAuthzConfigGenerator

	// configStore holds state that is shared across multiple functions within the plugin
	configStore *ConfigStore
}

func NewPlugin() *plugin {
	return &plugin{}
}

func GetAuthHeader(e *extauthapi.Settings) string {
	if e != nil {
		if e.UserIdHeader != "" {
			return e.UserIdHeader
		}
	}
	return DefaultAuthHeader
}

func (p *plugin) Name() string {
	return extauth.ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
	settings := params.Settings.GetExtauth()
	p.userIdHeader = GetAuthHeader(settings)
	p.namedExtAuthSettings = params.Settings.GetNamedExtauth()
	p.extAuthzConfigGenerator = getEnterpriseConfigGenerator(settings, p.namedExtAuthSettings)
	p.configStore = NewConfigStore()
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter

	// Read ExtAuth HttpFilters from the ConfigStore
	// During the processing of VirtualHosts and Routes, we generate the actual HttpFilter configuration
	// and persist it in the store. Our translation loop relies upon the expectation that
	// ProcessVirtualHost and ProcessRoute will be executed before HttpFilters
	stagedFilters, err := p.configStore.getStagedFiltersForHttpListener(listener)
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
func (p *plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoy_config_route_v3.VirtualHost) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.HttpListener, params.Snapshot.Upstreams) {
		return nil
	}

	// VirtualHosts may define requests transformations which write data to dynamic metadata
	// The ExtAuth HttpFilter only has access to namespaces it has explicitly been configured with
	// Therefore, when we process the VirtualHost we store the namespaces that request transformations may write data to
	p.configStore.appendMetadataNamespacesForHttpListener(params.HttpListener, getMetadataNamespacesFromVirtualHostTransformations(in))

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
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.HttpListener, params.Snapshot.Upstreams) {
		return nil
	}

	// Routes may define requests transformations which write data to dynamic metadata
	// The ExtAuth HttpFilter only has access to namespaces it has explicitly been configured with
	// Therefore, when we process the Route, we store the namespaces that request transformations may write data to
	p.configStore.appendMetadataNamespacesForHttpListener(params.HttpListener, getMetadataNamespacesFromRouteTransformations(in))

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
func (p *plugin) ProcessWeightedDestination(params plugins.RouteParams, in *v1.WeightedDestination, out *envoy_config_route_v3.WeightedCluster_ClusterWeight) error {

	// Ext_authz filter is not configured, do nothing
	if !p.isExtAuthzFilterConfigured(params.HttpListener, params.Snapshot.Upstreams) {
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

// isExtAuthzFilterConfigured returns true if ExtAuth Filters will be generated for a particular HttpListener
// Historically, we rebuilt the HttpFilters every time, but that is a costly operation
// Therefore, we memoize the filters and query the ConfigStore for the current state
func (p *plugin) isExtAuthzFilterConfigured(listener *v1.HttpListener, upstreams v1.UpstreamList) bool {
	if p.configStore.hasFiltersForHttpListener(listener) {
		// the store already has the filters built
		return true
	}

	httpFilters, err := p.extAuthzConfigGenerator.GenerateListenerExtAuthzConfig(listener, upstreams)

	p.configStore.setFiltersForHttpListener(listener, httpFilters, err)
	return len(httpFilters) > 0
}

func getMetadataNamespacesFromRouteTransformations(route *v1.Route) []string {
	var requestMatches []*transformation.RequestMatch

	if route.GetOptions().GetStagedTransformations().GetEarly().GetRequestTransforms() != nil {
		requestMatches = append(requestMatches, route.GetOptions().GetStagedTransformations().GetEarly().GetRequestTransforms()...)
	}
	if route.GetOptions().GetStagedTransformations().GetRegular().GetRequestTransforms() != nil {
		requestMatches = append(requestMatches, route.GetOptions().GetStagedTransformations().GetRegular().GetRequestTransforms()...)
	}

	return getDynamicMetadataNamespacesForRequestTransformations(requestMatches)
}

func getMetadataNamespacesFromVirtualHostTransformations(virtualHost *v1.VirtualHost) []string {
	var requestMatches []*transformation.RequestMatch

	if virtualHost.GetOptions().GetStagedTransformations().GetEarly().GetRequestTransforms() != nil {
		requestMatches = append(requestMatches, virtualHost.GetOptions().GetStagedTransformations().GetEarly().GetRequestTransforms()...)
	}
	if virtualHost.GetOptions().GetStagedTransformations().GetRegular().GetRequestTransforms() != nil {
		requestMatches = append(requestMatches, virtualHost.GetOptions().GetStagedTransformations().GetRegular().GetRequestTransforms()...)
	}

	return getDynamicMetadataNamespacesForRequestTransformations(requestMatches)
}

// getDynamicMetadataNamespacesForRequestTransformations returns the set of dynamic metadata namespaces
// that request transformations may persist data.
// NOTE: This set of namespaces will be sorted before being applied to the HttpFilter,
// which is why we are ok returning an unsorted list here.
func getDynamicMetadataNamespacesForRequestTransformations(requestMatches []*transformation.RequestMatch) []string {
	namespaces := sets.NewString()

	for _, requestMatch := range requestMatches {
		for _, metaData := range requestMatch.GetRequestTransformation().GetTransformationTemplate().GetDynamicMetadataValues() {
			namespace := metaData.GetMetadataNamespace()
			if namespace != "" {
				namespaces.Insert(namespace)
			}
		}
	}

	return namespaces.UnsortedList()
}
