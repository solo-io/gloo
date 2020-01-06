package upstreamsvc

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/truncate"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"go.uber.org/zap"
)

type upstreamGrpcService struct {
	ctx              context.Context
	clientCache      client.ClientCache
	licenseClient    license.Client
	settingsValues   settings.ValuesClient
	rawGetter        rawgetter.RawGetter
	upstreamSearcher search.UpstreamSearcher
	truncator        truncate.UpstreamTruncator
}

// this client is not mocked by gloo, so mock it ourselves here
//go:generate mockgen -destination mocks/upstream_group_client_mock.go -self_package github.com/solo-io/gloo/projects/gloo/pkg/api/v1 -package mocks github.com/solo-io/gloo/projects/gloo/pkg/api/v1 UpstreamGroupClient

func NewUpstreamGrpcService(
	ctx context.Context,
	clientCache client.ClientCache,
	licenseClient license.Client,
	settingsValues settings.ValuesClient,
	rawGetter rawgetter.RawGetter,
	upstreamSearcher search.UpstreamSearcher,
	truncator truncate.UpstreamTruncator,
) v1.UpstreamApiServer {

	return &upstreamGrpcService{
		ctx:              ctx,
		clientCache:      clientCache,
		settingsValues:   settingsValues,
		rawGetter:        rawGetter,
		licenseClient:    licenseClient,
		upstreamSearcher: upstreamSearcher,
		truncator:        truncator,
	}
}

func (s *upstreamGrpcService) GetUpstream(ctx context.Context, request *v1.GetUpstreamRequest) (*v1.GetUpstreamResponse, error) {
	upstream, err := s.readUpstream(request.GetRef())
	if err != nil {
		wrapped := FailedToReadUpstreamError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.GetUpstreamResponse{UpstreamDetails: s.getDetails(upstream)}, nil
}

func (s *upstreamGrpcService) ListUpstreams(ctx context.Context, request *v1.ListUpstreamsRequest) (*v1.ListUpstreamsResponse, error) {
	var upstreamList gloov1.UpstreamList
	for _, ns := range s.settingsValues.GetWatchNamespaces() {
		upstreams, err := s.clientCache.GetUpstreamClient().List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListUpstreamsError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		upstreamList = append(upstreamList, upstreams...)
	}

	detailsList := make([]*v1.UpstreamDetails, 0, len(upstreamList))
	for _, u := range upstreamList {
		s.truncator.Truncate(u)
		detailsList = append(detailsList, s.getDetails(u))
	}

	return &v1.ListUpstreamsResponse{UpstreamDetails: detailsList}, nil
}

func (s *upstreamGrpcService) CreateUpstream(ctx context.Context, request *v1.CreateUpstreamRequest) (*v1.CreateUpstreamResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	written, err := s.clientCache.GetUpstreamClient().Write(request.GetUpstreamInput(), clients.WriteOpts{Ctx: s.ctx})
	if err != nil {
		ref := request.GetUpstreamInput().GetMetadata().Ref()
		wrapped := FailedToCreateUpstreamError(err, &ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.CreateUpstreamResponse{UpstreamDetails: s.getDetails(written)}, nil
}

func (s *upstreamGrpcService) UpdateUpstream(ctx context.Context, request *v1.UpdateUpstreamRequest) (*v1.UpdateUpstreamResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	written, err := s.clientCache.GetUpstreamClient().Write(request.GetUpstreamInput(), clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	if err != nil {
		ref := request.GetUpstreamInput().GetMetadata().Ref()
		wrapped := FailedToUpdateUpstreamError(err, &ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.UpdateUpstreamResponse{UpstreamDetails: s.getDetails(written)}, nil
}

func (s *upstreamGrpcService) DeleteUpstream(ctx context.Context, request *v1.DeleteUpstreamRequest) (*v1.DeleteUpstreamResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	var (
		namespace   = request.GetRef().GetNamespace()
		name        = request.GetRef().GetName()
		upstreamRef = request.GetRef()
	)

	// make sure we aren't trying to delete an upstream that's referenced in some virtual service
	containingVirtualServiceRefs, err := s.upstreamSearcher.FindContainingVirtualServices(s.ctx, upstreamRef)

	if err != nil {
		return nil, s.wrapAndLogDeletionError(FailedToCheckIsUpstreamReferencedError(err, upstreamRef), request)
	}

	if len(containingVirtualServiceRefs) > 0 {
		return nil, s.wrapAndLogDeletionError(CannotDeleteReferencedUpstreamError(upstreamRef, containingVirtualServiceRefs), request)
	}

	err = s.clientCache.GetUpstreamClient().Delete(namespace, name, clients.DeleteOpts{Ctx: s.ctx})
	if err != nil {
		return nil, s.wrapAndLogDeletionError(err, request)
	}
	return &v1.DeleteUpstreamResponse{}, nil
}

func (s *upstreamGrpcService) wrapAndLogDeletionError(baseError error, request *v1.DeleteUpstreamRequest) error {
	wrapped := FailedToDeleteUpstreamError(baseError, request.GetRef())
	contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(baseError), zap.Any("request", request))
	return wrapped
}

func (s *upstreamGrpcService) readUpstream(ref *core.ResourceRef) (*gloov1.Upstream, error) {
	return s.clientCache.GetUpstreamClient().Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: s.ctx})
}

func (s *upstreamGrpcService) writeUpstream(upstream *gloov1.Upstream, overwriteExisting bool) (*gloov1.Upstream, error) {
	return s.clientCache.GetUpstreamClient().Write(upstream, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: overwriteExisting})
}

func (s *upstreamGrpcService) getDetails(upstream *gloov1.Upstream) *v1.UpstreamDetails {
	return &v1.UpstreamDetails{
		Upstream: upstream,
		Raw:      s.rawGetter.GetRaw(s.ctx, upstream, gloov1.UpstreamCrd),
	}
}
