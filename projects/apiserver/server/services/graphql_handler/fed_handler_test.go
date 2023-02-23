package graphql_handler_test

import (
	"context"
	"os"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	mock_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/mocks"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	mock_graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1/mocks"
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

		mockGlooInstanceClient *mock_fed_v1.MockGlooInstanceClient
		mockSettingsClient     *mock_gloo_v1.MockSettingsClient
		mockMCClientset        *mock_graphql_v1beta1.MockMulticlusterClientset
		mockGraphqlClientset1  *mock_graphql_v1beta1.MockClientset
		mockGraphqlClientset2  *mock_graphql_v1beta1.MockClientset
		mockGraphqlApiClient1  *mock_graphql_v1beta1.MockGraphQLApiClient
		mockGraphqlApiClient2  *mock_graphql_v1beta1.MockGraphQLApiClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockGlooInstanceClient = mock_fed_v1.NewMockGlooInstanceClient(ctrl)
		mockSettingsClient = mock_gloo_v1.NewMockSettingsClient(ctrl)
		mockMCClientset = mock_graphql_v1beta1.NewMockMulticlusterClientset(ctrl)
		mockGraphqlClientset1 = mock_graphql_v1beta1.NewMockClientset(ctrl)
		mockGraphqlClientset2 = mock_graphql_v1beta1.NewMockClientset(ctrl)
		mockGraphqlApiClient1 = mock_graphql_v1beta1.NewMockGraphQLApiClient(ctrl)
		mockGraphqlApiClient2 = mock_graphql_v1beta1.NewMockGraphQLApiClient(ctrl)

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
		mockGraphqlClientset1.EXPECT().GraphQLApis().Return(mockGraphqlApiClient1).AnyTimes()
		mockGraphqlClientset2.EXPECT().GraphQLApis().Return(mockGraphqlApiClient2).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetGraphqlApi", func() {
		It("can get graphql api by ref", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient1.EXPECT().GetGraphQLApi(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlApi, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			resp, err := handler.GetGraphqlApi(ctx, &rpc_edge_v1.GetGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.GetGraphqlApiResponse{
				GraphqlApi: &rpc_edge_v1.GraphqlApi{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
					Spec:         &petstoreGraphqlApi.Spec,
					Status:       &petstoreGraphqlApi.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
				},
			}))
		})
	})

	Context("GetGraphqlApiYaml", func() {
		It("can get graphql api yaml", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient1.EXPECT().GetGraphQLApi(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlApi, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			resp, err := handler.GetGraphqlApiYaml(ctx, &rpc_edge_v1.GetGraphqlApiYamlRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.YamlData.Yaml).To(MatchYAML(petstoreYaml))
		})
	})

	Context("ListGraphqlApis", func() {
		It("can list graphql apis", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			bookinfoYaml, err := os.ReadFile("graphql_test_data/bookinfo.yaml")
			Expect(err).NotTo(HaveOccurred())
			bookinfoGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(bookinfoYaml, bookinfoGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			stitchedYaml, err := os.ReadFile("graphql_test_data/stitched.yaml")
			Expect(err).NotTo(HaveOccurred())
			stitchedGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(stitchedYaml, stitchedGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient1.EXPECT().ListGraphQLApi(ctx).Return(&graphql_v1beta1.GraphQLApiList{
				Items: []graphql_v1beta1.GraphQLApi{
					*petstoreGraphqlApi,
					*stitchedGraphqlApi,
				},
			}, nil)
			mockGraphqlApiClient2.EXPECT().ListGraphQLApi(ctx).Return(&graphql_v1beta1.GraphQLApiList{
				Items: []graphql_v1beta1.GraphQLApi{
					*bookinfoGraphqlApi,
				},
			}, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			resp, err := handler.ListGraphqlApis(ctx, &rpc_edge_v1.ListGraphqlApisRequest{})
			Expect(err).NotTo(HaveOccurred())
			// within each gloo instance, graphqlapis are sorted by namespace then name
			Expect(resp).To(MatchProto(&rpc_edge_v1.ListGraphqlApisResponse{
				GraphqlApis: []*rpc_edge_v1.GraphqlApiSummary{
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "stitched-api", Namespace: "gloo-system", ClusterName: "local-cluster"},
						Status:       &stitchedGraphqlApi.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
						ApiTypeSummary: &rpc_edge_v1.GraphqlApiSummary_Stitched{
							Stitched: &rpc_edge_v1.GraphqlApiSummary_StitchedSchemaSummary{
								NumSubschemas: 3,
							},
						},
					},
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
						Status:       &petstoreGraphqlApi.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
						ApiTypeSummary: &rpc_edge_v1.GraphqlApiSummary_Executable{
							Executable: &rpc_edge_v1.GraphqlApiSummary_ExecutableSchemaSummary{
								NumResolvers: 2,
							},
						},
					},
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "bookinfo", Namespace: "gloo-system", ClusterName: "remote-cluster"},
						Status:       &bookinfoGraphqlApi.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "remote-test", Namespace: "gloo-system"},
						ApiTypeSummary: &rpc_edge_v1.GraphqlApiSummary_Executable{
							Executable: &rpc_edge_v1.GraphqlApiSummary_ExecutableSchemaSummary{
								NumResolvers: 6,
							},
						},
					},
				},
			}))
		})
	})

	Context("CreateGraphqlApi", func() {
		It("can create a graphqlapi", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)
			mockGraphqlApiClient1.EXPECT().CreateGraphQLApi(ctx, gomock.Any()).Return(nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			resp, err := handler.CreateGraphqlApi(ctx, &rpc_edge_v1.CreateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:          &petstoreGraphqlApi.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.CreateGraphqlApiResponse{
				GraphqlApi: &rpc_edge_v1.GraphqlApi{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
					Spec:         &petstoreGraphqlApi.Spec,
					Status:       &petstoreGraphqlApi.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
				},
			}))
		})
		It("cannot create graphqlapi if readonly is true", func() {
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: true},
					},
				},
			}, nil)
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err := handler.CreateGraphqlApi(ctx, &rpc_edge_v1.CreateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &graphql_v1beta1.GraphQLApiSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot perform update: UI is read-only."))
		})
		It("cannot create graphqlapi if validation fails", func() {
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{
				Spec: graphql_v1beta1.GraphQLApiSpec{
					Schema: &graphql_v1beta1.GraphQLApiSpec_ExecutableSchema{
						ExecutableSchema: &graphql_v1beta1.ExecutableSchema{
							SchemaDefinition: `
	type Query {
		productsForHome: asdfas[Product]
	}
`,
						},
					},
				},
			}

			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err := handler.CreateGraphqlApi(ctx, &rpc_edge_v1.CreateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:          &petstoreGraphqlApi.Spec,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to parse graphql schema"))
		})
	})

	Context("UpdateGraphqlApi", func() {
		It("can update a graphqlapi", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			updatedPetstoreGraphqlApi := petstoreGraphqlApi.DeepCopy()
			updatedPetstoreGraphqlApi.Spec.GetExecutableSchema().SchemaDefinition = `
      type Query {
        productsForHome: [Product]
      }

      type Product {
        "this is an id description"
        id: String
        "this is an author description"
        author: String
		"this is a new field"
		newField: String
      }
`

			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil).AnyTimes()
			mockGraphqlApiClient1.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(petstoreGraphqlApi, nil)
			mockGraphqlApiClient1.EXPECT().UpdateGraphQLApi(ctx, gomock.Any()).Return(nil)

			// change spec
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			resp, err := handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:          &updatedPetstoreGraphqlApi.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.UpdateGraphqlApiResponse{
				GraphqlApi: &rpc_edge_v1.GraphqlApi{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
					Spec:         &updatedPetstoreGraphqlApi.Spec,
					Status:       &updatedPetstoreGraphqlApi.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "local-test", Namespace: "gloo-system"},
				},
			}))
		})
		It("cannot update graphqlapi if readonly is true", func() {
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: true},
					},
				},
			}, nil)
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err := handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &graphql_v1beta1.GraphQLApiSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot perform update: UI is read-only."))
		})
		It("errors if ref points to a nonexistent graphqlapi", func() {
			mockGraphqlApiClient1.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(nil, eris.New("not found!"))
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err := handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:          &graphql_v1beta1.GraphQLApiSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Cannot edit a graphqlapi that does not exist: not found!"))
		})
		It("cannot update graphqlapi if validation fails", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			updatedPetstoreGraphqlApi := petstoreGraphqlApi.DeepCopy()
			updatedPetstoreGraphqlApi.Spec.GetExecutableSchema().SchemaDefinition = `
      type Query {
        productsForHome: asdfasdf[Product]
      }

      type Product {
        "this is an id description"
        id: String
        "this is an author description"
        author: String
      }
`

			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)
			mockGraphqlApiClient1.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(petstoreGraphqlApi, nil)

			// change spec
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err = handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:          &updatedPetstoreGraphqlApi.Spec,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("unable to parse graphql schema"))
		})
		It("cannot update graphqlapi with breaking changes", func() {
			petstoreYaml, err := os.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			updatedPetstoreGraphqlApi := petstoreGraphqlApi.DeepCopy()
			// a field was removed from Product, which is a breaking change
			updatedPetstoreGraphqlApi.Spec.GetExecutableSchema().SchemaDefinition = `
      type Query {
        productsForHome: [Product]
      }

      type Product {
        "this is an id description"
        id: String
      }
`
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
					GraphqlOptions: &gloo_v1.GraphqlOptions{
						SchemaChangeValidationOptions: &gloo_v1.GraphqlOptions_SchemaChangeValidationOptions{
							RejectBreakingChanges: &wrappers.BoolValue{
								Value: true,
							},
						},
					},
				},
			}, nil).AnyTimes()
			mockGraphqlApiClient1.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(petstoreGraphqlApi, nil)

			// change spec
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err = handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
				Spec:          &updatedPetstoreGraphqlApi.Spec,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Rejected graphqlapi update due to validation error"))
		})
	})

	Context("DeleteGraphqlApi", func() {
		It("can delete a graphqlapi", func() {
			mockGraphqlApiClient1.EXPECT().DeleteGraphQLApi(ctx, gomock.Any()).Return(nil)
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)

			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			resp, err := handler.DeleteGraphqlApi(ctx, &rpc_edge_v1.DeleteGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.DeleteGraphqlApiResponse{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns", ClusterName: "local-cluster"},
			}))
		})
		It("cannot delete graphqlapi if readonly is true", func() {
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: true},
					},
				},
			}, nil)
			handler := graphql_handler.NewFedGraphqlHandler(mockGlooInstanceClient, mockSettingsClient, mockMCClientset)
			_, err := handler.DeleteGraphqlApi(ctx, &rpc_edge_v1.DeleteGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot perform update: UI is read-only."))
		})
	})
})
