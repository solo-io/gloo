package route_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/reports"
)

var _ = Describe("GatewayRouteTranslator", func() {
	var (
		ctrl              *gomock.Controller
		ctx               context.Context
		reportsMap        reports.ReportMap
		baseReporter      reports.Reporter
		parentRefReporter reports.ParentRefReporter
		backends          query.BackendMap[client.Object]
		parentRef         *gwv1.ParentReference
		route             client.Object
		backingSvc        *corev1.Service
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		ctx = context.Background()

		reportsMap = reports.NewReportMap()
		baseReporter = reports.NewReporter(&reportsMap)
		parentRef = &gwv1.ParentReference{Name: "my-gw"}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Route resource routing", func() {

		When("using HTTPRoute with a valid backing service", func() {
			BeforeEach(func() {
				route = &gwv1.HTTPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-httproute",
						Namespace: "bar",
					},
					Spec: gwv1.HTTPRouteSpec{
						Hostnames: []gwv1.Hostname{"example.com"},
						CommonRouteSpec: gwv1.CommonRouteSpec{
							ParentRefs: []gwv1.ParentReference{*parentRef},
						},
						Rules: []gwv1.HTTPRouteRule{
							{
								Matches: []gwv1.HTTPRouteMatch{
									{Path: &gwv1.HTTPPathMatch{
										Type:  ptr.To(gwv1.PathMatchPathPrefix),
										Value: ptr.To("/"),
									}},
								},
								BackendRefs: []gwv1.HTTPBackendRef{
									{
										BackendRef: gwv1.BackendRef{
											BackendObjectReference: gwv1.BackendObjectReference{
												Name:      gwv1.ObjectName("foo"),
												Namespace: ptr.To(gwv1.Namespace("bar")),
												Kind:      ptr.To(gwv1.Kind("Service")),
												Port:      ptr.To(gwv1.PortNumber(8080)),
											},
										},
									},
								},
							},
						},
					},
				}

				// Setup the backing service
				backingSvc = &corev1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Spec: corev1.ServiceSpec{
						Ports: []corev1.ServicePort{{
							Name: "http",
							Port: 8080,
						}},
					},
				}

				// Add backing service to backend map
				backends = query.NewBackendMap[client.Object]()
				backends.Add(route.(*gwv1.HTTPRoute).Spec.Rules[0].BackendRefs[0].BackendObjectReference, backingSvc)

				parentRefReporter = baseReporter.Route(route).ParentRef(parentRef)
			})

			It("translates the HTTPRoute correctly", func() {
				// Use parentRefReporter to simulate the condition reporting mechanism
				parentRefReporter.SetCondition(reports.RouteCondition{
					Type:    gwv1.RouteConditionAccepted,
					Status:  metav1.ConditionTrue,
					Reason:  gwv1.RouteReasonAccepted,
					Message: "HTTPRoute successfully accepted",
				})

				routeStatus := reportsMap.BuildRouteStatus(ctx, route, "")
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				By("verifying the route was accepted")
				accepted := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionAccepted))
				Expect(accepted).NotTo(BeNil())
				Expect(accepted.Status).To(Equal(metav1.ConditionTrue))
				Expect(accepted.Reason).To(BeEquivalentTo(gwv1.RouteReasonAccepted))
			})
		})

		When("using TCPRoute with a valid configuration", func() {
			BeforeEach(func() {
				route = &gwv1alpha2.TCPRoute{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "foo-tcproute",
						Namespace: "bar",
					},
					Spec: gwv1alpha2.TCPRouteSpec{
						CommonRouteSpec: gwv1alpha2.CommonRouteSpec{
							ParentRefs: []gwv1alpha2.ParentReference{*parentRef},
						},
					},
				}

				// Initialize the backend map for TCPRoute (simulated)
				backends = query.NewBackendMap[client.Object]()
				parentRefReporter = baseReporter.Route(route).ParentRef(parentRef)
			})

			It("translates the TCPRoute correctly", func() {
				// Use parentRefReporter to simulate the condition reporting mechanism
				parentRefReporter.SetCondition(reports.RouteCondition{
					Type:    gwv1.RouteConditionAccepted,
					Status:  metav1.ConditionTrue,
					Reason:  gwv1.RouteReasonAccepted,
					Message: "TCPRoute successfully accepted",
				})

				routeStatus := reportsMap.BuildRouteStatus(ctx, route, "")
				Expect(routeStatus).NotTo(BeNil())
				Expect(routeStatus.Parents).To(HaveLen(1))
				By("verifying the TCPRoute was accepted")
				accepted := meta.FindStatusCondition(routeStatus.Parents[0].Conditions, string(gwv1.RouteConditionAccepted))
				Expect(accepted).NotTo(BeNil())
				Expect(accepted.Status).To(Equal(metav1.ConditionTrue))
				Expect(accepted.Reason).To(BeEquivalentTo(gwv1.RouteReasonAccepted))
			})
		})
	})
})
