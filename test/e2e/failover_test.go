package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	kubeconverters "github.com/solo-io/gloo/projects/gloo/pkg/api/converters/kube"
	corev2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("Failover", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients

		timeout = 1 * time.Second

		simpleProxy = func(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
			var vhosts []*gloov1.VirtualHost

			vhost := &gloov1.VirtualHost{
				Name:    "gloo-system.virt1",
				Domains: []string{"*"},
				Routes: []*gloov1.Route{{
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: upstream,
									},
								},
							},
						},
					},
				}},
			}

			vhosts = append(vhosts, vhost)

			p := &gloov1.Proxy{
				Metadata: &core.Metadata{
					Name:      "proxy",
					Namespace: "default",
				},
				Listeners: []*gloov1.Listener{{
					Name:        "listener",
					BindAddress: net.IPv4zero.String(),
					BindPort:    envoyPort,
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: vhosts,
						},
					},
				}},
			}

			return p
		}
	)

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{
			Ctx:           ctx,
			Cache:         cache,
			LocalGlooPort: int32(testClients.GlooPort),
			What:          what,
			Namespace:     defaults.GlooSystem,
		},
		)
	})

	AfterEach(func() {
		cancel()
	})

	Context("Local Envoy", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
		)

		var testRequest = func() string {
			var bodyStr string
			EventuallyWithOffset(3, func() (int, error) {
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
				Expect(err).NotTo(HaveOccurred())
				resp, err := client.Do(&http.Request{
					Method: http.MethodGet,
					URL:    reqUrl,
				})
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return 0, err
				}
				bodyStr = string(body)
				return resp.StatusCode, nil
			}, "20s", "1s").Should(Equal(http.StatusOK))
			return bodyStr
		}

		var testRequestReturns = func(result string) {
			bodyStr := testRequest()
			ExpectWithOffset(2, bodyStr).To(ContainSubstring(result))
		}

		Context("Failover", func() {

			var testFailover = func(address string) {
				unhealthyCtx, unhealthyCancel := context.WithCancel(context.Background())

				secret := helpers.GetKubeSecret("tls", "gloo-system")
				glooSecret, err := (&kubeconverters.TLSSecretConverter{}).FromKubeSecret(ctx, nil, secret)
				_, err = testClients.SecretClient.Write(glooSecret.(*gloov1.Secret), clients.WriteOpts{})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				testUpstream = v1helpers.NewTestHttpUpstreamWithReply(unhealthyCtx, envoyInstance.LocalAddr(), "hello")
				testUpstream2 := v1helpers.NewTestHttpsUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "world")
				testUpstream.Upstream.HealthChecks = []*corev2.HealthCheck{
					{
						HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
							HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
								Path: "/health",
							},
						},
						HealthyThreshold: &wrappers.UInt32Value{
							Value: 1,
						},
						UnhealthyThreshold: &wrappers.UInt32Value{
							Value: 1,
						},
						NoTrafficInterval: ptypes.DurationProto(time.Second / 2),
						Timeout:           ptypes.DurationProto(timeout),
						Interval:          ptypes.DurationProto(timeout),
					},
				}
				testUpstream.Upstream.Failover = &gloov1.Failover{
					PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
						{
							LocalityEndpoints: []*gloov1.LocalityLbEndpoints{
								{
									LbEndpoints: []*gloov1.LbEndpoint{
										{
											Address: address,
											Port:    testUpstream2.Port,
											UpstreamSslConfig: &gloossl.UpstreamSslConfig{
												SslSecrets: &gloossl.UpstreamSslConfig_SecretRef{
													SecretRef: &core.ResourceRef{
														Name:      "tls",
														Namespace: "gloo-system",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				proxy := simpleProxy(envoyPort, testUpstream.Upstream.Metadata.Ref())

				_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				helpers.EventuallyResourceStatusMatchesState(1, func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				}, core.Status_Accepted)

				testRequestReturns("hello")
				unhealthyCancel()
				testRequestReturns("world")
			}

			BeforeEach(func() {
				var err error
				envoyInstance, err = envoyFactory.NewEnvoyInstance()
				Expect(err).NotTo(HaveOccurred())

				err = envoyInstance.Run(testClients.GlooPort)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				if envoyInstance != nil {
					envoyInstance.Clean()
				}
			})

			It("Will use health check path specified on failover endpoint", func() {
				unhealthyCtx, unhealthyCancel := context.WithCancel(context.Background())

				secret := helpers.GetKubeSecret("tls", "gloo-system")
				glooSecret, err := (&kubeconverters.TLSSecretConverter{}).FromKubeSecret(ctx, nil, secret)
				_, err = testClients.SecretClient.Write(glooSecret.(*gloov1.Secret), clients.WriteOpts{})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				testUpstream = v1helpers.NewTestHttpUpstreamWithReply(unhealthyCtx, envoyInstance.LocalAddr(), "hello")
				testUpstream2 := v1helpers.NewTestHttpsUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "world")
				testUpstream.Upstream.HealthChecks = []*corev2.HealthCheck{
					{
						HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
							HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
								Path: "/health",
							},
						},
						HealthyThreshold: &wrappers.UInt32Value{
							Value: 1,
						},
						UnhealthyThreshold: &wrappers.UInt32Value{
							Value: 1,
						},
						NoTrafficInterval: ptypes.DurationProto(time.Second / 2),
						Timeout:           ptypes.DurationProto(timeout),
						Interval:          ptypes.DurationProto(timeout),
					},
				}
				testUpstream.Upstream.Failover = &gloov1.Failover{
					PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
						{
							LocalityEndpoints: []*gloov1.LocalityLbEndpoints{
								{
									LbEndpoints: []*gloov1.LbEndpoint{
										{
											HealthCheckConfig: &gloov1.LbEndpoint_HealthCheckConfig{
												Path:   "/lbendpointhealth",
												Method: "POST",
											},
											Address: envoyInstance.LocalAddr(),
											Port:    testUpstream2.Port,
											UpstreamSslConfig: &gloossl.UpstreamSslConfig{
												SslSecrets: &gloossl.UpstreamSslConfig_SecretRef{
													SecretRef: &core.ResourceRef{
														Name:      "tls",
														Namespace: "gloo-system",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				}
				_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				proxy := simpleProxy(envoyPort, testUpstream.Upstream.Metadata.Ref())

				_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				helpers.EventuallyResourceStatusMatchesState(1, func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				}, core.Status_Accepted)

				testRequestReturns("hello")
				unhealthyCancel()
				// Ensure that testUpstream2 recieved a health check request at the failover-config endpoint
				r := <-testUpstream2.C
				Expect(r.Method).To(Equal("POST"))
				Expect(r.URL.Path).To(Equal("/lbendpointhealth"))
				testRequestReturns("world")
			})

			It("Will failover to testUpstream2 when the first is unhealthy", func() {
				testFailover(envoyInstance.LocalAddr())
			})

			It("Will failover to testUpstream2 when the first is unhealthy with DNS resolution", func() {
				if envoyInstance.LocalAddr() == "127.0.0.1" {
					// Domain which resolves to "127.0.0.1"
					testFailover("thing.solo.io")
				} else {
					testFailover(fmt.Sprintf("%s.xip.io", envoyInstance.LocalAddr()))
				}
			})

		})

		Context("Failover Locality Load Balancing", func() {

			var (
				// represents the context that will be cancelled, triggering an unhealthy upstream
				unhealthyCtx    context.Context
				unhealthyCancel context.CancelFunc
			)

			var prepareTestUpstreamWithFailover = func() {
				secret := helpers.GetKubeSecret("tls", "gloo-system")
				glooSecret, err := (&kubeconverters.TLSSecretConverter{}).FromKubeSecret(ctx, nil, secret)
				_, err = testClients.SecretClient.Write(glooSecret.(*gloov1.Secret), clients.WriteOpts{})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				// configure upstreams
				testUpstream = v1helpers.NewTestHttpUpstreamWithReply(unhealthyCtx, envoyInstance.LocalAddr(), "hello")
				testUpstreamEast := v1helpers.NewTestHttpsUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "east")
				testUpstreamWest := v1helpers.NewTestHttpsUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "west")

				// configure failover
				testUpstream.Upstream.HealthChecks = []*corev2.HealthCheck{
					{
						HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
							HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
								Path: "/health",
							},
						},
						HealthyThreshold: &wrappers.UInt32Value{
							Value: 1,
						},
						UnhealthyThreshold: &wrappers.UInt32Value{
							Value: 1,
						},
						NoTrafficInterval: ptypes.DurationProto(time.Second / 2),
						Timeout:           ptypes.DurationProto(timeout),
						Interval:          ptypes.DurationProto(timeout),
					},
				}
				testUpstream.Upstream.Failover = &gloov1.Failover{
					PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
						{
							LocalityEndpoints: []*gloov1.LocalityLbEndpoints{
								{
									LbEndpoints: []*gloov1.LbEndpoint{
										{
											Address: envoyInstance.LocalAddr(),
											Port:    testUpstreamEast.Port,
											UpstreamSslConfig: &gloossl.UpstreamSslConfig{
												SslSecrets: &gloossl.UpstreamSslConfig_SecretRef{
													SecretRef: &core.ResourceRef{
														Name:      "tls",
														Namespace: "gloo-system",
													},
												},
											},
										},
									},
									Locality: &gloov1.Locality{
										Region:  "east_region",
										Zone:    "east_zone",
										SubZone: "east_sub_zone",
									},
									LoadBalancingWeight: &wrappers.UInt32Value{
										Value: 75,
									},
								},
								{
									LbEndpoints: []*gloov1.LbEndpoint{
										{
											Address: envoyInstance.LocalAddr(),
											Port:    testUpstreamWest.Port,
											UpstreamSslConfig: &gloossl.UpstreamSslConfig{
												SslSecrets: &gloossl.UpstreamSslConfig_SecretRef{
													SecretRef: &core.ResourceRef{
														Name:      "tls",
														Namespace: "gloo-system",
													},
												},
											},
										},
									},
									Locality: &gloov1.Locality{
										Region:  "west_region",
										Zone:    "west_zone",
										SubZone: "west_sub_zone",
									},
									LoadBalancingWeight: &wrappers.UInt32Value{
										Value: 25,
									},
								},
							},
						},
					},
				}
			}

			var persistTestUpstream = func() {
				var err error

				_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				proxy := simpleProxy(envoyPort, testUpstream.Upstream.Metadata.Ref())

				_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
				ExpectWithOffset(1, err).NotTo(HaveOccurred())

				helpers.EventuallyResourceStatusMatchesState(1, func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				}, core.Status_Accepted)
			}

			var testRequestReturnsPercentOfTime = func(expectedResponsePercentages map[string]int) {
				var requestCount = 100       // for simplicity, assume everything out of 100
				var percentMarginOfError = 5 // allow a margin of error
				var responseCounter = make(map[string]int)

				// execute requests, and track count of each response
				for i := 0; i < requestCount; i++ {
					response := testRequest()
					responseCounter[response]++
				}

				for response, actualResponsePercentage := range responseCounter {
					expectedResponsePercentage := expectedResponsePercentages[response]

					// validate that the actual percentage falls within a margin of error
					ExpectWithOffset(1, actualResponsePercentage).Should(BeNumerically(">=", expectedResponsePercentage-percentMarginOfError))
					ExpectWithOffset(1, actualResponsePercentage).Should(BeNumerically("<=", expectedResponsePercentage+percentMarginOfError))
				}
			}

			BeforeEach(func() {
				var err error
				envoyInstance, err = envoyFactory.NewEnvoyInstance()
				Expect(err).NotTo(HaveOccurred())

				err = envoyInstance.Run(testClients.GlooPort)
				Expect(err).NotTo(HaveOccurred())

				unhealthyCtx, unhealthyCancel = context.WithCancel(context.Background())
			})

			AfterEach(func() {
				if envoyInstance != nil {
					envoyInstance.Clean()
				}
			})

			It("Will distribute load equally across localities when locality_weighted_lb_config is not set", func() {
				// Prepare the upstream with failover
				prepareTestUpstreamWithFailover()

				// Persist it, ensuring we have an up to date proxy
				persistTestUpstream()

				// Routes are handled by the upstream
				testRequestReturns("hello")

				// Cause the upstream to become unhealthy
				unhealthyCancel()

				// we have not set the locality_weighted_lb_config on the testUpstream, so we expect
				// requests to be equally distributed
				testRequestReturnsPercentOfTime(map[string]int{
					"east": 50,
					"west": 50,
				})
			})

			It("Will distribute load according to load balancing weight when locality_weighted_lb_config is set", func() {
				// Prepare the upstream with failover
				prepareTestUpstreamWithFailover()

				// Set locality_weighted_lb_config
				testUpstream.Upstream.LoadBalancerConfig = &gloov1.LoadBalancerConfig{
					LocalityConfig: &gloov1.LoadBalancerConfig_LocalityWeightedLbConfig{
						LocalityWeightedLbConfig: &empty.Empty{},
					},
				}

				// Persist it, ensuring we have an up to date proxy
				persistTestUpstream()

				// Routes are handled by the upstream
				testRequestReturns("hello")

				// Cause the upstream to become unhealthy
				unhealthyCancel()

				// we have set the locality_weighted_lb_config on the testUpstream, so we expect
				// requests to be distributed per the load balancing weight
				testRequestReturnsPercentOfTime(map[string]int{
					"east": 75,
					"west": 25,
				})
			})

		})

	})

})
