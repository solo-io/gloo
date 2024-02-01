package extauth_test

import (
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

var _ = Describe("Basic Auth", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		authConfig := &extauth.AuthConfig{
			Metadata: &core.Metadata{
				Name:      GetBasicAuthExtension().GetConfigRef().Name,
				Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
			},
			Configs: []*extauth.AuthConfig_Config{{
				AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
					BasicAuth: getBasicAuthConfig(),
				},
			}},
		}

		vs := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
		vs.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
			Extauth: GetBasicAuthExtension(),
		})

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

	It("should deny ext auth envoy without credentials", func() {
		httpRequestBuilder := testContext.GetHttpRequestBuilder()
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))
	})

	It("should allow ext auth envoy", func() {
		httpRequestBuilder := testContext.GetHttpRequestBuilder().WithHostname("user:password@localhost")
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusOK))
	})

	It("should deny ext auth with wrong password", func() {
		httpRequestBuilder := testContext.GetHttpRequestBuilder().WithHostname("user:password2@localhost")
		Eventually(func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(http.StatusUnauthorized))
	})
})