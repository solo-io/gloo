package setup

import (
	"context"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/proxyprotocol"

	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/advanced_http"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/aws"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/dlp"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/leftmost_xff_address"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/proxylatency"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/rbac"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/sanitize_cluster_header"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/transformer"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/waf"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/wasm"

	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
)

func GetPluginRegistryFactory(
	opts bootstrap.Opts,
	apiEmitterChan chan struct{},
	licensedFeatureProvider *license.LicensedFeatureProvider,
) plugins.PluginRegistryFactory {
	return func(ctx context.Context) plugins.PluginRegistry {
		// Start with open source plugins
		availablePlugins := getOpenSourcePlugins(opts)

		// Process Enterprise feature
		enterpriseFeatureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
		if enterpriseFeatureState.Enabled {
			availablePlugins = reconcileEnterprisePlugins(availablePlugins, getEnterprisePlugins(apiEmitterChan))
		}
		if enterpriseFeatureState.Reason != "" {
			log.Debugf("%s", enterpriseFeatureState.Reason)
		}

		// Process GraphQL feature
		graphQLFeatureState := licensedFeatureProvider.GetStateForLicensedFeature(license.GraphQL)
		if graphQLFeatureState.Enabled {
			availablePlugins = reconcileGraphQLPlugins(availablePlugins, getGraphQLPlugins())
		}
		if graphQLFeatureState.Reason != "" {
			log.Debugf("%s", graphQLFeatureState.Reason)
		}

		// Load the reconciled set of plugins into the registry
		return registry.NewPluginRegistry(availablePlugins)
	}
}

func getOpenSourcePlugins(opts bootstrap.Opts) []plugins.Plugin {
	return registry.Plugins(opts)
}

func getEnterprisePlugins(apiEmitterChan chan struct{}) []plugins.Plugin {
	return []plugins.Plugin{
		ratelimit.NewPlugin(),
		extauth.NewPlugin(),
		sanitize_cluster_header.NewPlugin(),
		rbac.NewPlugin(),
		jwt.NewPlugin(),
		waf.NewPlugin(),
		dlp.NewPlugin(),
		aws.NewPlugin(),
		proxylatency.NewPlugin(),
		failover.NewFailoverPlugin(
			utils.NewSslConfigTranslator(),
			failover.NewDnsResolver(),
			apiEmitterChan,
		),
		advanced_http.NewPlugin(),
		wasm.NewPlugin(),
		leftmost_xff_address.NewPlugin(),
		transformer.NewPlugin(),
		proxyprotocol.NewPlugin(),
	}
}

func getGraphQLPlugins() []plugins.Plugin {
	return []plugins.Plugin{
		graphql.NewPlugin(),
	}
}

func reconcileEnterprisePlugins(currentPlugins, enterprisePlugins []plugins.Plugin) []plugins.Plugin {
	enterprisePluginsByName := make(map[string]bool)
	for _, enterprisePlugin := range enterprisePlugins {
		enterprisePluginsByName[enterprisePlugin.Name()] = true
	}

	var pluginsToDrop []int
	for i, currentPlugin := range currentPlugins {
		if _, inMap := enterprisePluginsByName[currentPlugin.Name()]; inMap {
			// An upgraded version of this plug exists,
			// mark this one for removal
			pluginsToDrop = append(pluginsToDrop, i)
		}
	}

	// Walk back through the currentPlugins and remove the redundant plugins
	for i := len(pluginsToDrop) - 1; i >= 0; i-- {
		badIndex := pluginsToDrop[i]
		currentPlugins = append(currentPlugins[:badIndex], currentPlugins[badIndex+1:]...)
	}

	// Append all enterprise plugins to the end of the list
	return append(currentPlugins, enterprisePlugins...)
}

func reconcileGraphQLPlugins(currentPlugins, graphQLPlugins []plugins.Plugin) []plugins.Plugin {
	// GraphQL plugins are only defined in GlooE and thus can just be appended to the end of the list
	return append(currentPlugins, graphQLPlugins...)
}
