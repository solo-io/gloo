package glooinstance_handler_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	mock_config_getter "github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/config_getter/mocks"
	mock_glooinstance_handler "github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/mocks"
	"k8s.io/client-go/discovery"
)

var _ = Describe("single cluster gloo instance handler", func() {
	var (
		ctrl *gomock.Controller
		ctx  context.Context

		mockGlooInstanceLister *mock_glooinstance_handler.MockSingleClusterGlooInstanceLister
		mockDiscoveryClient    *discovery.DiscoveryClient
		mockEnvoyConfigClient  *mock_config_getter.MockEnvoyConfigDumpGetter
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		mockGlooInstanceLister = mock_glooinstance_handler.NewMockSingleClusterGlooInstanceLister(ctrl)
		mockDiscoveryClient = &discovery.DiscoveryClient{}
		mockEnvoyConfigClient = mock_config_getter.NewMockEnvoyConfigDumpGetter(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("ConfigDumpGetter", func() {
		It("can get config dumps", func() {
			mockGlooInstanceLister.EXPECT().GetGlooInstance(ctx, &v1.ObjectRef{
				Name:      "test",
				Namespace: "gloo-system",
			}).Return(testGlooInstance, nil)
			mockEnvoyConfigClient.EXPECT().GetConfigs(ctx, testGlooInstance, gomock.Any()).Return([]*rpc_edge_v1.ConfigDump{
				{
					Name: "gateway-proxy",
					Raw:  "test-proxy-config-dump",
				},
			}, nil)
			handler := glooinstance_handler.NewSingleClusterGlooInstanceHandler(mockGlooInstanceLister, *mockDiscoveryClient, mockEnvoyConfigClient)
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

		It("handles an error in the EnvoyConfigDumpGetter", func() {
			mockGlooInstanceLister.EXPECT().GetGlooInstance(ctx, &v1.ObjectRef{
				Name:      "test",
				Namespace: "gloo-system",
			}).Return(testGlooInstance, nil)
			mockEnvoyConfigClient.EXPECT().GetConfigs(ctx, testGlooInstance, gomock.Any()).Return([]*rpc_edge_v1.ConfigDump{}, eris.New("test"))
			handler := glooinstance_handler.NewSingleClusterGlooInstanceHandler(mockGlooInstanceLister, *mockDiscoveryClient, mockEnvoyConfigClient)
			_, err := handler.GetConfigDumps(ctx, &rpc_edge_v1.GetConfigDumpsRequest{
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
