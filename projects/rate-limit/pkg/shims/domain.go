package shims

import (
	"context"

	"github.com/solo-io/rate-limiter/pkg/config"
	solo_api_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/internal"
)

type domainGenerator struct{}

func NewRateLimitDomainGenerator() RateLimitDomainGenerator {
	return domainGenerator{}
}

func (domainGenerator) NewRateLimitDomain(ctx context.Context, domain string, descriptors []*solo_api_rl.Descriptor) (config.RateLimitDomain, error) {
	// Convert descriptors from the solo-api type to the rate-limiter type
	convertedDescriptors, err := internal.ToRateLimiterDescriptors(descriptors)
	if err != nil {
		return nil, err
	}

	return config.NewRateLimitDomain(ctx, domain, convertedDescriptors)
}
