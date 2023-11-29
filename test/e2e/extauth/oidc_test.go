package extauth_test

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/onsi/gomega/gstruct"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/solo-io/solo-projects/test/services/redis"

	"github.com/solo-io/solo-projects/test/gomega/matchers"
	"github.com/solo-io/solo-projects/test/gomega/transforms"

	. "github.com/solo-io/solo-projects/test/services/extauth/servers"

	"github.com/golang-jwt/jwt"
	"github.com/golang/protobuf/ptypes"
	"github.com/solo-io/ext-auth-service/pkg/utils/cipher"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("OIDC", func() {

	var (
		testContext           *e2e.TestContextWithExtensions
		authConfig            *extauth.AuthConfig
		oauth2                *extauth.OAuth2_OidcAuthorizationCode
		privateKey            *rsa.PrivateKey
		discoveryServer       FakeDiscoveryServer
		secret                *gloov1.Secret
		token                 string
		cookies               []*http.Cookie
		virtualServiceBuilder *gloohelpers.VirtualServiceBuilder
		jwtRegex              = regexp.MustCompile("(\\w*.\\w*.\\w*)")
	)

	const (
		accessTokenValue  = "SlAV32hkKG"
		refreshTokenValue = "8xLOxBtZp8"
		// The id_token is a jwt so the only thing that would match is the first portion of the jwt(header)
		idTokenSubstring = "eyJhbGciOiJSUzI1NiIsImtpZCI6ImtpZC0xIiwidHlwIjoiSldUIn0"
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		cookies = nil
		discoveryServer = FakeDiscoveryServer{
			AccessTokenValue:  accessTokenValue,
			RefreshTokenValue: refreshTokenValue,
		}
		privateKey = discoveryServer.Start("localhost")

		secret = &gloov1.Secret{
			Metadata: &core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &gloov1.Secret_Oauth{
				Oauth: &extauth.OauthSecret{
					ClientSecret: "test",
				},
			},
		}

		oauth2 = discoveryServer.GetOidcAuthCodeConfig(testContext.EnvoyInstance().HttpPort, "localhost", secret.Metadata.Ref())

		testContext.ResourcesToCreate().Secrets = append(testContext.ResourcesToCreate().Secrets, secret)

		// get id token
		token = discoveryServer.Token
	})

	JustBeforeEach(func() {
		vs := virtualServiceBuilder.Build()
		// add the app url to the virtual service domains to allow the oauth redirects to work
		appUrlHttp := fmt.Sprintf("%s:%d", "localhost", testContext.EnvoyInstance().HttpPort)
		vs.GetVirtualHost().Domains = append(vs.GetVirtualHost().Domains, appUrlHttp)
		testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
			vs,
		}

		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		discoveryServer.Stop()
		testContext.AfterEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	makeSingleRequest := func(client *http.Client) (*http.Response, error) {
		// the default host causes issues, likely with redirects as we get redirected to "test.com/" and we are unable to reach the server.
		// we remove the host to avoid this issue.
		httpReqBuilder := testContext.GetHttpRequestBuilder().WithPath("success?foo=bar").WithHost("")
		resp, err := client.Do(httpReqBuilder.Build())
		if err != nil {
			return nil, err
		}
		return resp, nil
	}

	ExpectHappyPathToWork := func(makeSingleRequest func(client *http.Client) (*http.Response, error), loginSuccessExpectation func()) {
		// do auth flow and make sure we have a cookie named cookie:
		appPageReq := testContext.GetHttpRequestBuilder().Build()

		var finalUrl *url.URL
		jar, err := cookiejar.New(nil)
		Expect(err).NotTo(HaveOccurred())
		cookieJar := &unsecureCookieJar{CookieJar: jar}
		client := testutils.DefaultClientBuilder().Build()
		client.Jar = cookieJar
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			finalUrl = req.URL
			if len(via) > 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		}

		Eventually(func() (*http.Response, error) {
			return makeSingleRequest(client)
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

		Expect(finalUrl).NotTo(BeNil())
		Expect(finalUrl.Path).To(Equal("/success"))
		// make sure query is passed through as well
		Expect(finalUrl.RawQuery).To(Equal("foo=bar"))

		// check the cookie jar
		tmpCookies := jar.Cookies(appPageReq.URL)
		Expect(tmpCookies).NotTo(BeEmpty())

		// grab the original cookies for these cookies, as `jar.Cookies` doesn't return
		// all the properties of the cookies
		for _, c := range tmpCookies {
			cookies = append(cookies, cookieJar.OriginalCookies[c.Name])
		}

		// make sure login is successful
		loginSuccessExpectation()

		// try to logout:
		logoutReqBuilder := testContext.GetHttpRequestBuilder().WithPath("logout")
		resp, err := client.Do(logoutReqBuilder.Build())
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(HaveHTTPStatus(http.StatusOK))

		// Verify that the logout resulted in a redirect to the default url
		Expect(finalUrl).NotTo(BeNil())
		Expect(finalUrl.Path).To(Equal("/"))
	}

	ExpectHappyPathToWorkWithGomega := func(makeSingleRequest func(client *http.Client) (*http.Response, error), loginSuccessExpectation func(g Gomega), g Gomega) {
		// do auth flow and make sure we have a cookie named cookie:
		appPageReq := testContext.GetHttpRequestBuilder().Build()

		var finalUrl *url.URL
		jar, err := cookiejar.New(nil)
		g.Expect(err).NotTo(HaveOccurred())
		cookieJar := &unsecureCookieJar{CookieJar: jar}
		client := testutils.DefaultClientBuilder().Build()
		client.Jar = cookieJar
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			finalUrl = req.URL
			if len(via) > 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		}

		g.Eventually(func() (*http.Response, error) {
			return makeSingleRequest(client)
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

		g.Expect(finalUrl).NotTo(BeNil())
		g.Expect(finalUrl.Path).To(Equal("/success"))
		// make sure query is passed through as well
		Expect(finalUrl.RawQuery).To(Equal("foo=bar"))

		// check the cookie jar
		tmpCookies := jar.Cookies(appPageReq.URL)
		g.Expect(tmpCookies).NotTo(BeEmpty())

		// grab the original cookies for these cookies, as jar.Cookies doesn't return
		// all the properties of the cookies
		for _, c := range tmpCookies {
			cookies = append(cookies, cookieJar.OriginalCookies[c.Name])
		}

		// make sure login is successful
		loginSuccessExpectation(g)

		// try to logout:
		logoutReqBuilder := testContext.GetHttpRequestBuilder().WithPath("logout")
		resp, err := client.Do(logoutReqBuilder.Build())
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(resp).To(HaveHTTPStatus(http.StatusOK))
		// Verify that the logout resulted in a redirect to the default url
		g.Expect(finalUrl).NotTo(BeNil())
		g.Expect(finalUrl.Path).To(Equal("/"))
	}

	ExpectUpstreamRequest := func() {
		httpReqBuilder := testContext.GetHttpRequestBuilder().WithHeader("Authorization", "Bearer "+token).WithHost("")
		EventuallyWithOffset(1, func() (*http.Response, error) {
			resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
			if err != nil {
				return nil, err
			}
			return resp, nil
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

		select {
		case r := <-testContext.TestUpstream().C:
			ExpectWithOffset(1, r.Headers).To(WithTransform(HeaderStripper(),
				HaveKeyWithValue("X-User-Id", fmt.Sprintf("http://%s:%d;user", discoveryServer.ServerAddress, discoveryServer.Port)),
			))
		case <-time.After(time.Second):
			Fail("expected a message to be received")
		}
	}

	Context("oidc (no opa)", func() {
		BeforeEach(func() {
			// Set the extauth extension for the default VS to be oidc
			virtualServiceBuilder = gloohelpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
			virtualServiceBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
				Extauth: getOidcExtAuthExtension(),
			})
		})

		// The AuthConfig should be written after the tests in this context modify the oauth2 config
		// We can't use the ResourcesToCreate() method because that is called in testContext.JustBeforeEach(), which is before this.
		JustBeforeEach(func() {
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
			_, err := testContext.TestClients().AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: testContext.Ctx()})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("redis for session store", func() {
			const (
				cookieName = "cookie"
			)
			var (
				redisInstance *redis.Instance
			)

			BeforeEach(func() {
				redisInstance = redisFactory.NewInstance()
				redisInstance.Run(testContext.Ctx())

				// update the config to use redis
				oauth2.OidcAuthorizationCode.Session = getRedisUserSession(redisInstance.Url(), cookieName)
			})

			AfterEach(func() {
				redisInstance.Clean()
			})

			It("should work", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {
					Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
						HaveLen(1),
						HaveKey(cookieName),
					)))
				})
			})

			It("should refresh token", func() {
				discoveryServer.CreateExpiredToken = true
				discoveryServer.UpdateToken("")
				ExpectHappyPathToWork(makeSingleRequest, func() {
					Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
						HaveLen(1),
						HaveKey(cookieName),
					)))
				})
				Expect(discoveryServer.LastGrant).To(Equal("refresh_token"))
			})

			It("should auth successfully after refreshing token", func() {
				forceTokenRefresh := func(client *http.Client) (*http.Response, error) {
					// Create token that will expire in 1 second
					discoveryServer.CreateNearlyExpiredToken = true
					discoveryServer.UpdateToken("")
					discoveryServer.CreateNearlyExpiredToken = false
					Expect(discoveryServer.HandlerStats.Get(TokenEndpoint)).To(BeEquivalentTo(0))

					// execute first request.
					Eventually(func() (*http.Response, error) {
						return makeSingleRequest(client)
					}, "10s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

					// sleep for 1 second, so the token expires
					time.Sleep(time.Second)

					// execute second request.
					resp, err := makeSingleRequest(client)
					Expect(err).NotTo(HaveOccurred())

					// execute third request. We should not hit the /token handler, because the refreshed token should be in the store.
					baseRefreshes := discoveryServer.HandlerStats.Get(TokenEndpoint)
					resp, err = makeSingleRequest(client)
					Expect(err).NotTo(HaveOccurred())
					Expect(discoveryServer.HandlerStats.Get(TokenEndpoint)).To(BeNumerically("==", baseRefreshes))

					return resp, nil
				}

				ExpectHappyPathToWork(forceTokenRefresh, func() {
					Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
						HaveLen(1),
						HaveKey(cookieName),
					)))
				})
			})

			Context("no refreshing", func() {
				BeforeEach(func() {
					oauth2.OidcAuthorizationCode.Session.Session.(*extauth.UserSession_Redis).Redis.AllowRefreshing = &wrappers.BoolValue{Value: false}
				})

				It("should NOT refresh token", func() {
					discoveryServer.CreateExpiredToken = true
					discoveryServer.UpdateToken("")

					jar, err := cookiejar.New(nil)
					Expect(err).NotTo(HaveOccurred())
					client := testutils.DefaultClientBuilder().Build()
					client.Jar = &unsecureCookieJar{CookieJar: jar}

					// removing the default host for the http builder, because the redirect(s) cause it to try to
					// connect to "http://test.com/success?foo=bar" and causes the context to time out and fail.
					httpReqBuilder := testContext.GetHttpRequestBuilder().WithPath("success?foo=bar").WithHost("")

					// as we will always provide an expired token, this will result in a redirect loop.
					Eventually(func() error {
						resp, err := client.Do(httpReqBuilder.Build())
						if err != nil {
							return err
						}
						defer resp.Body.Close()
						_, _ = io.ReadAll(resp.Body)
						return nil
					}, "5s", "0.5s").Should(MatchError(ContainSubstring("stopped after 10 redirects")))
					Expect(discoveryServer.LastGrant).To(Equal(""))
				})
			})

			// add context with refresh; get an expired token going and make sure it was refreshed.
		})

		Context("forward header token with Bearer schema", func() {
			BeforeEach(func() {
				// update the config to use redis
				oauth2.OidcAuthorizationCode.Headers = &extauth.HeaderConfiguration{
					IdTokenHeader:                   "foo",
					AccessTokenHeader:               "Authorization",
					UseBearerSchemaForAuthorization: &wrappers.BoolValue{Value: true},
				}
			})

			It("should use Bearer schema if using Authorization access token header", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {})

				select {
				case r := <-testContext.TestUpstream().C:
					Expect(r.Headers).To(WithTransform(HeaderStripper(), And(
						HaveKeyWithValue("Foo", discoveryServer.Token),
						HaveKeyWithValue("Authorization", fmt.Sprintf("Bearer %s", accessTokenValue)),
					)))
				case <-time.After(time.Second):
					Fail("timed out")
				}
			})
		})

		Context("does NOT forward header token with Bearer schema if not enabled", func() {
			BeforeEach(func() {
				// update the config to use redis
				oauth2.OidcAuthorizationCode.Headers = &extauth.HeaderConfiguration{
					IdTokenHeader:                   "foo",
					AccessTokenHeader:               "Authorization",
					UseBearerSchemaForAuthorization: &wrappers.BoolValue{Value: false},
				}
			})

			It("should not use Bearer schema", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {})

				select {
				case r := <-testContext.TestUpstream().C:
					Expect(r.Headers).To(WithTransform(HeaderStripper(), And(
						HaveKeyWithValue("Foo", discoveryServer.Token),
						HaveKeyWithValue("Authorization", accessTokenValue),
					)))
				case <-time.After(time.Second):
					Fail("timed out")
				}
			})
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
				case r := <-testContext.TestUpstream().C:
					Expect(r.Headers).To(WithTransform(HeaderStripper(), And(
						HaveKeyWithValue("Foo", discoveryServer.Token),
						HaveKeyWithValue("Bar", accessTokenValue),
					)))
				case <-time.After(time.Second):
					Fail("timed out")
				}
			})
		})

		Context("forward id token normally even if bearer addition enabled", func() {
			BeforeEach(func() {
				// update the config to use redis
				oauth2.OidcAuthorizationCode.Headers = &extauth.HeaderConfiguration{
					IdTokenHeader:                   "foo",
					AccessTokenHeader:               "bar",
					UseBearerSchemaForAuthorization: &wrappers.BoolValue{Value: true},
				}
			})

			It("should work", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {})

				select {
				case r := <-testContext.TestUpstream().C:
					Expect(r.Headers).To(WithTransform(HeaderStripper(), And(
						HaveKeyWithValue("Foo", discoveryServer.Token),
						HaveKeyWithValue("Bar", accessTokenValue),
					)))
				case <-time.After(time.Second):
					Fail("timed out")
				}
			})
		})

		Context("discovery override", func() {

			BeforeEach(func() {
				oauth2.OidcAuthorizationCode.DiscoveryOverride = &extauth.DiscoveryOverride{
					AuthEndpoint: fmt.Sprintf("http://%s:%d%s", discoveryServer.ServerAddress, discoveryServer.Port, AlternateAuthPath),
				}
			})

			It("should redirect to different auth endpoint with auth override", func() {
				client := testutils.DefaultClientBuilder().Build()
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					// stop at the auth point
					if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
						return http.ErrUseLastResponse
					}
					return nil
				}
				// Confirm that the response matches the one set by the /alternate-auth endpoint
				httpReqBuilder := testContext.GetHttpRequestBuilder()
				Eventually(func() (*http.Response, error) {
					resp, err := client.Do(httpReqBuilder.Build())
					if err != nil {
						return nil, err
					}
					_, _ = fmt.Fprintf(GinkgoWriter, "headers are %v \n", resp.Header)
					return resp, nil
				}, "5s", "0.5s").Should(HaveHTTPBody("alternate-auth"))
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
				httpReqBuilder := testContext.GetHttpRequestBuilder().WithHeader("Authorization", "Bearer "+token)
				EventuallyWithOffset(offset+1, func() (*http.Response, error) {
					resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
					if err != nil {
						return nil, err
					}
					return resp, nil
				}, "10s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
			}

			// A request with invalid token will be redirected to the /auth endoint
			expectRequestWithTokenFails := func(offset int, token string) {
				client := testutils.DefaultClientBuilder().Build()
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					// stop at the auth point
					if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
						return http.ErrUseLastResponse
					}
					return nil
				}

				httpReqBuilder := testContext.GetHttpRequestBuilder().WithHeader("Authorization", "Bearer "+token)
				EventuallyWithOffset(offset+1, func() (*http.Response, error) {
					resp, err := client.Do(httpReqBuilder.Build())
					if err != nil {
						return nil, err
					}
					return resp, nil
				}, "5s", "0.5s").Should(HaveHTTPBody("auth"))
			}

			JustBeforeEach(func() {
				// Ensure that keys have been loaded properly
				validToken := discoveryServer.GenerateValidIdToken(privateKey)
				expectRequestWithTokenSucceeds(0, validToken)
			})

			When("policy is nil or NEVER", func() {
				BeforeEach(func() {
					oauth2.OidcAuthorizationCode.JwksCacheRefreshPolicy = nil
				})

				It("should accept token with valid kid", func() {
					validToken := discoveryServer.GenerateValidIdToken(privateKey)
					expectRequestWithTokenSucceeds(0, validToken)
				})

				It("should deny token with new kid", func() {
					invalidToken := discoveryServer.GenerateIdTokenWithKid("kid-2", privateKey)
					expectRequestWithTokenFails(0, invalidToken)
				})

				It("should deny token with new kid after keys rotate", func() {
					// rotate the keys
					discoveryServer.UpdateKeyIds([]string{"kid-2"})

					// execute a request with the valid token
					// it should be denied because the local cache is never updated
					newToken := discoveryServer.GenerateIdTokenWithKid("kid-2", privateKey)
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
					validToken := discoveryServer.GenerateValidIdToken(privateKey)
					expectRequestWithTokenSucceeds(0, validToken)
				})

				It("should deny token with new kid", func() {
					invalidToken := discoveryServer.GenerateIdTokenWithKid("kid-2", privateKey)
					expectRequestWithTokenFails(0, invalidToken)
				})

				It("should accept token with new kid after keys rotate", func() {
					for i := 0; i < 5; i++ {
						// rotate the keys
						newKid := fmt.Sprintf("kid-new-%d", i)
						discoveryServer.UpdateKeyIds([]string{newKid})

						// execute a request using the new token
						// it should be accepted because the local cache gets updated
						validToken := discoveryServer.GenerateIdTokenWithKid(newKid, privateKey)
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
						discoveryServer.UpdateKeyIds([]string{newKid})

						// execute a request using the new token
						// it should be accepted because the local cache gets updated
						validToken := discoveryServer.GenerateIdTokenWithKid(newKid, privateKey)
						expectRequestWithTokenSucceeds(1, validToken)
					}

					// rotate the keys one more time
					discoveryServer.UpdateKeyIds([]string{rateLimitedKid})

					// execute a request using the new token
					// it should be rejected because the local cache no longer will be updated
					newToken := discoveryServer.GenerateIdTokenWithKid(rateLimitedKid, privateKey)
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
					validToken := discoveryServer.GenerateValidIdToken(privateKey)
					expectRequestWithTokenSucceeds(0, validToken)
				})

				It("should deny token with new kid", func() {
					invalidToken := discoveryServer.GenerateIdTokenWithKid("kid-2", privateKey)
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
						newToken := discoveryServer.GenerateIdTokenWithKid(rateLimitedKid, privateKey)
						expectRequestWithTokenSucceeds(0, newToken)

						expectRequestWithNewKidAcceptedNTimes()
					})
				})

			})
		})

		Context("happy path with default settings (no redis)", func() {
			It("should work", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {
					Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
						HaveKeyWithValue("id_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: true,
							Value:    ContainSubstring(idTokenSubstring),
						})),
						HaveKeyWithValue("access_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: true,
							Value:    accessTokenValue,
						})),
					)))
				})
			})
		})

		Context("happy path with default settings and http only set to false", func() {
			BeforeEach(func() {
				oauth2.OidcAuthorizationCode.Session = &extauth.UserSession{
					Session: &extauth.UserSession_Cookie{Cookie: &extauth.UserSession_InternalSession{}},
					CookieOptions: &extauth.UserSession_CookieOptions{
						HttpOnly: &wrappers.BoolValue{Value: false},
					},
				}
			})

			It("should work", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {
					Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
						HaveKeyWithValue("id_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: false,
							Value:    ContainSubstring(idTokenSubstring),
						})),
						HaveKeyWithValue("access_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: false,
							Value:    accessTokenValue,
						})),
					)))
				})
			})
		})

		Context("happy path with default settings and allowing refreshing", func() {
			BeforeEach(func() {
				oauth2.OidcAuthorizationCode.Session = &extauth.UserSession{
					Session: &extauth.UserSession_Cookie{Cookie: &extauth.UserSession_InternalSession{
						AllowRefreshing: &wrappers.BoolValue{Value: true},
					}},
				}
			})

			It("should work", func() {
				ExpectHappyPathToWork(makeSingleRequest, func() {
					Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
						HaveKeyWithValue("id_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: true,
							Value:    ContainSubstring(idTokenSubstring),
						})),
						HaveKeyWithValue("access_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: true,
							Value:    accessTokenValue,
						})),
						HaveKeyWithValue("refresh_token", matchers.MatchCookieData(&transforms.CookieData{
							HttpOnly: true,
							Value:    refreshTokenValue,
						})),
					)))
				})
			})
		})

		Context("happy path encryption session cookie", func() {
			var cookieValues = map[string]string{
				"id_token":      idTokenSubstring,
				"access_token":  accessTokenValue,
				"refresh_token": refreshTokenValue,
			}
			const (
				encryptionKey  = "this is a example encryption key"
				encryptionKey2 = "this is the second encryption ke"
			)

			validateEncryptedKeys := func(encryptionKey string) {
				ExpectHappyPathToWork(makeSingleRequest, func() {
					cookieMap := make(map[string]string)
					for _, c := range cookies {
						Expect(c.HttpOnly).To(BeFalse())
						cookieMap[c.Name] = c.Value
					}

					// The cookies are encrypted, so we can't assert the values equal to their original values
					Expect(cookieMap).To(And(
						HaveKey("id_token"),
						HaveKey("access_token"),
						HaveKey("refresh_token"),
					))

					cipher, err := cipher.NewGCMEncryption([]byte(encryptionKey))
					Expect(err).ToNot(HaveOccurred())

					for name, val := range cookieMap {
						expectedValue := cookieValues[name]
						unencryptedValue, err := cipher.Decrypt(val)
						Expect(err).ToNot(HaveOccurred())
						if name == "id_token" {
							// id_token is a jwt so the only thing that should match is the first portion of the jwt(header)
							// since the header will not change and is just a base64 of the json used
							// because of this we also check that the value is has a jwt format
							splits := strings.Split(unencryptedValue, ".")
							Expect(splits[0]).To(Equal(expectedValue))
							Expect(jwtRegex.Match([]byte(unencryptedValue))).To(BeTrue())
						} else {
							Expect(unencryptedValue).To(Equal(expectedValue))
						}
					}
				})
			}

			validateEncryptedKeysWithGomega := func(encryptionKey string, g Gomega) {
				ExpectHappyPathToWorkWithGomega(makeSingleRequest, func(g Gomega) {
					cookieMap := make(map[string]string)
					for _, c := range cookies {
						Expect(c.HttpOnly).To(BeFalse())
						cookieMap[c.Name] = c.Value
					}

					// The cookies are encrypted, so we can't assert the values equal to their original values
					Expect(cookieMap).To(And(
						HaveKey("id_token"),
						HaveKey("access_token"),
						HaveKey("refresh_token"),
					))

					c, err := cipher.NewGCMEncryption([]byte(encryptionKey))
					g.Expect(err).ToNot(HaveOccurred())

					for name, val := range cookieMap {
						expectedValue := cookieValues[name]
						unencryptedValue, err := c.Decrypt(val)
						g.Expect(err).ToNot(HaveOccurred())
						if name == "id_token" {
							// id_token is a jwt so the only thing that should match is the first portion of the jwt(header)
							// since the header will not change and is just a base64 of the json used
							// because of this we also check that the value is has a jwt format
							splits := strings.Split(unencryptedValue, ".")
							g.Expect(splits[0]).To(Equal(expectedValue))
							g.Expect(jwtRegex.Match([]byte(unencryptedValue))).To(BeTrue())
						} else {
							g.Expect(unencryptedValue).To(Equal(expectedValue))
						}
					}
				}, g)
			}

			Context("using Cipher Key Ref", func() {
				var encryptionSecret *gloov1.Secret

				BeforeEach(func() {
					encryptionSecret = &gloov1.Secret{
						Metadata: &core.Metadata{
							Name:      "encryption-secret",
							Namespace: "default",
						},
						Kind: &gloov1.Secret_Encryption{
							Encryption: &gloov1.EncryptionKeySecret{
								Key: encryptionKey,
							},
						},
					}

					oauth2.OidcAuthorizationCode.Session = &extauth.UserSession{
						Session: &extauth.UserSession_Cookie{Cookie: &extauth.UserSession_InternalSession{
							AllowRefreshing: &wrappers.BoolValue{Value: true},
						}},
						CipherConfig: &extauth.UserSession_CipherConfig{
							Key: &extauth.UserSession_CipherConfig_KeyRef{
								KeyRef: encryptionSecret.Metadata.Ref(),
							},
						},
						CookieOptions: &extauth.UserSession_CookieOptions{
							HttpOnly: &wrappers.BoolValue{Value: false},
						},
					}
				})

				JustBeforeEach(func() {
					createdSecret, err := testContext.TestClients().SecretClient.Write(encryptionSecret, clients.WriteOpts{Ctx: testContext.Ctx()})
					Expect(err).ToNot(HaveOccurred())
					encryptionSecret = createdSecret
				})

				It("should work with the secret, and rotating the secrets", func() {
					validateEncryptedKeys(encryptionKey)
					encryptionSecret.GetEncryption().Key = encryptionKey2
					newSecret, err := testContext.TestClients().SecretClient.Write(encryptionSecret, clients.WriteOpts{OverwriteExisting: true, Ctx: testContext.Ctx()})
					Expect(err).ToNot(HaveOccurred())
					encryptionSecret = newSecret
					Eventually(func(g Gomega) {
						cookies = nil
						validateEncryptedKeysWithGomega(encryptionKey2, g)
					}, "5s", "1s").Should(Succeed())
				})
			})

		})

		Context("Oidc callbackPath test", func() {
			BeforeEach(func() {
				oauth2.OidcAuthorizationCode.ParseCallbackPathAsRegex = true
				oauth2.OidcAuthorizationCode.CallbackPath = "/callback\\d"
			})

			It("should exchange token with regex callbackPath", func() {
				finalPageReq := testContext.GetHttpRequestBuilder().WithPath("success").Build()
				client := testutils.DefaultClientBuilder().Build()
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}

				st := oidc.NewStateSigner([]byte(testContext.ExtAuthInstance().GetServerSettings().SigningKey))
				signedState, err := st.Sign(finalPageReq.URL.String())
				Expect(err).NotTo(HaveOccurred())

				callbackReqBuilder := testContext.GetHttpRequestBuilder().WithPath(fmt.Sprintf("callback1?code=1234&state=%s", string(signedState)))

				Eventually(func(g Gomega) {
					resp, err := client.Do(callbackReqBuilder.Build())
					g.Expect(err).NotTo(HaveOccurred())

					expectedResponse := &testmatchers.HttpResponse{
						StatusCode: http.StatusFound,
						Headers: map[string]interface{}{
							"Location": finalPageReq.URL.String(),
						},
					}
					g.Expect(resp).To(testmatchers.HaveHttpResponse(expectedResponse))
				}, "5s", "0.5s").Should(Succeed())
			})
		})

		Context("Oidc tests that don't forward to upstream", func() {

			When("Authorization fails", func() {

				DescribeTable("failOnRedirect changes the response behavior",
					func(failOnRedirect bool, expectedStatusCode int) {
						err := gloohelpers.PatchResource(
							context.Background(),
							authConfig.GetMetadata().Ref(),
							func(resource resources.Resource) resources.Resource {
								ac := resource.(*extauth.AuthConfig)
								ac.FailOnRedirect = failOnRedirect

								// FailOnRedirect changes the behavior only when there are multiple Configs defined
								// so we add a second config that is a duplicate of the original
								ac.Configs = []*extauth.AuthConfig_Config{
									{
										Name: &wrappers.StringValue{
											Value: "config1",
										},
										AuthConfig: &extauth.AuthConfig_Config_Oauth2{
											Oauth2: &extauth.OAuth2{
												OauthType: oauth2,
											},
										},
									},
									{
										Name: &wrappers.StringValue{
											Value: "config2",
										},
										AuthConfig: &extauth.AuthConfig_Config_Oauth2{
											Oauth2: &extauth.OAuth2{
												OauthType: oauth2,
											},
										},
									},
								}
								ac.BooleanExpr = &wrappers.StringValue{
									Value: "(config1 || config2)",
								}
								return ac
							},
							testContext.TestClients().AuthConfigClient.BaseClient())
						Expect(err).NotTo(HaveOccurred())

						client := testutils.DefaultClientBuilder().Build()
						client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
							// stop at the auth point
							if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
								return http.ErrUseLastResponse
							}
							return nil
						}

						httpReqBuilder := testContext.GetHttpRequestBuilder()
						Eventually(func(g Gomega) {
							resp, err := client.Do(httpReqBuilder.Build())
							g.Expect(err).NotTo(HaveOccurred())
							g.Expect(resp).To(testmatchers.HaveHttpResponse(&testmatchers.HttpResponse{
								StatusCode: expectedStatusCode,
								Body:       gstruct.Ignore(),
							}))
						}, "10s", "0.5s").Should(Succeed())

					},
					Entry("failOnRedirect is false", false, http.StatusFound),

					// https://github.com/solo-io/ext-auth-service/issues/669
					// In reality, we should be getting a 401, but due to the above behavior in the ext-auth-service
					// a 403 is returned
					Entry("failOnRedirect is true", true, http.StatusForbidden),
				)

			})

			It("should include email scope in url", func() {
				client := testutils.DefaultClientBuilder().Build()
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}

				httpReqBuilder := testContext.GetHttpRequestBuilder()
				Eventually(func(g Gomega) {
					resp, err := client.Do(httpReqBuilder.Build())
					g.Expect(err).NotTo(HaveOccurred())

					expectedResponse := &testmatchers.HttpResponse{
						StatusCode: http.StatusFound,
						Headers: map[string]interface{}{
							"Location": ContainSubstring("email"),
						},
					}
					g.Expect(resp).To(testmatchers.HaveHttpResponse(expectedResponse))
				}, "5s", "0.5s").Should(Succeed())
			})

			It("should exchange token", func() {
				finalPageUrl := testContext.GetHttpRequestBuilder().WithPath("success").Build().URL.String()

				client := testutils.DefaultClientBuilder().Build()
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}

				st := oidc.NewStateSigner([]byte(testContext.ExtAuthInstance().GetServerSettings().SigningKey))
				signedState, err := st.Sign(finalPageUrl)
				Expect(err).NotTo(HaveOccurred())

				callbackReqBuilder := testContext.GetHttpRequestBuilder().WithPath(fmt.Sprintf("callback?code=1234&state=%s", string(signedState)))
				Eventually(func(g Gomega) {
					resp, err := client.Do(callbackReqBuilder.Build())
					g.Expect(err).NotTo(HaveOccurred())

					expectedResponse := &testmatchers.HttpResponse{
						StatusCode: http.StatusFound,
						Headers: map[string]interface{}{
							"Location": finalPageUrl,
						},
					}
					g.Expect(resp).To(testmatchers.HaveHttpResponse(expectedResponse))
				}, "5s", "0.5s").Should(Succeed())
			})

			Context("oidc + tls is already terminated", func() {
				Context("listener is http, appUrl is https", func() {
					BeforeEach(func() {
						oauth2.OidcAuthorizationCode.AppUrl = strings.Replace(oauth2.OidcAuthorizationCode.AppUrl,
							"http:",
							"https:",
							1)
						authConfig.Configs = []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_Oauth2{
								Oauth2: &extauth.OAuth2{
									OauthType: oauth2,
								},
							},
						}}
					})

					It("should prefer appUrl scheme to http request scheme", func() {
						client := testutils.DefaultClientBuilder().Build()
						client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
							return http.ErrUseLastResponse
						}

						httpReqBuilder := testContext.GetHttpRequestBuilder()
						// There is a delay between when we write configuration using our resource clients,
						// and when the ext-auth-service receives it over xDS from Gloo.
						// To handle this latency and prevent flakes, we wrap the test in an Eventually block
						Eventually(func(g Gomega) {
							resp, err := client.Do(httpReqBuilder.Build())
							g.Expect(err).NotTo(HaveOccurred())
							defer resp.Body.Close()
							locHdr := resp.Header.Get("Location")
							g.Expect(locHdr).NotTo(BeEmpty())
							locUrl, err := url.Parse(locHdr)
							g.Expect(err).NotTo(HaveOccurred())

							stateVals := locUrl.Query().Get("state")
							g.Expect(stateVals).NotTo(BeEmpty())

							var stateClaim struct {
								jwt.StandardClaims
								State string
							}
							_, err = jwt.ParseWithClaims(stateVals,
								&stateClaim,
								func(*jwt.Token) (interface{}, error) {
									return nil, nil
								},
							)

							// state URI has been upgraded to https:
							g.Expect(stateClaim.State).To(ContainSubstring("https://"))
						}, time.Second*3, time.Millisecond*250).Should(Succeed()) // originally ShouldNot(HaveOccurred())
					})
				})

				Context("listener is https, appUrl is http", func() {
					BeforeEach(func() {
						// overwrite the default gateway with a TLS gateway
						httpsGateway := gatewaydefaults.DefaultSslGateway(e2e.WriteNamespace)
						testContext.ResourcesToCreate().Gateways = gatewayv1.GatewayList{
							httpsGateway,
						}

						tlsSecret := &gloov1.Secret{
							Metadata: &core.Metadata{
								Name:      "tls-secret",
								Namespace: "default",
							},
							Kind: &gloov1.Secret_Tls{
								Tls: &gloov1.TlsSecret{
									CertChain:  gloohelpers.Certificate(),
									PrivateKey: gloohelpers.PrivateKey(),
								},
							},
						}
						testContext.ResourcesToCreate().Secrets = append(testContext.ResourcesToCreate().Secrets, tlsSecret)

						virtualServiceBuilder.WithSslConfig(&gloossl.SslConfig{
							SslSecrets: &gloossl.SslConfig_SecretRef{
								SecretRef: tlsSecret.Metadata.Ref(),
							},
						})
					})

					It("should prefer https scheme when appUrl scheme is a downgrade", func() {
						client := testutils.DefaultClientBuilder().WithTLSRootCa(gloohelpers.Certificate()).Build()
						client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
							return http.ErrUseLastResponse
						}

						httpsReqBuilder := testContext.GetHttpsRequestBuilder()

						// There is a delay between when we write configuration using our resource clients,
						// and when the ext-auth-service receives it over xDS from Gloo.
						// To handle this latency and prevent flakes, we wrap the test in an Eventually block
						Eventually(func(g Gomega) {
							resp, err := client.Do(httpsReqBuilder.Build())
							g.Expect(err).NotTo(HaveOccurred())

							locHdr := resp.Header.Get("Location")
							g.Expect(locHdr).NotTo(BeEmpty())
							locUrl, err := url.Parse(locHdr)
							g.Expect(err).NotTo(HaveOccurred())

							stateVals := locUrl.Query().Get("state")
							g.Expect(stateVals).NotTo(BeEmpty())

							var stateClaim struct {
								jwt.StandardClaims
								State string
							}
							_, err = jwt.ParseWithClaims(stateVals,
								&stateClaim,
								func(*jwt.Token) (interface{}, error) {
									return nil, nil
								},
							)

							// state URI has not been downgraded to http:
							g.Expect(stateClaim.State).To(ContainSubstring("https://"))
						}, time.Second*3, time.Millisecond*250).Should(Succeed())

					})

				})
			})
		})

		Context("Oidc tests that do forward to upstream", func() {
			It("should allow access with proper jwt token", func() {
				ExpectUpstreamRequest()
			})
		})

		Context("end session properties", func() {
			When("end session properties are set to POST", func() {
				BeforeEach(func() {
					oauth2.OidcAuthorizationCode.EndSessionProperties = &extauth.EndSessionProperties{
						MethodType: extauth.EndSessionProperties_PostMethod,
					}
				})

				It("does a POST request to the logout endpoint", func() {
					ExpectHappyPathToWork(makeSingleRequest, func() {})
					// the internal handler should have done a POST request to the /logout endpoint
					Expect(discoveryServer.HandlerStats.Get(LogoutEndpoint)).To(BeEquivalentTo(1))
					Expect(discoveryServer.EndpointData.GetMethodCount(LogoutEndpoint, http.MethodPost)).To(Equal(1))
				})
			})

			When("end session properties are set to GET (default)", func() {
				BeforeEach(func() {
					oauth2.OidcAuthorizationCode.EndSessionProperties = &extauth.EndSessionProperties{
						MethodType: extauth.EndSessionProperties_GetMethod,
					}
				})

				It("does a GET request to the logout endpoint", func() {
					ExpectHappyPathToWork(makeSingleRequest, func() {})
					// the internal handler should have done a GET request to the /logout endpoint
					Expect(discoveryServer.HandlerStats.Get(LogoutEndpoint)).To(BeEquivalentTo(1))
					Expect(discoveryServer.EndpointData.GetMethodCount(LogoutEndpoint, http.MethodGet)).To(Equal(1))
				})
			})
		})
	})

	Context("oidc + opa", func() {
		// OPA tests use the token generated by the oidc server for authorization, which includes a claim of [foo=bar]
		var (
			policy *gloov1.Artifact
		)

		BeforeEach(func() {
			virtualServiceBuilder = gloohelpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
			virtualServiceBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
				Extauth: getOidcAndOpaExtAuthExtension(),
			})
		})

		// The AuthConfig + Policy should be written after the tests in this context modify the oauth2 config
		// We can't use the ResourcesToCreate() method because that is called in testContext.JustBeforeEach(), which is before this.
		JustBeforeEach(func() {
			_, err := testContext.TestClients().ArtifactClient.Write(policy, clients.WriteOpts{Ctx: testContext.Ctx()})
			Expect(err).NotTo(HaveOccurred())
			_, err = testContext.TestClients().AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: testContext.Ctx()})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("with policy expecting jwt token to have [foo=bar] claim", func() {
			BeforeEach(func() {
				policy = &gloov1.Artifact{
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

				modules := []*core.ResourceRef{{
					Name:      policy.GetMetadata().GetName(),
					Namespace: policy.GetMetadata().GetNamespace(),
				}}
				options := &extauth.OpaAuthOptions{FastInputConversion: true}
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
						Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
					},
					Configs: []*extauth.AuthConfig_Config{
						{
							AuthConfig: &extauth.AuthConfig_Config_Oauth2{
								Oauth2: &extauth.OAuth2{
									OauthType: discoveryServer.GetOidcAuthCodeConfig(testContext.EnvoyInstance().HttpPort, "localhost", secret.Metadata.Ref()),
								},
							},
						},
						{
							AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
								OpaAuth: getOpaConfig(modules, options),
							},
						},
					},
				}
			})

			It("should allow access", func() {
				ExpectUpstreamRequest()
			})
		})

		Context("with policy expecting jwt token to have [foo=not-bar] claim", func() {
			BeforeEach(func() {
				policy = &gloov1.Artifact{
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

				modules := []*core.ResourceRef{{
					Name:      policy.GetMetadata().GetName(),
					Namespace: policy.GetMetadata().GetNamespace(),
				}}
				options := &extauth.OpaAuthOptions{FastInputConversion: true}
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().GetName(),
						Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().GetNamespace(),
					},
					Configs: []*extauth.AuthConfig_Config{
						{
							AuthConfig: &extauth.AuthConfig_Config_Oauth2{
								Oauth2: &extauth.OAuth2{
									OauthType: discoveryServer.GetOidcAuthCodeConfig(testContext.EnvoyInstance().HttpPort, "localhost", secret.GetMetadata().Ref()),
								},
							},
						},
						{
							AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
								OpaAuth: getOpaConfig(modules, options),
							},
						},
					},
				}
			})

			It("should NOT allow access", func() {
				httpReqBuilder := testContext.GetHttpRequestBuilder().WithHeader("Authorization", "Bearer "+token)
				Eventually(func(g Gomega) *http.Response {
					resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})

		})
	})
})

func getOidcExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oidc-auth",
				Namespace: e2e.WriteNamespace,
			},
		},
	}
}

func getOidcAndOpaExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oidcand-opa-auth",
				Namespace: e2e.WriteNamespace,
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

// Header stripper returns a Gomega Transform that maps an array of Headers (map[string][]string) to a map[string]string
// only keeping the first value of each header.
func HeaderStripper() func(c map[string][]string) map[string]string {
	return func(h map[string][]string) map[string]string {
		strippedHeaders := make(map[string]string)
		for header, values := range h {
			strippedHeaders[header] = values[0]
		}

		return strippedHeaders
	}
}

func getRedisUserSession(host, cookie string) *extauth.UserSession {
	return &extauth.UserSession{
		FailOnFetchFailure: true,
		Session: &extauth.UserSession_Redis{
			Redis: &extauth.UserSession_RedisSession{
				Options: &extauth.RedisOptions{
					Host: host,
				},
				KeyPrefix:       "key",
				CookieName:      cookie,
				AllowRefreshing: &wrappers.BoolValue{Value: true},
				PreExpiryBuffer: &duration.Duration{Seconds: 2, Nanos: 0},
			},
		},
	}
}
