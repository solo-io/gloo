package glooinstance_handler_test

import (
	"context"

	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	mock_config_getter "github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/config_getter/mocks"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	mock_multicluster "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster/mocks"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("gloo instance and cluster handler", func() {
	var (
		ctrl                 *gomock.Controller
		ctx                  context.Context
		clusterClient        *mock_multicluster.MockClusterSet
		instanceClient       *mock_v1.MockGlooInstanceClient
		mockGetter           *mock_config_getter.MockEnvoyConfigDumpGetter
		testGlooInstanceList *fedv1.GlooInstanceList
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		clusterClient = mock_multicluster.NewMockClusterSet(ctrl)
		instanceClient = mock_v1.NewMockGlooInstanceClient(ctrl)
		mockGetter = mock_config_getter.NewMockEnvoyConfigDumpGetter(ctrl)

		local := fedv1.GlooInstance{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "local-test",
				Namespace: "gloo-system",
			},
			Spec: types.GlooInstanceSpec{
				Cluster: "local",
			},
			Status: types.GlooInstanceStatus{},
		}
		remote := fedv1.GlooInstance{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "remote-test",
				Namespace: "gloo-system",
			},
			Spec: types.GlooInstanceSpec{
				Cluster: "remote-cluster",
			},
			Status: types.GlooInstanceStatus{},
		}
		testGlooInstanceList = &fedv1.GlooInstanceList{
			Items: []fedv1.GlooInstance{
				local,
				remote,
			},
		}
		clusterClient.EXPECT().ListClusters().Return([]string{"local", "remote-cluster"}).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can list gloo instances", func() {

		instanceClient.EXPECT().
			ListGlooInstance(ctx).
			Return(testGlooInstanceList, nil)

		clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
		resp, err := clusterServer.ListGlooInstances(ctx, &rpc_v1.ListGlooInstancesRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(&rpc_v1.ListGlooInstancesResponse{
			GlooInstances: []*rpc_v1.GlooInstance{
				&rpc_v1.GlooInstance{
					Metadata: &rpc_v1.ObjectMeta{
						Name:      "local-test",
						Namespace: "gloo-system",
					},
					Spec: &types.GlooInstanceSpec{
						Cluster: "local",
					},
					Status: &types.GlooInstanceStatus{},
				},
				&rpc_v1.GlooInstance{
					Metadata: &rpc_v1.ObjectMeta{
						Name:      "remote-test",
						Namespace: "gloo-system",
					},
					Spec: &types.GlooInstanceSpec{
						Cluster: "remote-cluster",
					},
					Status: &types.GlooInstanceStatus{},
				},
			}}))
	})

	It("can list cluster details", func() {
		instanceClient.EXPECT().
			ListGlooInstance(ctx).
			Return(testGlooInstanceList, nil)

		clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
		resp, err := clusterServer.ListClusterDetails(ctx, &rpc_v1.ListClusterDetailsRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(&rpc_v1.ListClusterDetailsResponse{
			ClusterDetails: []*rpc_v1.ClusterDetails{
				&rpc_v1.ClusterDetails{
					Cluster: "local",
					GlooInstances: []*rpc_v1.GlooInstance{
						&rpc_v1.GlooInstance{
							Metadata: &rpc_v1.ObjectMeta{
								Name:      "local-test",
								Namespace: "gloo-system",
							},
							Spec: &types.GlooInstanceSpec{
								Cluster: "local",
							},
							Status: &types.GlooInstanceStatus{},
						},
					},
				},
				&rpc_v1.ClusterDetails{
					Cluster: "remote-cluster",
					GlooInstances: []*rpc_v1.GlooInstance{
						&rpc_v1.GlooInstance{
							Metadata: &rpc_v1.ObjectMeta{
								Name:      "remote-test",
								Namespace: "gloo-system",
							},
							Spec: &types.GlooInstanceSpec{
								Cluster: "remote-cluster",
							},
							Status: &types.GlooInstanceStatus{},
						},
					},
				},
			},
		}))
	})

	Context("ConfigDumpGetter", func() {
		It("can get config dumps", func() {
			glooInstance := fedv1.GlooInstance{
				Spec: types.GlooInstanceSpec{
					Cluster:      "mgmt",
					IsEnterprise: false,
					ControlPlane: &types.GlooInstanceSpec_ControlPlane{
						Namespace: "gloo-system",
					},
					Proxies: []*types.GlooInstanceSpec_Proxy{&types.GlooInstanceSpec_Proxy{
						Name: "gateway-proxy",
					}},
				},
				Status: types.GlooInstanceStatus{},
			}

			mockGetter.EXPECT().GetConfigs(ctx, glooInstance).Return([]*rpc_v1.ConfigDump{
				&rpc_v1.ConfigDump{
					Name: "gateway-proxy",
					Raw:  "test-proxy-config-dump",
				},
			}, nil)
			instanceClient.EXPECT().GetGlooInstance(ctx, client.ObjectKey{Name: "test", Namespace: "gloo-system"}).Return(
				&glooInstance, nil)
			clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
			resp, err := clusterServer.GetConfigDumps(ctx, &rpc_v1.GetConfigDumpsRequest{
				GlooInstanceRef: &v1.ObjectRef{
					Name:      "test",
					Namespace: "gloo-system",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_v1.GetConfigDumpsResponse{
				ConfigDumps: []*rpc_v1.ConfigDump{
					&rpc_v1.ConfigDump{
						Name: "gateway-proxy",
						Raw:  "test-proxy-config-dump",
					},
				},
			}))
		})

		It("handles an error in the EnvoyConfigDumpGetter", func() {
			mockGetter := mock_config_getter.NewMockEnvoyConfigDumpGetter(ctrl)
			glooInstance := fedv1.GlooInstance{
				Spec: types.GlooInstanceSpec{
					Cluster:      "mgmt",
					IsEnterprise: false,
					ControlPlane: &types.GlooInstanceSpec_ControlPlane{
						Namespace: "gloo-system",
					},
					Proxies: []*types.GlooInstanceSpec_Proxy{&types.GlooInstanceSpec_Proxy{}},
				},
				Status: types.GlooInstanceStatus{},
			}

			mockGetter.EXPECT().GetConfigs(ctx, glooInstance).Return([]*rpc_v1.ConfigDump{}, eris.New("test"))
			instanceClient.EXPECT().GetGlooInstance(ctx, client.ObjectKey{Name: "test", Namespace: "gloo-system"}).Return(
				&glooInstance, nil)
			clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
			_, err := clusterServer.GetConfigDumps(ctx, &rpc_v1.GetConfigDumpsRequest{
				GlooInstanceRef: &v1.ObjectRef{
					Name:      "test",
					Namespace: "gloo-system",
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("test"))
		})
	})

})
