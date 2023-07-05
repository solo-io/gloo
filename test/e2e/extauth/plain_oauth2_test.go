package extauth_test

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/solo-io/solo-projects/test/services/redis"

	matchers2 "github.com/solo-io/solo-projects/test/gomega/matchers"
	"github.com/solo-io/solo-projects/test/gomega/transforms"

	testServers "github.com/solo-io/solo-projects/test/services/extauth/servers"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	oauth2Service "github.com/solo-io/ext-auth-service/pkg/config/oauth2"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("Plain OAuth2", func() {

	var (
		testContext  *e2e.TestContextWithExtensions
		authConfig   *extauth.AuthConfig
		oauth2Server *testServers.FakeOAuth2Server
		oauth2       *extauth.OAuth2_Oauth2
		secret       *gloov1.Secret
		cookies      []*http.Cookie
	)

	const (
		accessTokenValue  = "SlAV32hkKG"
		refreshTokenValue = "8xLOxBtZp8"
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		cookies = nil
		oauth2Server = &testServers.FakeOAuth2Server{
			AccessTokenValue:  accessTokenValue,
			RefreshTokenValue: refreshTokenValue,
		}
		oauth2Server.Start()

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
		testContext.ResourcesToCreate().Secrets = append(testContext.ResourcesToCreate().Secrets, secret)

		// initialize the auth config
		oauth2 = oauth2Server.GetOAuth2Config(testContext.EnvoyInstance().HttpPort, "localhost", secret.Metadata.Ref())

		vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
		vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
			Extauth: getOAuth2ExtAuthExtension(),
		})
		vs := vsBuilder.Build()

		// Adding the app url as a domain. Without this, we 404 issues during the auth server redirects + callback.
		appUrlDomain := fmt.Sprintf("%s:%d", "localhost", testContext.EnvoyInstance().HttpPort)
		vs.GetVirtualHost().Domains = append(vs.GetVirtualHost().Domains, appUrlDomain)
		testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
			vs,
		}
	})

	JustBeforeEach(func() {
		authConfig = &extauth.AuthConfig{
			Metadata: &core.Metadata{
				Name:      getOAuth2ExtAuthExtension().GetConfigRef().Name,
				Namespace: getOAuth2ExtAuthExtension().GetConfigRef().Namespace,
			},
			Configs: []*extauth.AuthConfig_Config{{
				AuthConfig: &extauth.AuthConfig_Config_Oauth2{
					Oauth2: &extauth.OAuth2{
						OauthType: oauth2,
					},
				},
			}},
		}
		testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
			authConfig,
		}
		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
		oauth2Server.Stop()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	makeSingleRequest := func(g Gomega, client *http.Client) *http.Response {
		httpReqBuilder := testContext.GetHttpRequestBuilder().WithPath("success?foo=bar")
		resp, err := client.Do(httpReqBuilder.Build())
		g.Expect(err).NotTo(HaveOccurred())
		return resp
	}

	ExpectHappyPathToWork := func(makeSingleRequest func(g Gomega, client *http.Client) *http.Response, loginSuccessExpectation func()) {
		// do auth flow and make sure we have a cookie named cookie:
		appPageBuilder := testContext.GetHttpRequestBuilder()
		appPage := appPageBuilder.Build()

		var finalUrl *url.URL
		jar, err := cookiejar.New(nil)
		Expect(err).NotTo(HaveOccurred())
		cookieJar := &unsecureCookieJar{CookieJar: jar}

		// Using the default http client for its timeout, but using the builder since the DefaultHttpClient is a var which
		// causes potential issues when modifying it directly in tests.
		client := testutils.DefaultClientBuilder().Build()
		client.Jar = cookieJar
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			finalUrl = req.URL
			if len(via) > 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		}

		Eventually(func(g Gomega) *http.Response {
			return makeSingleRequest(g, client)
		}, "15s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))

		Expect(finalUrl).NotTo(BeNil())
		Expect(finalUrl.Path).To(Equal("/success"))
		// make sure query is passed through as well
		Expect(finalUrl.RawQuery).To(Equal("foo=bar"))

		// check the cookie jar
		tmpCookies := jar.Cookies(appPage.URL)
		Expect(tmpCookies).NotTo(BeEmpty())

		// grab the original cookies for these cookies, as `jar.Cookies` doesn't return
		// all the properties of the cookies
		for _, c := range tmpCookies {
			cookies = append(cookies, cookieJar.OriginalCookies[c.Name])
		}

		// make sure login is successful
		loginSuccessExpectation()
	}

	Context("redis for session store", func() {
		const (
			cookieName = "cookie"
		)
		var (
			redisInstance *redis.Instance
		)

		BeforeEach(func() {
			// update the config to use redis
			redisInstance = redisFactory.NewInstance()
			redisInstance.Run(testContext.Ctx())

			oauth2.Oauth2.Session = getRedisUserSession(redisInstance.Url(), cookieName)
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
	})

	Context("forward access token", func() {
		It("should work", func() {
			ExpectHappyPathToWork(makeSingleRequest, func() {})

			select {
			case r := <-testContext.TestUpstream().C:
				Expect(r.Headers).To(HaveKeyWithValue("Set-Cookie", ContainElements(ContainSubstring(fmt.Sprintf("access_token=%s", accessTokenValue)))))
			case <-time.After(time.Second):
				Fail("timedout")
			}
		})
	})

	Context("happy path with default settings (no redis)", func() {
		It("should get expected cookies", func() {
			ExpectHappyPathToWork(makeSingleRequest, func() {
				Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
					HaveKeyWithValue("access_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: true,
						Value:    accessTokenValue,
					})),
					HaveKeyWithValue("id_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: true,
					})),
				)))
			})
		})
	})

	Context("happy path with default settings and http only set to false", func() {
		BeforeEach(func() {
			oauth2.Oauth2.Session = &extauth.UserSession{
				Session: &extauth.UserSession_Cookie{Cookie: &extauth.UserSession_InternalSession{}},
				CookieOptions: &extauth.UserSession_CookieOptions{
					HttpOnly: &wrappers.BoolValue{Value: false},
				},
			}
		})

		It("should get expected cookies with HttpOnly set to false", func() {
			ExpectHappyPathToWork(makeSingleRequest, func() {
				Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
					HaveKeyWithValue("access_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: false,
						Value:    accessTokenValue,
					})),
					HaveKeyWithValue("id_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: false,
					})),
				)))
			})
		})
	})

	Context("happy path with default settings and client secret disabled", func() {
		BeforeEach(func() {
			oauth2.Oauth2.DisableClientSecret = &wrappers.BoolValue{Value: true}
			oauth2.Oauth2.ClientId = "no-secret-id"
			oauth2.Oauth2.Session = &extauth.UserSession{
				Session: &extauth.UserSession_Cookie{Cookie: &extauth.UserSession_InternalSession{}},
				CookieOptions: &extauth.UserSession_CookieOptions{
					HttpOnly: &wrappers.BoolValue{Value: false},
				},
			}
		})

		It("should work", func() {
			ExpectHappyPathToWork(makeSingleRequest, func() {
				Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
					HaveKeyWithValue("access_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: false,
						Value:    accessTokenValue,
					})),
					HaveKeyWithValue("id_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: false,
					})),
				)))
			})
		})
	})

	Context("happy path with default settings and allowing refreshing", func() {
		BeforeEach(func() {
			oauth2.Oauth2.Session = &extauth.UserSession{
				Session: &extauth.UserSession_Cookie{Cookie: &extauth.UserSession_InternalSession{
					AllowRefreshing: &wrappers.BoolValue{Value: true},
				}},
			}
		})

		It("should work", func() {
			ExpectHappyPathToWork(makeSingleRequest, func() {
				Expect(cookies).To(WithTransform(transforms.CookieDataMapper(), And(
					HaveKeyWithValue("access_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: true,
						Value:    accessTokenValue,
					})),
					HaveKeyWithValue("id_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: true,
					})),
					HaveKeyWithValue("refresh_token", matchers2.MatchCookieData(&transforms.CookieData{
						HttpOnly: true,
						Value:    refreshTokenValue,
					})),
				)))
			})
		})
	})

	Context("callbackPath test", func() {
		BeforeEach(func() {
			oauth2.Oauth2.CallbackPath = "/callback"
		})

		It("should exchange token with callbackPath", func() {
			finalPage := testContext.GetHttpRequestBuilder().WithPath("success").Build()
			finalPageUrl := finalPage.URL.String()

			client := testutils.DefaultClientBuilder().Build()
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			st := oauth2Service.NewStateSigner([]byte(testContext.ExtAuthInstance().GetServerSettings().SigningKey))
			signedState, err := st.Sign(oauth2.Oauth2.GetAppUrl())
			Expect(err).NotTo(HaveOccurred())

			callbackReqBuilder := testContext.GetHttpRequestBuilder().WithPath(fmt.Sprintf("callback?code=1234&gloo_urlToRedirect=%s&state=%s", url.QueryEscape(finalPageUrl), string(signedState)))

			Eventually(func(g Gomega) {
				resp, err := client.Do(callbackReqBuilder.Build())
				g.Expect(err).NotTo(HaveOccurred())

				expectedResponse := &matchers.HttpResponse{
					StatusCode: http.StatusFound,
					Headers: map[string]interface{}{
						"Location": finalPageUrl,
					},
				}
				g.Expect(resp).To(matchers.HaveHttpResponse(expectedResponse))
			}, "5s", "0.5s").Should(Succeed())
		})
	})

	Context("OAuth2 tests that don't forward to upstream", func() {
		var client *http.Client

		BeforeEach(func() {
			// initialize the client
			client = testutils.DefaultClientBuilder().Build()
		})

		It("should redirect to auth page", func() {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				// stop at the auth point
				if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
					return http.ErrUseLastResponse
				}
				return nil
			}

			httpReqBuilder := testContext.GetHttpRequestBuilder()
			Eventually(func(g Gomega) *http.Response {
				resp, err := client.Do(httpReqBuilder.Build())
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "15s", "0.5s").Should(HaveHTTPBody("auth"))
		})

		It("should include email scope in url", func() {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			httpReqBuilder := testContext.GetHttpRequestBuilder()
			Eventually(func(g Gomega) {
				resp, err := client.Do(httpReqBuilder.Build())
				g.Expect(err).NotTo(HaveOccurred())

				expectedResponse := &matchers.HttpResponse{
					StatusCode: http.StatusFound,
					Headers: map[string]interface{}{
						"Location": ContainSubstring("scope=email"),
					},
				}
				g.Expect(resp).To(matchers.HaveHttpResponse(expectedResponse))
			}, "15s", "0.5s").Should(Succeed())
		})

		It("should exchange token", func() {
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}

			finalPageReqBuilder := testContext.GetHttpRequestBuilder().WithPath("success")
			finalPageUrl := finalPageReqBuilder.Build().URL.String()
			st := oidc.NewStateSigner([]byte(testContext.ExtAuthInstance().GetServerSettings().SigningKey))
			signedState, err := st.Sign(finalPageUrl)
			Expect(err).NotTo(HaveOccurred())

			callbackReqBuilder := testContext.GetHttpRequestBuilder().WithPath(fmt.Sprintf("callback?code=1234&state=%s&gloo_urlToRedirect=%s", string(signedState), url.QueryEscape(finalPageUrl)))

			Eventually(func(g Gomega) {
				resp, err := client.Do(callbackReqBuilder.Build())
				g.Expect(err).NotTo(HaveOccurred())

				expectedResponse := &matchers.HttpResponse{
					StatusCode: http.StatusFound,
					Headers: map[string]interface{}{
						"Location": finalPageUrl,
					},
				}
				g.Expect(resp).To(matchers.HaveHttpResponse(expectedResponse))
			}, "5s", "0.5s").Should(Succeed())
		})
	})
})

func getOAuth2ExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oauth2-auth",
				Namespace: e2e.WriteNamespace,
			},
		},
	}
}
