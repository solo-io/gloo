package e2e_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/fgrosse/zaptest"
	"github.com/form3tech-oss/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	. "github.com/onsi/gomega/gstruct"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/test_utils"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/user_info"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	grpcPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/grpc"
	passthrough_test_utils "github.com/solo-io/ext-auth-service/pkg/config/passthrough/test_utils"
	"github.com/solo-io/ext-auth-service/pkg/controller/translation"
	"github.com/solo-io/ext-auth-service/pkg/server"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/trace/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tracing"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/matchers"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	baseExtauthPort = uint32(27000)
)

var _ = Describe("External auth", func() {

	var (
		ctx          context.Context
		cancel       context.CancelFunc
		testClients  services.TestClients
		settings     extauthrunner.Settings
		glooSettings *gloov1.Settings
		cache        memory.InMemoryResourceCache
	)

	BeforeEach(func() {
		extAuthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)
		extAuthHealthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache = memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

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

		_, err := testClients.UpstreamClient.Write(extAuthServer, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		ref := extAuthServer.Metadata.Ref()
		extauthSettings := &extauth.Settings{
			ExtauthzServerRef: ref,
			RequestBody: &extauth.BufferSettings{
				MaxRequestBytes:     0,
				AllowPartialMessage: false,
				PackAsBytes:         false,
			},
		}

		settings = extauthrunner.Settings{
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
					// Note(yuval-k): Disable debug logs as they are noisy. If you are writing new
					// tests, uncomment this while developing to increase verbosity. I couldn't find
					// a good way to wire this to GinkgoWriter
					//DebugMode:  "1",
					LoggerName: "extauth-service-test",
				},
			},
		}
		glooSettings = &gloov1.Settings{Extauth: extauthSettings}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}
		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, glooSettings)
		go func(testCtx context.Context) {
			defer GinkgoRecover()
			os.Setenv("DEBUG_MODE", "1")
			err := extauthrunner.RunWithSettings(testCtx, settings)
			if testCtx.Err() == nil {
				Expect(err).NotTo(HaveOccurred())
			}
		}(ctx)
	})

	AfterEach(func() {
		cancel()
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

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

			var opts clients.WriteOpts
			up := testUpstream.Upstream
			_, err = testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())

		})

		AfterEach(func() {
			if envoyInstance != nil {
				_ = envoyInstance.Clean()
			}
		})

		var basicConfigSetup = func() {
			_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetBasicAuthExtension().GetConfigRef().Name,
					Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
						BasicAuth: getBasicAuthConfig(),
					},
				}},
			}, clients.WriteOpts{Ctx: ctx})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			proxy := getProxyExtAuthBasicAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

			_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})
		}

		Context("using new config format", func() {

			Context("basic auth sanity tests", func() {

				BeforeEach(func() {
					basicConfigSetup()
				})

				It("should deny ext auth envoy", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})

				It("should allow ext auth envoy", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://user:password@%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})

				It("should deny ext auth with wrong password", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://user:password2@%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})
			})

			Context("oidc sanity", func() {

				var (
					authConfig      *extauth.AuthConfig
					oauth2          *extauth.OAuth2_OidcAuthorizationCode
					privateKey      *rsa.PrivateKey
					discoveryServer fakeDiscoveryServer
					handlerStats    map[string]int
					secret          *gloov1.Secret
					proxy           *gloov1.Proxy
					token           string
					cookies         []*http.Cookie
				)

				BeforeEach(func() {
					handlerStats = make(map[string]int)
					discoveryServer = fakeDiscoveryServer{
						handlerStats: handlerStats,
					}
					privateKey = discoveryServer.Start()

					clientSecret := &extauth.OauthSecret{
						ClientSecret: "test",
					}

					secret = &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      "secret",
							Namespace: "default",
						},
						Kind: &gloov1.Secret_Oauth{
							Oauth: clientSecret,
						},
					}
					_, err := testClients.SecretClient.Write(secret, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					oauth2 = getOidcAuthCodeConfig(envoyPort, secret.Metadata.Ref())
					authConfig = &extauth.AuthConfig{
						Metadata: &core.Metadata{
							Name:      getOidcExtAuthExtension().GetConfigRef().Name,
							Namespace: getOidcExtAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_Oauth2{
								Oauth2: &extauth.OAuth2{
									OauthType: oauth2,
								},
							},
						}},
					}

					proxy = getProxyExtAuthOIDC(envoyPort, testUpstream.Upstream.Metadata.Ref())
					//get id token
					token = discoveryServer.token
				})

				JustBeforeEach(func() {
					_, err := testClients.AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
						return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
					})
				})

				AfterEach(func() {
					discoveryServer.Stop()
				})

				makeSingleRequest := func(client *http.Client) (http.Response, error) {
					req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/success?foo=bar", "localhost", envoyPort), nil)
					Expect(err).NotTo(HaveOccurred())
					r, err := client.Do(req)
					if err != nil {
						return http.Response{}, err
					}

					return *r, err
				}

				ExpectHappyPathToWork := func(makeSingleRequest func(client *http.Client) (http.Response, error), loginSuccessExpectation func()) {
					// do auth flow and make sure we have a cookie named cookie:
					appPage, err := url.Parse(fmt.Sprintf("http://%s:%d/", "localhost", envoyPort))
					Expect(err).NotTo(HaveOccurred())

					var finalurl *url.URL
					jar, err := cookiejar.New(nil)
					Expect(err).NotTo(HaveOccurred())
					client := &http.Client{
						Jar: &unsecureCookieJar{CookieJar: jar},
						CheckRedirect: func(req *http.Request, via []*http.Request) error {
							finalurl = req.URL
							if len(via) > 10 {
								return errors.New("stopped after 10 redirects")
							}
							return nil
						},
					}

					Eventually(func() (http.Response, error) {
						return makeSingleRequest(client)
					}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
						"StatusCode": Equal(http.StatusOK),
					}))

					Expect(finalurl).NotTo(BeNil())
					Expect(finalurl.Path).To(Equal("/success"))
					// make sure query is passed through as well
					Expect(finalurl.RawQuery).To(Equal("foo=bar"))

					// check the cookie jar
					cookies = jar.Cookies(appPage)
					Expect(cookies).NotTo(BeEmpty())

					// make sure login is successful
					loginSuccessExpectation()

					// try to logout:

					logout := fmt.Sprintf("http://%s:%d/logout", "localhost", envoyPort)
					req, err := http.NewRequest("GET", logout, nil)
					Expect(err).NotTo(HaveOccurred())
					resp, err := client.Do(req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					// Verify that the logout resulted in a redirect to the defaul url
					Expect(finalurl).NotTo(BeNil())
					Expect(finalurl.Path).To(Equal("/"))
				}

				Context("redis for session store", func() {

					const (
						redisaddr  = "127.0.0.1"
						redisport  = uint32(6379)
						cookieName = "cookie"
					)
					var (
						redisSession *gexec.Session
					)
					BeforeEach(func() {
						// update the config to use redis
						oauth2.OidcAuthorizationCode.Session = &extauth.UserSession{
							FailOnFetchFailure: true,
							Session: &extauth.UserSession_Redis{
								Redis: &extauth.UserSession_RedisSession{
									Options: &extauth.RedisOptions{
										Host: fmt.Sprintf("%s:%d", redisaddr, redisport),
									},
									KeyPrefix:       "key",
									CookieName:      cookieName,
									AllowRefreshing: &wrappers.BoolValue{Value: true},
								},
							},
						}

						command := exec.Command(getRedisPath(), "--port", fmt.Sprintf("%d", redisport))
						var err error
						redisSession, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
						Expect(err).NotTo(HaveOccurred())
						// give redis a chance to start
						Eventually(redisSession.Out, "5s").Should(gbytes.Say("Ready to accept connections"))
					})

					AfterEach(func() {
						redisSession.Kill()
					})

					It("should work", func() {
						ExpectHappyPathToWork(makeSingleRequest, func() {
							Expect(cookies[0].Name).To(Equal(cookieName))
						})
					})

					It("should refresh token", func() {
						discoveryServer.createExpiredToken = true
						discoveryServer.updateToken("")
						ExpectHappyPathToWork(makeSingleRequest, func() {
							Expect(cookies[0].Name).To(Equal(cookieName))
						})
						Expect(discoveryServer.lastGrant).To(Equal("refresh_token"))
					})

					It("should auth successfully after refreshing token", func() {
						forceTokenRefresh := func(client *http.Client) (http.Response, error) {
							// Create token that will expire in 1 second
							discoveryServer.createNearlyExpiredToken = true
							discoveryServer.updateToken("")
							discoveryServer.createNearlyExpiredToken = false
							Expect(handlerStats["/token"]).To(BeEquivalentTo(0))

							// execute first request.
							Eventually(func() (http.Response, error) {
								return makeSingleRequest(client)
							}, "10s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
								"StatusCode": Equal(http.StatusOK),
							}))

							// sleep for 1 second, so the token expires
							time.Sleep(time.Second)

							// execute second request.
							r, err := makeSingleRequest(client)
							Expect(err).NotTo(HaveOccurred())

							// execute third request. We should not hit the /token handler, because the refreshed token should be in the store.
							baseRefreshes := handlerStats["/token"]
							r, err = makeSingleRequest(client)
							Expect(err).NotTo(HaveOccurred())
							Expect(handlerStats["/token"]).To(BeNumerically("==", baseRefreshes))

							return r, err
						}

						ExpectHappyPathToWork(forceTokenRefresh, func() {
							Expect(cookies[0].Name).To(Equal(cookieName))
						})
					})

					Context("no refreshing", func() {
						BeforeEach(func() {
							// update the config to use redis
							oauth2.OidcAuthorizationCode.Session.Session.(*extauth.UserSession_Redis).Redis.AllowRefreshing = &wrappers.BoolValue{Value: false}
						})

						It("should NOT refresh token", func() {
							discoveryServer.createExpiredToken = true
							discoveryServer.updateToken("")

							jar, err := cookiejar.New(nil)
							Expect(err).NotTo(HaveOccurred())
							client := &http.Client{
								Jar: &unsecureCookieJar{CookieJar: jar},
							}

							// as we will always provide an expired token, this will result in a
							// redirect loop.
							Eventually(func() error {
								req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/success?foo=bar", "localhost", envoyPort), nil)
								Expect(err).NotTo(HaveOccurred())
								_, err = client.Do(req)
								return err
							}, "5s", "0.5s").Should(MatchError(ContainSubstring("stopped after 10 redirects")))
							Expect(discoveryServer.lastGrant).To(Equal(""))
						})
					})

					// add context with refresh; get an expired token going and make sure it was refreshed.
				})

				Context("forward id token", func() {

					BeforeEach(func() {
						// update the config to use redis
						oauth2.OidcAuthorizationCode.Headers = &extauth.HeaderConfiguration{
							IdTokenHeader:     "foo",
							AccessTokenHeader: "bar",
						}
					})

					It("should work", func() {
						ExpectHappyPathToWork(makeSingleRequest, func() {})

						select {
						case r := <-testUpstream.C:
							Expect(r.Headers.Get("foo")).To(Equal(discoveryServer.token))
							Expect(r.Headers.Get("bar")).To(Equal("SlAV32hkKG"))
						case <-time.After(time.Second):
							Fail("timedout")
						}
					})
				})

				Context("discovery override", func() {

					BeforeEach(func() {
						oauth2.OidcAuthorizationCode.DiscoveryOverride = &extauth.DiscoveryOverride{
							AuthEndpoint: "http://localhost:5556/alternate-auth",
						}
					})

					It("should redirect to different auth endpoint with auth override", func() {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								// stop at the auth point
								if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
									return http.ErrUseLastResponse
								}
								return nil
							},
						}
						// Confirm that the response matches the one set by the /alternate-auth endpoint
						Eventually(func() (string, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							Expect(err).NotTo(HaveOccurred())
							resp, err := client.Do(req)
							if err != nil {
								return "", err
							}
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								return "", err
							}
							fmt.Fprintf(GinkgoWriter, "headers are %v \n", resp.Header)
							return string(body), nil
						}, "5s", "0.5s").Should(Equal("alternate-auth"))
					})
				})

				Context("jwks on demand cache refresh policy", func() {
					// The JWKS on demand cache refresh policy defines the behavior
					// when an id token is provided with a key id that is not in the local OIDC store
					//
					// The tests make an assumption:
					//	OIDC polls the discovery endpoint at an interval (default=15 minutes). That poll
					//	will update the local store with the freshest keys.
					//	These tests assume that a poll will not occur during a test. The reason
					//  we're comfortable with this assumption is that each test writes a new AuthConfig,
					//	which in turn creates a new AuthService, which restarts the polling. Therefore,
					//	a discovery poll should not occur during a test.

					// A request with valid token will return 200
					expectRequestWithTokenSucceeds := func(offset int, token string) {
						EventuallyWithOffset(offset+1, func() (int, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							req.Header.Add("Authorization", "Bearer "+token)

							resp, err := http.DefaultClient.Do(req)
							if err != nil {
								return 0, err
							}
							return resp.StatusCode, nil
						}, "5s", "0.5s").Should(Equal(http.StatusOK))
					}

					// A request with invalid token will be redirected to the /auth endoint
					expectRequestWithTokenFails := func(offset int, token string) {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								// stop at the auth point
								if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
									return http.ErrUseLastResponse
								}
								return nil
							},
						}

						EventuallyWithOffset(offset+1, func() (string, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							req.Header.Add("Authorization", "Bearer "+token)

							resp, err := client.Do(req)
							if err != nil {
								return "", err
							}
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								return "", err
							}
							return string(body), nil
						}, "5s", "0.5s").Should(Equal("auth"))
					}

					// create an id token
					createIdTokenWithKid := func(kid string) string {
						tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
							"foo": "bar",
							"aud": "test-clientid",
							"sub": "user",
							"iss": "http://localhost:5556",
						})
						tokenToSign.Header["kid"] = kid
						token, err := tokenToSign.SignedString(privateKey)
						Expect(err).NotTo(HaveOccurred())

						return token
					}

					JustBeforeEach(func() {
						// Ensure that keys have been loaded properly
						validToken := createIdTokenWithKid(discoveryServer.getValidKeyId())
						expectRequestWithTokenSucceeds(0, validToken)
					})

					When("policy is nil or NEVER", func() {

						BeforeEach(func() {
							oauth2.OidcAuthorizationCode.JwksCacheRefreshPolicy = nil
						})

						It("should accept token with valid kid", func() {
							validToken := createIdTokenWithKid(discoveryServer.getValidKeyId())
							expectRequestWithTokenSucceeds(0, validToken)
						})

						It("should deny token with new kid", func() {
							invalidToken := createIdTokenWithKid("kid-2")
							expectRequestWithTokenFails(0, invalidToken)
						})

						It("should deny token with new kid after keys rotate", func() {
							// rotate the keys
							discoveryServer.updateKeyIds([]string{"kid-2"})

							// execute a request with the valid token
							// it should be denied because the local cache is never updated
							newToken := createIdTokenWithKid("kid-2")
							expectRequestWithTokenFails(0, newToken)
						})

					})

					When("policy is ALWAYS", func() {

						BeforeEach(func() {
							oauth2.OidcAuthorizationCode.JwksCacheRefreshPolicy = &extauth.JwksOnDemandCacheRefreshPolicy{
								Policy: &extauth.JwksOnDemandCacheRefreshPolicy_Always{},
							}
						})

						It("should accept token with valid kid", func() {
							validToken := createIdTokenWithKid(discoveryServer.getValidKeyId())
							expectRequestWithTokenSucceeds(0, validToken)
						})

						It("should deny token with new kid", func() {
							invalidToken := createIdTokenWithKid("kid-2")
							expectRequestWithTokenFails(0, invalidToken)
						})

						It("should accept token with new kid after keys rotate", func() {
							for i := 0; i < 5; i++ {
								// rotate the keys
								newKid := fmt.Sprintf("kid-new-%d", i)
								discoveryServer.updateKeyIds([]string{newKid})

								// execute a request using the new token
								// it should be accepted because the local cache gets updated
								validToken := createIdTokenWithKid(newKid)
								expectRequestWithTokenSucceeds(0, validToken)
							}
						})
					})

					When("policy is MAX_IDP_REQUESTS_PER_POLLING_INTERVAL", func() {

						// The number of refreshes to allow before rate limiting
						// This test is subject to flakes because it relies on the behavior of OIDC
						// polling the discovery endpoint at an interval. To ensure that we complete the maxRequests
						// before the polling occurs, set the value to 1.
						const maxRequests = 1

						// The kid that will be rate limited
						const rateLimitedKid = "kid-ratelimited"

						expectRequestWithNewKidAcceptedNTimes := func() {
							// The first n times should succeed
							for i := 1; i <= maxRequests; i++ {
								// rotate the keys
								newKid := fmt.Sprintf("kid-new-%d", i)
								discoveryServer.updateKeyIds([]string{newKid})

								// execute a request using the new token
								// it should be accepted because the local cache gets updated
								validToken := createIdTokenWithKid(newKid)
								expectRequestWithTokenSucceeds(1, validToken)
							}

							// rotate the keys one more time
							discoveryServer.updateKeyIds([]string{rateLimitedKid})

							// execute a request using the new token
							// it should be rejected because the local cache no longer will be updated
							newToken := createIdTokenWithKid(rateLimitedKid)
							expectRequestWithTokenFails(1, newToken)
						}

						BeforeEach(func() {
							oauth2.OidcAuthorizationCode.JwksCacheRefreshPolicy = &extauth.JwksOnDemandCacheRefreshPolicy{
								Policy: &extauth.JwksOnDemandCacheRefreshPolicy_MaxIdpReqPerPollingInterval{
									MaxIdpReqPerPollingInterval: maxRequests,
								},
							}
						})

						It("should accept token with valid kid", func() {
							validToken := createIdTokenWithKid(discoveryServer.getValidKeyId())
							expectRequestWithTokenSucceeds(0, validToken)
						})

						It("should deny token with new kid", func() {
							invalidToken := createIdTokenWithKid("kid-2")
							expectRequestWithTokenFails(0, invalidToken)
						})

						It("should accept token with new kid after keys rotate first n times", func() {
							expectRequestWithNewKidAcceptedNTimes()
						})

						Context("after discovery poll interval", func() {
							// The rate limit is set per discovery interval, so after that interval
							// the rate limit should be reset.

							BeforeEach(func() {
								// Set the poll interval to a relatively short period
								oauth2.OidcAuthorizationCode.DiscoveryPollInterval = ptypes.DurationProto(time.Second * 4)
							})

							It("should reset rate limit", func() {
								expectRequestWithNewKidAcceptedNTimes()

								// Create a new token with the rate limited kid
								// Eventually this should be accepted, which indicates that a discovery poll occurred
								newToken := createIdTokenWithKid(rateLimitedKid)
								expectRequestWithTokenSucceeds(0, newToken)

								expectRequestWithNewKidAcceptedNTimes()
							})
						})

					})
				})

				Context("happy path with default settings (no redis)", func() {
					It("should work", func() {
						ExpectHappyPathToWork(makeSingleRequest, func() {
							Expect(cookies).ToNot(BeEmpty())
							var cookienames []string
							for _, c := range cookies {
								cookienames = append(cookienames, c.Name)
							}
							Expect(cookienames).To(ConsistOf("id_token", "access_token"))
						})
					})
				})

				Context("Oidc tests that don't forward to upstream", func() {
					It("should redirect to auth page", func() {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								// stop at the auth point
								if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
									return http.ErrUseLastResponse
								}
								return nil
							},
						}
						Eventually(func() (string, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							Expect(err).NotTo(HaveOccurred())
							resp, err := client.Do(req)
							if err != nil {
								return "", err
							}
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								return "", err
							}
							fmt.Fprintf(GinkgoWriter, "headers are %v \n", resp.Header)
							return string(body), nil
						}, "10s", "0.5s").Should(Equal("auth"))
					})

					It("should include email scope in url", func() {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
						}
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						Expect(err).NotTo(HaveOccurred())

						Eventually(func() (http.Response, error) {
							r, err := client.Do(req)
							if err != nil {
								return http.Response{}, err
							}
							return *r, err
						}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
							"StatusCode": Equal(http.StatusFound),
							"Header":     HaveKeyWithValue("Location", ContainElement(ContainSubstring("email"))),
						}))
					})

					It("should exchange token", func() {
						finalpage := fmt.Sprintf("http://%s:%d/success", "localhost", envoyPort)
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
						}

						st := oidc.NewStateSigner([]byte(settings.ExtAuthSettings.SigningKey))
						signedState, err := st.Sign(finalpage)
						Expect(err).NotTo(HaveOccurred())
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/callback?code=1234&state="+string(signedState), "localhost", envoyPort), nil)
						Expect(err).NotTo(HaveOccurred())

						Eventually(func() (http.Response, error) {
							r, err := client.Do(req)
							if err != nil {
								return http.Response{}, err
							}
							return *r, err
						}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
							"StatusCode": Equal(http.StatusFound),
							"Header":     HaveKeyWithValue("Location", []string{finalpage}),
						}))
					})

					Context("oidc + opa sanity", func() {
						BeforeEach(func() {
							policy := &gloov1.Artifact{
								Metadata: &core.Metadata{
									Name:      "jwt",
									Namespace: "default",
									Labels:    map[string]string{"team": "infrastructure"},
								},
								Data: map[string]string{
									"jwt.rego": `package test
	
				default allow = false
				allow {
					[header, payload, signature] = io.jwt.decode(input.state.jwt)
					payload["foo"] = "not-bar"
				}
				`}}
							modules := []*core.ResourceRef{{Name: policy.Metadata.Name}}
							options := &extauth.OpaAuthOptions{FastInputConversion: true}

							_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
								Metadata: &core.Metadata{
									Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
									Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
								},
								Configs: []*extauth.AuthConfig_Config{
									{
										AuthConfig: &extauth.AuthConfig_Config_Oauth2{
											Oauth2: &extauth.OAuth2{
												OauthType: getOidcAuthCodeConfig(envoyPort, secret.Metadata.Ref()),
											},
										},
									},
									{
										AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
											OpaAuth: getOpaConfig(modules, options),
										},
									},
								},
							}, clients.WriteOpts{Ctx: ctx})
							Expect(err).NotTo(HaveOccurred())

							proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

							_, err = testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
							Expect(err).ToNot(HaveOccurred())
						})

						It("should NOT allow access", func() {
							EventuallyWithOffset(1, func() (int, error) {
								req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
								req.Header.Add("Authorization", "Bearer "+token)

								resp, err := http.DefaultClient.Do(req)
								if err != nil {
									return 0, err
								}
								return resp.StatusCode, nil
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))

						})

					})
				})

				ExpectUpstreamRequest := func() {
					EventuallyWithOffset(1, func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("Authorization", "Bearer "+token)

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))

					select {
					case r := <-testUpstream.C:
						ExpectWithOffset(1, r.Headers["X-User-Id"]).To(HaveLen(1))
						ExpectWithOffset(1, r.Headers["X-User-Id"][0]).To(Equal("http://localhost:5556;user"))
					case <-time.After(time.Second):
						Fail("expected a message to be received")
					}
				}

				Context("Oidc tests that do forward to upstream", func() {
					It("should allow access with proper jwt token", func() {
						ExpectUpstreamRequest()
					})
				})

				Context("oidc + opa sanity", func() {
					BeforeEach(func() {
						policy := &gloov1.Artifact{
							Metadata: &core.Metadata{
								Name:      "jwt",
								Namespace: "default",
								Labels:    map[string]string{"team": "infrastructure"},
							},
							Data: map[string]string{
								"jwt.rego": `package test

			default allow = false
			allow {
				[header, payload, signature] = io.jwt.decode(input.state.jwt)
				payload["foo"] = "bar"
			}
			`}}
						modules := []*core.ResourceRef{{Name: policy.Metadata.Name, Namespace: policy.Metadata.Namespace}}
						options := &extauth.OpaAuthOptions{FastInputConversion: true}
						_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
								Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{
								{
									AuthConfig: &extauth.AuthConfig_Config_Oauth2{
										Oauth2: &extauth.OAuth2{
											OauthType: getOidcAuthCodeConfig(envoyPort, secret.Metadata.Ref()),
										},
									},
								},
								{
									AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
										OpaAuth: getOpaConfig(modules, options),
									},
								},
							},
						}, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())
						proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

						_, err = testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
						Expect(err).ToNot(HaveOccurred())
					})
					It("should allow access", func() {
						ExpectUpstreamRequest()
					})
				})

			})

			Context("oauth2 token introspection sanity", func() {

				var (
					proxy      *gloov1.Proxy
					authServer *test_utils.AuthServer
					authConfig *extauth.AuthConfig
				)

				// Execute a request with an access token, against an endpoint that requires token authentication
				// and return the status code of the response
				requestWithAccessToken := func(token string) (int, error) {
					getReq, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
					Expect(err).ToNot(HaveOccurred())
					getReq.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

					client := &http.Client{
						CheckRedirect: func(req *http.Request, via []*http.Request) error {
							return http.ErrUseLastResponse
						},
					}
					var resp *http.Response
					resp, err = client.Do(getReq)
					if err != nil {
						return 0, err
					}
					return resp.StatusCode, nil
				}

				JustBeforeEach(func() {
					// Start the auth server
					authServer.Start()

					// Write the auth configuration and ensure it is accepted
					_, err := testClients.AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())
					v1helpers.EventuallyResourceAcceptedWithOffset(1, func() (resources.InputResource, error) {
						return testClients.AuthConfigClient.Read(authConfig.GetMetadata().GetNamespace(), authConfig.GetMetadata().GetName(), clients.ReadOpts{Ctx: ctx})
					})

					// Write the proxy and ensure it is accepted
					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())
					v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
						return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
					})
				})

				AfterEach(func() {
					authServer.Stop()
				})

				Context("using IntrospectionUrl", func() {

					BeforeEach(func() {
						authServer = test_utils.NewAuthServer(
							fmt.Sprintf(":%d", 5556),
							&test_utils.AuthEndpoints{
								TokenIntrospectionEndpoint: "/introspection",
								UserInfoEndpoint:           "/userinfo",
							},
							&test_utils.AuthHandlers{
								// Use default auth handlers
							},
							sets.NewString("valid-access-token"),
							map[string]user_info.UserInfo{})

						authConfig = &extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
								Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_Oauth2{
									Oauth2: &extauth.OAuth2{
										OauthType: getOauthTokenIntrospectionUrlConfig(),
									},
								},
							}},
						}

						proxy = getProxyExtAuthOauthTokenIntrospection(envoyPort, testUpstream.Upstream.Metadata.Ref())
					})

					It("should accept introspection url with valid access token", func() {
						Eventually(func() (int, error) {
							return requestWithAccessToken("valid-access-token")
						}, "5s", "0.5s").Should(Equal(http.StatusOK))
						Consistently(func() (int, error) {
							return requestWithAccessToken("valid-access-token")
						}, "3s", "0.5s").Should(Equal(http.StatusOK))
					})

					It("should deny introspection url with invalid access token", func() {
						Eventually(func() (int, error) {
							return requestWithAccessToken("invalid-access-token")
						}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
						Consistently(func() (int, error) {
							return requestWithAccessToken("invalid-access-token")
						}, "3s", "0.5s").Should(Equal(http.StatusForbidden))
					})

				})

				Context("using Introspection", func() {

					var (
						secret *gloov1.Secret
					)

					createBasicAuthHandler := func(validToken, clientId, clientSecret string) func(http.ResponseWriter, *http.Request) {

						return func(writer http.ResponseWriter, request *http.Request) {
							err := request.ParseForm()
							if err != nil {
								panic(err)
							}

							requestedToken := request.Form.Get("token")
							requestedClientId := request.Form.Get("client_id")
							requestedClientSecret := request.Form.Get("client_secret")

							response := &token_validation.IntrospectionResponse{}

							// Request is only validated if all criteria match
							if validToken == requestedToken && requestedClientId == clientId && requestedClientSecret == clientSecret {
								response.Active = true
							}

							bytes, err := json.Marshal(response)
							if err != nil {
								panic(err)
							}
							writer.Write(bytes)
						}
					}

					getOauthSecret := func(name, value string) *gloov1.Secret {
						clientSecret := &extauth.OauthSecret{
							ClientSecret: value,
						}
						return &gloov1.Secret{
							Metadata: &core.Metadata{
								Name:      name,
								Namespace: "default",
							},
							Kind: &gloov1.Secret_Oauth{
								Oauth: clientSecret,
							},
						}
					}

					BeforeEach(func() {
						// Create the client secret
						secret = getOauthSecret("secret", "client-secret")
						_, err := testClients.SecretClient.Write(secret, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())

						// Create an auth server that requires clients to provide credentials (client-id and client-secret)
						authServer = test_utils.NewAuthServer(
							fmt.Sprintf(":%d", 5556),
							&test_utils.AuthEndpoints{
								TokenIntrospectionEndpoint: "/introspection",
								UserInfoEndpoint:           "/userinfo",
							},
							&test_utils.AuthHandlers{
								TokenIntrospectionHandler: createBasicAuthHandler("valid-access-token", "client-id", "client-secret"),
							},
							sets.NewString("valid-access-token"),
							map[string]user_info.UserInfo{})

						// Create an auth config, with proper references to the client credentials
						authConfig = &extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
								Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_Oauth2{
									Oauth2: &extauth.OAuth2{
										OauthType: getOauthTokenIntrospectionConfig("client-id", secret.Metadata.Ref()),
									},
								},
							}},
						}

						proxy = getProxyExtAuthOauthTokenIntrospection(envoyPort, testUpstream.Upstream.Metadata.Ref())
					})

					When("auth config includes valid credentials", func() {
						// The default auth config that we initialize is valid

						It("should accept introspection with valid access token", func() {
							Eventually(func() (int, error) {
								return requestWithAccessToken("valid-access-token")
							}, "5s", "0.5s").Should(Equal(http.StatusOK))
							Consistently(func() (int, error) {
								return requestWithAccessToken("valid-access-token")
							}, "3s", "0.5s").Should(Equal(http.StatusOK))
						})

						It("should deny introspection with invalid access token", func() {
							Eventually(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
							Consistently(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "3s", "0.5s").Should(Equal(http.StatusForbidden))
						})
					})

					When("auth config includes invalid credentials", func() {

						var invalidSecret *gloov1.Secret

						BeforeEach(func() {
							// Create a client secret with the wrong value
							invalidSecret = getOauthSecret("invalid-secret", "invalid-client-secret")
							_, err := testClients.SecretClient.Write(invalidSecret, clients.WriteOpts{})
							Expect(err).NotTo(HaveOccurred())

							// Set the auth config to reference that invalid secret
							authConfig = &extauth.AuthConfig{
								Metadata: &core.Metadata{
									Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
									Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
								},
								Configs: []*extauth.AuthConfig_Config{{
									AuthConfig: &extauth.AuthConfig_Config_Oauth2{
										Oauth2: &extauth.OAuth2{
											OauthType: getOauthTokenIntrospectionConfig("client-id", invalidSecret.Metadata.Ref()),
										},
									},
								}},
							}
						})

						It("should deny introspection with valid access token", func() {
							Eventually(func() (int, error) {
								return requestWithAccessToken("valid-access-token")
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
							Consistently(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "3s", "0.5s").Should(Equal(http.StatusForbidden))
						})

						It("should deny introspection with invalid access token", func() {
							Eventually(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
							Consistently(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "3s", "0.5s").Should(Equal(http.StatusForbidden))
						})
					})

					When("auth config is missing credentials and auth server doesn't require credentials", func() {

						BeforeEach(func() {
							authServer = test_utils.NewAuthServer(
								fmt.Sprintf(":%d", 5556),
								&test_utils.AuthEndpoints{
									TokenIntrospectionEndpoint: "/introspection",
									UserInfoEndpoint:           "/userinfo",
								},
								&test_utils.AuthHandlers{
									// Use the default handlers, which do not require client credentials in the introspection request
								},
								sets.NewString("valid-access-token"),
								map[string]user_info.UserInfo{})

							authConfig = &extauth.AuthConfig{
								Metadata: &core.Metadata{
									Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
									Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
								},
								Configs: []*extauth.AuthConfig_Config{{
									AuthConfig: &extauth.AuthConfig_Config_Oauth2{
										Oauth2: &extauth.OAuth2{
											OauthType: getOauthTokenIntrospectionConfig("", nil),
										},
									},
								}},
							}
						})

						It("should accept introspection with valid access token", func() {
							Eventually(func() (int, error) {
								return requestWithAccessToken("valid-access-token")
							}, "5s", "0.5s").Should(Equal(http.StatusOK))
							Consistently(func() (int, error) {
								return requestWithAccessToken("valid-access-token")
							}, "3s", "0.5s").Should(Equal(http.StatusOK))
						})

						It("should deny introspection with invalid access token", func() {
							Eventually(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))
							Consistently(func() (int, error) {
								return requestWithAccessToken("invalid-access-token")
							}, "3s", "0.5s").Should(Equal(http.StatusForbidden))
						})
					})

				})

			})

			Context("api key sanity tests", func() {
				BeforeEach(func() {
					_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
						Metadata: &core.Metadata{
							Name:      getApiKeyExtAuthExtension().GetConfigRef().Name,
							Namespace: getApiKeyExtAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
								ApiKeyAuth: getApiKeyAuthConfig(),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					apiKeySecret1 := &extauth.ApiKeySecret{
						ApiKey: "secretApiKey1",
					}

					secret1 := &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      "secret1",
							Namespace: "default",
						},
						Kind: &gloov1.Secret_ApiKey{
							ApiKey: apiKeySecret1,
						},
					}

					apiKeySecret2 := &extauth.ApiKeySecret{
						ApiKey: "secretApiKey2",
					}

					secret2 := &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      "secret2",
							Namespace: "default",
							Labels:    map[string]string{"team": "infrastructure"},
						},
						Kind: &gloov1.Secret_ApiKey{
							ApiKey: apiKeySecret2,
						},
					}

					_, err = testClients.SecretClient.Write(secret1, clients.WriteOpts{})
					Expect(err).ToNot(HaveOccurred())

					_, err = testClients.SecretClient.Write(secret2, clients.WriteOpts{})
					Expect(err).ToNot(HaveOccurred())

					proxy := getProxyExtAuthApiKeyAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
						return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
					})
				})

				It("should deny ext auth envoy without apikey", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})

				It("should deny ext auth envoy with incorrect apikey", func() {
					Eventually(func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("api-key", "badApiKey")
						resp, err := http.DefaultClient.Do(req)

						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})

				It("should accept ext auth envoy with correct apikey -- secret ref match", func() {
					Eventually(func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("api-key", "secretApiKey1")
						resp, err := http.DefaultClient.Do(req)

						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})

				It("should accept ext auth envoy with correct apikey -- label match", func() {
					Eventually(func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("api-key", "secretApiKey2")
						resp, err := http.DefaultClient.Do(req)

						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})
			})

			Context("http passthrough", func() {

				expectStatusCodeWithHeaders := func(responseCode int, reqHeadersToAdd map[string][]string, responseHeadersToExpect map[string]string) *http.Response {
					var resp *http.Response
					EventuallyWithOffset(1, func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						if err != nil {
							return 0, err
						}
						req.Header = reqHeadersToAdd
						resp, err = http.DefaultClient.Do(req)
						if err != nil {
							return 0, err
						}
						for headerName, headerValue := range responseHeadersToExpect {
							ExpectWithOffset(2, resp.Header.Get(headerName)).To(Equal(headerValue))
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(responseCode))
					return resp
				}

				expectStatusCodeWithBody := func(responseCode int, body string) *http.Response {
					var (
						resp *http.Response
						err  error
					)
					EventuallyWithOffset(1, func() (int, error) {
						resp, err = http.Post(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), "text/plain", strings.NewReader(body))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}).Should(Equal(responseCode))
					return resp
				}

				Context("passthrough sanity", func() {
					var (
						proxy                 *gloov1.Proxy
						httpAuthServer        *v1helpers.TestUpstream
						httpPassthroughConfig *extauth.PassThroughHttp
						handler               v1helpers.ExtraHandlerFunc
						authconfigCfg         *structpb.Struct
						authConfigRequestPath string
						protocol              string
					)

					BeforeEach(func() {
						httpPassthroughConfig = &extauth.PassThroughHttp{
							Request:  &extauth.PassThroughHttp_Request{},
							Response: &extauth.PassThroughHttp_Response{},
							ConnectionTimeout: &duration.Duration{
								Seconds: 10,
							},
						}
						authconfigCfg = nil
						handler = nil
						authConfigRequestPath = ""
						protocol = "http"
					})

					JustBeforeEach(func() {
						httpAuthServer = v1helpers.NewTestHttpUpstreamWithHandler(ctx, "127.0.0.1", handler)
						up := httpAuthServer.Upstream
						_, err := testClients.UpstreamClient.Write(up, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())
						httpPassthroughConfig.Url = fmt.Sprintf("%s://%s%s", protocol, httpAuthServer.Address, authConfigRequestPath)
						ac := &extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
								Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
									PassThroughAuth: &extauth.PassThroughAuth{
										Protocol: &extauth.PassThroughAuth_Http{
											Http: httpPassthroughConfig,
										},
										Config: authconfigCfg,
									},
								},
							}},
						}
						_, err = testClients.AuthConfigClient.Write(ac, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())

						// ensure auth config is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.AuthConfigClient.Read(ac.Metadata.Namespace, ac.Metadata.Name, clients.ReadOpts{})
						})

						proxy = getProxyExtAuthPassThroughAuth(envoyPort, testUpstream.Upstream.Metadata.Ref(), false)
						_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})

						// ensure proxy is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						})
					})

					It("works", func() {

						expectStatusCodeWithHeaders(200, nil, nil)
						select {
						case received := <-httpAuthServer.C:
							Expect(received.Method).To(Equal("POST"))
						case <-time.After(time.Second * 5):
							Fail("request didn't make it upstream")
						}
					})

					Context("setting path on URL", func() {
						BeforeEach(func() {
							authConfigRequestPath += "/auth"
						})
						It("has correct path to auth server", func() {
							expectStatusCodeWithHeaders(200, nil, nil)
							select {
							case received := <-httpAuthServer.C:
								Expect(received.Method).To(Equal("POST"))
								Expect(received.URL.Path).To(Equal("/auth"))
							case <-time.After(time.Second * 5):
								Fail("request didn't make it upstream")
							}
						})
					})

					Context("Request", func() {
						BeforeEach(func() {
							httpPassthroughConfig.Request = &extauth.PassThroughHttp_Request{
								AllowedHeaders: []string{"x-passthrough-1", "x-passthrough-2"},
								HeadersToAdd: map[string]string{
									"x-added-header-1": "net new header",
								},
							}
						})
						It("copies `allowed_headers` request headers and adds `headers_to_add` headers to auth request", func() {
							expectStatusCodeWithHeaders(200, map[string][]string{
								"x-passthrough-1":    {"some header from request"},
								"x-passthrough-2":    {"some header from request 2"},
								"x-dont-passthrough": {"some header from request that shouldn't be passed through to auth server"},
							}, nil)
							select {
							case received := <-httpAuthServer.C:
								Expect(received.Method).To(Equal("POST"))
								Expect(received.Headers.Get("X-Passthrough-1")).To(Equal("some header from request"))
								Expect(received.Headers.Get("X-Passthrough-2")).To(Equal("some header from request 2"))
								Expect(received.Headers["X-Dont-Passthrough-1"]).To(BeNil())

								Expect(received.Headers.Get("x-added-header-1")).To(Equal("net new header"))
							case <-time.After(time.Second * 5):
								Fail("request didn't make it upstream")
							}
						})
					})

					Context("Response", func() {
						BeforeEach(func() {
							httpPassthroughConfig.Response = &extauth.PassThroughHttp_Response{
								AllowedUpstreamHeaders:       []string{"x-auth-header-1", "x-auth-header-2"},
								AllowedClientHeadersOnDenied: []string{"x-auth-header-1"},
							}
						})
						Context("On authorized response", func() {
							BeforeEach(func() {
								handler = func(rw http.ResponseWriter, r *http.Request) bool {
									rw.Header().Set("x-auth-header-1", "some value")
									rw.Header().Set("x-auth-header-2", "some value 2")
									rw.Header().Set("x-shouldnt-upstream", "shouldn't upstream")
									return true
								}
							})
							It("copies `allowed_headers` request headers and adds `headers_to_add` headers to auth request", func() {
								expectStatusCodeWithHeaders(200, map[string][]string{
									"x-auth-header-1": {"hello"},
									"x-passthrough-1": {"some header from request that should go to upstream"},
								}, nil)
								select {
								case received := <-testUpstream.C:
									Expect(received.Method).To(Equal("GET"))
									// This header should have an appended value since it exists on the original request
									Expect(received.Headers.Get("x-auth-header-1")).To(Equal("hello,some value"))
									Expect(received.Headers.Get("x-auth-header-2")).To(Equal("some value 2"))
									Expect(received.Headers.Get("x-shouldnt-upstream")).To(BeEmpty())

								case <-time.After(time.Second * 5):
									Fail("request didn't make it upstream")
								}
							})
						})

						Context("on authorized response", func() {
							BeforeEach(func() {
								handler = func(rw http.ResponseWriter, r *http.Request) bool {
									rw.Header().Set("x-auth-header-1", "some value")
									rw.WriteHeader(http.StatusUnauthorized)
									return true
								}
							})
							It("sends allowed authorization headers back to downstream", func() {
								expectStatusCodeWithHeaders(http.StatusUnauthorized, nil, map[string]string{"x-auth-header-1": "some value"})
							})
						})
					})

					Context("Request to Auth Server Body", func() {
						BeforeEach(func() {
							// We need these settings so envoy buffers the request body and sends it to the ext-auth-service
							glooSettings.Extauth.RequestBody = &extauth.BufferSettings{
								MaxRequestBytes:     uint32(1024),
								AllowPartialMessage: true,
							}
						})
						Context("passes through http request body", func() {
							BeforeEach(func() {
								httpPassthroughConfig.Request.PassThroughBody = true
							})
							It("correctly", func() {
								expectStatusCodeWithBody(http.StatusOK, "some body")
								select {
								case received := <-httpAuthServer.C:
									Expect(string(received.Body)).To(Equal(`{"body":"some body"}`))
								case <-time.After(time.Second * 5):
									Fail("request didn't make it upstream")
								}
							})
						})

						Context("doesn't pass through http request body", func() {
							BeforeEach(func() {
								httpPassthroughConfig.Request.PassThroughBody = false
							})
							It("body is empty if no body, config, state, or filtermetadata passthrough is set", func() {
								expectStatusCodeWithBody(http.StatusOK, "some body")
								select {
								case received := <-httpAuthServer.C:
									Expect(string(received.Body)).To(BeEmpty())
								case <-time.After(time.Second * 5):
									Fail("request didn't make it upstream")
								}
							})
						})
					})

					Context("pass config specified on auth config in auth request body", func() {

						BeforeEach(func() {
							authconfigCfg = &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"nestedStruct": {
										Kind: &structpb.Value_StructValue{
											StructValue: &structpb.Struct{
												Fields: map[string]*structpb.Value{
													"list": {
														Kind: &structpb.Value_ListValue{
															ListValue: &structpb.ListValue{
																Values: []*structpb.Value{
																	{
																		Kind: &structpb.Value_StringValue{
																			StringValue: "some string",
																		},
																	},
																	{
																		Kind: &structpb.Value_NumberValue{
																			NumberValue: float64(23),
																		},
																	},
																},
															},
														},
													},
													"string": {
														Kind: &structpb.Value_StringValue{
															StringValue: "some string",
														},
													},
													"int": {
														Kind: &structpb.Value_NumberValue{
															NumberValue: float64(23),
														},
													},
													"bool": {
														Kind: &structpb.Value_BoolValue{
															BoolValue: true,
														},
													},
												},
											},
										},
									},
								},
							}
						})
						It("passes through config from auth config to passthrough server", func() {
							type ConfigStruct struct {
								Config *structpb.Struct `json:"config"`
							}
							expectStatusCodeWithHeaders(http.StatusOK, nil, nil)
							select {
							case received := <-httpAuthServer.C:
								cfgStruct := &ConfigStruct{}
								err := json.Unmarshal(received.Body, cfgStruct)
								Expect(err).NotTo(HaveOccurred())
								Expect(cfgStruct.Config).To(matchers.MatchProto(authconfigCfg))
							case <-time.After(time.Second * 5):
								Fail("request didn't make it upstream")
							}
						})
					})

				})

				Context("http passthrough chaining sanity", func() {
					var (
						proxy      *gloov1.Proxy
						authConfig *extauth.AuthConfig
					)

					var (
						httpAuthServerA,
						httpAuthServerB *v1helpers.TestUpstream
						httpPassthroughConfigA,
						httpPassthroughConfigB *extauth.PassThroughHttp
						handlerA,
						handlerB v1helpers.ExtraHandlerFunc
					)

					BeforeEach(func() {
						httpPassthroughConfigA = &extauth.PassThroughHttp{
							Request:  &extauth.PassThroughHttp_Request{},
							Response: &extauth.PassThroughHttp_Response{},
							ConnectionTimeout: &duration.Duration{
								Seconds: 10,
							},
						}
						httpPassthroughConfigB = &extauth.PassThroughHttp{
							Request:  &extauth.PassThroughHttp_Request{},
							Response: &extauth.PassThroughHttp_Response{},
							ConnectionTimeout: &duration.Duration{
								Seconds: 10,
							},
						}
						authConfig = &extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
								Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
									PassThroughAuth: &extauth.PassThroughAuth{
										Protocol: &extauth.PassThroughAuth_Http{
											Http: httpPassthroughConfigA,
										},
									},
								},
							},
								{
									AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
										PassThroughAuth: &extauth.PassThroughAuth{
											Protocol: &extauth.PassThroughAuth_Http{
												Http: httpPassthroughConfigB,
											},
										},
									},
								},
							},
						}
						handlerA, handlerB = nil, nil
					})

					JustBeforeEach(func() {
						httpAuthServerA = v1helpers.NewTestHttpUpstreamWithHandler(ctx, "127.0.0.1", handlerA)
						up := httpAuthServerA.Upstream
						_, err := testClients.UpstreamClient.Write(up, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())
						httpPassthroughConfigA.Url = fmt.Sprintf("http://%s", httpAuthServerA.Address)
						httpAuthServerB = v1helpers.NewTestHttpUpstreamWithHandler(ctx, "127.0.0.1", handlerB)
						up = httpAuthServerB.Upstream
						_, err = testClients.UpstreamClient.Write(up, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())
						httpPassthroughConfigB.Url = fmt.Sprintf("http://%s", httpAuthServerB.Address)

						_, err = testClients.AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())

						// ensure auth config is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.AuthConfigClient.Read(authConfig.Metadata.Namespace, authConfig.Metadata.Name, clients.ReadOpts{})
						})

						proxy = getProxyExtAuthPassThroughAuth(envoyPort, testUpstream.Upstream.Metadata.Ref(), false)
						_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})

						// ensure proxy is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						})
					})

					_ = func(responseCode int, body string) *http.Response {
						var (
							resp *http.Response
							err  error
						)
						EventuallyWithOffset(1, func() (int, error) {
							resp, err = http.Post(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), "text/plain", strings.NewReader(body))
							if err != nil {
								return 0, err
							}
							return resp.StatusCode, nil
						}).Should(Equal(responseCode))
						return resp
					}

					It("works", func() {
						expectStatusCodeWithHeaders(200, nil, nil)
					})
					Context("can modify state", func() {
						BeforeEach(func() {
							httpPassthroughConfigA.Request.PassThroughState = true
							httpPassthroughConfigA.Response.ReadStateFromResponse = true
							httpPassthroughConfigB.Request.PassThroughState = true
							handlerA = func(rw http.ResponseWriter, r *http.Request) bool {
								rw.Write([]byte(`{"state":{"list": ["item1", "item2", 3, {"item4":""}], "string": "hello", "integer": 9, "nestedObject":{"key":"value"}}}`))
								return false
							}
						})
						It("modifies state in authServerA and authServerB can see the new state", func() {
							expectStatusCodeWithHeaders(200, nil, nil)

							select {
							case received := <-httpAuthServerB.C:
								Expect(string(received.Body)).To(Equal(`{"state":{"integer":9,"list":["item1","item2",3,{"item4":""}],"nestedObject":{"key":"value"},"string":"hello"}}`))
							case <-time.After(time.Second * 5):
								Fail("request didn't make it upstream")
							}
						})
					})
				})

				Context("https", func() {
					var (
						proxy                 *gloov1.Proxy
						httpAuthServer        *v1helpers.TestUpstream
						httpPassthroughConfig *extauth.PassThroughHttp
						handler               v1helpers.ExtraHandlerFunc
						authconfigCfg         *structpb.Struct
						protocol              string
						rootCaBytes           []byte
					)

					BeforeEach(func() {
						httpPassthroughConfig = &extauth.PassThroughHttp{
							Request:  &extauth.PassThroughHttp_Request{},
							Response: &extauth.PassThroughHttp_Response{},
							ConnectionTimeout: &duration.Duration{
								Seconds: 10,
							},
						}
						authconfigCfg = nil
						handler = nil
						protocol = "https"
					})

					JustBeforeEach(func() {
						rootCaBytes, httpAuthServer = v1helpers.NewTestHttpsUpstreamWithHandler(ctx, "127.0.0.1", handler)
						// set environment variable for ext auth server passthrough https
						err := os.Setenv(translation.HttpsPassthroughCaCert, base64.StdEncoding.EncodeToString(rootCaBytes))
						Expect(err).NotTo(HaveOccurred())
						up := httpAuthServer.Upstream
						_, err = testClients.UpstreamClient.Write(up, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())
						httpPassthroughConfig.Url = fmt.Sprintf("%s://%s", protocol, httpAuthServer.Address)
						ac := &extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
								Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
									PassThroughAuth: &extauth.PassThroughAuth{
										Protocol: &extauth.PassThroughAuth_Http{
											Http: httpPassthroughConfig,
										},
										Config: authconfigCfg,
									},
								},
							}},
						}
						_, err = testClients.AuthConfigClient.Write(ac, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())

						// ensure auth config is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.AuthConfigClient.Read(ac.Metadata.Namespace, ac.Metadata.Name, clients.ReadOpts{})
						})

						proxy = getProxyExtAuthPassThroughAuth(envoyPort, testUpstream.Upstream.Metadata.Ref(), false)
						_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})

						// ensure proxy is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						})
					})

					It("works", func() {

						expectStatusCodeWithHeaders(200, nil, nil)
						select {
						case received := <-httpAuthServer.C:
							Expect(received.Method).To(Equal("POST"))
						case <-time.After(time.Second * 5):
							Fail("request didn't make it upstream")
						}
					})
					AfterEach(func() {
						err := os.Unsetenv(translation.HttpsPassthroughCaCert)
						Expect(err).NotTo(HaveOccurred())
					})
				})
			})

			Context("grpc passthrough", func() {

				Context("passthrough sanity", func() {
					var (
						proxy          *gloov1.Proxy
						authServer     *passthrough_test_utils.GrpcAuthServer
						authServerPort = 5556
						zipkinTracing  bool
					)

					expectRequestEventuallyReturnsResponseCode := func(responseCode int) {
						EventuallyWithOffset(1, func() (int, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							if err != nil {
								return 0, err
							}

							resp, err := http.DefaultClient.Do(req)
							if err != nil {
								return 0, err
							}

							return resp.StatusCode, nil
						}, "5s", "0.5s").Should(Equal(responseCode))
					}

					JustBeforeEach(func() {
						// start auth server
						err := authServer.Start(authServerPort)
						Expect(err).NotTo(HaveOccurred())

						// write auth configuration
						_, err = testClients.AuthConfigClient.Write(&extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
								Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
									PassThroughAuth: getPassThroughAuthConfig(authServer.GetAddress()),
								},
							}},
						}, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())

						// get proxy with pass through auth extension
						proxy = getProxyExtAuthPassThroughAuth(envoyPort, testUpstream.Upstream.Metadata.Ref(), zipkinTracing)

						// write proxy
						_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())

						// ensure proxy is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						})
					})

					AfterEach(func() {
						authServer.Stop()
					})

					Context("when auth server returns ok response", func() {

						BeforeEach(func() {
							authServerResponse := passthrough_test_utils.OkResponse()
							authServer = passthrough_test_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
						})

						It("should accept extauth passthrough", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusOK)
						})

					})

					Context("when auth server returns denied response", func() {

						BeforeEach(func() {
							authServerResponse := passthrough_test_utils.DeniedResponse()
							authServer = passthrough_test_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
						})

						It("should deny extauth passthrough", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
						})

					})

					Context("when auth server errors", func() {

						BeforeEach(func() {
							authServerError := errors.New("auth server internal server error")
							authServer = passthrough_test_utils.NewGrpcAuthServerWithResponse(nil, authServerError)
						})

						It("should deny extauth passthrough", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusForbidden)
						})

					})

					Context("when auth server returns ok response with valid dynamic metadata properties", func() {

						BeforeEach(func() {
							authServerResponse := passthrough_test_utils.OkResponseWithDynamicMetadata(&structpb.Struct{
								Fields: map[string]*structpb.Value{
									"current-state-key": {
										Kind: &structpb.Value_StringValue{
											StringValue: "new-state-value",
										},
									},
									"new-state-key": {
										Kind: &structpb.Value_StringValue{
											StringValue: "new-state-value",
										},
									},
								},
							})
							authServer = passthrough_test_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
						})

						It("should accept extauth passthrough", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusOK)
						})

					})

					Context("when auth server returns ok response when tracing metadata is present", func() {

						BeforeEach(func() {
							// create zipkin upstream
							zipkinUs := &gloov1.Upstream{
								Metadata: &core.Metadata{
									Name:      "zipkin",
									Namespace: "default",
								},
								UpstreamType: &gloov1.Upstream_Static{
									Static: &gloov1static.UpstreamSpec{
										Hosts: []*gloov1static.Host{
											{
												Addr: envoyInstance.LocalAddr(),
												Port: 9411,
											},
										},
									},
								},
							}
							_, err := testClients.UpstreamClient.Write(zipkinUs, clients.WriteOpts{})
							Expect(err).NotTo(HaveOccurred())

							zipkinTracing = true
							authServer = passthrough_test_utils.NewGrpcAuthServerWithTracingRequired()
						})

						AfterEach(func() {
							zipkinTracing = false
						})

						It("should accept extauth passthrough", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusOK)
						})

					})

				})

				Context("passthrough chaining sanity", func() {
					// These tests are used to validate that state is passed properly to and from the passthrough service

					var (
						proxy           *gloov1.Proxy
						authServerA     *passthrough_test_utils.GrpcAuthServer
						authServerAPort = 5556

						authServerB     *passthrough_test_utils.GrpcAuthServer
						authServerBPort = 5557
					)

					expectRequestEventuallyReturnsResponseCode := func(responseCode int) {
						EventuallyWithOffset(1, func() (int, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							if err != nil {
								return 0, nil
							}

							resp, err := http.DefaultClient.Do(req)
							if err != nil {
								return 0, err
							}

							return resp.StatusCode, nil
						}, "5s", "0.5s").Should(Equal(responseCode))
					}

					JustBeforeEach(func() {
						// start auth servers
						err := authServerA.Start(authServerAPort)
						Expect(err).NotTo(HaveOccurred())

						err = authServerB.Start(authServerBPort)
						Expect(err).NotTo(HaveOccurred())

						// write auth configuration
						_, err = testClients.AuthConfigClient.Write(&extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
								Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{
								{
									// Ordering is important here, AuthServerA is listed first so it is earlier in the chain
									AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
										PassThroughAuth: getPassThroughAuthConfig(authServerA.GetAddress()),
									},
								},
								{
									AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
										PassThroughAuth: getPassThroughAuthConfig(authServerB.GetAddress()),
									},
								},
							},
						}, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())

						// get proxy with pass through auth extension
						proxy = getProxyExtAuthPassThroughAuth(envoyPort, testUpstream.Upstream.Metadata.Ref(), false)

						// write proxy
						_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())

						// ensure proxy is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						})
					})

					AfterEach(func() {
						authServerA.Stop()
						authServerB.Stop()
					})

					Context("first auth server writes metadata, second requires it", func() {

						BeforeEach(func() {
							// Configure AuthServerA (first in chain) to return DynamicMetadata.
							authServerAResponse := passthrough_test_utils.OkResponseWithDynamicMetadata(&structpb.Struct{
								Fields: map[string]*structpb.Value{
									"key": {
										Kind: &structpb.Value_StringValue{
											StringValue: "value",
										},
									},
									"non-string-value": {
										Kind: &structpb.Value_StructValue{
											StructValue: &structpb.Struct{
												Fields: map[string]*structpb.Value{
													"nested-key": {
														Kind: &structpb.Value_StringValue{
															StringValue: "nested-value",
														},
													},
												},
											},
										},
									},
								},
							})
							authServerA = passthrough_test_utils.NewGrpcAuthServerWithResponse(authServerAResponse, nil)

							// Configure AuthServerB (second in chain) to expect those dynamic metadata keys
							authServerB = passthrough_test_utils.NewGrpcAuthServerWithRequiredMetadata([]string{
								"key",
								"non-string-value",
							})
						})

						It("should accept extauth passthrough", func() {
							// This will pass only if the following events occur:
							//		1. AuthServerA returns DynamicMetadata under PassThrough Key and that data is stored on AuthorizationRequest
							//		2. State on AuthorizationRequest is parsed and sent on subsequent request to AuthServerB
							//		3. AuthServerB receives the Metadata and returns ok if all keys are present.
							expectRequestEventuallyReturnsResponseCode(http.StatusOK)
						})

					})

					Context("first auth server does not write metadata, second requires it", func() {

						BeforeEach(func() {
							// Configure AuthServerA (first in chain) to NOT return DynamicMetadata.
							authServerAResponse := passthrough_test_utils.OkResponse()
							authServerA = passthrough_test_utils.NewGrpcAuthServerWithResponse(authServerAResponse, nil)

							// Configure AuthServerB (second in chain) to expect dynamic metadata keys
							authServerB = passthrough_test_utils.NewGrpcAuthServerWithRequiredMetadata([]string{
								"key",
								"non-string-value",
							})
						})

						It("should deny extauth passthrough", func() {
							// This will deny the request because:
							//		1. AuthServerA does not return DynamicMetadata under PassThrough Key. So there is not AuthorizationRequest.State
							//		2. Since there is no AuthorizationRequest.State, no Metadata is sent in request to AuthServerB
							//		3. AuthServerB receives no Metadata, but requires certain fields and returns 401 since there are missing properties
							expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
						})

					})

				})

				Context("passthrough auth config sanity", func() {
					// These tests are used to validate that custom config is passed properly to the passthrough service

					var (
						proxy           *gloov1.Proxy
						authServerA     *passthrough_test_utils.GrpcAuthServer
						authServerAPort = 5556
					)

					expectRequestEventuallyReturnsResponseCode := func(responseCode int) {
						EventuallyWithOffset(1, func() (int, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							if err != nil {
								return 0, nil
							}

							resp, err := http.DefaultClient.Do(req)
							if err != nil {
								return 0, err
							}

							return resp.StatusCode, nil
						}, "5s", "0.5s").Should(Equal(responseCode))
					}

					newGrpcAuthServerwithRequiredConfig := func() *passthrough_test_utils.GrpcAuthServer {
						return &passthrough_test_utils.GrpcAuthServer{
							AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
								// Check if config exists in the FilterMetadata under the MetadataConfigKey.
								if passThroughFilterMetadata, ok := req.GetAttributes().GetMetadataContext().GetFilterMetadata()[grpcPassthrough.MetadataConfigKey]; ok {
									passThroughFields := passThroughFilterMetadata.GetFields()
									if value, ok := passThroughFields["customConfig1"]; ok && value.GetBoolValue() == true {
										// Required key was in FilterMetadata, succeed request
										return passthrough_test_utils.OkResponse(), nil
									}
									// Required key was not in FilterMetadata, deny fail request
									return passthrough_test_utils.DeniedResponse(), nil
								}
								// No passthrough properties were sent in FilterMetadata, fail request
								return passthrough_test_utils.DeniedResponse(), nil
							},
						}
					}

					JustBeforeEach(func() {
						// start auth server
						err := authServerA.Start(authServerAPort)
						Expect(err).NotTo(HaveOccurred())

						authConfig := &extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
								Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{{
								AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
									PassThroughAuth: getPassThroughAuthWithCustomConfig(authServerA.GetAddress()),
								},
							}},
						}
						// write auth configuration
						_, err = testClients.AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())

						// get proxy with pass through auth extension
						proxy = getProxyExtAuthPassThroughAuth(envoyPort, testUpstream.Upstream.Metadata.Ref(), false)

						// write proxy
						_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
						Expect(err).NotTo(HaveOccurred())

						// ensure proxy is accepted
						v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
							return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						})
					})

					AfterEach(func() {
						authServerA.Stop()
					})

					Context("passes config block to passthrough auth service", func() {
						BeforeEach(func() {
							authServerA = newGrpcAuthServerwithRequiredConfig()
						})

						It("correctly", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusOK)
						})
					})
				})
			})
		})

		Context("using old config format", func() {

			Context("oidc sanity", func() {
				var (
					privateKey      *rsa.PrivateKey
					discoveryServer fakeDiscoveryServer
					secret          *gloov1.Secret
					proxy           *gloov1.Proxy
					token           string
				)

				// create an id token with a particular key id
				createIdTokenWithKid := func(kid string) string {
					tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
						"foo": "bar",
						"aud": "test-clientid",
						"sub": "user",
						"iss": "http://localhost:5556",
					})
					tokenToSign.Header["kid"] = kid
					idToken, err := tokenToSign.SignedString(privateKey)
					ExpectWithOffset(1, err).NotTo(HaveOccurred())

					return idToken
				}

				BeforeEach(func() {
					discoveryServer = fakeDiscoveryServer{}
					privateKey = discoveryServer.Start()

					clientSecret := &extauth.OauthSecret{
						ClientSecret: "test",
					}

					secret = &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      "secret",
							Namespace: "default",
						},
						Kind: &gloov1.Secret_Oauth{
							Oauth: clientSecret,
						},
					}
					_, err := testClients.SecretClient.Write(secret, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					_, err = testClients.AuthConfigClient.Write(&extauth.AuthConfig{
						Metadata: &core.Metadata{
							Name:      getOidcExtAuthExtension().GetConfigRef().Name,
							Namespace: getOidcExtAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_Oauth{
								Oauth: getOauthConfig(envoyPort, secret.Metadata.Ref()),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					proxy = getProxyExtAuthOIDC(envoyPort, testUpstream.Upstream.Metadata.Ref())

					token = createIdTokenWithKid(discoveryServer.getValidKeyId())
				})

				JustBeforeEach(func() {
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					v1helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
						return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
					})
				})

				AfterEach(func() {
					discoveryServer.Stop()
				})

				Context("Oidc tests that don't forward to upstream", func() {
					It("should redirect to auth page", func() {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								// stop at the auth point
								if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
									return http.ErrUseLastResponse
								}
								return nil
							},
						}
						Eventually(func() (string, error) {
							req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
							Expect(err).NotTo(HaveOccurred())
							resp, err := client.Do(req)
							if err != nil {
								return "", err
							}
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								return "", err
							}
							return string(body), nil
						}, "10s", "0.5s").Should(Equal("auth"))
					})

					It("should include email scope in url", func() {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
						}
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						Expect(err).NotTo(HaveOccurred())

						Eventually(func() (http.Response, error) {
							r, err := client.Do(req)
							if err != nil {
								return http.Response{}, err
							}
							return *r, err
						}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
							"StatusCode": Equal(http.StatusFound),
							"Header":     HaveKeyWithValue("Location", ContainElement(ContainSubstring("email"))),
						}))
					})

					It("should exchange token", func() {
						finalpage := fmt.Sprintf("http://%s:%d/success", "localhost", envoyPort)
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
						}

						st := oidc.NewStateSigner([]byte(settings.ExtAuthSettings.SigningKey))
						signedState, err := st.Sign(finalpage)
						Expect(err).NotTo(HaveOccurred())
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/callback?code=1234&state="+string(signedState), "localhost", envoyPort), nil)
						Expect(err).NotTo(HaveOccurred())

						Eventually(func() (http.Response, error) {
							r, err := client.Do(req)
							if err != nil {
								return http.Response{}, err
							}
							return *r, err
						}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
							"StatusCode": Equal(http.StatusFound),
							"Header":     HaveKeyWithValue("Location", []string{finalpage}),
						}))
					})

					Context("oidc + opa sanity", func() {
						BeforeEach(func() {
							policy := &gloov1.Artifact{
								Metadata: &core.Metadata{
									Name:      "jwt",
									Namespace: "default",
									Labels:    map[string]string{"team": "infrastructure"},
								},
								Data: map[string]string{
									"jwt.rego": `package test
	
				default allow = false
				allow {
					[header, payload, signature] = io.jwt.decode(input.state.jwt)
					payload["foo"] = "not-bar"
				}
				`}}
							modules := []*core.ResourceRef{{Name: policy.Metadata.Name}}
							options := &extauth.OpaAuthOptions{FastInputConversion: true}

							_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
								Metadata: &core.Metadata{
									Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
									Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
								},
								Configs: []*extauth.AuthConfig_Config{
									{
										AuthConfig: &extauth.AuthConfig_Config_Oauth{
											Oauth: getOauthConfig(envoyPort, secret.Metadata.Ref()),
										},
									},
									{
										AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
											OpaAuth: getOpaConfig(modules, options),
										},
									},
								},
							}, clients.WriteOpts{Ctx: ctx})
							Expect(err).NotTo(HaveOccurred())

							proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

							_, err = testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
							Expect(err).ToNot(HaveOccurred())
						})

						It("should NOT allow access", func() {
							EventuallyWithOffset(1, func() (int, error) {
								req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
								req.Header.Add("Authorization", "Bearer "+token)

								resp, err := http.DefaultClient.Do(req)
								if err != nil {
									return 0, err
								}
								return resp.StatusCode, nil
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))

						})

					})
				})

				ExpectUpstreamRequest := func() {
					EventuallyWithOffset(1, func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("Authorization", "Bearer "+token)

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))

					select {
					case r := <-testUpstream.C:
						ExpectWithOffset(1, r.Headers["X-User-Id"]).To(HaveLen(1))
						ExpectWithOffset(1, r.Headers["X-User-Id"][0]).To(Equal("http://localhost:5556;user"))
					case <-time.After(time.Second):
						Fail("expected a message to be received")
					}
				}

				Context("Oidc tests that do forward to upstream", func() {
					It("should allow access with proper jwt token", func() {
						ExpectUpstreamRequest()
					})
				})

				Context("oidc + opa sanity", func() {
					BeforeEach(func() {
						policy := &gloov1.Artifact{
							Metadata: &core.Metadata{
								Name:      "jwt",
								Namespace: "default",
								Labels:    map[string]string{"team": "infrastructure"},
							},
							Data: map[string]string{
								"jwt.rego": `package test

			default allow = false
			allow {
				[header, payload, signature] = io.jwt.decode(input.state.jwt)
				payload["foo"] = "bar"
			}
			`}}
						modules := []*core.ResourceRef{{Name: policy.Metadata.Name, Namespace: policy.Metadata.Namespace}}
						options := &extauth.OpaAuthOptions{FastInputConversion: true}

						_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
							Metadata: &core.Metadata{
								Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
								Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{
								{
									AuthConfig: &extauth.AuthConfig_Config_Oauth{
										Oauth: getOauthConfig(envoyPort, secret.Metadata.Ref()),
									},
								},
								{
									AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
										OpaAuth: getOpaConfig(modules, options),
									},
								},
							},
						}, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())
						proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

						_, err = testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
						Expect(err).ToNot(HaveOccurred())
					})
					It("should allow access", func() {
						ExpectUpstreamRequest()
					})
				})

			})

		})

		Context("health checker", func() {

			var healthCheckClient grpc_health_v1.HealthClient

			getHealthCheckClient := func() grpc_health_v1.HealthClient {
				if healthCheckClient != nil {
					return healthCheckClient
				}

				extAuthHealthServerAddr := "localhost:" + strconv.Itoa(settings.ExtAuthSettings.ServerPort)
				conn, err := grpc.Dial(extAuthHealthServerAddr, grpc.WithInsecure())
				Expect(err).ToNot(HaveOccurred())

				healthCheckClient = grpc_health_v1.NewHealthClient(conn)

				go func() {
					select {
					case <-ctx.Done():
						healthCheckClient = nil
						conn.Close()

						return
					}
				}()

				return healthCheckClient
			}

			getServiceHealthStatus := func() (grpc_health_v1.HealthCheckResponse_ServingStatus, error) {
				client := getHealthCheckClient()

				var header metadata.MD
				resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{
					Service: settings.ExtAuthSettings.ServiceName,
				}, grpc.Header(&header))

				return resp.GetStatus(), err
			}

			Context("should pass after receiving xDS config from gloo", func() {

				It("without auth configs", func() {
					Eventually(getServiceHealthStatus, "10s", ".1s").Should(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
					Consistently(getServiceHealthStatus, "3s", ".1s").Should(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
				})

				It("with auth configs", func() {
					// Creates a proxy with an auth configuration
					basicConfigSetup()

					Eventually(getServiceHealthStatus, "10s", ".1s").Should(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
					Consistently(getServiceHealthStatus, "3s", ".1s").Should(Equal(grpc_health_v1.HealthCheckResponse_SERVING))
				})

			})

			// NOTE: This test MUST run last, since it runs cancel()
			Context("shutdown", func() {

				It("should fail healthcheck immediately on shutdown", func() {

					Eventually(getServiceHealthStatus, "10s", ".1s").Should(Equal(grpc_health_v1.HealthCheckResponse_SERVING))

					// Start sending health checking requests continuously
					waitForHealthcheck := make(chan struct{})
					go func(waitForHealthcheck chan struct{}) {
						defer GinkgoRecover()
						Eventually(func() bool {
							ctx = context.Background()
							var header metadata.MD
							getHealthCheckClient().Check(ctx, &grpc_health_v1.HealthCheckRequest{
								Service: settings.ExtAuthSettings.ServiceName,
							}, grpc.Header(&header))
							return len(header.Get("x-envoy-immediate-health-check-fail")) == 1
						}, "5s", ".1s").Should(BeTrue())
						waitForHealthcheck <- struct{}{}
					}(waitForHealthcheck)

					// Start the health checker first, then cancel
					time.Sleep(200 * time.Millisecond)
					cancel()
					Eventually(waitForHealthcheck, "5s", ".1s").Should(Receive())
				})
			})

		})

	})

})

