package gatewaysvc

import (
	"context"

	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/gateway_client_mock.go -package mocks github.com/solo-io/gloo/projects/gateway/pkg/api/v2 GatewayClient

type gatewayGrpcService struct {
	ctx           context.Context
	gatewayClient gatewayv2.GatewayClient
	rawGetter     rawgetter.RawGetter
}

func NewGatewayGrpcService(ctx context.Context, gatewayClient gatewayv2.GatewayClient, rawGetter rawgetter.RawGetter) v1.GatewayApiServer {
	return &gatewayGrpcService{
		ctx:           ctx,
		gatewayClient: gatewayClient,
		rawGetter:     rawGetter,
	}
}

func (s *gatewayGrpcService) GetGateway(ctx context.Context, request *v1.GetGatewayRequest) (*v1.GetGatewayResponse, error) {
	gateway, err := s.gatewayClient.Read(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToGetGatewayError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.GetGatewayResponse{GatewayDetails: s.getDetails(gateway)}, nil
}

func (s *gatewayGrpcService) ListGateways(ctx context.Context, request *v1.ListGatewaysRequest) (*v1.ListGatewaysResponse, error) {
	var gatewayDetailsList []*v1.GatewayDetails
	for _, ns := range request.GetNamespaces() {
		gatewaysInNamespace, err := s.gatewayClient.List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListGatewaysError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, gw := range gatewaysInNamespace {
			gatewayDetailsList = append(gatewayDetailsList, s.getDetails(gw))
		}
	}
	return &v1.ListGatewaysResponse{GatewayDetails: gatewayDetailsList}, nil
}

func (s *gatewayGrpcService) getDetails(gateway *v2.Gateway) *v1.GatewayDetails {
	raw, err := s.rawGetter.GetRaw(gateway, gatewayv2.GatewayCrd)
	if err != nil {
		// Failed to generate yaml for resource -- not worth propagating
		contextutils.LoggerFrom(s.ctx).Errorw(err.Error(), zap.Error(err), zap.Any("gateway", gateway))
	}

	return &v1.GatewayDetails{
		Gateway: gateway,
		Raw:     raw,
	}
}
