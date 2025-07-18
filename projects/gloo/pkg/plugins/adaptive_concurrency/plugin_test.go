package adaptiveconcurrency_test

import (
	envoy_adaptive_concurrency_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/adaptive_concurrency/v3"
	envoyhcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/adaptive_concurrency"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/plugins/adaptive_concurrency"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/test/matchers"
)

// baseConfigParams holds the parameters for creating a base expected config
type baseConfigParams struct {
	ConcurrencyUpdateIntervalMillis uint32
	MinRttIntervalMillis            uint32
	MinRttFixedValueMillis          uint32
}

// createBaseExpectedConfig creates a base expected config with defaults
// params.ConcurrencyUpdateIntervalMillis: the interval for concurrency updates in milliseconds
// params.MinRttIntervalMillis: the interval for min RTT calculation in milliseconds (0 to use fixed value instead)
// params.MinRttFixedValueMillis: the fixed value for min RTT calculation in milliseconds (only used if MinRttIntervalMillis is 0)
func createBaseExpectedConfig(params baseConfigParams) *envoy_adaptive_concurrency_v3.AdaptiveConcurrency {
	concurrencyUpdateDuration := &durationpb.Duration{
		Seconds: int64(params.ConcurrencyUpdateIntervalMillis / 1000),
		Nanos:   int32(params.ConcurrencyUpdateIntervalMillis%1000) * 1_000_000,
	}

	minRttCalcParams := &envoy_adaptive_concurrency_v3.GradientControllerConfig_MinimumRTTCalculationParams{
		RequestCount:   &wrapperspb.UInt32Value{Value: DefaultRequestCount},
		MinConcurrency: &wrapperspb.UInt32Value{Value: DefaultMinConcurrency},
		Jitter:         &typev3.Percent{Value: DefaultJitterPercentile},
		Buffer:         &typev3.Percent{Value: DefaultBufferPercentile},
	}

	if params.MinRttIntervalMillis > 0 {
		minRttCalcParams.Interval = &durationpb.Duration{
			Seconds: int64(params.MinRttIntervalMillis / 1000),
			Nanos:   int32(params.MinRttIntervalMillis%1000) * 1_000_000,
		}
	} else {
		minRttCalcParams.FixedValue = &durationpb.Duration{
			Seconds: int64(params.MinRttFixedValueMillis / 1000),
			Nanos:   int32(params.MinRttFixedValueMillis%1000) * 1_000_000,
		}
	}

	return &envoy_adaptive_concurrency_v3.AdaptiveConcurrency{
		ConcurrencyControllerConfig: &envoy_adaptive_concurrency_v3.AdaptiveConcurrency_GradientControllerConfig{
			GradientControllerConfig: &envoy_adaptive_concurrency_v3.GradientControllerConfig{
				SampleAggregatePercentile: &typev3.Percent{Value: DefaultSampleAggregatePercentile},
				ConcurrencyLimitParams: &envoy_adaptive_concurrency_v3.GradientControllerConfig_ConcurrencyLimitCalculationParams{
					ConcurrencyUpdateInterval: concurrencyUpdateDuration,
					MaxConcurrencyLimit:       &wrapperspb.UInt32Value{Value: DefaultMaxConcurrencyLimit},
				},
				MinRttCalcParams: minRttCalcParams,
			},
		},
		ConcurrencyLimitExceededStatus: &typev3.HttpStatus{
			Code: typev3.StatusCode(DefaultConcurrencyLimitExceededStatus),
		},
	}
}