var startDiscoveryServerOnce sync.Once
var cachedPrivateKey *rsa.PrivateKey

type fakeDiscoveryServer struct {
	s                        http.Server
	createExpiredToken       bool
	createNearlyExpiredToken bool
	token                    string
	lastGrant                string
	handlerStats             map[string]int

	// The set of key IDs that are supported by the server
	keyIds []string
}

func (f *fakeDiscoveryServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = f.s.Shutdown(ctx)
}

func (f *fakeDiscoveryServer) updateToken(grantType string) {
	startDiscoveryServerOnce.Do(func() {
		var err error
		cachedPrivateKey, err = rsa.GenerateKey(rand.Reader, 512)
		Expect(err).NotTo(HaveOccurred())
	})

	claims := jwt.MapClaims{
		"foo": "bar",
		"aud": "test-clientid",
		"sub": "user",
		"iss": "http://localhost:5556",
	}
	f.lastGrant = grantType
	if grantType == "" && f.createExpiredToken {
		// create expired token so we can test refresh
		claims["exp"] = time.Now().Add(-time.Minute).Unix()
	} else if grantType == "" && f.createNearlyExpiredToken {
		// create token that expires ten ms from now
		claims["exp"] = time.Now().Add(10 * time.Millisecond).Unix()
	}

	tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenToSign.Header["kid"] = f.getValidKeyId()
	token, err := tokenToSign.SignedString(cachedPrivateKey)
	Expect(err).NotTo(HaveOccurred())
	f.token = token
}

