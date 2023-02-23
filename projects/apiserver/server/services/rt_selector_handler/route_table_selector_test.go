package rt_selector_handler_test

import (
	"context"

	gateway_solo_io_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	mock_gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1/mocks"
	mock_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1/mocks"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/rt_selector_handler"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("route table selector", func() {
	var (
		ctrl               *gomock.Controller
		mcGatewayCRDClient *mock_v1.MockMulticlusterClientset
		gatewayClientSet   *mock_v1.MockClientset
		rtClient           *mock_gateway_v1.MockRouteTableClient
	)

	BeforeEach(func() {
		ctrl, _ = gomock.WithContext(context.TODO(), GinkgoT())

		mcGatewayCRDClient = mock_v1.NewMockMulticlusterClientset(ctrl)
		gatewayClientSet = mock_v1.NewMockClientset(ctrl)
		rtClient = mock_gateway_v1.NewMockRouteTableClient(ctrl)

		mcGatewayCRDClient.EXPECT().Cluster(gomock.Any()).Return(gatewayClientSet, nil).AnyTimes()
		gatewayClientSet.EXPECT().RouteTables().Return(rtClient).AnyTimes()

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("can SelectRouteTables with a delegate action by ref", func() {
		defaultRt1 := &gateway_solo_io_v1.RouteTable{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "test-route-table",
				Namespace: "not-gloo-system",
			},
		}
		rtClient.EXPECT().GetRouteTable(context.Background(), client.ObjectKey{
			Name:      "test-route-table",
			Namespace: "not-gloo-system",
		}).Return(defaultRt1, nil)
		selector := rt_selector_handler.NewRouteTableSelector(context.Background(), rtClient)
		routeTables, err := selector.SelectRouteTables(&gatewayv1.DelegateAction{
			DelegationType: &gatewayv1.DelegateAction_Ref{
				Ref: &core.ResourceRef{
					Name:      "test-route-table",
					Namespace: "not-gloo-system",
				},
			},
		}, "gloo-system")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(routeTables.Items)).To(Equal(1))
	})

	It("can SelectRouteTables with a delegate action with a namespace selector and labels", func() {
		correctLabels := map[string]string{"test": "labels"}
		defaultRt1 := gateway_solo_io_v1.RouteTable{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "route-table-1",
				Namespace: "default",
				Labels:    correctLabels,
			},
		}
		defaultRt2 := gateway_solo_io_v1.RouteTable{
			ObjectMeta: meta_v1.ObjectMeta{
				Name:      "route-table-1",
				Namespace: "wrong",
				Labels:    correctLabels,
			},
		}
		rtClient.EXPECT().ListRouteTable(gomock.Any(), client.MatchingLabels(correctLabels)).
			Return(&gateway_solo_io_v1.RouteTableList{
				Items: []gateway_solo_io_v1.RouteTable{
					defaultRt1,
					defaultRt2,
				},
			}, nil)
		selector := rt_selector_handler.NewRouteTableSelector(context.Background(), rtClient)
		routeTables, err := selector.SelectRouteTables(&gatewayv1.DelegateAction{
			DelegationType: &gatewayv1.DelegateAction_Selector{
				Selector: &gatewayv1.RouteTableSelector{
					Namespaces: []string{"default", "gloo-system"},
					Labels:     correctLabels,
				},
			},
		}, "gloo-system")

		Expect(err).NotTo(HaveOccurred())
		Expect(len(routeTables.Items)).To(Equal(1))
	})

})
