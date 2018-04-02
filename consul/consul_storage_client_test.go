package consul_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-storage/consul"
	"github.com/solo-io/gloo-testing/helpers"
)

var _ = Describe("ConsulStorageClient", func() {
	var rootPath string
	var consul *api.Client
	BeforeEach(func() {
		rootPath = helpers.RandString(4)
		c, err := api.NewClient(api.DefaultConfig())
		Expect(err).NotTo(HaveOccurred())
		consul = c
	})
	AfterEach(func() {
		consul.KV().DeleteTree(rootPath, nil)
	})
	Context("Upstreams", func() {
		Context("create", func() {
			It("creates the upstream as a consul key", func() {
				client, err := NewStorage(api.DefaultConfig(), rootPath, time.Millisecond)
				Expect(err).NotTo(HaveOccurred())
				input := &v1.Upstream{
					Name:              "myupstream",
					Type:              "foo",
					ConnectionTimeout: time.Second,
				}
				us, err := client.V1().Upstreams().Create(input)
				Expect(err).NotTo(HaveOccurred())
				Expect(us).NotTo(Equal(input))
				p, _, err := consul.KV().Get(rootPath+"/upstreams/"+input.Name, nil)
				Expect(err).NotTo(HaveOccurred())
				var unmarshalledUpstream v1.Upstream
				err = proto.Unmarshal(p.Value, &unmarshalledUpstream)
				Expect(err).NotTo(HaveOccurred())
				Expect(&unmarshalledUpstream).To(Equal(input))
				resourceVersion := fmt.Sprintf("%v", p.CreateIndex)
				Expect(us.Metadata.ResourceVersion).To(Equal(resourceVersion))
				input.Metadata = us.Metadata
				Expect(us).To(Equal(input))
			})
		})
	})
})
