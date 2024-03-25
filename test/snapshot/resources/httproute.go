package resources

import (
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/snapshot/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HttpbinHTTPRoute is a simple HTTPRoute that routes traffic to the httpbin upstream
var HttpbinHTTPRoute = &gwv1.HTTPRoute{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "httpbin-route",
		Namespace: "httpbin",
	},
	Spec: gwv1.HTTPRouteSpec{
		CommonRouteSpec: gwv1.CommonRouteSpec{
			ParentRefs: []gwv1.ParentReference{
				{
					Name:      "example-gateway",
					Namespace: utils.PtrTo(gwv1.Namespace(gloodefaults.GlooSystem)),
				},
			},
		},
		Hostnames: []gwv1.Hostname{"httpbin.example.com"},
		Rules: []gwv1.HTTPRouteRule{
			{
				BackendRefs: []gwv1.HTTPBackendRef{
					{
						BackendRef: gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Name:      "httpbin-v1",
								Namespace: utils.PtrTo(gwv1.Namespace("httpbin")),
								Port:      utils.PtrTo(gwv1.PortNumber(8000)),
							},
						},
					},
				},
			},
		},
	},
}
