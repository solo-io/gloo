package registry

import (
	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/headermodifier"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/mirror"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/redirect"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/urlrewrite"
	plugins2 "github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PluginRegistry interface {

}

func (h *pluginRegistry) GetRoutePlugins() []plugins.RoutePlugin {
	return h.routePlugins
}

type pluginRegistry struct {
	routePlugins []plugins.RoutePlugin
	namespacePlugins []plugins.NamespacePlugin
}

func (h *pluginRegistry) GetRoutePlugins() []plugins.RoutePlugin {
	return h.routePlugins
}

func NewPluginRegistry(queries query.GatewayQueries, cli client.Client, scheme *runtime.Scheme, makeExtensionPlugins func() []plugins2.Plugin) PluginRegistry {
	var (
		routePlugins []plugins.RoutePlugin
		namespacePlugins []plugins.NamespacePlugin
	)


	allPlugins := append( buildPlugins(queries),

	for _, plugin := range allPlugins {
		if routePlugin, ok := plugin.(plugins.RoutePlugin); ok {
			routePlugins = append(routePlugins, routePlugin)
		}
		if namespacePlugin, ok := plugin.(plugins.NamespacePlugin); ok {
			namespacePlugins = append(namespacePlugins, namespacePlugin)
		}
	}
	return &PluginRegistry{
		routePlugins: routePlugins,
		namespacePlugins: namespacePlugins,
	}
}

// This function returns the full set of plugins to be registered.
// New plugins should be added to this list (and only this list).
// If modification of this list is needed for testing etc,
// we can add a new registry constructor that accepts this function
func buildPlugins(queries query.GatewayQueries) []plugins.Plugin {
	return []plugins.Plugin{
		headermodifier.NewPlugin(),
		mirror.NewPlugin(queries),
		redirect.NewPlugin(),
		routeoptions.NewPlugin(queries),
		urlrewrite.NewPlugin(),
	}
}
