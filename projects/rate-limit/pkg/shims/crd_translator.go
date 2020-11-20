package shims

import (
	"github.com/solo-io/rate-limiter/pkg/config/translation"
	solo_apis "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/internal"
)

type crd_translator struct {
	rlTranslator translation.RateLimitConfigTranslator
}

func NewRateLimitConfigTranslator() RateLimitConfigTranslator {
	return &crd_translator{
		rlTranslator: translation.NewRateLimitConfigTranslator(),
	}
}

func (t *crd_translator) ToDescriptors(soloApiConfig *solo_apis.RateLimitConfig) (*solo_apis.RateLimitConfigSpec_Raw, error) {
	rlConfig, err := internal.ToRateLimiterResource(soloApiConfig)
	if err != nil {
		return nil, err
	}

	rlDescriptors, err := t.rlTranslator.ToDescriptors(rlConfig)
	if err != nil {
		return nil, err
	}

	return internal.ToSoloAPIsResourceSpec_Raw(rlDescriptors)
}

func (t *crd_translator) ToActions(soloApiConfig *solo_apis.RateLimitConfig) ([]*solo_apis.RateLimitActions, error) {
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
