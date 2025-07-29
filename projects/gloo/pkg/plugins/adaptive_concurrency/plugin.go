package adaptiveconcurrency

import (
	"context"
	"fmt"
	"net/http"

	envoy_adaptive_concurrency_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/adaptive_concurrency/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/solo-io/gloo/pkg/utils/protoutils/duration"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/adaptive_concurrency"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/go-utils/contextutils"
)

var (
	pluginStage                          = plugins.DuringStage(plugins.RateLimitStage)
	ErrConcurrencyLimitCalcParamsMissing = fmt.Errorf("concurrency_limit_calc_params is required")
	ErrMinRttCalcParamsMissing           = fmt.Errorf("min_rtt_calc_params is required")
	ErrIntervalMissing                   = fmt.Errorf("Either interval or fixed_value must be set")
	ErrConcurrencyUpdateIntervalMissing  = fmt.Errorf("concurrency_update_interval is required")
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
	DefaultConcurrencyLimitExceededStatus = uint32(http.StatusServiceUnavailable)
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

	concurrencyLimitParams, err := translateConcurrencyLimitParams(adaptiveConcurrency.GetConcurrencyLimitCalculationParams())
	if err != nil {
		return nil, err
	}

	minRttCalcParams, err := translateMinRttCalcParams(adaptiveConcurrency.GetMinRttCalculationParams())
	if err != nil {
		return nil, err
	}

	// Variable declarations with defaults
	var (
		sampleAggregatePercentile      float64 = DefaultSampleAggregatePercentile
		concurrencyLimitExceededStatus uint32  = DefaultConcurrencyLimitExceededStatus
	)

	// Set sample_aggregate_percentile with default
	if adaptiveConcurrency.GetSampleAggregatePercentile() != nil {
		sampleAggregatePercentile = adaptiveConcurrency.GetSampleAggregatePercentile().GetValue()
	}

	// Set concurrency limit exceeded status with default of 503 for non-error codes
	status := adaptiveConcurrency.GetConcurrencyLimitExceededStatus()
	if status >= 400 {
		concurrencyLimitExceededStatus = status
	} else {
		if status != 0 { // Don't log if status is unset
			// The envoy treats this as a non-error condition, and that logic will be applied here as well.
			// However, we want to log this as a warning, so that users are aware of the override of their setting.
			logger := contextutils.LoggerFrom(context.Background())
			logger.Warnf("concurrencyLimitExceededStatus is %d, which is not an error code. The default setting of 503 Unavailable will be applied.", status)
		}
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

func translateConcurrencyLimitParams(in *adaptive_concurrency.FilterConfig_ConcurrencyLimitCalculationParams) (*envoy_adaptive_concurrency_v3.GradientControllerConfig_ConcurrencyLimitCalculationParams, error) {
	if in == nil {
		return nil, ErrConcurrencyLimitCalcParamsMissing
	}

	// Variable declarations with defaults
	var (
		maxConcurrencyLimit       uint32 = DefaultMaxConcurrencyLimit
		concurrencyUpdateInterval uint32
	)

	// Set concurrency_update_interval with default
	if in.GetConcurrencyUpdateInterval() != 0 {
		concurrencyUpdateInterval = in.GetConcurrencyUpdateInterval()
	} else {
		return nil, ErrConcurrencyUpdateIntervalMissing
	}

	// Set max_concurrency_limit with default
	if in.GetMaxConcurrencyLimit() != nil {
		maxConcurrencyLimit = in.GetMaxConcurrencyLimit().GetValue()
	}

	out := &envoy_adaptive_concurrency_v3.GradientControllerConfig_ConcurrencyLimitCalculationParams{
		ConcurrencyUpdateInterval: duration.MillisToDuration(concurrencyUpdateInterval),
		MaxConcurrencyLimit:       &wrapperspb.UInt32Value{Value: maxConcurrencyLimit},
	}

	return out, nil
}

func translateMinRttCalcParams(in *adaptive_concurrency.FilterConfig_MinRoundtripTimeCalculationParams) (*envoy_adaptive_concurrency_v3.GradientControllerConfig_MinimumRTTCalculationParams, error) {
	if in == nil {
		return nil, ErrMinRttCalcParamsMissing
	}

	interval := in.GetInterval()
	fixedValue := in.GetFixedValue()

	// If both are set, Interval is used
	if interval == 0 && fixedValue == 0 {
		return nil, ErrIntervalMissing
	}

	// Variable declarations with defaults
	var (
		requestCount     uint32  = DefaultRequestCount
		minConcurrency   uint32  = DefaultMinConcurrency
		jitterPercentile float64 = DefaultJitterPercentile
		bufferPercentile float64 = DefaultBufferPercentile
	)

	// Set request_count with default
	if in.GetRequestCount() != nil {
		requestCount = in.GetRequestCount().GetValue()
	}

	// Set min_concurrency with default
	if in.GetMinConcurrency() != nil {
		minConcurrency = in.GetMinConcurrency().GetValue()
	}

	// Set jitter_percentile with default
	if in.GetJitterPercentile() != nil {
		jitterPercentile = in.GetJitterPercentile().GetValue()
	}

	// Set buffer_percentile with default
	if in.GetBufferPercentile() != nil {
		bufferPercentile = in.GetBufferPercentile().GetValue()
	}

	out := &envoy_adaptive_concurrency_v3.GradientControllerConfig_MinimumRTTCalculationParams{
		RequestCount:   &wrapperspb.UInt32Value{Value: requestCount},
		MinConcurrency: &wrapperspb.UInt32Value{Value: minConcurrency},
		Jitter:         &typev3.Percent{Value: jitterPercentile},
		Buffer:         &typev3.Percent{Value: bufferPercentile},
	}

	if interval > 0 {
		out.Interval = duration.MillisToDuration(interval)
	} else {
		out.FixedValue = duration.MillisToDuration(fixedValue)
	}

	return out, nil
}
