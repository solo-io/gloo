package extauth_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-service/pkg/server"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
	envoy2 "github.com/solo-io/solo-projects/test/services/envoy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("Health Checker", func() {

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

	Context("health checker", func() {

		Context("without auth configs", func() {
			// This test is redundant due to how we run the ExtAuth service locally
			// Our new TestContextWithExtensions validates that the ExtAuth service is healthy
			// before proceeding.
		})

		Context("with auth configs", func() {
			BeforeEach(func() {
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
				testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{authConfig}
			})

			JustBeforeEach(func() {
				testContext.PatchDefaultVirtualService(func(vs *v1.VirtualService) *v1.VirtualService {
					builder := helpers.BuilderFromVirtualService(vs)
					builder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
						Extauth: GetBasicAuthExtension(),
					})
					return builder.Build()
				})
			})

			It("should pass after receiving xDS config from gloo", func() {
				testContext.ExtAuthInstance().EventuallyIsHealthy()
				Consistently(testContext.ExtAuthInstance().IsHealthy(), "3s", ".1s").Should(Equal(true))
			})
		})

		Context("shutdown", func() {
			BeforeEach(func() {
				// Allow the health check to return `envoy2.HealthCheckFailHeader` for two seconds, before shutting down
				testContext.ExtAuthInstance().UpdateServerSettings(func(settings *server.Settings) {
					settings.HealthCheckFailTimeout = 2 // seconds
				})
			})

			It("should fail healthcheck immediately on shutdown", func() {
				testContext.ExtAuthInstance().EventuallyIsHealthy()

				// Start sending health checking requests continuously
				waitForHealthcheck := make(chan struct{})
				go func(waitForHealthcheck chan struct{}) {
					defer GinkgoRecover()
					Eventually(func() bool {
						var header metadata.MD
						testContext.ExtAuthInstance().GetHealthCheckResponse(grpc.Header(&header))
						return len(header.Get(envoy2.HealthCheckFailHeader)) == 1
					}, "5s", ".1s").Should(BeTrue())
					waitForHealthcheck <- struct{}{}
				}(waitForHealthcheck)

				// Start the health checker first, then cancel
				time.Sleep(200 * time.Millisecond)
				testContext.CancelContext()
				Eventually(waitForHealthcheck, "5s", ".1s").Should(Receive())
			})
		})
	})

})
