package protoutil_test

import (
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/glue/pkg/protoutil"
)

type testType struct {
	A string
	B b
}

type b struct {
	C string
	D string
}

var _ = Describe("Protoutil Funcs", func() {
	Describe("MarshalStruct", func() {
		It("returns a pb struct for the given type", func() {
			t := testType{
				A: "a",
				B: b{
					C: "c",
					D: "d",
				},
			}
			pb, err := MarshalStruct(t)
			Expect(err).NotTo(HaveOccurred())
			expectedPB := &types.Struct{
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
			}
			Expect(pb).To(Equal(expectedPB))
		})
	})
})
