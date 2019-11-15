package converter

import (
	"context"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
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

	return details
}
