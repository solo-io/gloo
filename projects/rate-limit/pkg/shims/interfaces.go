package shims

import (
	"context"

	"github.com/solo-io/rate-limiter/pkg/config"
	solo_apis "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
)

//go:generate mockgen -destination ./mocks/mock_interfaces.go -source ./interfaces.go

// A shim for `github.com/solo-io/rate-limiter/pkg/config/translation.RateLimitConfigTranslator`
type RateLimitConfigTranslator interface {
	// Translate the `RateLimitConfig` into the descriptors and set descriptors that will be used to configure the server.
	ToDescriptors(soloApiConfig *solo_apis.RateLimitConfig) (*solo_apis.RateLimitConfigSpec_Raw, error)

	// Translate the `RateLimitConfig` into the rate limit actions that will be used to configure Envoy.
	ToActions(config *solo_apis.RateLimitConfig) ([]*solo_apis.RateLimitActions, error)
}

// A shim for `github.com/solo-io/rate-limiter/pkg/config.RateLimitDomain`
type RateLimitDomainGenerator interface {
	NewRateLimitDomain(ctx context.Context, configId, domain string, rateLimitConfig *solo_apis.RateLimitConfigSpec_Raw) (config.RateLimitDomain, error)
}

// A shim for `github.com/solo-io/rate-limiter/pkg/config/translation.GlobalRateLimitTranslator`
type GlobalRateLimitTranslator interface {
	// Translate the rate limit setDescriptors into the set descriptors that will be used to configure the server.
	// This function also checks the tree descriptors to ensure they don't contain the special purpose set descriptor genericKey and errors if they do.
	ToSetDescriptors(descriptors []*solo_apis.Descriptor, setDescriptors []*solo_apis.SetDescriptor) ([]*solo_apis.SetDescriptor, error)

	// Translate the rate limit actions from a `VirtualService` or `Route` into the rate limit actions that will be used to configure Envoy.
	ToActions(rlActions []*solo_apis.RateLimitActions) ([]*solo_apis.RateLimitActions, error)
}
