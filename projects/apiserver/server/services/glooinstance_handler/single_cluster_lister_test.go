package glooinstance_handler_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	mock_apps_v1 "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/mocks"
	mock_core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/mocks"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	enterprise_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	mock_enterprise_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1/mocks"
	gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	mock_gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1/mocks"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	mock_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/mocks"
	ratelimit_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	mock_ratelimit_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1/mocks"
	. "github.com/solo-io/solo-kit/test/matchers"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	apps_v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testGlooInstance = &rpc_edge_v1.GlooInstance{
		Metadata: &rpc_edge_v1.ObjectMeta{
			Name:      "gloo",
			Namespace: "gloo-system",
		},
		Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
			Cluster:      "cluster",
			IsEnterprise: true,
			ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
				Version:           "1.2.3",
				Namespace:         "gloo-system",
				WatchedNamespaces: []string{"ns1", "ns2", "gloo-system"},
			},
			Proxies: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{},
			Check: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{
				Gateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 2,
				},
				MatchableHttpGateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 1,
					Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
						{
							Ref: &v1.ObjectRef{Name: "httpgw1", Namespace: "gloo-system"},
						},
					},
				},
				MatchableTcpGateways: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 1,
					Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
						{
							Ref: &v1.ObjectRef{Name: "tcpgw1", Namespace: "gloo-system"},
						},
					},
				},
				RateLimitConfigs: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 2,
					Errors: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
						{
							Ref: &v1.ObjectRef{Name: "rlc1", Namespace: "gloo-system"},
						},
					},
				},
				VirtualServices: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{},
				RouteTables:     &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{},
				AuthConfigs:     &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{},
				Settings: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 1,
				},
				Upstreams: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 3,
					Errors: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
						{
							Ref:     &v1.ObjectRef{Name: "us1", Namespace: "ns2"},
							Message: "i don't need a reason",
						},
					},
					Warnings: []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
						{
							Ref:     &v1.ObjectRef{Name: "us2", Namespace: "ns1"},
							Message: "uh oh",
						},
					},
				},
				UpstreamGroups: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{},
				Proxies:        &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{},
				Deployments: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
					Total: 2,
				},
				Pods: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{},
			},
		},
	}
)

