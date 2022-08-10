package clientside_sharding_test

import (
	. "github.com/solo-io/solo-projects/test/kube2e/internal"

	. "github.com/onsi/ginkgo"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Regular rate_limit tests as run in other suites
var _ = Describe("RateLimit tests", func() {
	sharedInputs := RateLimitTestInputs{}

	BeforeEach(func() {
		sharedInputs.TestHelper = testHelper
	})

	Context("running rateLimit tests", func() {
		RunRateLimitTests(&sharedInputs)
	})

})
