package grpc

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugin"
)

var _ = Describe("Plugin", func() {
	//FIt("unmarshal file descriptor", func() {
	//	b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
	//	Expect(err).NotTo(HaveOccurred())
	//	descriptor, err := convertProto(b)
	//	Expect(err).NotTo(HaveOccurred())
	//	log.Printf("%v", FuncsForProto("Bookstore", descriptor))
	//})
	//It("unmarshal file descriptor", func() {
	//	b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
	//	Expect(err).NotTo(HaveOccurred())
	//	descriptor, err := convertProto(b)
	//	Expect(err).NotTo(HaveOccurred())
	//	//addHttpRulesToProto("my-upstream", "Bookstore", descriptor)
	//	log.Printf("%v", FuncsForProto("Bookstore", descriptor))
	//})
	It("returns a dependency with a file ref for each descriptor file "+
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
	Describe("ProcessUpstream", func() {
		It("Marks the cluster metadata with the transformation filter", func() {
			in := &v1.Upstream{
				Name: "myupstream",
				ServiceInfo: &v1.ServiceInfo{
					Type: ServiceTypeGRPC,
					Properties: EncodeServiceProperties(ServiceProperties{
						DescriptorsFileRef: "file_1",
						GRPCServiceName:    "Bookstore",
					}),
				},
			}
			p := NewPlugin()
			b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
			Expect(err).NotTo(HaveOccurred())
			params := &plugin.UpstreamPluginParams{
				Files: map[string][]byte{"file_1": b},
			}
			out := &envoyapi.Cluster{}
			err = p.ProcessUpstream(params, in, out)
			Expect(err).To(BeNil())
			Expect(out.Metadata).NotTo(BeNil())
			Expect(out.Metadata.FilterMetadata).NotTo(BeNil())
			Expect(out.Metadata.FilterMetadata).To(HaveKey("io.solo.transformation"))
		})
		It("Stores the descriptors proto in the plugin memory and adds to it http rules", func() {
			in := &v1.Upstream{
				Name: "myupstream",
				ServiceInfo: &v1.ServiceInfo{
					Type: ServiceTypeGRPC,
					Properties: EncodeServiceProperties(ServiceProperties{
						DescriptorsFileRef: "file_1",
						GRPCServiceName:    "Bookstore",
					}),
				},
			}
			p := NewPlugin()
			b, err := ioutil.ReadFile("grpc-test-service/descriptors/proto.pb")
			Expect(err).NotTo(HaveOccurred())
			params := &plugin.UpstreamPluginParams{
				Files: map[string][]byte{"file_1": b},
			}
			out := &envoyapi.Cluster{}
			err = p.ProcessUpstream(params, in, out)
			Expect(err).To(BeNil())
			Expect(p.upstreamServices["myupstream"]).To(Equal("Bookstore"))
			Expect(p.serviceDescriptors["Bookstore"]).NotTo(BeNil())
			//for _, file := range p.serviceDescriptors["Bookstore"].File {
			//log.Printf("%v", p.serviceDescriptors["Bookstore"])
			//}
		})
	})
})
