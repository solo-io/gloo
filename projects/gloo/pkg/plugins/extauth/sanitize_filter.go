package extauth

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/extauth"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
)

const SanitizeFilterName = "io.solo.filters.http.sanitize"

func buildSanitizeFilter(userIdHeader string, includeCustomAuthServiceName bool) (plugins.StagedHttpFilter, error) {
	sanitizeConf := &Sanitize{}

	if userIdHeader != "" {
		sanitizeConf.HeadersToRemove = []string{userIdHeader}
	}

	if includeCustomAuthServiceName {
		sanitizeConf.CustomAuthServerName = DefaultCustomAuthServiceName
	}

	return plugins.NewStagedFilterWithConfig(SanitizeFilterName, sanitizeConf, sanitizeFilterStage)
}

func setVirtualHostCustomAuth(out *envoy_config_route_v3.VirtualHost, customAuth *extauthapi.CustomAuth, availableCustomAuth map[string]*extauthapi.Settings) error {
	customAuthConfig := buildSanitizePerRouteConfig(customAuth, availableCustomAuth)
	if customAuthConfig == nil {
		return nil
	}
	return pluginutils.SetVhostPerFilterConfig(out, SanitizeFilterName, customAuthConfig)
}

func setWeightedClusterCustomAuth(out *envoy_config_route_v3.WeightedCluster_ClusterWeight, customAuth *extauthapi.CustomAuth, availableCustomAuth map[string]*extauthapi.Settings) error {
	customAuthConfig := buildSanitizePerRouteConfig(customAuth, availableCustomAuth)
	if customAuthConfig == nil {
		return nil
	}
	return pluginutils.SetWeightedClusterPerFilterConfig(out, SanitizeFilterName, customAuthConfig)
}

func setRouteCustomAuth(out *envoy_config_route_v3.Route, customAuth *extauthapi.CustomAuth, availableCustomAuth map[string]*extauthapi.Settings) error {
	customAuthConfig := buildSanitizePerRouteConfig(customAuth, availableCustomAuth)
	if customAuthConfig == nil {
		return nil
	}
	return pluginutils.SetRoutePerFilterConfig(out, SanitizeFilterName, customAuthConfig)
}

func buildSanitizePerRouteConfig(customAuth *extauthapi.CustomAuth, availableCustomAuth map[string]*extauthapi.Settings) *SanitizePerRoute {
	if customAuth == nil || availableCustomAuth == nil {
		return nil
	}

	customAuthServiceName := customAuth.GetName()

	// If name is not provided, or doesn't match any available names
	// rely on the default configuration as the per-route config
	if _, ok := availableCustomAuth[customAuthServiceName]; !ok {
		customAuthServiceName = DefaultCustomAuthServiceName
	}

	return &SanitizePerRoute{
		CustomAuthServerName: customAuthServiceName,
	}
}
