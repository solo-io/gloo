package configsvc

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/pkg/license"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube"
	"go.uber.org/zap"
)

type BuildVersion string

type configGrpcService struct {
	ctx             context.Context
	settingsClient  gloov1.SettingsClient
	licenseClient   license.Client
	namespaceClient kube.NamespaceClient
	oAuthEndpoint   v1.OAuthEndpoint
	version         BuildVersion
	podNamespace    string
}

func NewConfigGrpcService(
	ctx context.Context,
	settingsClient gloov1.SettingsClient,
	licenseClient license.Client,
	namespaceClient kube.NamespaceClient,
	oAuthEndpoint v1.OAuthEndpoint,
	version BuildVersion,
	podNamespace string) v1.ConfigApiServer {

	return &configGrpcService{
		ctx:             ctx,
		settingsClient:  settingsClient,
		licenseClient:   licenseClient,
		namespaceClient: namespaceClient,
		oAuthEndpoint:   oAuthEndpoint,
		version:         version,
		podNamespace:    podNamespace,
	}
}

func (s *configGrpcService) GetVersion(context.Context, *v1.GetVersionRequest) (*v1.GetVersionResponse, error) {
	return &v1.GetVersionResponse{Version: string(s.version)}, nil
}

func (s *configGrpcService) GetOAuthEndpoint(context.Context, *v1.GetOAuthEndpointRequest) (*v1.GetOAuthEndpointResponse, error) {
	return &v1.GetOAuthEndpointResponse{OAuthEndpoint: &s.oAuthEndpoint}, nil
}

func (s *configGrpcService) GetIsLicenseValid(context.Context, *v1.GetIsLicenseValidRequest) (*v1.GetIsLicenseValidResponse, error) {
	if err := s.licenseClient.IsLicenseValid(); err != nil {
		wrapped := LicenseIsInvalidError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}

	return &v1.GetIsLicenseValidResponse{IsLicenseValid: true}, nil
}

func (s *configGrpcService) GetSettings(ctx context.Context, request *v1.GetSettingsRequest) (*v1.GetSettingsResponse, error) {
	namespace := s.podNamespace
	name := defaults.SettingsName
	settings, err := s.settingsClient.Read(namespace, name, clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToReadSettingsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}

	return &v1.GetSettingsResponse{Settings: settings}, nil
}

func (s *configGrpcService) UpdateSettings(ctx context.Context, request *v1.UpdateSettingsRequest) (*v1.UpdateSettingsResponse, error) {
	existing, err := s.settingsClient.Read(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToUpdateSettingsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	if request.GetRefreshRate() != nil {
		if err := validateRefreshRate(request.GetRefreshRate()); err != nil {
			wrapped := FailedToUpdateSettingsError(err)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		existing.RefreshRate = request.GetRefreshRate()
	}
	existing.WatchNamespaces = request.GetWatchNamespaces()

	existing.Status = core.Status{}
	written, err := s.settingsClient.Write(existing, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	if err != nil {
		wrapped := FailedToUpdateSettingsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.UpdateSettingsResponse{Settings: written}, nil
}

func (s *configGrpcService) ListNamespaces(ctx context.Context, request *v1.ListNamespacesRequest) (*v1.ListNamespacesResponse, error) {
	namespaceList, err := s.namespaceClient.ListNamespaces()
	if err != nil {
		wrapped := FailedToListNamespacesError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.ListNamespacesResponse{Namespaces: namespaceList}, nil
}

func (s *configGrpcService) GetPodNamespace(context.Context, *v1.GetPodNamespaceRequest) (*v1.GetPodNamespaceResponse, error) {
	return &v1.GetPodNamespaceResponse{Namespace: s.podNamespace}, nil
}

func validateRefreshRate(rr *types.Duration) error {
	duration, err := types.DurationFromProto(rr)
	if err != nil {
		return err
	}

	if duration < time.Second {
		return InvalidRefreshRateError(duration)
	}
	return nil
}
