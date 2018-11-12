package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-projects/projects/nginx"
)

var _ = Describe("HTTP", func() {
	Context("when it is empty", func() {
		It("should equal the empty HTTP", func() {
			http := Http{}
			Expect(http).To(Equal(Http{}))
		})
	})
})
