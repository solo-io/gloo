package gatewaysvc_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	. "github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	mock_rawgetter "github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter/mocks"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/gatewaysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/gatewaysvc/mocks"
	. "github.com/solo-io/solo-projects/projects/grpcserver/server/service/internal/testutils"
	"google.golang.org/grpc"
)

var (
	grpcServer    *grpc.Server
	conn          *grpc.ClientConn
	apiserver     v1.GatewayApiServer
	client        v1.GatewayApiClient
	mockCtrl      *gomock.Controller
	gatewayClient *mocks.MockGatewayClient
	rawGetter     *mock_rawgetter.MockRawGetter
	testErr       = errors.Errorf("test-err")
)

var _ = Describe("ServiceTest", func() {

	getRaw := func(gateway *gatewayv2.Gateway) *v1.Raw {
		return &v1.Raw{FileName: gateway.GetMetadata().Name}
	}

	getGatewayDetails := func(gateway *gatewayv2.Gateway) *v1.GatewayDetails {
		return &v1.GatewayDetails{
			Gateway: gateway,
			Raw:     getRaw(gateway),
		}
	}

	getGatewayDetailsList := func(gateways ...*gatewayv2.Gateway) []*v1.GatewayDetails {
		list := make([]*v1.GatewayDetails, 0, len(gateways))
		for _, g := range gateways {
			list = append(list, &v1.GatewayDetails{
				Gateway: g,
				Raw:     getRaw(g),
			})
		}
		return list
	}

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		gatewayClient = mocks.NewMockGatewayClient(mockCtrl)
		rawGetter = mock_rawgetter.NewMockRawGetter(mockCtrl)
		apiserver = gatewaysvc.NewGatewayGrpcService(context.TODO(), gatewayClient, rawGetter)

		grpcServer, conn = MustRunGrpcServer(func(s *grpc.Server) { v1.RegisterGatewayApiServer(s, apiserver) })
		client = v1.NewGatewayApiClient(conn)
	})

	AfterEach(func() {
		grpcServer.Stop()
		mockCtrl.Finish()
	})

	Describe("GetGateway", func() {
		It("works when the gateway client works", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()
			gateway := &gatewayv2.Gateway{
				Metadata: metadata,
			}

			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(gateway, nil)
			rawGetter.EXPECT().
				GetRaw(gateway, gatewayv2.GatewayCrd).
				Return(getRaw(gateway), nil)

			request := &v1.GetGatewayRequest{Ref: &ref}
			actual, err := client.GetGateway(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.GetGatewayResponse{GatewayDetails: getGatewayDetails(gateway)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the gateway client errors", func() {
			metadata := core.Metadata{
				Namespace: "ns",
				Name:      "name",
			}
			ref := metadata.Ref()

			gatewayClient.EXPECT().
				Read(metadata.Namespace, metadata.Name, clients.ReadOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.GetGatewayRequest{Ref: &ref}
			_, err := client.GetGateway(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := gatewaysvc.FailedToGetGatewayError(testErr, &ref)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})

	Describe("ListGateways", func() {
		It("works when the gateway client works", func() {
			ns1, ns2 := "one", "two"
			gateway1 := &gatewayv2.Gateway{
				Metadata: core.Metadata{Namespace: ns1},
			}
			gateway2 := &gatewayv2.Gateway{
				Metadata: core.Metadata{Namespace: ns2},
			}

			gatewayClient.EXPECT().
				List(ns1, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv2.Gateway{gateway1}, nil)
			rawGetter.EXPECT().
				GetRaw(gateway1, gatewayv2.GatewayCrd).
				Return(getRaw(gateway1), nil)
			gatewayClient.EXPECT().
				List(ns2, clients.ListOpts{Ctx: context.TODO()}).
				Return([]*gatewayv2.Gateway{gateway2}, nil)
			rawGetter.EXPECT().
				GetRaw(gateway2, gatewayv2.GatewayCrd).
				Return(getRaw(gateway2), nil)

			request := &v1.ListGatewaysRequest{Namespaces: []string{ns1, ns2}}
			actual, err := client.ListGateways(context.TODO(), request)
			Expect(err).NotTo(HaveOccurred())
			expected := &v1.ListGatewaysResponse{GatewayDetails: getGatewayDetailsList(gateway1, gateway2)}
			ExpectEqualProtoMessages(actual, expected)
		})

		It("errors when the gateway client errors", func() {
			ns := "ns"

			gatewayClient.EXPECT().
				List(ns, clients.ListOpts{Ctx: context.TODO()}).
				Return(nil, testErr)

			request := &v1.ListGatewaysRequest{Namespaces: []string{ns}}
			_, err := client.ListGateways(context.TODO(), request)
			Expect(err).To(HaveOccurred())
			expectedErr := gatewaysvc.FailedToListGatewaysError(testErr, ns)
			Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
		})
	})
})
