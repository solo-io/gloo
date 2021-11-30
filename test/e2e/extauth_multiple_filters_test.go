package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/solo-io/gloo/test/helpers"

	"github.com/golang/protobuf/ptypes/wrappers"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/fgrosse/zaptest"
	"github.com/solo-io/ext-auth-service/pkg/config/passthrough/test_utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/v1helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("External auth with multiple auth servers", func() {

	// The tests validate that we can configure multiple ext_authz filters
	// in the filter chain, and selectively opt in or out of a filter
	// based on the dynamic metadata of a request
	//
	// To do this, we generate multiple grpc services that validate
	// requests by comparing the bearer token prefix, with an expected
	// value. By configuring multiple services, each looking for a separate
	// prefix, we can be sure that a token that is accepted, was only
	// validated by a single service
	// (meaning that only a single ext_authz filter was enabled)

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		envoyPort     uint32

		cache    memory.InMemoryResourceCache
		settings *gloov1.Settings
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache = memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		envoyPort = defaults.HttpPort
	})

	AfterEach(func() {
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
		cancel()
	})

	JustBeforeEach(func() {
		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, settings)
	})

	Context("default auth service and 1 named auth service", func() {

		const (
			invalidToken = "invalid-token"
			defaultToken = "default-token"
			namedTokenA  = "named-A-token"
			namedTokenB  = "named-B-token"

			namedAuthServerA = "named-A"
			namedAuthServerB = "named-B"
		)

		var (
			// The upstream that handles requests
			testUpstream *v1helpers.TestUpstream
			proxy        *gloov1.Proxy

			// A running instance of an authServer
			authServerDefault         *test_utils.GrpcAuthServer
			authServerDefaultPort     = 5556
			authServerDefaultUpstream *gloov1.Upstream

			// A running instance of an authServer
			authServerNamedA         *test_utils.GrpcAuthServer
			authServerNamedAPort     = 5557
			authServerNamedAUpstream *gloov1.Upstream

			// A running instance of an authServer
			authServerNamedB         *test_utils.GrpcAuthServer
			authServerNamedBPort     = 5558
			authServerNamedBUpstream *gloov1.Upstream
		)

		expectRequestEventuallyReturnsResponseCodeOffset := func(offset int, path, bearerToken string, responseCode int) {
			EventuallyWithOffset(offset+1, func() (int, error) {
				req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/%s", "localhost", envoyPort, path), nil)
				if err != nil {
					return 0, nil
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return 0, err
				}

				return resp.StatusCode, nil
			}, "5s", "0.5s").Should(Equal(responseCode))
		}

		expectRequestEventuallyReturnsResponseCode := func(bearerToken string, responseCode int) {
			expectRequestEventuallyReturnsResponseCodeOffset(1, "1", bearerToken, responseCode)
		}

		expectRequestPathEventuallyReturnsResponseCode := func(path, bearerToken string, responseCode int) {
			expectRequestEventuallyReturnsResponseCodeOffset(1, path, bearerToken, responseCode)
		}

		BeforeEach(func() {
			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

			authServerDefaultUpstream = &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "extauth-default",
					Namespace: "default",
				},
				UseHttp2: &wrappers.BoolValue{
					Value: true,
				},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: envoyInstance.LocalAddr(),
							Port: uint32(authServerDefaultPort),
						}},
					},
				},
			}

			authServerNamedAUpstream = &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "extauth-named-a",
					Namespace: "default",
				},
				UseHttp2: &wrappers.BoolValue{
					Value: true,
				},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: envoyInstance.LocalAddr(),
							Port: uint32(authServerNamedAPort),
						}},
					},
				},
			}

			authServerNamedBUpstream = &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "extauth-named-b",
					Namespace: "default",
				},
				UseHttp2: &wrappers.BoolValue{
					Value: true,
				},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: envoyInstance.LocalAddr(),
							Port: uint32(authServerNamedBPort),
						}},
					},
				},
			}

			settings = &gloov1.Settings{
				Extauth: &v1.Settings{
					ExtauthzServerRef: authServerDefaultUpstream.Metadata.Ref(),
				},
				NamedExtauth: map[string]*v1.Settings{
					namedAuthServerA: {
						ExtauthzServerRef: authServerNamedAUpstream.Metadata.Ref(),
					},
					namedAuthServerB: {
						ExtauthzServerRef: authServerNamedBUpstream.Metadata.Ref(),
					},
				},
			}

			// authServerDefault accepts tokens with the `default-` bearer token prefix
			authServerDefault = startLocalGrpcExtAuthServer(authServerDefaultPort, "default-")

			// authServerNamedA accepts tokens with the `named-A-` bearer token prefix
			authServerNamedA = startLocalGrpcExtAuthServer(authServerNamedAPort, "named-A-")

			// authServerNamedB accepts tokens with the `named-B-` bearer token prefix
			authServerNamedB = startLocalGrpcExtAuthServer(authServerNamedBPort, "named-B-")
		})

		AfterEach(func() {
			authServerDefault.Stop()
			authServerNamedA.Stop()
			authServerNamedB.Stop()
		})

		JustBeforeEach(func() {
			// configure upstream for authServerDefault
			_, err := testClients.UpstreamClient.Write(authServerDefaultUpstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			// Configure upstream for authServerNamed
			_, err = testClients.UpstreamClient.Write(authServerNamedAUpstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			// Configure upstream for authServerFallback
			_, err = testClients.UpstreamClient.Write(authServerNamedBUpstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			// configure upstream to handle all requests
			_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			// write proxy and ensure it is accepted
			_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})
		})

		Context("auth config is set on virtual host", func() {

			When("vhost=unset", func() {

				BeforeEach(func() {
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstream("/", testUpstream.Upstream.Metadata.Ref()),
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("vhost=disabled", func() {

				BeforeEach(func() {
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstream("/", testUpstream.Upstream.Metadata.Ref()),
						},
						Options: &gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_Disable{
									Disable: true,
								},
							},
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("vhost=custom (default)", func() {

				BeforeEach(func() {
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstream("/", testUpstream.Upstream.Metadata.Ref()),
						},
						Options: &gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{},
								},
							},
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("token should be validated against default server", func() {
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusOK)

					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusUnauthorized)
				})
			})

			When("vhost=custom (named)", func() {

				BeforeEach(func() {
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstream("/", testUpstream.Upstream.Metadata.Ref()),
						},
						Options: &gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{
										Name: namedAuthServerA, // Matches the key in Settings.NamedExtauth
									},
								},
							},
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("token should be validated against named server", func() {
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusOK)

					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusUnauthorized)
				})
			})

		})

		Context("auth config is set on route", func() {

			When("route=unset", func() {

				BeforeEach(func() {
					var routeAuthConfig *v1.ExtAuthExtension // unset

					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstreamWithAuth("/", testUpstream.Upstream.Metadata.Ref(), routeAuthConfig),
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("route=disabled", func() {

				BeforeEach(func() {
					routeAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_Disable{
							Disable: true,
						},
					}
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstreamWithAuth("/", testUpstream.Upstream.Metadata.Ref(), routeAuthConfig),
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("route=custom (default)", func() {

				BeforeEach(func() {
					routeAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{},
						},
					}
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstreamWithAuth("/", testUpstream.Upstream.Metadata.Ref(), routeAuthConfig),
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("token should be validated against default server", func() {
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusOK)

					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusUnauthorized)
				})
			})

			When("route=custom (named)", func() {

				BeforeEach(func() {
					routeAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{
								Name: namedAuthServerA,
							},
						},
					}
					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstreamWithAuth("/", testUpstream.Upstream.Metadata.Ref(), routeAuthConfig),
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("token should be validated against named server", func() {
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusOK)

					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusUnauthorized)
				})
			})

		})

		Context("auth config is set on virtual host and route", func() {

			// ensure that using a default customauth server at the virtualhost level does not override extauth config at the route level
			When("vhost=custom (default), routeA=custom (default), routeB=custom (named)", func() {

				BeforeEach(func() {
					defaultRouteAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{},
						},
					}
					namedRouteAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{
								Name: namedAuthServerA,
							},
						},
					}

					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstreamWithAuth("/default", testUpstream.Upstream.Metadata.Ref(), defaultRouteAuthConfig),
							getRouteToUpstreamWithAuth("/named", testUpstream.Upstream.Metadata.Ref(), namedRouteAuthConfig),
							getRouteToUpstream("/other", testUpstream.Upstream.Metadata.Ref()),
						},
						Options: &gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{},
								},
							},
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("/default route should validate token against default server", func() {
					defaultRoute := "default"
					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, defaultToken, http.StatusOK)

					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, invalidToken, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, namedTokenA, http.StatusUnauthorized)
				})

				It("/named route should validate token against named server", func() {
					namedRoute := "named"
					expectRequestPathEventuallyReturnsResponseCode(namedRoute, namedTokenA, http.StatusOK)

					expectRequestPathEventuallyReturnsResponseCode(namedRoute, invalidToken, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(namedRoute, defaultToken, http.StatusUnauthorized)
				})

				It("/other route should validate token against default server", func() {
					otherRoute := "other"
					expectRequestPathEventuallyReturnsResponseCode(otherRoute, defaultToken, http.StatusOK)

					expectRequestPathEventuallyReturnsResponseCode(otherRoute, invalidToken, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(otherRoute, namedTokenA, http.StatusUnauthorized)
				})
			})

			// ensure that using a named customauth server at the virtualhost level does not override extauth config at the route level
			When("vhost=custom (named), routeA=custom (default), routeB=custom (named)", func() {
				BeforeEach(func() {
					defaultRouteAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{},
						},
					}
					namedRouteAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{
								Name: namedAuthServerA,
							},
						},
					}

					vhost := &gloov1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Routes: []*gloov1.Route{
							getRouteToUpstreamWithAuth("/default", testUpstream.Upstream.Metadata.Ref(), defaultRouteAuthConfig),
							getRouteToUpstreamWithAuth("/named", testUpstream.Upstream.Metadata.Ref(), namedRouteAuthConfig),
							getRouteToUpstream("/other", testUpstream.Upstream.Metadata.Ref()),
						},
						Options: &gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{
										Name: namedAuthServerB,
									},
								},
							},
						},
					}
					proxy = getProxyWithVirtualHost(envoyPort, vhost)
				})

				It("/default route should validate token against default server", func() {
					defaultRoute := "default"
					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, defaultToken, http.StatusOK)
					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, namedTokenA, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, namedTokenB, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(defaultRoute, invalidToken, http.StatusUnauthorized)
				})

				It("/named route should validate token against named server", func() {
					namedRoute := "named"
					expectRequestPathEventuallyReturnsResponseCode(namedRoute, namedTokenA, http.StatusOK)
					expectRequestPathEventuallyReturnsResponseCode(namedRoute, defaultToken, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(namedRoute, namedTokenB, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(namedRoute, invalidToken, http.StatusUnauthorized)
				})

				It("/other route should validate token against fallback server", func() {
					otherRoute := "other"
					expectRequestPathEventuallyReturnsResponseCode(otherRoute, namedTokenB, http.StatusOK)
					expectRequestPathEventuallyReturnsResponseCode(otherRoute, defaultToken, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(otherRoute, namedTokenA, http.StatusUnauthorized)
					expectRequestPathEventuallyReturnsResponseCode(otherRoute, invalidToken, http.StatusUnauthorized)
				})
			})

		})

	})

})

