package plugin

import (
	"context"
	"fmt"

	envoycorev2 "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errors"
	"go.uber.org/zap"
)

type HeaderValuePlugin struct{}

type Config struct {
	RequiredHeader string
	AllowedValues  []string
}

func (p *HeaderValuePlugin) NewConfigInstance(ctx context.Context) (interface{}, error) {
	logger(ctx).Infow("Called 'NewConfigInstance' on HeaderValuePlugin")
	return &Config{}, nil
}

func (p *HeaderValuePlugin) GetAuthService(ctx context.Context, configInstance interface{}) (api.AuthService, error) {
	config, ok := configInstance.(*Config)
	if !ok {
		return nil, errors.New(fmt.Sprintf("unexpected config type %T", configInstance))
	}
	logger(ctx).Infow("Returning HeaderValueAuthService instance",
		zap.Any("requiredHeader", config.RequiredHeader),
		zap.Any("allowedHeaderValues", config.AllowedValues),
	)

	valueMap := map[string]bool{}
	for _, v := range config.AllowedValues {
		valueMap[v] = true
	}

	return &HeaderValueAuthService{
		RequiredHeader: config.RequiredHeader,
		AllowedValues:  valueMap,
	}, nil
}

type HeaderValueAuthService struct {
	RequiredHeader string
	AllowedValues  map[string]bool
}

func (c *HeaderValueAuthService) Start(ctx context.Context) error {
	logger(ctx).Infow("Called 'Start' on HeaderValueAuthService")
	return nil
}

func (c *HeaderValueAuthService) Authorize(ctx context.Context, request *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {
	for key, value := range request.CheckRequest.Attributes.Request.Http.Headers {
		if key == c.RequiredHeader {
			logger(ctx).Infow("Found required header, checking value.", "header", key, "value", value)

			if _, ok := c.AllowedValues[value]; ok {
				logger(ctx).Infow("Header value match. Allowing request.")
				response := api.AuthorizedResponse()

				// Append extra header
				response.CheckResponse.HttpResponse = &envoyauthv2.CheckResponse_OkResponse{
					OkResponse: &envoyauthv2.OkHttpResponse{
						Headers: []*envoycorev2.HeaderValueOption{{
							Header: &envoycorev2.HeaderValue{
								Key:   "matched-allowed-headers",
								Value: "true",
							},
						}},
					},
				}
				return response, nil
			}
			logger(ctx).Infow("Header value does not match allowed values, denying access.")
			return api.UnauthorizedResponse(), nil
		}
	}
	logger(ctx).Infow("Required header not found, denying access")
	return api.UnauthorizedResponse(), nil
}

func logger(ctx context.Context) *zap.SugaredLogger {
	return contextutils.LoggerFrom(contextutils.WithLogger(ctx, "header_value_plugin"))
}