func (f *fakeDiscoveryServer) updateKeyIds(keyIds []string) {
	if len(keyIds) > 0 {
		f.keyIds = keyIds
	}
}

func (f *fakeDiscoveryServer) getValidKeyId() string {
	// If there is more than one valid kid, return the first one
	return f.keyIds[0]
}

func (f *fakeDiscoveryServer) Start() *rsa.PrivateKey {
	// Initialize the server with 1 valid kid
	f.keyIds = []string{"kid-1"}

	f.updateToken("")
	n := base64.RawURLEncoding.EncodeToString(cachedPrivateKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(0).SetUint64(uint64(cachedPrivateKey.E)).Bytes())

	f.s = http.Server{
		Addr: ":5556",
	}

	f.s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer GinkgoRecover()
		rw.Header().Set("content-type", "application/json")

		if f.handlerStats != nil {
			f.handlerStats[r.URL.Path] += 1
		}
		switch r.URL.Path {
		case "/auth":
			// redirect back immediately. This simulates a user that's already logged in by the IDP.
			redirect_uri := r.URL.Query().Get("redirect_uri")
			state := r.URL.Query().Get("state")
			u, err := url.Parse(redirect_uri)
			Expect(err).NotTo(HaveOccurred())

			u.RawQuery = "code=1234&state=" + state
			fmt.Fprintf(GinkgoWriter, "redirecting to %s\n", u.String())
			rw.Header().Add("Location", u.String())
			rw.Header().Add("x-auth", "auth")
			rw.WriteHeader(http.StatusFound)

			_, _ = rw.Write([]byte(`auth`))
		case "/alternate-auth":
			// redirect back immediately. This simulates a user that's already logged in by the IDP.
			redirect_uri := r.URL.Query().Get("redirect_uri")
			state := r.URL.Query().Get("state")
			u, err := url.Parse(redirect_uri)
			Expect(err).NotTo(HaveOccurred())

			u.RawQuery = "code=9876&state=" + state
			fmt.Fprintf(GinkgoWriter, "redirecting to %s\n", u.String())
			rw.Header().Add("Location", u.String())
			rw.Header().Add("x-auth", "alternate-auth")
			rw.WriteHeader(http.StatusFound)

			_, _ = rw.Write([]byte(`alternate-auth`))
		case "/.well-known/openid-configuration":
			_, _ = rw.Write([]byte(`
		{
			"issuer": "http://localhost:5556",
			"authorization_endpoint": "http://localhost:5556/auth",
			"token_endpoint": "http://localhost:5556/token",
			"jwks_uri": "http://localhost:5556/keys",
			"response_types_supported": [
			  "code"
			],
			"subject_types_supported": [
			  "public"
			],
			"id_token_signing_alg_values_supported": [
			  "RS256"
			],
			"scopes_supported": [
			  "openid",
			  "email",
			  "profile"
			]
		  }
		`))
		case "/token":
			r.ParseForm()
			fmt.Fprintln(GinkgoWriter, "got request for token. query:", r.URL.RawQuery, r.URL.String(), "form:", r.Form.Encode())
			if r.URL.Query().Get("grant_type") == "refresh_token" || r.Form.Get("grant_type") == "refresh_token" {
				f.updateToken("refresh_token")
			}
			_, _ = rw.Write([]byte(`
			{
				"access_token": "SlAV32hkKG",
				"token_type": "Bearer",
				"refresh_token": "8xLOxBtZp8",
				"expires_in": 3600,
				"id_token": "` + f.token + `"
			 }
	`))
		case "/keys":
			var keyListBuffer bytes.Buffer
			for _, kid := range f.keyIds {
				keyListBuffer.WriteString(`
				{
					"use": "sig",
					"kty": "RSA",
					"kid": "` + kid + `",
					"alg": "RS256",
					"n": "` + n + `",
					"e": "` + e + `"
				},`)
			}
			// Remove the last comma so it's valid json
			keyList := strings.TrimSuffix(keyListBuffer.String(), ",")
			keysResponse := `
			{
				"keys": [
				    ` + keyList + `
				]
			}
			`
			_, _ = rw.Write([]byte(keysResponse))
		}
	})

	go func() {
		defer GinkgoRecover()
		err := f.s.ListenAndServe()
		if err != http.ErrServerClosed {
			Expect(err).NotTo(HaveOccurred())
		}
	}()

	return cachedPrivateKey
}

