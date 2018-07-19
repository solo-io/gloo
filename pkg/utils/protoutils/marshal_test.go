package protoutils_test

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-kit/pkg/utils/protoutils"
)

type testType struct {
	A string
	B b
}

type b struct {
	C string
	D string
}

var tests = []struct {
	in       interface{}
	expected proto.Message
}{
	{
		in: testType{
			A: "a",
			B: b{
				C: "c",
				D: "d",
			},
		},
		expected: &types.Struct{
			Fields: map[string]*types.Value{
				"A": {
					Kind: &types.Value_StringValue{StringValue: "a"},
				},
				"B": {
					Kind: &types.Value_StructValue{
						StructValue: &types.Struct{
							Fields: map[string]*types.Value{
								"C": {
									Kind: &types.Value_StringValue{StringValue: "c"},
								},
								"D": {
									Kind: &types.Value_StringValue{StringValue: "d"},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		in: map[string]interface{}{
			"a": "b",
			"c": "d",
		},
		expected: &types.Struct{
			Fields: map[string]*types.Value{
				"a": {
					Kind: &types.Value_StringValue{StringValue: "b"},
				},
				"c": {
					Kind: &types.Value_StringValue{StringValue: "d"},
				},
			},
		},
	},
}

var _ = Describe("Protoutil Funcs", func() {
	Describe("MarshalStruct", func() {
		for _, test := range tests {
			It("returns a pb struct for object of the given type", func() {
				pb, err := MarshalStruct(test.in)
				Expect(err).NotTo(HaveOccurred())
				Expect(pb).To(Equal(test.expected))
			})
		}
	})
})
