package reports_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"github.com/solo-io/gloo/projects/gateway2/reports"
)

var _ = Describe("Reporting Infrastructure", func() {
	BeforeEach(func() {
	})

	Describe("building gateway status", func() {
		It("should build all positive conditions with an empty report", func() {
			gw := gw()
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize GatewayReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.Gateway(gw)

			status := rm.BuildGWStatus(context.Background(), *gw)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))
		})

		It("should preserve conditions set externally", func() {
			gw := gw()
			meta.SetStatusCondition(&gw.Status.Conditions, metav1.Condition{
				Type:   "gloo.solo.io/SomeCondition",
				Status: metav1.ConditionFalse,
			})
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize GatewayReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.Gateway(gw)

			status := rm.BuildGWStatus(context.Background(), *gw)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(3)) // 2 from the report, 1 from the original status
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))
		})

		It("should correctly set negative gateway conditions from report and not add extra conditions", func() {
			gw := gw()
			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			reporter.Gateway(gw).SetCondition(reports.GatewayCondition{
				Type:   gwv1.GatewayConditionProgrammed,
				Status: metav1.ConditionFalse,
				Reason: gwv1.GatewayReasonAddressNotUsable,
			})
			status := rm.BuildGWStatus(context.Background(), *gw)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			programmed := meta.FindStatusCondition(status.Conditions, string(gwv1.GatewayConditionProgrammed))
			Expect(programmed.Status).To(Equal(metav1.ConditionFalse))
		})

		It("should correctly set negative listener conditions from report and not add extra conditions", func() {
			gw := gw()
			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			reporter.Gateway(gw).Listener(listener()).SetCondition(reports.ListenerCondition{
				Type:   gwv1.ListenerConditionResolvedRefs,
				Status: metav1.ConditionFalse,
				Reason: gwv1.ListenerReasonInvalidRouteKinds,
			})
			status := rm.BuildGWStatus(context.Background(), *gw)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			resolvedRefs := meta.FindStatusCondition(status.Listeners[0].Conditions, string(gwv1.ListenerConditionResolvedRefs))
			Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
		})

		It("should not modify LastTransitionTime for existing conditions that have not changed", func() {
			gw := gw()
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize GatewayReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.Gateway(gw)

			status := rm.BuildGWStatus(context.Background(), *gw)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			acceptedCond := meta.FindStatusCondition(status.Listeners[0].Conditions, string(gwv1.ListenerConditionAccepted))
			oldTransitionTime := acceptedCond.LastTransitionTime

			gw.Status = *status
			status = rm.BuildGWStatus(context.Background(), *gw)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			acceptedCond = meta.FindStatusCondition(status.Listeners[0].Conditions, string(gwv1.ListenerConditionAccepted))
			newTransitionTime := acceptedCond.LastTransitionTime
			Expect(newTransitionTime).To(Equal(oldTransitionTime))
		})

		// TODO(Law): add multiple gws/listener tests
		// TODO(Law): add test confirming transitionTime change when status change
	})

	Describe("building route status", func() {
		DescribeTable("should build all positive route conditions with an empty report",
			func(obj client.Object) {
				rm := reports.NewReportMap()

				reporter := reports.NewReporter(&rm)
				// initialize RouteReporter to mimic translation loop (i.e. report gets initialized for all Routes)
				reporter.Route(obj)

				status := rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(2))
			},
			Entry("regular httproute", httpRoute()),
			Entry("regular tcproute", tcpRoute()),
			Entry("regular tlsroute", tlsRoute()),
			Entry("delegatee route", delegateeRoute()),
		)

		DescribeTable("should preserve conditions set externally",
			func(obj client.Object) {
				rm := reports.NewReportMap()

				reporter := reports.NewReporter(&rm)
				// initialize RouteReporter to mimic translation loop (i.e. report gets initialized for all Routes)
				reporter.Route(obj)

				status := rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(3)) // 2 from the report, 1 from the original status
			},
			Entry("regular httproute", httpRoute(
				metav1.Condition{
					Type: "gloo.solo.io/SomeCondition",
				},
			)),
			Entry("regular tcproute", tcpRoute(
				metav1.Condition{
					Type: "gloo.solo.io/SomeCondition",
				},
			)),
			Entry("regular tlsroute", tlsRoute(
				metav1.Condition{
					Type: "gloo.solo.io/SomeCondition",
				},
			)),
			Entry("delegatee route", delegateeRoute(
				metav1.Condition{
					Type: "gloo.solo.io/SomeCondition",
				},
			)),
		)

		DescribeTable("should correctly set negative route conditions from report and not add extra conditions",
			func(obj client.Object, parentRef *gwv1.ParentReference) {
				rm := reports.NewReportMap()
				reporter := reports.NewReporter(&rm)
				reporter.Route(obj).ParentRef(parentRef).SetCondition(reports.RouteCondition{
					Type:   gwv1.RouteConditionResolvedRefs,
					Status: metav1.ConditionFalse,
					Reason: gwv1.RouteReasonBackendNotFound,
				})

				status := rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(2))

				resolvedRefs := meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
			},
			Entry("regular httproute", httpRoute(), parentRef()),
			Entry("regular tcproute", tcpRoute(), parentRef()),
			Entry("regular tlsroute", tlsRoute(), parentRef()),
			Entry("delegatee route", delegateeRoute(), parentRouteRef()),
		)

		DescribeTable("should filter out multiple negative route conditions of the same type from report",
			func(obj client.Object, parentRef *gwv1.ParentReference) {
				rm := reports.NewReportMap()
				reporter := reports.NewReporter(&rm)
				reporter.Route(obj).ParentRef(parentRef).SetCondition(reports.RouteCondition{
					Type:   gwv1.RouteConditionResolvedRefs,
					Status: metav1.ConditionFalse,
					Reason: gwv1.RouteReasonBackendNotFound,
				})
				reporter.Route(obj).ParentRef(parentRef).SetCondition(reports.RouteCondition{
					Type:   gwv1.RouteConditionResolvedRefs,
					Status: metav1.ConditionFalse,
					Reason: gwv1.RouteReasonBackendNotFound,
				})

				status := rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(2))

				resolvedRefs := meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
			},
			Entry("regular httproute", httpRoute(), parentRef()),
			Entry("regular tcproute", tcpRoute(), parentRef()),
			Entry("regular tlsroute", tlsRoute(), parentRef()),
			Entry("delegatee route", delegateeRoute(), parentRouteRef()),
		)

		DescribeTable("should not modify LastTransitionTime for existing conditions that have not changed",
			func(obj client.Object) {
				rm := reports.NewReportMap()

				reporter := reports.NewReporter(&rm)
				// initialize RouteReporter to mimic translation loop (i.e. report gets initialized for all Routes)
				reporter.Route(obj)

				status := rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(2))

				resolvedRefs := meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				oldTransitionTime := resolvedRefs.LastTransitionTime

				// Type assert the object to update the Status field based on its type
				switch route := obj.(type) {
				case *gwv1.HTTPRoute:
					route.Status.RouteStatus = *status
				case *gwv1a2.TCPRoute:
					route.Status.RouteStatus = *status
				case *gwv1a2.TLSRoute:
					route.Status.RouteStatus = *status
				default:
					Fail(fmt.Sprintf("unsupported route type: %T", obj))
				}

				status = rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(2))

				resolvedRefs = meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				newTransitionTime := resolvedRefs.LastTransitionTime
				Expect(newTransitionTime).To(Equal(oldTransitionTime))
			},
			Entry("regular httproute", httpRoute()),
			Entry("delegatee route", delegateeRoute()),
			Entry("regular tcproute", tcpRoute()),
			Entry("regular tlsroute", tlsRoute()),
		)

		DescribeTable("should correctly handle multiple ParentRefs on a route",
			func(obj client.Object) {
				// Add an additional ParentRef to test multiple parent references handling
				switch route := obj.(type) {
				case *gwv1.HTTPRoute:
					route.Spec.ParentRefs = append(route.Spec.ParentRefs, gwv1.ParentReference{
						Name: "additional-gateway",
					})
				case *gwv1a2.TCPRoute:
					route.Spec.ParentRefs = append(route.Spec.ParentRefs, gwv1.ParentReference{
						Name: "additional-gateway",
					})
				case *gwv1a2.TLSRoute:
					route.Spec.ParentRefs = append(route.Spec.ParentRefs, gwv1.ParentReference{
						Name: "additional-gateway",
					})
				default:
					Fail(fmt.Sprintf("unsupported route type: %T", obj))
				}

				rm := reports.NewReportMap()
				reporter := reports.NewReporter(&rm)

				// Initialize RouteReporter to mimic translation loop
				reporter.Route(obj)

				status := rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(2))

				// Check that each parent has the correct number of conditions
				for _, parent := range status.Parents {
					Expect(parent.Conditions).To(HaveLen(2))
				}
			},
			Entry("regular HTTPRoute", httpRoute()),
			Entry("regular TCPRoute", tcpRoute()),
			Entry("regular TLSRoute", tlsRoute()),
		)

		DescribeTable("should correctly associate multiple routes with shared and separate listeners",
			func(route1, route2 client.Object, listener1, listener2 gwv1.Listener) {
				gw := gw()
				gw.Spec.Listeners = []gwv1.Listener{listener1, listener2}

				// Assign the first listener to the first route's parent ref
				switch r1 := route1.(type) {
				case *gwv1.HTTPRoute:
					r1.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName(listener1.Name))
				case *gwv1a2.TCPRoute:
					r1.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName(listener1.Name))
				case *gwv1a2.TLSRoute:
					r1.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName(listener1.Name))
				}

				// Assign the second listener to the second route's parent ref
				switch r2 := route2.(type) {
				case *gwv1.HTTPRoute:
					r2.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName(listener2.Name))
				case *gwv1a2.TCPRoute:
					r2.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName(listener2.Name))
				case *gwv1a2.TLSRoute:
					r2.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName(listener2.Name))
				}

				rm := reports.NewReportMap()
				reporter := reports.NewReporter(&rm)

				// Initialize RouteReporter to mimic translation loop
				reporter.Route(route1)
				reporter.Route(route2)

				status1 := rm.BuildRouteStatus(context.Background(), route1, "gloo-gateway")
				status2 := rm.BuildRouteStatus(context.Background(), route2, "gloo-gateway")

				Expect(status1).NotTo(BeNil())
				Expect(status1.Parents[0].Conditions).To(HaveLen(2))

				Expect(status2).NotTo(BeNil())
				Expect(status2.Parents[0].Conditions).To(HaveLen(2))
			},
			Entry("HTTPRoutes with shared and separate listeners",
				httpRoute(), httpRoute(),
				gwv1.Listener{Name: "foo-http", Protocol: gwv1.HTTPProtocolType},
				gwv1.Listener{Name: "bar-http", Protocol: gwv1.HTTPProtocolType},
			),
			Entry("TCPRoutes with shared and separate listeners",
				tcpRoute(), tcpRoute(),
				gwv1.Listener{Name: "foo-tcp", Protocol: gwv1.TCPProtocolType},
				gwv1.Listener{Name: "bar-tcp", Protocol: gwv1.TCPProtocolType},
			),
			Entry("TLSRoutes with shared and separate listeners",
				tlsRoute(), tlsRoute(),
				gwv1.Listener{Name: "foo-tls", Protocol: gwv1.TLSProtocolType},
				gwv1.Listener{Name: "bar-tls", Protocol: gwv1.TLSProtocolType},
			),
		)
	})

	Describe("building listener set status", func() {
		It("should build all positive conditions with an empty report", func() {
			ls := ls()
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize ListenerSetReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.ListenerSet(ls)

			status := rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))
		})

		It("should preserve conditions set externally", func() {
			ls := ls()
			meta.SetStatusCondition(&ls.Status.Conditions, metav1.Condition{
				Type:   "gloo.solo.io/SomeCondition",
				Status: metav1.ConditionFalse,
			})
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize ListenerSetReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.ListenerSet(ls)

			status := rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(3)) // 2 from the report, 1 from the original status
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))
		})

		It("should correctly set negative gateway conditions from report and not add extra conditions", func() {
			ls := ls()
			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			reporter.ListenerSet(ls).SetCondition(reports.GatewayCondition{
				Type:   gwv1.GatewayConditionProgrammed,
				Status: metav1.ConditionFalse,
				Reason: gwv1.GatewayReasonAddressNotUsable,
			})
			status := rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			programmed := meta.FindStatusCondition(status.Conditions, string(gwv1.GatewayConditionProgrammed))
			Expect(programmed.Status).To(Equal(metav1.ConditionFalse))
		})

		It("should correctly set negative listener conditions from report and not add extra conditions", func() {
			ls := ls()
			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			reporter.ListenerSet(ls).Listener(listener()).SetCondition(reports.ListenerCondition{
				Type:   gwv1.ListenerConditionResolvedRefs,
				Status: metav1.ConditionFalse,
				Reason: gwv1.ListenerReasonInvalidRouteKinds,
			})
			status := rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			resolvedRefs := meta.FindStatusCondition(status.Listeners[0].Conditions, string(gwv1.ListenerConditionResolvedRefs))
			Expect(resolvedRefs.Status).To(Equal(metav1.ConditionFalse))
		})

		It("should not modify LastTransitionTime for existing conditions that have not changed", func() {
			ls := ls()
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize ListenerSetReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.ListenerSet(ls)

			status := rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			acceptedCond := meta.FindStatusCondition(status.Listeners[0].Conditions, string(gwv1.ListenerConditionAccepted))
			oldTransitionTime := acceptedCond.LastTransitionTime

			ls.Status = *status
			status = rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(HaveLen(1))
			Expect(status.Listeners[0].Conditions).To(HaveLen(4))

			acceptedCond = meta.FindStatusCondition(status.Listeners[0].Conditions, string(gwv1.ListenerConditionAccepted))
			newTransitionTime := acceptedCond.LastTransitionTime
			Expect(newTransitionTime).To(Equal(oldTransitionTime))
		})

		It("should not add status for listeners on a rejected listener set", func() {
			ls := ls()
			rm := reports.NewReportMap()

			reporter := reports.NewReporter(&rm)
			// initialize ListenerSetReporter to mimic translation loop (i.e. report gets initialized for all GWs)
			reporter.ListenerSet(ls).SetCondition(reports.GatewayCondition{
				Type:   gwv1.GatewayConditionAccepted,
				Status: metav1.ConditionFalse,
				Reason: gwv1.GatewayConditionReason(gwxv1a1.ListenerSetReasonNotAllowed),
			})
			reporter.ListenerSet(ls).SetCondition(reports.GatewayCondition{
				Type:   gwv1.GatewayConditionProgrammed,
				Status: metav1.ConditionFalse,
				Reason: gwv1.GatewayConditionReason(gwxv1a1.ListenerSetReasonNotAllowed),
			})

			status := rm.BuildListenerSetStatus(context.Background(), *ls)

			Expect(status).NotTo(BeNil())
			Expect(status.Conditions).To(HaveLen(2))
			Expect(status.Listeners).To(BeEmpty())
		})
	})

	DescribeTable("should handle routes with missing parent references gracefully",
		func(route client.Object) {
			// Remove ParentRefs from the route.
			switch r := route.(type) {
			case *gwv1.HTTPRoute:
				r.Spec.ParentRefs = nil
			case *gwv1a2.TCPRoute:
				r.Spec.ParentRefs = nil
			case *gwv1a2.TLSRoute:
				r.Spec.ParentRefs = nil
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)

			// Initialize RouteReporter to mimic translation loop
			reporter.Route(route)
			status := rm.BuildRouteStatus(context.Background(), route, "gloo-gateway")

			Expect(status).NotTo(BeNil())
			Expect(status.Parents).To(BeEmpty())
		},
		Entry("HTTPRoute with missing parent reference", httpRoute()),
		Entry("TCPRoute with missing parent reference", tcpRoute()),
		Entry("TLSRoute with missing parent reference", tlsRoute()),
	)
})

