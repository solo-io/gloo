package v1helpers

import (
	ratelimit2 "github.com/solo-io/gloo/projects/gloo/api/external/solo/ratelimit"
	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSimpleRateLimitConfig returns a RateLimitConfig with a 1 request per minute
// value and actionValue must be the same for the config to be valid
func NewSimpleRateLimitConfig(name, namespace, key, value, actionValue string) *v1alpha1.RateLimitConfig {
	rateLimitConfig := &v1alpha1.RateLimitConfig{
		RateLimitConfig: ratelimit2.RateLimitConfig{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: rlv1alpha1.RateLimitConfigSpec{
				ConfigType: &rlv1alpha1.RateLimitConfigSpec_Raw_{
					Raw: &rlv1alpha1.RateLimitConfigSpec_Raw{
						Descriptors: []*rlv1alpha1.Descriptor{{
							Key:   key,
							Value: value,
							RateLimit: &rlv1alpha1.RateLimit{
								Unit:            rlv1alpha1.RateLimit_MINUTE,
								RequestsPerUnit: 1,
							},
						}},
						RateLimits: []*rlv1alpha1.RateLimitActions{{
							Actions: []*rlv1alpha1.Action{{
								ActionSpecifier: &rlv1alpha1.Action_GenericKey_{
									GenericKey: &rlv1alpha1.Action_GenericKey{
										DescriptorValue: actionValue,
									},
								},
							}},
						}},
					},
				},
			},
		},
	}
	return rateLimitConfig
}
