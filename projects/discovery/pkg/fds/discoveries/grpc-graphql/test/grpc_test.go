package test

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/grpc-graphql"
	"google.golang.org/protobuf/types/descriptorpb"
)

var _ = Describe("Grpc reflection - graphql schema discovery test", func() {
	Context("translates proto to graphql schema", func() {
		var (
			refClient         grpc.GrpcReflectionClient
			functionDiscovery *grpc.GraphqlSchemaDiscovery
		)

		initialize := func(testProto string) {
			refClient = NewMockGrpcReflectionClient(testProto)
			upstream := &v1.Upstream{
				Metadata: &core.Metadata{
					Name:      "test",
					Namespace: "test-system",
				},
			}
			functionDiscovery = grpc.NewFunctionDiscoveryFactory().NewFunctionDiscovery(upstream, fds.AdditionalClients{}).(*grpc.GraphqlSchemaDiscovery)
		}

		getGraphqlDocForProto := func(offset int, proto string) *ast.Document {
			initialize(proto)
			desc := &descriptorpb.FileDescriptorSet{}
			// swallow error as this shouldn't occur
			services, _ := refClient.ListServices()
			schemaBuilder, _, err := functionDiscovery.GetSchemaBuilderForProtoFileDescriptor(refClient, desc, services)
			ExpectWithOffset(offset+1, err).NotTo(HaveOccurred())
			doc := schemaBuilder.Build()
			return doc
		}

		DescribeTable("translates basic gRPC types", func(grpcType, expectedGraphqlType string) {
			sampleGrpcProto := fmt.Sprintf(`
syntax = "proto3";
package foo;
message Pet {
	%s testType = 1;
}

message PetRequest {
	%s testType = 1;
}
service Foo {
	rpc GetPet(PetRequest) returns (Pet) {};
}
`, grpcType, grpcType)
			doc := getGraphqlDocForProto(1, sampleGrpcProto)
			ExpectWithOffset(1, getInputFieldType(doc)).To(Equal(expectedGraphqlType))
			ExpectWithOffset(1, getFieldType(doc)).To(Equal(expectedGraphqlType))
			fmt.Println(printer.Print(doc))
		},
			Entry("double", "double", "Float"),
			Entry("float", "float", "Float"),
			Entry("int32", "int32", "Int"),
			Entry("int64", "int64", "Int"),
			Entry("uint32", "uint32", "Int"),
			Entry("uint64", "uint64", "Int"),
			Entry("sint32", "sint32", "Int"),
			Entry("sint64", "sint64", "Int"),
			Entry("fixed32", "fixed32", "Int"),
			Entry("fixed64", "fixed64", "Int"),
			Entry("sfixed32", "sfixed32", "Int"),
			Entry("sfixed64", "sfixed64", "Int"),
			Entry("bool", "bool", "Boolean"),
			Entry("string", "string", "String"),
			Entry("repeated string", "repeated string", "[String]"),
			Entry("repeated double", "repeated double", "[Float]"),
			Entry("repeated bool", "repeated bool", "[Boolean]"),
			Entry("repeated int32", "repeated int32", "[Int]"),
		)

		Context("translates proto messages", func() {
			var (
				doc   *ast.Document
				proto string
			)

			JustBeforeEach(func() {
				doc = getGraphqlDocForProto(0, proto)
			})
			Context("translates Messages to Input types and regular types", func() {
				BeforeEach(func() {
					proto =
						`syntax = "proto3";
						package foo;
						message Pet {
							int32 id = 1;
						}
						
						message PetRequest {
							Pet pet = 1;
						}
						
						service Foo {
							rpc GetPet(PetRequest) returns (Pet) {};
						}`
				})
				It("creates PetRequestInput and PetInput graphql objects", func() {
					// input types
					fmt.Println(printer.Print(doc))
					petRequestInputType := getTypeDefinition(doc, true, "PetRequestInput")
					Expect(petRequestInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petRequestInputType, "pet", "PetInput")).To(BeTrue())
					petInputType := getTypeDefinition(doc, true, "PetInput")
					Expect(petInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petInputType, "id", "Int")).To(BeTrue())

					// non input types
					petType := getTypeDefinition(doc, false, "Pet")
					Expect(petType).NotTo(BeNil())
					Expect(getFieldsWithType(petType, "id", "Int")).To(BeTrue())
				})
			})

			Context("translates nested Messages to Input types and regular types", func() {
				BeforeEach(func() {
					proto =
						`syntax = "proto3";
package foo;
message Pet {
	int32 id = 1;
}

message PetRequest {
	message Pet {
		float input_id = 1;
  }
	Pet pet = 1;
}

service Foo {
	rpc GetPet(PetRequest) returns (Pet) {};
}`
				})

				It("creates PetRequestInput and PetInput graphql objects", func() {
					fmt.Println(printer.Print(doc))
					petRequestInputType := getTypeDefinition(doc, true, "PetRequestInput")
					Expect(petRequestInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petRequestInputType, "pet", "PetRequest_PetInput")).To(BeTrue())
					petInputType := getTypeDefinition(doc, true, "PetRequest_PetInput")
					Expect(petInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petInputType, "input_id", "Float")).To(BeTrue())
				})
			})

			Context("translates enums and nested enums", func() {
				BeforeEach(func() {
					proto =
						`syntax = "proto3";
package foo;
message Pet {
	enum Status {
		AVAILABLE = 0;
		SOLD = 1;
	}
	int32 id = 1;
	Status status = 2;
}

enum Status {
	IN_STORE = 0;
	IN_TRANSIT = 1;
}

message PetRequest {
	Status pet_status = 1;
}

service Foo {
	rpc GetPet(PetRequest) returns (Pet) {};
}`
				})

				getEnumValues := func(definition *ast.EnumDefinition) []string {
					var ret []string
					for _, val := range definition.Values {
						ret = append(ret, val.Name.Value)
					}
					return ret
				}
				It("creates PetRequestInput and PetInput graphql objects", func() {
					petStatusEnum := getEnumDefinition(doc, "Status")
					Expect(petStatusEnum).NotTo(BeNil())
					Expect(getEnumValues(petStatusEnum)).To(ContainElements("IN_STORE", "IN_TRANSIT"))

					petStatusEnum2 := getEnumDefinition(doc, "Pet_Status")
					Expect(petStatusEnum2).NotTo(BeNil())
					Expect(getEnumValues(petStatusEnum2)).To(ContainElements("AVAILABLE", "SOLD"))

				})
			})
		})
	})
})

func NewMockGrpcReflectionClient(proto string) *MockGrpcReflectionClient {
	p := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"foo/bar.proto": proto}),
	}
	fileDescriptors, err := p.ParseFiles("foo/bar.proto")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return &MockGrpcReflectionClient{
		fd: fileDescriptors,
	}
}

type MockGrpcReflectionClient struct {
	fd []*desc.FileDescriptor
}

func (*MockGrpcReflectionClient) ListServices() ([]string, error) {
	return []string{"Foo"}, nil
}

func (m *MockGrpcReflectionClient) FileContainingSymbol(string) (*desc.FileDescriptor, error) {
	return m.fd[0], nil
}
