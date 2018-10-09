package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("HTTP context", func() {
	It("should be an empty context", func() {
		httpContext, err := GenerateHttpContext()
		Expect(err).NotTo(HaveOccurred())
		expected := `
http {
}
`
		Expect(string(httpContext)).To(Equal(expected))
	})
})
