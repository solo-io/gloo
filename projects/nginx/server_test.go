package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-projects/projects/nginx"
)

var _ = Describe("Server", func() {
	It("should be empty", func() {
		server := Server{}
		Expect(server).To(Equal(Server{}))
	})
})
