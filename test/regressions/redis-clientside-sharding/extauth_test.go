package clientside_sharding_test

import (
	. "github.com/solo-io/solo-projects/test/regressions/internal"

	. "github.com/onsi/ginkgo"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

// Regular extAuth as run in other suites
var _ = Describe("ExtAuth tests", func() {
	sharedInputs := ExtAuthTestInputs{}

	BeforeEach(func() {
		sharedInputs.TestHelper = testHelper
	})

	Context("running ExtAuth tests", func() {
		RunExtAuthTests(&sharedInputs)
	})

})
