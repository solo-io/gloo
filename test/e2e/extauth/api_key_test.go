package extauth_test

import (
	"net/http"

	"github.com/solo-io/gloo/test/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("API Key", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		secret := &gloov1.Secret{
			Metadata: &core.Metadata{
				Name:      "secret1",
				Namespace: "default",
			},
			Kind: &gloov1.Secret_ApiKey{
				ApiKey: &extauth.ApiKey{
					ApiKey: "secretApiKey1",
				},
			},
		}
		authConfig := &extauth.AuthConfig{
			Metadata: &core.Metadata{
				Name:      getApiKeyExtAuthExtension().GetConfigRef().Name,
				Namespace: getApiKeyExtAuthExtension().GetConfigRef().Namespace,
			},
			Configs: []*extauth.AuthConfig_Config{{
				AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
					ApiKeyAuth: &extauth.ApiKeyAuth{
						ApiKeySecretRefs: []*core.ResourceRef{
							secret.Metadata.Ref(),
						},
						LabelSelector: map[string]string{"team": "infrastructure"},
					},
				},
			}},
		}
		labeledSecret := &gloov1.Secret{
			Metadata: &core.Metadata{
				Name:      "secret2",
				Namespace: "default",
				Labels:    map[string]string{"team": "infrastructure"},
			},
			Kind: &gloov1.Secret_ApiKey{
				ApiKey: &extauth.ApiKey{
					ApiKey: "secretApiKey2",
				},
			},
		}

		// add the secrets from above to the to-be-written secret list
		secrets := []*gloov1.Secret{secret, labeledSecret}

		// use the default virtual service and add the ext auth extension
		vs := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
		vs.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
			Extauth: getApiKeyExtAuthExtension(),
		})

		testContext.ResourcesToCreate().Secrets = secrets
		testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
			authConfig,
		}
		testContext.ResourcesToCreate().VirtualServices = v1.VirtualServiceList{
			vs.Build(),
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

	It("should deny ext auth envoy without apikey", func() {
		requestBuilder := testContext.GetHttpRequestBuilder()
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(requestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))
	})

	It("should deny ext auth envoy with incorrect apikey", func() {
		requestBuilder := testContext.GetHttpRequestBuilder().WithHeader("api-key", "badApiKey")
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(requestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))
	})

	It("should accept ext auth envoy with correct apikey -- secret ref match", func() {
		requestBuilder := testContext.GetHttpRequestBuilder().WithHeader("api-key", "secretApiKey1")
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(requestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
	})

	It("should accept ext auth envoy with correct apikey -- label match", func() {
		requestBuilder := testContext.GetHttpRequestBuilder().WithHeader("api-key", "secretApiKey2")
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(requestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
	})

})

func getApiKeyExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "apikey-auth",
				Namespace: e2e.WriteNamespace,
			},
		},
	}
}
