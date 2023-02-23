//go:build linux

package e2e_test

import (
	"context"
	"net"
	"sort"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/testutils/benchmarking"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("dlp", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
	)

	var getProxyDlp = func(envoyPort uint32, upstream *core.ResourceRef, dlpListenerSettings *dlp.FilterConfig,
		dlpVhostSettings *dlp.Config, dlpRouteSettings *dlp.Config) *gloov1.Proxy {

		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			Options: &gloov1.VirtualHostOptions{
				Dlp: dlpVhostSettings,
			},
			Routes: []*gloov1.Route{
				{
					// create a direct response action for benchmarking tests.
					// the upstream used in routeaction tests were resource
					// constrained causing benchmarks to crash
					Options: &gloov1.RouteOptions{
						Dlp: dlpRouteSettings,
					},
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/visa",
						},
					}},
					Action: &gloov1.Route_DirectResponseAction{
						DirectResponseAction: &gloov1.DirectResponseAction{
							Status: 200,
							Body:   "4397-9453-4034-4828",
						},
					},
				},
			},
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
						Options: &gloov1.HttpListenerOptions{
							Dlp: dlpListenerSettings,
						},
					},
				},
			}},
		}

		return p
	}

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

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem})
	})

	AfterEach(func() {
		cancel()
	})
	Context("With envoy", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)

			proxy *gloov1.Proxy
		)

		var testRequestVisa = func(result string) {
			testRequestPath("/visa/1", result, envoyPort)
		}

		var configureProxy = func() {
			Expect(proxy).NotTo(BeNil())
			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})
		}

		BeforeEach(func() {
			proxy = nil
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "hello")
			_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("listener rules", func() {

			var configureListenerProxy = func(actions []*dlp.Action, matcher *matchers.Matcher) {
				if matcher == nil {
					matcher = &matchers.Matcher{
						PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
					}
				}
				dlpCfg := &dlp.FilterConfig{
					DlpRules: []*dlp.DlpRule{
						{
							Matcher: matcher,
							Actions: actions,
						},
					},
				}
				proxy = getProxyDlp(envoyPort, testUpstream.Upstream.Metadata.Ref(), dlpCfg, nil, nil)
				configureProxy()
			}

			Context("Benchmarking", func() {
				JustBeforeEach(func() {
					testUpstream = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "4397-9453-4034-4828")
					_, err := testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())
				})

				It("ALL_CREDIT_CARDS action benchmarking", func() {
					configureListenerProxy([]*dlp.Action{{
						ActionType: dlp.Action_ALL_CREDIT_CARDS,
					}}, nil)

					const times = 750
					f := func() {
						for i := 0; i < times; i++ {
							testRequestVisa("XXXX-XXXX-XXXX-4828")
						}
					}

					samples := 10
					results := make([]float64, samples)
					sum := float64(0)
					for i := 0; i < samples; i++ {
						results[i] = benchmarking.TimeForFuncToComplete(f)
						sum += results[i]
					}
					sort.Float64s(results)
					Expect(sum / float64(samples)).To(BeNumerically("<=", 0.25))
					Expect(results[samples-2]).To(BeNumerically("<=", 0.375))
					// reference values for each benchmark sample:
					// utime: 0.040102
					// utime: 0.032279
					// utime: 0.038947
					// utime: 0.021455
					// utime: 0.012275
					// utime: 0.017617
					// utime: 0.034505
					// utime: 0.035638
					// utime: 0.020242
					// utime: 0.020336
					// utime: 0.038511
					// utime: 0.034466
					// utime: 0.037583
					// utime: 0.038202
					// utime: 0.030788
					// utime: 0.038159
					// utime: 0.037271
					// utime: 0.029028
					// utime: 0.035937
					// utime: 0.025540
					// utime: 0.037520
					// utime: 0.047230
					// utime: 0.028943
					// utime: 0.039649
					// utime: 0.039377
				})
			})

		})
	})
})
