package shims

import (
	"github.com/solo-io/rate-limiter/pkg/config/translation"
	solo_apis "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	solo_apis_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/internal"
)

type translator struct {
	rlTranslator translation.RateLimitConfigTranslator
}

func NewRateLimitConfigTranslator() RateLimitConfigTranslator {
	return &translator{
		rlTranslator: translation.NewRateLimitConfigTranslator(),
	}
}

func (t *translator) ToDescriptor(soloApiConfig *solo_apis.RateLimitConfig) (*solo_apis_types.Descriptor, error) {
	rlConfig, err := internal.ToRateLimiterResource(soloApiConfig)
	if err != nil {
		return nil, err
	}

	rlDescriptor, err := t.rlTranslator.ToDescriptor(rlConfig)
	if err != nil {
		return nil, err
	}

	return internal.ToSoloAPIsDescriptor(rlDescriptor)
}

func (t *translator) ToActions(soloApiConfig *solo_apis.RateLimitConfig) ([]*solo_apis_types.RateLimitActions, error) {
	rlConfig, err := internal.ToRateLimiterResource(soloApiConfig)
	if err != nil {
		return nil, err
	}

	rlActions, err := t.rlTranslator.ToActions(rlConfig)
	if err != nil {
		return nil, err
	}

	return internal.ToSoloAPIsActionsSlice(rlActions)
}
