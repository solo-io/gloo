package e2e_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/test/services"

	jwtplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"

	"github.com/fgrosse/zaptest"
	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/v1helpers"
	"gopkg.in/square/go-jose.v2"

	"github.com/dgrijalva/jwt-go"
)

var (
	baseJwksPort = uint32(28000)
)

const (
	issuer   = "issuer"
	audience = "thats-us"

	admin = "admin"
	user  = "user"
)

func jwks(ctx context.Context) (uint32, *rsa.PrivateKey) {
	priv, err := rsa.GenerateKey(rand.Reader, 512)
	Expect(err).NotTo(HaveOccurred())
	key := jose.JSONWebKey{
		Key:       priv.Public(),
		Algorithm: "RS256",
		Use:       "sig",
	}

	keySet := jose.JSONWebKeySet{
		Keys: []jose.JSONWebKey{key},
	}

	jwksBytes, err := json.Marshal(keySet)
	Expect(err).NotTo(HaveOccurred())

	jwksPort := atomic.AddUint32(&baseJwksPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)
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
	return jwksPort, priv
}

func getToken(claims jwt.StandardClaims, key *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	s, err := token.SignedString(key)
	Expect(err).NotTo(HaveOccurred())
	return s
}

var _ = Describe("JWT + RBAC", func() {

	var (
		ctx            context.Context
		cancel         context.CancelFunc
		testClients    services.TestClients
		jwksPort       uint32
		privateKey     *rsa.PrivateKey
		jwtksServerRef core.ResourceRef
		envoyInstance  *services.EnvoyInstance
		testUpstream   *v1helpers.TestUpstream
		envoyPort      = uint32(8080)
	)

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		jwksPort, privateKey = jwks(ctx)

		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		jwksServer := &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "jwks-server",
				Namespace: "default",
			},
			UseHttp2: true,
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: envoyInstance.GlooAddr,
						Port: jwksPort,
					}},
				},
			},
		}

		testClients.UpstreamClient.Write(jwksServer, clients.WriteOpts{})
		jwtksServerRef = jwksServer.Metadata.Ref()
		rbacSettings := &rbac.Settings{
			RequireRbac: true,
		}

		settings := &gloov1.Settings{Rbac: rbacSettings}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, settings)

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		var opts clients.WriteOpts
		up := testUpstream.Upstream
		_, err = testClients.UpstreamClient.Write(up, opts)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		cancel()
		if envoyInstance != nil {
			envoyInstance.Clean()
		}
	})

	ExpectAccess := func(bar, fooget, foopost int, augmentRequest func(*http.Request)) {
		query := func(method, path string) (*http.Response, error) {
			url := fmt.Sprintf("http://%s:%d%s", "localhost", envoyPort, path)
			By("Querying " + url)
			req, err := http.NewRequest(method, url, nil)
			if err != nil {
				return nil, err
			}
			augmentRequest(req)
			return http.DefaultClient.Do(req)
		}

		// test public route in eventually to let the proxy time to start
		Eventually(func() (int, error) {
			resp, err := query("GET", "/public_route")
			if err != nil {
				return 0, err
			}
			return resp.StatusCode, nil
		}, "5s", "0.5s").Should(Equal(http.StatusOK))

		// No need to do eventually here as all is initialized.
		resp, err := query("GET", "/private_route")
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, resp.StatusCode).To(Equal(http.StatusForbidden))

		resp, err = query("GET", "/bar")
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, resp.StatusCode).To(Equal(bar))

		resp, err = query("GET", "/foo")
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, resp.StatusCode).To(Equal(fooget))

		resp, err = query("POST", "/foo")
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
		ExpectWithOffset(1, resp.StatusCode).To(Equal(foopost))
	}

	getTokenFor := func(sub string) string {
		claims := jwt.StandardClaims{
			Issuer:   issuer,
			Audience: audience,
			Subject:  sub,
		}
		tok := getToken(claims, privateKey)
		By("using token " + tok)
		return tok
	}

	addBearer := func(req *http.Request, token string) {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	addToken := func(req *http.Request, sub string) {
		addBearer(req, getTokenFor(sub))
	}

	Context("jwt tests", func() {
		BeforeEach(func() {
			proxy := getProxyJwt(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() (core.Status, error) {
				proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				if err != nil {
					return core.Status{}, err
				}

				return proxy.Status, nil
			}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
				"Reason": BeEmpty(),
				"State":  Equal(core.Status_Accepted),
			}))

			// wait for key service to start
			Eventually(func() error {
				_, err := http.Get(fmt.Sprintf("http://%s:%d/", "localhost", jwksPort))
				return err
			}, "5s", "0.5s").ShouldNot(HaveOccurred())
		})

		Context("forward token", func() {
			It("should forward token upstream", func() {
				token := getTokenFor("user")

				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/authnonly", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					req.Header.Add("x-jwt", "JWT "+token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))

				select {
				case received := <-testUpstream.C:
					Expect(received.Headers).To(HaveKeyWithValue("X-Jwt", []string{"JWT " + token}))
				default:
					Fail("request didnt make it upstream")
				}

			})
		})
		Context("token source", func() {
			BeforeEach(func() {
				// drain channel as we dont care about it
				go func() {
					for range testUpstream.C {
					}
				}()
			})
			It("should get token from custom header", func() {
				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/authnonly", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					token := getTokenFor("user")
					req.Header.Add("x-jwt", "JWT "+token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})
			It("should get token from custom query param", func() {
				Eventually(func() (int, error) {
					token := getTokenFor("user")

					url := fmt.Sprintf("http://%s:%d/authnonly?jwttoken="+token, "localhost", envoyPort)
					By("Querying " + url)
					resp, err := http.Get(url)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
			})
		})

		Context("claims to headers", func() {
			It("should should move the sub claim to a header", func() {
				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/authnonly", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					token := getTokenFor("user")
					req.Header.Add("x-jwt", "JWT "+token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))

				select {
				case received := <-testUpstream.C:
					Expect(received.Headers).To(HaveKeyWithValue("X-Sub", []string{"user", "user"}))
				default:
					Fail("request didnt make it upstream")
				}
			})
			It("should re-route based on the new header added", func() {
				Eventually(func() (int, error) {
					token := getTokenFor("teatime")
					url := fmt.Sprintf("http://%s:%d/authnonly?jwttoken=%s", "localhost", envoyPort, token)
					By("Querying " + url)
					resp, err := http.Get(url)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))

				select {
				case received := <-testUpstream.C:
					Expect(received.Headers).To(HaveKeyWithValue("X-New-Header", []string{"new"}))
				default:
					Fail("request didnt make it upstream")
				}
			})
		})

	})
	Context("user access tests", func() {
		BeforeEach(func() {

			// drain channel as we dont care about it
			go func() {
				for range testUpstream.C {
				}
			}()

			proxy := getProxyJwtRbac(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref())

			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() (core.Status, error) {
				proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				if err != nil {
					return core.Status{}, err
				}

				return proxy.Status, nil
			}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
				"Reason": BeEmpty(),
				"State":  Equal(core.Status_Accepted),
			}))

			// wait for key service to start
			Eventually(func() error {
				_, err := http.Get(fmt.Sprintf("http://%s:%d/", "localhost", jwksPort))
				return err
			}, "5s", "0.5s").ShouldNot(HaveOccurred())

		})

		Context("non admin user", func() {
			It("should allow non admin user access to GET foo", func() {
				ExpectAccess(http.StatusForbidden, http.StatusOK, http.StatusForbidden,
					func(req *http.Request) { addToken(req, "user") })
			})

		})

		Context("admin user", func() {
			It("should allow everything", func() {
				ExpectAccess(http.StatusOK, http.StatusOK, http.StatusOK,
					func(req *http.Request) { addToken(req, "admin") })
			})
		})

		Context("anonymous user", func() {
			It("should only allow public route", func() {
				ExpectAccess(http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized,
					func(req *http.Request) {})
			})
		})

		Context("bad token user", func() {
			It("should only allow public route", func() {
				token := getTokenFor("admin")
				// remove some stuff to make the signature invalid
				badToken := token[:len(token)-10]
				ExpectAccess(http.StatusUnauthorized, http.StatusUnauthorized, http.StatusUnauthorized,
					func(req *http.Request) { addBearer(req, badToken) })
			})
		})

	})
})

