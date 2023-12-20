package tap

import (
	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/config/common/matcher/v3"
	envoycoreconfig "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoytapconfig "github.com/envoyproxy/go-control-plane/envoy/config/tap/v3"
	envoytapcommon "github.com/envoyproxy/go-control-plane/envoy/extensions/common/tap/v3"
	envoytap "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/tap/v3"
	envoyhcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	pany "github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	. "github.com/onsi/ginkgo/v2"
	solotapsinks "github.com/solo-io/gloo/projects/gloo/pkg/api/config/tap/output_sink/v3"
	envoycore "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"

	. "github.com/onsi/gomega"
	solotap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Plugin", func() {
	Context("tap filter output sinks", func() {
		var (
			httpUpstream = &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "test-http-tap-server",
					Namespace: "test-ns",
				},
			}
			grpcUpstream = &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "test-grpc-tap-server",
					Namespace: "test-ns",
				},
			}
		)

		sinkTableFunc := func(
			upstream *gloov1.Upstream,
			glooSinks []*solotap.Sink,
			envoySinkFunc func() (*pany.Any, error)) {
			params := plugins.Params{
				Snapshot: &gloov1snap.ApiSnapshot{
					Upstreams: []*gloov1.Upstream{upstream},
				},
			}
			filters, err := NewPlugin().HttpFilters(params, &gloov1.HttpListener{
				Options: &gloov1.HttpListenerOptions{
					Tap: &solotap.Tap{
						Sinks: glooSinks,
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			typedExtensionConfig, err := envoySinkFunc()
			Expect(err).NotTo(HaveOccurred())

			envoyConfig := &envoytap.Tap{
				CommonConfig: &envoytapcommon.CommonExtensionConfig{
					ConfigType: &envoytapcommon.CommonExtensionConfig_StaticConfig{
						StaticConfig: &envoytapconfig.TapConfig{
							Match: &envoymatcher.MatchPredicate{
								Rule: &envoymatcher.MatchPredicate_AnyMatch{
									AnyMatch: true,
								},
							},
							OutputConfig: &envoytapconfig.OutputConfig{
								Sinks: []*envoytapconfig.OutputSink{{
									Format: envoytapconfig.OutputSink_PROTO_BINARY,
									OutputSinkType: &envoytapconfig.OutputSink_CustomSink{
										CustomSink: &envoycoreconfig.TypedExtensionConfig{
											Name:        EnvoyExtensionName,
											TypedConfig: typedExtensionConfig,
										},
									},
								}},
							},
						},
					},
				},
			}
			typedConfig, err := utils.MessageToAny(envoyConfig)
			Expect(err).NotTo(HaveOccurred())
			envoyFilter := plugins.StagedHttpFilter{
				HttpFilter: &envoyhcm.HttpFilter{
					Name: FilterName,
					ConfigType: &envoyhcm.HttpFilter_TypedConfig{
						TypedConfig: typedConfig,
					},
				},
				Stage: filterStage,
			}
			Expect(filters).To(Equal([]plugins.StagedHttpFilter{envoyFilter}))
		}
		DescribeTable("output sink table tests", sinkTableFunc,
			Entry(
				"configures an http tap sink",
				httpUpstream,
				[]*solotap.Sink{{
					SinkType: &solotap.Sink_HttpService{
						HttpService: &solotap.HttpService{
							TapServer: &core.ResourceRef{
								Name:      httpUpstream.GetMetadata().GetName(),
								Namespace: httpUpstream.GetMetadata().GetNamespace(),
							},
							Timeout: &duration.Duration{Seconds: 10},
						},
					},
				}},
				func() (*pany.Any, error) {
					return utils.MessageToAny(
						&solotapsinks.HttpOutputSink{
							ServerUri: &envoycore.HttpUri{
								Uri: DefaultTraceUri,
								HttpUpstreamType: &envoycore.HttpUri_Cluster{
									Cluster: translator.UpstreamToClusterName(&core.ResourceRef{
										Name:      httpUpstream.GetMetadata().GetName(),
										Namespace: httpUpstream.GetMetadata().GetNamespace(),
									}),
								},
								Timeout: &duration.Duration{Seconds: 10},
							},
						})
				}),
			Entry(
				"configures a grpc tap sink",
				grpcUpstream,
				[]*solotap.Sink{{
					SinkType: &solotap.Sink_GrpcService{
						GrpcService: &solotap.GrpcService{
							TapServer: &core.ResourceRef{
								Name:      grpcUpstream.GetMetadata().GetName(),
								Namespace: grpcUpstream.GetMetadata().GetNamespace(),
							},
						},
					},
				}},
				func() (*pany.Any, error) {
					return utils.MessageToAny(&solotapsinks.GrpcOutputSink{
						GrpcService: &envoycore.GrpcService{
							TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
								EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
									ClusterName: translator.UpstreamToClusterName(&core.ResourceRef{
										Name:      grpcUpstream.GetMetadata().GetName(),
										Namespace: grpcUpstream.GetMetadata().GetNamespace(),
									}),
								},
							},
						},
					})
				}),
		)

		It("returns nil if no sinks are provided", func() {
			params := plugins.Params{
				Snapshot: &gloov1snap.ApiSnapshot{
					Upstreams: []*gloov1.Upstream{grpcUpstream},
				},
			}
			filters, err := NewPlugin().HttpFilters(params, &gloov1.HttpListener{
				Options: &gloov1.HttpListenerOptions{
					Tap: &solotap.Tap{
						Sinks: []*solotap.Sink{},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(BeNil())
		})

		It("returns an error if multiple sinks are provided", func() {
			params := plugins.Params{
				Snapshot: &gloov1snap.ApiSnapshot{
					Upstreams: []*gloov1.Upstream{grpcUpstream},
				},
			}
			sink := &solotap.Sink_GrpcService{
				GrpcService: &solotap.GrpcService{
					TapServer: &core.ResourceRef{
						Name:      params.Snapshot.Upstreams[0].GetMetadata().GetName(),
						Namespace: params.Snapshot.Upstreams[0].GetMetadata().GetNamespace(),
					},
				},
			}
			filters, err := NewPlugin().HttpFilters(params, &gloov1.HttpListener{
				Options: &gloov1.HttpListenerOptions{
					Tap: &solotap.Tap{
						Sinks: []*solotap.Sink{
							{SinkType: sink},
							{SinkType: sink},
						},
					},
				},
			})
			Expect(err).To(MatchError(Equal("exactly one sink must be specified for tap filter")))
			Expect(filters).To(BeNil())
		})
	})
})
