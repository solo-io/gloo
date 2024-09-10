package destrule

import (
	"context"
	"slices"
	"sync"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/wrapperspb"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	networkingv1 "istio.io/api/networking/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pilot/pkg/features"
	"istio.io/istio/pkg/config/labels"
	"istio.io/istio/pkg/config/schema/gvr"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
)

var (
	_ plugins.Plugin         = new(plugin)
	_ plugins.UpstreamPlugin = new(plugin)
)

const (
	ExtensionName = "destrule"
)

type plugin struct {
	destRuleIndex *DestruleIndex
}

func NewPlugin() *plugin {
	p := &plugin{}

	return p
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
	onceIndex.Do(func() {
		// TODO: this is very hacky, don't do this
		kubeRestConfig, err := kube.DefaultRestConfig("", "")
		if err != nil {
			return
		}
		cfg := kube.NewClientConfigForRestConfig(kubeRestConfig)
		client, err := kube.NewClient(cfg, "gloo-gateway")
		if err != nil {
			return
		}
		singleDestRuleIndex = newDestruleIndex(client)

	})
	p.destRuleIndex = singleDestRuleIndex
	if p.destRuleIndex == nil {
		panic("dest rule index not initialized")
	}
}

var (
	onceIndex           sync.Once
	singleDestRuleIndex *DestruleIndex
)

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *envoy_config_cluster_v3.Cluster) error {
	// 1. find out if the upstream is a service entry, return if not
	// 2. if it is a service entry, find all the destination rules that apply to gloo
	//   - these are:
	//     - destination rules that are global, on our namespace, or if they have selectors, they select us.
	// 3. if the destination rule applies to the upstream host name, then apply settings from the destination rule to us.

	// krt can help with finding the correct destination rules

	hostname := getHostname(in)
	if hostname == "" {
		return nil
	}

	destrules := p.destRuleIndex.OurDestRulesByHostName.Lookup(hostname)
	if len(destrules) == 0 {
		return nil
	}

	// use oldest. TODO - do we need to merge them?
	oldestDestRule := slices.MinFunc(destrules, func(i networkingclient.DestinationRule, j networkingclient.DestinationRule) int {
		return i.CreationTimestamp.Time.Compare(j.CreationTimestamp.Time)
	})

	// apply settings from the oldest one

	// apply mtls settings

	applyPolicy(params.Ctx, out, oldestDestRule.Spec.TrafficPolicy)

	return nil
}

func applyPolicy(ctx context.Context, out *envoy_config_cluster_v3.Cluster, policy *networkingv1.TrafficPolicy) {
	// apply mtls settings
	if policy == nil {
		return
	}
	// TODO: handle port level policy later. should we handle subsets?
	applyOutlierDetection(out, policy.OutlierDetection)
	applyUpstreamTLSSettings(ctx, out, policy)

}

// TODO: this is copied from istio, we may need to adjust this
// FIXME: there isn't a way to distinguish between unset values and zero values
func applyOutlierDetection(c *envoy_config_cluster_v3.Cluster, outlier *networkingv1.OutlierDetection) {
	if outlier == nil {
		return
	}

	out := &envoy_config_cluster_v3.OutlierDetection{}

	// SuccessRate based outlier detection should be disabled.
	out.EnforcingSuccessRate = &wrapperspb.UInt32Value{Value: 0}

	if e := outlier.Consecutive_5XxErrors; e != nil {
		v := e.GetValue()

		out.Consecutive_5Xx = &wrapperspb.UInt32Value{Value: v}

		if v > 0 {
			v = 100
		}
		out.EnforcingConsecutive_5Xx = &wrapperspb.UInt32Value{Value: v}
	}
	if e := outlier.ConsecutiveGatewayErrors; e != nil {
		v := e.GetValue()

		out.ConsecutiveGatewayFailure = &wrapperspb.UInt32Value{Value: v}

		if v > 0 {
			v = 100
		}
		out.EnforcingConsecutiveGatewayFailure = &wrapperspb.UInt32Value{Value: v}
	}

	if outlier.Interval != nil {
		out.Interval = outlier.Interval
	}
	if outlier.BaseEjectionTime != nil {
		out.BaseEjectionTime = outlier.BaseEjectionTime
	}
	if outlier.MaxEjectionPercent > 0 {
		out.MaxEjectionPercent = &wrapperspb.UInt32Value{Value: uint32(outlier.MaxEjectionPercent)}
	}

	if outlier.SplitExternalLocalOriginErrors {
		out.SplitExternalLocalOriginErrors = true
		if outlier.ConsecutiveLocalOriginFailures.GetValue() > 0 {
			out.ConsecutiveLocalOriginFailure = &wrapperspb.UInt32Value{Value: outlier.ConsecutiveLocalOriginFailures.Value}
			out.EnforcingConsecutiveLocalOriginFailure = &wrapperspb.UInt32Value{Value: 100}
		}
		// SuccessRate based outlier detection should be disabled.
		out.EnforcingLocalOriginSuccessRate = &wrapperspb.UInt32Value{Value: 0}
	}

	c.OutlierDetection = out

	// Disable panic threshold by default as its not typically applicable in k8s environments
	// with few pods per service.
	// To do so, set the healthy_panic_threshold field even if its value is 0 (defaults to 50 in Envoy).
	// FIXME: we can't distinguish between it being unset or being explicitly set to 0
	minHealthPercent := outlier.MinHealthPercent
	if minHealthPercent >= 0 {
		// When we are sending unhealthy endpoints, we should disable Panic Threshold. Otherwise
		// Envoy will send traffic to "Unready" pods when the percentage of healthy hosts fall
		// below minimum health percentage.
		if features.SendUnhealthyEndpoints.Load() { // TODO: what is the gloo equivalent?
			minHealthPercent = 0
		}
		c.CommonLbConfig.HealthyPanicThreshold = &envoy_type_v3.Percent{Value: float64(minHealthPercent)}
	}
}