func httpRoute(conditions ...metav1.Condition) client.Object {
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRef())
	if len(conditions) > 0 {
		route.Status.Parents = append(route.Status.Parents, gwv1.RouteParentStatus{
			ParentRef:  *parentRef(),
			Conditions: conditions,
		})
	}

	return route
}

func tcpRoute(conditions ...metav1.Condition) client.Object {
	route := &gwv1a2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRef())
	if len(conditions) > 0 {
		route.Status.Parents = append(route.Status.Parents, gwv1.RouteParentStatus{
			ParentRef:  *parentRef(),
			Conditions: conditions,
		})
	}
	return route
}

func tlsRoute(conditions ...metav1.Condition) client.Object {
	route := &gwv1a2.TLSRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRef())
	if len(conditions) > 0 {
		route.Status.Parents = append(route.Status.Parents, gwv1.RouteParentStatus{
			ParentRef:  *parentRef(),
			Conditions: conditions,
		})
	}
	return route
}

func parentRef() *gwv1.ParentReference {
	return &gwv1.ParentReference{
		Name: "parent",
	}
}

func delegateeRoute(conditions ...metav1.Condition) client.Object {
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "child-route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRouteRef())
	if len(conditions) > 0 {
		route.Status.Parents = append(route.Status.Parents, gwv1.RouteParentStatus{
			ParentRef:  *parentRouteRef(),
			Conditions: conditions,
		})
	}
	return route
}

