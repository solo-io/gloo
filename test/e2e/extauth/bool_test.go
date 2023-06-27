package extauth_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/golang-jwt/jwt"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	jwtplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"gopkg.in/square/go-jose.v2"
)

const (
	issuer    = "issuer"
	audience  = "thats-us"
	JwtAuth   = "JwtAuth"
	BasicAuth = "BasicAuth"
)

var (
	baseJwksPort = uint32(28000)
)

var _ = Describe("Staged JWT + ExtAuth", func() {

	const (
		authRoutePath   = "auth"
		publicRoutePath = "public"
	)

	var (
		testContext    *e2e.TestContextWithExtensions
		jwksPort       uint32
		privateKey     *rsa.PrivateKey
		jwtksServerRef *core.ResourceRef
		extauthConfig  *extauth.AuthConfig
		virtualService *gatewayv1.VirtualService
		token          string
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		// JWT authentication server (jwksServer) setup
		jwksPort, privateKey, _, _ = runJwkServer(testContext.Ctx())
		jwksUpstream := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "jwks-server",
				Namespace: "default",
			},
			UseHttp2: &wrappers.BoolValue{Value: true},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: testContext.EnvoyInstance().GlooAddr,
						Port: jwksPort,
					}},
				},
			},
		}

		testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, jwksUpstream)
		jwtksServerRef = jwksUpstream.Metadata.Ref()
	})

	JustBeforeEach(func() {
		testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
			virtualService,
		}
		testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
			extauthConfig,
		}
		token = getJwtTokenFor("user", privateKey)

		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	// Using the parent's Gomega instance allows us to fail here without failing the entire test suite if this was called within an Eventually.
	doRequest := func(g Gomega, request *http.Request) *http.Response {
		By("Querying " + request.RequestURI)
		resp, err := testutils.DefaultHttpClient.Do(request)
		g.ExpectWithOffset(1, err).NotTo(HaveOccurred())
		return resp
	}

	// This updates the `virtualService` variable that is used to overwrite the default VirtualService.
	// It should be called after the default VirtualService is created through testContext.BeforeEach().
	// If allowMissingOrFailedJwt is true, JWT will not immediately send an unauthorized response and allow the rest of the filter chain to run.
	// If forwardTokenUpstream is true, the jwt token will be forwarded to the upstream. If false, KeepToken on the AfterExtAuth jwt will be set to false and the token will not be forwarded upstream.
	updateVirtualService := func(allowMissingOrFailedJwt, forwardTokenUpstream bool) {
		authRouteName := e2e.DefaultRouteName + "-auth"
		// Use the default virtual service and modify it with the JWT and ExtAuth values set in test's BeforeEach
		vs := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
		virtualService = vs.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
			// Include BasicAuth and Jwt Extensions
			Extauth: GetBasicAuthExtension(),
			JwtConfig: &gloov1.VirtualHostOptions_JwtStaged{
				JwtStaged: &jwtplugin.JwtStagedVhostExtension{
					BeforeExtAuth: getJwtVhostCfg(jwtksServerRef, allowMissingOrFailedJwt, true),
					AfterExtAuth:  getJwtVhostCfg(jwtksServerRef, allowMissingOrFailedJwt, forwardTokenUpstream),
				},
			},
		}).
			WithRouteOptions(e2e.DefaultRouteName, &gloov1.RouteOptions{
				//Disable RBAC and JWT for publicly accessibly route
				Rbac: &rbac.ExtensionSettings{
					Disable: true,
				},
				JwtConfig: &gloov1.RouteOptions_JwtStaged{
					JwtStaged: &jwtplugin.JwtStagedRouteExtension{
						AfterExtAuth: &jwtplugin.RouteExtension{
							Disable: true,
						},
					},
				},
			}).
			WithRoutePrefixMatcher(e2e.DefaultRouteName, "/"+publicRoutePath).
			WithRouteActionToUpstreamRef(e2e.DefaultRouteName, testContext.TestUpstream().Upstream.GetMetadata().Ref()).
			WithRouteOptions(authRouteName, &gloov1.RouteOptions{
				//Disable RBAC and not JWT, for authn only tests
				Rbac: &rbac.ExtensionSettings{
					Disable: true,
				},
			}).
			WithRoutePrefixMatcher(authRouteName, "/"+authRoutePath).
			WithRouteActionToUpstreamRef(authRouteName, testContext.TestUpstream().Upstream.GetMetadata().Ref()).Build()
	}

	Context("Staged jwt tests", func() {
		Context("basic auth AND jwt", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth + " && " + BasicAuth)

				// Wait for jwks server to start
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPort(jwksPort).WithPath(authRoutePath)
				Eventually(func() error {
					_, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
					return err
				}).ShouldNot(HaveOccurred())

				updateVirtualService(true, true)
			})

			It("Jwt AND Basic Auth", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().
					WithHostname("user:password@localhost").WithPath(authRoutePath).
					WithHeader("x-jwt", "JWT "+token)

				// Include basic auth and Jwt token in request
				Eventually(func(g Gomega) *http.Response {
					return doRequest(g, httpRequestBuilder.Build())
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

				select {
				case received := <-testContext.TestUpstream().C:
					Expect(received.Headers).To(HaveKeyWithValue("X-Jwt", []string{"JWT " + token}))
				default:
					Fail("request didnt make it upstream")
				}
			})

			It("Jwt only fails", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPath(authRoutePath).WithHeader("x-jwt", "JWT "+token)

				Eventually(func(g Gomega) *http.Response {
					return doRequest(g, httpRequestBuilder.Build())
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))

			})
		})

		Context("Basic Auth OR Jwt", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth + " || " + BasicAuth)
				updateVirtualService(true, true)
			})

			It("Jwt only", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPath(authRoutePath).WithHeader("x-jwt", "JWT "+token)
				req := httpRequestBuilder.Build()

				Eventually(func(g Gomega) *http.Response {
					By("Querying " + req.URL.String())
					resp, err := testutils.DefaultHttpClient.Do(req)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

				select {
				case received := <-testContext.TestUpstream().C:
					Expect(received.Headers).To(HaveKeyWithValue("X-Jwt", []string{"JWT " + token}))
				default:
					Fail("request didnt make it upstream")
				}

			})

			It("Basic Auth only", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithHostname("user:password@localhost").WithPath(authRoutePath)

				Eventually(func(g Gomega) *http.Response {
					return doRequest(g, httpRequestBuilder.Build())
				}, "15s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

			})
		})

		Context("don't allow missing or failed Jwt", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth + " || " + BasicAuth)
				updateVirtualService(false, true)
			})

			It("Basic Auth fails because missing JWT immediately sends unauthorized response", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithHostname("user:password@localhost").WithPath(authRoutePath)
				req := httpRequestBuilder.Build()

				// This Eventually ensures that the new config has been picked up by envoy and returns the appropriate response
				Eventually(func(g Gomega) *http.Response {
					return doRequest(g, req)
				}, "15s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))

				Consistently(func(g Gomega) *http.Response {
					return doRequest(g, req)
				}, "2s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))
			})
		})

		Context("jwt vhost stages are assigned correct config", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth)
				// By only setting KeepToken as false on the second jwt auth stage, we should see that the request is authenticated,
				// but the upstream should not see the token header
				updateVirtualService(true, false)
			})

			It("jwt only", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPath(authRoutePath).
					WithHeader("x-jwt", "JWT "+token).WithHeader("x-additional-header", "should be seen by upstream")

				Eventually(func(g Gomega) *http.Response {
					return doRequest(g, httpRequestBuilder.Build())
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

				select {
				case received := <-testContext.TestUpstream().C:
					// test that jwt token header was sanitized by second jwt filter
					// but make sure that other headers (x-additional-header) were not sanitized
					Expect(received.Headers).To(Not(HaveKeyWithValue("X-Jwt", []string{"JWT " + token})))
					Expect(received.Headers).To(HaveKeyWithValue("X-Additional-Header", []string{"should be seen by upstream"}))
				default:
					Fail("request didnt make it upstream")
				}
			})
		})

		Context("jwt stages are assigned correct route config when specified", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth)
				// Though we are setting KeepToken false on AfterExtAuth here, we are disabling the after AfterExtauth on the '/public' route
				// so the token should be kept and forwarded upstream
				updateVirtualService(true, false)
			})

			It("should disable JWT per route", func() {
				httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPath(publicRoutePath).WithHeader("x-jwt", "JWT "+token)

				Eventually(func(g Gomega) *http.Response {
					return doRequest(g, httpRequestBuilder.Build())
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

				select {
				case received := <-testContext.TestUpstream().C:
					// test that jwt token header was sanitized by second jwt filter, but not
					Expect(received.Headers).To(HaveKeyWithValue("X-Jwt", []string{"JWT " + token}))
				default:
					Fail("request didnt make it upstream")
				}
			})

		})

	})

})

