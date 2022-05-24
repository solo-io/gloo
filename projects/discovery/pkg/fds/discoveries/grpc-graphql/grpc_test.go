package grpc_test

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/grpc-graphql"
)

var _ = Describe("Grpc reflection - graphql schema discovery test", func() {
	Context("translates proto to graphql schema", func() {
		var (
			refClient         grpc.GrpcReflectionClient
			functionDiscovery *grpc.GraphqlSchemaDiscovery
		)
		BeforeEach(func() {
			refClient = NewMockGrpcReflectionClient()
			upstream := &v1.Upstream{
				Metadata: &core.Metadata{
					Name:      "test",
					Namespace: "test-system",
				},
			}
			functionDiscovery = grpc.NewFunctionDiscoveryFactory().NewFunctionDiscovery(upstream, fds.AdditionalClients{}).(*grpc.GraphqlSchemaDiscovery)
		})

		var (
			schemaDef string
		)

		BeforeEach(func() {
			schema, err := functionDiscovery.BuildGraphQLApiFromGrpcReflection(refClient)
			Expect(err).NotTo(HaveOccurred())

			schemaDef = schema.GetExecutableSchema().GetSchemaDefinition()
		})

		AfterEach(func() {
			schemaDef = ""
		})

		It("should translate a simple proto to graphql schema", func() {

			// translates empty types
			Expect(schemaDef).To(ContainSubstring(`
"""Created from protobuf type google.protobuf.Empty"""
type google_protobuf_Empty {

  """This GraphQL type was generated from an empty proto message. This empty field exists to keep the schema GraphQL spec compliant. If queried, this field will always return false."""
  _: Boolean
}`))
		})

		It("should translate a basic Pet type", func() {

			// translates empty types
			Expect(schemaDef).To(ContainSubstring(`"""Created from protobuf type foo.Pet"""
type foo_Pet {
  id: Int
  full_name: String
  first_name: String
  empty: google_protobuf_Empty
}`))
		})
		It("should translate an input type", func() {
			Expect(schemaDef).To(ContainSubstring(`"""Created from protobuf type foo.PetRequest"""
input foo_PetRequestInput {
  id: Int
}
`))
		})

		It("should translate the Query types", func() {
			Expect(schemaDef).To(ContainSubstring(`
type Query {
  GetBar(PetRequest: PetRequestInput): Pet @resolve(name: "Query|Foo.GetBar")
}`))
		})
	})

})

func NewMockGrpcReflectionClient() *MockGrpcReflectionClient {
	p := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"foo/bar.proto": `
syntax = "proto3";
package foo;

import "google/protobuf/empty.proto";

message Pet {
	
	uint64 id = 1;
	oneof name {
		string full_name = 2;
		string first_name = 3;
	}
	// This is an empty type, which is illegal in GraphQL Schemas. 
	// In the GraphQL translation, we add an empty boolean field to make the GraphQL Schema spec-compliant.
	google.protobuf.Empty empty = 4;
}
message PetRequest {
	uint64 id = 1;
}
service Foo {
	rpc GetBar(PetRequest) returns (Pet) {};
}
				
				`}),
	}
	fds, err := p.ParseFiles("foo/bar.proto")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return &MockGrpcReflectionClient{
		fd: fds,
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
