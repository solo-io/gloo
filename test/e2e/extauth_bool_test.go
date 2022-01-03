package e2e_test

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"

	"github.com/onsi/ginkgo/config"
	"github.com/solo-io/ext-auth-service/pkg/server"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	jwtplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"

	"github.com/fgrosse/zaptest"
	"github.com/form3tech-oss/jwt-go"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Staged JWT + extauth ", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		testUpstream  *v1helpers.TestUpstream
		envoyPort     = uint32(8080)

		jwksPort       uint32
		privateKey     *rsa.PrivateKey
		jwtksServerRef *core.ResourceRef
	)

	BeforeEach(func() {
		// Test client and logger setup
		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		// Extauth service setup

		extAuthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)
		extAuthHealthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)
		extauthAddr := "localhost"
		if runtime.GOOS == "darwin" {
			extauthAddr = "host.docker.internal"
		}

		extAuthServer := &gloov1.Upstream{
			Metadata: &core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UseHttp2: &wrappers.BoolValue{Value: true},
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: extauthAddr,
						Port: extAuthPort,
					}},
				},
			},
		}

		ref := extAuthServer.Metadata.Ref()
		glooExtauthSettings := &extauth.Settings{
			ExtauthzServerRef: ref,
		}

		glooSettings := &gloov1.Settings{Extauth: glooExtauthSettings}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		_, err = testClients.UpstreamClient.Write(extAuthServer, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		up := testUpstream.Upstream
		_, err = testClients.UpstreamClient.Write(up, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		// JWT authentication server (jwksServer) setup
		jwksPort, privateKey = jwks(ctx)

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

		_, err = testClients.UpstreamClient.Write(jwksServer, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		jwtksServerRef = jwksServer.Metadata.Ref()

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, Settings: glooSettings})

		extauthRunnerSettings := extauthrunner.Settings{
			GlooAddress: fmt.Sprintf("localhost:%d", testClients.GlooPort),
			ExtAuthSettings: server.Settings{
				DebugPort:              0,
				ServerPort:             int(extAuthPort),
				SigningKey:             "hello",
				UserIdHeader:           "X-User-Id",
				HealthCheckFailTimeout: 2, // seconds
				HealthCheckHttpPort:    int(extAuthHealthPort),
				HealthCheckHttpPath:    "/healthcheck",
				LogSettings: server.LogSettings{
					// Disable debug logs as they are noisy. If you are writing new
					// tests, uncomment this while developing to increase verbosity. I couldn't find
					// a good way to wire this to GinkgoWriter
					//DebugMode:  "1",
					LoggerName: "extauth-service-test",
				},
			},
		}
		// Run extauth server
		go func(testCtx context.Context) {
			defer GinkgoRecover()
			os.Setenv("DEBUG_MODE", "1")
			err := extauthrunner.RunWithSettings(testCtx, extauthRunnerSettings)
			if testCtx.Err() == nil {
				Expect(err).NotTo(HaveOccurred())
			}
		}(ctx)

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		// clean up envoy
		if envoyInstance != nil {
			envoyInstance.Clean()
		}
		cancel()
	})

	Context("staged jwt tests", func() {

		var (
			extauthConfig           *extauth.AuthConfig
			allowMissingOrFailedJwt = true // If this is true, JWT will not immediately send an unauthorized response and allow the rest of the filter chain to run
			forwardTokenUpstream    = true // If this is true, the jwt token will be forwarded to the upstream, if this is false, KeepToken on the AfterExtAuth jwt will be set to false and the token will not be forwarded upstream
		)

		JustBeforeEach(func() {
			// Write AuthConfig for extauth
			_, err := testClients.AuthConfigClient.Write(extauthConfig, clients.WriteOpts{Ctx: ctx})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			proxy := getJwtAndBasicAuthProxy(envoyPort, jwtksServerRef, testUpstream.Upstream.Metadata.Ref(),
				allowMissingOrFailedJwt, forwardTokenUpstream)

			_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			// Wait for proxy to be ready
			Eventually(func() error {
				_, err := http.Get(fmt.Sprintf("http://localhost:%d", envoyPort))
				return err
			}, "5s", "0.2s").ShouldNot(HaveOccurred())
		})

		Context("basic auth AND jwt", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth + " && " + BasicAuth)

				// Wait for jwks server to start
				Eventually(func() error {
					_, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", jwksPort))
					return err
				}, "5s", "0.5s").ShouldNot(HaveOccurred())
			})

			It("Jwt AND Basic Auth", func() {
				token := getJwtTokenFor("user", privateKey)

				Eventually(func() (int, error) {
					// Include basic auth and Jwt token in request
					basicAuthPrefix := "user:password"
					url := fmt.Sprintf("http://%s@%s:%d/1", basicAuthPrefix, "localhost", envoyPort)
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

			It("Jwt only fails", func() {
				token := getJwtTokenFor("user", privateKey)

				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					req.Header.Add("x-jwt", "JWT "+token)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))

			})
		})

		Context("Basic Auth OR Jwt", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth + " || " + BasicAuth)
				allowMissingOrFailedJwt = true
			})

			It("Jwt only", func() {
				token := getJwtTokenFor("user", privateKey)

				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort)
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

			It("Basic Auth only", func() {
				Eventually(func() (int, error) {
					basicAuthPrefix := "user:password"
					url := fmt.Sprintf("http://%s@%s:%d/1", basicAuthPrefix, "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "15s", "0.5s").Should(Equal(http.StatusOK))

			})
		})

		Context("don't allow missing or failed Jwt", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth + " || " + BasicAuth)
				allowMissingOrFailedJwt = false
			})

			It("Basic Auth fails because missing JWT immediately sends unauthorized response", func() {

				Consistently(func() (int, error) {
					basicAuthPrefix := "user:password"
					url := fmt.Sprintf("http://%s@%s:%d/1", basicAuthPrefix, "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "2s", "0.5s").Should(Equal(http.StatusUnauthorized))

			})
		})

		Context("jwt vhost stages are assigned correct config", func() {
			BeforeEach(func() {
				extauthConfig = getExtauthConfig(JwtAuth)
				// By only setting KeepToken as false on the second jwt auth stage, we should see that the request is authenticated,
				// but the upstream should not see the token header
				forwardTokenUpstream = false
			})

			It("jwt only", func() {
				token := getJwtTokenFor("user", privateKey)

				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort)
					By("Querying " + url)
					req, err := http.NewRequest("GET", url, nil)
					Expect(err).NotTo(HaveOccurred())
					req.Header.Add("x-jwt", "JWT "+token)
					req.Header.Add("x-additional-header", "should be seen by upstream")
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}, "5s", "0.5s").Should(Equal(http.StatusOK))
				select {
				case received := <-testUpstream.C:
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
				forwardTokenUpstream = false
				// Though we are setting KeepToken false on AfterExtAuth here, we are disabling the after AfterExtauth on the '/public' route
				// so the token should be kept and forwarded upstream
			})

			It("should disable JWT per route", func() {
				token := getJwtTokenFor("user", privateKey)

				Eventually(func() (int, error) {
					url := fmt.Sprintf("http://%s:%d/public", "localhost", envoyPort)
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
					// test that jwt token header was sanitized by second jwt filter, but not
					Expect(received.Headers).To(HaveKeyWithValue("X-Jwt", []string{"JWT " + token}))
				default:
					Fail("request didnt make it upstream")
				}
			})

		})

	})

})

