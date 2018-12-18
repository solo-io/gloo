package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/services"
	ratelimitservice "github.com/solo-io/solo-projects/test/services/ratelimit"

	rlservice "github.com/solo-io/rate-limiter/pkg/service"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gloov1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Rate Limit", func() {

	var (
		ctx          context.Context
		cancel       context.CancelFunc
		testClients  services.TestClients
		redisSession *gexec.Session
		rlService    rlservice.RateLimitServiceServer
	)
	const (
		redisaddr = "127.0.0.1"
		redisport = uint32(6379)
		rladdr    = "127.0.0.1"
		rlport    = uint32(18081)
	)

	getRedisPath := func() string {
		binaryPath := os.Getenv("REDIS_BINARY")
		if binaryPath != "" {
			return binaryPath
		}
		return "redis_server"
	}

	BeforeEach(func() {
		var err error
		os.Setenv("REDIS_URL", fmt.Sprintf("%s:%d", redisaddr, redisport))
		os.Setenv("REDIS_SOCKET_TYPE", "tcp")

		command := exec.Command(getRedisPath(), "--port", "6379")
		redisSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		// give redis a chance to start
		Eventually(redisSession.Out, "5s").Should(gbytes.Say("Ready to accept connections"))

		ctx, cancel = context.WithCancel(context.Background())
		t := services.RunGateway(ctx, true)
		testClients = t

		rlService = ratelimitservice.RunRatelimit(ctx, t.GlooPort)
	})

	AfterEach(func() {
		cancel()
		redisSession.Kill()
	})

	It("rlserver receives config", func() {

		tu := v1helpers.NewTestHttpUpstream(ctx, "fake-addr")
		var opts clients.WriteOpts
		up := tu.Upstream
		_, err := testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())

		proxycli := testClients.ProxyClient
		envoyPort := uint32(8080)
		proxy := &gloov1.Proxy{
			Metadata: core.Metadata{
				Name:      "proxy",
				Namespace: "default",
			},
			Listeners: []*gloov1.Listener{{
				Name:        "listener",
				BindAddress: "127.0.0.1",
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{{
							Name:    "virt1",
							Domains: []string{"*"},
							Routes: []*gloov1.Route{{
								Matcher: &gloov1.Matcher{
									PathSpecifier: &gloov1.Matcher_Prefix{
										Prefix: "/",
									},
								},
								Action: &gloov1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Single{
											Single: &gloov1.Destination{
												Upstream: up.Metadata.Ref(),
											},
										},
									},
								},
							}},
							VirtualHostPlugins: &gloov1.VirtualHostPlugins{
								RateLimits: &ratelimit.IngressRateLimit{
									AnonymousLimits: &ratelimit.RateLimit{
										RequestsPerUnit: 1,
										Unit:            ratelimit.RateLimit_SECOND,
									},
								},
							},
						}},
					},
				},
			}},
		}

		_, err = proxycli.Write(proxy, opts)
		Expect(err).NotTo(HaveOccurred())

		Eventually(rlService.GetCurrentConfig, "5s").Should(Not(BeNil()))

	})

	Context("With envoy", func() {

		var (
			envoyInstance *services.EnvoyInstance
		)

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			envoyInstance.RatelimitAddr = rladdr
			envoyInstance.RatelimitPort = rlport

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		It("should should rate limit envoy", func() {
			tu := v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
			// drain channel as we dont care about it
			go func() {
				for range tu.C {
				}
			}()
			var opts clients.WriteOpts
			up := tu.Upstream
			_, err := testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())

			proxycli := testClients.ProxyClient
			envoyPort := uint32(8080)
			proxy := &gloov1.Proxy{
				Metadata: core.Metadata{
					Name:      "proxy",
					Namespace: "default",
				},
				Listeners: []*gloov1.Listener{{
					Name:        "listener",
					BindAddress: "127.0.0.1",
					BindPort:    envoyPort,
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: []*gloov1.VirtualHost{{
								Name:    "virt1",
								Domains: []string{"*"},
								Routes: []*gloov1.Route{{
									Matcher: &gloov1.Matcher{
										PathSpecifier: &gloov1.Matcher_Prefix{
											Prefix: "/",
										},
									},
									Action: &gloov1.Route_RouteAction{
										RouteAction: &gloov1.RouteAction{
											Destination: &gloov1.RouteAction_Single{
												Single: &gloov1.Destination{
													Upstream: up.Metadata.Ref(),
												},
											},
										},
									},
								}},
								VirtualHostPlugins: &gloov1.VirtualHostPlugins{
									RateLimits: &ratelimit.IngressRateLimit{
										AnonymousLimits: &ratelimit.RateLimit{
											RequestsPerUnit: 1,
											Unit:            ratelimit.RateLimit_SECOND,
										},
									},
								},
							}},
						},
					},
				}},
			}

			_, err = proxycli.Write(proxy, opts)
			Expect(err).NotTo(HaveOccurred())

			rls := rlService
			Eventually(rls.GetCurrentConfig, "5s").Should(Not(BeNil()))

			Eventually(func() error {
				res, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
				if err != nil {
					return err
				}
				if res.StatusCode != http.StatusTooManyRequests {
					return errors.New(fmt.Sprintf("%v is not TooManyRequests", res.StatusCode))
				}
				return nil
			}, "5s", ".1s").Should(BeNil())
		})
	})
})
