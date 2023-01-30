package gateway_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	"github.com/solo-io/solo-projects/test/kube2e/internal"
)

var _ = Describe("Failover Regression", func() {
	var (
		failoverTest *internal.FailoverTest
		ctx          context.Context
		cancel       context.CancelFunc
	)
	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		failoverTest = internal.FailoverBeforeEach(testHelper)
	})

	AfterEach(func() {
		internal.FailoverAfterEach(ctx, failoverTest, testHelper)
		cancel()
	})

	It("can failover to kubernetes EDS endpoints", func() {
		internal.FailoverSpec(failoverTest, testHelper)
	})

})
