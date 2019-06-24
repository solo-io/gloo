package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/proto"

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

	"github.com/solo-io/gloo/pkg/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/v1helpers"

	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	extauthpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("Rate Limit", func() {

	var (
		ctx            context.Context
		cancel         context.CancelFunc
		testClients    services.TestClients
		redisSession   *gexec.Session
		rlService      rlservice.RateLimitServiceServer
		glooExtensions map[string]*types.Struct
		cache          memory.InMemoryResourceCache
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
		cache = memory.NewInMemoryResourceCache()

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

		glooExtensions = map[string]*types.Struct{
			ratelimit2.ExtensionName: settingsStruct,
		}

		rlService = ratelimitservice.RunRatelimit(ctx, testClients.GlooPort)
	})

	JustBeforeEach(func() {
		extensions := &gloov1.Extensions{
			Configs: glooExtensions,
		}
		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, extensions)
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

		envoyPort := uint32(8080)
		proxy := getProxy(envoyPort, up.Metadata.Ref(), map[string]bool{"host1": true})

		proxycli := testClients.ProxyClient
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

		Context("with auth", func() {

			BeforeEach(func() {
				// start the ext auth server
				extauthport := uint32(9100)

				extauthserver := &gloov1.Upstream{
					Metadata: core.Metadata{
						Name:      "extauth-server",
						Namespace: "default",
					},
					UpstreamSpec: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Static{
							Static: &gloov1static.UpstreamSpec{
								Hosts: []*gloov1static.Host{{
									Addr: "localhost",
									Port: extauthport,
								}},
								UseHttp2: true,
							},
						},
					},
				}

				testClients.UpstreamClient.Write(extauthserver, clients.WriteOpts{})
				ref := extauthserver.Metadata.Ref()
				extauthSettings := &extauthpb.Settings{
					ExtauthzServerRef: &ref,
				}
				settingsStruct, err := envoyutil.MessageToStruct(extauthSettings)
				Expect(err).NotTo(HaveOccurred())

				glooExtensions[extauth.ExtensionName] = settingsStruct

				settings := extauthrunner.Settings{
					GlooAddress:  fmt.Sprintf("localhost:%d", testClients.GlooPort),
					DebugPort:    0,
					ServerPort:   int(extauthport),
					SigningKey:   "hello",
					UserIdHeader: "X-User-Id",
				}
				go func(testctx context.Context) {
					defer GinkgoRecover()
					err := extauthrunner.RunWithSettings(testctx, settings)
					if testctx.Err() == nil {
						Expect(err).NotTo(HaveOccurred())
					}
				}(ctx)

			})

			It("should ratelimit authorized users", func() {

				ingressRateLimit := &ratelimit.IngressRateLimit{
					AuthorizedLimits: &ratelimit.RateLimit{
						RequestsPerUnit: 1,
						Unit:            ratelimit.RateLimit_SECOND,
					},
				}
				rlb := RlProxyBuilder{
					envoyPort:         envoyPort,
					upstream:          testUpstream.Upstream.Metadata.Ref(),
					hostsToRateLimits: map[string]bool{"host1": true},
					ingressRateLimit:  ingressRateLimit,
				}
				proxy := rlb.getProxy()
				vhost := proxy.Listeners[0].ListenerType.(*gloov1.Listener_HttpListener).HttpListener.VirtualHosts[0]
				vhost.VirtualHostPlugins.Extensions.Configs[extauth.ExtensionName] = toStruct(GetBasicAuthExtension())
				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
				Expect(err).NotTo(HaveOccurred())

				rls := rlService
				Eventually(rls.GetCurrentConfig, "5s").Should(Not(BeNil()))
				// do the eventually first to give envoy a chance to start
				EventuallyRateLimited("user:password@host1", envoyPort)
				ConsistentlyNotRateLimited("host1/noauth", envoyPort)
			})
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
	ConsistentlyWithOffset(2, func() error {
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
	parts := strings.SplitN(hostname, "/", 2)
	hostname = parts[0]
	path := "1"
	if len(parts) > 1 {
		path = parts[1]
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/"+path, "localhost", port), nil)
	Expect(err).NotTo(HaveOccurred())

	// remove password part if exists
	parts = strings.SplitN(hostname, "@", 2)
	if len(parts) > 1 {
		hostname = parts[1]
		auth := strings.Split(parts[0], ":")
		req.SetBasicAuth(auth[0], auth[1])
	}

	req.Host = hostname
	return http.DefaultClient.Do(req)
}

func getProxy(envoyPort uint32, upstream core.ResourceRef, hostsToRateLimits map[string]bool) *gloov1.Proxy {

	ingressRateLimit := &ratelimit.IngressRateLimit{
		AnonymousLimits: &ratelimit.RateLimit{
			RequestsPerUnit: 1,
			Unit:            ratelimit.RateLimit_SECOND,
		},
	}
	rlb := RlProxyBuilder{
		envoyPort:         envoyPort,
		upstream:          upstream,
		hostsToRateLimits: hostsToRateLimits,
		ingressRateLimit:  ingressRateLimit,
	}
	return rlb.getProxy()
}

type RlProxyBuilder struct {
	ingressRateLimit  *ratelimit.IngressRateLimit
	upstream          core.ResourceRef
	hostsToRateLimits map[string]bool
	envoyPort         uint32
}

func (b *RlProxyBuilder) getProxy() *gloov1.Proxy {
	var extensions *gloov1.Extensions

	rateLimitStruct, err := envoyutil.MessageToStruct(b.ingressRateLimit)
	Expect(err).NotTo(HaveOccurred())
	protos := map[string]*types.Struct{
		ratelimit2.ExtensionName: rateLimitStruct,
	}

	extensions = &gloov1.Extensions{
		Configs: protos,
	}

	var vhosts []*gloov1.VirtualHost

	for hostname, enableRateLimits := range b.hostsToRateLimits {
		vhost := &gloov1.VirtualHost{
			Name:    "virt" + hostname,
			Domains: []string{hostname},
			Routes: []*gloov1.Route{
				{
					Matcher: &gloov1.Matcher{
						PathSpecifier: &gloov1.Matcher_Prefix{
							Prefix: "/noauth",
						},
					},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(b.upstream),
									},
								},
							},
						},
					},
					RoutePlugins: &gloov1.RoutePlugins{
						Extensions: &gloov1.Extensions{
							Configs: map[string]*types.Struct{
								extauth.ExtensionName: toStruct(&extauthpb.RouteExtension{
									Disable: true,
								}),
							},
						},
					},
				},
				{
					Matcher: &gloov1.Matcher{
						PathSpecifier: &gloov1.Matcher_Prefix{
							Prefix: "/",
						},
					},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: utils.ResourceRefPtr(b.upstream),
									},
								},
							},
						},
					},
				},
			},
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
			BindPort:    b.envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}

func toStruct(msg proto.Message) *types.Struct {
	strct, err := envoyutil.MessageToStruct(msg)
	Expect(err).NotTo(HaveOccurred())
	return strct
}
