package service

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type upstreamGrpcService struct {
	upstreamClient gloov1.UpstreamClient
}

func NewUpstreamGrpcService(upstreamClient gloov1.UpstreamClient) v1.UpstreamApiServer {
	return &upstreamGrpcService{
		upstreamClient: upstreamClient,
	}
}

func (s *upstreamGrpcService) GetUpstream(context.Context, *v1.GetUpstreamRequest) (*v1.GetUpstreamResponse, error) {
	panic("implement me")
}

func (s *upstreamGrpcService) ListUpstreams(context.Context, *v1.ListUpstreamsRequest) (*v1.ListUpstreamsResponse, error) {
	panic("implement me")
}

func (s *upstreamGrpcService) StreamUpstreamList(context.Context, *v1.StreamUpstreamListRequest) (*v1.StreamUpstreamListResponse, error) {
	panic("implement me")
}

func (s *upstreamGrpcService) CreateUpstream(context.Context, *v1.CreateUpstreamRequest) (*v1.CreateUpstreamResponse, error) {
	panic("implement me")
}

func (s *upstreamGrpcService) UpdateUpstream(context.Context, *v1.UpdateUpstreamRequest) (*v1.UpdateUpstreamResponse, error) {
	panic("implement me")
}

func (s *upstreamGrpcService) DeleteUpstream(context.Context, *v1.DeleteUpstreamRequest) (*v1.DeleteUpstreamResponse, error) {
	panic("implement me")
}
