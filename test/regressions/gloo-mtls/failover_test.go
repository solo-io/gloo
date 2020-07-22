package gloo_mtls_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/solo-io/solo-projects/test/regressions/internal"
)

var _ = Describe("Failover Regression", func() {
	var (
		failoverTest *internal.FailoverTest
	)
	BeforeEach(func() {
		failoverTest = internal.FailoverBeforeEach(testHelper)
	})

	AfterEach(func() {
		internal.FailoverAfterEach(failoverTest, testHelper)
	})

	It("can failover to kubernetes EDS endpoints", func() {
		internal.FailoverSpec(failoverTest, testHelper)
	})

})
