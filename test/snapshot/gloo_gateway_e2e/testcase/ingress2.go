package testcase

import (
	. "github.com/onsi/ginkgo/v2"
)

var Test = func(input TestInput) {
	BeforeEach(func() {
		input.ApplyConfig()
	})

	It("Get 200 response", func() {

	})
}
