package gloo_test

import (
	"fmt"
	"net"
	"strings"

	kubetestclients "github.com/solo-io/gloo/test/kubernetes/testutils/clients"

	"github.com/solo-io/gloo/test/services/envoy"
	corev1 "k8s.io/api/core/v1"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	testhelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/services"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skkubeutils "github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/test/helpers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("Happy path", func() {

	var (
		testClients   services.TestClients
		envoyInstance *envoy.Instance
		tu            *v1helpers.TestUpstream
	)

	BeforeEach(func() {
		envoyInstance = envoyFactory.NewInstance()

		tu = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
	})

	AfterEach(func() {
		envoyInstance.Clean()
	})

	TestUpstreamReachable := func() {
		v1helpers.TestUpstreamReachableWithOffset(3, envoyInstance.HttpPort, tu, nil)
	}

	Describe("kubernetes happy path", func() {

		var (
			testNamespace  string
			writeNamespace string
			kubeClient     kubernetes.Interface
			svc            *corev1.Service
		)

		BeforeEach(func() {
			testNamespace = ""
			writeNamespace = ""
			svc = nil
			kubeClient = kubetestclients.MustClientset()
		})

		AfterEach(func() {
			if testNamespace != "" {
				err := kubeClient.CoreV1().Namespaces().Delete(ctx, testNamespace, metav1.DeleteOptions{})
				Expect(err).NotTo(HaveOccurred())
			}
		})

		prepNamespace := func() {
			if testNamespace == "" {
				testNamespace = "gloo-e2e-" + helpers.RandString(8)
			}

			_, err := kubeClient.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: testNamespace,
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			svc, err = kubeClient.CoreV1().Services(testNamespace).Create(ctx, &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      "headlessservice",
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						{
							Name: "foo",
							Port: int32(tu.Port),
						},
					},
				},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			_, err = kubeClient.CoreV1().Endpoints(testNamespace).Create(ctx, &corev1.Endpoints{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: testNamespace,
					Name:      svc.Name,
				},
				Subsets: []corev1.EndpointSubset{{
					Addresses: []corev1.EndpointAddress{{
						IP:       getNonSpecialIP(envoyInstance),
						Hostname: "localhost",
					}},
					Ports: []corev1.EndpointPort{{
						Port: int32(tu.Port),
					}},
				}},
			}, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())
		}

		getUpstream := func() (*gloov1.Upstream, error) {
			l, err := testClients.UpstreamClient.List(writeNamespace, clients.ListOpts{})
			if err != nil {
				return nil, err
			}
			for _, u := range l {
				if strings.Contains(u.Metadata.Name, svc.Name) && strings.Contains(u.Metadata.Name, svc.Namespace) {
					return u, nil
				}
			}
			return nil, fmt.Errorf("not found")
		}

		Context("specific namespace", func() {

			BeforeEach(func() {
				prepNamespace()
				writeNamespace = testNamespace
				ro := &services.RunOptions{
					NsToWrite: writeNamespace,
					NsToWatch: []string{testNamespace},
					WhatToRun: services.What{
						DisableGateway: true,
					},
					KubeClient: kubeClient,
					Settings: &gloov1.Settings{
						Gloo: &gloov1.GlooOptions{
							EnableRestEds: &wrappers.BoolValue{
								Value: false,
							},
						},
					},
				}

				testClients = services.RunGlooGatewayUdsFds(ctx, ro)
				role := testNamespace + "~" + gatewaydefaults.GatewayProxyName
				err := envoyInstance.RunWithRole(role, testClients.GlooPort)
				Expect(err).NotTo(HaveOccurred())

				// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
				// other syncers that have translated them, so we can only detect that the objects exist here
				testhelpers.EventuallyResourceExists(func() (resources.Resource, error) {
					return getUpstream()
				}, "20s", ".5s")
			})

			It("should discover service", func() {
				up, err := getUpstream()
				Expect(err).NotTo(HaveOccurred())

				proxycli := testClients.ProxyClient
				proxy := getTrivialProxyForUpstream(testNamespace, envoyInstance.HttpPort, up.Metadata.Ref())
				var opts clients.WriteOpts
				_, err = proxycli.Write(proxy, opts)
				Expect(err).NotTo(HaveOccurred())

				TestUpstreamReachable()
			})

			It("should create appropriate config for discovered service", func() {
				up, err := getUpstream()
				Expect(err).NotTo(HaveOccurred())

				svc.Annotations = map[string]string{
					"gloo.solo.io/upstream_config": "{\"initial_stream_window_size\": 2048}",
				}
				svc, err = kubeClient.CoreV1().Services(testNamespace).Update(ctx, svc, metav1.UpdateOptions{})
				Expect(err).NotTo(HaveOccurred())

				proxycli := testClients.ProxyClient
				proxy := getTrivialProxyForUpstream(testNamespace, envoyInstance.HttpPort, up.Metadata.Ref())
				var opts clients.WriteOpts
				_, err = proxycli.Write(proxy, opts)
				Expect(err).NotTo(HaveOccurred())

				TestUpstreamReachable()
				upstream, err := getUpstream()
				Expect(err).NotTo(HaveOccurred())
				Expect(int(upstream.GetInitialStreamWindowSize().GetValue())).To(Equal(2048))
			})

			It("correctly routes requests to a service destination", func() {
				svcRef := skkubeutils.FromKubeMeta(svc.ObjectMeta, true).Ref()
				svcPort := svc.Spec.Ports[0].Port
				proxy := getTrivialProxyForService(testNamespace, envoyInstance.HttpPort, svcRef, uint32(svcPort))

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				TestUpstreamReachable()
			})
		})

		Context("all namespaces", func() {
			BeforeEach(func() {
				testNamespace = "gloo-e2e-" + helpers.RandString(8)
				prepNamespace()

				writeNamespace = testNamespace
				ro := &services.RunOptions{
					NsToWrite: writeNamespace,
					NsToWatch: []string{},
					WhatToRun: services.What{
						DisableGateway: true,
					},
					KubeClient: kubeClient,
					Settings: &gloov1.Settings{
						Gloo: &gloov1.GlooOptions{
							EnableRestEds: &wrappers.BoolValue{
								Value: false,
							},
						},
					},
				}

				testClients = services.RunGlooGatewayUdsFds(ctx, ro)
				role := testNamespace + "~" + gatewaydefaults.GatewayProxyName
				err := envoyInstance.RunWithRole(role, testClients.GlooPort)
				Expect(err).NotTo(HaveOccurred())

			})

			It("watch all namespaces", func() {
				// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
				// other syncers that have translated them, so we can only detect that the objects exist here
				testhelpers.EventuallyResourceExists(func() (resources.Resource, error) {
					return getUpstream()
				}, "20s", ".5s")

				up, err := getUpstream()
				Expect(err).NotTo(HaveOccurred())

				proxycli := testClients.ProxyClient
				proxy := getTrivialProxyForUpstream(testNamespace, envoyInstance.HttpPort, up.Metadata.Ref())
				var opts clients.WriteOpts
				_, err = proxycli.Write(proxy, opts)
				Expect(err).NotTo(HaveOccurred())

				TestUpstreamReachable()
			})
		})
	})
})

