package extauth_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

/*
	TODO: Fix the skipped test. It fails because the HMAC signature is invalid, at least based off of the response "HMAC signature missing or invalid".
	- The original tests in extauth_test.go passed, so this is likely due to my setup. Will ask @seth &| @nathan for help.
	- I used https://www.devglan.com/online-tools/hmac-sha256-online to generate the HMAC signature and tried different order
		- Plaintext
			date: Thu, 22 Jun 2017 17:15:21 GMT
			GET /requests HTTP/1.1
		- Secret key: secret
		- SHA-256
*/

var _ = Describe("HMAC", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("hmac tests with sha1 and list of secret refs", func() {
		BeforeEach(func() {
			secret := &gloov1.Secret{
				Metadata: &core.Metadata{
					Name:      "secret",
					Namespace: "default",
				},
				Kind: &gloov1.Secret_Credentials{
					Credentials: &gloov1.AccountCredentialsSecret{
						Username: "alice",
						Password: "secret123",
					},
				},
			}
			authConfig := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      "hmac-auth",
					Namespace: e2e.WriteNamespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_HmacAuth{
						HmacAuth: &extauth.HmacAuth{
							ImplementationType: &extauth.HmacAuth_ParametersInHeaders{
								ParametersInHeaders: &extauth.HmacParametersInHeaders{},
							},
							SecretStorage: &extauth.HmacAuth_SecretRefs{
								SecretRefs: &extauth.SecretRefList{
									SecretRefs: []*core.ResourceRef{
										secret.Metadata.Ref(),
									},
								},
							},
						},
					},
				}},
			}

			// Update the default virtual service to use the hmac auth extension
			vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
			vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
				Extauth: &extauth.ExtAuthExtension{
					Spec: &extauth.ExtAuthExtension_ConfigRef{
						ConfigRef: authConfig.Metadata.Ref(),
					},
				},
			})
			virtualService := vsBuilder.Build()

			testContext.ResourcesToCreate().Secrets = gloov1.SecretList{
				secret,
			}
			testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
				authConfig,
			}
			testContext.ResourcesToCreate().VirtualServices = v1.VirtualServiceList{
				virtualService,
			}
		})

		buildHmacRequest := func(user, signature string) *http.Request {
			reqBuilder := testContext.GetHttpRequestBuilder().WithPath("requests")
			req := reqBuilder.Build()

			// Adding the authorization and date header outside the builder. The builder's `WithHeader` method uses the `,` as a delimiter to create multi-value headers.
			req.Header.Add("date", "Thu, 22 Jun 2017 17:15:21 GMT")
			req.Header.Add("Authorization", fmt.Sprintf("hmac username=\"%s\", algorithm=\"hmac-sha256\", headers=\"date @request-target\", signature=\"%s\"", user, signature))
			//req.Header.Add("Host", "hmac.com")
			//req.Body = nil
			return req
		}

		// TODO: Uncomment once fixed
		//It("Allows requests with valid signature", func() {
		//	req := buildHmacRequest("alice123", "ujWCGHeec9Xd6UD2zlyxiNMCiXnDOWeVFMu5VeRUxtw=")
		//	Eventually(func(g Gomega) *http.Response {
		//		resp, err := http.DefaultClient.Do(req)
		//		g.Expect(err).NotTo(HaveOccurred())
		//		return resp
		//	}, "15s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
		//})

		It("Denies requests without valid signature", func() {
			req := buildHmacRequest("alice123", "notreallyalice")
			Eventually(func(g Gomega) *http.Response {
				resp, err := testutils.DefaultHttpClient.Do(req)
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}, "15s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))
		})
	})

})
