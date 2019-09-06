package envoydetails_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/mocks"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails"
	kubev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	getter      envoydetails.ProxyStatusGetter
	proxyClient *mocks.MockProxyClient

	namespace = "ns"
	name      = "name"
	id        = "id"
)

var _ = Describe("ProxyStatusGetter Test", func() {

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

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		proxyClient = mocks.NewMockProxyClient(mockCtrl)
		getter = envoydetails.NewProxyStatusGetter(proxyClient)
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("GetProxyStatus", func() {
		It("reports OK when pod is running and proxy is accepted", func() {
			proxyClient.EXPECT().
				Read(namespace, id, clients.ReadOpts{Ctx: context.Background()}).
				Return(getProxy(core.Status_Accepted, ""), nil)

			actual := getter.GetProxyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code: v1.Status_OK,
			}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("reports error when pod is in Failed phase", func() {
			actual := getter.GetProxyStatus(context.Background(), getPod(kubev1.PodFailed))
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

			actual := getter.GetProxyStatus(context.Background(), getPod(kubev1.PodRunning))
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

			actual := getter.GetProxyStatus(context.Background(), getPod(kubev1.PodRunning))
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

			actual := getter.GetProxyStatus(context.Background(), getPod(kubev1.PodRunning))
			expected := &v1.Status{
				Code:    v1.Status_ERROR,
				Message: envoydetails.ProxyResourceRejected(namespace, id, "test-reason"),
			}
			ExpectEqualProtoMessages(actual, expected)
		})
	})

})
