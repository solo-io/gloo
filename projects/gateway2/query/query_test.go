package query_test

import (
	"context"

	"github.com/solo-io/gloo/pkg/schemes"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"
	apiv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	apiv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"

	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/mock_queries.go -package mocks github.com/solo-io/gloo/projects/gateway2/query GatewayQueries

var _ = Describe("Query", func() {
	var (
		scheme  *runtime.Scheme
		builder *fake.ClientBuilder
	)

	tofrom := func(o client.Object) query.From {
		return query.FromObject{Scheme: scheme, Object: o}
	}

	BeforeEach(func() {
		scheme = schemes.TestingScheme()
		builder = fake.NewClientBuilder().WithScheme(scheme)
		err := query.IterateIndices(func(o client.Object, f string, fun client.IndexerFunc) error {
			builder.WithIndex(o, f, fun)
			return nil
		})
		Expect(err).NotTo(HaveOccurred())
	})
	Describe("GetBackendForRef", func() {
		It("should get service from same namespace", func() {
			fakeClient := fake.NewFakeClient(svc("default"))

			gq := query.NewData(fakeClient, scheme)
			ref := &apiv1.BackendObjectReference{
				Name: "foo",
			}

			backend, err := gq.GetBackendForRef(context.Background(), tofrom(httpRoute()), ref)
			Expect(err).NotTo(HaveOccurred())
			Expect(backend).NotTo(BeNil())
			Expect(backend.GetName()).To(Equal("foo"))
			Expect(backend.GetNamespace()).To(Equal("default"))
		})

		It("should get service from different ns if we have a ref grant", func() {
			rg := refGrant()
			fakeClient := builder.WithObjects(svc("default2"), rg).Build()
			gq := query.NewData(fakeClient, scheme)
			ref := &apiv1.BackendObjectReference{
				Name:      "foo",
				Namespace: nsptr("default2"),
			}

			backend, err := gq.GetBackendForRef(context.Background(), tofrom(httpRoute()), ref)
			Expect(err).NotTo(HaveOccurred())
			Expect(backend).NotTo(BeNil())
			Expect(backend.GetName()).To(Equal("foo"))
			Expect(backend.GetNamespace()).To(Equal("default2"))
		})

		It("should fail with service not found if we have a ref grant", func() {
			rg := refGrant()
			fakeClient := builder.WithObjects(rg).Build()
			gq := query.NewData(fakeClient, scheme)
			ref := &apiv1.BackendObjectReference{
				Name:      "foo",
				Namespace: nsptr("default2"),
			}
			backend, err := gq.GetBackendForRef(context.Background(), tofrom(httpRoute()), ref)
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
			Expect(backend).To(BeNil())
		})

		It("should fail getting a service with ref grant with wrong from", func() {
			ref := &apiv1.BackendObjectReference{
				Name:      "foo",
				Namespace: nsptr("default2"),
			}
			rg := &apiv1beta1.ReferenceGrant{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default2",
					Name:      "foo",
				},
				Spec: apiv1beta1.ReferenceGrantSpec{
					From: []apiv1beta1.ReferenceGrantFrom{
						{
							Group:     apiv1.Group("gateway.networking.k8s.io"),
							Kind:      apiv1.Kind("NotGateway"),
							Namespace: apiv1.Namespace("default"),
						},
						{
							Group:     apiv1.Group("gateway.networking.k8s.io"),
							Kind:      apiv1.Kind("Gateway"),
							Namespace: apiv1.Namespace("default2"),
						},
					},
					To: []apiv1beta1.ReferenceGrantTo{
						{
							Group: apiv1.Group("core"),
							Kind:  apiv1.Kind("Service"),
						},
					},
				},
			}
			fakeClient := builder.WithObjects(rg, svc("default2")).Build()

			gq := query.NewData(fakeClient, scheme)
			backend, err := gq.GetBackendForRef(context.Background(), tofrom(httpRoute()), ref)
			Expect(err).To(MatchError(query.ErrMissingReferenceGrant))
			Expect(backend).To(BeNil())
		})

		It("should fail getting a service with no ref grant", func() {
			fakeClient := builder.WithObjects(svc("default3")).Build()
			gq := query.NewData(fakeClient, scheme)
			ref := &apiv1.BackendObjectReference{
				Name:      "foo",
				Namespace: nsptr("default3"),
			}

			backend, err := gq.GetBackendForRef(context.Background(), tofrom(httpRoute()), ref)
			Expect(err).To(MatchError(query.ErrMissingReferenceGrant))
			Expect(backend).To(BeNil())
		})

		It("should fail getting a service with ref grant in wrong ns", func() {
			rg := refGrant()
			fakeClient := builder.WithObjects(svc("default3"), rg).Build()

			gq := query.NewData(fakeClient, scheme)
			ref := &apiv1.BackendObjectReference{
				Name:      "foo",
				Namespace: nsptr("default3"),
			}
			backend, err := gq.GetBackendForRef(context.Background(), tofrom(httpRoute()), ref)
			Expect(err).To(MatchError(query.ErrMissingReferenceGrant))
			Expect(backend).To(BeNil())
		})
	})

	Describe("GetSecretRef", func() {
		It("should get secret from different ns if we have a ref grant", func() {
			rg := refGrantSecret()
			fakeClient := builder.WithObjects(secret("default2"), rg).Build()
			gq := query.NewData(fakeClient, scheme)
			ref := apiv1.SecretObjectReference{
				Name:      "foo",
				Namespace: nsptr("default2"),
			}
			backend, err := gq.GetSecretForRef(context.Background(), tofrom(gw()), ref)
			Expect(err).NotTo(HaveOccurred())
			Expect(backend).NotTo(BeNil())
			Expect(backend.GetName()).To(Equal("foo"))
			Expect(backend.GetNamespace()).To(Equal("default2"))
		})
	})

	Describe("Get Routes", func() {
		It("should get http routes for listener", func() {
			gwWithListener := gw()
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
				},
			}
			hr := httpRoute()
			hr.Spec.ParentRefs = []apiv1.ParentReference{
				{
					Name: "test",
				},
			}

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(BeEmpty())
			Expect(routes.ListenerResults["foo"].Error).NotTo(HaveOccurred())
			Expect(routes.ListenerResults["foo"].Routes).To(HaveLen(1))
		})

		It("should get http routes in other ns for listener", func() {
			gwWithListener := gw()
			all := apiv1.NamespacesFromAll
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Namespaces: &apiv1.RouteNamespaces{
							From: &all,
						},
					},
				},
			}
			hr := httpRoute()
			hr.Namespace = "default2"
			hr.Spec.ParentRefs = []apiv1.ParentReference{
				{
					Name:      "test",
					Namespace: nsptr("default"),
				},
			}

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(BeEmpty())
			Expect(routes.ListenerResults["foo"].Error).NotTo(HaveOccurred())
			Expect(routes.ListenerResults["foo"].Routes).To(HaveLen(1))
		})

		It("should error with invalid label selector", func() {
			gwWithListener := gw()
			selector := apiv1.NamespacesFromSelector
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Namespaces: &apiv1.RouteNamespaces{
							From:     &selector,
							Selector: nil,
						},
					},
				},
			}
			hr := httpRoute()
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.ListenerResults["foo"].Error).To(MatchError("selector must be set"))
		})

		It("should error when listeners do not allow route", func() {
			gwWithListener := gw()
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Kinds: []apiv1.RouteGroupKind{{Kind: "FakeKind"}},
					},
				},
				{
					Name:     "foo2",
					Protocol: apiv1.HTTPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Kinds: []apiv1.RouteGroupKind{{Kind: "FakeKind2"}},
					},
				},
			}
			hr := httpRoute()
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors[0].Error.E).To(MatchError(query.ErrNotAllowedByListeners))
			Expect(routes.RouteErrors[0].Error.Reason).To(Equal(apiv1.RouteReasonNotAllowedByListeners))
			Expect(routes.RouteErrors[0].ParentRef).To(Equal(hr.Spec.ParentRefs[0]))
		})

		It("should NOT error when one listener allows route", func() {
			gwWithListener := gw()
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Kinds: []apiv1.RouteGroupKind{{Kind: "FakeKind"}},
					},
				},
				{
					Name:     "foo2",
					Protocol: apiv1.HTTPProtocolType,
				},
			}
			hr := httpRoute()
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(BeEmpty())
			Expect(routes.ListenerResults["foo2"].Routes).To(HaveLen(1))
			Expect(routes.ListenerResults["foo2"].Error).NotTo(HaveOccurred())
			Expect(routes.ListenerResults["foo"].Routes).To(BeEmpty())
			Expect(routes.ListenerResults["foo"].Error).NotTo(HaveOccurred())
		})

		It("should error when listeners don't match route", func() {
			gwWithListener := gw()
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
				},
				{
					Name:     "bar",
					Protocol: apiv1.HTTPProtocolType,
					Port:     81,
				},
			}
			hr := httpRoute()
			var port apiv1.PortNumber = 1234
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
				Port: &port,
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors[0].Error.E).To(MatchError(query.ErrNoMatchingParent))
			Expect(routes.RouteErrors[0].Error.Reason).To(Equal(apiv1.RouteReasonNoMatchingParent))
			Expect(routes.RouteErrors[0].ParentRef).To(Equal(hr.Spec.ParentRefs[0]))
		})

		It("should NOT error when one listener match route", func() {
			gwWithListener := gw()
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
				},
				{
					Name:     "foo2",
					Protocol: apiv1.HTTPProtocolType,
					Port:     81,
				},
			}
			hr := httpRoute()
			var port apiv1.PortNumber = 81
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
				Port: &port,
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(BeEmpty())
			Expect(routes.ListenerResults["foo2"].Routes).To(HaveLen(1))
			Expect(routes.ListenerResults["foo"].Routes).To(BeEmpty())
		})

		It("should error when listeners hostnames don't intersect", func() {
			gwWithListener := gw()
			var hostname apiv1.Hostname = "foo.com"
			var hostname2 apiv1.Hostname = "foo2.com"
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
					Hostname: &hostname,
				},
				{
					Name:     "foo2",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
					Hostname: &hostname2,
				},
			}
			hr := httpRoute()
			hr.Spec.Hostnames = append(hr.Spec.Hostnames, "bar.com")
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors[0].Error.E).To(MatchError(query.ErrNoMatchingListenerHostname))
			Expect(routes.RouteErrors[0].Error.Reason).To(Equal(apiv1.RouteReasonNoMatchingListenerHostname))
			Expect(routes.RouteErrors[0].ParentRef).To(Equal(hr.Spec.ParentRefs[0]))
		})

		It("should NOT error when one listener hostname do intersect", func() {
			gwWithListener := gw()
			var hostname apiv1.Hostname = "foo.com"
			var hostname2 apiv1.Hostname = "bar.com"
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
					Hostname: &hostname,
				},
				{
					Name:     "foo2",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
					Hostname: &hostname2,
				},
			}
			hr := httpRoute()
			hr.Spec.Hostnames = append(hr.Spec.Hostnames, "bar.com")
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(BeEmpty())
			Expect(routes.ListenerResults["foo2"].Routes).To(HaveLen(1))
			Expect(routes.ListenerResults["foo"].Routes).To(BeEmpty())
		})

		It("should error for one parent ref but not the other", func() {
			gwWithListener := gw()
			var hostname apiv1.Hostname = "foo.com"
			gwWithListener.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo",
					Protocol: apiv1.HTTPProtocolType,
					Port:     80,
					Hostname: &hostname,
				},
			}
			hr := httpRoute()
			var badPort apiv1.PortNumber = 81
			hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
				Port: &badPort,
			}, apiv1.ParentReference{
				Name: apiv1.ObjectName(gwWithListener.Name),
			})

			fakeClient := builder.WithObjects(hr).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(HaveLen(1))
			Expect(routes.ListenerResults["foo"].Routes).To(HaveLen(1))
			Expect(routes.ListenerResults["foo"].Routes[0].ParentRef).To(Equal(apiv1.ParentReference{
				Name: hr.Spec.ParentRefs[1].Name,
			}))
			Expect(routes.RouteErrors[0].Error.E).To(MatchError(query.ErrNoMatchingParent))
			Expect(routes.RouteErrors[0].Error.Reason).To(Equal(apiv1.RouteReasonNoMatchingParent))
			Expect(routes.RouteErrors[0].ParentRef).To(Equal(hr.Spec.ParentRefs[0]))
		})

		Context("test host intersection", func() {
			expectHostnamesToMatch := func(lh string, rh []string, expectedHostnames ...string) {
				gwWithListener := gw()
				gwWithListener.Spec.Listeners = []apiv1.Listener{
					{
						Name:     "foo",
						Protocol: apiv1.HTTPProtocolType,
					},
				}
				if lh != "" {
					h := apiv1.Hostname(lh)
					gwWithListener.Spec.Listeners[0].Hostname = &h

				}

				hr := httpRoute()
				for _, h := range rh {
					hr.Spec.Hostnames = append(hr.Spec.Hostnames, apiv1.Hostname(h))
				}
				hr.Spec.ParentRefs = append(hr.Spec.ParentRefs, apiv1.ParentReference{
					Name: apiv1.ObjectName(gwWithListener.Name),
				})

				fakeClient := builder.WithObjects(hr).Build()
				gq := query.NewData(fakeClient, scheme)
				routes, err := gq.GetRoutesForGateway(context.Background(), gwWithListener)

				Expect(err).NotTo(HaveOccurred())
				if expectedHostnames == nil {
					expectedHostnames = []string{}
				}
				Expect(routes.ListenerResults["foo"].Routes[0].Hostnames()).To(Equal(expectedHostnames))
			}

			It("should work with identical names", func() {
				expectHostnamesToMatch("foo.com", []string{"foo.com"}, "foo.com")
			})
			It("should work with specific listeners and prefix http", func() {
				expectHostnamesToMatch("bar.foo.com", []string{"*.foo.com", "foo.com", "example.com"}, "bar.foo.com")
			})
			It("should work with prefix listeners and specific http", func() {
				expectHostnamesToMatch("*.foo.com", []string{"bar.foo.com", "foo.com", "far.foo.com", "blah.com"}, "bar.foo.com", "far.foo.com")
			})
			It("should work with catch all listener hostname", func() {
				expectHostnamesToMatch("", []string{"foo.com", "blah.com"}, "foo.com", "blah.com")
			})
			It("should work with catch all http hostname", func() {
				expectHostnamesToMatch("foo.com", nil, "foo.com")
			})
			It("should work with listener prefix and catch all http hostname", func() {
				expectHostnamesToMatch("*.foo.com", nil, "*.foo.com")
			})
			It("should work with double catch all", func() {
				expectHostnamesToMatch("", nil)
			})
		})

		It("should match TCPRoutes for Listener", func() {
			gw := gw()
			gw.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo-tcp",
					Protocol: apiv1.TCPProtocolType,
				},
			}

			tcpRoute := tcpRoute("test-tcp-route", gw.Namespace)
			tcpRoute.Spec = apiv1a2.TCPRouteSpec{
				CommonRouteSpec: apiv1.CommonRouteSpec{
					ParentRefs: []apiv1.ParentReference{
						{
							Name: apiv1.ObjectName(gw.Name),
						},
					},
				},
			}

			fakeClient := builder.WithObjects(tcpRoute).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gw)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.ListenerResults[string(gw.Spec.Listeners[0].Name)].Routes).To(HaveLen(1))
			Expect(routes.ListenerResults[string(gw.Spec.Listeners[0].Name)].Error).NotTo(HaveOccurred())
		})

		It("should get TCPRoutes in other namespace for listener", func() {
			gw := gw()
			gw.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo-tcp",
					Protocol: apiv1.TCPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Namespaces: &apiv1.RouteNamespaces{
							From: ptr.To(apiv1.NamespacesFromAll),
						},
					},
				},
			}

			tcpRoute := tcpRoute("test-tcp-route", "other-ns")
			tcpRoute.Spec = apiv1a2.TCPRouteSpec{
				CommonRouteSpec: apiv1.CommonRouteSpec{
					ParentRefs: []apiv1.ParentReference{
						{
							Name:      apiv1.ObjectName(gw.Name),
							Namespace: ptr.To(apiv1.Namespace(gw.Namespace)),
						},
					},
				},
			}

			fakeClient := builder.WithObjects(tcpRoute).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gw)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.ListenerResults["foo-tcp"].Error).NotTo(HaveOccurred())
			Expect(routes.ListenerResults["foo-tcp"].Routes).To(HaveLen(1))
		})

		It("should error when listeners don't match TCPRoute", func() {
			gw := gw()
			gw.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo-tcp",
					Protocol: apiv1.TCPProtocolType,
					Port:     8080,
				},
				{
					Name:     "bar-tcp",
					Protocol: apiv1.TCPProtocolType,
					Port:     8081,
				},
			}

			tcpRoute := tcpRoute("test-tcp-route", gw.Namespace)
			var badPort apiv1.PortNumber = 9999
			tcpRoute.Spec = apiv1a2.TCPRouteSpec{
				CommonRouteSpec: apiv1.CommonRouteSpec{
					ParentRefs: []apiv1.ParentReference{
						{
							Name: apiv1.ObjectName(gw.Name),
							Port: &badPort,
						},
					},
				},
			}

			fakeClient := builder.WithObjects(tcpRoute).Build()
			gq := query.NewData(fakeClient, scheme)
			routes, err := gq.GetRoutesForGateway(context.Background(), gw)

			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(HaveLen(1))
			Expect(routes.RouteErrors[0].Error.E).To(MatchError(query.ErrNoMatchingParent))
			Expect(routes.RouteErrors[0].Error.Reason).To(Equal(apiv1.RouteReasonNoMatchingParent))
			Expect(routes.RouteErrors[0].ParentRef).To(Equal(tcpRoute.Spec.ParentRefs[0]))
		})

		It("should error when listener does not allow TCPRoute kind", func() {
			gw := gw()
			gw.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo-tcp",
					Protocol: apiv1.TCPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Kinds: []apiv1.RouteGroupKind{{Kind: "FakeKind"}},
					},
				},
			}

			tcpRoute := tcpRoute("test-tcp-route", gw.Namespace)
			tcpRoute.Spec = apiv1a2.TCPRouteSpec{
				CommonRouteSpec: apiv1.CommonRouteSpec{
					ParentRefs: []apiv1.ParentReference{
						{
							Name: apiv1.ObjectName(gw.Name),
						},
					},
				},
			}

			fakeClient := builder.WithObjects(tcpRoute).Build()
			gq := query.NewData(fakeClient, scheme)

			routes, err := gq.GetRoutesForGateway(context.Background(), gw)
			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(HaveLen(1))
			Expect(routes.RouteErrors[0].Error.E).To(MatchError(query.ErrNotAllowedByListeners))
		})

		It("should allow TCPRoute for one listener", func() {
			gw := gw()
			gw.Spec.Listeners = []apiv1.Listener{
				{
					Name:     "foo-tcp",
					Protocol: apiv1.TCPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Kinds: []apiv1.RouteGroupKind{{Kind: wellknown.TCPRouteKind}},
					},
				},
				{
					Name:     "bar",
					Protocol: apiv1.TCPProtocolType,
					AllowedRoutes: &apiv1.AllowedRoutes{
						Kinds: []apiv1.RouteGroupKind{{Kind: "FakeKind"}},
					},
				},
			}

			tcpRoute := tcpRoute("test-tcp-route", gw.Namespace)
			tcpRoute.Spec = apiv1a2.TCPRouteSpec{
				CommonRouteSpec: apiv1.CommonRouteSpec{
					ParentRefs: []apiv1.ParentReference{
						{
							Name: apiv1.ObjectName(gw.Name),
						},
					},
				},
			}

			fakeClient := builder.WithObjects(tcpRoute).Build()
			gq := query.NewData(fakeClient, scheme)

			routes, err := gq.GetRoutesForGateway(context.Background(), gw)
			Expect(err).NotTo(HaveOccurred())
			Expect(routes.RouteErrors).To(BeEmpty())
			Expect(routes.ListenerResults["foo-tcp"].Routes).To(HaveLen(1))
			Expect(routes.ListenerResults["bar"].Routes).To(BeEmpty())
		})
	})
})

