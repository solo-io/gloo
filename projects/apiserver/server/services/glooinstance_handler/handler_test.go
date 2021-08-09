package glooinstance_handler_test

import (
	"context"

	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
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
				Cluster:      "local",
				IsEnterprise: true,
				ControlPlane: &types.GlooInstanceSpec_ControlPlane{
					Version:           "v123",
					Namespace:         "ns1",
					WatchedNamespaces: []string{"a", "b", "c"},
				},
				Proxies: []*types.GlooInstanceSpec_Proxy{
					{
						Replicas:                      5,
						AvailableReplicas:             3,
						ReadyReplicas:                 2,
						WasmEnabled:                   false,
						ReadConfigMulticlusterEnabled: false,
						Version:                       "xyz",
						Name:                          "myname",
						Namespace:                     "myns",
						WorkloadControllerType:        types.GlooInstanceSpec_Proxy_DAEMON_SET,
						Zones:                         []string{"zone1", "zone2"},
						IngressEndpoints: []*types.GlooInstanceSpec_Proxy_IngressEndpoint{
							{
								Address: "1.2.3.4",
								Ports: []*types.GlooInstanceSpec_Proxy_IngressEndpoint_Port{
									{
										Port: 555,
										Name: "port1",
									},
								},
								ServiceName: "myservice",
							},
						},
					},
				},
				Region: "region6",
				Check: &types.GlooInstanceSpec_Check{
					Gateways: &types.GlooInstanceSpec_Check_Summary{
						Total: 5,
						Errors: []*types.GlooInstanceSpec_Check_Summary_ResourceReport{
							{
								Ref:     &v1.ObjectRef{Name: "some", Namespace: "ref"},
								Message: "my message",
							},
						},
					},
					RouteTables: &types.GlooInstanceSpec_Check_Summary{
						Total: 2,
						Warnings: []*types.GlooInstanceSpec_Check_Summary_ResourceReport{
							{
								Ref:     &v1.ObjectRef{Name: "another", Namespace: "ref"},
								Message: "this is a warning",
							},
						},
					},
				},
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
				ControlPlane: &types.GlooInstanceSpec_ControlPlane{
					Version:           "v123",
					Namespace:         "ns2",
					WatchedNamespaces: []string{"d", "e"},
				},
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
		resp, err := clusterServer.ListGlooInstances(ctx, &rpc_edge_v1.ListGlooInstancesRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(&rpc_edge_v1.ListGlooInstancesResponse{
			GlooInstances: []*rpc_edge_v1.GlooInstance{
				{
					Metadata: &rpc_edge_v1.ObjectMeta{
						Name:      "local-test",
						Namespace: "gloo-system",
					},
					Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
						Cluster:      "local",
						IsEnterprise: true,
						ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
							Version:           "v123",
							Namespace:         "ns1",
							WatchedNamespaces: []string{"a", "b", "c"},
						},
						Proxies: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
							{
								Replicas:                      5,
								AvailableReplicas:             3,
								ReadyReplicas:                 2,
								WasmEnabled:                   false,
								ReadConfigMulticlusterEnabled: false,
								Version:                       "xyz",
								Name:                          "myname",
								Namespace:                     "myns",
								WorkloadControllerType:        rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DAEMON_SET,
								Zones:                         []string{"zone1", "zone2"},
								IngressEndpoints: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint{
									{
										Address: "1.2.3.4",
										Ports: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port{
											{
												Port: 555,
												Name: "port1",
											},
										},
										ServiceName: "myservice",
									},
								},
							},
						},
						Region: "region6",
						Check: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{
							Gateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Total: 5,
								Errors: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
									{
										Ref:     &v1.ObjectRef{Name: "some", Namespace: "ref"},
										Message: "my message",
									},
								},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							VirtualServices: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							RouteTables: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Total:  2,
								Errors: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
									{
										Ref:     &v1.ObjectRef{Name: "another", Namespace: "ref"},
										Message: "this is a warning",
									},
								},
							},
							AuthConfigs: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Settings: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Upstreams: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							UpstreamGroups: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Proxies: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Deployments: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Pods: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
						},
					},
					Status: &rpc_edge_v1.GlooInstance_GlooInstanceStatus{},
				},
				{
					Metadata: &rpc_edge_v1.ObjectMeta{
						Name:      "remote-test",
						Namespace: "gloo-system",
					},
					Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
						Cluster: "remote-cluster",
						ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
							Version:           "v123",
							Namespace:         "ns2",
							WatchedNamespaces: []string{"d", "e"},
						},
						Proxies: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{},
						Check: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{
							Gateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							VirtualServices: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							RouteTables: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							AuthConfigs: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Settings: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Upstreams: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							UpstreamGroups: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Proxies: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Deployments: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
							Pods: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
								Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
								Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
							},
						},
					},
					Status: &rpc_edge_v1.GlooInstance_GlooInstanceStatus{},
				},
			}}))
	})

	It("can list cluster details", func() {
		instanceClient.EXPECT().
			ListGlooInstance(ctx).
			Return(testGlooInstanceList, nil)

		clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
		resp, err := clusterServer.ListClusterDetails(ctx, &rpc_edge_v1.ListClusterDetailsRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).To(Equal(&rpc_edge_v1.ListClusterDetailsResponse{
			ClusterDetails: []*rpc_edge_v1.ClusterDetails{
				{
					Cluster: "local",
					GlooInstances: []*rpc_edge_v1.GlooInstance{
						{
							Metadata: &rpc_edge_v1.ObjectMeta{
								Name:      "local-test",
								Namespace: "gloo-system",
							},
							Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
								Cluster:      "local",
								IsEnterprise: true,
								ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
									Version:           "v123",
									Namespace:         "ns1",
									WatchedNamespaces: []string{"a", "b", "c"},
								},
								Proxies: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
									{
										Replicas:                      5,
										AvailableReplicas:             3,
										ReadyReplicas:                 2,
										WasmEnabled:                   false,
										ReadConfigMulticlusterEnabled: false,
										Version:                       "xyz",
										Name:                          "myname",
										Namespace:                     "myns",
										WorkloadControllerType:        rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DAEMON_SET,
										Zones:                         []string{"zone1", "zone2"},
										IngressEndpoints: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint{
											{
												Address: "1.2.3.4",
												Ports: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port{
													{
														Port: 555,
														Name: "port1",
													},
												},
												ServiceName: "myservice",
											},
										},
									},
								},
								Region: "region6",
								Check: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{
									Gateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Total: 5,
										Errors: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
											{
												Ref:     &v1.ObjectRef{Name: "some", Namespace: "ref"},
												Message: "my message",
											},
										},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									VirtualServices: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									RouteTables: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Total:  2,
										Errors: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
											{
												Ref:     &v1.ObjectRef{Name: "another", Namespace: "ref"},
												Message: "this is a warning",
											},
										},
									},
									AuthConfigs: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Settings: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Upstreams: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									UpstreamGroups: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Proxies: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Deployments: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Pods: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
								},
							},
							Status: &rpc_edge_v1.GlooInstance_GlooInstanceStatus{},
						},
					},
				},
				{
					Cluster: "remote-cluster",
					GlooInstances: []*rpc_edge_v1.GlooInstance{
						{
							Metadata: &rpc_edge_v1.ObjectMeta{
								Name:      "remote-test",
								Namespace: "gloo-system",
							},
							Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
								Cluster: "remote-cluster",
								ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
									Version:           "v123",
									Namespace:         "ns2",
									WatchedNamespaces: []string{"d", "e"},
								},
								Proxies: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{},
								Check: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{
									Gateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									VirtualServices: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									RouteTables: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									AuthConfigs: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Settings: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Upstreams: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									UpstreamGroups: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Proxies: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Deployments: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
									Pods: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
										Errors:   []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
										Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{},
									},
								},
							},
							Status: &rpc_edge_v1.GlooInstance_GlooInstanceStatus{},
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
					Proxies: []*types.GlooInstanceSpec_Proxy{
						{
							Name: "gateway-proxy",
						},
					},
				},
				Status: types.GlooInstanceStatus{},
			}

			mockGetter.EXPECT().GetConfigs(ctx, glooInstance).Return([]*rpc_edge_v1.ConfigDump{
				{
					Name: "gateway-proxy",
					Raw:  "test-proxy-config-dump",
				},
			}, nil)
			instanceClient.EXPECT().GetGlooInstance(ctx, client.ObjectKey{Name: "test", Namespace: "gloo-system"}).Return(
				&glooInstance, nil)
			clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
			resp, err := clusterServer.GetConfigDumps(ctx, &rpc_edge_v1.GetConfigDumpsRequest{
				GlooInstanceRef: &v1.ObjectRef{
					Name:      "test",
					Namespace: "gloo-system",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.GetConfigDumpsResponse{
				ConfigDumps: []*rpc_edge_v1.ConfigDump{
					{
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
					Proxies: []*types.GlooInstanceSpec_Proxy{{}},
				},
				Status: types.GlooInstanceStatus{},
			}

			mockGetter.EXPECT().GetConfigs(ctx, glooInstance).Return([]*rpc_edge_v1.ConfigDump{}, eris.New("test"))
			instanceClient.EXPECT().GetGlooInstance(ctx, client.ObjectKey{Name: "test", Namespace: "gloo-system"}).Return(
				&glooInstance, nil)
			clusterServer := glooinstance_handler.NewGlooInstanceHandler(clusterClient, mockGetter, instanceClient)
			_, err := clusterServer.GetConfigDumps(ctx, &rpc_edge_v1.GetConfigDumpsRequest{
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
