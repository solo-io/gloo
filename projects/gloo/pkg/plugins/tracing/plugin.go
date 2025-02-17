package tracing

import (
	"strings"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_config_trace_v3 "github.com/envoyproxy/go-control-plane/envoy/config/trace/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_type_metadata_v3 "github.com/envoyproxy/go-control-plane/envoy/type/metadata/v3"
	envoytracing "github.com/envoyproxy/go-control-plane/envoy/type/tracing/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/trace/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tracing"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/internal/common"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	translatorutil "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"istio.io/istio/pkg/slices"
)

var (
	_ plugins.Plugin                      = new(plugin)
	_ plugins.HttpConnectionManagerPlugin = new(plugin)
	_ plugins.RoutePlugin                 = new(plugin)
)

const (
	ExtensionName = "tracing"

	// default all tracing percentages to 100%
	oneHundredPercent float32 = 100.0
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(_ plugins.InitParams) {
}

// Manage the tracing portion of the HCM settings
func (p *plugin) ProcessHcmNetworkFilter(params plugins.Params, parent *v1.Listener, listener *v1.HttpListener, out *envoyhttp.HttpConnectionManager) error {
	// only apply tracing config to the listener is using the HCM plugin
	in := listener.GetOptions().GetHttpConnectionManagerSettings()
	if in == nil {
		return nil
	}

	tracingSettings := in.GetTracing()
	if tracingSettings == nil {
		return nil
	}

	// this plugin will overwrite any prior tracing config
	trCfg := &envoyhttp.HttpConnectionManager_Tracing{}

	customTags := customTags(tracingSettings)
	trCfg.CustomTags = customTags
	trCfg.Verbose = tracingSettings.GetVerbose().GetValue()

	tracingProvider, err := processEnvoyTracingProvider(params, tracingSettings, parent)
	if err != nil {
		return err
	}
	trCfg.Provider = tracingProvider

	// Gloo configures envoy as an ingress, rather than an egress
	// 06/2020 removing below- OperationName field is being deprecated, and we set it to the default value anyway
	// trCfg.OperationName = envoyhttp.HttpConnectionManager_Tracing_INGRESS
	if percentages := tracingSettings.GetTracePercentages(); percentages != nil {
		trCfg.ClientSampling = envoySimplePercentWithDefault(percentages.GetClientSamplePercentage(), oneHundredPercent)
		trCfg.RandomSampling = envoySimplePercentWithDefault(percentages.GetRandomSamplePercentage(), oneHundredPercent)
		trCfg.OverallSampling = envoySimplePercentWithDefault(percentages.GetOverallSamplePercentage(), oneHundredPercent)
	} else {
		trCfg.ClientSampling = envoySimplePercent(oneHundredPercent)
		trCfg.RandomSampling = envoySimplePercent(oneHundredPercent)
		trCfg.OverallSampling = envoySimplePercent(oneHundredPercent)
	}
	trCfg.SpawnUpstreamSpan = wrapperspb.Bool(tracingSettings.GetSpawnUpstreamSpan())
	out.Tracing = trCfg

	return nil
}

func customTags(tracingSettings *tracing.ListenerTracingSettings) []*envoytracing.CustomTag {
	var customTags []*envoytracing.CustomTag

	for _, requestHeaderTag := range tracingSettings.GetRequestHeadersForTags() {
		tag := &envoytracing.CustomTag{
			Tag: requestHeaderTag.GetValue(),
			Type: &envoytracing.CustomTag_RequestHeader{
				RequestHeader: &envoytracing.CustomTag_Header{
					Name: requestHeaderTag.GetValue(),
				},
			},
		}
		customTags = append(customTags, tag)
	}
	for _, envVarTag := range tracingSettings.GetEnvironmentVariablesForTags() {
		tag := &envoytracing.CustomTag{
			Tag: envVarTag.GetTag().GetValue(),
			Type: &envoytracing.CustomTag_Environment_{
				Environment: &envoytracing.CustomTag_Environment{
					Name:         envVarTag.GetName().GetValue(),
					DefaultValue: envVarTag.GetDefaultValue().GetValue(),
				},
			},
		}
		customTags = append(customTags, tag)
	}
	for _, literalTag := range tracingSettings.GetLiteralsForTags() {
		tag := &envoytracing.CustomTag{
			Tag: literalTag.GetTag().GetValue(),
			Type: &envoytracing.CustomTag_Literal_{
				Literal: &envoytracing.CustomTag_Literal{
					Value: literalTag.GetValue().GetValue(),
				},
			},
		}
		customTags = append(customTags, tag)
	}
	for _, metadataTag := range tracingSettings.GetMetadataForTags() {
		keyDelimiter := "."
		if metadataTag.GetValue().GetNestedFieldDelimiter() != "" {
			keyDelimiter = metadataTag.GetValue().GetNestedFieldDelimiter()
		}
		tag := &envoytracing.CustomTag_Metadata{
			MetadataKey: &envoy_type_metadata_v3.MetadataKey{
				Key: metadataTag.GetValue().GetNamespace(),
				Path: slices.Map(strings.Split(metadataTag.GetValue().GetKey(), keyDelimiter), func(key string) *envoy_type_metadata_v3.MetadataKey_PathSegment {
					return &envoy_type_metadata_v3.MetadataKey_PathSegment{Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{Key: key}}
				}),
			},
			DefaultValue: metadataTag.GetDefaultValue(),
		}
		switch metadataTag.GetKind() {
		case tracing.TracingTagMetadata_REQUEST:
			tag.Kind = &envoy_type_metadata_v3.MetadataKind{
				Kind: &envoy_type_metadata_v3.MetadataKind_Request_{Request: &envoy_type_metadata_v3.MetadataKind_Request{}},
			}
		case tracing.TracingTagMetadata_ENDPOINT:
			tag.Kind = &envoy_type_metadata_v3.MetadataKind{
				Kind: &envoy_type_metadata_v3.MetadataKind_Host_{Host: &envoy_type_metadata_v3.MetadataKind_Host{}},
			}

		}
		customTags = append(customTags, &envoytracing.CustomTag{
			Tag: metadataTag.GetTag(),
			Type: &envoytracing.CustomTag_Metadata_{
				Metadata: tag,
			},
		})
	}

	return customTags
}

func processEnvoyTracingProvider(
	params plugins.Params,
	tracingSettings *tracing.ListenerTracingSettings,
	parent *v1.Listener,
) (*envoy_config_trace_v3.Tracing_Http, error) {
	if tracingSettings.GetProviderConfig() == nil {
		return nil, nil
	}

	snapshot := params.Snapshot

	switch typed := tracingSettings.GetProviderConfig().(type) {
	case *tracing.ListenerTracingSettings_ZipkinConfig:
		return processEnvoyZipkinTracing(snapshot, typed)

	case *tracing.ListenerTracingSettings_DatadogConfig:
		return processEnvoyDatadogTracing(snapshot, typed)

	case *tracing.ListenerTracingSettings_OpenTelemetryConfig:
		return processEnvoyOpenTelemetryTracing(params, typed, parent)

	default:
		return nil, errors.Errorf("Unsupported Tracing.ProviderConfiguration: %v", typed)
	}
}

func processEnvoyZipkinTracing(
	snapshot *v1snap.ApiSnapshot,
	zipkinTracingSettings *tracing.ListenerTracingSettings_ZipkinConfig,
) (*envoy_config_trace_v3.Tracing_Http, error) {
	var collectorClusterName string

	switch collectorCluster := zipkinTracingSettings.ZipkinConfig.GetCollectorCluster().(type) {
	case *v3.ZipkinConfig_CollectorUpstreamRef:
		// Support upstreams as the collector cluster
		var err error
		collectorClusterName, err = getEnvoyTracingCollectorClusterName(snapshot, collectorCluster.CollectorUpstreamRef)
		if err != nil {
			return nil, err
		}
	case *v3.ZipkinConfig_ClusterName:
		// Support static clusters as the collector cluster
		collectorClusterName = collectorCluster.ClusterName
	}

	envoyConfig, err := api_conversion.ToEnvoyZipkinConfiguration(zipkinTracingSettings.ZipkinConfig, collectorClusterName)
	if err != nil {
		return nil, err
	}

	marshalledEnvoyConfig, err := ptypes.MarshalAny(envoyConfig)
	if err != nil {
		return nil, err
	}

	return &envoy_config_trace_v3.Tracing_Http{
		Name: "envoy.tracers.zipkin",
		ConfigType: &envoy_config_trace_v3.Tracing_Http_TypedConfig{
			TypedConfig: marshalledEnvoyConfig,
		},
	}, nil
}

func processEnvoyDatadogTracing(
	snapshot *v1snap.ApiSnapshot,
	datadogTracingSettings *tracing.ListenerTracingSettings_DatadogConfig,
) (*envoy_config_trace_v3.Tracing_Http, error) {
	var collectorClusterName string

	switch collectorCluster := datadogTracingSettings.DatadogConfig.GetCollectorCluster().(type) {
	case *v3.DatadogConfig_CollectorUpstreamRef:
		// Support upstreams as the collector cluster
		var err error
		collectorClusterName, err = getEnvoyTracingCollectorClusterName(snapshot, collectorCluster.CollectorUpstreamRef)
		if err != nil {
			return nil, err
		}
	case *v3.DatadogConfig_ClusterName:
		// Support static clusters as the collector cluster
		collectorClusterName = collectorCluster.ClusterName
	default:
		return nil, errors.Errorf("Unsupported Tracing.ProviderConfiguration: %v", collectorCluster)
	}

	envoyConfig, err := api_conversion.ToEnvoyDatadogConfiguration(datadogTracingSettings.DatadogConfig, collectorClusterName)
	if err != nil {
		return nil, err
	}

	marshalledEnvoyConfig, err := ptypes.MarshalAny(envoyConfig)
	if err != nil {
		return nil, err
	}

	return &envoy_config_trace_v3.Tracing_Http{
		Name: "envoy.tracers.datadog",
		ConfigType: &envoy_config_trace_v3.Tracing_Http_TypedConfig{
			TypedConfig: marshalledEnvoyConfig,
		},
	}, nil
}

func processEnvoyOpenTelemetryTracing(
	params plugins.Params,
	openTelemetryTracingSettings *tracing.ListenerTracingSettings_OpenTelemetryConfig,
	parent *v1.Listener,
) (*envoy_config_trace_v3.Tracing_Http, error) {
	var collectorClusterName string

	snapshot := params.Snapshot
	switch collectorCluster := openTelemetryTracingSettings.OpenTelemetryConfig.GetCollectorCluster().(type) {
	case *v3.OpenTelemetryConfig_CollectorUpstreamRef:
		// Support upstreams as the collector cluster
		var err error
		collectorClusterName, err = getEnvoyTracingCollectorClusterName(snapshot, collectorCluster.CollectorUpstreamRef)
		if err != nil {
			return nil, err
		}
	case *v3.OpenTelemetryConfig_ClusterName:
		// Support static clusters as the collector cluster
		collectorClusterName = collectorCluster.ClusterName
	default:
		return nil, errors.Errorf("Unsupported Tracing.ProviderConfiguration: %v", collectorCluster)
	}

	serviceName := openTelemetryTracingSettings.OpenTelemetryConfig.GetServiceName()
	if serviceName == "" {
		serviceName = api_conversion.GetGatewayNameFromParent(params.Ctx, parent)
	}

	authority := ""
	grpcService := openTelemetryTracingSettings.OpenTelemetryConfig.GetGrpcService()
	if grpcService != nil {
		authority = grpcService.GetAuthority()
	}

	envoyConfig := api_conversion.ToEnvoyOpenTelemetryConfiguration(collectorClusterName, serviceName, authority)
	marshalledEnvoyConfig, err := anypb.New(envoyConfig)
	if err != nil {
		return nil, err
	}

	return &envoy_config_trace_v3.Tracing_Http{
		Name: "envoy.tracers.opentelemetry",
		ConfigType: &envoy_config_trace_v3.Tracing_Http_TypedConfig{
			TypedConfig: marshalledEnvoyConfig,
		},
	}, nil
}

func getEnvoyTracingCollectorClusterName(snapshot *v1snap.ApiSnapshot, collectorUpstreamRef *core.ResourceRef) (string, error) {
	if snapshot == nil {
		return "", errors.Errorf("Invalid Snapshot (nil provided)")
	}

	if collectorUpstreamRef == nil {
		return "", errors.Errorf("Invalid CollectorUpstreamRef (nil ref provided)")
	}

	// Make sure the upstream exists
	_, err := snapshot.Upstreams.Find(collectorUpstreamRef.GetNamespace(), collectorUpstreamRef.GetName())
	if err != nil {
		return "", pluginutils.NewUpstreamNotFoundErr(collectorUpstreamRef)
	}

	return translatorutil.UpstreamToClusterName(collectorUpstreamRef), nil
}

func envoySimplePercent(numerator float32) *envoy_type.Percent {
	return &envoy_type.Percent{Value: float64(numerator)}
}

// use FloatValue to detect when nil (avoids error-prone float comparisons)
func envoySimplePercentWithDefault(numerator *wrappers.FloatValue, defaultValue float32) *envoy_type.Percent {
	if numerator == nil {
		return envoySimplePercent(defaultValue)
	}
	return envoySimplePercent(numerator.GetValue())
}

func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	if in.GetOptions() == nil || in.GetOptions().GetTracing() == nil {
		return nil
	}
	if percentages := in.GetOptions().GetTracing().GetTracePercentages(); percentages != nil {
		out.Tracing = &envoy_config_route_v3.Tracing{
			ClientSampling:  common.ToEnvoyPercentageWithDefault(percentages.GetClientSamplePercentage(), oneHundredPercent),
			RandomSampling:  common.ToEnvoyPercentageWithDefault(percentages.GetRandomSamplePercentage(), oneHundredPercent),
			OverallSampling: common.ToEnvoyPercentageWithDefault(percentages.GetOverallSamplePercentage(), oneHundredPercent),
		}
	} else {
		out.Tracing = &envoy_config_route_v3.Tracing{
			ClientSampling:  common.ToEnvoyPercentage(oneHundredPercent),
			RandomSampling:  common.ToEnvoyPercentage(oneHundredPercent),
			OverallSampling: common.ToEnvoyPercentage(oneHundredPercent),
		}
	}
	descriptor := in.GetOptions().GetTracing().GetRouteDescriptor()
	if descriptor != "" {
		out.Decorator = &envoy_config_route_v3.Decorator{
			Operation: descriptor,
			Propagate: in.GetOptions().GetTracing().GetPropagate(),
		}
	}
	return nil
}
