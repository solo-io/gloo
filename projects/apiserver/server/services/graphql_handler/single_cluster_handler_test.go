package graphql_handler_test

import (
	"context"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	mock_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/mocks"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	mock_graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1/mocks"
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

		mockGlooInstanceLister *mock_glooinstance_handler.MockSingleClusterGlooInstanceLister
		mockGraphqlClientset   *mock_graphql_v1beta1.MockClientset
		mockGraphqlApiClient   *mock_graphql_v1beta1.MockGraphQLApiClient
		mockSettingsClient     *mock_gloo_v1.MockSettingsClient

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
		mockGraphqlClientset = mock_graphql_v1beta1.NewMockClientset(ctrl)
		mockGraphqlApiClient = mock_graphql_v1beta1.NewMockGraphQLApiClient(ctrl)
		mockSettingsClient = mock_gloo_v1.NewMockSettingsClient(ctrl)

		mockGraphqlClientset.EXPECT().GraphQLApis().Return(mockGraphqlApiClient).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetGraphqlApi", func() {
		It("can get graphql api by ref", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient.EXPECT().GetGraphQLApi(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlApi, nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			resp, err := handler.GetGraphqlApi(ctx, &rpc_edge_v1.GetGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.GetGraphqlApiResponse{
				GraphqlApi: &rpc_edge_v1.GraphqlApi{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &petstoreGraphqlApi.Spec,
					Status:       &petstoreGraphqlApi.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
				},
			}))
		})

		It("returns error if graphql api not found", func() {
			mockGraphqlApiClient.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(nil, eris.New("error!"))
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			_, err := handler.GetGraphqlApi(ctx, &rpc_edge_v1.GetGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Failed to get graphqlapi: error!"))
		})

		It("returns error if gloo instance lister returns error", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient.EXPECT().GetGraphQLApi(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlApi, nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return(nil, eris.New("uh oh!"))

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			_, err = handler.GetGraphqlApi(ctx, &rpc_edge_v1.GetGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("uh oh!"))
		})
	})

	Context("GetGraphqlApiYaml", func() {
		It("can get graphql api yaml", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient.EXPECT().GetGraphQLApi(ctx, client.ObjectKey{
				Namespace: "ns",
				Name:      "petstore",
			}).Return(petstoreGraphqlApi, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			resp, err := handler.GetGraphqlApiYaml(ctx, &rpc_edge_v1.GetGraphqlApiYamlRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.YamlData.Yaml).To(MatchYAML(petstoreYaml))
		})
	})

	Context("ListGraphqlApis", func() {
		It("can list graphql apis", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			bookinfoYaml, err := ioutil.ReadFile("graphql_test_data/bookinfo.yaml")
			Expect(err).NotTo(HaveOccurred())
			bookinfoGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(bookinfoYaml, bookinfoGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			stitchedYaml, err := ioutil.ReadFile("graphql_test_data/stitched.yaml")
			Expect(err).NotTo(HaveOccurred())
			stitchedGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(stitchedYaml, stitchedGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockGraphqlApiClient.EXPECT().ListGraphQLApi(ctx).Return(&graphql_v1beta1.GraphQLApiList{
				Items: []graphql_v1beta1.GraphQLApi{
					*petstoreGraphqlApi,
					*bookinfoGraphqlApi,
					*stitchedGraphqlApi,
				},
			}, nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			resp, err := handler.ListGraphqlApis(ctx, &rpc_edge_v1.ListGraphqlApisRequest{})
			Expect(err).NotTo(HaveOccurred())
			// within each gloo instance, graphqlapis are sorted by namespace then name
			Expect(resp).To(MatchProto(&rpc_edge_v1.ListGraphqlApisResponse{
				GraphqlApis: []*rpc_edge_v1.GraphqlApiSummary{
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "bookinfo", Namespace: "gloo-system"},
						Status:       &bookinfoGraphqlApi.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
						ApiTypeSummary: &rpc_edge_v1.GraphqlApiSummary_Executable{
							Executable: &rpc_edge_v1.GraphqlApiSummary_ExecutableSchemaSummary{
								NumResolvers: 6,
							},
						},
					},
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "stitched-api", Namespace: "gloo-system"},
						Status:       &stitchedGraphqlApi.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
						ApiTypeSummary: &rpc_edge_v1.GraphqlApiSummary_Stitched{
							Stitched: &rpc_edge_v1.GraphqlApiSummary_StitchedSchemaSummary{
								NumSubschemas: 3,
							},
						},
					},
					{
						Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
						Status:       &petstoreGraphqlApi.Status,
						GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
						ApiTypeSummary: &rpc_edge_v1.GraphqlApiSummary_Executable{
							Executable: &rpc_edge_v1.GraphqlApiSummary_ExecutableSchemaSummary{
								NumResolvers: 2,
							},
						},
					},
				},
			}))
		})
	})

	Context("CreateGraphqlApi", func() {
		It("can create a graphqlapi", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
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
			mockGraphqlApiClient.EXPECT().CreateGraphQLApi(ctx, gomock.Any()).Return(nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			resp, err := handler.CreateGraphqlApi(ctx, &rpc_edge_v1.CreateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &petstoreGraphqlApi.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.CreateGraphqlApiResponse{
				GraphqlApi: &rpc_edge_v1.GraphqlApi{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &petstoreGraphqlApi.Spec,
					Status:       &petstoreGraphqlApi.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
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
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			_, err := handler.CreateGraphqlApi(ctx, &rpc_edge_v1.CreateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &graphql_v1beta1.GraphQLApiSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot perform update: UI is read-only."))
		})
	})

	Context("UpdateGraphqlApi", func() {
		It("can update a graphqlapi", func() {
			petstoreYaml, err := ioutil.ReadFile("graphql_test_data/petstore.yaml")
			Expect(err).NotTo(HaveOccurred())
			petstoreGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(petstoreYaml, petstoreGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			bookinfoYaml, err := ioutil.ReadFile("graphql_test_data/bookinfo.yaml")
			Expect(err).NotTo(HaveOccurred())
			bookinfoGraphqlApi := &graphql_v1beta1.GraphQLApi{}
			err = yaml.Unmarshal(bookinfoYaml, bookinfoGraphqlApi)
			Expect(err).NotTo(HaveOccurred())

			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)
			mockGraphqlApiClient.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(petstoreGraphqlApi, nil)
			mockGraphqlApiClient.EXPECT().UpdateGraphQLApi(ctx, gomock.Any()).Return(nil)
			mockGlooInstanceLister.EXPECT().ListGlooInstances(ctx).Return([]*rpc_edge_v1.GlooInstance{glooInstance}, nil)

			// change spec
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			resp, err := handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &bookinfoGraphqlApi.Spec,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.UpdateGraphqlApiResponse{
				GraphqlApi: &rpc_edge_v1.GraphqlApi{
					Metadata:     &rpc_edge_v1.ObjectMeta{Name: "petstore", Namespace: "ns"},
					Spec:         &bookinfoGraphqlApi.Spec,
					Status:       &petstoreGraphqlApi.Status,
					GlooInstance: &skv2_v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"},
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
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			_, err := handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &graphql_v1beta1.GraphQLApiSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot perform update: UI is read-only."))
		})
		It("errors if ref points to a nonexistent graphqlapi", func() {
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)
			mockGraphqlApiClient.EXPECT().GetGraphQLApi(ctx, gomock.Any()).Return(nil, eris.New("not found!"))

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			_, err := handler.UpdateGraphqlApi(ctx, &rpc_edge_v1.UpdateGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
				Spec:          &graphql_v1beta1.GraphQLApiSpec{},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Cannot edit a graphqlapi that does not exist: not found!"))
		})
	})

	Context("DeleteGraphqlApi", func() {
		It("can delete a graphqlapi", func() {
			mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
				Spec: gloo_v1.SettingsSpec{
					ConsoleOptions: &gloo_v1.ConsoleOptions{
						ReadOnly: &wrappers.BoolValue{Value: false},
					},
				},
			}, nil)
			mockGraphqlApiClient.EXPECT().DeleteGraphQLApi(ctx, gomock.Any()).Return(nil)

			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			resp, err := handler.DeleteGraphqlApi(ctx, &rpc_edge_v1.DeleteGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.DeleteGraphqlApiResponse{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
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
			handler := graphql_handler.NewSingleClusterGraphqlHandler(mockGraphqlClientset, mockGlooInstanceLister, mockSettingsClient)
			_, err := handler.DeleteGraphqlApi(ctx, &rpc_edge_v1.DeleteGraphqlApiRequest{
				GraphqlApiRef: &skv2_v1.ClusterObjectRef{Name: "petstore", Namespace: "ns"},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Cannot perform update: UI is read-only."))
		})
	})
})