func parentRouteRef() *gwv1.ParentReference {
	return &gwv1.ParentReference{
		Group:     ptr.To(gwv1.Group("gateway.networking.k8s.io")),
		Kind:      ptr.To(gwv1.Kind("HTTPRoute")),
		Name:      "parent-route",
		Namespace: ptr.To(gwv1.Namespace("default")),
	}
}

var _ = Describe("ReportMap.Equals", func() {
	// buildReportMap mimics one translation pass: it allocates a fresh
	// ReportMap (with fresh report pointers) and populates a Gateway and an
	// HTTPRoute parentRef condition.
	buildReportMap := func() reports.ReportMap {
		gw := gw()
		route := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "route"},
		}
		rm := reports.NewReportMap()
		reporter := reports.NewReporter(&rm)
		reporter.Gateway(gw).SetCondition(reports.GatewayCondition{
			Type:   gwv1.GatewayConditionProgrammed,
			Status: metav1.ConditionTrue,
			Reason: gwv1.GatewayReasonProgrammed,
		})
		reporter.Route(route).ParentRef(&gwv1.ParentReference{
			Name: "test",
		}).SetCondition(reports.RouteCondition{
			Type:   gwv1.RouteConditionAccepted,
			Status: metav1.ConditionTrue,
			Reason: gwv1.RouteReasonAccepted,
		})
		return rm
	}

	It("returns true for two independently-built, identical report maps", func() {
		// Regression: Equals must compare report content, not pointer identity.
		// The maps hold *RouteReport/*GatewayReport, so a pointer-based compare
		// always reports "changed" and pegs the control plane re-translating.
		a := buildReportMap()
		b := buildReportMap()
		Expect(a.Equals(b)).To(BeTrue(), "identical report content must compare equal")
	})

	It("returns false when a condition differs", func() {
		a := buildReportMap()
		b := buildReportMap()

		route := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "route"},
		}
		reports.NewReporter(&b).Route(route).ParentRef(&gwv1.ParentReference{
			Name: "test",
		}).SetCondition(reports.RouteCondition{
			Type:   gwv1.RouteConditionAccepted,
			Status: metav1.ConditionFalse,
			Reason: gwv1.RouteReasonBackendNotFound,
		})

		Expect(a.Equals(b)).To(BeFalse(), "differing condition must compare not-equal")
	})

	It("returns false when route sets differ", func() {
		a := buildReportMap()
		b := buildReportMap()
		extra := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "route2"},
		}
		reports.NewReporter(&b).Route(extra)
		Expect(a.Equals(b)).To(BeFalse(), "different route key set must compare not-equal")
	})

	It("returns false when route observedGeneration differs", func() {
		buildWithGeneration := func(gen int64) reports.ReportMap {
			route := &gwv1.HTTPRoute{
				ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "route", Generation: gen},
			}
			rm := reports.NewReportMap()
			reports.NewReporter(&rm).Route(route).ParentRef(&gwv1.ParentReference{
				Name: "test",
			}).SetCondition(reports.RouteCondition{
				Type:   gwv1.RouteConditionAccepted,
				Status: metav1.ConditionTrue,
				Reason: gwv1.RouteReasonAccepted,
			})
			return rm
		}
		a := buildWithGeneration(1)
		b := buildWithGeneration(2)
		Expect(a.Equals(b)).To(BeFalse(), "differing route observedGeneration must compare not-equal")
	})

	// supportedKinds mimics translation building RouteGroupKinds with FRESH Group
	// pointers each pass (see buildDefaultRouteKindsForProtocol). Each call returns
	// equal-by-value kinds backed by distinct *Group pointers.
	supportedKinds := func() []gwv1.RouteGroupKind {
		http := gwv1.Group(gwv1.GroupName)
		grpc := gwv1.Group(gwv1.GroupName)
		return []gwv1.RouteGroupKind{
			{Group: &http, Kind: gwv1.Kind("HTTPRoute")},
			{Group: &grpc, Kind: gwv1.Kind("GRPCRoute")},
		}
	}

	buildWithSupportedKinds := func(kinds []gwv1.RouteGroupKind) reports.ReportMap {
		gw := gw()
		rm := reports.NewReportMap()
		reports.NewReporter(&rm).Gateway(gw).Listener(listener()).SetSupportedKinds(kinds)
		return rm
	}

	It("returns true when SupportedKinds are equal by value but use fresh Group pointers", func() {
		// Regression: RouteGroupKind.Group is a *Group; comparing the struct with
		// == compares the pointer, so identical kinds with fresh pointers (as
		// translation produces every pass) would otherwise look "changed" and keep
		// the hot loop alive for normal valid listeners.
		a := buildWithSupportedKinds(supportedKinds())
		b := buildWithSupportedKinds(supportedKinds())
		Expect(a.Equals(b)).To(BeTrue(), "value-equal SupportedKinds with fresh Group pointers must compare equal")
	})

	It("returns true when SupportedKinds match but are in a different order", func() {
		// The kinds builder ranges a map, so slice order is not stable across passes.
		// Build the reversed slice from a fresh supportedKinds() call so it also uses
		// distinct *Group pointers; equality must hold across BOTH order and pointer
		// differences, as translation produces every pass.
		fresh := supportedKinds()
		reversed := []gwv1.RouteGroupKind{fresh[1], fresh[0]}
		a := buildWithSupportedKinds(supportedKinds())
		b := buildWithSupportedKinds(reversed)
		Expect(a.Equals(b)).To(BeTrue(), "SupportedKinds must compare order-insensitively")
	})

	It("returns false when SupportedKinds differ in value", func() {
		http := gwv1.Group(gwv1.GroupName)
		a := buildWithSupportedKinds([]gwv1.RouteGroupKind{{Group: &http, Kind: gwv1.Kind("HTTPRoute")}})
		tcp := gwv1.Group(gwv1.GroupName)
		b := buildWithSupportedKinds([]gwv1.RouteGroupKind{{Group: &tcp, Kind: gwv1.Kind("TCPRoute")}})
		Expect(a.Equals(b)).To(BeFalse(), "different SupportedKinds must compare not-equal")
	})

	It("remains equal to fresh translation output after gateway status rendering", func() {
		gateway := gw()
		build := func() reports.ReportMap {
			rm := reports.NewReportMap()
			reports.NewReporter(&rm).Gateway(gateway)
			return rm
		}

		rm := build()
		expected := build()

		Expect(rm.BuildGWStatus(context.Background(), *gateway)).NotTo(BeNil())
		Expect(rm.Equals(expected)).To(BeTrue(), "status rendering must not add default conditions to the report")
	})

	It("remains equal to fresh translation output after listener set status rendering", func() {
		listenerSet := ls()
		build := func() reports.ReportMap {
			rm := reports.NewReportMap()
			reports.NewReporter(&rm).ListenerSet(listenerSet)
			return rm
		}

		rm := build()
		expected := build()

		Expect(rm.BuildListenerSetStatus(context.Background(), *listenerSet)).NotTo(BeNil())
		Expect(rm.Equals(expected)).To(BeTrue(), "status rendering must not add listener set defaults to the report")
	})

	It("remains equal to fresh translation output after route status rendering", func() {
		route := httpRoute()
		build := func() reports.ReportMap {
			rm := reports.NewReportMap()
			reports.NewReporter(&rm).Route(route)
			return rm
		}

		rm := build()
		expected := build()

		Expect(rm.BuildRouteStatus(context.Background(), route, "gloo-gateway")).NotTo(BeNil())
		Expect(rm.Equals(expected)).To(BeTrue(), "status rendering must not add parent defaults to the report")
	})
})

func gw() *gwv1.Gateway {
	gw := &gwv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test",
		},
	}
	gw.Spec.Listeners = append(gw.Spec.Listeners, *listener())
	return gw
}

func listener() *gwv1.Listener {
	return &gwv1.Listener{
		Name: "http",
	}
}

func ls() *gwxv1a1.XListenerSet {
	ls := &gwxv1a1.XListenerSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test",
		},
	}
	ls.Spec.Listeners = []gwxv1a1.ListenerEntry{
		{
			Name: "http",
		},
	}
	return ls
}
