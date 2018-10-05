package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("Ingress Nginx config", func() {
	It("should be able to contain a server", func() {
		server := Server{}
		config := IngressNginxConfig{
			Server: server,
		}
		Expect(config).NotTo(BeNil())
		Expect(config.Server).NotTo(BeNil())
		Expect(config.Server).To(Equal(Server{}))
	})
})
