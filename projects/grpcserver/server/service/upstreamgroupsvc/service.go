package upstreamgroupsvc

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search"

	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"
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
}

// this client is not mocked by gloo, so mock it ourselves here
//go:generate mockgen -destination mocks/upstream_group_client_mock.go -package mocks github.com/solo-io/gloo/projects/gloo/pkg/api/v1 UpstreamGroupClient

func NewUpstreamGroupGrpcService(
	ctx context.Context,
	clientCache client.ClientCache,
	licenseClient license.Client,
	settingsValues settings.ValuesClient,
	rawGetter rawgetter.RawGetter,
	upstreamSearcher search.UpstreamSearcher,
) v1.UpstreamGroupApiServer {

	return &upstreamGrpcService{
		ctx:              ctx,
		clientCache:      clientCache,
		settingsValues:   settingsValues,
		rawGetter:        rawGetter,
		licenseClient:    licenseClient,
		upstreamSearcher: upstreamSearcher,
	}
}

func (s *upstreamGrpcService) GetUpstreamGroup(ctx context.Context, request *v1.GetUpstreamGroupRequest) (*v1.GetUpstreamGroupResponse, error) {
	upstream, err := s.readUpstreamGroup(request.GetRef())
	if err != nil {
		wrapped := FailedToReadUpstreamGroupError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.GetUpstreamGroupResponse{UpstreamGroupDetails: s.getDetails(upstream)}, nil
}

func (s *upstreamGrpcService) ListUpstreamGroups(ctx context.Context, request *v1.ListUpstreamGroupsRequest) (*v1.ListUpstreamGroupsResponse, error) {
	var upstreamList gloov1.UpstreamGroupList
	for _, ns := range s.settingsValues.GetWatchNamespaces() {
		upstreams, err := s.clientCache.GetUpstreamGroupClient().List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListUpstreamGroupsError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		upstreamList = append(upstreamList, upstreams...)
	}

	detailsList := make([]*v1.UpstreamGroupDetails, 0, len(upstreamList))
	for _, u := range upstreamList {
		detailsList = append(detailsList, s.getDetails(u))
	}

	return &v1.ListUpstreamGroupsResponse{UpstreamGroupDetails: detailsList}, nil
}

func (s *upstreamGrpcService) CreateUpstreamGroup(ctx context.Context, request *v1.CreateUpstreamGroupRequest) (*v1.CreateUpstreamGroupResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	written, err := s.writeUpstreamGroup(request.UpstreamGroup, false)
	if err != nil {
		wrapped := FailedToCreateUpstreamGroupError(err, request.UpstreamGroup.Metadata.Namespace, request.UpstreamGroup.Metadata.Name)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.CreateUpstreamGroupResponse{UpstreamGroupDetails: s.getDetails(written)}, nil
}

func (s *upstreamGrpcService) UpdateUpstreamGroup(ctx context.Context, request *v1.UpdateUpstreamGroupRequest) (*v1.UpdateUpstreamGroupResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	written, err := s.writeUpstreamGroup(request.UpstreamGroup, true)
	if err != nil {
		wrapped := FailedToUpdateUpstreamGroupError(err, request.UpstreamGroup.Metadata.Namespace, request.UpstreamGroup.Metadata.Name)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.UpdateUpstreamGroupResponse{UpstreamGroupDetails: s.getDetails(written)}, nil
}

func (s *upstreamGrpcService) DeleteUpstreamGroup(ctx context.Context, request *v1.DeleteUpstreamGroupRequest) (*v1.DeleteUpstreamGroupResponse, error) {
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
		return nil, s.wrapAndLogDeletionError(FailedToCheckIsUpstreamGroupReferencedError(err, upstreamRef), request)
	}

	if len(containingVirtualServiceRefs) > 0 {
		return nil, s.wrapAndLogDeletionError(CannotDeleteReferencedUpstreamGroupError(upstreamRef, containingVirtualServiceRefs), request)
	}

	err = s.clientCache.GetUpstreamGroupClient().Delete(namespace, name, clients.DeleteOpts{Ctx: s.ctx})
	if err != nil {
		return nil, s.wrapAndLogDeletionError(err, request)
	}
	return &v1.DeleteUpstreamGroupResponse{}, nil
}

func (s *upstreamGrpcService) wrapAndLogDeletionError(baseError error, request *v1.DeleteUpstreamGroupRequest) error {
	wrapped := FailedToDeleteUpstreamGroupError(baseError, request.GetRef())
	contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(baseError), zap.Any("request", request))
	return wrapped
}

func (s *upstreamGrpcService) readUpstreamGroup(ref *core.ResourceRef) (*gloov1.UpstreamGroup, error) {
	return s.clientCache.GetUpstreamGroupClient().Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: s.ctx})
}

func (s *upstreamGrpcService) writeUpstreamGroup(upstreamGroup *gloov1.UpstreamGroup, overwriteExisting bool) (*gloov1.UpstreamGroup, error) {
	return s.clientCache.GetUpstreamGroupClient().Write(upstreamGroup, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: overwriteExisting})
}

func (s *upstreamGrpcService) getDetails(upstreamGroup *gloov1.UpstreamGroup) *v1.UpstreamGroupDetails {
	return &v1.UpstreamGroupDetails{
		UpstreamGroup: upstreamGroup,
		Raw:           s.rawGetter.GetRaw(s.ctx, upstreamGroup, gloov1.UpstreamGroupCrd),
	}
}
