package extauth_test

import (
	"context"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("ExtAuth Translation - Concurrency Tests", func() {

	var (
		_      context.Context
		cancel context.CancelFunc
	)

	BeforeEach(func() {
		_, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	Context("TODO", func() {
		// TODO (sam-heilbron)
		// Write tests to demonstrate the our extension is not (yet) thread-safe
		// This should fail if someone makes a change that assumes thread safety
	})

})
