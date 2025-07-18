package adaptiveconcurrency

import (
	"fmt"

	envoy_adaptive_concurrency_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/adaptive_concurrency/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/adaptive_concurrency"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var (
	pluginStage                               = plugins.DuringStage(plugins.RateLimitStage)
	ErrConcurrencyUpdateIntervalMillisMissing = fmt.Errorf("concurrency_update_interval_millis is required")
	ErrMinRttCalcParamsMissing                = fmt.Errorf("min_rtt_calc_params is required")
	ErrIntervalMissing                        = fmt.Errorf("Either interval_millis or fixed_value_millis must be set")
)

const (
	ExtensionName = "envoy.extensions.filters.http.adaptive_concurrency.v3.AdaptiveConcurrency"

	// Default values as documented in the proto file
	// These match the defaults in the envoy docs: https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/adaptive_concurrency/v3/adaptive_concurrency.proto#envoy-v3-api-msg-extensions-filters-http-adaptive-concurrency-v3-adaptiveconcurrency
	// These are set in the plugin to avoid drift in our documentation in the Envoy values change.
	DefaultSampleAggregatePercentile      = 50.0
	DefaultMaxConcurrencyLimit            = 1000
	DefaultRequestCount                   = 50
	DefaultMinConcurrency                 = 3
	DefaultJitterPercentile               = 15.0
	DefaultBufferPercentile               = 25.0
	DefaultConcurrencyLimitExceededStatus = 503
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}
func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	in := listener.GetOptions()
	adaptiveConcurrencyConfig, err := translateAdaptiveConcurrency(in)

	if err != nil {
		return nil, err
	}

	if adaptiveConcurrencyConfig == nil {
		return []plugins.StagedHttpFilter{}, nil
	}

	return []plugins.StagedHttpFilter{plugins.MustNewStagedFilter(ExtensionName, adaptiveConcurrencyConfig, pluginStage)}, nil
}

func translateAdaptiveConcurrency(in *v1.HttpListenerOptions) (*envoy_adaptive_concurrency_v3.AdaptiveConcurrency, error) {
	adaptiveConcurrency := in.GetAdaptiveConcurrency()

	if adaptiveConcurrency == nil {
		return nil, nil
	}

	concurrencyLimitParams, err := translateConcurrencyLimitParams(adaptiveConcurrency)
	if err != nil {
		return nil, err
	}

	minRttCalcParams, err := translateMinRttCalcParams(adaptiveConcurrency.GetMinRttCalcParams())
	if err != nil {
		return nil, err
	}

	// Variable declarations with defaults
	var (
		sampleAggregatePercentile      float64
		concurrencyLimitExceededStatus uint32
	)

	// Set sample_aggregate_percentile with default
	if adaptiveConcurrency.GetSampleAggregatePercentile() != nil {
		sampleAggregatePercentile = adaptiveConcurrency.GetSampleAggregatePercentile().GetValue()
	} else {
		sampleAggregatePercentile = DefaultSampleAggregatePercentile
	}

	// Set concurrency limit exceeded status with default of 503 for non-error codes
	status := adaptiveConcurrency.GetConcurrencyLimitExceededStatus()
	if status == 0 || status < 400 {
		concurrencyLimitExceededStatus = DefaultConcurrencyLimitExceededStatus
	} else {
		concurrencyLimitExceededStatus = status
	}

	gradientControllerConfig := &envoy_adaptive_concurrency_v3.GradientControllerConfig{
		ConcurrencyLimitParams:    concurrencyLimitParams,
		MinRttCalcParams:          minRttCalcParams,
		SampleAggregatePercentile: &typev3.Percent{Value: sampleAggregatePercentile},
	}

	out := &envoy_adaptive_concurrency_v3.AdaptiveConcurrency{
		ConcurrencyControllerConfig: &envoy_adaptive_concurrency_v3.AdaptiveConcurrency_GradientControllerConfig{
			GradientControllerConfig: gradientControllerConfig,
		},
		ConcurrencyLimitExceededStatus: &typev3.HttpStatus{
			Code: typev3.StatusCode(concurrencyLimitExceededStatus),
		},
	}

	return out, nil
}

func translateConcurrencyLimitParams(in *adaptive_concurrency.FilterConfig) (*envoy_adaptive_concurrency_v3.GradientControllerConfig_ConcurrencyLimitCalculationParams, error) {
	if in.GetConcurrencyUpdateIntervalMillis() == 0 {
		return nil, ErrConcurrencyUpdateIntervalMillisMissing
	}

	// Variable declarations with defaults
	var maxConcurrencyLimit uint32

	// Set max_concurrency_limit with default
	if in.GetMaxConcurrencyLimit() != nil {
		maxConcurrencyLimit = in.GetMaxConcurrencyLimit().GetValue()
	} else {
		maxConcurrencyLimit = DefaultMaxConcurrencyLimit
	}

	out := &envoy_adaptive_concurrency_v3.GradientControllerConfig_ConcurrencyLimitCalculationParams{
		ConcurrencyUpdateInterval: millisToDuration(in.GetConcurrencyUpdateIntervalMillis()),
		MaxConcurrencyLimit:       &wrapperspb.UInt32Value{Value: maxConcurrencyLimit},
	}

	return out, nil
}

func translateMinRttCalcParams(in *adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams) (*envoy_adaptive_concurrency_v3.GradientControllerConfig_MinimumRTTCalculationParams, error) {
	if in == nil {
		return nil, ErrMinRttCalcParamsMissing
	}

	intervalMillis := in.GetIntervalMillis()
	fixedValueMillis := in.GetFixedValueMillis()

	// If both are set, Interval is used
	if intervalMillis == 0 && fixedValueMillis == 0 {
		return nil, ErrIntervalMissing
	}

	// Variable declarations with defaults
	var (
		requestCount     uint32
		minConcurrency   uint32
		jitterPercentile float64
		bufferPercentile float64
	)

	// Set request_count with default
	if in.GetRequestCount() != nil {
		requestCount = in.GetRequestCount().GetValue()
	} else {
		requestCount = DefaultRequestCount
	}

	// Set min_concurrency with default
	if in.GetMinConcurrency() != nil {
		minConcurrency = in.GetMinConcurrency().GetValue()
	} else {
		minConcurrency = DefaultMinConcurrency
	}

	// Set jitter_percentile with default
	if in.GetJitterPercentile() != nil {
		jitterPercentile = in.GetJitterPercentile().GetValue()
	} else {
		jitterPercentile = DefaultJitterPercentile
	}

	// Set buffer_percentile with default
	if in.GetBufferPercentile() != nil {
		bufferPercentile = in.GetBufferPercentile().GetValue()
	} else {
		bufferPercentile = DefaultBufferPercentile
	}

	out := &envoy_adaptive_concurrency_v3.GradientControllerConfig_MinimumRTTCalculationParams{
		RequestCount:   &wrapperspb.UInt32Value{Value: requestCount},
		MinConcurrency: &wrapperspb.UInt32Value{Value: minConcurrency},
		Jitter:         &typev3.Percent{Value: jitterPercentile},
		Buffer:         &typev3.Percent{Value: bufferPercentile},
	}

	if intervalMillis > 0 {
		out.Interval = millisToDuration(intervalMillis)
	} else {
		out.FixedValue = millisToDuration(fixedValueMillis)
	}

	return out, nil
}

// Convert milliseconds to durationpb.Duration
func millisToDuration(millis uint32) *durationpb.Duration {
	nanos := millis % 1000 * 1_000_000
	seconds := millis / 1000
	return &durationpb.Duration{
		Seconds: int64(seconds),
		Nanos:   int32(nanos),
	}
}
