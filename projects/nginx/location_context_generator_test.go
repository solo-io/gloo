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
			locationContext, err := GenerateLocationContext(location)
			Expect(err).NotTo(HaveOccurred())
			expected := `
location / {
    root /data/www;
}
`
			Expect(string(locationContext)).To(Equal(expected))
		})
	})
})
