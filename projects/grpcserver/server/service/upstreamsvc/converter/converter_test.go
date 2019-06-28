package converter_test

import (
	. "github.com/onsi/ginkgo"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	. "github.com/solo-io/go-utils/testutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/converter"
)

var _ = Describe("ConverterTest", func() {
	Describe("ConvertInputToUpstreamSpec", func() {
		It("works", func() {
			upstreamInputConverter := converter.NewUpstreamInputConverter()

			testCases := []struct {
				in  *v1.UpstreamInput
				out *gloov1.UpstreamSpec
			}{
				{
					in: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
					out: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
				},
				{
					in: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Azure{
							Azure: &azure.UpstreamSpec{FunctionAppName: "test"},
						},
					},
					out: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Azure{
							Azure: &azure.UpstreamSpec{FunctionAppName: "test"},
						},
					},
				},
				{
					in: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Consul{
							Consul: &consul.UpstreamSpec{ServiceName: "test"},
						},
					},
					out: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Consul{
							Consul: &consul.UpstreamSpec{ServiceName: "test"},
						},
					},
				},
				{
					in: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Kube{
							Kube: &kubernetes.UpstreamSpec{ServiceName: "test"},
						},
					},
					out: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Kube{
							Kube: &kubernetes.UpstreamSpec{ServiceName: "test"},
						},
					},
				},
				{
					in: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Static{
							Static: &static.UpstreamSpec{UseTls: true},
						},
					},
					out: &gloov1.UpstreamSpec{
						UpstreamType: &gloov1.UpstreamSpec_Static{
							Static: &static.UpstreamSpec{UseTls: true},
						},
					},
				},
				{
					in:  &v1.UpstreamInput{},
					out: nil,
				},
			}

			for _, tc := range testCases {
				actual := upstreamInputConverter.ConvertInputToUpstreamSpec(tc.in)
				ExpectEqualProtoMessages(actual, tc.out)
			}
		})
	})
})
