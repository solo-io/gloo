package local_ratelimit

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	envoy_extensions_filters_http_local_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/local_ratelimit/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	local_ratelimit "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/local_ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Local Rate Limit Plugin", func() {
	var p *plugin
	var httpListener *v1.HttpListener
	var tokenBucket *local_ratelimit.TokenBucket
	var expectedEarlyFilter *envoy_extensions_filters_http_local_ratelimit_v3.LocalRateLimit

	BeforeEach(func() {
		p = NewPlugin()
		p.Init(plugins.InitParams{})

		httpListener = &v1.HttpListener{
			Options: &v1.HttpListenerOptions{
				HttpLocalRatelimit: &local_ratelimit.Settings{
					LocalRateLimitPerDownstreamConnection: &wrapperspb.BoolValue{Value: true},
					EnableXRatelimitHeaders:               &wrapperspb.BoolValue{Value: true},
					DefaultLimit: &local_ratelimit.TokenBucket{
						MaxTokens: 10,
						TokensPerFill: &wrapperspb.UInt32Value{
							Value: 10,
						},
						FillInterval: &durationpb.Duration{
							Seconds: 10,
						},
					},
				},
			},
		}

		tokenBucket = &local_ratelimit.TokenBucket{
			MaxTokens: 10,
			TokensPerFill: &wrapperspb.UInt32Value{
				Value: 10,
			},
			FillInterval: &durationpb.Duration{
				Seconds: 10,
			},
		}

		expectedEarlyFilter = &envoy_extensions_filters_http_local_ratelimit_v3.LocalRateLimit{
			StatPrefix:                            HTTPFilterStatPrefix,
			LocalRateLimitPerDownstreamConnection: true,
			EnableXRatelimitHeaders:               envoyratelimit.XRateLimitHeadersRFCVersion_DRAFT_VERSION_03,
			Stage:                                 CustomStageBeforeAuth,
			FilterEnabled: &corev3.RuntimeFractionalPercent{
				DefaultValue: &envoy_type_v3.FractionalPercent{
					Numerator:   100,
					Denominator: envoy_type_v3.FractionalPercent_HUNDRED,
				},
			},
			FilterEnforced: &corev3.RuntimeFractionalPercent{
				DefaultValue: &envoy_type_v3.FractionalPercent{
					Numerator:   100,
					Denominator: envoy_type_v3.FractionalPercent_HUNDRED,
				},
			},
			TokenBucket: &envoy_type_v3.TokenBucket{
				MaxTokens: 10,
				TokensPerFill: &wrapperspb.UInt32Value{
					Value: 10,
				},
				FillInterval: &durationpb.Duration{
					Seconds: 10,
				},
			},
		}
	})

	When("The local rate limit early config is defined on the CR", func() {
		It("Should copy the config from the vHost", func() {
			out := &envoy_config_route_v3.VirtualHost{}
			err := p.ProcessVirtualHost(plugins.VirtualHostParams{
				HttpListener: httpListener,
			}, &v1.VirtualHost{
				Options: &v1.VirtualHostOptions{
					RateLimitEarlyConfigType: &v1.VirtualHostOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitVhostExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).NotTo(HaveOccurred())
			typedConfig, err := utils.MessageToAny(expectedEarlyFilter)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(&envoy_config_route_v3.VirtualHost{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}))
		})

		It("Should copy the config from the route", func() {
			out := &envoy_config_route_v3.Route{}
			err := p.ProcessRoute(plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					HttpListener: httpListener,
				},
			}, &v1.Route{
				Options: &v1.RouteOptions{
					RateLimitEarlyConfigType: &v1.RouteOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).NotTo(HaveOccurred())
			typedConfig, err := utils.MessageToAny(expectedEarlyFilter)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(&envoy_config_route_v3.Route{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}))
		})
	})

	When("The local rate limit is already configured in OSS", func() {
		It("Should error if the filter is already defined in OSS on the vHost and not change the existing config", func() {
			typedConfig, err := utils.MessageToAny(expectedEarlyFilter)
			Expect(err).NotTo(HaveOccurred())
			out := &envoy_config_route_v3.VirtualHost{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}

			err = p.ProcessVirtualHost(plugins.VirtualHostParams{
				HttpListener: httpListener,
			}, &v1.VirtualHost{
				Options: &v1.VirtualHostOptions{
					RateLimitEarlyConfigType: &v1.VirtualHostOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitVhostExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).To(Equal(ErrFilterDefinedInOSS))
			Expect(out).To(Equal(&envoy_config_route_v3.VirtualHost{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}))
		})

		It("Should error if the filter is already defined in OSS on the route and not change the existing config", func() {
			typedConfig, err := utils.MessageToAny(expectedEarlyFilter)
			Expect(err).NotTo(HaveOccurred())
			out := &envoy_config_route_v3.Route{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}

			err = p.ProcessRoute(plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					HttpListener: httpListener,
				},
			}, &v1.Route{
				Options: &v1.RouteOptions{
					RateLimitEarlyConfigType: &v1.RouteOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).To(Equal(ErrFilterDefinedInOSS))
			Expect(out).To(Equal(&envoy_config_route_v3.Route{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}))

		})
	})

	When("The local rate limit is configured in ratelimitRegular", func() {
		It("Should error if the filter is already defined in OSS on the vHost and not change the existing config", func() {
			typedConfig, err := utils.MessageToAny(expectedEarlyFilter)
			Expect(err).NotTo(HaveOccurred())
			out := &envoy_config_route_v3.VirtualHost{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}

			err = p.ProcessVirtualHost(plugins.VirtualHostParams{
				HttpListener: httpListener,
			}, &v1.VirtualHost{
				Options: &v1.VirtualHostOptions{
					RateLimitRegularConfigType: &v1.VirtualHostOptions_RatelimitRegular{
						RatelimitRegular: &ratelimit.RateLimitVhostExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).To(Equal(ErrFilterDefinedInRegular))
			Expect(out).To(Equal(&envoy_config_route_v3.VirtualHost{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}))
		})

		It("Should error if the filter is already defined in early on the vHost and not change the existing config", func() {
			out := &envoy_config_route_v3.Route{}
			err := p.ProcessRoute(plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					HttpListener: httpListener,
				},
			}, &v1.Route{
				Options: &v1.RouteOptions{
					RateLimitEarlyConfigType: &v1.RouteOptions_RatelimitEarly{
						RatelimitEarly: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: tokenBucket,
						},
					},
					RateLimitRegularConfigType: &v1.RouteOptions_RatelimitRegular{
						RatelimitRegular: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).To(Equal(ErrFilterDefinedInRegular))
			typedConfig, err := utils.MessageToAny(expectedEarlyFilter)
			Expect(err).NotTo(HaveOccurred())
			Expect(out).To(Equal(&envoy_config_route_v3.Route{
				TypedPerFilterConfig: map[string]*anypb.Any{
					HTTPFilterName: typedConfig,
				},
			}))
		})

		It("Should error out and not add any limits if only configured in ratelimitRegular", func() {
			out := &envoy_config_route_v3.Route{}
			err := p.ProcessRoute(plugins.RouteParams{
				VirtualHostParams: plugins.VirtualHostParams{
					HttpListener: httpListener,
				},
			}, &v1.Route{
				Options: &v1.RouteOptions{
					RateLimitRegularConfigType: &v1.RouteOptions_RatelimitRegular{
						RatelimitRegular: &ratelimit.RateLimitRouteExtension{
							LocalRatelimit: tokenBucket,
						},
					},
				},
			}, out)

			Expect(err).To(Equal(ErrFilterDefinedInRegular))
			Expect(out).To(Equal(&envoy_config_route_v3.Route{}))
		})
	})
})
