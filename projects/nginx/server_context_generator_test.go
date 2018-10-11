package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("Server context", func() {
	Context("when the `Server` instance is empty", func() {
		It("should be an empty context", func() {
			server := &Server{}
			serverContext, err := GenerateServerContext(server)
			Expect(err).NotTo(HaveOccurred())
			expected := `
server {
}
`
			Expect(string(serverContext)).To(Equal(expected))
		})
	})
	Context("when the `Server` instance contains a single `Location`", func() {
		It("should contain a single location context", func() {
			location := &Location{
				Prefix: "/",
				Root:   "/data/www",
			}
			server := &Server{
				Location: location,
			}
			serverContext, err := GenerateServerContext(server)
			Expect(err).NotTo(HaveOccurred())
			expected := `
server {
    location / {
        root /data/www;
    }
}
`
			Expect(string(serverContext)).To(Equal(expected))
		})
	})
})
