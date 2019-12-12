package truncate_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/grpc"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/truncate"
)

var _ = Describe("UpstreamTruncator", func() {

	var subject = truncate.Truncator{}

	Describe("Truncate", func() {
		It("strips the Descriptors field on gRPC Kube upstreams", func() {
			inputServiceSpec := &options.ServiceSpec{
				PluginType: &options.ServiceSpec_Grpc{
					Grpc: &grpc.ServiceSpec{
						Descriptors: []byte("non empty string that we want to clear out"),
						GrpcServices: []*grpc.ServiceSpec_GrpcService{{
							PackageName: "non-empty",
						}},
					},
				},
			}

			expectedServiceSpec := &options.ServiceSpec{
				PluginType: &options.ServiceSpec_Grpc{
					Grpc: &grpc.ServiceSpec{
						GrpcServices: []*grpc.ServiceSpec_GrpcService{{
							PackageName: "non-empty",
						}},
					},
				},
			}

			// Define a base upstream from which we will derive input and expected output
			baseUpstream := &gloov1.Upstream{
				Metadata: core.Metadata{Namespace: "foo", Name: "bar"},
				UpstreamType: &gloov1.Upstream_Kube{
					Kube: &kubernetes.UpstreamSpec{
						ServiceName:      "my-service",
						ServiceNamespace: "my-namespace",
						ServicePort:      99,
						ServiceSpec: &options.ServiceSpec{
							PluginType: &options.ServiceSpec_Grpc{
								Grpc: &grpc.ServiceSpec{
									GrpcServices: []*grpc.ServiceSpec_GrpcService{{
										PackageName: "non-empty",
									}},
								},
							},
						},
					},
				},
			}

			// Input is a copy of the base upstream with non-empty descriptors
			input := &gloov1.Upstream{}
			baseUpstream.DeepCopyInto(input)
			input.GetKube().ServiceSpec = inputServiceSpec

			// Expected is a copy of the base upstream which should have nil descriptors
			expected := &gloov1.Upstream{}
			baseUpstream.DeepCopyInto(expected)
			expected.GetKube().ServiceSpec = expectedServiceSpec

			Expect(input).NotTo(Equal(expected), "sanity check that input and output are different initially")

			// Truncate updates the input in place
			subject.Truncate(input)

			ExpectEqualProtoMessages(input, expected, "input should have been modified to match expected")
			Expect(input.GetKube().GetServiceSpec().GetGrpc().GetDescriptors()).To(BeNil())
		})
	})
})