// Returns config for AuthConfig resource
func getExtauthConfig(booleanExpr string) *extauth.AuthConfig {
	jwtextauth := &extauth.AuthConfig_Config{
		Name:       &wrappers.StringValue{Value: JwtAuth},
		AuthConfig: &extauth.AuthConfig_Config_Jwt{},
	}
	basicAuth := &extauth.AuthConfig_Config{
		Name: &wrappers.StringValue{Value: BasicAuth},
		AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
			BasicAuth: getBasicAuthConfig(),
		},
	}

	return &extauth.AuthConfig{
		Metadata: &core.Metadata{
			Name:      GetBasicAuthExtension().GetConfigRef().Name,
			Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
		},
		Configs:     []*extauth.AuthConfig_Config{basicAuth, jwtextauth},
		BooleanExpr: &wrappers.StringValue{Value: booleanExpr},
	}
}

// Returns Jwt Config for virtual host perfilterconfig, taken from rbac_jwt_test
func getJwtVhostCfg(jwtksServerRef *core.ResourceRef, allowMissingFailed, keepToken bool) *jwtplugin.VhostExtension {
	return &jwtplugin.VhostExtension{
		AllowMissingOrFailedJwt: allowMissingFailed,
		Providers: map[string]*jwtplugin.Provider{
			"provider1": {
				Jwks: &jwtplugin.Jwks{
					Jwks: &jwtplugin.Jwks_Remote{
						Remote: &jwtplugin.RemoteJwks{
							Url:         "http://test/keys",
							UpstreamRef: jwtksServerRef,
						},
					},
				},
				Issuer:           issuer,
				ClockSkewSeconds: &wrappers.UInt32Value{Value: 120},
				Audiences:        []string{audience},
				KeepToken:        keepToken,
				TokenSource: &jwtplugin.TokenSource{
					Headers: []*jwtplugin.TokenSource_HeaderSource{{
						Header: "x-jwt",
						Prefix: "JWT ",
					}},
					QueryParams: []string{"jwttoken"},
				},
				ClaimsToHeaders: []*jwtplugin.ClaimToHeader{{
					Claim:  "sub",
					Header: "x-sub",
				}, {
					Claim:  "sub",
					Header: "x-sub",
					Append: true,
				}},
			},
		},
	}
}

