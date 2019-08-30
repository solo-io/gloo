package gatewaysvc

import (
	"context"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status"

	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/gateway_client_mock.go -package mocks github.com/solo-io/gloo/projects/gateway/pkg/api/v2 GatewayClient

type gatewayGrpcService struct {
	ctx             context.Context
	gatewayClient   gatewayv2.GatewayClient
	rawGetter       rawgetter.RawGetter
	statusConverter status.InputResourceStatusGetter
	licenseClient   license.Client
}

func NewGatewayGrpcService(
	ctx context.Context,
	gatewayClient gatewayv2.GatewayClient,
	rawGetter rawgetter.RawGetter,
	statusConverter status.InputResourceStatusGetter,
	licenseClient license.Client,
) v1.GatewayApiServer {
	return &gatewayGrpcService{
		ctx:             ctx,
		gatewayClient:   gatewayClient,
		rawGetter:       rawGetter,
		statusConverter: statusConverter,
		licenseClient:   licenseClient,
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

func (s *gatewayGrpcService) UpdateGateway(ctx context.Context, request *v1.UpdateGatewayRequest) (*v1.UpdateGatewayResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	gateway := request.GetGateway()
	ref := gateway.GetMetadata().Ref()
	read, err := s.gatewayClient.Read(ref.Namespace, ref.Name, clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToUpdateGatewayError(err, &ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	gateway.Metadata.ResourceVersion = read.Metadata.ResourceVersion
	gateway.Status = core.Status{}
	written, err := s.gatewayClient.Write(gateway, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	if err != nil {
		wrapped := FailedToUpdateGatewayError(err, &ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.UpdateGatewayResponse{GatewayDetails: s.getDetails(written)}, nil
}

func (s *gatewayGrpcService) getDetails(gateway *v2.Gateway) *v1.GatewayDetails {
	return &v1.GatewayDetails{
		Gateway: gateway,
		Raw:     s.rawGetter.GetRaw(s.ctx, gateway, gatewayv2.GatewayCrd),
		Status:  s.statusConverter.GetApiStatusFromResource(gateway),
	}
}
