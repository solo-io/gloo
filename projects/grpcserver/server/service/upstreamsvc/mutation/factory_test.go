package mutation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws/ec2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/azure"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	. "github.com/solo-io/go-utils/testutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation"
)

var _ = Describe("FactoryTest", func() {
	Describe("ConfigureUpstream", func() {
		It("works", func() {
			testCases := []struct {
				desc          string
				input         *v1.UpstreamInput
				before, after *gloov1.Upstream
				err           error
			}{
				{
					desc: "aws",
					input: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
					before: &gloov1.Upstream{},
					after: &gloov1.Upstream{
						UpstreamType: &gloov1.Upstream_Aws{
							Aws: &aws.UpstreamSpec{Region: "test"},
						},
					},
				},
				{
					desc: "azure",
					input: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Azure{
							Azure: &azure.UpstreamSpec{FunctionAppName: "test"},
						},
					},
					before: &gloov1.Upstream{},
					after: &gloov1.Upstream{
						UpstreamType: &gloov1.Upstream_Azure{
							Azure: &azure.UpstreamSpec{FunctionAppName: "test"},
						},
					},
				},
				{
					desc: "consul",
					input: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Consul{
							Consul: &consul.UpstreamSpec{ServiceName: "test"},
						},
					},
					before: &gloov1.Upstream{},
					after: &gloov1.Upstream{
						UpstreamType: &gloov1.Upstream_Consul{
							Consul: &consul.UpstreamSpec{ServiceName: "test"},
						},
					},
				},
				{
					desc: "kube",
					input: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Kube{
							Kube: &kubernetes.UpstreamSpec{ServiceName: "test"},
						},
					},
					before: &gloov1.Upstream{},
					after: &gloov1.Upstream{
						UpstreamType: &gloov1.Upstream_Kube{
							Kube: &kubernetes.UpstreamSpec{ServiceName: "test"},
						},
					},
				},
				{
					desc: "static",
					input: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_Static{
							Static: &static.UpstreamSpec{UseTls: true},
						},
					},
					before: &gloov1.Upstream{},
					after: &gloov1.Upstream{
						UpstreamType: &gloov1.Upstream_Static{
							Static: &static.UpstreamSpec{UseTls: true},
						},
					},
				},
				{
					desc: "ec2",
					input: &v1.UpstreamInput{
						Spec: &v1.UpstreamInput_AwsEc2{
							AwsEc2: &ec2.UpstreamSpec{Region: "test"},
						},
					},
					before: &gloov1.Upstream{},
					after: &gloov1.Upstream{
						UpstreamType: &gloov1.Upstream_AwsEc2{
							AwsEc2: &ec2.UpstreamSpec{Region: "test"},
						},
					},
				},
				{
					desc:   "noop: empty input",
					input:  &v1.UpstreamInput{},
					before: &gloov1.Upstream{},
					after:  &gloov1.Upstream{},
					err:    mutation.EmptyInputError,
				},
				{
					desc:   "noop: nil input",
					input:  nil,
					before: &gloov1.Upstream{},
					after:  &gloov1.Upstream{},
					err:    mutation.EmptyInputError,
				},
			}

			for _, tc := range testCases {
				err := mutation.NewFactory().ConfigureUpstream(tc.input)(tc.before)
				if tc.err != nil {
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(tc.err))
				} else {
					Expect(err).NotTo(HaveOccurred())
				}
				ExpectEqualProtoMessages(tc.before, tc.after, tc.desc)
			}
		})
	})
})
