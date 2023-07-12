package gateway_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/solo-io/solo-projects/test/kube2e/internal"
)

var _ = Describe("Failover Regression", func() {

	var (
		failoverTestContext *internal.FailoverTestContext
	)

	BeforeEach(func() {
		failoverTestContext = &internal.FailoverTestContext{
			TestHelper:        testHelper,
			ResourceClientset: resourceClientset,
			SnapshotWriter:    snapshotWriter,
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
