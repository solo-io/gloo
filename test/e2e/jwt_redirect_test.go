package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	jwtplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Http Sanitize Headers Local E2E", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		testUpstream  *v1helpers.TestUpstream
		envoyPort     uint32
		vhosts        []*gloov1.VirtualHost

		jwksPort uint32
		// privateKey     *rsa.PrivateKey
		jwtksServerRef *core.ResourceRef
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
			DisableFds:     true,
			DisableUds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: "gloo-system", Settings: &gloov1.Settings{}})
	})

	runEnvoy := func() {
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())
		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())
	}

	setupProxy := func(headerSanitation bool) {

		envoyPort = defaults.HttpPort
		proxy := getProxyWithVhosts(envoyPort, vhosts)

		_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", envoyInstance.AdminPort), nil)
		Expect(err).NotTo(HaveOccurred())
		client := &http.Client{}
		Eventually(func() (int, error) {
			response, err := client.Do(request)
			if err != nil {
				return 0, err
			}
			defer response.Body.Close()
			_, _ = io.ReadAll(response.Body)
			return response.StatusCode, nil
		}, 20*time.Second, 1*time.Second).Should(Equal(200))

	}

	setupJwt := func() {
		// JWT authentication server (jwksServer) setup
		// jwksPort, privateKey = jwks(ctx)
		jwksPort, _ = jwks(ctx)

		Eventually(envoyInstance.GlooAddr, 20*time.Second, 1*time.Second).Should(Not(BeNil()))

		jwksServer := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "jwks-server",
				Namespace: "default",
			},
			UseHttp2: &wrappers.BoolValue{Value: true},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: envoyInstance.GlooAddr,
						Port: jwksPort,
					}},
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(jwksServer, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		jwtksServerRef = jwksServer.Metadata.Ref()
	}

	setupTestUpstream := func() {
		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())
		_, err := testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
	}

	AfterEach(func() {
		cancel()
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
	})

	Context("With vhost-level jwt config", func() {

		var (
			vhost *gloov1.VirtualHost
			route *gloov1.Route
		)

		BeforeEach(func() {
			runEnvoy()
			setupJwt()
		})

		Context("direct response route", func() {
			It("should reject unauthenticated requests", func() {
				// write a vhost with vhost level jwt authn
				vhost = getJwtVhost(jwtksServerRef)
				route = getDirectResponseRoute()
				vhost.Routes = []*gloov1.Route{route}
				vhosts = []*gloov1.VirtualHost{vhost}
				setupProxy(false)

				expectStatusCode(envoyPort, 401)
			})

			It("should not reject unauthenticated requests when jwt authn is disabled at route-level", func() {
				// write a vhost with vhost level jwt authn, but disable the jwt authn on the only route
				vhost = getJwtVhost(jwtksServerRef)
				route = disableAuthOnRoute(getDirectResponseRoute())
				vhost.Routes = []*gloov1.Route{route}
				vhosts = []*gloov1.VirtualHost{vhost}
				setupProxy(false)
				expectStatusCode(envoyPort, 200)
			})
		})

		Context("redirectAction route", func() {
			It("should reject unauthenticated requests", func() {
				// write a vhost with vhost level jwt authn
				vhost = getJwtVhost(jwtksServerRef)
				route = getRedirectActionRoute()
				vhost.Routes = []*gloov1.Route{route}
				vhosts = []*gloov1.VirtualHost{vhost}
				setupProxy(false)

				expectStatusCode(envoyPort, 401)
			})

			It("should not reject unauthenticated requests when jwt authn is disabled at route-level", func() {
				// write a vhost with vhost level jwt authn, but disable the jwt authn on the only route
				vhost = getJwtVhost(jwtksServerRef)
				route = disableAuthOnRoute(getRedirectActionRoute())
				vhost.Routes = []*gloov1.Route{route}
				vhosts = []*gloov1.VirtualHost{vhost}
				setupProxy(false)

				expectStatusCode(envoyPort, 301)
			})
		})

		Context("routeAction route", func() {
			It("should reject unauthenticated requests", func() {
				// write a vhost with vhost level jwt authn
				vhost = getJwtVhost(jwtksServerRef)
				setupTestUpstream()
				route = getRouteActionRoute(testUpstream.Upstream.Metadata.Ref())
				vhost.Routes = []*gloov1.Route{route}
				vhosts = []*gloov1.VirtualHost{vhost}
				setupProxy(false)

				expectStatusCode(envoyPort, 401)
			})

			It("should not reject unauthenticated requests when jwt authn is disabled at route-level", func() {
				// write a vhost with vhost level jwt authn, but disable the jwt authn on the only route
				vhost = getJwtVhost(jwtksServerRef)
				setupTestUpstream()
				route = disableAuthOnRoute(getRouteActionRoute(testUpstream.Upstream.Metadata.Ref()))
				vhost.Routes = []*gloov1.Route{route}
				vhosts = []*gloov1.VirtualHost{vhost}
				setupProxy(false)

				expectStatusCode(envoyPort, 200)
			})
		})
	})
})

func getJwtVhost(jwtksServerRef *core.ResourceRef) *gloov1.VirtualHost {
	vhost := &gloov1.VirtualHost{
		Name:    "virt1",
		Domains: []string{"*"},
		Routes:  []*gloov1.Route{},
		Options: &gloov1.VirtualHostOptions{
			JwtConfig: &gloov1.VirtualHostOptions_JwtStaged{
				JwtStaged: &jwtplugin.JwtStagedVhostExtension{
					BeforeExtAuth: getJwtVhostCfg(jwtksServerRef, false, true),
				},
			},
		},
	}

	return vhost
}

func getDirectResponseRoute() *gloov1.Route {
	route := gloov1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: "/",
			},
		}},
		Action: &gloov1.Route_DirectResponseAction{
			DirectResponseAction: &gloov1.DirectResponseAction{
				Status: 200,
			},
		},
	}

	return &route
}

func getRedirectActionRoute() *gloov1.Route {
	route := gloov1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: "/",
			},
		}},
		Action: &gloov1.Route_RedirectAction{
			RedirectAction: &gloov1.RedirectAction{
				PathRewriteSpecifier: &gloov1.RedirectAction_PathRedirect{
					PathRedirect: "/redirect",
				},
			},
		},
	}

	return &route
}

func getRouteActionRoute(upstream *core.ResourceRef) *gloov1.Route {
	route := gloov1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: "/",
			},
		}},
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
	}

	return &route
}

func disableAuthOnRoute(route *gloov1.Route) *gloov1.Route {
	route.Options = &gloov1.RouteOptions{
		JwtConfig: &gloov1.RouteOptions_JwtStaged{
			JwtStaged: &jwtplugin.JwtStagedRouteExtension{
				BeforeExtAuth: &jwtplugin.RouteExtension{
					Disable: true,
				},
			},
		},
	}

	return route
}

func getProxyWithVhosts(envoyPort uint32, vhosts []*gloov1.VirtualHost) *gloov1.Proxy {
	return &gloov1.Proxy{
		Metadata: &core.Metadata{
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
}

func expectStatusCode(envoyPort uint32, statusCode int) {
	request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", envoyPort), nil)
	Expect(err).NotTo(HaveOccurred())

	Eventually(func() int {
		client := &http.Client{
			// do not follow redirects
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		response, err := client.Do(request)
		if err != nil {
			return 0
		}
		defer response.Body.Close()
		_, _ = io.ReadAll(response.Body)
		return response.StatusCode
	}, 10*time.Second, 1*time.Second).Should(Equal(statusCode))
}
