package consul_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-storage/consul"
)

var _ = Describe("ConsulStorageClient", func() {
	Context("Upstreams", func() {
		Context("create", func() {
			It("creates the upstream as a consul key", func() {
				client, err := NewStorage(api.DefaultConfig(), "foo", time.Millisecond)
				Expect(err).NotTo(HaveOccurred())
				input := &v1.Upstream{
					Name:              "myupstream",
					Type:              "foo",
					ConnectionTimeout: time.Second,
				}
				us, err := client.V1().Upstreams().Create(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(us).NotTo(Equal(input))
				input.Metadata = us.Metadata
				Expect(us).To(Equal(input))
			})
		})
	})
})
