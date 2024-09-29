package kubernetes

import (
	"context"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/projects/gloo/constants"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	kubev1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	mock_kubernetes "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes/mocks"
	mock_cache "github.com/solo-io/gloo/test/mocks/cache"
)

var _ = Describe("EDS", func() {
	var (
		ctx               context.Context
		controller        *gomock.Controller
		mockCache         *mock_cache.MockKubeCoreCache
		mockSharedFactory *mock_kubernetes.MockKubePluginSharedFactory
	)
	BeforeEach(func() {
		ctx = context.Background()
		ctx = settingsutil.WithSettings(ctx, &gloov1.Settings{WatchNamespaces: []string{"foo"}})

		controller = gomock.NewController(GinkgoT())
		mockCache = mock_cache.NewMockKubeCoreCache(controller)
		mockSharedFactory = mock_kubernetes.NewMockKubePluginSharedFactory(controller)
	})
	AfterEach(func() {
		controller.Finish()
	})

	Context("EDS watcher", func() {
		It("should ignore upstreams in non-watched namesapces", func() {
			up := gloov1.NewUpstream("foo", "name")
			up.UpstreamType = &gloov1.Upstream_Kube{
				Kube: &kubev1.UpstreamSpec{
					ServiceName:      "name",
					ServiceNamespace: "bar",
				},
			}
			upstreamsToTrack := gloov1.UpstreamList{up}

			mockCache.EXPECT().NamespacedServiceLister("bar").Return(nil)

			watcher, err := newEndpointWatcherForUpstreams(func([]string) KubePluginSharedFactory {
				return mockSharedFactory
			}, mockCache,
				upstreamsToTrack,
				clients.WatchOpts{Ctx: ctx},
				nil,
			)
			Expect(err).NotTo(HaveOccurred())
			watcher.List("foo", clients.ListOpts{Ctx: ctx})
			Expect(func() {}).NotTo(Panic())
		})

		It("should default to watchNamespaces if no upstreams exist", func() {
			watchNamespaces := []string{"gloo-system"}
			_, err := newEndpointWatcherForUpstreams(func(namespaces []string) KubePluginSharedFactory {
				Expect(namespaces).To(Equal(watchNamespaces))
				return mockSharedFactory
			},
				mockCache, gloov1.UpstreamList{},
				clients.WatchOpts{Ctx: ctx}, &gloov1.Settings{WatchNamespaces: watchNamespaces},
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	// TODO: Add test cases for ports filtering. This might include servicePort being nil
	// if that's even supported? Need to still investigate the right logic based on how
	// we handle regular Endpoint CRs.
	Context("EndpointSlices", func() {
		var (
			eps         *discoveryv1.EndpointSlice
			epList      []*corev1.Endpoints
			endpoints   []*gloov1.Endpoint
			warnsToLog  []string
			errorsToLog []string
		)
		JustBeforeEach(func() {
			// common test setup used by multiple test cases. each test will
			// populate the eps with different values, so we need to set it up
			// with a JustBeforeEach block to ensure the correct values are used.
			writeNamespace := "foo"
			up := gloov1.NewUpstream(writeNamespace, "name")
			up.UpstreamType = &gloov1.Upstream_Kube{
				Kube: &kubev1.UpstreamSpec{
					ServiceName:      "bar",
					ServiceNamespace: "foo",
					ServicePort:      9080,
					Selector:         map[string]string{"app": "bar"},
				},
			}

			endpoints, warnsToLog, errorsToLog = filterEndpoints(
				writeNamespace,
				epList,
				[]*discoveryv1.EndpointSlice{eps},
				[]*corev1.Service{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "foo",
					},
					Spec: corev1.ServiceSpec{
						Selector: map[string]string{"app": "bar"},
						Ports: []corev1.ServicePort{
							{
								Name:     "http",
								Port:     9080,
								Protocol: "TCP",
								TargetPort: intstr.IntOrString{
									Type:   intstr.Int,
									IntVal: 9080,
								},
							},
						},
					},
				}},
				[]*corev1.Pod{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar-7d4d7c7b4b-4z5zv",
						Namespace: "foo",
						Labels: map[string]string{
							"app": "bar",
						},
					},
					Spec: corev1.PodSpec{},
				}},
				map[*core.ResourceRef]*kubeplugin.UpstreamSpec{
					up.Metadata.Ref(): up.GetKube(),
				},
			)
		})

		When("matching EndpointSlices are present", func() {
			BeforeEach(func() {
				eps = &discoveryv1.EndpointSlice{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "foo",
						Labels: map[string]string{
							discoveryv1.LabelServiceName: "bar",
						},
					},
					AddressType: discoveryv1.AddressTypeIPv4,
					Ports: []discoveryv1.EndpointPort{{
						Name:     ptr.To("http"),
						Port:     ptr.To(int32(6443)),
						Protocol: ptr.To(corev1.ProtocolTCP),
					}},
					Endpoints: []discoveryv1.Endpoint{{
						Addresses: []string{"172.19.0.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready: ptr.To(true),
						},
					}},
				}
			})
			It("should translate EDS for matching EndpointSlices", func() {
				Expect(endpoints).To(HaveLen(1), "expected to have 1 endpoint")
				Expect(warnsToLog).To(BeEmpty(), "expected no warnings")
				Expect(errorsToLog).To(BeEmpty(), "expected no errors")

				endpoint := endpoints[0]
				Expect(endpoint.Address).To(Equal(eps.Endpoints[0].Addresses[0]))
				Expect(endpoint.Port).To(BeEquivalentTo(*eps.Ports[0].Port))
			})
		})

		When("an EndpointSlices does not configure the kubernetes.io/service-name label", func() {
			BeforeEach(func() {
				eps = &discoveryv1.EndpointSlice{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "foo",
						Labels:    map[string]string{},
					},
					AddressType: discoveryv1.AddressTypeIPv4,
					Ports: []discoveryv1.EndpointPort{{
						Name:     ptr.To("http"),
						Port:     ptr.To(int32(6443)),
						Protocol: ptr.To(corev1.ProtocolTCP),
					}},
					Endpoints: []discoveryv1.Endpoint{{
						Addresses: []string{"172.19.0.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready: ptr.To(false),
						},
					}},
				}
			})

			It("should not translate EDS", func() {
				Expect(endpoints).To(BeEmpty(), "expected to have no endpoints")
				Expect(warnsToLog).To(BeEmpty(), "expected no warnings")
				Expect(errorsToLog).To(BeEmpty(), "expected no errors")
			})
		})

		When("a matching EndpointSlices is not ready", func() {
			BeforeEach(func() {
				eps = &discoveryv1.EndpointSlice{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "bar",
						Namespace: "foo",
						Labels:    map[string]string{discoveryv1.LabelServiceName: "bar"},
					},
					AddressType: discoveryv1.AddressTypeIPv4,
					Ports: []discoveryv1.EndpointPort{{
						Name:     ptr.To("http"),
						Port:     ptr.To(int32(6443)),
						Protocol: ptr.To(corev1.ProtocolTCP),
					}},
					Endpoints: []discoveryv1.Endpoint{{
						Addresses: []string{"172.19.0.2"},
						Conditions: discoveryv1.EndpointConditions{
							Ready: ptr.To(false),
						},
					}},
				}
			})

			It("should not translate EDS for a non-ready endpoint", func() {
				Expect(endpoints).To(BeEmpty(), "expected to have no endpoints")
				Expect(warnsToLog).To(BeEmpty(), "expected no warnings")
				Expect(errorsToLog).To(BeEmpty(), "expected no errors")
			})
		})

		// TODO: add test for the port filtering too.

		// TODO: fix test.
		// When("matching Endpoints and EndpointSlices are present", func() {
		// 	BeforeEach(func() {
		// 		eps = &discoveryv1.EndpointSlice{
		// 			ObjectMeta: metav1.ObjectMeta{
		// 				Name:      "bar-7d4d7c7b4b-4z5zv",
		// 				Namespace: "foo",
		// 				Labels: map[string]string{
		// 					discoveryv1.LabelServiceName: "bar",
		// 				},
		// 			},
		// 			AddressType: discoveryv1.AddressTypeIPv4,
		// 			Ports: []discoveryv1.EndpointPort{{
		// 				Name:     ptr.To("http"),
		// 				Port:     ptr.To(int32(6443)),
		// 				Protocol: ptr.To(corev1.ProtocolTCP),
		// 			}},
		// 			Endpoints: []discoveryv1.Endpoint{{
		// 				Addresses: []string{""},
		// 			}},
		// 		}
		// 		epList = []*corev1.Endpoints{{
		// 			ObjectMeta: metav1.ObjectMeta{
		// 				Name:      "bar",
		// 				Namespace: "foo",
		// 			},
		// 			Subsets: []corev1.EndpointSubset{{
		// 				Ports: []corev1.EndpointPort{{
		// 					Port:     9080,
		// 					Name:     "http",
		// 					Protocol: "TCP",
		// 				}},
		// 				Addresses: []corev1.EndpointAddress{{
		// 					IP: "10.244.0.14",
		// 					TargetRef: &corev1.ObjectReference{
		// 						Kind:      "Pod",
		// 						Name:      "bar-7d4d7c7b4b-4z5zv",
		// 						Namespace: "foo",
		// 					},
		// 				}},
		// 			}},
		// 		}}
		// 	})

		// 	It("should translate EDS for matching Endpoints and EndpointSlices", func() {
		// 		Expect(endpoints).To(HaveLen(2), "expected to have 2 endpoint")
		// 		Expect(warnsToLog).To(BeEmpty(), "expected no warnings")
		// 		Expect(errorsToLog).To(BeEmpty(), "expected no errors")

		// 		endpoint := endpoints[0]
		// 		Expect(endpoint.Address).To(Equal(eps.Endpoints[0].Addresses[0]))
		// 		Expect(endpoint.Port).To(BeEquivalentTo(*eps.Ports[0].Port))
		// 	})
		// })
	})

	Context("Istio integration", func() {
		// TODO(tim): add test for EndpointSlices.

		It("isIstioInjectionEnabled should respond correctly to ENABLE_ISTIO_SIDECAR_ON_GATEWAY env var", func() {
			os.Setenv(constants.IstioInjectionEnabled, "true")
			istioEnabled, warnsToLog := isIstioInjectionEnabled()
			Expect(istioEnabled).To(BeTrue())
			Expect(warnsToLog).To(HaveLen(1), "expected to have 1 warning")
			Expect(warnsToLog).To(ContainElements(enableIstioSidecarOnGatewayDeprecatedWarning), "expected deprecation warning for enableIstioSidecarOnGateway")

			os.Setenv(constants.IstioInjectionEnabled, "TRUE")
			istioEnabled, warnsToLog = isIstioInjectionEnabled()
			Expect(istioEnabled).To(BeTrue())
			Expect(warnsToLog).To(HaveLen(1), "expected to have 1 warning")
			Expect(warnsToLog).To(ContainElements(enableIstioSidecarOnGatewayDeprecatedWarning), "expected deprecation warning for enableIstioSidecarOnGateway")

			os.Unsetenv(constants.IstioInjectionEnabled)
			istioEnabled, warnsToLog = isIstioInjectionEnabled()
			Expect(istioEnabled).To(BeFalse())
			Expect(warnsToLog).To(BeEmpty(), "expected to have no warning")

			os.Setenv(constants.IstioInjectionEnabled, "false")
			istioEnabled, warnsToLog = isIstioInjectionEnabled()
			Expect(istioEnabled).To(BeFalse())
			Expect(warnsToLog).To(BeEmpty(), "expected to have no warning")
		})

		It("should translate EDS metadata", func() {
			writeNamespace := "foo"
			up := gloov1.NewUpstream(writeNamespace, "name")
			up.UpstreamType = &gloov1.Upstream_Kube{
				Kube: &kubev1.UpstreamSpec{
					ServiceName:      "bar",
					ServiceNamespace: "foo",
					ServicePort:      9080,
					Selector:         map[string]string{"app": "bar"},
				},
			}

			endpoints, warnsToLog, errorsToLog := filterEndpoints(
				writeNamespace,
				[]*corev1.Endpoints{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "bar",
							Namespace: "foo",
						},
						Subsets: []corev1.EndpointSubset{
							{
								Ports: []corev1.EndpointPort{
									{
										Port:     9080,
										Name:     "http",
										Protocol: "TCP",
									},
								},
								Addresses: []corev1.EndpointAddress{
									{
										IP: "10.244.0.14",
										TargetRef: &corev1.ObjectReference{
											Kind:      "Pod",
											Name:      "bar-7d4d7c7b4b-4z5zv",
											Namespace: "foo",
										},
									},
								},
							},
						},
					},
				},
				[]*discoveryv1.EndpointSlice{},
				[]*corev1.Service{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "bar",
							Namespace: "foo",
						},
						Spec: corev1.ServiceSpec{
							Selector: map[string]string{"app": "bar"},
							Ports: []corev1.ServicePort{
								{
									Name:     "http",
									Port:     9080,
									Protocol: "TCP",
									TargetPort: intstr.IntOrString{
										Type:   intstr.Int,
										IntVal: 9080,
									},
								},
							},
						},
					},
				},
				[]*corev1.Pod{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "bar-7d4d7c7b4b-4z5zv",
							Namespace: "foo",
							Labels: map[string]string{
								"app":                       "bar",
								"security.istio.io/tlsMode": "istio",
							},
						},
						Spec: corev1.PodSpec{},
					},
				},
				map[*core.ResourceRef]*kubeplugin.UpstreamSpec{
					up.Metadata.Ref(): up.GetKube(),
				})

			Expect(endpoints).To(HaveLen(1), "expected to have 1 endpoint")
			Expect(warnsToLog).To(BeEmpty(), "expected no warnings")
			Expect(errorsToLog).To(BeEmpty(), "expected no errors")

			// Check endpoint has automtls metadata
			Expect(endpoints[0].Metadata.Labels).To(HaveKeyWithValue(constants.IstioTlsModeLabel, "istio"))
		})
	})
})
