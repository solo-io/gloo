package gateway_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/solo-io/solo-projects/test/kube2e/internal"
)

// Regular extAuth as run in other suites
var _ = Describe("ExtAuth tests", func() {
	var (
		sharedInputs ExtAuthTestInputs
	)

	BeforeEach(func() {
		testContext := testContextFactory.NewTestContext()
		testContext.BeforeEach()

		sharedInputs = ExtAuthTestInputs{
			TestContext:    testContext,
			ShouldTestLDAP: true,
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

	Context("running ExtAuth tests", func() {
		RunExtAuthTests(&sharedInputs)
	})

})
