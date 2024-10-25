package listener_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"

	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/listener"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

var _ = Describe("Translator Listener", func() {
	BeforeEach(func() {
	})

	Describe("translates gateway api resources to gloo proxy listeners", func() {
		It("should create a TCP listener with multiple backend references", func() {
			gwListener := gwv1.Listener{
				Name:     "foo-tcp",
				Protocol: gwv1.TCPProtocolType,
				Port:     8080,
			}

			tcpRoute := tcpRoute("test-tcp-route", "default")
			tcpRoute.Spec = gwv1a2.TCPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{
						{
							Name:      gwv1.ObjectName("test-gateway"),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Kind:      ptr.To(gwv1.Kind(wellknown.GatewayKind)),
						},
					},
				},
				Rules: []gwv1a2.TCPRouteRule{
					{
						BackendRefs: []gwv1.BackendRef{
							{
								BackendObjectReference: gwv1.BackendObjectReference{
									Name:      "backend-svc1",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Port:      ptr.To(gwv1.PortNumber(8081)),
								},
								Weight: ptr.To(int32(50)),
							},
							{
								BackendObjectReference: gwv1.BackendObjectReference{
									Name:      "backend-svc2",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Port:      ptr.To(gwv1.PortNumber(8082)),
								},
								Weight: ptr.To(int32(50)),
							},
						},
					},
				},
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			gateway := &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gateway", Namespace: "default"},
			}
			gatewayReporter := reporter.Gateway(gateway)
			listenerReporter := gatewayReporter.Listener(&gwListener)

			routes := []*query.RouteInfo{{Object: tcpRoute}}

			ml := &listener.MergedListeners{}
			ml.AppendTcpListener(gwListener, routes, listenerReporter)

			// Validate that the TCP listener is properly created with multiple backend references
			Expect(ml.Listeners).To(HaveLen(1))
			Expect(ml.Listeners[0].TcpListeners).To(HaveLen(1))
			Expect(ml.Listeners[0].TcpListeners[0].TcpHosts).To(HaveLen(1))
			Expect(ml.Listeners[0].TcpListeners[0].TcpHosts[0].Destination.GetMulti().Destinations).To(HaveLen(2))
		})

		It("should log an error for TCPRoute with missing parent reference", func() {
			gwListener := gwv1.Listener{
				Name:     "foo-tcp",
				Protocol: gwv1.TCPProtocolType,
				Port:     8080,
			}

			tcpRoute := tcpRoute("test-tcp-route", "default")
			tcpRoute.Spec = gwv1a2.TCPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{}, // Empty ParentRefs to trigger the error
				},
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			gateway := &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gateway", Namespace: "default"},
			}
			gatewayReporter := reporter.Gateway(gateway)
			listenerReporter := gatewayReporter.Listener(&gwListener)

			routes := []*query.RouteInfo{{Object: tcpRoute}}

			ml := &listener.MergedListeners{}
			ml.AppendTcpListener(gwListener, routes, listenerReporter)

			// Validate that no listeners are created due to missing ParentRefs
			Expect(ml.Listeners).To(BeEmpty(), "Expected no listeners due to missing ParentRefs")
		})

		It("should handle TCPRoute with empty backend references", func() {
			gwListener := gwv1.Listener{
				Name:     "foo-tcp",
				Protocol: gwv1.TCPProtocolType,
				Port:     8080,
			}

			tcpRoute := tcpRoute("test-empty-backend", "default")
			tcpRoute.Spec = gwv1a2.TCPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{
						{
							Name:      gwv1.ObjectName("test-gateway"),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Kind:      ptr.To(gwv1.Kind(wellknown.GatewayKind)),
						},
					},
				},
				Rules: []gwv1a2.TCPRouteRule{
					{BackendRefs: []gwv1.BackendRef{}}, // Empty BackendRefs
				},
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			gateway := &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gateway", Namespace: "default"},
			}
			gatewayReporter := reporter.Gateway(gateway)
			listenerReporter := gatewayReporter.Listener(&gwListener)

			routes := []*query.RouteInfo{{Object: tcpRoute}}

			ml := &listener.MergedListeners{}
			ml.AppendTcpListener(gwListener, routes, listenerReporter)

			// Validate that no valid listener was created
			Expect(ml.Listeners).To(BeEmpty(), "Expected no listeners due to empty backend references")
		})

		It("should not append a listener for an unsupported protocol", func() {
			gwListener := gwv1.Listener{
				Name:     "foo-unsupported",
				Protocol: gwv1.ProtocolType("UNSUPPORTED"),
				Port:     8080,
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			gateway := &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gateway", Namespace: "default"},
			}
			gatewayReporter := reporter.Gateway(gateway)
			listenerReporter := gatewayReporter.Listener(&gwListener)

			ml := &listener.MergedListeners{}
			ml.AppendTcpListener(gwListener, nil, listenerReporter)

			// Validate that no listeners are created due to the unsupported protocol
			Expect(ml.Listeners).To(BeEmpty(), "Expected no listeners due to unsupported protocol")
		})

		It("should skip routes with invalid parent references and process valid ones", func() {
			gwListener := gwv1.Listener{
				Name:     "foo-tcp",
				Protocol: gwv1.TCPProtocolType,
				Port:     8080,
			}

			// Valid TCP route with a correct parent reference
			validRoute := tcpRoute("valid-tcp-route", "default")
			validRoute.Spec = gwv1a2.TCPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{
						{
							Name:      gwv1.ObjectName("test-gateway"),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Kind:      ptr.To(gwv1.Kind(wellknown.GatewayKind)),
						},
					},
				},
				Rules: []gwv1a2.TCPRouteRule{
					{
						BackendRefs: []gwv1.BackendRef{
							{
								BackendObjectReference: gwv1.BackendObjectReference{
									Name:      "backend-svc1",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Port:      ptr.To(gwv1.PortNumber(8081)),
								},
							},
						},
					},
				},
			}

			// Invalid TCP route with missing parent reference
			invalidRoute := tcpRoute("invalid-tcp-route", "default")
			invalidRoute.Spec = gwv1a2.TCPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{}, // No parent reference provided
				},
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			gateway := &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gateway", Namespace: "default"},
			}
			gatewayReporter := reporter.Gateway(gateway)
			listenerReporter := gatewayReporter.Listener(&gwListener)

			routes := []*query.RouteInfo{
				{Object: validRoute},
				{Object: invalidRoute},
			}

			ml := &listener.MergedListeners{}
			ml.AppendTcpListener(gwListener, routes, listenerReporter)

			// Verify that only the valid route is processed
			Expect(ml.Listeners).To(HaveLen(1))                             // One valid listener
			Expect(ml.Listeners[0].TcpListeners).To(HaveLen(1))             // One TCP listener
			Expect(ml.Listeners[0].TcpListeners[0].TcpHosts).To(HaveLen(1)) // One TCP host
			Expect(ml.Listeners[0].TcpListeners[0].TcpHosts[0].Destination.GetSingle()).NotTo(BeNil())
			Expect(ml.Listeners[0].TcpListeners[0].TcpHosts[0].Destination.GetSingle().GetKube().GetRef().Name).
				To(Equal("backend-svc1"))
		})

		It("should create a TCP listener with a single weighted backend reference", func() {
			gwListener := gwv1.Listener{
				Name:     "foo-tcp",
				Protocol: gwv1.TCPProtocolType,
				Port:     8080,
			}

			tcpRoute := tcpRoute("test-tcp-route", "default")
			tcpRoute.Spec = gwv1a2.TCPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{
						{
							Name:      gwv1.ObjectName("test-gateway"),
							Namespace: ptr.To(gwv1.Namespace("default")),
							Kind:      ptr.To(gwv1.Kind(wellknown.GatewayKind)),
						},
					},
				},
				Rules: []gwv1a2.TCPRouteRule{
					{
						BackendRefs: []gwv1.BackendRef{
							{
								BackendObjectReference: gwv1.BackendObjectReference{
									Name:      "backend-svc1",
									Namespace: ptr.To(gwv1.Namespace("default")),
									Port:      ptr.To(gwv1.PortNumber(8081)),
								},
								Weight: ptr.To(int32(100)),
							},
						},
					},
				},
			}

			rm := reports.NewReportMap()
			reporter := reports.NewReporter(&rm)
			gateway := &gwv1.Gateway{
				ObjectMeta: metav1.ObjectMeta{Name: "test-gateway", Namespace: "default"},
			}
			gatewayReporter := reporter.Gateway(gateway)
			listenerReporter := gatewayReporter.Listener(&gwListener)

			routes := []*query.RouteInfo{{Object: tcpRoute}}

			ml := &listener.MergedListeners{}
			ml.AppendTcpListener(gwListener, routes, listenerReporter)

			// Validate the listener creation
			Expect(ml.Listeners).To(HaveLen(1))
			Expect(ml.Listeners[0].TcpListeners).To(HaveLen(1))
			Expect(ml.Listeners[0].TcpListeners[0].TcpHosts).To(HaveLen(1))

			// Access the destination field properly
			tcpHost := ml.Listeners[0].TcpListeners[0].TcpHosts[0]
			Expect(tcpHost.Destination.GetSingle()).NotTo(BeNil())
			Expect(tcpHost.Destination.GetSingle().GetKube().GetRef().Name).To(Equal("backend-svc1"))
		})
	})
})

func tcpRoute(name, ns string) *gwv1a2.TCPRoute {
	return &gwv1a2.TCPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       wellknown.TCPRouteKind,
			APIVersion: gwv1a2.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}
