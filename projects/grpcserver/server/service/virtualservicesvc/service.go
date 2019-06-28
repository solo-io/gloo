package virtualservicesvc

import (
	"context"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type virtualServiceGrpcService struct {
	virtualServiceClient gatewayv1.VirtualServiceClient
}

func NewVirtualServiceGrpcService(virtualServiceClient gatewayv1.VirtualServiceClient) v1.VirtualServiceApiServer {
	return &virtualServiceGrpcService{
		virtualServiceClient: virtualServiceClient,
	}
}

func (s *virtualServiceGrpcService) GetVirtualService(context.Context, *v1.GetVirtualServiceRequest) (*v1.GetVirtualServiceResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) ListVirtualServices(context.Context, *v1.ListVirtualServicesRequest) (*v1.ListVirtualServicesResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) StreamVirtualServiceList(*v1.StreamVirtualServiceListRequest, v1.VirtualServiceApi_StreamVirtualServiceListServer) error {
	panic("implement me")
}

func (s *virtualServiceGrpcService) CreateVirtualService(context.Context, *v1.CreateVirtualServiceRequest) (*v1.CreateVirtualServiceResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) UpdateVirtualService(context.Context, *v1.UpdateVirtualServiceRequest) (*v1.UpdateVirtualServiceResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) DeleteVirtualService(context.Context, *v1.DeleteVirtualServiceRequest) (*v1.DeleteVirtualServiceResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) CreateRoute(context.Context, *v1.CreateRouteRequest) (*v1.CreateRouteResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) UpdateRoute(context.Context, *v1.UpdateRouteRequest) (*v1.UpdateRouteResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) DeleteRoute(context.Context, *v1.DeleteRouteRequest) (*v1.DeleteRouteResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) SwapRoutes(context.Context, *v1.SwapRoutesRequest) (*v1.SwapRoutesResponse, error) {
	panic("implement me")
}

func (s *virtualServiceGrpcService) ShiftRoutes(context.Context, *v1.ShiftRoutesRequest) (*v1.ShiftRoutesResponse, error) {
	panic("implement me")
}
