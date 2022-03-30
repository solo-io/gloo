package grpc_test

import (
	"fmt"

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

		It("should translate a simple proto to graphql schema", func() {
			schema, err := functionDiscovery.BuildGraphQLApiFromGrpcReflection(refClient)
			fmt.Printf("%s\n", schema.GetExecutableSchema().GetSchemaDefinition())
			Expect(err).NotTo(HaveOccurred())

		})
	})
})

func NewMockGrpcReflectionClient() *MockGrpcReflectionClient {
	p := protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"foo/bar.proto": `
syntax = "proto3";
package foo;
message Pet {
	
	uint64 id = 1;
	oneof name {
		string full_name = 2;
		string first_name = 3;
	}
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
