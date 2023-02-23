package test

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/grpc-graphql"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Grpc reflection - graphql schema discovery test", func() {
	Context("translates proto to graphql schema", func() {
		var (
			opts bootstrap.Opts
		)
		BeforeEach(func() {
			opts = bootstrap.Opts{
				Settings: &v1.Settings{
					Discovery: &v1.Settings_DiscoveryOptions{
						FdsOptions: &v1.Settings_DiscoveryOptions_FdsOptions{
							GraphqlEnabled: &wrapperspb.BoolValue{Value: true},
						},
					},
				},
			}
		})
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
			functionDiscovery = grpc.NewFunctionDiscoveryFactory(opts).NewFunctionDiscovery(upstream, fds.AdditionalClients{}).(*grpc.GraphqlSchemaDiscovery)
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

		DescribeTable("translates basic gRPC types",
			func(grpcType, expectedGraphqlType string, expectedGraphqlInputType ...string) {
				expectedInputType := expectedGraphqlType
				if len(expectedGraphqlInputType) > 0 {
					expectedInputType = expectedGraphqlInputType[0]
				}
				sampleGrpcProto := fmt.Sprintf(`
syntax = "proto3";

import "google/protobuf/wrappers.proto";
import "google/protobuf/duration.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/field_mask.proto";


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
				fmt.Println(printer.Print(doc))
				Expect(getInputFieldType(doc)).To(Equal(grpc.GetNonNullType(expectedInputType)))
				Expect(getFieldType(doc)).To(Equal(expectedGraphqlType))
			},
			Entry("double", "double", grpc.GRAPHQL_NON_NULL_FLOAT),
			Entry("float", "float", grpc.GRAPHQL_NON_NULL_FLOAT),
			Entry("int32", "int32", grpc.GRAPHQL_NON_NULL_INT),
			Entry("int64", "int64", grpc.GRAPHQL_NON_NULL_INT),
			Entry("uint32", "uint32", grpc.GRAPHQL_NON_NULL_INT),
			Entry("uint64", "uint64", grpc.GRAPHQL_NON_NULL_INT),
			Entry("sint32", "sint32", grpc.GRAPHQL_NON_NULL_INT),
			Entry("sint64", "sint64", grpc.GRAPHQL_NON_NULL_INT),
			Entry("fixed32", "fixed32", grpc.GRAPHQL_NON_NULL_INT),
			Entry("fixed64", "fixed64", grpc.GRAPHQL_NON_NULL_INT),
			Entry("sfixed32", "sfixed32", grpc.GRAPHQL_NON_NULL_INT),
			Entry("sfixed64", "sfixed64", grpc.GRAPHQL_NON_NULL_INT),
			Entry("bool", "bool", grpc.GRAPHQL_NON_NULL_BOOLEAN),
			Entry("string", "string", grpc.GRAPHQL_NON_NULL_STRING),
			Entry("repeated string", "repeated string", "[String!]", "[String]"),
			Entry("repeated double", "repeated double", "[Float!]", "[Float]"),
			Entry("repeated bool", "repeated bool", "[Boolean!]", "[Boolean]"),
			Entry("repeated int32", "repeated int32", "[Int!]", "[Int]"),
			// google protobuf wrapper types
			Entry("string wrapper", "google.protobuf.StringValue", grpc.GRAPHQL_STRING, "StringValueInput"),
			Entry("bytes wrapper", "google.protobuf.BytesValue", grpc.GRAPHQL_STRING, "BytesValueInput"),
			Entry("double wrapper", "google.protobuf.DoubleValue", grpc.GRAPHQL_FLOAT, "DoubleValueInput"),
			Entry("bool wrapper", "google.protobuf.BoolValue", grpc.GRAPHQL_BOOLEAN, "BoolValueInput"),
			Entry("int32 wrapper", "google.protobuf.Int32Value", grpc.GRAPHQL_INT, "Int32ValueInput"),
			Entry("uint32 wrapper", "google.protobuf.UInt32Value", grpc.GRAPHQL_INT, "UInt32ValueInput"),
			Entry("int64 wrapper", "google.protobuf.Int64Value", grpc.GRAPHQL_STRING, "Int64ValueInput"),
			Entry("uint64 wrapper", "google.protobuf.UInt64Value", grpc.GRAPHQL_STRING, "UInt64ValueInput"),
			Entry("float wrapper", "google.protobuf.FloatValue", grpc.GRAPHQL_FLOAT, "FloatValueInput"),
			Entry("repeated int wrapper", "repeated google.protobuf.Int32Value", "[Int]", "[Int32ValueInput]"),
			// misc google types
			Entry("duration", "google.protobuf.Duration", grpc.GRAPHQL_STRING, "DurationInput"),
			Entry("field mask", "google.protobuf.FieldMask", grpc.GRAPHQL_STRING, "FieldMaskInput"),
			Entry("time stamp", "google.protobuf.Timestamp", grpc.GRAPHQL_STRING, "TimestampInput"),
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
				It("creates both PetRequestInput and PetInput graphql objects", func() {
					// input types
					fmt.Println(printer.Print(doc))
					petRequestInputType := getTypeDefinition(doc, true, "PetRequestInput")
					Expect(petRequestInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petRequestInputType, "pet", "PetInput")).To(BeTrue())
					petInputType := getTypeDefinition(doc, true, "PetInput")
					Expect(petInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petInputType, "id", grpc.GRAPHQL_INT)).To(BeTrue())

					// non input types
					petType := getTypeDefinition(doc, false, "Pet")
					Expect(petType).NotTo(BeNil())
					Expect(getFieldsWithType(petType, "id", grpc.GRAPHQL_NON_NULL_INT)).To(BeTrue())
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

				It("creates PetRequestInput and nested PetInput graphql objects", func() {
					fmt.Println(printer.Print(doc))
					petRequestInputType := getTypeDefinition(doc, true, "PetRequestInput")
					Expect(petRequestInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petRequestInputType, "pet", "PetRequest_PetInput")).To(BeTrue())
					petInputType := getTypeDefinition(doc, true, "PetRequest_PetInput")
					Expect(petInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petInputType, "input_id", grpc.GRAPHQL_FLOAT)).To(BeTrue())
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
				It("creates PetRequestInput, PetInput, and nested enum Status graphql objects", func() {
					petStatusEnum := getEnumDefinition(doc, "Status")
					Expect(petStatusEnum).NotTo(BeNil())
					Expect(getEnumValues(petStatusEnum)).To(ContainElements("IN_STORE", "IN_TRANSIT"))

					petStatusEnum2 := getEnumDefinition(doc, "Pet_Status")
					Expect(petStatusEnum2).NotTo(BeNil())
					Expect(getEnumValues(petStatusEnum2)).To(ContainElements("AVAILABLE", "SOLD"))

				})
			})

			Context("translates oneof", func() {
				BeforeEach(func() {
					proto =
						`syntax = "proto3";
package foo;
message Pet {
	int32 id = 1;
	string name = 2;
}

message PetRequest {
	oneof identifier {
		string name = 1;
		int32 id = 2;
		bool active = 3;
		float f = 4;
	  }	
}

service Foo {
	rpc GetPet(PetRequest) returns (Pet) {};
}`
				})

				It("creates PetRequestInput with oneof and PetInput graphql objects", func() {
					fmt.Println(printer.Print(doc))
					petRequestInputType := getTypeDefinition(doc, true, "PetRequestInput")
					Expect(petRequestInputType).NotTo(BeNil())
					Expect(getFieldsWithType(petRequestInputType, "name", grpc.GRAPHQL_STRING)).To(BeTrue())
					Expect(getFieldsWithType(petRequestInputType, "id", grpc.GRAPHQL_INT)).To(BeTrue())
					Expect(getFieldsWithType(petRequestInputType, "active", grpc.GRAPHQL_BOOLEAN)).To(BeTrue())
					Expect(getFieldsWithType(petRequestInputType, "f", grpc.GRAPHQL_FLOAT)).To(BeTrue())
				})
			})

			Context("translates input message that is also output message", func() {
				// in this case the input message should be different than the output message. So that the values of the input message
				// will be non-null all the time.
				BeforeEach(func() {
					proto =
						`syntax = "proto3";
package foo;
message Pet {
	int32 id = 1;
	string name = 2;
}

service Foo {
	rpc GetPet(Pet) returns (Pet) {};
}`
				})

				It("creates PetInput and Pet graphql objects with only one message type", func() {
					fmt.Println(printer.Print(doc))
					petInput := getTypeDefinition(doc, true, "PetInput")
					Expect(petInput).NotTo(BeNil())
					Expect(getFieldsWithType(petInput, "id", grpc.GRAPHQL_INT)).To(BeTrue())
					Expect(getFieldsWithType(petInput, "name", grpc.GRAPHQL_STRING)).To(BeTrue())
					petOuput := getTypeDefinition(doc, false, "Pet")
					Expect(petOuput).NotTo(BeNil())
					Expect(getFieldsWithType(petOuput, "id", grpc.GRAPHQL_NON_NULL_INT)).To(BeTrue())
					Expect(getFieldsWithType(petOuput, "name", grpc.GRAPHQL_NON_NULL_STRING)).To(BeTrue())
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
