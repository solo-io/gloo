package plugin

import (
	"context"
	"fmt"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

// This plugin will authorize any requests that include a given header, regardless of its value.
type IsHeaderPresentPlugin struct{}

type Config struct {
	RequiredHeader string
}

func (p *IsHeaderPresentPlugin) NewConfigInstance(ctx context.Context) (interface{}, error) {
	logger(ctx).Infow("Called 'NewConfigInstance' on IsHeaderPresentPlugin")
	return &Config{}, nil
}

func (p *IsHeaderPresentPlugin) GetAuthService(ctx context.Context, configInstance interface{}) (api.AuthService, error) {
	config, ok := configInstance.(*Config)
	if !ok {
		return nil, errors.New(fmt.Sprintf("unexpected config type %T", configInstance))
	}
	logger(ctx).Infow("Returning IsHeaderPresentAuthService instance", zap.Any("requiredHeader", config.RequiredHeader))
	return &IsHeaderPresentAuthService{RequiredHeader: config.RequiredHeader}, nil
}

type IsHeaderPresentAuthService struct {
	RequiredHeader string
}

func (c *IsHeaderPresentAuthService) Start(ctx context.Context) error {
	logger(ctx).Infow("Called 'Start' on IsHeaderPresentAuthService")
	return nil
}

func (c *IsHeaderPresentAuthService) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {
	for key, value := range request.CheckRequest.Attributes.Request.Http.Headers {
		if key == c.RequiredHeader {
			logger(ctx).Infow("Found required header", "header", key, "value", value)
			response := api.AuthorizedResponse()

			// Append extra header
			response.CheckResponse.HttpResponse = &envoy_service_auth_v3.CheckResponse_OkResponse{
				OkResponse: &envoy_service_auth_v3.OkHttpResponse{
					Headers: []*envoy_config_core_v3.HeaderValueOption{{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "auth-header-found",
							Value: "true",
						},
					}},
				},
			}

			return response, nil
		}
	}
	logger(ctx).Infow("Required header not found, denying access")
	return api.UnauthorizedResponse(), nil
}

func logger(ctx context.Context) *zap.SugaredLogger {
	return contextutils.LoggerFrom(contextutils.WithLogger(ctx, "required_header_plugin"))
}
