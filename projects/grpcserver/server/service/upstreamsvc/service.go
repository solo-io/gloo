package upstreamsvc

import (
	"context"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation"
	"go.uber.org/zap"
)

type upstreamGrpcService struct {
	ctx             context.Context
	upstreamClient  gloov1.UpstreamClient
	licenseClient   license.Client
	settingsValues  settings.ValuesClient
	mutator         mutation.Mutator
	mutationFactory mutation.Factory
	rawGetter       rawgetter.RawGetter
}

func NewUpstreamGrpcService(
	ctx context.Context,
	upstreamClient gloov1.UpstreamClient,
	licenseClient license.Client,
	settingsValues settings.ValuesClient,
	mutator mutation.Mutator,
	factory mutation.Factory,
	rawGetter rawgetter.RawGetter) v1.UpstreamApiServer {

	return &upstreamGrpcService{
		ctx:             ctx,
		upstreamClient:  upstreamClient,
		settingsValues:  settingsValues,
		mutator:         mutator,
		mutationFactory: factory,
		rawGetter:       rawGetter,
		licenseClient:   licenseClient,
	}
}

func (s *upstreamGrpcService) GetUpstream(ctx context.Context, request *v1.GetUpstreamRequest) (*v1.GetUpstreamResponse, error) {
	upstream, err := s.readUpstream(request.GetRef())
	if err != nil {
		wrapped := FailedToReadUpstreamError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.GetUpstreamResponse{Upstream: upstream, UpstreamDetails: s.getDetails(upstream)}, nil
}

func (s *upstreamGrpcService) ListUpstreams(ctx context.Context, request *v1.ListUpstreamsRequest) (*v1.ListUpstreamsResponse, error) {
	var upstreamList gloov1.UpstreamList
	for _, ns := range request.GetNamespaces() {
		upstreams, err := s.upstreamClient.List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListUpstreamsError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		upstreamList = append(upstreamList, upstreams...)
	}

	detailsList := make([]*v1.UpstreamDetails, 0, len(upstreamList))
	for _, u := range upstreamList {
		detailsList = append(detailsList, s.getDetails(u))
	}

	return &v1.ListUpstreamsResponse{Upstreams: upstreamList, UpstreamDetails: detailsList}, nil
}

func (s *upstreamGrpcService) CreateUpstream(ctx context.Context, request *v1.CreateUpstreamRequest) (*v1.CreateUpstreamResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	var (
		written *gloov1.Upstream
		err     error
	)
	if request.GetUpstreamInput() == nil {
		written, err = s.mutator.Create(s.ctx, request.GetInput().GetRef(), s.mutationFactory.ConfigureUpstream(request.GetInput()))
	} else {
		written, err = s.mutator.CreateUpstream(s.ctx, request.GetUpstreamInput())
	}

	if err != nil {
		wrapped := FailedToCreateUpstreamError(err, request.GetInput().GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.CreateUpstreamResponse{Upstream: written, UpstreamDetails: s.getDetails(written)}, nil
}

func (s *upstreamGrpcService) UpdateUpstream(ctx context.Context, request *v1.UpdateUpstreamRequest) (*v1.UpdateUpstreamResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	var (
		written *gloov1.Upstream
		err     error
	)

	if request.GetUpstreamInput() == nil {
		written, err = s.mutator.Update(s.ctx, request.GetInput().GetRef(), s.mutationFactory.ConfigureUpstream(request.GetInput()))
	} else {
		written, err = s.mutator.UpdateUpstream(s.ctx, request.GetUpstreamInput())
	}

	if err != nil {
		wrapped := FailedToUpdateUpstreamError(err, request.GetInput().GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.UpdateUpstreamResponse{Upstream: written, UpstreamDetails: s.getDetails(written)}, nil
}

func (s *upstreamGrpcService) DeleteUpstream(ctx context.Context, request *v1.DeleteUpstreamRequest) (*v1.DeleteUpstreamResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
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

func (s *upstreamGrpcService) writeUpstream(upstream *gloov1.Upstream, overwriteExisting bool) (*gloov1.Upstream, error) {
	return s.upstreamClient.Write(upstream, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: overwriteExisting})
}

func (s *upstreamGrpcService) getDetails(upstream *gloov1.Upstream) *v1.UpstreamDetails {
	return &v1.UpstreamDetails{
		Upstream: upstream,
		Raw:      s.rawGetter.GetRaw(s.ctx, upstream, gloov1.UpstreamCrd),
	}
}
