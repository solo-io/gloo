package gateway_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/solo-io/solo-projects/test/regressions/internal"
)

// Regular extAuth as run in other suites
var _ = Describe("ExtAuth tests", func() {
	sharedInputs := ExtAuthTestInputs{}

	BeforeEach(func() {
		sharedInputs.TestHelper = testHelper
		sharedInputs.ShouldTestLDAP = true
	})

	Context("running ExtAuth tests", func() {
		RunExtAuthTests(&sharedInputs)
	})

})
