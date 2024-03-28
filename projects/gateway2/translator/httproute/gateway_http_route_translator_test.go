package httproute_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/httproute"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/faultinjection"
	corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("GatewayHttpRouteTranslator", func() {
})

var _ = Describe("HTTPRoute translation with RouteOptions", func() {
	When("HTTPRoute with RouteOptions filter AND attached RouteOptions", func() {
		It("Only applies RouteOptions from filter", func() {
			ctx := context.Background()

			deps := []client.Object{
				&solokubev1.RouteOption{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "attached-policy",
						Namespace: "wu-tang",
					},
					Spec: sologatewayv1.RouteOption{
						TargetRef: &corev1.PolicyTargetReference{
							Group:     gwv1.GroupVersion.Group,
							Kind:      wellknown.HTTPRouteKind,
							Name:      "ghostface",
							Namespace: wrapperspb.String("wu-tang"),
						},
						Options: &v1.RouteOptions{
							Faults: &faultinjection.RouteFaults{
								Abort: &faultinjection.RouteAbort{
									Percentage: 1.00,
								},
							},
						},
					},
				},
				&solokubev1.RouteOption{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "filter-policy",
						Namespace: "wu-tang",
					},
					Spec: sologatewayv1.RouteOption{
						Options: &v1.RouteOptions{
							Faults: &faultinjection.RouteFaults{
								Abort: &faultinjection.RouteAbort{
									Percentage: 4.19,
								},
							},
						},
					},
				},
			}
			queries := testutils.BuildGatewayQueries(deps)
			pluginRegistry := registry.NewPluginRegistry(registry.BuildPlugins(queries))

			listener := gwv1.Listener{}
			parentRef := gwv1.ParentReference{
				Name: "my-gw",
			}
			route := gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "ghostface",
					Namespace: "wu-tang",
				},
				Spec: gwv1.HTTPRouteSpec{
					CommonRouteSpec: gwv1.CommonRouteSpec{
						ParentRefs: []gwv1.ParentReference{
							parentRef,
						},
					},
					Rules: []gwv1.HTTPRouteRule{{
						Filters: []gwv1.HTTPRouteFilter{{
							Type: "ExtensionRef",
							ExtensionRef: &gwv1.LocalObjectReference{
								Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
								Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
								Name:  "filter-policy",
							}},
						}},
					},
				},
			}

			reportsMap := reports.NewReportMap()
			reporter := reports.NewReporter(&reportsMap)
			parentRefReporter := reporter.Route(&route).ParentRef(&parentRef)

			glooRoutes := httproute.TranslateGatewayHTTPRouteRules(
				ctx,
				pluginRegistry,
				queries,
				listener,
				route,
				parentRefReporter,
			)
			Expect(glooRoutes).To(HaveLen(1))
			Expect(glooRoutes[0].GetOptions()).To(Not(BeNil()))

			expectedOptions := &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 4.19,
					},
				},
			}
			Expect(proto.Equal(glooRoutes[0].GetOptions(), expectedOptions)).To(BeTrue())
		})
	})
})
