package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("HTTP context", func() {
	Context("when the `Http` instance is empty", func() {
		It("should be an empty context", func() {
			http := &Http{}
			httpContext, err := GenerateHttpContext(http)
			Expect(err).NotTo(HaveOccurred())
			expected := `
http {
}
`
			Expect(string(httpContext)).To(Equal(expected))
		})
	})
	Context("when the `Http` instance contains a `Server`", func() {
		It("should contain a server context", func() {
			server := &Server{}
			http := &Http{
				Server: server,
			}
			httpContext, err := GenerateHttpContext(http)
			Expect(err).NotTo(HaveOccurred())
			expected := `
http {
	server {
	}
}
`
			Expect(string(httpContext)).To(Equal(expected))
		})
	})
})
