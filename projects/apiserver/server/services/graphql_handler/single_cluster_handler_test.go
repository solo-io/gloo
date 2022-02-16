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
	mock_glooinstance_handler "github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/mocks"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/graphql_handler"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("single cluster graphql handler", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context

		mockGlooInstanceLister  *mock_glooinstance_handler.MockSingleClusterGlooInstanceLister
		mockGraphqlClientset    *mock_graphql_v1alpha1.MockClientset
		mockGraphqlSchemaClient *mock_graphql_v1alpha1.MockGraphQLSchemaClient

		glooInstance = &rpc_edge_v1.GlooInstance{
			Metadata: &rpc_edge_v1.ObjectMeta{
				Name:      "gloo",
				Namespace: "gloo-system",
			},
		}
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockGlooInstanceLister = mock_glooinstance_handler.NewMockSingleClusterGlooInstanceLister(ctrl)
		mockGraphqlClientset = mock_graphql_v1alpha1.NewMockClientset(ctrl)
		mockGraphqlSchemaClient = mock_graphql_v1alpha1.NewMockGraphQLSchemaClient(ctrl)

		mockGraphqlClientset.EXPECT().GraphQLSchemas().Return(mockGraphqlSchemaClient).AnyTimes()
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

			mockGraphqlSchemaClient.EXPECT().GetGraphQLSchema(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlSchema, nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			resp, err := handler.GetGraphqlSchema(ctx, &rpc_edge_v1.GetGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.GetGraphqlSchemaResponse{
				GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &petstoreGraphqlSchema.Spec,
					Status:       &petstoreGraphqlSchema.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
				},
			}))
		})

		It("returns error if graphql schema not found", func() {
			mockGraphqlSchemaClient.EXPECT().GetGraphQLSchema(ctx, gomock.Any()).Return(nil, eris.New("error!"))
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			_, err := handler.GetGraphqlSchema(ctx, &rpc_edge_v1.GetGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Failed to get graphqlschema: error!"))
		})

		It("returns error if gloo instance lister returns error", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient.EXPECT().GetGraphQLSchema(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlSchema, nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return(nil, eris.New("uh oh!"))

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			_, err = handler.GetGraphqlSchema(ctx, &rpc_edge_v1.GetGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("uh oh!"))
		})
	})

	Context("GetGraphqlSchemaYaml", func() {
		It("can get graphql schema yaml", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlSchema := &graphql_v1alpha1.GraphQLSchema{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlSchema)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlSchemaClient.EXPECT().GetGraphQLSchema(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlSchema, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			resp, err := handler.GetGraphqlSchemaYaml(ctx, &rpc_edge_v1.GetGraphqlSchemaYamlRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
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

			mockGraphqlSchemaClient.EXPECT().ListGraphQLSchema(ctx).Return(&graphql_v1alpha1.GraphQLSchemaList{
				Items: []graphql_v1alpha1.GraphQLSchema{
					*petstoreGraphqlSchema,
					*bookinfoGraphqlSchema,
				},
			}, nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			resp, err := handler.ListGraphqlSchemas(ctx, &rpc_edge_v1.ListGraphqlSchemasRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(MatchProto(&rpc_edge_v1.ListGraphqlSchemasResponse{
				GraphqlSchemas: []*rpc_edge_v1.GraphqlSchema{
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "bookinfo", Namespace: "gloo-system"},
						Spec:         &bookinfoGraphqlSchema.Spec,
						Status:       &bookinfoGraphqlSchema.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
					},
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
						Spec:         &petstoreGraphqlSchema.Spec,
						Status:       &petstoreGraphqlSchema.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
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

			mockGraphqlSchemaClient.EXPECT().CreateGraphQLSchema(ctx, gomock.Any()).Return(nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			resp, err := handler.CreateGraphqlSchema(ctx, &rpc_edge_v1.CreateGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:             &petstoreGraphqlSchema.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.CreateGraphqlSchemaResponse{
				GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &petstoreGraphqlSchema.Spec,
					Status:       &petstoreGraphqlSchema.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
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

			mockGraphqlSchemaClient.EXPECT().GetGraphQLSchema(ctx, gomock.Any()).Return(petstoreGraphqlSchema, nil)
			mockGraphqlSchemaClient.EXPECT().UpdateGraphQLSchema(ctx, gomock.Any()).Return(nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			// change spec
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			resp, err := handler.UpdateGraphqlSchema(ctx, &rpc_edge_v1.UpdateGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:             &bookinfoGraphqlSchema.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.UpdateGraphqlSchemaResponse{
				GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &bookinfoGraphqlSchema.Spec,
					Status:       &petstoreGraphqlSchema.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
				},
			}))
		})
		It("errors if ref points to a nonexistent graphqlschema", func() {
			mockGraphqlSchemaClient.EXPECT().GetGraphQLSchema(ctx, gomock.Any()).Return(nil, eris.New("not found!"))

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			_, err := handler.UpdateGraphqlSchema(ctx, &rpc_edge_v1.UpdateGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:             &graphql_v1alpha1.GraphQLSchemaSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Cannot edit a graphqlschema that does not exist: not found!"))
		})
	})

	Context("DeleteGraphqlSchema", func() {
		It("can delete a graphqlschema", func() {
			mockGraphqlSchemaClient.EXPECT().DeleteGraphQLSchema(ctx, gomock.Any()).Return(nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister)
			resp, err := handler.DeleteGraphqlSchema(ctx, &rpc_edge_v1.DeleteGraphqlSchemaRequest{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.DeleteGraphqlSchemaResponse{
				GraphqlSchemaRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			}))
		})
	})
})
