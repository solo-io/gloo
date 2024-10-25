package reports_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

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
			Entry("delegatee route", delegateeRoute()),
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
				default:
					Fail(fmt.Sprintf("unsupported route type: %T", obj))
				}

				status = rm.BuildRouteStatus(context.Background(), obj, "gloo-gateway")

				Expect(status).NotTo(BeNil())
				Expect(status.Parents).To(HaveLen(1))
				Expect(status.Parents[0].Conditions).To(HaveLen(2))

				resolvedRefs = meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionResolvedRefs))
				newTransitionTime := resolvedRefs.LastTransitionTime

				// Use a comparison that accounts for minor precision differences in the transition times
				Expect(newTransitionTime.Time).To(BeTemporally("~", oldTransitionTime.Time, time.Millisecond),
					"Expected LastTransitionTime to not change significantly")
			},
			Entry("regular httproute", httpRoute()),
			Entry("delegatee route", delegateeRoute()),
			Entry("regular tcproute", tcpRoute()),
		)
	})

	It("should correctly associate multiple TCPRoutes with shared and separate listeners", func() {
		gw := gw()
		gw.Spec.Listeners = []gwv1.Listener{
			{Name: "foo-tcp", Protocol: gwv1.TCPProtocolType},
			{Name: "bar-tcp", Protocol: gwv1.TCPProtocolType},
		}

		tcpRoute1 := tcpRoute().(*gwv1a2.TCPRoute)
		tcpRoute1.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName("foo-tcp"))

		tcpRoute2 := tcpRoute().(*gwv1a2.TCPRoute)
		tcpRoute2.Spec.ParentRefs[0].SectionName = ptr.To(gwv1.SectionName("bar-tcp"))

		rm := reports.NewReportMap()
		reporter := reports.NewReporter(&rm)

		reporter.Route(tcpRoute1)
		reporter.Route(tcpRoute2)

		status1 := rm.BuildRouteStatus(context.Background(), tcpRoute1, "gloo-gateway")
		status2 := rm.BuildRouteStatus(context.Background(), tcpRoute2, "gloo-gateway")

		Expect(status1).NotTo(BeNil())
		Expect(status1.Parents[0].Conditions).To(HaveLen(2))

		Expect(status2).NotTo(BeNil())
		Expect(status2.Parents[0].Conditions).To(HaveLen(2))
	})

	It("should handle TCPRoute with missing parent reference gracefully", func() {
		tcpRoute := tcpRoute().(*gwv1a2.TCPRoute)
		tcpRoute.Spec.ParentRefs = nil

		rm := reports.NewReportMap()
		reporter := reports.NewReporter(&rm)

		reporter.Route(tcpRoute)
		status := rm.BuildRouteStatus(context.Background(), tcpRoute, "gloo-gateway")

		Expect(status).NotTo(BeNil())
		Expect(status.Parents).To(BeEmpty())
	})

	It("should correctly handle multiple ParentRefs on a TCPRoute", func() {
		tcpRoute := tcpRoute().(*gwv1a2.TCPRoute)
		tcpRoute.Spec.ParentRefs = append(tcpRoute.Spec.ParentRefs, gwv1.ParentReference{
			Name: "additional-gateway",
		})

		rm := reports.NewReportMap()
		reporter := reports.NewReporter(&rm)

		reporter.Route(tcpRoute)
		status := rm.BuildRouteStatus(context.Background(), tcpRoute, "gloo-gateway")

		Expect(status).NotTo(BeNil())
		Expect(status.Parents).To(HaveLen(2))
		for _, parent := range status.Parents {
			Expect(parent.Conditions).To(HaveLen(2))
		}
	})

	It("should update LastTransitionTime when TCPRoute condition changes", func() {
		tcpRoute := tcpRoute().(*gwv1a2.TCPRoute)

		rm := reports.NewReportMap()
		reporter := reports.NewReporter(&rm)

		reporter.Route(tcpRoute)
		status := rm.BuildRouteStatus(context.Background(), tcpRoute, "gloo-gateway")

		initialTime := meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionAccepted)).LastTransitionTime

		reporter.Route(tcpRoute).ParentRef(&gwv1.ParentReference{
			Name: gwv1.ObjectName("gloo-gateway"),
		}).SetCondition(reports.RouteCondition{
			Type:   gwv1.RouteConditionAccepted,
			Status: metav1.ConditionFalse,
			Reason: gwv1.RouteReasonNoMatchingParent,
		})

		status = rm.BuildRouteStatus(context.Background(), tcpRoute, "gloo-gateway")
		updatedTime := meta.FindStatusCondition(status.Parents[0].Conditions, string(gwv1.RouteConditionAccepted)).LastTransitionTime

		Expect(updatedTime).ToNot(Equal(initialTime))
	})

	It("should filter out duplicate conditions of the same type on TCPRoute", func() {
		tcpRoute := tcpRoute().(*gwv1a2.TCPRoute)

		rm := reports.NewReportMap()
		reporter := reports.NewReporter(&rm)

		reporter.Route(tcpRoute).ParentRef(&gwv1.ParentReference{
			Name: gwv1.ObjectName("gloo-gateway"),
		}).SetCondition(reports.RouteCondition{
			Type:   gwv1.RouteConditionResolvedRefs,
			Status: metav1.ConditionFalse,
			Reason: gwv1.RouteReasonBackendNotFound,
		})

		reporter.Route(tcpRoute).ParentRef(&gwv1.ParentReference{
			Name: gwv1.ObjectName("gloo-gateway"),
		}).SetCondition(reports.RouteCondition{
			Type:   gwv1.RouteConditionResolvedRefs,
			Status: metav1.ConditionFalse,
			Reason: gwv1.RouteReasonBackendNotFound,
		})

		status := rm.BuildRouteStatus(context.Background(), tcpRoute, "gloo-gateway")

		Expect(status).NotTo(BeNil())
		Expect(status.Parents[0].Conditions).To(HaveLen(2))
	})
})

func httpRoute() client.Object {
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRef())
	return route
}

func tcpRoute() client.Object {
	route := &gwv1a2.TCPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRef())
	return route
}

func parentRef() *gwv1.ParentReference {
	return &gwv1.ParentReference{
		Name: "http",
	}
}

func delegateeRoute() client.Object {
	route := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "child-route",
			Namespace: "default",
		},
	}
	route.Spec.CommonRouteSpec.ParentRefs = append(route.Spec.CommonRouteSpec.ParentRefs, *parentRouteRef())
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
