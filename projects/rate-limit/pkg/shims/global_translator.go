package shims

import (
	"github.com/solo-io/rate-limiter/pkg/config/translation"
	solo_apis "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/internal"
)

type global_translator struct {
	rlTranslator translation.GlobalRateLimitTranslator
}

func NewGlobalRateLimitTranslator() GlobalRateLimitTranslator {
	return &global_translator{
		rlTranslator: translation.NewGlobalRateLimitTranslator(),
	}
}

func (t *global_translator) ToSetDescriptors(soloApiDescriptors []*solo_apis.Descriptor, soloApiSetDescriptors []*solo_apis.SetDescriptor) ([]*solo_apis.SetDescriptor, error) {
	rlDescriptors, err := internal.ToRateLimiterDescriptors(soloApiDescriptors)
	if err != nil {
		return nil, err
	}
	rlSetDescriptors, err := internal.ToRateLimiterSetDescriptors(soloApiSetDescriptors)
	if err != nil {
		return nil, err
	}

	setDescriptors, err := t.rlTranslator.ToSetDescriptors(rlDescriptors, rlSetDescriptors)
	if err != nil {
		return nil, err
	}

	return internal.ToSoloAPIsSetDescriptors(setDescriptors)
}

func (t *global_translator) ToActions(soloApiActions []*solo_apis.RateLimitActions) ([]*solo_apis.RateLimitActions, error) {
	rlConfig, err := internal.ToRateLimiterActionsSlice(soloApiActions)
	if err != nil {
		return nil, err
	}

	rlActions, err := t.rlTranslator.ToActions(rlConfig)
	if err != nil {
		return nil, err
	}

	return internal.ToSoloAPIsActionsSlice(rlActions)
}
