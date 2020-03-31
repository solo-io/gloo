package envoydetails_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gwMocks "github.com/solo-io/gloo/projects/gateway/pkg/mocks"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	clientmocks "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	getter         envoydetails.EnvoyStatusGetter
	proxyClient    *mocks.MockProxyClient
	gatewayClient  *gwMocks.MockGatewayClient
	settingsClient *mocks.MockSettingsClient
	clientCache    *clientmocks.MockClientCache

	namespace = "ns"
	name      = "name"
	id        = "id"
)

var _ = Describe("EnvoyStatusGetter Test", func() {

	getPod := func(phase kubev1.PodPhase) kubev1.Pod {
		return kubev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      name,
				Labels:    map[string]string{envoydetails.GatewayProxyIdLabel: id},
			},
			Status: kubev1.PodStatus{Phase: phase},
		}
	}

	getProxy := func(state core.Status_State, reason string) *gloov1.Proxy {
		return &gloov1.Proxy{
			Status: core.Status{
				State:  state,
				Reason: reason,
			},
			Metadata: core.Metadata{
				Namespace: namespace,
				Name:      id,
			},
		}
	}

	getGateway := func(state core.Status_State, reason string) *gatewayv1.Gateway {
		return &gatewayv1.Gateway{
			Status: core.Status{
				State:  state,
				Reason: reason,
			},
			Metadata: core.Metadata{
				Namespace: namespace,
				Name:      id,
			},
		}
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		proxyClient = mocks.NewMockProxyClient(mockCtrl)
		gatewayClient = gwMocks.NewMockGatewayClient(mockCtrl)
		settingsClient = mocks.NewMockSettingsClient(mockCtrl)
		clientCache = clientmocks.NewMockClientCache(mockCtrl)
		clientCache.EXPECT().GetProxyClient().Return(proxyClient).AnyTimes()
		clientCache.EXPECT().GetGatewayClient().Return(gatewayClient).AnyTimes()
		clientCache.EXPECT().GetSettingsClient().Return(settingsClient).AnyTimes()
		settingsClient.EXPECT().Read(gomock.Any(), defaults.SettingsName, gomock.Any()).
			Return(&gloov1.Settings{Gateway: &gloov1.GatewayOptions{ReadGatewaysFromAllNamespaces: false}}, nil).AnyTimes()
		getter = envoydetails.NewEnvoyStatusGetter(clientCache)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetEnvoyStatus", func() {
		It("reports OK when pod is running and proxy is accepted", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Accepted, ""), nil)
			gatewayClient.EXPECT().
				List(namespace, clients.ListOpts{Ctx: context.Background()}).
				Return([]*gatewayv1.Gateway{getGateway(core.Status_Accepted, "")}, nil)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code: v1.Status_OK,
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when pod is in Failed phase", func() {
			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodFailed))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: envoydetails.GatewayProxyPodIsNotRunning(namespace, name, string(kubev1.PodFailed)),
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when proxy cannot be found", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(nil, testErr)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: envoydetails.ProxyResourceNotFound(id),
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports warning when proxy is pending", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Pending, ""), nil)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_WARNING,
				Message: envoydetails.ProxyResourcePending(namespace, id),
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when proxy is not accepted or pending", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Rejected, "test-reason"), nil)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: envoydetails.ProxyResourceRejected(namespace, id, "test-reason"),
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when gateways cannot be found", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Accepted, ""), nil)
			gatewayClient.EXPECT().
				List(namespace, clients.ListOpts{Ctx: context.Background()}).
				Return(nil, testErr)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: envoydetails.GatewayResourcesNotFound(namespace),
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when gateways report warnings", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Accepted, ""), nil)
			gatewayClient.EXPECT().
				List(namespace, clients.ListOpts{Ctx: context.Background()}).
				Return([]*gatewayv1.Gateway{getGateway(core.Status_Warning, "test-reason")}, nil)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: "Gateways are in a bad state",
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when gateways report both errors and warnings", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Accepted, ""), nil)
			gatewayClient.EXPECT().
				List(namespace, clients.ListOpts{Ctx: context.Background()}).
				Return([]*gatewayv1.Gateway{
					getGateway(core.Status_Rejected, "test-reason"),
					getGateway(core.Status_Rejected, "test-reason"),
					getGateway(core.Status_Warning, "test-reason"),
					getGateway(core.Status_Accepted, "test-reason")}, nil)

			actual := getter.GetEnvoyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: "Gateways are in a bad state",
			}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

})
