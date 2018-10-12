package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("Location context", func() {
	Context("when the `Location` instance contains a `Prefix` and a `Root`", func() {
		It("should contain a prefix and a root", func() {
			location := &Location{
				Prefix: "/",
				Root:   "/data/www",
			}
			locationContext, err := location.GenerateContext()
			Expect(err).NotTo(HaveOccurred())
			expected := `location / {
    root /data/www;
}`
			Expect(string(locationContext)).To(Equal(expected))
		})
	})
	Context("when the `Location` instance contains a `Prefix` and a `ProxyPass`", func() {
		It("should contain a prefix and a proxy pass URL", func() {
			location := &Location{
				Prefix:    "/",
				ProxyPass: "http://localhost:8080/",
			}
			locationContext, err := location.GenerateContext()
			Expect(err).NotTo(HaveOccurred())
			expected := `location / {
    proxy_pass http://localhost:8080/;
}`
			Expect(string(locationContext)).To(Equal(expected))
		})
	})
})