const (
	JwtAuth   = "JwtAuth"
	BasicAuth = "BasicAuth"
)

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
				Issuer:    issuer,
				Audiences: []string{audience},
				KeepToken: keepToken,
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
		Audience: []string{audience},
		Subject:  sub,
	}
	tok := getToken(claims, privateKey)
	By("using token " + tok)
	return tok
}

func getJwtAndBasicAuthProxy(envoyPort uint32, jwtksServerRef, upstream *core.ResourceRef, allowFailedJwt, forwardTokenUpstream bool) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "gloo-system.virt1",
		Domains: []string{"*"},
		Options: &gloov1.VirtualHostOptions{
			// Include BasicAuth and Jwt Extensions
			Extauth: GetBasicAuthExtension(),
			JwtConfig: &gloov1.VirtualHostOptions_JwtStaged{
				JwtStaged: &jwtplugin.JwtStagedVhostExtension{
					BeforeExtAuth: getJwtVhostCfg(jwtksServerRef, allowFailedJwt, true),
					AfterExtAuth:  getJwtVhostCfg(jwtksServerRef, allowFailedJwt, forwardTokenUpstream),
				},
			},
		},
		Routes: []*gloov1.Route{
			{
				Options: &gloov1.RouteOptions{
					//Disable RBAC and JWT for publicly accessibly route
					Rbac: getDisabledRbac(),
					JwtConfig: &gloov1.RouteOptions_JwtStaged{
						JwtStaged: &jwtplugin.JwtStagedRouteExtension{
							AfterExtAuth: &jwtplugin.RouteExtension{
								Disable: true,
							},
						},
					},
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/public",
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
			},
			{
				Options: &gloov1.RouteOptions{
					//Disable RBAC and not JWT, for authn only tests
					Rbac: getDisabledRbac(),
				},
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/1",
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
