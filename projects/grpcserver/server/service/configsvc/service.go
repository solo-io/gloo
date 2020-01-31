package configsvc

import (
	"context"
	"time"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	"github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/pkg/license"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube"
	"go.uber.org/zap"
)

type configGrpcService struct {
	ctx             context.Context
	clientCache     client.ClientCache
	licenseClient   license.Client
	namespaceClient kube.NamespaceClient
	oAuthEndpoint   v1.OAuthEndpoint
	rawGetter       rawgetter.RawGetter
	version         setup.BuildVersion
	podNamespace    string
}

func NewConfigGrpcService(
	ctx context.Context,
	clientCache client.ClientCache,
	licenseClient license.Client,
	namespaceClient kube.NamespaceClient,
	rawGetter rawgetter.RawGetter,
	oAuthEndpoint v1.OAuthEndpoint,
	version setup.BuildVersion,
	podNamespace string) (v1.ConfigApiServer, error) {

	return &configGrpcService{
		ctx:             ctx,
		clientCache:     clientCache,
		licenseClient:   licenseClient,
		namespaceClient: namespaceClient,
		rawGetter:       rawGetter,
		oAuthEndpoint:   oAuthEndpoint,
		version:         version,
		podNamespace:    podNamespace,
	}, nil
}

func (s *configGrpcService) GetVersion(context.Context, *v1.GetVersionRequest) (*v1.GetVersionResponse, error) {
	return &v1.GetVersionResponse{Version: string(s.version)}, nil
}

func (s *configGrpcService) GetOAuthEndpoint(context.Context, *v1.GetOAuthEndpointRequest) (*v1.GetOAuthEndpointResponse, error) {
	return &v1.GetOAuthEndpointResponse{OAuthEndpoint: &s.oAuthEndpoint}, nil
}

func (s *configGrpcService) GetIsLicenseValid(context.Context, *v1.GetIsLicenseValidRequest) (*v1.GetIsLicenseValidResponse, error) {
	var wrappedError error
	isValid := true
	var invalidReason string
	if err := s.licenseClient.IsLicenseValid(); err != nil {
		wrappedError = LicenseIsInvalidError(err)
		contextutils.LoggerFrom(s.ctx).Debugw("license is not valid", zap.Error(wrappedError))
		isValid = false
	}
	if wrappedError != nil {
		invalidReason = wrappedError.Error()
	}

	return &v1.GetIsLicenseValidResponse{IsLicenseValid: isValid, InvalidReason: invalidReason}, nil
}

func (s *configGrpcService) GetSettings(ctx context.Context, request *v1.GetSettingsRequest) (*v1.GetSettingsResponse, error) {
	namespace := s.podNamespace
	name := defaults.SettingsName
	settings, err := s.clientCache.GetSettingsClient().Read(namespace, name, clients.ReadOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToReadSettingsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}

	return &v1.GetSettingsResponse{SettingsDetails: s.getDetails(settings)}, nil
}

func (s *configGrpcService) UpdateSettings(ctx context.Context, request *v1.UpdateSettingsRequest) (*v1.UpdateSettingsResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	if request.GetSettings().GetRefreshRate() != nil {
		if err := validateRefreshRate(request.GetSettings().GetRefreshRate()); err != nil {
			wrapped := FailedToUpdateSettingsError(err)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
	}

	written, err := s.clientCache.GetSettingsClient().Write(request.GetSettings(), clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	if err != nil {
		wrapped := FailedToUpdateSettingsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.UpdateSettingsResponse{SettingsDetails: s.getDetails(written)}, nil
}

func (s *configGrpcService) UpdateSettingsYaml(ctx context.Context, request *v1.UpdateSettingsYamlRequest) (*v1.UpdateSettingsResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}

	var (
		editedYaml  = request.GetEditedYamlData().GetEditedYaml()
		settingsRef = request.GetEditedYamlData().GetRef()
	)

	settingsFromYaml := &gloov1.Settings{}
	err := s.rawGetter.InitResourceFromYamlString(s.ctx, editedYaml, settingsRef, settingsFromYaml)

	if err != nil {
		wrapped := FailedToParseSettingsFromYamlError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	written, err := s.clientCache.GetSettingsClient().Write(settingsFromYaml, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: true})
	if err != nil {
		wrapped := FailedToUpdateSettingsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &v1.UpdateSettingsResponse{SettingsDetails: s.getDetails(written)}, nil
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

func (s *configGrpcService) getDetails(settings *gloov1.Settings) *v1.SettingsDetails {
	return &v1.SettingsDetails{
		Settings: settings,
		Raw:      s.rawGetter.GetRaw(s.ctx, settings, gloov1.SettingsCrd),
	}
}
