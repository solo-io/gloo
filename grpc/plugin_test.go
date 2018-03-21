package grpc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-plugins/grpc"
	"github.com/solo-io/gloo/pkg/log"
)

var _ = Describe("Plugin", func() {
	It("unmarshal file descriptor", func() {
		b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
		Expect(err).NotTo(HaveOccurred())
		descriptor, err := ConvertProto(b)
		Expect(err).NotTo(HaveOccurred())
		log.Printf("%v", FuncsForProto("Bookstore", descriptor))
	})
	It("unmarshal file descriptor", func() {
		b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
		Expect(err).NotTo(HaveOccurred())
		descriptor, err := ConvertProto(b)
		Expect(err).NotTo(HaveOccurred())
		AddHttpRulesToProto("my-upstream", "Bookstore", descriptor)
		log.Printf("%v", FuncsForProto("Bookstore", descriptor))
	})
	FIt("returns a dependency with a file ref for each descriptor file "+
		"specified in the spec", func() {
		us := &v1.Upstream{
			ServiceInfo: &v1.ServiceInfo{
				Type: ServiceTypeGRPC,
				Properties: EncodeServiceProperties(ServiceProperties{
					DescriptorsFileRef: "file_1",
				}),
			},
		}
		plugin := &Plugin{}
		deps := plugin.GetDependencies(&v1.Config{Upstreams: []*v1.Upstream{us}})
		Expect(deps.FileRefs).To(HaveLen(1))
		Expect(deps.FileRefs[0]).To(Equal("file_1"))
	})
})