func refGrantSecret() *apiv1beta1.ReferenceGrant {
	return &apiv1beta1.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default2",
			Name:      "foo",
		},
		Spec: apiv1beta1.ReferenceGrantSpec{
			From: []apiv1beta1.ReferenceGrantFrom{
				{
					Group:     apiv1.Group("gateway.networking.k8s.io"),
					Kind:      apiv1.Kind("Gateway"),
					Namespace: apiv1.Namespace("default"),
				},
			},
			To: []apiv1beta1.ReferenceGrantTo{
				{
					Group: apiv1.Group("core"),
					Kind:  apiv1.Kind("Secret"),
				},
			},
		},
	}
}

func refGrant() *apiv1beta1.ReferenceGrant {
	return &apiv1beta1.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default2",
			Name:      "foo",
		},
		Spec: apiv1beta1.ReferenceGrantSpec{
			From: []apiv1beta1.ReferenceGrantFrom{
				{
					Group:     apiv1.Group("gateway.networking.k8s.io"),
					Kind:      apiv1.Kind("HTTPRoute"),
					Namespace: apiv1.Namespace("default"),
				},
			},
			To: []apiv1beta1.ReferenceGrantTo{
				{
					Group: apiv1.Group("core"),
					Kind:  apiv1.Kind("Service"),
				},
			},
		},
	}
}

func httpRoute() *apiv1.HTTPRoute {
	return &apiv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       wellknown.HTTPRouteKind,
			APIVersion: apiv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test",
		},
	}
}

func gw() *apiv1.Gateway {
	return &apiv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "test",
		},
	}
}

func secret(ns string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "foo",
		},
	}
}

func svc(ns string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "foo",
		},
	}
}

func tcpRoute(name, ns string) *apiv1a2.TCPRoute {
	return &apiv1a2.TCPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       wellknown.TCPRouteKind,
			APIVersion: apiv1a2.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}
}

func nsptr(s string) *apiv1.Namespace {
	var ns apiv1.Namespace = apiv1.Namespace(s)
	return &ns
}
