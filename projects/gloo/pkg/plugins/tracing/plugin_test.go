package tracing

import (
	v12 "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoytrace "github.com/envoyproxy/go-control-plane/envoy/config/trace/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoy_type_metadata_v3 "github.com/envoyproxy/go-control-plane/envoy/type/metadata/v3"
	envoytracing "github.com/envoyproxy/go-control-plane/envoy/type/tracing/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	envoytrace_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/trace/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tracing"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Plugin", func() {

	var (
		plugin       plugins.Plugin
		pluginParams plugins.Params

		hcmSettings *hcm.HttpConnectionManagerSettings
	)

	BeforeEach(func() {
		plugin = NewPlugin()
		pluginParams = plugins.Params{
			Snapshot: nil,
		}
		hcmSettings = &hcm.HttpConnectionManagerSettings{}
	})

	processHcmNetworkFilter := func(cfg *envoyhttp.HttpConnectionManager) error {
		httpListener := &v1.HttpListener{
			Options: &v1.HttpListenerOptions{
				HttpConnectionManagerSettings: hcmSettings,
			},
		}
		listener := &v1.Listener{
			OpaqueMetadata: &v1.Listener_MetadataStatic{
				MetadataStatic: &v1.SourceMetadata{
					Sources: []*v1.SourceMetadata_SourceRef{
						{
							ResourceRef: &core.ResourceRef{
								Name:      "delegate-1",
								Namespace: "gloo-system",
							},
							ResourceKind:       "*v1.RouteTable",
							ObservedGeneration: 0,
						},
						{
							ResourceRef: &core.ResourceRef{
								Name:      "gateway",
								Namespace: "gloo-system",
							},
							ResourceKind:       "*v1.Gateway",
							ObservedGeneration: 0,
						},
					},
				},
			},
		}
		return plugin.(plugins.HttpConnectionManagerPlugin).ProcessHcmNetworkFilter(pluginParams, listener, httpListener, cfg)
	}

	It("should update listener properly", func() {
		cfg := &envoyhttp.HttpConnectionManager{}
		hcmSettings = &hcm.HttpConnectionManagerSettings{
			Tracing: &tracing.ListenerTracingSettings{
				RequestHeadersForTags: []*wrappers.StringValue{{Value: "header1"}, {Value: "header2"}},
				EnvironmentVariablesForTags: []*tracing.TracingTagEnvironmentVariable{
					{
						Tag:  &wrappers.StringValue{Value: "k8s.pod.name"},
						Name: &wrappers.StringValue{Value: "POD_NAME"},
					},
					{
						Tag:          &wrappers.StringValue{Value: "k8s.pod.ip"},
						Name:         &wrappers.StringValue{Value: "POD_IP"},
						DefaultValue: &wrappers.StringValue{Value: "NO_POD_IP"},
					},
				},
				LiteralsForTags: []*tracing.TracingTagLiteral{
					{
						Tag:   &wrappers.StringValue{Value: "foo"},
						Value: &wrappers.StringValue{Value: "bar"},
					},
				},
				MetadataForTags: []*tracing.TracingTagMetadata{
					{
						Tag:  "envoy.metadata.foo",
						Kind: tracing.TracingTagMetadata_REQUEST,
						Value: &tracing.TracingTagMetadata_MetadataValue{
							Namespace: "namespace",
							Key:       "nested.key",
						},
					},
					{
						Tag:  "envoy.metadata.bar",
						Kind: tracing.TracingTagMetadata_ENDPOINT,
						Value: &tracing.TracingTagMetadata_MetadataValue{
							Namespace: "namespace",
							Key:       "nested.key",
						},
					},
					{
						Tag:          "envoy.metadata.baz",
						Kind:         tracing.TracingTagMetadata_REQUEST,
						DefaultValue: "default",
						Value: &tracing.TracingTagMetadata_MetadataValue{
							Namespace:            "namespace",
							Key:                  "nested:key",
							NestedFieldDelimiter: ":",
						},
					},
				},
				Verbose: &wrappers.BoolValue{Value: true},
				TracePercentages: &tracing.TracePercentages{
					ClientSamplePercentage:  &wrappers.FloatValue{Value: 10},
					RandomSamplePercentage:  &wrappers.FloatValue{Value: 20},
					OverallSamplePercentage: &wrappers.FloatValue{Value: 30},
				},
				ProviderConfig: nil,
			},
		}

		err := processHcmNetworkFilter(cfg)
		Expect(err).NotTo(HaveOccurred())
		expected := &envoyhttp.HttpConnectionManager{
			Tracing: &envoyhttp.HttpConnectionManager_Tracing{
				SpawnUpstreamSpan: &wrappers.BoolValue{Value: false},
				CustomTags: []*envoytracing.CustomTag{
					{
						Tag: "header1",
						Type: &envoytracing.CustomTag_RequestHeader{
							RequestHeader: &envoytracing.CustomTag_Header{
								Name: "header1",
							},
						},
					},
					{
						Tag: "header2",
						Type: &envoytracing.CustomTag_RequestHeader{
							RequestHeader: &envoytracing.CustomTag_Header{
								Name: "header2",
							},
						},
					},
					{
						Tag: "k8s.pod.name",
						Type: &envoytracing.CustomTag_Environment_{
							Environment: &envoytracing.CustomTag_Environment{
								Name: "POD_NAME",
							},
						},
					},
					{
						Tag: "k8s.pod.ip",
						Type: &envoytracing.CustomTag_Environment_{
							Environment: &envoytracing.CustomTag_Environment{
								Name:         "POD_IP",
								DefaultValue: "NO_POD_IP",
							},
						},
					},
					{
						Tag: "foo",
						Type: &envoytracing.CustomTag_Literal_{
							Literal: &envoytracing.CustomTag_Literal{
								Value: "bar",
							},
						},
					},
					{
						Tag: "envoy.metadata.foo",
						Type: &envoytracing.CustomTag_Metadata_{
							Metadata: &envoytracing.CustomTag_Metadata{
								MetadataKey: &envoy_type_metadata_v3.MetadataKey{
									Key: "namespace",
									Path: []*envoy_type_metadata_v3.MetadataKey_PathSegment{
										{
											Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{
												Key: "nested",
											},
										},
										{
											Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{
												Key: "key",
											},
										},
									},
								},
								Kind: &envoy_type_metadata_v3.MetadataKind{
									Kind: &envoy_type_metadata_v3.MetadataKind_Request_{
										Request: &envoy_type_metadata_v3.MetadataKind_Request{},
									},
								},
							},
						},
					},
					{
						Tag: "envoy.metadata.bar",
						Type: &envoytracing.CustomTag_Metadata_{
							Metadata: &envoytracing.CustomTag_Metadata{
								MetadataKey: &envoy_type_metadata_v3.MetadataKey{
									Key: "namespace",
									Path: []*envoy_type_metadata_v3.MetadataKey_PathSegment{
										{
											Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{
												Key: "nested",
											},
										},
										{
											Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{
												Key: "key",
											},
										},
									},
								},
								Kind: &envoy_type_metadata_v3.MetadataKind{
									Kind: &envoy_type_metadata_v3.MetadataKind_Host_{
										Host: &envoy_type_metadata_v3.MetadataKind_Host{},
									},
								},
							},
						},
					},
					{
						Tag: "envoy.metadata.baz",
						Type: &envoytracing.CustomTag_Metadata_{
							Metadata: &envoytracing.CustomTag_Metadata{
								MetadataKey: &envoy_type_metadata_v3.MetadataKey{
									Key: "namespace",
									Path: []*envoy_type_metadata_v3.MetadataKey_PathSegment{
										{
											Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{
												Key: "nested",
											},
										},
										{
											Segment: &envoy_type_metadata_v3.MetadataKey_PathSegment_Key{
												Key: "key",
											},
										},
									},
								},
								DefaultValue: "default",
								Kind: &envoy_type_metadata_v3.MetadataKind{
									Kind: &envoy_type_metadata_v3.MetadataKind_Request_{
										Request: &envoy_type_metadata_v3.MetadataKind_Request{},
									},
								},
							},
						},
					},
				},
				ClientSampling:  &envoy_type.Percent{Value: 10},
				RandomSampling:  &envoy_type.Percent{Value: 20},
				OverallSampling: &envoy_type.Percent{Value: 30},
				Verbose:         true,
				Provider:        nil,
			},
		}
		Expect(cfg).To(Equal(expected))
	})

	It("should update listener properly - with defaults", func() {
		cfg := &envoyhttp.HttpConnectionManager{}
		hcmSettings = &hcm.HttpConnectionManagerSettings{
			Tracing: &tracing.ListenerTracingSettings{},
		}

		err := processHcmNetworkFilter(cfg)
		Expect(err).NotTo(HaveOccurred())
		expected := &envoyhttp.HttpConnectionManager{
			Tracing: &envoyhttp.HttpConnectionManager_Tracing{
				ClientSampling:    &envoy_type.Percent{Value: 100},
				RandomSampling:    &envoy_type.Percent{Value: 100},
				OverallSampling:   &envoy_type.Percent{Value: 100},
				Verbose:           false,
				Provider:          nil,
				SpawnUpstreamSpan: &wrappers.BoolValue{Value: false},
			},
		}
		Expect(cfg).To(Equal(expected))
	})

	It("should properly set spawn_upstream_span", func() {
		cfg := &envoyhttp.HttpConnectionManager{}
		hcmSettings = &hcm.HttpConnectionManagerSettings{
			Tracing: &tracing.ListenerTracingSettings{
				SpawnUpstreamSpan: true,
			},
		}

		err := processHcmNetworkFilter(cfg)
		Expect(err).NotTo(HaveOccurred())
		expected := &envoyhttp.HttpConnectionManager{
			Tracing: &envoyhttp.HttpConnectionManager_Tracing{
				ClientSampling:    &envoy_type.Percent{Value: 100},
				RandomSampling:    &envoy_type.Percent{Value: 100},
				OverallSampling:   &envoy_type.Percent{Value: 100},
				Verbose:           false,
				Provider:          nil,
				SpawnUpstreamSpan: &wrappers.BoolValue{Value: true},
			},
		}
		Expect(cfg).To(Equal(expected))
	})

	Context("should handle tracing provider config", func() {

		It("when provider config is nil", func() {
			cfg := &envoyhttp.HttpConnectionManager{}
			hcmSettings = &hcm.HttpConnectionManagerSettings{
				Tracing: &tracing.ListenerTracingSettings{
					ProviderConfig: nil,
				},
			}
			err := processHcmNetworkFilter(cfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Tracing.Provider).To(BeNil())
		})

		Describe("when zipkin provider config", func() {
			It("references invalid upstream", func() {
				pluginParams = plugins.Params{
					Snapshot: &v1snap.ApiSnapshot{
						Upstreams: v1.UpstreamList{
							// No valid upstreams
						},
					},
				}
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_ZipkinConfig{
							ZipkinConfig: &envoytrace_gloo.ZipkinConfig{
								CollectorCluster: &envoytrace_gloo.ZipkinConfig_CollectorUpstreamRef{
									CollectorUpstreamRef: &core.ResourceRef{
										Name:      "invalid-name",
										Namespace: "invalid-namespace",
									},
								},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).To(HaveOccurred())
			})

			It("references valid upstream", func() {
				us := v1.NewUpstream("default", "valid")
				pluginParams = plugins.Params{
					Snapshot: &v1snap.ApiSnapshot{
						Upstreams: v1.UpstreamList{us},
					},
				}

				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_ZipkinConfig{
							ZipkinConfig: &envoytrace_gloo.ZipkinConfig{
								CollectorCluster: &envoytrace_gloo.ZipkinConfig_CollectorUpstreamRef{
									CollectorUpstreamRef: &core.ResourceRef{
										Name:      "valid",
										Namespace: "default",
									},
								},
								CollectorEndpoint:        "/api/v2/spans",
								CollectorEndpointVersion: envoytrace_gloo.ZipkinConfig_HTTP_JSON,
								SharedSpanContext:        nil,
								TraceId_128Bit:           &wrappers.BoolValue{Value: false},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.ZipkinConfig{
					CollectorCluster:         "valid_default",
					CollectorEndpoint:        "/api/v2/spans",
					CollectorEndpointVersion: envoytrace.ZipkinConfig_HTTP_JSON,
					SharedSpanContext:        nil,
					TraceId_128Bit:           false,
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.zipkin",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("references cluster name", func() {
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_ZipkinConfig{
							ZipkinConfig: &envoytrace_gloo.ZipkinConfig{
								CollectorCluster: &envoytrace_gloo.ZipkinConfig_ClusterName{
									ClusterName: "zipkin-cluster-name",
								},
								CollectorEndpoint:        "/api/v2/spans",
								CollectorEndpointVersion: envoytrace_gloo.ZipkinConfig_HTTP_JSON,
								SharedSpanContext:        nil,
								TraceId_128Bit:           &wrappers.BoolValue{Value: false},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.ZipkinConfig{
					CollectorCluster:         "zipkin-cluster-name",
					CollectorEndpoint:        "/api/v2/spans",
					CollectorEndpointVersion: envoytrace.ZipkinConfig_HTTP_JSON,
					SharedSpanContext:        nil,
					TraceId_128Bit:           false,
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.zipkin",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})
		})

		Describe("when datadog provider config", func() {
			It("references invalid upstream", func() {
				pluginParams = plugins.Params{
					Snapshot: &v1snap.ApiSnapshot{
						Upstreams: v1.UpstreamList{
							// No valid upstreams
						},
					},
				}
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_CollectorUpstreamRef{
									CollectorUpstreamRef: &core.ResourceRef{
										Name:      "invalid-name",
										Namespace: "invalid-namespace",
									},
								},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).To(HaveOccurred())
			})

			It("references valid upstream", func() {
				us := v1.NewUpstream("default", "valid")
				pluginParams = plugins.Params{
					Snapshot: &v1snap.ApiSnapshot{
						Upstreams: v1.UpstreamList{us},
					},
				}
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_CollectorUpstreamRef{
									CollectorUpstreamRef: &core.ResourceRef{
										Name:      "valid",
										Namespace: "default",
									},
								},
								ServiceName: &wrappers.StringValue{Value: "datadog-gloo"},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.DatadogConfig{
					CollectorCluster: "valid_default",
					ServiceName:      "datadog-gloo",
					// This would turn on RemoteConfig when nothing is set on the gloo settings side
					// to preserve existing behavior for envoy >= v1.31
					RemoteConfig: &envoytrace.DatadogRemoteConfig{},
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.datadog",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("references cluster name", func() {
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_ClusterName{
									ClusterName: "datadog-cluster-name",
								},
								ServiceName: &wrappers.StringValue{Value: "datadog-gloo"},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.DatadogConfig{
					CollectorCluster: "datadog-cluster-name",
					ServiceName:      "datadog-gloo",
					// This would turn on RemoteConfig when nothing is set on the gloo settings side
					// to preserve existing behavior for envoy >= v1.31
					RemoteConfig: &envoytrace.DatadogRemoteConfig{},
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.datadog",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("overrides cluster hostname", func() {
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_ClusterName{
									ClusterName: "datadog-cluster-name",
								},
								ServiceName:       &wrappers.StringValue{Value: "datadog-gloo"},
								CollectorHostname: "foodog.com",
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.DatadogConfig{
					CollectorCluster:  "datadog-cluster-name",
					ServiceName:       "datadog-gloo",
					CollectorHostname: "foodog.com",
					// This would turn on RemoteConfig when nothing is set on the gloo settings side
					// to preserve existing behavior for envoy >= v1.31
					RemoteConfig: &envoytrace.DatadogRemoteConfig{},
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.datadog",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("disables remote config", func() {
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_ClusterName{
									ClusterName: "datadog-cluster-name",
								},
								ServiceName:  &wrappers.StringValue{Value: "datadog-gloo"},
								RemoteConfig: &envoytrace_gloo.DatadogRemoteConfig{Disabled: &wrappers.BoolValue{Value: true}},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.DatadogConfig{
					CollectorCluster: "datadog-cluster-name",
					ServiceName:      "datadog-gloo",
					// Note that RemoteConfig will be missing because it's disabled
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.datadog",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("sets poll interval", func() {
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_ClusterName{
									ClusterName: "datadog-cluster-name",
								},
								ServiceName:  &wrappers.StringValue{Value: "datadog-gloo"},
								RemoteConfig: &envoytrace_gloo.DatadogRemoteConfig{PollingInterval: &duration.Duration{Seconds: 12}},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.DatadogConfig{
					CollectorCluster: "datadog-cluster-name",
					ServiceName:      "datadog-gloo",
					RemoteConfig:     &envoytrace.DatadogRemoteConfig{PollingInterval: &duration.Duration{Seconds: 12}},
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.datadog",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("ignore poll interval when remote config is disabled", func() {
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_DatadogConfig{
							DatadogConfig: &envoytrace_gloo.DatadogConfig{
								CollectorCluster: &envoytrace_gloo.DatadogConfig_ClusterName{
									ClusterName: "datadog-cluster-name",
								},
								ServiceName: &wrappers.StringValue{Value: "datadog-gloo"},
								RemoteConfig: &envoytrace_gloo.DatadogRemoteConfig{
									Disabled:        &wrappers.BoolValue{Value: true},
									PollingInterval: &duration.Duration{Seconds: 12},
								},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.DatadogConfig{
					CollectorCluster: "datadog-cluster-name",
					ServiceName:      "datadog-gloo",
					// Note that RemoteConfig will be missing because it's disabled
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)

				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.datadog",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})
		})

		Describe("when opencensus provider config", func() {
			It("translates the plugin correctly using OcagentAddress", func() {

				expectedHttpAddress := "localhost:10000"
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_OpenCensusConfig{
							OpenCensusConfig: &envoytrace_gloo.OpenCensusConfig{
								TraceConfig: &envoytrace_gloo.TraceConfig{
									Sampler: &envoytrace_gloo.TraceConfig_ConstantSampler{
										ConstantSampler: &envoytrace_gloo.ConstantSampler{
											Decision: envoytrace_gloo.ConstantSampler_ALWAYS_ON,
										},
									},
									MaxNumberOfAttributes:    5,
									MaxNumberOfAnnotations:   10,
									MaxNumberOfMessageEvents: 15,
									MaxNumberOfLinks:         20,
								},
								OcagentExporterEnabled: true,
								OcagentAddress: &envoytrace_gloo.OpenCensusConfig_HttpAddress{
									HttpAddress: expectedHttpAddress,
								},
								IncomingTraceContext: nil,
								OutgoingTraceContext: nil,
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.OpenCensusConfig{
					TraceConfig: &v12.TraceConfig{
						Sampler: &v12.TraceConfig_ConstantSampler{
							ConstantSampler: &v12.ConstantSampler{
								Decision: v12.ConstantSampler_ALWAYS_ON,
							},
						},
						MaxNumberOfAttributes:    5,
						MaxNumberOfAnnotations:   10,
						MaxNumberOfMessageEvents: 15,
						MaxNumberOfLinks:         20,
					},
					OcagentExporterEnabled: true,
					OcagentAddress:         expectedHttpAddress,
					OcagentGrpcService:     nil,
					IncomingTraceContext:   nil,
					OutgoingTraceContext:   nil,
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)
				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.opencensus",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("translates the plugin correctly using OcagentGrpcService", func() {

				sampleGrpcTargetUri := "sampleGrpcTargetUri"
				sampleGrpcStatPrefix := "sampleGrpcStatPrefix"
				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_OpenCensusConfig{
							OpenCensusConfig: &envoytrace_gloo.OpenCensusConfig{
								TraceConfig: &envoytrace_gloo.TraceConfig{
									Sampler: &envoytrace_gloo.TraceConfig_ConstantSampler{
										ConstantSampler: &envoytrace_gloo.ConstantSampler{
											Decision: envoytrace_gloo.ConstantSampler_ALWAYS_ON,
										},
									},
									MaxNumberOfAttributes:    5,
									MaxNumberOfAnnotations:   10,
									MaxNumberOfMessageEvents: 15,
									MaxNumberOfLinks:         20,
								},
								OcagentExporterEnabled: true,
								OcagentAddress: &envoytrace_gloo.OpenCensusConfig_GrpcAddress{
									GrpcAddress: &envoytrace_gloo.OpenCensusConfig_OcagentGrpcAddress{
										TargetUri:  sampleGrpcTargetUri,
										StatPrefix: sampleGrpcStatPrefix,
									},
								},
								IncomingTraceContext: []envoytrace_gloo.OpenCensusConfig_TraceContext{
									envoytrace_gloo.OpenCensusConfig_NONE,
									envoytrace_gloo.OpenCensusConfig_TRACE_CONTEXT,
									envoytrace_gloo.OpenCensusConfig_GRPC_TRACE_BIN,
									envoytrace_gloo.OpenCensusConfig_CLOUD_TRACE_CONTEXT,
									envoytrace_gloo.OpenCensusConfig_B3,
								},
								OutgoingTraceContext: []envoytrace_gloo.OpenCensusConfig_TraceContext{
									envoytrace_gloo.OpenCensusConfig_B3,
									envoytrace_gloo.OpenCensusConfig_CLOUD_TRACE_CONTEXT,
									envoytrace_gloo.OpenCensusConfig_GRPC_TRACE_BIN,
									envoytrace_gloo.OpenCensusConfig_TRACE_CONTEXT,
									envoytrace_gloo.OpenCensusConfig_NONE,
								},
							},
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.OpenCensusConfig{
					TraceConfig: &v12.TraceConfig{
						Sampler: &v12.TraceConfig_ConstantSampler{
							ConstantSampler: &v12.ConstantSampler{
								Decision: v12.ConstantSampler_ALWAYS_ON,
							},
						},
						MaxNumberOfAttributes:    5,
						MaxNumberOfAnnotations:   10,
						MaxNumberOfMessageEvents: 15,
						MaxNumberOfLinks:         20,
					},
					OcagentExporterEnabled: true,
					OcagentAddress:         "",
					OcagentGrpcService: &envoy_config_core_v3.GrpcService{
						TargetSpecifier: &envoy_config_core_v3.GrpcService_GoogleGrpc_{
							GoogleGrpc: &envoy_config_core_v3.GrpcService_GoogleGrpc{
								TargetUri:  sampleGrpcTargetUri,
								StatPrefix: sampleGrpcStatPrefix,
							},
						},
					},
					IncomingTraceContext: []envoytrace.OpenCensusConfig_TraceContext{
						envoytrace.OpenCensusConfig_NONE,
						envoytrace.OpenCensusConfig_TRACE_CONTEXT,
						envoytrace.OpenCensusConfig_GRPC_TRACE_BIN,
						envoytrace.OpenCensusConfig_CLOUD_TRACE_CONTEXT,
						envoytrace.OpenCensusConfig_B3,
					},
					OutgoingTraceContext: []envoytrace.OpenCensusConfig_TraceContext{
						envoytrace.OpenCensusConfig_B3,
						envoytrace.OpenCensusConfig_CLOUD_TRACE_CONTEXT,
						envoytrace.OpenCensusConfig_GRPC_TRACE_BIN,
						envoytrace.OpenCensusConfig_TRACE_CONTEXT,
						envoytrace.OpenCensusConfig_NONE,
					},
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)
				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.opencensus",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}

				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})
		})

		Describe("when opentelemetry provider config", func() {
			It("translates the plugin correctly", func() {
				testClusterName := "test-cluster"
				otelConfig := &envoytrace_gloo.OpenTelemetryConfig{
					CollectorCluster: &envoytrace_gloo.OpenTelemetryConfig_ClusterName{
						ClusterName: testClusterName,
					},
				}

				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_OpenTelemetryConfig{
							OpenTelemetryConfig: otelConfig,
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.OpenTelemetryConfig{
					GrpcService: &envoy_config_core_v3.GrpcService{
						TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
								ClusterName: testClusterName,
							},
						},
					},
					ServiceName: "gateway",
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)
				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.opentelemetry",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}
				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})

			It("can override the service_name", func() {
				testClusterName := "test-cluster"
				otelConfig := &envoytrace_gloo.OpenTelemetryConfig{
					CollectorCluster: &envoytrace_gloo.OpenTelemetryConfig_ClusterName{
						ClusterName: testClusterName,
					},
					ServiceName: "custom-service-name",
				}

				cfg := &envoyhttp.HttpConnectionManager{}
				hcmSettings = &hcm.HttpConnectionManagerSettings{
					Tracing: &tracing.ListenerTracingSettings{
						ProviderConfig: &tracing.ListenerTracingSettings_OpenTelemetryConfig{
							OpenTelemetryConfig: otelConfig,
						},
					},
				}
				err := processHcmNetworkFilter(cfg)
				Expect(err).NotTo(HaveOccurred())

				expectedEnvoyConfig := &envoytrace.OpenTelemetryConfig{
					GrpcService: &envoy_config_core_v3.GrpcService{
						TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
								ClusterName: testClusterName,
							},
						},
					},
					ServiceName: "custom-service-name",
				}
				expectedEnvoyConfigMarshalled, _ := ptypes.MarshalAny(expectedEnvoyConfig)
				expectedEnvoyTracingProvider := &envoytrace.Tracing_Http{
					Name: "envoy.tracers.opentelemetry",
					ConfigType: &envoytrace.Tracing_Http_TypedConfig{
						TypedConfig: expectedEnvoyConfigMarshalled,
					},
				}
				Expect(cfg.Tracing.Provider.GetName()).To(Equal(expectedEnvoyTracingProvider.GetName()))
				Expect(cfg.Tracing.Provider.GetTypedConfig()).To(Equal(expectedEnvoyTracingProvider.GetTypedConfig()))
			})
		})

	})

	It("should update routes properly", func() {
		in := &v1.Route{}
		out := &envoy_config_route_v3.Route{}
		err := plugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{}, in, out)
		Expect(err).NotTo(HaveOccurred())

		inFull := &v1.Route{
			Options: &v1.RouteOptions{
				Tracing: &tracing.RouteTracingSettings{
					RouteDescriptor: "hello",
					Propagate:       &wrappers.BoolValue{Value: false},
				},
			},
		}
		outFull := &envoy_config_route_v3.Route{}
		err = plugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{}, inFull, outFull)
		Expect(err).NotTo(HaveOccurred())
		Expect(outFull.Decorator.Operation).To(Equal("hello"))
		Expect(outFull.Decorator.Propagate).To(Equal(&wrappers.BoolValue{Value: false}))
		Expect(outFull.Tracing.ClientSampling.Numerator / 10000).To(Equal(uint32(100)))
		Expect(outFull.Tracing.RandomSampling.Numerator / 10000).To(Equal(uint32(100)))
		Expect(outFull.Tracing.OverallSampling.Numerator / 10000).To(Equal(uint32(100)))
	})

	It("should update routes properly - with defaults", func() {
		in := &v1.Route{}
		out := &envoy_config_route_v3.Route{}
		err := plugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{}, in, out)
		Expect(err).NotTo(HaveOccurred())

		inFull := &v1.Route{
			Options: &v1.RouteOptions{
				Tracing: &tracing.RouteTracingSettings{
					RouteDescriptor: "hello",
					TracePercentages: &tracing.TracePercentages{
						ClientSamplePercentage:  &wrappers.FloatValue{Value: 10},
						RandomSamplePercentage:  &wrappers.FloatValue{Value: 20},
						OverallSamplePercentage: &wrappers.FloatValue{Value: 30},
					},
				},
			},
		}
		outFull := &envoy_config_route_v3.Route{}
		err = plugin.(plugins.RoutePlugin).ProcessRoute(plugins.RouteParams{}, inFull, outFull)
		Expect(err).NotTo(HaveOccurred())
		Expect(outFull.Decorator.Operation).To(Equal("hello"))
		Expect(outFull.Decorator.Propagate).To(BeNil())
		Expect(outFull.Tracing.ClientSampling.Numerator / 10000).To(Equal(uint32(10)))
		Expect(outFull.Tracing.RandomSampling.Numerator / 10000).To(Equal(uint32(20)))
		Expect(outFull.Tracing.OverallSampling.Numerator / 10000).To(Equal(uint32(30)))
	})

})
