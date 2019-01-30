package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	ratelimit2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"

	"github.com/gogo/protobuf/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/services"
	ratelimitservice "github.com/solo-io/solo-projects/test/services/ratelimit"

	rlservice "github.com/solo-io/rate-limiter/pkg/service"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
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
		return "redis-server"
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
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		// add the rl service as a static upstream
		rlserver := &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "rl-server",
				Namespace: "default",
			},
			UpstreamSpec: &gloov1.UpstreamSpec{
				UpstreamType: &gloov1.UpstreamSpec_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: "localhost",
							Port: rlport,
						}},
						UseHttp2: true,
					},
				},
			},
		}

		testClients.UpstreamClient.Write(rlserver, clients.WriteOpts{})
		ref := rlserver.Metadata.Ref()
		rlSettings := &ratelimit.Settings{
			RatelimitServerRef: &ref,
		}
		settingsStruct, err := envoyutil.MessageToStruct(rlSettings)
		Expect(err).NotTo(HaveOccurred())

		extensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				ratelimit2.ExtensionName: settingsStruct,
			},
		}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, extensions)
		rlService = ratelimitservice.RunRatelimit(ctx, testClients.GlooPort)
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

		rateLimits := &ratelimit.IngressRateLimit{
			AnonymousLimits: &ratelimit.RateLimit{
				RequestsPerUnit: 1,
				Unit:            ratelimit.RateLimit_SECOND,
			},
		}
		rateLimitStruct, err := envoyutil.MessageToStruct(rateLimits)
		Expect(err).NotTo(HaveOccurred())
		protos := map[string]*types.Struct{
			ratelimit2.ExtensionName: rateLimitStruct,
		}

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
								Extensions: &gloov1.Extensions{
									Configs: protos,
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
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
		)

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			envoyInstance.RatelimitAddr = rladdr
			envoyInstance.RatelimitPort = rlport

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
			// drain channel as we dont care about it
			go func() {
				for range testUpstream.C {
				}
			}()
			var opts clients.WriteOpts
			up := testUpstream.Upstream
			_, err = testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		It("should should rate limit envoy", func() {

			hosts := map[string]bool{"host1": true}
			proxy := getProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts)

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			rls := rlService
			Eventually(rls.GetCurrentConfig, "5s").Should(Not(BeNil()))
			EventuallyRateLimited("host1", envoyPort)
		})

		It("should should rate limit two vhosts", func() {

			hosts := map[string]bool{"host1": true, "host2": true}
			proxy := getProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts)

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			rls := rlService
			Eventually(rls.GetCurrentConfig, "5s").Should(Not(BeNil()))

			EventuallyRateLimited("host1", envoyPort)
			EventuallyRateLimited("host2", envoyPort)
		})
		It("should should rate limit one vhosts", func() {

			hosts := map[string]bool{"host1": false, "host2": true}
			proxy := getProxy(envoyPort, testUpstream.Upstream.Metadata.Ref(), hosts)

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			rls := rlService
			Eventually(rls.GetCurrentConfig, "5s").Should(Not(BeNil()))

			// waiting for envoy to start, so that consistently works
			EventuallyOk("host1", envoyPort)

			ConsistentlyNotRateLimited("host1", envoyPort)
			EventuallyRateLimited("host2", envoyPort)
		})
	})
})

func EventuallyOk(hostname string, port uint32) {
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func ConsistentlyNotRateLimited(hostname string, port uint32) {
	ConsistentlyWithOffset(1, func() error {
		res, err := get(hostname, port)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return errors.New(fmt.Sprintf("%v is not OK", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}
func EventuallyRateLimited(hostname string, port uint32) {
	EventuallyWithOffset(1, func() error {
		res, err := get(hostname, port)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusTooManyRequests {
			return errors.New(fmt.Sprintf("%v is not TooManyRequests", res.StatusCode))
		}
		return nil
	}, "5s", ".1s").Should(BeNil())
}

func get(hostname string, port uint32) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", port), nil)
	Expect(err).NotTo(HaveOccurred())
	req.Host = hostname
	return http.DefaultClient.Do(req)
}

func getProxy(envoyPort uint32, upstream core.ResourceRef, hostsToRateLimits map[string]bool) *gloov1.Proxy {
	var extensions *gloov1.Extensions

	rateLimits := &ratelimit.IngressRateLimit{
		AnonymousLimits: &ratelimit.RateLimit{
			RequestsPerUnit: 1,
			Unit:            ratelimit.RateLimit_SECOND,
		},
	}
	rateLimitStruct, err := envoyutil.MessageToStruct(rateLimits)
	Expect(err).NotTo(HaveOccurred())
	protos := map[string]*types.Struct{
		ratelimit2.ExtensionName: rateLimitStruct,
	}

	extensions = &gloov1.Extensions{
		Configs: protos,
	}

	var vhosts []*gloov1.VirtualHost

	for hostname, enableRateLimits := range hostsToRateLimits {
		vhost := &gloov1.VirtualHost{
			Name:    "virt1",
			Domains: []string{hostname},
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
								Upstream: upstream,
							},
						},
					},
				},
			}},
		}

		if enableRateLimits {
			vhost.VirtualHostPlugins = &gloov1.VirtualHostPlugins{
				Extensions: extensions,
			}
		}
		vhosts = append(vhosts, vhost)
	}

	p := &gloov1.Proxy{
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
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}