func applyUpstreamTLSSettings(ctx context.Context,
	out *envoy_config_cluster_v3.Cluster, policy *networkingv1.TrafficPolicy,
) {
	tls := policy.Tls
	if tls == nil {
		return
	}
	if tls.Mode != networkingv1.ClientTLSSettings_ISTIO_MUTUAL {
		// TODO: deal with other modes
		return
	}

	settings := settingsutil.MaybeFromContext(ctx)
	if settings == nil {
		return
	}

	// make sure if we have a sidecar
	// TODO: this should be enableintegration
	haveSidecar := settings.GetGloo().GetIstioOptions().GetEnableAutoMtls().GetValue()

	if !haveSidecar {
		return
	}

	// TODO: this is copied from istio plugin, we may need to un-copypaste this.
	sslSds := &tlsv3.UpstreamTlsContext{
		CommonTlsContext: &tlsv3.CommonTlsContext{
			AlpnProtocols: []string{"istio"},
			TlsParams:     &tlsv3.TlsParameters{},
			ValidationContextType: &tlsv3.CommonTlsContext_ValidationContextSdsSecretConfig{
				ValidationContextSdsSecretConfig: &tlsv3.SdsSecretConfig{
					Name: constants.IstioValidationContext,
					SdsConfig: &envoy_config_core_v3.ConfigSource{
						ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
						ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_ApiConfigSource{
							ApiConfigSource: &envoy_config_core_v3.ApiConfigSource{
								// Istio sets this to skip the node identifier in later discovery requests
								SetNodeOnFirstMessageOnly: true,
								ApiType:                   envoy_config_core_v3.ApiConfigSource_GRPC,
								TransportApiVersion:       envoy_config_core_v3.ApiVersion_V3,
								GrpcServices: []*envoy_config_core_v3.GrpcService{
									{
										TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
											EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{ClusterName: constants.SdsClusterName},
										},
									},
								},
							},
						},
					},
				},
			},
			TlsCertificateSdsSecretConfigs: []*tlsv3.SdsSecretConfig{
				{
					Name: constants.IstioCertSecret,
					SdsConfig: &envoy_config_core_v3.ConfigSource{
						ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
						ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_ApiConfigSource{
							ApiConfigSource: &envoy_config_core_v3.ApiConfigSource{
								ApiType: envoy_config_core_v3.ApiConfigSource_GRPC,
								// Istio sets this to skip the node identifier in later discovery requests
								SetNodeOnFirstMessageOnly: true,
								TransportApiVersion:       envoy_config_core_v3.ApiVersion_V3,
								GrpcServices: []*envoy_config_core_v3.GrpcService{
									{
										TargetSpecifier: &envoy_config_core_v3.GrpcService_EnvoyGrpc_{
											EnvoyGrpc: &envoy_config_core_v3.GrpcService_EnvoyGrpc{
												ClusterName: constants.SdsClusterName,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	typedConfig, _ := glooutils.MessageToAny(sslSds)
	transportSocket := &envoy_config_core_v3.TransportSocket{
		Name:       wellknown.TransportSocketTls,
		ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{TypedConfig: typedConfig},
	}
	out.TransportSocket = transportSocket
}

func getHostname(upstream *v1.Upstream) string {
	if len(upstream.GetStatic().GetHosts()) != 0 {
		return upstream.GetStatic().GetHosts()[0].Addr
	}
	return ""
}

type DestruleIndex struct {
	OurDestRulesByHostName *krt.Index[networkingclient.DestinationRule, string]
}

func newDestruleIndex(c kube.Client) *DestruleIndex {
	var d DestruleIndex

	ourNs := utils.GetPodNamespace()
	ourLabels := utils.GetPodLabels()

	destinationRules := kclient.NewDelayedInformer[*networkingclient.DestinationRule](c,
		gvr.DestinationRule, kubetypes.StandardInformer, kclient.Filter{Namespace: ourNs})
	DestinationRules := krt.WrapClient[*networkingclient.DestinationRule](destinationRules, krt.WithName("DestinationRules"))

	// filter the ones that apply to us
	// look for ones in the config namespace (for now we ignore this), our namespace with no selectors, or with selectors that select us
	ourDestRules := krt.NewCollection(DestinationRules, func(ctx krt.HandlerContext, i *networkingclient.DestinationRule) *networkingclient.DestinationRule {
		// make sure this either doesn't have selectors, or they select us:
		selector := i.Spec.WorkloadSelector
		if selector == nil {
			return i
		}
		// see if selectors select us
		if labels.Instance(i.Spec.WorkloadSelector.MatchLabels).SubsetOf(ourLabels) {
			return i
		}
		return nil
	})
	// index by hostname
	d.OurDestRulesByHostName = krt.NewIndex(ourDestRules, func(s networkingclient.DestinationRule) []string {
		return []string{s.Spec.Host}
	})

	ourDestRules.Register(func(o krt.Event[networkingclient.DestinationRule]) {
		d.kickTranslation()
	})
	return &d
}

func (d *DestruleIndex) kickTranslation() {
	// TODO - implement
}