func getProxyJwtRbac(envoyPort uint32, jwtksServerRef, upstream core.ResourceRef) *gloov1.Proxy {

	jwtCfg := &jwtplugin.VhostExtension{
		Providers: map[string]*jwtplugin.Provider{
			"testprovider": {
				Jwks: &jwtplugin.Jwks{
					Jwks: &jwtplugin.Jwks_Remote{
						Remote: &jwtplugin.RemoteJwks{
							Url:         "http://test/keys",
							UpstreamRef: &jwtksServerRef,
						},
					},
				},
				Audiences: []string{audience},
				Issuer:    issuer,
			}},
	}

	rbacCfg := &rbac.ExtensionSettings{
		Policies: map[string]*rbac.Policy{
			"user": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"iss": issuer,
							"sub": user,
						},
					},
				}},
				Permissions: &rbac.Permissions{
					PathPrefix: "/foo",
					Methods:    []string{"GET"},
				},
			},
			"admin": {
				Principals: []*rbac.Principal{{
					JwtPrincipal: &rbac.JWTPrincipal{
						Claims: map[string]string{
							"iss": issuer,
							"sub": admin,
						},
					},
				}},
				Permissions: &rbac.Permissions{},
			},
		},
	}

	return getProxyJwtRbacWithExtensions(envoyPort, jwtksServerRef, upstream, jwtCfg, rbacCfg)
}