var _ = Describe("single cluster gloo instance lister", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context

		// clientsets
		mockCoreClientset           *mock_core_v1.MockClientset
		mockAppsClientset           *mock_apps_v1.MockClientset
		mockGatewayClientset        *mock_gateway_v1.MockClientset
		mockGlooClientset           *mock_gloo_v1.MockClientset
		mockEnterpriseGlooClientset *mock_enterprise_gloo_v1.MockClientset
		mockRateLimitClientset      *mock_ratelimit_v1alpha1.MockClientset

		// clients
		mockServiceClient              *mock_core_v1.MockServiceClient
		mockPodClient                  *mock_core_v1.MockPodClient
		mockNodeClient                 *mock_core_v1.MockNodeClient
		mockDeploymentClient           *mock_apps_v1.MockDeploymentClient
		mockDaemonSetClient            *mock_apps_v1.MockDaemonSetClient
		mockGatewayClient              *mock_gateway_v1.MockGatewayClient
		mockVirtualServiceClient       *mock_gateway_v1.MockVirtualServiceClient
		mockRouteTableClient           *mock_gateway_v1.MockRouteTableClient
		mockMatchableHttpGatewayClient *mock_gateway_v1.MockMatchableHttpGatewayClient
		mockMatchableTcpGatewayClient  *mock_gateway_v1.MockMatchableTcpGatewayClient
		mockSettingsClient             *mock_gloo_v1.MockSettingsClient
		mockUpstreamClient             *mock_gloo_v1.MockUpstreamClient
		mockUpstreamGroupClient        *mock_gloo_v1.MockUpstreamGroupClient
		mockProxyClient                *mock_gloo_v1.MockProxyClient
		mockAuthConfigClient           *mock_enterprise_gloo_v1.MockAuthConfigClient
		mockRateLimitConfigClient      *mock_ratelimit_v1alpha1.MockRateLimitConfigClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())

		// core clientset
		mockCoreClientset = mock_core_v1.NewMockClientset(ctrl)
		mockServiceClient = mock_core_v1.NewMockServiceClient(ctrl)
		mockPodClient = mock_core_v1.NewMockPodClient(ctrl)
		mockNodeClient = mock_core_v1.NewMockNodeClient(ctrl)
		mockCoreClientset.EXPECT().Services().Return(mockServiceClient)
		mockCoreClientset.EXPECT().Pods().Return(mockPodClient).AnyTimes()
		mockCoreClientset.EXPECT().Nodes().Return(mockNodeClient)

		// apps clientset
		mockAppsClientset = mock_apps_v1.NewMockClientset(ctrl)
		mockDeploymentClient = mock_apps_v1.NewMockDeploymentClient(ctrl)
		mockDaemonSetClient = mock_apps_v1.NewMockDaemonSetClient(ctrl)
		mockAppsClientset.EXPECT().Deployments().Return(mockDeploymentClient).AnyTimes()
		mockAppsClientset.EXPECT().DaemonSets().Return(mockDaemonSetClient)

		// gateway clientset
		mockGatewayClientset = mock_gateway_v1.NewMockClientset(ctrl)
		mockGatewayClient = mock_gateway_v1.NewMockGatewayClient(ctrl)
		mockVirtualServiceClient = mock_gateway_v1.NewMockVirtualServiceClient(ctrl)
		mockRouteTableClient = mock_gateway_v1.NewMockRouteTableClient(ctrl)
		mockMatchableHttpGatewayClient = mock_gateway_v1.NewMockMatchableHttpGatewayClient(ctrl)
		mockMatchableTcpGatewayClient = mock_gateway_v1.NewMockMatchableTcpGatewayClient(ctrl)
		mockGatewayClientset.EXPECT().Gateways().Return(mockGatewayClient)
		mockGatewayClientset.EXPECT().VirtualServices().Return(mockVirtualServiceClient)
		mockGatewayClientset.EXPECT().RouteTables().Return(mockRouteTableClient)
		mockGatewayClientset.EXPECT().MatchableHttpGateways().Return(mockMatchableHttpGatewayClient)
		mockGatewayClientset.EXPECT().MatchableTcpGateways().Return(mockMatchableTcpGatewayClient)

		// gloo clientset
		mockGlooClientset = mock_gloo_v1.NewMockClientset(ctrl)
		mockSettingsClient = mock_gloo_v1.NewMockSettingsClient(ctrl)
		mockUpstreamClient = mock_gloo_v1.NewMockUpstreamClient(ctrl)
		mockUpstreamGroupClient = mock_gloo_v1.NewMockUpstreamGroupClient(ctrl)
		mockProxyClient = mock_gloo_v1.NewMockProxyClient(ctrl)
		mockGlooClientset.EXPECT().Settings().Return(mockSettingsClient)
		mockGlooClientset.EXPECT().Upstreams().Return(mockUpstreamClient)
		mockGlooClientset.EXPECT().UpstreamGroups().Return(mockUpstreamGroupClient)
		mockGlooClientset.EXPECT().Proxies().Return(mockProxyClient)

		// enterprise gloo clientset
		mockEnterpriseGlooClientset = mock_enterprise_gloo_v1.NewMockClientset(ctrl)
		mockAuthConfigClient = mock_enterprise_gloo_v1.NewMockAuthConfigClient(ctrl)
		mockEnterpriseGlooClientset.EXPECT().AuthConfigs().Return(mockAuthConfigClient)

		// ratelimit clientset
		mockRateLimitClientset = mock_ratelimit_v1alpha1.NewMockClientset(ctrl)
		mockRateLimitConfigClient = mock_ratelimit_v1alpha1.NewMockRateLimitConfigClient(ctrl)
		mockRateLimitClientset.EXPECT().RateLimitConfigs().Return(mockRateLimitConfigClient)

		// mock data
		mockServiceClient.EXPECT().ListService(ctx).Return(&core_v1.ServiceList{
			Items: []core_v1.Service{},
		}, nil).AnyTimes()
		mockPodClient.EXPECT().ListPod(ctx).Return(&core_v1.PodList{
			Items: []core_v1.Pod{},
		}, nil).AnyTimes()
		mockNodeClient.EXPECT().GetNode(ctx, gomock.Any()).Return(&core_v1.Node{
			ObjectMeta: metav1.ObjectMeta{},
			Spec:       core_v1.NodeSpec{},
		}, nil).AnyTimes()
		mockNodeClient.EXPECT().ListNode(ctx).Return(&core_v1.NodeList{
			Items: []core_v1.Node{},
		}, nil).AnyTimes()
		deployment := &apps_v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "gloo",
				Namespace: "gloo-system",
				Labels:    map[string]string{"app": "gloo", "gloo": "gloo"},
			},
			Spec: apps_v1.DeploymentSpec{
				Template: core_v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{},
					Spec: core_v1.PodSpec{
						Containers: []core_v1.Container{
							{
								Image: "quay.io/solo-io/gloo-ee:1.2.3",
							},
						},
					},
				},
			},
		}
		mockDeploymentClient.EXPECT().GetDeployment(ctx, gomock.Any()).Return(deployment, nil)
		mockDeploymentClient.EXPECT().ListDeployment(ctx).Return(&apps_v1.DeploymentList{
			Items: []apps_v1.Deployment{
				*deployment,
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "another-gloo",
						Namespace: "gloo-system",
						Labels:    map[string]string{"app": "gloo", "gloo": "gloo"},
					},
					Spec: apps_v1.DeploymentSpec{},
				},
			},
		}, nil).AnyTimes()
		mockDaemonSetClient.EXPECT().ListDaemonSet(ctx).Return(&apps_v1.DaemonSetList{
			Items: []apps_v1.DaemonSet{},
		}, nil).AnyTimes()
		mockGatewayClient.EXPECT().ListGateway(ctx).Return(&gateway_v1.GatewayList{
			Items: []gateway_v1.Gateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gw1",
						Namespace: "ns1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "gw2",
						Namespace: "ns2",
					},
				},
			},
		}, nil)
		mockVirtualServiceClient.EXPECT().ListVirtualService(ctx).Return(&gateway_v1.VirtualServiceList{
			Items: []gateway_v1.VirtualService{},
		}, nil)
		mockRouteTableClient.EXPECT().ListRouteTable(ctx).Return(&gateway_v1.RouteTableList{
			Items: []gateway_v1.RouteTable{},
		}, nil)
		mockMatchableHttpGatewayClient.EXPECT().ListMatchableHttpGateway(ctx).Return(&gateway_v1.MatchableHttpGatewayList{
			Items: []gateway_v1.MatchableHttpGateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "httpgw1",
						Namespace: "gloo-system",
					},
					Status: gateway_v1.MatchableHttpGatewayStatus{
						State: gateway_v1.MatchableHttpGatewayStatus_Warning,
					},
				},
			},
		}, nil)
		mockMatchableTcpGatewayClient.EXPECT().ListMatchableTcpGateway(ctx).Return(&gateway_v1.MatchableTcpGatewayList{
			Items: []gateway_v1.MatchableTcpGateway{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "tcpgw1",
						Namespace: "gloo-system",
					},
					Status: gateway_v1.MatchableTcpGatewayStatus{
						State: gateway_v1.MatchableTcpGatewayStatus_Warning,
					},
				},
			},
		}, nil)
		mockRateLimitConfigClient.EXPECT().ListRateLimitConfig(ctx).Return(&ratelimit_v1alpha1.RateLimitConfigList{
			Items: []ratelimit_v1alpha1.RateLimitConfig{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rlc1",
						Namespace: "gloo-system",
					},
					Status: ratelimit_v1alpha1.RateLimitConfigStatus{
						State: ratelimit_v1alpha1.RateLimitConfigStatus_REJECTED,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rlc2",
						Namespace: "gloo-system",
					},
					Status: ratelimit_v1alpha1.RateLimitConfigStatus{
						State: ratelimit_v1alpha1.RateLimitConfigStatus_ACCEPTED,
					},
				},
			},
		}, nil)
		mockSettingsClient.EXPECT().GetSettings(ctx, gomock.Any()).Return(&gloo_v1.Settings{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "default",
				Namespace: "gloo-system",
			},
			Spec: gloo_v1.SettingsSpec{
				WatchNamespaces: []string{"ns1", "ns2", "gloo-system"},
			},
		}, nil)
		mockSettingsClient.EXPECT().ListSettings(ctx).Return(&gloo_v1.SettingsList{
			Items: []gloo_v1.Settings{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "default",
						Namespace: "gloo-system",
					},
					Spec: gloo_v1.SettingsSpec{
						WatchNamespaces: []string{"ns1", "ns2", "gloo-system"},
					},
				},
			},
		}, nil)
		mockUpstreamClient.EXPECT().ListUpstream(ctx).Return(&gloo_v1.UpstreamList{
			Items: []gloo_v1.Upstream{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "us1",
						Namespace: "ns2",
					},
					Status: gloo_v1.UpstreamStatus{State: gloo_v1.UpstreamStatus_Rejected, Reason: "i don't need a reason"},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "us2",
						Namespace: "ns1",
					},
					Status: gloo_v1.UpstreamStatus{State: gloo_v1.UpstreamStatus_Warning, Reason: "uh oh"},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "us3",
						Namespace: "ns2",
					},
					Status: gloo_v1.UpstreamStatus{State: gloo_v1.UpstreamStatus_Pending},
				},
			},
		}, nil)
		mockUpstreamGroupClient.EXPECT().ListUpstreamGroup(ctx).Return(&gloo_v1.UpstreamGroupList{
			Items: []gloo_v1.UpstreamGroup{},
		}, nil)
		mockProxyClient.EXPECT().ListProxy(ctx).Return(&gloo_v1.ProxyList{
			Items: []gloo_v1.Proxy{},
		}, nil)
		mockAuthConfigClient.EXPECT().ListAuthConfig(ctx).Return(&enterprise_gloo_v1.AuthConfigList{
			Items: []enterprise_gloo_v1.AuthConfig{},
		}, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can get a gloo instance", func() {
		lister := glooinstance_handler.NewSingleClusterGlooInstanceLister(
			mockCoreClientset,
			mockAppsClientset,
			mockGatewayClientset,
			mockGlooClientset,
			mockEnterpriseGlooClientset,
			mockRateLimitClientset,
		)
		glooInstance, err := lister.GetGlooInstance(ctx, &v1.ObjectRef{Name: "gloo", Namespace: "gloo-system"})
		Expect(err).NotTo(HaveOccurred())
		Expect(glooInstance).To(MatchProto(testGlooInstance))
	})

	It("can list gloo instances", func() {
		lister := glooinstance_handler.NewSingleClusterGlooInstanceLister(
			mockCoreClientset,
			mockAppsClientset,
			mockGatewayClientset,
			mockGlooClientset,
			mockEnterpriseGlooClientset,
			mockRateLimitClientset,
		)
		glooInstances, err := lister.ListGlooInstances(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(glooInstances).To(HaveLen(1))
		Expect(glooInstances[0]).To(MatchProto(testGlooInstance))
	})
})
