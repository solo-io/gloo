package glooinstance_handler_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	mock_envoy_admin "github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/envoy_admin/mocks"
	mock_glooinstance_handler "github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/mocks"
	"k8s.io/client-go/discovery"
)

var _ = Describe("single cluster gloo instance handler", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context

		mockGlooInstanceLister *mock_glooinstance_handler.MockSingleClusterGlooInstanceLister
		mockDiscoveryClient    *discovery.DiscoveryClient
		mockEnvoyAdminClient   *mock_envoy_admin.MockEnvoyAdminClient
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockGlooInstanceLister = mock_glooinstance_handler.NewMockSingleClusterGlooInstanceLister(ctrl)
		mockDiscoveryClient = &discovery.DiscoveryClient{}
		mockEnvoyAdminClient = mock_envoy_admin.NewMockEnvoyAdminClient(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("EnvoyAdminClient", func() {
		It("can get config dumps", func() {
			mockGlooInstanceLister.EXPECT().GetGlooInstance(ctx, &v1.ObjectRef{
				Name:      "test",
				Namespace: "gloo-system",
			}).Return(testGlooInstance, nil)
			mockEnvoyAdminClient.EXPECT().GetConfigs(ctx, testGlooInstance, gomock.Any()).Return([]*rpc_edge_v1.ConfigDump{
				{
					Name: "gateway-proxy",
					Raw:  "test-proxy-config-dump",
				},
			}, nil)
			handler := glooinstance_handler.NewSingleClusterGlooInstanceHandler(mockGlooInstanceLister, mockDiscoveryClient.RESTClient(), mockEnvoyAdminClient)
			resp, err := handler.GetConfigDumps(ctx, &rpc_edge_v1.GetConfigDumpsRequest{
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

		It("handles an error in the EnvoyAdminClient", func() {
			mockGlooInstanceLister.EXPECT().GetGlooInstance(ctx, &v1.ObjectRef{
				Name:      "test",
				Namespace: "gloo-system",
			}).Return(testGlooInstance, nil)
			mockEnvoyAdminClient.EXPECT().GetConfigs(ctx, testGlooInstance, gomock.Any()).Return([]*rpc_edge_v1.ConfigDump{}, eris.New("test"))
			handler := glooinstance_handler.NewSingleClusterGlooInstanceHandler(mockGlooInstanceLister, mockDiscoveryClient.RESTClient(), mockEnvoyAdminClient)
			_, err := handler.GetConfigDumps(ctx, &rpc_edge_v1.GetConfigDumpsRequest{
				GlooInstanceRef: &v1.ObjectRef{
					Name:      "test",
					Namespace: "gloo-system",
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("test"))
		})

		It("can get upstream hosts", func() {
			mockGlooInstanceLister.EXPECT().GetGlooInstance(ctx, &v1.ObjectRef{
				Name:      "test",
				Namespace: "gloo-system",
			}).Return(testGlooInstance, nil)
			upstreamHosts := map[string]*rpc_edge_v1.HostList{
				"gloo-system.upstream1": {
					Hosts: []*rpc_edge_v1.Host{
						{Address: "1.2.3.4", Port: 80, Weight: 7, ProxyRef: &v1.ObjectRef{Name: "proxy1", Namespace: "ns1"}},
						{Address: "4.5.6.7", Port: 12, Weight: 3, ProxyRef: &v1.ObjectRef{Name: "proxy2", Namespace: "ns2"}},
					},
				},
			}
			mockEnvoyAdminClient.EXPECT().GetHostsByUpstream(ctx, testGlooInstance, gomock.Any()).Return(upstreamHosts, nil)
			handler := glooinstance_handler.NewSingleClusterGlooInstanceHandler(mockGlooInstanceLister, mockDiscoveryClient.RESTClient(), mockEnvoyAdminClient)
			resp, err := handler.GetUpstreamHosts(ctx, &rpc_edge_v1.GetUpstreamHostsRequest{
				GlooInstanceRef: &v1.ObjectRef{
					Name:      "test",
					Namespace: "gloo-system",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).To(Equal(&rpc_edge_v1.GetUpstreamHostsResponse{
				UpstreamHosts: upstreamHosts,
			}))
		})
	})

})
