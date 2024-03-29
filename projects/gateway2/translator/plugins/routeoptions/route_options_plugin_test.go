package routeoptions

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	rtoptquery "github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/testutils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/faultinjection"
	corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = Describe("RouteOptionsPlugin", func() {
	When("applying RouteOptions as Filter", func() {
		It("applies fault injecton RouteOptions directly from resource to output route", func() {
			deps := []client.Object{routeOption()}
			fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
			gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
			plugin := NewPlugin(gwQueries, fakeClient)

			rtCtx := &plugins.RouteContext{
				Route: &gwv1.HTTPRoute{},
				Rule: &gwv1.HTTPRouteRule{
					Filters: []gwv1.HTTPRouteFilter{{
						Type: gwv1.HTTPRouteFilterExtensionRef,
						ExtensionRef: &gwv1.LocalObjectReference{
							Group: gwv1.Group(sologatewayv1.RouteOptionGVK.Group),
							Kind:  gwv1.Kind(sologatewayv1.RouteOptionGVK.Kind),
							Name:  "policy",
						}}}}}

			outputRoute := &v1.Route{
				Options: &v1.RouteOptions{},
			}
			plugin.ApplyRoutePlugin(context.Background(), rtCtx, outputRoute)

			expectedRoute := &v1.Route{
				Options: &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 1,
							HttpStatus: 500,
						},
					},
				},
			}
			Expect(proto.Equal(outputRoute, expectedRoute)).To(BeTrue())
		})
	})

	Describe("Attaching RouteOptions via policy attachemnt as a RoutePlugin", func() {
		When("RouteOptions exist in the same namespace and are attached correctly", func() {
			It("correctly adds faultinjection", func() {
				deps := []client.Object{attachedRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					Route: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "ghostface",
							Namespace: "wu-tang",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 1.00,
							HttpStatus: 500,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())
			})
		})

		When("RouteOptions exist in the same namespace and are attached correctly but omit the namespace in targetRef", func() {
			It("correctly adds faultinjection", func() {
				deps := []client.Object{attachedRouteOptionOmitNamespace()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					Route: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "ghostface",
							Namespace: "wu-tang",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 1.00,
							HttpStatus: 500,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())
			})
		})

		When("Two RouteOptions are attached correctly with different creation timestamps", func() {
			It("correctly adds faultinjection from the earliest created object", func() {
				deps := []client.Object{attachedRouteOption(), attachedRouteOptionBefore()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					Route: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "ghostface",
							Namespace: "wu-tang",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				expectedOptions := &v1.RouteOptions{
					Faults: &faultinjection.RouteFaults{
						Abort: &faultinjection.RouteAbort{
							Percentage: 4.19,
							HttpStatus: 500,
						},
					},
				}
				Expect(proto.Equal(outputRoute.GetOptions(), expectedOptions)).To(BeTrue())
			})
		})

		When("RouteOptions exist in the same namespace but are not attached correctly", func() {
			It("does not add faultinjection", func() {
				deps := []client.Object{nonAttachedRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					Route: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "ghostface",
							Namespace: "wu-tang",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				Expect(proto.Equal(outputRoute.GetOptions(), &v1.RouteOptions{})).To(BeTrue())
			})
		})

		When("RouteOptions exist in a different namespace than the provided routeCtx", func() {
			It("does not add faultinjection", func() {
				deps := []client.Object{attachedRouteOption()}
				fakeClient := testutils.BuildIndexedFakeClient(deps, gwquery.IterateIndices, rtoptquery.IterateIndices)
				gwQueries := testutils.BuildGatewayQueriesWithClient(fakeClient)
				plugin := NewPlugin(gwQueries, fakeClient)

				ctx := context.Background()
				routeCtx := &plugins.RouteContext{
					Route: &gwv1.HTTPRoute{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "ghostface",
							Namespace: "default",
						},
					},
				}

				outputRoute := &v1.Route{
					Options: &v1.RouteOptions{},
				}
				plugin.ApplyRoutePlugin(ctx, routeCtx, outputRoute)

				Expect(proto.Equal(outputRoute.GetOptions(), &v1.RouteOptions{})).To(BeTrue())
			})
		})
	})
})

func routeOption() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "default",
		},
		Spec: sologatewayv1.RouteOption{
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 1.00,
						HttpStatus: 500,
					},
				},
			},
		},
	}
}

func attachedRouteOption() *solokubev1.RouteOption {
	now := metav1.Now()
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "policy",
			Namespace:         "wu-tang",
			CreationTimestamp: now,
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
						HttpStatus: 500,
					},
				},
			},
		},
	}
}

func attachedRouteOptionBefore() *solokubev1.RouteOption {
	anHourAgo := metav1.NewTime(time.Now().Add(-1 * time.Hour))
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "policy-older",
			Namespace:         "wu-tang",
			CreationTimestamp: anHourAgo,
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
						Percentage: 4.19,
						HttpStatus: 500,
					},
				},
			},
		},
	}
}

func attachedRouteOptionOmitNamespace() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "wu-tang",
		},
		Spec: sologatewayv1.RouteOption{
			TargetRef: &corev1.PolicyTargetReference{
				Group: gwv1.GroupVersion.Group,
				Kind:  wellknown.HTTPRouteKind,
				Name:  "ghostface",
			},
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 1.00,
						HttpStatus: 500,
					},
				},
			},
		},
	}
}

func nonAttachedRouteOption() *solokubev1.RouteOption {
	return &solokubev1.RouteOption{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "policy",
			Namespace: "wu-tang",
		},
		Spec: sologatewayv1.RouteOption{
			TargetRef: &corev1.PolicyTargetReference{
				Group:     gwv1.GroupVersion.Group,
				Kind:      wellknown.HTTPRouteKind,
				Name:      "my-route",
				Namespace: wrapperspb.String("wu-tang"),
			},
			Options: &v1.RouteOptions{
				Faults: &faultinjection.RouteFaults{
					Abort: &faultinjection.RouteAbort{
						Percentage: 1.00,
						HttpStatus: 500,
					},
				},
			},
		},
	}
}
