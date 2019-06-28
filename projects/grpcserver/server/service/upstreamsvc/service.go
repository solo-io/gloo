package upstreamsvc

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/converter"
	"go.uber.org/zap"
)

type upstreamGrpcService struct {
	ctx            context.Context
	upstreamClient gloov1.UpstreamClient
	settingsValues settings.ValuesClient
	inputConverter converter.UpstreamInputConverter
}

func NewUpstreamGrpcService(ctx context.Context, upstreamClient gloov1.UpstreamClient, inputConverter converter.UpstreamInputConverter, settingsValues settings.ValuesClient) v1.UpstreamApiServer {
	return &upstreamGrpcService{
		ctx:            ctx,
		upstreamClient: upstreamClient,
		inputConverter: inputConverter,
		settingsValues: settingsValues,
	}
}

func (s *upstreamGrpcService) GetUpstream(ctx context.Context, request *v1.GetUpstreamRequest) (*v1.GetUpstreamResponse, error) {
	upstream, err := s.readUpstream(request.GetRef())
	if err != nil {
		wrapped := FailedToReadUpstreamError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.GetUpstreamResponse{Upstream: upstream}, nil
}

func (s *upstreamGrpcService) ListUpstreams(ctx context.Context, request *v1.ListUpstreamsRequest) (*v1.ListUpstreamsResponse, error) {
	var upstreamList gloov1.UpstreamList
	for _, ns := range request.GetNamespaceList() {
		upstreams, err := s.upstreamClient.List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListUpstreamsError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		upstreamList = append(upstreamList, upstreams...)
	}

	return &v1.ListUpstreamsResponse{UpstreamList: upstreamList}, nil
}

func (s *upstreamGrpcService) StreamUpstreamList(request *v1.StreamUpstreamListRequest, stream v1.UpstreamApi_StreamUpstreamListServer) error {
	watch, errs, err := s.upstreamClient.Watch(request.GetNamespace(), clients.WatchOpts{
		RefreshRate: s.settingsValues.GetRefreshRate(),
		Ctx:         stream.Context(),
		Selector:    request.GetSelector(),
	})
	if err != nil {
		wrapped := FailedToStreamUpstreamsError(err, request.GetNamespace())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return wrapped
	}

	for {
		select {
		case list, ok := <-watch:
			if !ok {
				return nil
			}
			err := stream.Send(&v1.StreamUpstreamListResponse{UpstreamList: list})
			if err != nil {
				wrapped := ErrorWhileWatchingUpstreams(err, request.GetNamespace())
				contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
				return wrapped
			}
		case err, ok := <-errs:
			if !ok {
				return nil
			}
			wrapped := ErrorWhileWatchingUpstreams(err, request.GetNamespace())
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return wrapped
		case <-stream.Context().Done():
			return nil
		}
	}
}

func (s *upstreamGrpcService) CreateUpstream(ctx context.Context, request *v1.CreateUpstreamRequest) (*v1.CreateUpstreamResponse, error) {
	upstream := gloov1.Upstream{
		Metadata: core.Metadata{
			Namespace: request.GetInput().GetRef().GetNamespace(),
			Name:      request.GetInput().GetRef().GetName(),
		},
		UpstreamSpec: s.inputConverter.ConvertInputToUpstreamSpec(request.GetInput()),
	}

	written, err := s.writeUpstream(upstream, false)
	if err != nil {
		wrapped := FailedToCreateUpstreamError(err, request.GetInput().GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.CreateUpstreamResponse{Upstream: written}, nil
}

func (s *upstreamGrpcService) UpdateUpstream(ctx context.Context, request *v1.UpdateUpstreamRequest) (*v1.UpdateUpstreamResponse, error) {
	read, err := s.readUpstream(request.GetInput().GetRef())
	if err != nil {
		wrapped := FailedToUpdateUpstreamError(err, request.GetInput().GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	read.UpstreamSpec = s.inputConverter.ConvertInputToUpstreamSpec(request.GetInput())
	read.Status = core.Status{}

	written, err := s.writeUpstream(*read, true)
	if err != nil {
		wrapped := FailedToUpdateUpstreamError(err, request.GetInput().GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.UpdateUpstreamResponse{Upstream: written}, nil
}

func (s *upstreamGrpcService) DeleteUpstream(ctx context.Context, request *v1.DeleteUpstreamRequest) (*v1.DeleteUpstreamResponse, error) {
	err := s.upstreamClient.Delete(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.DeleteOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToDeleteUpstreamError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.DeleteUpstreamResponse{}, nil
}

func (s *upstreamGrpcService) readUpstream(ref *core.ResourceRef) (*gloov1.Upstream, error) {
	return s.upstreamClient.Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: s.ctx})
}

func (s *upstreamGrpcService) writeUpstream(upstream gloov1.Upstream, overwriteExisting bool) (*gloov1.Upstream, error) {
	return s.upstreamClient.Write(&upstream, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: overwriteExisting})
}
