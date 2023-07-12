package extauth_test

import (
	"fmt"
	"net/http"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/e2e"
)

/*
	- The translation logic for HMAC authorization is found here: https://github.com/solo-io/ext-auth-service/blob/cb9dd01fba4a805b19421c48c87b338961125963/pkg/controller/translation/hmac.go#L53-L81
	- https://www.devglan.com/online-tools/hmac-sha256-online can be used to generate the HMAC signature. The data below is
	  the data used to generate the valid HMAC signature for the secret we use in our tests:
		- Plaintext
			date: Thu, 22 Jun 2017 17:15:21 GMT
			GET /requests HTTP/1.1
		- Secret key:
			secret
		- Cryptographic Hash Function:
			SHA-256
*/

var _ = Describe("HMAC", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	const (
		validUser      = "alice"
		validSignature = "ujWCGHeec9Xd6UD2zlyxiNMCiXnDOWeVFMu5VeRUxtw="
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		secret := &gloov1.Secret{
			Metadata: &core.Metadata{
				Name:      "hmac-secret",
				Namespace: "default",
			},
			Kind: &gloov1.Secret_Credentials{
				Credentials: &gloov1.AccountCredentialsSecret{
					Username: validUser,
					Password: "secret",
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
		testContext.ResourcesToCreate().VirtualServices = v1.VirtualServiceList{
			vsBuilder.Build(),
		}
		testContext.ResourcesToCreate().Secrets = append(testContext.ResourcesToCreate().Secrets, secret)
		testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
			authConfig,
		}
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

	eventuallyConsistentResponse := func(req *http.Request, expectedStatus int) {
		EventuallyWithOffset(1, func() (*http.Response, error) {
			resp, err := testutils.DefaultHttpClient.Do(req)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}).Should(HaveHTTPStatus(expectedStatus))
		ConsistentlyWithOffset(1, func() (*http.Response, error) {
			resp, err := testutils.DefaultHttpClient.Do(req)
			if err != nil {
				return nil, err
			}
			return resp, nil
		}, "5s", "1s").Should(HaveHTTPStatus(expectedStatus))
	}

	Context("with SHA-256 and list of secret refs", func() {
		buildHmacRequest := func(user, signature string) *http.Request {
			reqBuilder := testContext.GetHttpRequestBuilder().WithPath("requests").
				WithRawHeader("date", "Thu, 22 Jun 2017 17:15:21 GMT").
				WithRawHeader("authorization", fmt.Sprintf("hmac username=\"%s\", algorithm=\"hmac-sha256\", headers=\"date @request-target\", signature=\"%s\"", user, signature))
			return reqBuilder.Build()
		}

		It("Allows requests with valid signature", func() {
			// the signature is explained in the comment at the top of this file
			req := buildHmacRequest(validUser, validSignature)
			eventuallyConsistentResponse(req, http.StatusOK)
		})

		It("Denies requests without valid signature", func() {
			// Generated this signature same as above, but with the wrong secret (`secret123`).
			req := buildHmacRequest(validUser, "jypxis61NLyOvHGayXVvp/TiTR96d1as8cJLdqltIGk=")
			eventuallyConsistentResponse(req, http.StatusUnauthorized)
		})

		It("Denies requests without valid user", func() {
			req := buildHmacRequest(validUser+"-invalid", validSignature)
			eventuallyConsistentResponse(req, http.StatusUnauthorized)
		})
	})
})
