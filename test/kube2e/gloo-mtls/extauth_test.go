package gloo_mtls_test

import (
	. "github.com/solo-io/solo-projects/test/kube2e/internal"

	. "github.com/onsi/ginkgo/v2"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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

	Context("running ExtAuth tests", func() {
		RunExtAuthTests(&sharedInputs)
	})

})
