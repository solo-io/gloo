package grpc_test

import (
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/grpc-graphql"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Grpc reflection - graphql schema discovery test", func() {
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
	Context("Graphql FDS Disabled", func() {
		BeforeEach(func() {
			opts.Settings.GetDiscovery().GetFdsOptions().GraphqlEnabled.Value = false
		})

		It("should return a nil Discovery Factory when disabled in settings", func() {
			Expect(grpc.NewFunctionDiscoveryFactory(opts)).To(BeNil())
		})

	})
	Context("Graphql FDS Enabled", func() {
		It("should enable GraphQL FDS when no setting is present", func() {
			opts.Settings.GetDiscovery().FdsOptions = nil
			Expect(grpc.NewFunctionDiscoveryFactory(opts)).NotTo(BeNil())
		})

		Context("translates proto to graphql schema", func() {
			var (
				refClient         grpc.GrpcReflectionClient
				functionDiscovery *grpc.GraphqlSchemaDiscovery
				schemaDef         string
				localExecutor     *v1beta1.Executor_Local_
			)
			BeforeEach(func() {
				refClient = NewMockGrpcReflectionClient()
				upstream := &v1.Upstream{
					Metadata: &core.Metadata{
						Name:      "test",
						Namespace: "test-system",
					},
				}
				functionDiscovery = grpc.NewFunctionDiscoveryFactory(opts).NewFunctionDiscovery(upstream, fds.AdditionalClients{}).(*grpc.GraphqlSchemaDiscovery)
				schema, err := functionDiscovery.BuildGraphQLApiFromGrpcReflection(refClient)
				Expect(err).NotTo(HaveOccurred())

				schemaDef = schema.GetExecutableSchema().GetSchemaDefinition()
				Expect(schema.GetExecutableSchema().GetExecutor().Executor).To(BeAssignableToTypeOf(&v1beta1.Executor_Local_{}))
				localExecutor = schema.GetExecutableSchema().GetExecutor().Executor.(*v1beta1.Executor_Local_)
			})

			AfterEach(func() {
				schemaDef = ""
				localExecutor = nil
			})

			It("should translate a simple proto to graphql schema", func() {

				// translates empty types
				Expect(schemaDef).To(ContainSubstring(`
"""Created from protobuf type google.protobuf.Empty"""
type Empty {

  """This GraphQL type was generated from an empty proto message. This empty field exists to keep the schema GraphQL spec compliant. If queried, this field will always return false."""
  _: Boolean
}`))
			})

			It("should translate a basic Pet type", func() {

				// translates empty types
				Expect(schemaDef).To(ContainSubstring(`"""Created from protobuf type foo.Pet"""
type Pet {
  id: Int
  full_name: String
  first_name: String
  empty: Empty
}`))
			})
			It("should translate an input type", func() {
				Expect(schemaDef).To(ContainSubstring(`"""Created from protobuf type foo.PetRequest"""
input PetRequestInput {
  id: Int
  recursive: PetRequest_RecursiveMessageInput
}`))
				Expect(schemaDef).To(ContainSubstring(`"""Created from protobuf type foo.PetRequest.RecursiveMessage"""
input PetRequest_RecursiveMessageInput {
  request: PetRequestInput
}`))
			})

			It("should translate the Query types", func() {
				Expect(schemaDef).To(ContainSubstring(`
type Query {
  GetBar(PetRequest: PetRequestInput): Pet @resolve(name: "Query|Foo.GetBar")
}`))
			})

			It("should construct outgoing message", func() {
				resolution := &v1beta1.Resolution{
					Resolver: &v1beta1.Resolution_GrpcResolver{
						GrpcResolver: &v1beta1.GrpcResolver{
							UpstreamRef: &core.ResourceRef{Name: "test", Namespace: "test-system"},
							RequestTransform: &v1beta1.GrpcRequestTemplate{
								OutgoingMessageJson: &structpb.Value{
									Kind: &structpb.Value_StructValue{
										StructValue: &structpb.Struct{
											Fields: map[string]*structpb.Value{
												"id": {Kind: &structpb.Value_StringValue{StringValue: "{$args.PetRequest.id}"}},
												// shouldn't recurse into the recursive message
												"recursive": {Kind: &structpb.Value_StringValue{StringValue: "{$args.PetRequest.recursive}"}},
											},
										},
									},
								},
								ServiceName: "foo.Foo",
								MethodName:  "GetBar",
							},
						},
					},
				}
				Expect(localExecutor.Local.GetResolutions()).To(HaveKey("Query|Foo.GetBar"))
				Expect(localExecutor.Local.GetResolutions()["Query|Foo.GetBar"]).To(matchers.MatchProto(resolution))
			})
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
	RecursiveMessage recursive = 2;
	message RecursiveMessage {
		PetRequest request = 1;
	}
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
