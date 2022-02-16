package graphql_handler_test

import (
	"context"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1"
	mock_graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1/mocks"
	. "github.com/solo-io/solo-kit/test/matchers"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/graphql_handler"
	fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	fed_v1_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("fed graphql handler", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context

		mockGlooInstanceClient   *mock_fed_v1.MockGlooInstanceClient
		mockMCClientset          *mock_graphql_v1alpha1.MockMulticlusterClientset
		mockGraphqlClientset1    *mock_graphql_v1alpha1.MockClientset
		mockGraphqlClientset2    *mock_graphql_v1alpha1.MockClientset
		mockGraphqlSchemaClient1 *mock_graphql_v1alpha1.MockGraphQLSchemaClient
		mockGraphqlSchemaClient2 *mock_graphql_v1alpha1.MockGraphQLSchemaClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockGlooInstanceClient = mock_fed_v1.NewMockGlooInstanceClient(ctrl)
		mockMCClientset = mock_graphql_v1alpha1.NewMockMulticlusterClientset(ctrl)
		mockGraphqlClientset1 = mock_graphql_v1alpha1.NewMockClientset(ctrl)
		mockGraphqlClientset2 = mock_graphql_v1alpha1.NewMockClientset(ctrl)
		mockGraphqlSchemaClient1 = mock_graphql_v1alpha1.NewMockGraphQLSchemaClient(ctrl)
		mockGraphqlSchemaClient2 = mock_graphql_v1alpha1.NewMockGraphQLSchemaClient(ctrl)

		local := fed_v1.GlooInstance{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "local-test",
				Namespace: "gloo-system",
			},
			Spec: fed_v1_types.GlooInstanceSpec{
				Cluster: "local-cluster",
			},
			Status: fed_v1_types.GlooInstanceStatus{},
		}
		remote := fed_v1.GlooInstance{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "remote-test",
				Namespace: "gloo-system",
			},
			Spec: fed_v1_types.GlooInstanceSpec{
				Cluster: "remote-cluster",
			},
			Status: fed_v1_types.GlooInstanceStatus{},
		}
		mockGlooInstanceClient.EXPECT().ListGlooInstance(ctx).Return(&fed_v1.GlooInstanceList{
			Items: []fed_v1.GlooInstance{
				local,
				remote,
			},
		}, nil).AnyTimes()
		mockMCClientset.EXPECT().Cluster("local-cluster").Return(mockGraphqlClientset1, nil).AnyTimes()
		mockMCClientset.EXPECT().Cluster("remote-cluster").Return(mockGraphqlClientset2, nil).AnyTimes()
		mockGraphqlClientset1.EXPECT().GraphQLSchemas().Return(mockGraphqlSchemaClient1).AnyTimes()
		mockGraphqlClientset2.EXPECT().GraphQLSchemas().Return(mockGraphqlSchemaClient2).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetGraphqlSchema", func() {
		It("can get graphql schema by ref", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient1.EXPECT().GetGraphQLSchema(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlSchema, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			resp, err := handler.GetGraphqlSchema(ctx, &rpc_edge_v1.GetGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.GetGraphqlSchemaResponse{
				GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &petstoreGraphqlSchema.Spec,
					Status:       &petstoreGraphqlSchema.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
				},
			}))
		})
	})

	Context("GetGraphqlSchemaYaml", func() {
		It("can get graphql schema yaml", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient1.EXPECT().GetGraphQLSchema(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlSchema, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			resp, err := handler.GetGraphqlSchemaYaml(ctx, &rpc_edge_v1.GetGraphqlSchemaYamlRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.YamlData.Yaml).To(MatchYAML(petstoreYaml))
		})
	})

	Context("ListGraphqlSchemas", func() {
		It("can list graphql schemas", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			bookinfoYaml, err := ioutil.ReadFile("graphql_test_data/bookinfo.yaml")
			Expect(err).NotTo(HaveOccurred())
			bookinfoGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(bookinfoYaml, bookinfoGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient1.EXPECT().ListGraphQLSchema(ctx).Return(&graphql_v1alpha1.GraphQLSchemaList{
				Items: []graphql_v1alpha1.GraphQLSchema{
					*petstoreGraphqlSchema,
				},
			}, nil)
			mockGraphqlSchemaClient2.EXPECT().ListGraphQLSchema(ctx).Return(&graphql_v1alpha1.GraphQLSchemaList{
				Items: []graphql_v1alpha1.GraphQLSchema{
					*bookinfoGraphqlSchema,
				},
			}, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			resp, err := handler.ListGraphqlSchemas(ctx, &rpc_edge_v1.ListGraphqlSchemasRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(MatchProto(&rpc_edge_v1.ListGraphqlSchemasResponse{
				GraphqlSchemas: []*rpc_edge_v1.GraphqlSchema{
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
						Spec:         &petstoreGraphqlSchema.Spec,
						Status:       &petstoreGraphqlSchema.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
					},
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "bookinfo", Namespace: "gloo-system", ClusterName: "remote-cluster"},
						Spec:         &bookinfoGraphqlSchema.Spec,
						Status:       &bookinfoGraphqlSchema.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "remote-test", Namespace: "gloo-system"},
					},
				},
			}))
		})
	})

	Context("CreateGraphqlSchema", func() {
		It("can create a graphqlschema", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient1.EXPECT().CreateGraphQLSchema(ctx, gomock.Any()).Return(nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			resp, err := handler.CreateGraphqlSchema(ctx, &rpc_edge_v1.CreateGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:             &petstoreGraphqlSchema.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.CreateGraphqlSchemaResponse{
				GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &petstoreGraphqlSchema.Spec,
					Status:       &petstoreGraphqlSchema.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
				},
			}))
		})
	})

	Context("UpdateGraphqlSchema", func() {
		It("can update a graphqlschema", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			bookinfoYaml, err := ioutil.ReadFile("graphql_test_data/bookinfo.yaml")
			Expect(err).NotTo(HaveOccurred())
			bookinfoGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(bookinfoYaml, bookinfoGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient1.EXPECT().GetGraphQLSchema(ctx, gomock.Any()).Return(petstoreGraphqlSchema, nil)
			mockGraphqlSchemaClient1.EXPECT().UpdateGraphQLSchema(ctx, gomock.Any()).Return(nil)

			// change spec
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			resp, err := handler.UpdateGraphqlSchema(ctx, &rpc_edge_v1.UpdateGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:             &bookinfoGraphqlSchema.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.UpdateGraphqlSchemaResponse{
				GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &bookinfoGraphqlSchema.Spec,
					Status:       &petstoreGraphqlSchema.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
				},
			}))
		})
		It("errors if ref points to a nonexistent graphqlschema", func() {
			mockGraphqlSchemaClient1.EXPECT().GetGraphQLSchema(ctx, gomock.Any()).Return(nil, eris.New("not found!"))

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			_, err := handler.UpdateGraphqlSchema(ctx, &rpc_edge_v1.UpdateGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:             &graphql_v1alpha1.GraphQLSchemaSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Cannot edit a graphqlschema that does not exist: not found!"))
		})
	})

	Context("DeleteGraphqlSchema", func() {
		It("can delete a graphqlschema", func() {
			mockGraphqlSchemaClient1.EXPECT().DeleteGraphQLSchema(ctx, gomock.Any()).Return(nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockMCClientset)
			resp, err := handler.DeleteGraphqlSchema(ctx, &rpc_edge_v1.DeleteGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.DeleteGraphqlSchemaResponse{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			}))
		})
	})
})
