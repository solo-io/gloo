package ratelimit_test

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
)

func authorizedRateLimits(requestsPerUnit uint32, unit rlv1alpha1.RateLimit_Unit) *ratelimit.IngressRateLimit {
	return &ratelimit.IngressRateLimit{
		AuthorizedLimits: &rlv1alpha1.RateLimit{
			RequestsPerUnit: requestsPerUnit,
			Unit:            unit,
		},
	}
}

func anonymousRateLimits(requestsPerUnit uint32, unit rlv1alpha1.RateLimit_Unit) *ratelimit.IngressRateLimit {
	return &ratelimit.IngressRateLimit{
		AnonymousLimits: &rlv1alpha1.RateLimit{
			RequestsPerUnit: requestsPerUnit,
			Unit:            unit,
		},
	}
}
