package extauth_test

import (
	"context"
	"fmt"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/e2e"

	"net/http"
	"strings"

	passthrough_utils "github.com/solo-io/ext-auth-service/pkg/config/passthrough/test_utils"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/golang/protobuf/ptypes/wrappers"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
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
		testContext *e2e.TestContext
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("default auth service and 2 named auth services", func() {

		const (
			invalidToken = "invalid-token"
			defaultToken = "default-token"
			namedTokenA  = "named-A-token"
			namedTokenB  = "named-B-token"

			namedAuthServerA = "named-A"
			namedAuthServerB = "named-B"
		)

		var (
			// A running instance of an authServer
			authServerDefault         *passthrough_utils.GrpcAuthServer
			authServerDefaultPort     = 5556
			authServerDefaultUpstream *gloov1.Upstream

			// A running instance of an authServer
			authServerNamedA         *passthrough_utils.GrpcAuthServer
			authServerNamedAPort     = 5557
			authServerNamedAUpstream *gloov1.Upstream

			// A running instance of an authServer
			authServerNamedB         *passthrough_utils.GrpcAuthServer
			authServerNamedBPort     = 5558
			authServerNamedBUpstream *gloov1.Upstream
		)

		expectRequestEventuallyReturnsResponseCodeOffset := func(offset int, path, bearerToken string, responseCode int) {
			httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPath(path).WithHeader("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
			EventuallyWithOffset(offset+1, func(g Gomega) *http.Response {
				resp, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "5s", "0.5s").Should(HaveHTTPStatus(responseCode))
		}

		expectRequestEventuallyReturnsResponseCode := func(bearerToken string, responseCode int) {
			expectRequestEventuallyReturnsResponseCodeOffset(1, "1", bearerToken, responseCode)
		}

		expectRequestPathEventuallyReturnsResponseCode := func(path, bearerToken string, responseCode int) {
			expectRequestEventuallyReturnsResponseCodeOffset(1, path, bearerToken, responseCode)
		}

		BeforeEach(func() {
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
							Addr: testContext.EnvoyInstance().LocalAddr(),
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
							Addr: testContext.EnvoyInstance().LocalAddr(),
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
							Addr: testContext.EnvoyInstance().LocalAddr(),
							Port: uint32(authServerNamedBPort),
						}},
					},
				},
			}

			// authServerDefault accepts tokens with the `default-` bearer token prefix
			authServerDefault = startLocalGrpcExtAuthServer(authServerDefaultPort, "default-")

			// authServerNamedA accepts tokens with the `named-A-` bearer token prefix
			authServerNamedA = startLocalGrpcExtAuthServer(authServerNamedAPort, "named-A-")

			// authServerNamedB accepts tokens with the `named-B-` bearer token prefix
			authServerNamedB = startLocalGrpcExtAuthServer(authServerNamedBPort, "named-B-")

			// configure upstream for authServerDefault
			testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, authServerDefaultUpstream, authServerNamedAUpstream, authServerNamedBUpstream)
			testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
				settings.Extauth = &v1.Settings{
					ExtauthzServerRef: authServerDefaultUpstream.Metadata.Ref(),
				}
				settings.NamedExtauth = map[string]*v1.Settings{
					namedAuthServerA: {
						ExtauthzServerRef: authServerNamedAUpstream.Metadata.Ref(),
					},
					namedAuthServerB: {
						ExtauthzServerRef: authServerNamedBUpstream.Metadata.Ref(),
					},
				}
			})
		})

		AfterEach(func() {
			authServerDefault.Stop()
			authServerNamedA.Stop()
			authServerNamedB.Stop()
		})

		Context("auth config is set on virtual host", func() {

			When("vhost=unset", func() {
				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("vhost=disabled", func() {
				JustBeforeEach(func() {
					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_Disable{
									Disable: true,
								},
							},
						})
						return vsBuilder.Build()
					})
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("vhost=custom (default)", func() {

				JustBeforeEach(func() {
					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{},
								},
							},
						}).WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/", testContext.TestUpstream().Upstream.GetMetadata().Ref(), nil),
						)
						return vsBuilder.Build()
					})
				})

				It("token should be validated against default server", func() {
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusOK)
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusUnauthorized)
				})
			})

			When("vhost=custom (named)", func() {

				JustBeforeEach(func() {
					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{
										Name: namedAuthServerA, // Matches the key in Settings.NamedExtauth
									},
								},
							},
						}).WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/", testContext.TestUpstream().Upstream.GetMetadata().Ref(), nil),
						)
						return vsBuilder.Build()
					})
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

				JustBeforeEach(func() {
					var routeAuthConfig *v1.ExtAuthExtension // unset
					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/", testContext.TestUpstream().Upstream.GetMetadata().Ref(), routeAuthConfig),
						)
						return vsBuilder.Build()
					})
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("route=disabled", func() {

				JustBeforeEach(func() {
					routeAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_Disable{
							Disable: true,
						},
					}

					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/", testContext.TestUpstream().Upstream.GetMetadata().Ref(), routeAuthConfig),
						)
						return vsBuilder.Build()
					})
				})

				It("auth should be disabled", func() {
					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusOK)
				})
			})

			When("route=custom (default)", func() {

				JustBeforeEach(func() {
					routeAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{},
						},
					}
					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/", testContext.TestUpstream().Upstream.GetMetadata().Ref(), routeAuthConfig),
						)
						return vsBuilder.Build()
					})
				})

				It("token should be validated against default server", func() {
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusOK)

					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusUnauthorized)
				})
			})

			When("route=custom (named)", func() {

				JustBeforeEach(func() {
					routeAuthConfig := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_CustomAuth{
							CustomAuth: &v1.CustomAuth{
								Name: namedAuthServerA,
							},
						},
					}
					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/", testContext.TestUpstream().Upstream.GetMetadata().Ref(), routeAuthConfig),
						)
						return vsBuilder.Build()
					})
				})

				It("token should be validated against named server", func() {
					expectRequestEventuallyReturnsResponseCode(namedTokenA, http.StatusOK)

					expectRequestEventuallyReturnsResponseCode(invalidToken, http.StatusUnauthorized)
					expectRequestEventuallyReturnsResponseCode(defaultToken, http.StatusUnauthorized)
				})
			})

		})

		Context("auth config is set on virtual host and route", func() {
			When("vhost=CustomAuth(named), route=ConfigRef", func() {
				BeforeEach(func() {
					authConfig := &v1.AuthConfig{
						Metadata: &core.Metadata{
							Name:      GetBasicAuthExtension().GetConfigRef().Name,
							Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*v1.AuthConfig_Config{{
							AuthConfig: &v1.AuthConfig_Config_BasicAuth{
								BasicAuth: getBasicAuthConfig(),
							},
						}},
					}
					testContext.ResourcesToCreate().AuthConfigs = v1.AuthConfigList{
						authConfig,
					}
				})

				JustBeforeEach(func() {
					authConfigToExtension := &v1.ExtAuthExtension{
						Spec: &v1.ExtAuthExtension_ConfigRef{
							ConfigRef: &core.ResourceRef{
								Name:      GetBasicAuthExtension().GetConfigRef().Name,
								Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
							},
						},
					}

					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithRoute(
							e2e.DefaultRouteName,
							getRouteToUpstream("/default", testContext.TestUpstream().Upstream.GetMetadata().Ref(), authConfigToExtension),
						).WithRoute(
							"other",
							getRouteToUpstream("/other", testContext.TestUpstream().Upstream.GetMetadata().Ref(), nil),
						).WithVirtualHostOptions(&gloov1.VirtualHostOptions{
							Extauth: &v1.ExtAuthExtension{
								Spec: &v1.ExtAuthExtension_CustomAuth{
									CustomAuth: &v1.CustomAuth{
										Name: namedAuthServerB,
									},
								},
							},
						})
						return vsBuilder.Build()
					})
				})

				It("/ route should validate token against default server", func() {
					// ensure that when we route traffic to an AuthConfig-configured route, we reject VH-level tokens
					expectRequestPathEventuallyReturnsResponseCode("default", namedTokenB, http.StatusUnauthorized)

					// verify that route-level && VH-level traffic is functional
					expectRequestPathEventuallyReturnsResponseCode("default", defaultToken, http.StatusOK)
					expectRequestPathEventuallyReturnsResponseCode("other", namedTokenB, http.StatusOK)
				})
			})

			// ensure that using a default customauth server at the virtualhost level does not override extauth config at the route level
			When("vhost=custom (default), routeA=custom (default), routeB=custom (named)", func() {

				JustBeforeEach(func() {
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

					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.
							WithRoute(e2e.DefaultRouteName, getRouteToUpstream("/default", testContext.TestUpstream().Upstream.GetMetadata().Ref(), defaultRouteAuthConfig)).
							WithRoute("named", getRouteToUpstream("/named", testContext.TestUpstream().Upstream.GetMetadata().Ref(), namedRouteAuthConfig)).
							WithRoute("other", getRouteToUpstream("/other", testContext.TestUpstream().Upstream.GetMetadata().Ref(), nil)).
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								Extauth: &v1.ExtAuthExtension{
									Spec: &v1.ExtAuthExtension_CustomAuth{
										CustomAuth: &v1.CustomAuth{},
									},
								},
							})
						return vsBuilder.Build()
					})
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
				JustBeforeEach(func() {
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

					testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.
							WithRoute(e2e.DefaultRouteName, getRouteToUpstream("/default", testContext.TestUpstream().Upstream.GetMetadata().Ref(), defaultRouteAuthConfig)).
							WithRoute("named", getRouteToUpstream("/named", testContext.TestUpstream().Upstream.GetMetadata().Ref(), namedRouteAuthConfig)).
							WithRoute("other", getRouteToUpstream("/other", testContext.TestUpstream().Upstream.GetMetadata().Ref(), nil)).
							WithVirtualHostOptions(&gloov1.VirtualHostOptions{
								Extauth: &v1.ExtAuthExtension{
									Spec: &v1.ExtAuthExtension_CustomAuth{
										CustomAuth: &v1.CustomAuth{
											Name: namedAuthServerB,
										},
									},
								},
							})
						return vsBuilder.Build()
					})
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
//
//	200 Ok - if presented with a Bearer token with the proper prefix
//	401 Unauthorized - Otherwise
func startLocalGrpcExtAuthServer(port int, expectedBearerTokenPrefix string) *passthrough_utils.GrpcAuthServer {
	authServer := &passthrough_utils.GrpcAuthServer{
		AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
			authorizationHeaders, ok := req.GetAttributes().GetRequest().GetHttp().GetHeaders()["authorization"]

			if !ok {
				return passthrough_utils.DeniedResponse(), nil
			}

			extracted := strings.Fields(authorizationHeaders)
			if len(extracted) == 2 && extracted[0] == "Bearer" {
				token := extracted[1]
				if strings.HasPrefix(token, expectedBearerTokenPrefix) {
					return passthrough_utils.OkResponse(), nil
				}
			}
			return passthrough_utils.DeniedResponse(), nil
		},
	}

	err := authServer.Start(port)
	Expect(err).NotTo(HaveOccurred())
	return authServer
}

func getRouteToUpstream(prefix string, upstreamRef *core.ResourceRef, extAuthExtension *v1.ExtAuthExtension) *gatewayv1.Route {
	return &gatewayv1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: prefix,
			},
		}},
		Action: &gatewayv1.Route_RouteAction{
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
