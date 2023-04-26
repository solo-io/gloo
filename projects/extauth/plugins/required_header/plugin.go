package main

import (
	"context"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/solo-io/ext-auth-plugins/api"
)

func main() {}

// Compile-time assertion
var _ api.ExtAuthPlugin = &mockPlugin{}

// This is the exported symbol that GlooE will look for.
//
//goland:noinspection GoUnusedGlobalVariable
var Plugin mockPlugin

type mockPlugin struct{}

func (m *mockPlugin) NewConfigInstance(ctx context.Context) (configInstance interface{}, err error) {
	return &structpb.Struct{}, nil
}

func (m *mockPlugin) GetAuthService(ctx context.Context, configInstance interface{}) (api.AuthService, error) {
	return &mockAuthService{}, nil
}

type mockAuthService struct{}

func (m mockAuthService) Start(ctx context.Context) error {
	return nil
}

func (m mockAuthService) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {
	if request.CheckRequest.GetAttributes().GetRequest().GetHttp().GetHeaders() != nil {
		return api.AuthorizedResponse(), nil
	}
	return api.UnauthorizedResponse(), nil
}

var _ api.ExtAuthPlugin = &mockPlugin{}
