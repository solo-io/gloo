package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/fgrosse/zaptest"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/pkg/utils"
	kubeconverters "github.com/solo-io/gloo/projects/gloo/pkg/api/converters/kube"
	corev2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
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

		simpleProxy = func(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
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
										Upstream: utils.ResourceRefPtr(upstream),
									},
								},
							},
						},
					},
				}},
			}

			vhosts = append(vhosts, vhost)

			p := &gloov1.Proxy{
				Metadata: core.Metadata{
					Name:      "proxy",
					Namespace: "default",
				},
				Listeners: []*gloov1.Listener{{
					Name:        "listener",
					BindAddress: "0.0.0.0",
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

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(
			ctx,
			cache,
			int32(testClients.GlooPort),
			what,
			defaults.GlooSystem,
			nil,
			nil,
			nil,
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

		var testRequest = func(result string) {
			var resp *http.Response
			Eventually(func() (int, error) {
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/hello/1", "localhost", envoyPort))
				Expect(err).NotTo(HaveOccurred())
				resp, err = client.Do(&http.Request{
					Method: http.MethodGet,
					URL:    reqUrl,
				})
				if resp == nil {
					return 0, nil
				}
				return resp.StatusCode, nil
			}, "10s", "1s").Should(Equal(http.StatusOK))
			bodyStr, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(bodyStr).To(ContainSubstring(result))
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

		It("Will failover to testUpstream2 when the first is unhealthy", func() {

			unhealthyCtx, unhealthyCancel := context.WithCancel(context.Background())

			secret := helpers.GetKubeSecret("tls", "gloo-system")
			glooSecret, err := (&kubeconverters.TLSSecretConverter{}).FromKubeSecret(ctx, nil, secret)
			_, err = testClients.SecretClient.Write(glooSecret.(*gloov1.Secret), clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstreamWithReply(unhealthyCtx, envoyInstance.LocalAddr(), "hello")
			testUpstream2 := v1helpers.NewTestHttpsUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "world")
			testUpstream.Upstream.HealthChecks = []*corev2.HealthCheck{
				{
					HealthChecker: &corev2.HealthCheck_HttpHealthCheck_{
						HttpHealthCheck: &corev2.HealthCheck_HttpHealthCheck{
							Path: "/health",
						},
					},
					HealthyThreshold: &types.UInt32Value{
						Value: 1,
					},
					UnhealthyThreshold: &types.UInt32Value{
						Value: 1,
					},
					NoTrafficInterval: types.DurationProto(time.Second / 2),
					Timeout:           &timeout,
					Interval:          &timeout,
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
										Port:    testUpstream2.Port,
										UpstreamSslConfig: &gloov1.UpstreamSslConfig{
											SslSecrets: &gloov1.UpstreamSslConfig_SecretRef{
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
			Expect(err).NotTo(HaveOccurred())

			proxy := simpleProxy(envoyPort, testUpstream.Upstream.Metadata.Ref())

			_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() (core.Status, error) {
				proxy, err = testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				if err != nil {
					return core.Status{}, err
				}

				return proxy.Status, nil
			}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
				"Reason": BeEmpty(),
				"State":  Equal(core.Status_Accepted),
			}))

			testRequest("hello")
			unhealthyCancel()
			testRequest("world")
		})
	})

})
