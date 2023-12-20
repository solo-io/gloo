package tap

import (
	"errors"
	"fmt"

	envoycoreconfig "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/rotisserie/eris"
	envoycore "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/types/known/anypb"

	envoymatcher "github.com/envoyproxy/go-control-plane/envoy/config/common/matcher/v3"
	envoytapconfig "github.com/envoyproxy/go-control-plane/envoy/config/tap/v3"
	envoytapcommon "github.com/envoyproxy/go-control-plane/envoy/extensions/common/tap/v3"
	envoytap "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/tap/v3"
	solotapsinks "github.com/solo-io/gloo/projects/gloo/pkg/api/config/tap/output_sink/v3"
	solotap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tap"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
)

const (
	FilterName         = "envoy.filters.http.tap"
	ExtensionName      = "tap"
	DefaultTraceUri    = "submit_trace"
	EnvoyExtensionName = "envoy.tap.sinks.solo.grpc_output_sink"
)

// TODO figure out what stage the plugin needs
// https://github.com/solo-io/solo-projects/issues/5513
// for now, put it after ratelimit which seems like the best choice for
// security (ie. prevents users from flooding the tap server)
var (
	filterStage = plugins.AfterStage(plugins.RateLimitStage)
)

type plugin struct {
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(initParams plugins.InitParams) {
}

func buildHttpOutputSink(upstreamRef *core.ResourceRef, timeout *duration.Duration) (*anypb.Any, error) {
	httpOutputSink := &solotapsinks.HttpOutputSink{
		ServerUri: &envoycore.HttpUri{
			Uri: DefaultTraceUri,
			HttpUpstreamType: &envoycore.HttpUri_Cluster{
				Cluster: translator.UpstreamToClusterName(upstreamRef),
			},
			Timeout: timeout,
		},
	}
	return utils.MessageToAny(httpOutputSink)
}

func buildGrpcOutputSink(upstreamRef *core.ResourceRef) (*anypb.Any, error) {
	grpcOutputSink := &solotapsinks.GrpcOutputSink{
		GrpcService: &envoycore.GrpcService{
			TargetSpecifier: &envoycore.GrpcService_EnvoyGrpc_{
				EnvoyGrpc: &envoycore.GrpcService_EnvoyGrpc{
					ClusterName: translator.UpstreamToClusterName(upstreamRef),
				},
			},
		},
	}
	return utils.MessageToAny(grpcOutputSink)
}

func buildEnvoyTapConfig(typedExtensionConfig *anypb.Any) *envoytap.Tap {
	return &envoytap.Tap{
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
							// Format is presently not used by the output sink
							// extensions - in order to support alternate
							// formats, we would have to explicitly implement
							// them ourselves in the http/grpc output sinks,
							// but since all data is written over the network,
							// we only support writing in the protobuf wireline
							// format for now.
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
}

// Placeholder function for retrieving tap sinks. In the future, when it is
// possible to define the tap filter on multiple levels (settings, gateway,
// virtual service, etc.) we can use this function to retrieve the tap sink
// from the correct source and reconcile differences between the higher and
// lower level settings. But since the tap filter can only be configured on
// gateways for now, this function is just a very straightforward retrieval
func (p *plugin) getTapSinks(listener *v1.HttpListener) []*solotap.Sink {
	return listener.GetOptions().GetTap().GetSinks()
}

func (p *plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	sinks := p.getTapSinks(listener)
	if len(sinks) == 0 {
		return nil, nil // tap sink is not configured
	}
	// Upstream envoy only supports a single tap filter for now, so we do as well:
	// https://github.com/envoyproxy/envoy/blob/234e408bc2f417b66d045450ce61b894528ec313/api/envoy/config/tap/v3/common.proto#L161-L163
	if len(sinks) > 1 {
		return nil, errors.New("exactly one sink must be specified for tap filter")
	}

	var typedExtensionConfig *anypb.Any
	switch sink := sinks[0].GetSinkType().(type) {
	case *solotap.Sink_GrpcService:
		upstreamRef := sink.GrpcService.GetTapServer()
		_, err := params.Snapshot.Upstreams.Find(upstreamRef.GetNamespace(), upstreamRef.GetName())
		if err != nil {
			return nil, eris.Wrapf(err, "tap filter server upstream not found: %s", upstreamRef.String())
		}
		// TODO can we return an error if the upstream is not configured to use http/2?
		typedExtensionConfig, err = buildGrpcOutputSink(upstreamRef)
		if err != nil {
			return nil, eris.Wrap(err, "failed to generate gRPC output sink configuration")
		}
	case *solotap.Sink_HttpService:
		upstreamRef := sink.HttpService.GetTapServer()
		_, err := params.Snapshot.Upstreams.Find(upstreamRef.GetNamespace(), upstreamRef.GetName())
		if err != nil {
			return nil, eris.Wrapf(err, "tap filter server upstream not found: %s", upstreamRef.String())
		}
		typedExtensionConfig, err = buildHttpOutputSink(upstreamRef, sink.HttpService.GetTimeout())
		if err != nil {
			return nil, eris.Wrap(err, "failed to generate HTTP output sink configuration")
		}
	default:
		return nil, errors.New(fmt.Sprintf("not implemented: tap sink %T\n", sink))
	}

	envoyTapConfig := buildEnvoyTapConfig(typedExtensionConfig)
	stagedFilter, err := plugins.NewStagedFilter(FilterName, envoyTapConfig, filterStage)
	if err != nil {
		return nil, eris.Wrap(err, "failed to generate staged filter configuration")
	}
	return []plugins.StagedHttpFilter{stagedFilter}, nil
}