func getOauthTokenIntrospectionUrlConfig() *extauth.OAuth2_AccessTokenValidation {
	return &extauth.OAuth2_AccessTokenValidation{
		AccessTokenValidation: &extauth.AccessTokenValidation{
			ValidationType: &extauth.AccessTokenValidation_IntrospectionUrl{
				IntrospectionUrl: "http://localhost:5556/introspection",
			},
			UserinfoUrl:  "http://localhost:5556/userinfo",
			CacheTimeout: nil,
		},
	}
}

func getOauthTokenIntrospectionConfig(clientId string, clientSecretRef *core.ResourceRef) *extauth.OAuth2_AccessTokenValidation {
	return &extauth.OAuth2_AccessTokenValidation{
		AccessTokenValidation: &extauth.AccessTokenValidation{
			ValidationType: &extauth.AccessTokenValidation_Introspection{
				Introspection: &extauth.AccessTokenValidation_IntrospectionValidation{
					IntrospectionUrl: "http://localhost:5556/introspection",
					ClientId:         clientId,
					ClientSecretRef:  clientSecretRef,
				},
			},
			UserinfoUrl:  "http://localhost:5556/userinfo",
			CacheTimeout: nil,
		},
	}
}

func getOauthTokenIntrospectionExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oauth-token-introspection",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getProxyExtAuthOauthTokenIntrospection(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getOauthTokenIntrospectionExtAuthExtension(), false)
}

