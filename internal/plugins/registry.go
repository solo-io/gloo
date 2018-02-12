package plugins

import (
	"github.com/solo-io/glue/internal/bootstrap"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/plugin"
)

var defaultRegistry = &registry{}

type EndpointDiscoveryInitFunc func(options bootstrap.Options, stopCh <-chan struct{}) (endpointdiscovery.Interface, error)

func Register(plugin plugin.TranslatorPlugin, startEndpointDiscovery EndpointDiscoveryInitFunc) {
	defaultRegistry.plugins = append(defaultRegistry.plugins, plugin)
	if startEndpointDiscovery != nil {
		defaultRegistry.endpointDiscoveries = append(defaultRegistry.endpointDiscoveries, startEndpointDiscovery)
	}
}

func RegisteredPlugins() []plugin.TranslatorPlugin {
	return defaultRegistry.plugins
}

func EndpointDiscoveryInitializers() []EndpointDiscoveryInitFunc {
	return defaultRegistry.endpointDiscoveries
}

type registry struct {
	plugins             []plugin.TranslatorPlugin
	endpointDiscoveries []EndpointDiscoveryInitFunc
}
