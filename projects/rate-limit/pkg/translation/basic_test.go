package translation_test

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation/internal"
)

var _ = Describe("Basic Rate Limit Config Translation", func() {

	var (
		virtualHostName string
		translator      translation.BasicRateLimitTranslator
	)

	BeforeEach(func() {
		virtualHostName = "my-virtual-host"
		translator = translation.NewBasicRateLimitTranslator()
	})

	Describe("generating server config", func() {
		It("works as expected", func() {
			input := ratelimit.IngressRateLimit{
				AuthorizedLimits: &v1alpha1.RateLimit{
					Unit:            v1alpha1.RateLimit_SECOND,
					RequestsPerUnit: 10,
				},
				AnonymousLimits: &v1alpha1.RateLimit{
					Unit:            v1alpha1.RateLimit_MINUTE,
					RequestsPerUnit: 1,
				},
			}

			expected := &v1alpha1.Descriptor{
				Key:   internal.GenericKey,
				Value: virtualHostName,
				Descriptors: []*v1alpha1.Descriptor{
					{
						Key:   internal.HeaderMatch,
						Value: internal.Anonymous,
						Descriptors: []*v1alpha1.Descriptor{
							{
								Key: internal.RemoteAddress,
								RateLimit: &v1alpha1.RateLimit{
									Unit:            v1alpha1.RateLimit_MINUTE,
									RequestsPerUnit: 1,
								},
							},
						},
					},
					{
						Key:   internal.HeaderMatch,
						Value: internal.Authenticated,
						Descriptors: []*v1alpha1.Descriptor{
							{
								Key: internal.UserId,
								RateLimit: &v1alpha1.RateLimit{
									Unit:            v1alpha1.RateLimit_SECOND,
									RequestsPerUnit: 10,
								},
							},
						},
					},
				},
			}

			actual, err := translator.GenerateServerConfig(virtualHostName, input)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(expected))
		})
	})

	Describe("generating client config", func() {
		It("works as expected", func() {
			headerName := "x-foo"
			stage := uint32(2)

			expected := []*envoy_config_route_v3.RateLimit{
				{
					Stage: &wrappers.UInt32Value{
						Value: stage,
					},
					Actions: []*envoy_config_route_v3.RateLimit_Action{
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
								GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
									DescriptorValue: virtualHostName,
								},
							},
						},
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_HeaderValueMatch_{
								HeaderValueMatch: &envoy_config_route_v3.RateLimit_Action_HeaderValueMatch{
									DescriptorValue: internal.Authenticated,
									ExpectMatch: &wrappers.BoolValue{
										Value: true,
									},
									Headers: []*envoy_config_route_v3.HeaderMatcher{
										{
											Name: headerName,
											HeaderMatchSpecifier: &envoy_config_route_v3.HeaderMatcher_PresentMatch{
												PresentMatch: true,
											},
										},
									},
								},
							},
						},
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_RequestHeaders_{
								RequestHeaders: &envoy_config_route_v3.RateLimit_Action_RequestHeaders{
									DescriptorKey: internal.UserId,
									HeaderName:    headerName,
								},
							},
						},
					},
				},
				{
					Stage: &wrappers.UInt32Value{
						Value: stage,
					},
					Actions: []*envoy_config_route_v3.RateLimit_Action{
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_GenericKey_{
								GenericKey: &envoy_config_route_v3.RateLimit_Action_GenericKey{
									DescriptorValue: virtualHostName,
								},
							},
						},
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_HeaderValueMatch_{
								HeaderValueMatch: &envoy_config_route_v3.RateLimit_Action_HeaderValueMatch{
									DescriptorValue: internal.Anonymous,
									ExpectMatch: &wrappers.BoolValue{
										Value: false,
									},
									Headers: []*envoy_config_route_v3.HeaderMatcher{
										{
											Name: headerName,
											HeaderMatchSpecifier: &envoy_config_route_v3.HeaderMatcher_PresentMatch{
												PresentMatch: true,
											},
										},
									},
								},
							},
						},
						{
							ActionSpecifier: &envoy_config_route_v3.RateLimit_Action_RemoteAddress_{
								RemoteAddress: &envoy_config_route_v3.RateLimit_Action_RemoteAddress{},
							},
						},
					},
				},
			}

			actual := translator.GenerateResourceConfig(virtualHostName, headerName, stage)
			Expect(actual).To(Equal(expected))
		})
	})
})
