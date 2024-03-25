package resources

import (
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/core/matchers"
	gloocore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Simple VirtualService that routes traffic to the httpbin upstream
var HttpbinVirtualService = &v1.VirtualService{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "httpbin-route",
		Namespace: "httpbin",
	},
	Spec: v1.VirtualServiceSpec{
		VirtualHost: &v1.VirtualHost{
			Domains: []string{"httpbin.example.com"},
			Routes: []*v1.Route{{
				Matchers: []*matchers.Matcher{
					{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/",
						},
					},
				},
				Action: &v1.Route_RouteAction{
					RouteAction: &gloov1.RouteAction{
						Destination: &gloov1.RouteAction_Single{
							Single: &gloov1.Destination{
								DestinationType: &gloov1.Destination_Upstream{
									Upstream: &gloocore.ResourceRef{
										Name:      "httpbin-v1-httpbin-8000",
										Namespace: gloodefaults.GlooSystem,
									},
								},
							},
						},
					},
				},
			}},
		},
	},
}
