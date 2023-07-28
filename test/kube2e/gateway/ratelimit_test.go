package gateway_test

import (
	. "github.com/solo-io/solo-projects/test/kube2e/internal"

	. "github.com/onsi/ginkgo/v2"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("RateLimit tests", func() {
	var (
		sharedInputs RateLimitTestInputs
	)

	BeforeEach(func() {
		testContext := testContextFactory.NewTestContext()
		testContext.BeforeEach()

		sharedInputs = RateLimitTestInputs{
			TestContext: testContext,
		}
	})

	AfterEach(func() {
		sharedInputs.TestContext.AfterEach()
	})

	JustBeforeEach(func() {
		sharedInputs.TestContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		sharedInputs.TestContext.JustAfterEach()
	})

	Context("running rateLimit tests", func() {
		RunRateLimitTests(&sharedInputs)
	})

})
