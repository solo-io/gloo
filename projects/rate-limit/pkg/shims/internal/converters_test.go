package internal_test

import (
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	rl_api "github.com/solo-io/rate-limiter/pkg/api/ratelimit.solo.io/v1alpha1"
	rl_api_types "github.com/solo-io/rate-limiter/pkg/api/ratelimit.solo.io/v1alpha1/types"
	solo_apis "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	solo_apis_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Converters", func() {

	Describe("converting a solo-apis RateLimitConfig type to its rate-limiter equivalent", func() {
		var (
			soloApiResource         *solo_apis.RateLimitConfig
			rlApiEquivalentResource *rl_api.RateLimitConfig
		)

		BeforeEach(func() {
			soloApiResource = &solo_apis.RateLimitConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "RateLimitConfig",
					APIVersion: "ratelimit.solo.io/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "foo",
					Namespace:       "bar",
					ResourceVersion: "123",
					Labels: map[string]string{
						"foo": "bar",
					},
					Annotations: map[string]string{
						"baz": "bar",
					},
				},
				Spec: solo_apis_types.RateLimitConfigSpec{
					ConfigType: &solo_apis_types.RateLimitConfigSpec_Raw_{
						Raw: &solo_apis_types.RateLimitConfigSpec_Raw{
							Descriptors: []*solo_apis_types.Descriptor{
								{
									Key:   "key",
									Value: "val",
									RateLimit: &solo_apis_types.RateLimit{
										Unit:            solo_apis_types.RateLimit_SECOND,
										RequestsPerUnit: 10,
									},
									Descriptors: []*solo_apis_types.Descriptor{
										{
											Key:   "nested-key",
											Value: "nested-val",
											RateLimit: &solo_apis_types.RateLimit{
												Unit:            solo_apis_types.RateLimit_SECOND,
												RequestsPerUnit: 20,
											},
										},
									},
									Weight:      42,
									AlwaysApply: true,
								},
							},
							RateLimits: []*solo_apis_types.RateLimitActions{
								{
									Actions: []*solo_apis_types.Action{
										{
											ActionSpecifier: &solo_apis_types.Action_GenericKey_{
												GenericKey: &solo_apis_types.Action_GenericKey{
													DescriptorValue: "foo",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Status: solo_apis_types.RateLimitConfigStatus{
					State:              solo_apis_types.RateLimitConfigStatus_ACCEPTED,
					Message:            "hello",
					ObservedGeneration: 2,
				},
			}

			rlApiEquivalentResource = &rl_api.RateLimitConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "RateLimitConfig",
					APIVersion: "ratelimit.solo.io/v1alpha1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "foo",
					Namespace:       "bar",
					ResourceVersion: "123",
					Labels: map[string]string{
						"foo": "bar",
					},
					Annotations: map[string]string{
						"baz": "bar",
					},
				},
				Spec: rl_api_types.RateLimitConfigSpec{
					ConfigType: &rl_api_types.RateLimitConfigSpec_Raw_{
						Raw: &rl_api_types.RateLimitConfigSpec_Raw{
							Descriptors: []*rl_api_types.Descriptor{
								{
									Key:   "key",
									Value: "val",
									RateLimit: &rl_api_types.RateLimit{
										Unit:            rl_api_types.RateLimit_SECOND,
										RequestsPerUnit: 10,
									},
									Descriptors: []*rl_api_types.Descriptor{
										{
											Key:   "nested-key",
											Value: "nested-val",
											RateLimit: &rl_api_types.RateLimit{
												Unit:            rl_api_types.RateLimit_SECOND,
												RequestsPerUnit: 20,
											},
										},
									},
									Weight:      42,
									AlwaysApply: true,
								},
							},
							RateLimits: []*rl_api_types.RateLimitActions{
								{
									Actions: []*rl_api_types.Action{
										{
											ActionSpecifier: &rl_api_types.Action_GenericKey_{
												GenericKey: &rl_api_types.Action_GenericKey{
													DescriptorValue: "foo",
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Status: rl_api_types.RateLimitConfigStatus{
					State:              rl_api_types.RateLimitConfigStatus_ACCEPTED,
					Message:            "hello",
					ObservedGeneration: 2,
				},
			}
		})

		It("should successfully convert the resource", func() {
			actual, err := internal.ToRateLimiterResource(soloApiResource)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(rlApiEquivalentResource))
		})

		Measure("1000 function calls", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				for i := 0; i < 1000; i++ {
					_, err := internal.ToRateLimiterResource(soloApiResource)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 50))
		}, 10)
	})

	Describe("converting rate-limiter actions to their solo-apis equivalents", func() {
		var (
			rlApiActions             []*rl_api_types.RateLimitActions
			equivalentSoloApiActions []*solo_apis_types.RateLimitActions
		)

		BeforeEach(func() {
			rlApiActions = []*rl_api_types.RateLimitActions{
				{
					Actions: []*rl_api_types.Action{
						{
							ActionSpecifier: &rl_api_types.Action_GenericKey_{
								GenericKey: &rl_api_types.Action_GenericKey{
									DescriptorValue: "foo",
								},
							},
						},
					},
				},
				{
					Actions: []*rl_api_types.Action{
						{
							ActionSpecifier: &rl_api_types.Action_HeaderValueMatch_{
								HeaderValueMatch: &rl_api_types.Action_HeaderValueMatch{
									DescriptorValue: "bar",
									ExpectMatch:     &types.BoolValue{Value: true},
									Headers: []*rl_api_types.Action_HeaderValueMatch_HeaderMatcher{
										{
											Name: "baz",
											HeaderMatchSpecifier: &rl_api_types.Action_HeaderValueMatch_HeaderMatcher_RegexMatch{
												RegexMatch: ".*",
											},
											InvertMatch: true,
										},
									},
								},
							},
						},
					},
				},
			}
			equivalentSoloApiActions = []*solo_apis_types.RateLimitActions{
				{
					Actions: []*solo_apis_types.Action{
						{
							ActionSpecifier: &solo_apis_types.Action_GenericKey_{
								GenericKey: &solo_apis_types.Action_GenericKey{
									DescriptorValue: "foo",
								},
							},
						},
					},
				},
				{
					Actions: []*solo_apis_types.Action{
						{
							ActionSpecifier: &solo_apis_types.Action_HeaderValueMatch_{
								HeaderValueMatch: &solo_apis_types.Action_HeaderValueMatch{
									DescriptorValue: "bar",
									ExpectMatch:     &types.BoolValue{Value: true},
									Headers: []*solo_apis_types.Action_HeaderValueMatch_HeaderMatcher{
										{
											Name: "baz",
											HeaderMatchSpecifier: &solo_apis_types.Action_HeaderValueMatch_HeaderMatcher_RegexMatch{
												RegexMatch: ".*",
											},
											InvertMatch: true,
										},
									},
								},
							},
						},
					},
				},
			}
		})

		It("should successfully convert the resources", func() {
			actual, err := internal.ToSoloAPIsActionsSlice(rlApiActions)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(equivalentSoloApiActions))
		})

		Measure("1000 function calls", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				for i := 0; i < 1000; i++ {
					_, err := internal.ToSoloAPIsActionsSlice(rlApiActions)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 50))
		}, 10)
	})

	Describe("converting solo-apis descriptors to their rate-limiter equivalents", func() {
		var (
			soloApiDescriptors         []*solo_apis_types.Descriptor
			rlApiEquivalentDescriptors []*rl_api_types.Descriptor
		)

		BeforeEach(func() {
			soloApiDescriptors = []*solo_apis_types.Descriptor{
				{
					Key:   "key-1",
					Value: "val",
					RateLimit: &solo_apis_types.RateLimit{
						Unit:            solo_apis_types.RateLimit_SECOND,
						RequestsPerUnit: 10,
					},
					Descriptors: []*solo_apis_types.Descriptor{
						{
							Key:   "nested-key",
							Value: "nested-val",
							RateLimit: &solo_apis_types.RateLimit{
								Unit:            solo_apis_types.RateLimit_SECOND,
								RequestsPerUnit: 20,
							},
						},
					},
					Weight:      42,
					AlwaysApply: true,
				},
				{
					Key:   "key-2",
					Value: "val-2",
					RateLimit: &solo_apis_types.RateLimit{
						Unit:            solo_apis_types.RateLimit_HOUR,
						RequestsPerUnit: 3600,
					},
				},
			}
			rlApiEquivalentDescriptors = []*rl_api_types.Descriptor{
				{
					Key:   "key-1",
					Value: "val",
					RateLimit: &rl_api_types.RateLimit{
						Unit:            rl_api_types.RateLimit_SECOND,
						RequestsPerUnit: 10,
					},
					Descriptors: []*rl_api_types.Descriptor{
						{
							Key:   "nested-key",
							Value: "nested-val",
							RateLimit: &rl_api_types.RateLimit{
								Unit:            rl_api_types.RateLimit_SECOND,
								RequestsPerUnit: 20,
							},
						},
					},
					Weight:      42,
					AlwaysApply: true,
				},
				{
					Key:   "key-2",
					Value: "val-2",
					RateLimit: &rl_api_types.RateLimit{
						Unit:            rl_api_types.RateLimit_HOUR,
						RequestsPerUnit: 3600,
					},
				},
			}
		})

		It("should successfully convert the resources", func() {
			actual, err := internal.ToRateLimiterDescriptors(soloApiDescriptors)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(rlApiEquivalentDescriptors))
		})

		Measure("1000 function calls", func(b Benchmarker) {
			runtime := b.Time("runtime", func() {
				for i := 0; i < 1000; i++ {
					_, err := internal.ToRateLimiterDescriptors(soloApiDescriptors)
					Expect(err).NotTo(HaveOccurred())
				}
			})

			Expect(runtime.Milliseconds()).Should(BeNumerically("<", 50))
		}, 10)
	})
})
