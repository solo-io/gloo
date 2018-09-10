package translator

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/todo"
)

// trnslate a snapshot of schemas and resolvermaps to:
// a 1 proxy for the snapshot, assigned to the sqoop sidecar
func Translate(writeNamespace string, snap *v1.ApiSnapshot, resourceErrs reporter.ResourceErrors) *gloov1.Proxy {
	ourRoutes := routesForResolverMaps(snap.ResolverMaps.List(), resourceErrs)

	var routes []*gloov1.Route

	for _, r := range  ourRoutes {
		routes = append(routes, &gloov1.Route{
			Action: &gloov1.Route_RouteAction{RouteAction: r.action},
		})
	}

	return &gloov1.Proxy{
		Metadata: core.Metadata{
			Name:        "-proxy",
			Namespace:   writeNamespace,
			Annotations: map[string]string{"owner_ref": "gateway"},
		},
		Listeners: []*gloov1.Listener{
			{
				// TODO (ilackarms): make this section configurable
				Name:        "sqoop-listener",
				BindAddress: TODO.SqoopSidecarBindAddr,
				BindPort:    TODO.SqoopSidecarBindPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{
							{
								Name:    "sqoop-vhost",
								Domains: []string{"*"},
								Routes:  routes,
							},
						},
					},
				},
				// TODO(ilackarms / yuval-k): decide if we need ssl for connecting to sidecar
				SslConfiguations: nil,
			},
		},
	}
}
