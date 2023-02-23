package wasmfilter_handler_test

import (
	"context"

	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/mock/gomock"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	mock_gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1/mocks"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/options/wasm"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/wasmfilter_handler"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	mock_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/mocks"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("wasm filter handler", func() {
	var (
		ctrl                  *gomock.Controller
		ctx                   context.Context
		instanceClient        *mock_v1.MockGlooInstanceClient
		mcGatewayCRDClientset *mock_gateway_v1.MockMulticlusterClientset
		mockGatewayClientset1 *mock_gateway_v1.MockClientset
		mockGatewayClientset2 *mock_gateway_v1.MockClientset
		mockGatewayClient1    *mock_gateway_v1.MockGatewayClient
		mockGatewayClient2    *mock_gateway_v1.MockGatewayClient
		testGlooInstanceList  *fedv1.GlooInstanceList
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		instanceClient = mock_v1.NewMockGlooInstanceClient(ctrl)
		mcGatewayCRDClientset = mock_gateway_v1.NewMockMulticlusterClientset(ctrl)
		mockGatewayClientset1 = mock_gateway_v1.NewMockClientset(ctrl)
		mockGatewayClientset2 = mock_gateway_v1.NewMockClientset(ctrl)
		mockGatewayClient1 = mock_gateway_v1.NewMockGatewayClient(ctrl)
		mockGatewayClient2 = mock_gateway_v1.NewMockGatewayClient(ctrl)

		local := fedv1.GlooInstance{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "local-test",
				Namespace: "gloo-system",
			},
			Spec: types.GlooInstanceSpec{
				Cluster: "local-cluster",
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
		instanceClient.EXPECT().ListGlooInstance(ctx).Return(testGlooInstanceList, nil).AnyTimes()
		mcGatewayCRDClientset.EXPECT().Cluster("local-cluster").Return(mockGatewayClientset1, nil).AnyTimes()
		mcGatewayCRDClientset.EXPECT().Cluster("remote-cluster").Return(mockGatewayClientset2, nil).AnyTimes()
		mockGatewayClientset1.EXPECT().Gateways().Return(mockGatewayClient1).AnyTimes()
		mockGatewayClientset2.EXPECT().Gateways().Return(mockGatewayClient2).AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	buildGateway := func(cluster string) gatewayv1.Gateway {
		return gatewayv1.Gateway{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:        "test-gateway",
				Namespace:   "test-namespace",
				ClusterName: cluster,
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayType: &gatewayv1.GatewaySpec_HttpGateway{
					HttpGateway: &gatewayv1.HttpGateway{
						Options: &gloov1.HttpListenerOptions{
							Wasm: &wasm.PluginSource{
								Filters: []*wasm.WasmFilter{
									{
										Src: &wasm.WasmFilter_Image{
											Image: "wasm-image",
										},
										Config:      nil,
										FilterStage: nil,
										Name:        "filter-1",
										RootId:      "filter-1",
										VmType:      0,
									},
								},
							},
						},
					},
				},
			},
			Status: gatewayv1.GatewayStatus{
				State: 1,
			},
		}
	}

	It("can list wasm filters", func() {
		mockGatewayClient1.EXPECT().ListGateway(ctx).Return(&gatewayv1.GatewayList{
			Items: []gatewayv1.Gateway{buildGateway("local-cluster")},
		}, nil)
		mockGatewayClient2.EXPECT().ListGateway(ctx).Return(&gatewayv1.GatewayList{
			Items: []gatewayv1.Gateway{buildGateway("remote-cluster")},
		}, nil)

		wasmFilterServer := wasmfilter_handler.NewFedWasmFilterHandler(instanceClient, mcGatewayCRDClientset)
		resp, err := wasmFilterServer.ListWasmFilters(ctx, &rpc_edge_v1.ListWasmFiltersRequest{})
		Expect(err).NotTo(HaveOccurred())
		localGatewayRef := apiserverutils.ToClusterObjectRef("test-gateway", "test-namespace", "local-cluster")
		remoteGatewayRef := apiserverutils.ToClusterObjectRef("test-gateway", "test-namespace", "remote-cluster")
		Expect(resp).To(Equal(&rpc_edge_v1.ListWasmFiltersResponse{
			WasmFilters: []*rpc_edge_v1.WasmFilter{
				{
					Name:   "filter-1",
					RootId: "filter-1",
					Source: "wasm-image",
					Config: "",
					Locations: []*rpc_edge_v1.WasmFilter_Location{
						{
							GatewayRef: &localGatewayRef,
							GatewayStatus: &gatewayv1.GatewayStatus{
								State: 1,
							},
							GlooInstanceRef: &v1.ObjectRef{
								Name:      "local-test",
								Namespace: "gloo-system",
							},
						},
						{
							GatewayRef: &remoteGatewayRef,
							GatewayStatus: &gatewayv1.GatewayStatus{
								State: 1,
							},
							GlooInstanceRef: &v1.ObjectRef{
								Name:      "remote-test",
								Namespace: "gloo-system",
							},
						},
					},
				},
			},
		}))
	})

	It("can get wasm filter", func() {
		localGateway := buildGateway("local-cluster")
		mockGatewayClient1.EXPECT().GetGateway(ctx, client.ObjectKey{Name: "test-gateway", Namespace: "test-namespace"}).
			Return(&localGateway, nil)
		remoteGateway := buildGateway("local-cluster")
		mockGatewayClient2.EXPECT().GetGateway(ctx, client.ObjectKey{Name: "test-gateway", Namespace: "test-namespace"}).
			Return(&remoteGateway, nil)

		wasmFilterServer := wasmfilter_handler.NewFedWasmFilterHandler(instanceClient, mcGatewayCRDClientset)
		resp, err := wasmFilterServer.DescribeWasmFilter(ctx, &rpc_edge_v1.DescribeWasmFilterRequest{
			Name:   "filter-1",
			RootId: "filter-1",
			GatewayRef: &v1.ObjectRef{
				Name:      "test-gateway",
				Namespace: "test-namespace",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		localGatewayRef := apiserverutils.ToClusterObjectRef("test-gateway", "test-namespace", "local-cluster")
		remoteGatewayRef := apiserverutils.ToClusterObjectRef("test-gateway", "test-namespace", "remote-cluster")
		Expect(resp).To(Equal(&rpc_edge_v1.DescribeWasmFilterResponse{
			WasmFilter: &rpc_edge_v1.WasmFilter{
				Name:   "filter-1",
				RootId: "filter-1",
				Source: "wasm-image",
				Locations: []*rpc_edge_v1.WasmFilter_Location{
					{
						GatewayRef: &localGatewayRef,
						GatewayStatus: &gatewayv1.GatewayStatus{
							State: 1,
						},
						GlooInstanceRef: &v1.ObjectRef{
							Name:      "local-test",
							Namespace: "gloo-system",
						},
					},
					{
						GatewayRef: &remoteGatewayRef,
						GatewayStatus: &gatewayv1.GatewayStatus{
							State: 1,
						},
						GlooInstanceRef: &v1.ObjectRef{
							Name:      "remote-test",
							Namespace: "gloo-system",
						},
					},
				},
			},
		}))
	})

})
