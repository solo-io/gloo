package extauth_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	testServers "github.com/solo-io/solo-projects/test/services/extauth/servers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	oauth_utils "github.com/solo-io/ext-auth-service/pkg/config/oauth/test_utils"
	"github.com/solo-io/ext-auth-service/pkg/config/oauth/token_validation"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("OAuth2 Token Introspection", func() {

	var (
		testContext         *e2e.TestContextWithExtensions
		introspectionServer *testServers.TokenIntrospectionServer
		authConfig          *extauth.AuthConfig
		virtualService      *v1.VirtualService
		// This test utilizes multiple secrets. We have this local variable to store N secrets and write them in the `JustBeforeEach` phase.
		secrets []*gloov1.Secret
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		// Initializes the server, which generates a port to run on.
		// This can then be used to get auth configs with correct addresses with `GetOauthTokenIntrospectionConfig`.
		introspectionServer = testServers.NewTokenIntrospectionServer()

		// initialize the secrets list to be empty
		secrets = []*gloov1.Secret{}
	})

	JustBeforeEach(func() {
		// Start the auth server
		introspectionServer.Start()

		testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
			authConfig,
		}
		testContext.ResourcesToCreate().VirtualServices = v1.VirtualServiceList{
			virtualService,
		}
		testContext.ResourcesToCreate().Secrets = secrets

		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
		introspectionServer.Stop()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	// Execute a request with an access token, against an endpoint that requires token authentication and return the response with error.
	requestWithAccessToken := func(token string) (*http.Response, error) {
		requestBuilder := testContext.GetHttpRequestBuilder().WithHeader("Authorization", fmt.Sprintf("Bearer %s", token))
		return testutils.DefaultHttpClient.Do(requestBuilder.Build())
	}

	Context("using IntrospectionUrl", func() {
		BeforeEach(func() {
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
					Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_Oauth2{
						Oauth2: &extauth.OAuth2{
							OauthType: introspectionServer.GetOauthTokenIntrospectionUrlConfig(),
						},
					},
				}},
			}

			vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
			vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
				Extauth: getOauthTokenIntrospectionExtAuthExtension(),
			})
			virtualService = vsBuilder.Build()
		})

		It("should accept introspection url with valid access token", func() {
			Eventually(func(g Gomega) *http.Response {
				resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
			Consistently(func(g Gomega) *http.Response {
				resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
		})

		It("should deny introspection url with invalid access token", func() {
			Eventually(func(g Gomega) *http.Response {
				resp, err := requestWithAccessToken("invalid-access-token")
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			Consistently(func(g Gomega) *http.Response {
				resp, err := requestWithAccessToken("invalid-access-token")
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
		})

	})

	Context("using Introspection", func() {

		createBasicAuthHandler := func(validToken, clientId, clientSecret string) func(http.ResponseWriter, *http.Request) {
			return func(writer http.ResponseWriter, request *http.Request) {
				err := request.ParseForm()
				if err != nil {
					panic(err)
				}

				requestedToken := request.Form.Get("token")
				requestedClientId := request.Form.Get("client_id")
				requestedClientSecret := request.Form.Get("client_secret")
				fmt.Fprintln(GinkgoWriter, "token request", request.Form)
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
			return &gloov1.Secret{
				Metadata: &core.Metadata{
					Name:      name,
					Namespace: "default",
				},
				Kind: &gloov1.Secret_Oauth{
					Oauth: &extauth.OauthSecret{
						ClientSecret: value,
					},
				},
			}
		}

		BeforeEach(func() {
			secret := getOauthSecret("secret", "client-secret")
			secrets = append(secrets, secret)

			// Create an auth server that requires clients to provide credentials (client-id and client-secret)
			introspectionServer.AuthHandlers = &oauth_utils.AuthHandlers{
				TokenIntrospectionHandler: createBasicAuthHandler(testServers.IntrospectionAccessToken, "client-id", "client-secret"),
			}

			// Create an auth config, with proper references to the client credentials
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
					Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_Oauth2{
						Oauth2: &extauth.OAuth2{
							OauthType: introspectionServer.GetOauthTokenIntrospectionConfig("client-id", secret.Metadata.Ref(), false),
						},
					},
				}},
			}

			vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
			vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
				Extauth: getOauthTokenIntrospectionExtAuthExtension(),
			})
			virtualService = vsBuilder.Build()
		})

		When("auth config includes valid credentials", func() {
			// The default auth config that we initialize is valid

			It("should accept introspection with valid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
			})

			It("should deny introspection with invalid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})
		})

		When("auth config includes invalid credentials", func() {
			BeforeEach(func() {
				// Create a client secret with the wrong value
				invalidSecret := getOauthSecret("invalid-secret", "invalid-client-secret")
				secrets = append(secrets, invalidSecret)

				// Set the auth config to reference that invalid secret
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
						Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
					},
					Configs: []*extauth.AuthConfig_Config{{
						AuthConfig: &extauth.AuthConfig_Config_Oauth2{
							Oauth2: &extauth.OAuth2{
								OauthType: introspectionServer.GetOauthTokenIntrospectionConfig("client-id", invalidSecret.Metadata.Ref(), false),
							},
						},
					}},
				}
			})

			It("should deny introspection with valid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})
			It("should deny introspection with invalid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})
		})

		When("auth config is missing credentials and auth server doesn't require credentials", func() {

			BeforeEach(func() {
				// Resetting the AuthHandlers, because the Context's BeforeEach sets them up, and we don't want them for these tests.
				introspectionServer.AuthHandlers = &oauth_utils.AuthHandlers{}
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
						Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
					},
					Configs: []*extauth.AuthConfig_Config{{
						AuthConfig: &extauth.AuthConfig_Config_Oauth2{
							Oauth2: &extauth.OAuth2{
								OauthType: introspectionServer.GetOauthTokenIntrospectionConfig("", nil, false),
							},
						},
					}},
				}
			})

			It("should accept introspection with valid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
			})

			It("should deny introspection with invalid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})
		})

		When("Client secret is disabled", func() {
			BeforeEach(func() {
				introspectionServer.AuthHandlers = &oauth_utils.AuthHandlers{
					TokenIntrospectionHandler: createBasicAuthHandler(testServers.IntrospectionAccessToken, "no-secret-id", ""),
				}

				// Create an auth config, with proper references to the client credentials
				authConfig = &extauth.AuthConfig{
					Metadata: &core.Metadata{
						Name:      getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Name,
						Namespace: getOauthTokenIntrospectionExtAuthExtension().GetConfigRef().Namespace,
					},
					Configs: []*extauth.AuthConfig_Config{{
						AuthConfig: &extauth.AuthConfig_Config_Oauth2{
							Oauth2: &extauth.OAuth2{
								OauthType: introspectionServer.GetOauthTokenIntrospectionConfig("no-secret-id", nil, true),
							},
						},
					}},
				}
			})
			It("should accept introspection with valid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken(testServers.IntrospectionAccessToken)
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
			})

			It("should deny introspection with invalid access token", func() {
				Eventually(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
				Consistently(func(g Gomega) *http.Response {
					resp, err := requestWithAccessToken("invalid-access-token")
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "3s", "0.5s").Should(HaveHTTPStatus(http.StatusForbidden))
			})
		})

	})
})

func getOauthTokenIntrospectionExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oauth-token-introspection",
				Namespace: e2e.WriteNamespace,
			},
		},
	}
}