func getJwtTokenFor(sub string, privateKey *rsa.PrivateKey) string {
	claims := jwt.StandardClaims{
		Issuer:   issuer,
		Audience: audience,
		Subject:  sub,
	}
	tok := getToken(claims, privateKey, jwt.SigningMethodRS256)
	By("using token " + tok)
	return tok
}

func runJwkServer(ctx context.Context) (uint32, *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey) {
	rsaPriv, err := rsa.GenerateKey(rand.Reader, 512)
	Expect(err).NotTo(HaveOccurred())
	ecdsaPriv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	Expect(err).NotTo(HaveOccurred())
	ed25519Pub, ed25519Priv, err := ed25519.GenerateKey(rand.Reader)
	Expect(err).NotTo(HaveOccurred())
	keySet := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{{
			Key:       rsaPriv.Public(),
			Algorithm: "RS256",
			Use:       "sig",
		}, {
			Key:       ecdsaPriv.Public(),
			Algorithm: "ES256",
			Use:       "sig",
		}, {
			Key:       ed25519Pub,
			Algorithm: "EdDSA",
			Use:       "sig",
		}},
	}

	jwksBytes, err := json.Marshal(keySet)
	Expect(err).NotTo(HaveOccurred())

	jwksPort := atomic.AddUint32(&baseJwksPort, 1) + uint32(parallel.GetPortOffset())
	jwtHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "application/json")
		w.Write(jwksBytes)
	}
	h2s := &http2.Server{}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", jwksPort),
		Handler: h2c.NewHandler(http.HandlerFunc(jwtHandler), h2s),
	}

	go s.ListenAndServe()
	go func() {
		<-ctx.Done()
		s.Shutdown(context.Background())
	}()

	// serialize json and show
	return jwksPort, rsaPriv, ecdsaPriv, ed25519Priv
}

func getToken(claims jwt.Claims, key interface{}, method jwt.SigningMethod) string {
	var s string
	var err error
	switch key.(type) {
	case *rsa.PublicKey:
		s, err = jwt.NewWithClaims(method, claims).SignedString(key.(*rsa.PublicKey))
	case *ecdsa.PublicKey:
		s, err = jwt.NewWithClaims(method, claims).SignedString(key.(*ecdsa.PublicKey))
	case *ed25519.PublicKey:
		s, err = jwt.NewWithClaims(method, claims).SignedString(key.(*ed25519.PublicKey))
	default:
		err = eris.New("Unsupported token type")
	}
	s, err = jwt.NewWithClaims(method, claims).SignedString(key)
	Expect(err).NotTo(HaveOccurred())
	return s
}
