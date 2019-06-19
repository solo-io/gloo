package service

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type configGrpcService struct {
	settingsClient gloov1.SettingsClient
}

func NewConfigGrpcService(settingsClient gloov1.SettingsClient) v1.ConfigApiServer {
	return &configGrpcService{
		settingsClient: settingsClient,
	}
}

func (s *configGrpcService) GetVersion(context.Context, *v1.GetVersionRequest) (*v1.GetVersionResponse, error) {
	panic("implement me")
}

func (s *configGrpcService) GetOAuthEndpoint(context.Context, *v1.GetOAuthEndpointRequest) (*v1.GetOAuthEndpointResponse, error) {
	panic("implement me")
}

func (s *configGrpcService) GetIsLicenseValid(context.Context, *v1.GetIsLicenseValidRequest) (*v1.GetIsLicenseValidResponse, error) {
	panic("implement me")
}

func (s *configGrpcService) GetSettings(context.Context, *v1.GetSettingsRequest) (*v1.GetSettingsResponse, error) {
	panic("implement me")
}

func (s *configGrpcService) UpdateSettings(context.Context, *v1.UpdateSettingsRequest) (*v1.UpdateSettingsResponse, error) {
	panic("implement me")
}

func (s *configGrpcService) ListNamespaces(context.Context, *v1.ListNamespacesRequest) (*v1.ListNamespacesResponse, error) {
	panic("implement me")
}
