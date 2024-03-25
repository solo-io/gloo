package resources

import (
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/snapshot/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// SimpleGateway is a simple Gateway that http traffic from all namespaces
var SimpleGateway = &gwv1.Gateway{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example-gateway",
		Namespace: gloodefaults.GlooSystem,
	},
	Spec: gwv1.GatewaySpec{
		GatewayClassName: "gloo-gateway",
		Listeners: []gwv1.Listener{
			{
				Name:     "http",
				Port:     8080,
				Protocol: "HTTP",
				AllowedRoutes: &gwv1.AllowedRoutes{
					Namespaces: &gwv1.RouteNamespaces{
						From: utils.PtrTo(gwv1.NamespacesFromAll),
					},
				},
			},
		},
	},
}
