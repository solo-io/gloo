package shims

import (
	"context"

	"github.com/solo-io/rate-limiter/pkg/config"
	solo_apis "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	solo_apis_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
)

//go:generate mockgen -destination ./mocks/mock_interfaces.go -source ./interfaces.go

// A shim for `github.com/solo-io/rate-limiter/pkg/config/translation.RateLimitConfigTranslator`
type RateLimitConfigTranslator interface {
	// Translate the `RateLimitConfig` into the descriptor that can will be used to configure the server.
	ToDescriptor(config *solo_apis.RateLimitConfig) (*solo_apis_types.Descriptor, error)

	// Translate the `RateLimitConfig` into the rate limit actions that can will be used to configure Envoy.
	ToActions(config *solo_apis.RateLimitConfig) ([]*solo_apis_types.RateLimitActions, error)
}

// A shim for `github.com/solo-io/rate-limiter/pkg/config.RateLimitDomain`
type RateLimitDomainGenerator interface {
	NewRateLimitDomain(ctx context.Context, domain string, descriptors []*solo_apis.Descriptor) (config.RateLimitDomain, error)
}
