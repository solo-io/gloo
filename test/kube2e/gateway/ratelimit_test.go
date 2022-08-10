package gateway_test

import (
	. "github.com/solo-io/solo-projects/test/kube2e/internal"

	. "github.com/onsi/ginkgo"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("RateLimit tests", func() {
	sharedInputs := RateLimitTestInputs{}

	BeforeEach(func() {
		sharedInputs.TestHelper = testHelper
	})

	Context("running rateLimit tests", func() {
		RunRateLimitTests(&sharedInputs)
	})

})