var _ = Describe("Adaptive Concurrency Plugin", func() {
	var (
		p        plugins.HttpFilterPlugin
		params   plugins.Params
		listener *v1.HttpListener
	)

	BeforeEach(func() {
		p = NewPlugin()
		p.Init(plugins.InitParams{})
		params = plugins.Params{}
		listener = &v1.HttpListener{
			Options: &v1.HttpListenerOptions{},
		}
	})

	Describe("Name", func() {
		It("should return the correct extension name", func() {
			Expect(p.Name()).To(Equal("envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency"))
		})
	})

	Describe("HttpFilters", func() {
		Context("when adaptive concurrency is not configured", func() {
			It("should return empty filters", func() {
				filters, err := p.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())
				Expect(filters).To(BeEmpty())
			})
		})

		Context("when adaptive concurrency is configured with valid settings", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					SampleAggregatePercentile:       &wrapperspb.DoubleValue{Value: 60.0},
					ConcurrencyUpdateIntervalMillis: 1000,
					MaxConcurrencyLimit:             &wrapperspb.UInt32Value{Value: 100},
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						IntervalMillis:   500,
						RequestCount:     &wrapperspb.UInt32Value{Value: 10},
						MinConcurrency:   &wrapperspb.UInt32Value{Value: 5},
						JitterPercentile: &wrapperspb.DoubleValue{Value: 20.0},
						BufferPercentile: &wrapperspb.DoubleValue{Value: 30.0},
					},
					ConcurrencyLimitExceededStatus: 503,
				}
			})

			It("should create the correct filter", func() {
				filters, err := p.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())
				Expect(filters).To(HaveLen(1))

				filter := filters[0]
				Expect(filter.Stage).To(Equal(plugins.DuringStage(plugins.RateLimitStage)))

				expectedConfig := createBaseExpectedConfig(baseConfigParams{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttIntervalMillis:            500,
				})
				// Override custom values
				gradientConfig := expectedConfig.ConcurrencyControllerConfig.(*envoy_adaptive_concurrency_v3.AdaptiveConcurrency_GradientControllerConfig).GradientControllerConfig
				gradientConfig.SampleAggregatePercentile.Value = 60.0
				gradientConfig.ConcurrencyLimitParams.MaxConcurrencyLimit.Value = 100
				gradientConfig.MinRttCalcParams.RequestCount.Value = 10
				gradientConfig.MinRttCalcParams.MinConcurrency.Value = 5
				gradientConfig.MinRttCalcParams.Jitter.Value = 20.0
				gradientConfig.MinRttCalcParams.Buffer.Value = 30.0

				typedConfig, err := utils.MessageToAny(expectedConfig)
				Expect(err).NotTo(HaveOccurred())

				expectedFilter := &envoyhcm.HttpFilter{
					Name: "envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency",
					ConfigType: &envoyhcm.HttpFilter_TypedConfig{
						TypedConfig: typedConfig,
					},
				}

				Expect(filter.Filter).To(matchers.MatchProto(expectedFilter))
			})
		})

		Context("when concurrency_update_interval_millis is missing", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 0, // Missing/zero
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						IntervalMillis: 500,
					},
				}
			})

			It("should return the exported error", func() {
				_, err := p.HttpFilters(params, listener)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrConcurrencyUpdateIntervalMillisMissing))
			})
		})

		Context("when min_rtt_calc_params is missing", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams:                nil, // Missing
				}
			})

			It("should return the exported error", func() {
				_, err := p.HttpFilters(params, listener)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrMinRttCalcParamsMissing))
			})
		})

		Context("when both interval_millis and fixed_value_millis are missing", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams:                &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{},
				}
			})

			It("should return the exported error", func() {
				_, err := p.HttpFilters(params, listener)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrIntervalMissing))
			})
		})

		Context("when both interval_millis and fixed_value_millis are 0", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						IntervalMillis:   0, // Interval is 0
						FixedValueMillis: 0, // FixedValue is 0 (invalid)
					},
				}
			})

			It("should return the fixed value validation error", func() {
				_, err := p.HttpFilters(params, listener)
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(ErrIntervalMissing))
			})
		})

		Context("when using fixed_value_millis instead of interval_millis", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						FixedValueMillis: 200, // Use fixed value instead of interval
					},
				}
			})

			It("should create filter with fixed value and defaults", func() {
				filters, err := p.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())
				Expect(filters).To(HaveLen(1))

				filter := filters[0]
				expectedConfig := createBaseExpectedConfig(baseConfigParams{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttFixedValueMillis:          200,
				})

				typedConfig, err := utils.MessageToAny(expectedConfig)
				Expect(err).NotTo(HaveOccurred())

				expectedFilter := &envoyhcm.HttpFilter{
					Name: "envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency",
					ConfigType: &envoyhcm.HttpFilter_TypedConfig{
						TypedConfig: typedConfig,
					},
				}

				Expect(filter.Filter).To(matchers.MatchProto(expectedFilter))
			})
		})

		Context("when both interval_millis and fixed_value_millis are set", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						IntervalMillis:   500,
						FixedValueMillis: 200, // Both set - interval should take precedence
					},
				}
			})

			It("should use interval_millis and ignore fixed_value_millis with defaults", func() {
				filters, err := p.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())
				Expect(filters).To(HaveLen(1))

				filter := filters[0]
				expectedConfig := createBaseExpectedConfig(baseConfigParams{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttIntervalMillis:            500,
				})

				typedConfig, err := utils.MessageToAny(expectedConfig)
				Expect(err).NotTo(HaveOccurred())

				expectedFilter := &envoyhcm.HttpFilter{
					Name: "envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency",
					ConfigType: &envoyhcm.HttpFilter_TypedConfig{
						TypedConfig: typedConfig,
					},
				}

				Expect(filter.Filter).To(matchers.MatchProto(expectedFilter))
			})
		})

		Context("when concurrency_limit_exceeded_status is set to non-error code", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						IntervalMillis: 500,
					},
					ConcurrencyLimitExceededStatus: 200, // Non-error code
				}
			})

			It("should default to 503 for non-error codes", func() {
				filters, err := p.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())
				Expect(filters).To(HaveLen(1))

				filter := filters[0]
				expectedConfig := createBaseExpectedConfig(baseConfigParams{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttIntervalMillis:            500,
				})

				typedConfig, err := utils.MessageToAny(expectedConfig)
				Expect(err).NotTo(HaveOccurred())

				expectedFilter := &envoyhcm.HttpFilter{
					Name: "envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency",
					ConfigType: &envoyhcm.HttpFilter_TypedConfig{
						TypedConfig: typedConfig,
					},
				}

				Expect(filter.Filter).To(matchers.MatchProto(expectedFilter))
			})
		})

		Context("when only minimal required fields are set", func() {
			BeforeEach(func() {
				listener.Options.AdaptiveConcurrency = &adaptive_concurrency.FilterConfig{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttCalcParams: &adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams{
						IntervalMillis: 500,
						// All other fields are optional and not set
					},
				}
			})

			It("should create filter with defaults for all optional fields", func() {
				filters, err := p.HttpFilters(params, listener)
				Expect(err).NotTo(HaveOccurred())
				Expect(filters).To(HaveLen(1))

				filter := filters[0]
				expectedConfig := createBaseExpectedConfig(baseConfigParams{
					ConcurrencyUpdateIntervalMillis: 1000,
					MinRttIntervalMillis:            500,
				})

				typedConfig, err := utils.MessageToAny(expectedConfig)
				Expect(err).NotTo(HaveOccurred())

				expectedFilter := &envoyhcm.HttpFilter{
					Name: "envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency",
					ConfigType: &envoyhcm.HttpFilter_TypedConfig{
						TypedConfig: typedConfig,
					},
				}

				Expect(filter.Filter).To(matchers.MatchProto(expectedFilter))
			})
		})
	})

})
