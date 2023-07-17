package extauth_test

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	testServers "github.com/solo-io/solo-projects/test/services/extauth/servers"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"

	gloohelpers "github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("Deprecated OAuth", func() {
	/*
		These tests are for the deprecated OAuth auth config which have been replaced by the OAuth2 auth config.
		These tests are similar to those in the OIDC tests (which use OAuth2), but have been moved to their own file to
		  lessen confusion, make it easier to remove in the future, and so the OIDC tests focus on our newer config.
		As this is deprecated, we shouldn't need to update these tests when adding functionality to the OAuth2 extauth.
	*/

	var (
		testContext           *e2e.TestContextWithExtensions
		authConfig            *extauth.AuthConfig
		oauth                 *extauth.OAuth
		privateKey            *rsa.PrivateKey
		discoveryServer       testServers.FakeDiscoveryServer
		secret                *gloov1.Secret
		token                 string
		virtualServiceBuilder *gloohelpers.VirtualServiceBuilder
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

		discoveryServer = testServers.FakeDiscoveryServer{
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

		oauth = discoveryServer.GetOauthConfig(testContext.EnvoyInstance().HttpPort, "localhost", secret.Metadata.Ref())

		testContext.ResourcesToCreate().Secrets = append(testContext.ResourcesToCreate().Secrets, secret)

		// get id token
		token = discoveryServer.GenerateValidIdToken(privateKey)
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

	Context("oidc (old config)", func() {
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
					AuthConfig: &extauth.AuthConfig_Config_Oauth{
						Oauth: oauth,
					},
				}},
			}
			_, err := testContext.TestClients().AuthConfigClient.Write(authConfig, clients.WriteOpts{Ctx: testContext.Ctx()})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("Oidc tests that don't forward to upstream", func() {
			It("should redirect to auth page", func() {
				client := testutils.DefaultClientBuilder().Build()
				client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					// stop at the auth point
					if req.Response != nil && req.Response.Header.Get("x-auth") != "" {
						return http.ErrUseLastResponse
					}
					return nil
				}

				httpReqBuilder := testContext.GetHttpRequestBuilder()
				Eventually(func() (*http.Response, error) {
					resp, err := client.Do(httpReqBuilder.Build())
					if err != nil {
						return nil, err
					}
					_, _ = fmt.Fprintf(GinkgoWriter, "headers are %v \n", resp.Header)
					return resp, nil
				}, "10s", "0.5s").Should(HaveHTTPBody("auth"))
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
		})

		Context("Oidc tests that do forward to upstream", func() {
			It("should allow access with proper jwt token", func() {
				ExpectUpstreamRequest()
			})
		})

	})

	Context("oidc + opa (old config)", func() {
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
							AuthConfig: &extauth.AuthConfig_Config_Oauth{
								Oauth: discoveryServer.GetOauthConfig(testContext.EnvoyInstance().HttpPort, "localhost", secret.Metadata.Ref()),
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
							AuthConfig: &extauth.AuthConfig_Config_Oauth{
								Oauth: discoveryServer.GetOauthConfig(testContext.EnvoyInstance().HttpPort, "localhost", secret.Metadata.Ref()),
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
				httpReqBuilder := testContext.GetHttpRequestBuilder().WithHeader("Authorization", "Bearer "+token).WithHost("")
				EventuallyWithOffset(1, func() (*http.Response, error) {
					resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
					if err != nil {
						return nil, err
					}
					return resp, nil
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})

		})
	})
})
