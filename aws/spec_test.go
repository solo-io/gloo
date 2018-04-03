package aws_test

import (
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	. "github.com/solo-io/gloo-plugins/aws"
)

var _ = Describe("Spec", func() {
	Describe("Decode Function", func() {
		DescribeTable("with invalid spec",
			func(f v1.FunctionSpec) {
				_, err := DecodeFunctionSpec(f)
				Expect(err).To(HaveOccurred())
			},
			Entry("empty", &types.Struct{}),
			Entry("empty fields", &types.Struct{Fields: map[string]*types.Value{}}),
			Entry("missing function name", funcSpec("", "my-qualifier")))

		DescribeTable("with valid spec",
			func(name, qualifier string) {
				f, err := DecodeFunctionSpec(funcSpec(name, qualifier))
				Expect(err).NotTo(HaveOccurred())
				Expect(f.FunctionName).To(Equal(name))
				Expect(f.Qualifier).To(Equal(qualifier))
			},
			Entry("empty qualifier", "func1", ""),
			Entry("name and qualifier", "func1", "v1"))
	})

	Describe("Decode upstream", func() {
		DescribeTable("with invalid spec ",
			func(u v1.UpstreamSpec) {
				_, err := DecodeUpstreamSpec(u)
				Expect(err).To(HaveOccurred())
			},
			Entry("empty spec", &types.Struct{Fields: map[string]*types.Value{}}),
			Entry("invalid region", upstreamSpec("solo", "gloo")),
			Entry("missing region", upstreamSpec("", "x23")),
			Entry("missing secret ref", upstreamSpec("us-east-1", "")))

		Context("with valid spec", func() {
			var (
				u   *UpstreamSpec
				err error
			)
			BeforeEach(func() {
				u, err = DecodeUpstreamSpec(upstreamSpec("us-east-1", "aws-secret"))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should have correct parameters", func() {
				Expect(u.Region).To(Equal("us-east-1"))
				Expect(u.SecretRef).To(Equal("aws-secret"))
			})

			It("should return a hostname with region", func() {
				Expect(u.GetLambdaHostname()).To(ContainSubstring("us-east-1"))
			})
		})
	})
})

func funcSpec(name, qualifier string) v1.FunctionSpec {
	return &types.Struct{
		Fields: map[string]*types.Value{
			"function_name": &types.Value{
				Kind: &types.Value_StringValue{StringValue: name},
			},
			"qualifier": &types.Value{
				Kind: &types.Value_StringValue{StringValue: qualifier},
			},
		},
	}
}

func upstreamSpec(region, secretRef string) v1.UpstreamSpec {
	return &types.Struct{
		Fields: map[string]*types.Value{
			"region": &types.Value{
				Kind: &types.Value_StringValue{StringValue: region},
			},
			"secret_ref": &types.Value{
				Kind: &types.Value_StringValue{StringValue: secretRef},
			},
		},
	}
}
