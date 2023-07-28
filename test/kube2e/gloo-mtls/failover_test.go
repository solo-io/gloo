package gloo_mtls_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/solo-projects/test/kube2e/internal"
)

var _ = Describe("Failover Regression", func() {

	var (
		failoverTestContext *internal.FailoverTestContext
	)

	BeforeEach(func() {
		// The internal FailoverTestContext handles calling the TestContext's Before/AfterEach.
		testContext := testContextFactory.NewTestContext()
		failoverTestContext = &internal.FailoverTestContext{
			TestContext: testContext,
		}

		failoverTestContext.BeforeEach()
	})

	AfterEach(func() {
		failoverTestContext.AfterEach()
	})

	JustBeforeEach(func() {
		failoverTestContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		failoverTestContext.JustAfterEach()
	})

	internal.FailoverTests(func() *internal.FailoverTestContext {
		return failoverTestContext
	})

})