func getProxyJwt(envoyPort uint32, jwtksServerRef, upstream core.ResourceRef) *gloov1.Proxy {
	jwtCfg := &jwtplugin.VhostExtension{
		Providers: map[string]*jwtplugin.Provider{
			"provider1": {
				Jwks: &jwtplugin.Jwks{
					Jwks: &jwtplugin.Jwks_Remote{
						Remote: &jwtplugin.RemoteJwks{
							Url:         "http://test/keys",
							UpstreamRef: &jwtksServerRef,
						},
					},
				},
				Issuer:    issuer,
				Audiences: []string{audience},
				KeepToken: true,
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

	return getProxyJwtRbacWithExtensions(envoyPort, jwtksServerRef, upstream, jwtCfg, nil)
}

func getProxyJwtRbacWithExtensions(envoyPort uint32, jwtksServerRef, upstream core.ResourceRef, jwtCfg *jwtplugin.VhostExtension, rbacCfg *rbac.ExtensionSettings) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "virt1",
		Domains: []string{"*"},
		Options: &gloov1.VirtualHostOptions{
			Rbac: rbacCfg,
			Jwt:  jwtCfg,
		},
		Routes: []*gloov1.Route{
			{
				Options: &gloov1.RouteOptions{
					Jwt:  getDisabledJwt(),
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/public_route",
					},
				}},
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
			}, {
				Options: &gloov1.RouteOptions{
					// Disable JWT and not RBAC, so that no one can get here
					Jwt: getDisabledJwt(),
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/private_route",
					},
				}},
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
			}, {
				Options: &gloov1.RouteOptions{
					Transformations: &transformation.RouteTransformations{
						RequestTransformation: &transformation.Transformation{
							TransformationType: &transformation.Transformation_TransformationTemplate{
								TransformationTemplate: &transformation.TransformationTemplate{
									Headers:            map[string]*transformation.InjaTemplate{"x-new-header": {Text: "new"}},
									BodyTransformation: &transformation.TransformationTemplate_Passthrough{Passthrough: &transformation.Passthrough{}},
								},
							},
						},
					},
					// Disable RBAC and not JWT, for authn only tests
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					Headers: []*matchers.HeaderMatcher{{
						Name:  "x-sub",
						Value: "teatime",
					}},
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/authnonly",
					},
				}},
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
			}, {
				Options: &gloov1.RouteOptions{
					// Disable RBAC and not JWT, for authn only tests
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/authnonly",
					},
				}},
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
			}, {
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

func getDisabledJwt() *jwtplugin.RouteExtension {
	return &jwtplugin.RouteExtension{
		Disable: true,
	}
}

func getDisabledRbac() *rbac.ExtensionSettings {
	return &rbac.ExtensionSettings{
		Disable: true,
	}
}
