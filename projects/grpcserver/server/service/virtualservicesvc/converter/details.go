package converter

import (
	"context"

	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	extauthapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
	ratelimitapi "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/details_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter VirtualServiceDetailsConverter

const (
	FailedToParseExtAuthConfig   = "Failed to parse extauth config"
	FailedToParseRateLimitConfig = "Failed to parse rate limit config"
)

type VirtualServiceDetailsConverter interface {
	GetDetails(ctx context.Context, vs *gatewayv1.VirtualService) *v1.VirtualServiceDetails
}

type virtualServiceDetailsConverter struct {
	rawGetter rawgetter.RawGetter
}

var _ VirtualServiceDetailsConverter = virtualServiceDetailsConverter{}

func NewVirtualServiceDetailsConverter(r rawgetter.RawGetter) VirtualServiceDetailsConverter {
	return virtualServiceDetailsConverter{rawGetter: r}
}

func (c virtualServiceDetailsConverter) GetDetails(ctx context.Context, vs *gatewayv1.VirtualService) *v1.VirtualServiceDetails {
	details := &v1.VirtualServiceDetails{
		VirtualService: vs,
		Raw:            c.rawGetter.GetRaw(ctx, vs, gatewayv1.VirtualServiceCrd),
	}

	var configs map[string]*types.Struct
	if configs = vs.GetVirtualHost().GetVirtualHostPlugins().GetExtensions().GetConfigs(); configs == nil {
		return details
	}

	details.Plugins = &v1.Plugins{}

	if extAuthStruct, ok := configs[extauth.ExtensionName]; ok {
		details.Plugins.ExtAuth = &v1.ExtAuthPlugin{}
		extAuth := &extauthapi.ExtAuthConfig{}
		err := util.StructToMessage(extAuthStruct, extAuth)
		if err != nil {
			details.Plugins.ExtAuth.Error = FailedToParseExtAuthConfig
			contextutils.LoggerFrom(ctx).Errorw(FailedToParseExtAuthConfig, zap.Error(err), zap.Any("virtualService", vs))
		} else {
			details.Plugins.ExtAuth.Value = extAuth
		}
	}

	if rateLimitStruct, ok := configs[ratelimit.ExtensionName]; ok {
		details.Plugins.RateLimit = &v1.RateLimitPlugin{}
		rateLimit := &ratelimitapi.IngressRateLimit{}
		err := util.StructToMessage(rateLimitStruct, rateLimit)
		if err != nil {
			details.Plugins.RateLimit.Error = FailedToParseRateLimitConfig
			contextutils.LoggerFrom(ctx).Errorw(FailedToParseRateLimitConfig, zap.Error(err), zap.Any("virtualService", vs))
		} else {
			details.Plugins.RateLimit.Value = rateLimit
		}
	}

	return details
}