func getOauthConfig(envoyPort uint32, secretRef *core.ResourceRef) *extauth.OAuth {
	return &extauth.OAuth{
		ClientId:        "test-clientid",
		ClientSecretRef: secretRef,
		IssuerUrl:       "http://localhost:5556/",
		AppUrl:          fmt.Sprintf("http://localhost:%d", envoyPort),
		CallbackPath:    "/callback",
		Scopes:          []string{"email"},
	}
}

func getOidcAuthCodeConfig(envoyPort uint32, secretRef *core.ResourceRef) *extauth.OAuth2_OidcAuthorizationCode {
	return &extauth.OAuth2_OidcAuthorizationCode{
		OidcAuthorizationCode: &extauth.OidcAuthorizationCode{
			ClientId:        "test-clientid",
			ClientSecretRef: secretRef,
			IssuerUrl:       "http://localhost:5556/",
			AppUrl:          fmt.Sprintf("http://localhost:%d", envoyPort),
			CallbackPath:    "/callback",
			LogoutPath:      "/logout",
			Scopes:          []string{"email"},
		},
	}
}

func getProxyExtAuthOIDC(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getOidcExtAuthExtension(), false)
}

func getOidcExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oidc-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getProxyExtAuthOIDCAndOpa(envoyPort uint32, secretRef, upstream *core.ResourceRef, modules []*core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getOidcAndOpaExtAuthExtension(), false)
}

func getOidcAndOpaExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oidcand-opa-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getOpaConfig(modules []*core.ResourceRef, options *extauth.OpaAuthOptions) *extauth.OpaAuth {
	return &extauth.OpaAuth{
		Modules: modules,
		Query:   "data.test.allow == true",
		Options: options,
	}
}

func getProxyExtAuthBasicAuth(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, GetBasicAuthExtension(), false)
}

func GetBasicAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "basic-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getBasicAuthConfig() *extauth.BasicAuth {
	return &extauth.BasicAuth{
		Realm: "gloo",
		Apr: &extauth.BasicAuth_Apr{
			Users: map[string]*extauth.BasicAuth_Apr_SaltedHashedPassword{
				"user": {
					// Password is password
					Salt:           "0adzfifo",
					HashedPassword: "14o4fMw/Pm2L34SvyyA2r.",
				},
			},
		},
	}
}

func getProxyExtAuthApiKeyAuth(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getApiKeyExtAuthExtension(), false)
}

func getApiKeyAuthConfig() *extauth.ApiKeyAuth {
	return &extauth.ApiKeyAuth{
		ApiKeySecretRefs: []*core.ResourceRef{
			{
				Namespace: "default",
				Name:      "secret1",
			},
		},
		LabelSelector: map[string]string{"team": "infrastructure"},
	}
}

func getApiKeyExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "apikey-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getProxyExtAuthPassThroughAuth(envoyPort uint32, upstream *core.ResourceRef, zipkinTracing bool) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, GetPassThroughExtAuthExtension(), zipkinTracing)
}

func GetPassThroughExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "passthrough-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

// This provides PassThroughAuth AuthConfig with Custom Config
func getPassThroughAuthWithCustomConfig(address string) *extauth.PassThroughAuth {
	passThroughAuth := getPassThroughAuthConfig(address)
	passThroughAuth.Config = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"customConfig1": {
				Kind: &structpb.Value_BoolValue{
					BoolValue: true,
				},
			},
		},
	}
	return passThroughAuth
}

func getPassThroughAuthConfig(address string) *extauth.PassThroughAuth {
	return &extauth.PassThroughAuth{
		Protocol: &extauth.PassThroughAuth_Grpc{
			Grpc: &extauth.PassThroughGrpc{
				Address: address,
				// use default connection timeout
			},
		},
	}
}

func getProxyExtAuth(envoyPort uint32, upstream *core.ResourceRef, extauthCfg *extauth.ExtAuthExtension, zipkinTracing bool) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "gloo-system.virt1",
		Domains: []string{"*"},
		Options: &gloov1.VirtualHostOptions{
			Extauth: extauthCfg,
		},
		Routes: []*gloov1.Route{{
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
		}},
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

	if zipkinTracing {
		p.Listeners[0] = &gloov1.Listener{
			Name:        "listener",
			BindAddress: "0.0.0.0",
			BindPort:    envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
					Options: &gloov1.HttpListenerOptions{
						HttpConnectionManagerSettings: &hcm.HttpConnectionManagerSettings{
							Tracing: &tracing.ListenerTracingSettings{
								ProviderConfig: &tracing.ListenerTracingSettings_ZipkinConfig{
									ZipkinConfig: &v3.ZipkinConfig{
										CollectorCluster: &v3.ZipkinConfig_CollectorUpstreamRef{
											CollectorUpstreamRef: &core.ResourceRef{
												Namespace: "default",
												Name:      "zipkin",
											},
										},
										CollectorEndpoint:        "/api/v2/spans",
										CollectorEndpointVersion: v3.ZipkinConfig_HTTP_JSON,
									},
								},
							},
						},
					},
				},
			},
		}
	}

	return p
}

type unsecureCookieJar struct {
	http.CookieJar
}

func (j *unsecureCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, c := range cookies {
		// hack to work around go client impl that doesn't consider localhost secure.
		c.Secure = false
	}
	j.CookieJar.SetCookies(u, cookies)
}
