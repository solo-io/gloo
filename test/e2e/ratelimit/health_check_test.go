package ratelimit_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/test/e2e"
	"github.com/solo-io/solo-projects/test/services/envoy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"time"
)

var _ = Describe("Rate Limiter Health Checks", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			RateLimit: true,
		})
		testContext.BeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	Context("should pass after receiving xDS config from gloo", func() {

		It("without rate limit configs", func() {
			// This test is redundant due to how we run the RateLimit service locally
			// Our new TestContextWithExtensions validates that the RateLimit service is healthy
			// before proceeding.
		})

		It("with rate limit configs", func() {
			// This test is redundant due to how our e2e tests work
			// Our new TestContextWithExtensions constructs the RateLimit filter with
			// DenyOnFail=true (https://github.com/solo-io/gloo/blob/e95be056dc12db0a6a0995540c0c7c43b388f0d2/projects/gloo/api/v1/enterprise/options/ratelimit/ratelimit.proto#L26)
			// This means that if the RateLimit service is unhealthy we would see tests fail
		})

	})

	Context("shutdown", func() {

		BeforeEach(func() {
			testContext.RateLimitInstance().UpdateServerSettings(func(settings *server.Settings) {
				settings.HealthFailTimeout = 2 // seconds
			})
		})

		It("should fail healthcheck immediately on shutdown", func() {
			testContext.RateLimitInstance().EventuallyIsHealthy()

			// Start sending health checking requests continuously
			waitForHealthcheckFail := make(chan struct{})
			go func(waitForHealthcheckFail chan struct{}) {
				defer GinkgoRecover()

				Eventually(func() (bool, error) {
					var header metadata.MD
					_, _ = testContext.RateLimitInstance().GetHealthCheckResponse(grpc.Header(&header))
					return len(header.Get(envoy.HealthCheckFailHeader)) == 1, nil
				}, "5s", ".1s").Should(BeTrue())
				waitForHealthcheckFail <- struct{}{}
			}(waitForHealthcheckFail)

			time.Sleep(100 * time.Millisecond) // Allow the above goroutine validating health checks to start
			testContext.CancelContext()

			Eventually(waitForHealthcheckFail, "5s", ".1s").Should(Receive())
		})

	})

})