func getTrivialProxyForUpstream(ns string, bindPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	proxy := getTrivialProxy(ns, bindPort)
	proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.
		VirtualHosts[0].Routes[0].Action.(*gloov1.Route_RouteAction).RouteAction.
		Destination.(*gloov1.RouteAction_Single).Single.DestinationType =
		&gloov1.Destination_Upstream{Upstream: upstream}
	return proxy
}

func getTrivialProxyForService(ns string, bindPort uint32, service *core.ResourceRef, svcPort uint32) *gloov1.Proxy {
	proxy := getTrivialProxy(ns, bindPort)
	proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.
		VirtualHosts[0].Routes[0].Action.(*gloov1.Route_RouteAction).RouteAction.
		Destination.(*gloov1.RouteAction_Single).Single.DestinationType =
		&gloov1.Destination_Kube{
			Kube: &gloov1.KubernetesServiceDestination{
				Ref:  service,
				Port: svcPort,
			},
		}
	return proxy
}

func getTrivialProxy(ns string, bindPort uint32) *gloov1.Proxy {
	return &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      gatewaydefaults.GatewayProxyName,
			Namespace: ns,
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "::",
			BindPort:    bindPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{{
							Action: &gloov1.Route_RouteAction{
								RouteAction: &gloov1.RouteAction{
									Destination: &gloov1.RouteAction_Single{
										Single: &gloov1.Destination{},
									},
								},
							},
						}},
					}},
				},
			},
		}},
	}
}

// getNonSpecialIP returns a non-special IP that Kubernetes will allow in an endpoint.
func getNonSpecialIP(instance *envoy.Instance) string {
	if instance.UseDocker {
		return instance.LocalAddr()
	}

	ifaces, err := net.Interfaces()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if isNonSpecialIP(ip) {
				return ip.String()
			}
		}
	}
	Fail("no ip address available", 1)
	return ""
}

// isNonSpecialIP is adapted from ValidateNonSpecialIP in k8s.io/kubernetes/pkg/apis/core/validation/validation.go
//
// Specifically disallowed are unspecified, loopback addresses, and link-local addresses
// which tend to be used for node-centric purposes (e.g. metadata service).
func isNonSpecialIP(ip net.IP) bool {
	if ip == nil {
		return false // must be a valid IP address
	}
	if ip.IsUnspecified() {
		return false // may not be unspecified
	}
	if ip.IsLoopback() {
		return false // may not be in the loopback range (127.0.0.0/8, ::1/128)
	}
	if ip.IsLinkLocalUnicast() {
		return false // may not be in the link-local range (169.254.0.0/16, fe80::/10)
	}
	if ip.IsLinkLocalMulticast() {
		return false // may not be in the link-local multicast range (224.0.0.0/24, ff02::/10)
	}
	return true
}
