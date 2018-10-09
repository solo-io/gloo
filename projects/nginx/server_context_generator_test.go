package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("Server context", func() {
	It("should be an empty context", func() {
		serverContext, err := GenerateServerContext()
		Expect(err).NotTo(HaveOccurred())
		expected := `
server {
}
`
		Expect(string(serverContext)).To(Equal(expected))
	})
})