// Represents an external auth service that returns:
// 	200 Ok - if presented with a Bearer token with the proper prefix
// 	401 Unauthorized - Otherwise
func startLocalGrpcExtAuthServer(port int, expectedBearerTokenPrefix string) *test_utils.GrpcAuthServer {
	authServer := &test_utils.GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			authorizationHeaders, ok := req.GetAttributes().GetRequest().GetHttp().GetHeaders()["authorization"]

			if !ok {
				return test_utils.DeniedResponse(), nil
			}

			extracted := strings.Fields(authorizationHeaders)
			if len(extracted) == 2 && extracted[0] == "Bearer" {
				token := extracted[1]
				if strings.HasPrefix(token, expectedBearerTokenPrefix) {
					return test_utils.OkResponse(), nil
				}
			}
			return test_utils.DeniedResponse(), nil
		},
	}

	err := authServer.Start(port)
	Expect(err).NotTo(HaveOccurred())
	return authServer
}

func getProxyWithVirtualHost(envoyPort uint32, vhost *gloov1.VirtualHost) *gloov1.Proxy {
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
					VirtualHosts: []*gloov1.VirtualHost{vhost},
				},
			},
		}},
	}
}

func getRouteToUpstream(prefix string, upstreamRef *core.ResourceRef) *gloov1.Route {
	return getRouteToUpstreamWithAuth(prefix, upstreamRef, nil)
}

func getRouteToUpstreamWithAuth(prefix string, upstreamRef *core.ResourceRef, extAuthExtension *v1.ExtAuthExtension) *gloov1.Route {
	return &gloov1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: prefix,
			},
		}},
		Action: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: upstreamRef,
						},
					},
				},
			},
		},
		Options: &gloov1.RouteOptions{
			Extauth: extAuthExtension,
		},
	}
}
