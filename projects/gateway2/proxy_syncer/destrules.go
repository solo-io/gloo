package proxy_syncer

import (
	"context"
	"fmt"
	"slices"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	networkingv1 "istio.io/api/networking/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/kube/krt"
)

const (
	ExtensionName = "destrule"
)

const (
	FailoverPriorityLabelDefaultSeparator = '='
)

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

type DestinationRuleWrapper struct {
	*networkingclient.DestinationRule
}

func (s DestinationRuleWrapper) GetLabelSelector() map[string]string {
	return s.Spec.WorkloadSelector.MatchLabels
}

type nsWithHostname struct {
	ns       string
	hostname string
}

type destRuleIndex struct {
	idx krt.Index[nsWithHostname, DestinationRuleWrapper]
}

func newDestruleIndex(destRuleCollection krt.Collection[DestinationRuleWrapper]) destRuleIndex {
	idx := krt.NewIndex(destRuleCollection, func(d DestinationRuleWrapper) []nsWithHostname {
		return []nsWithHostname{{
			ns:       d.Namespace,
			hostname: d.Spec.Host,
		}}
	})
	return destRuleIndex{idx: idx}
}

func (d *destRuleIndex) getDesrRules(ns string,
	hostname string) []DestinationRuleWrapper {
	return d.idx.Lookup(struct {
		ns       string
		hostname string
	}{
		ns:       ns,
		hostname: hostname,
	})
}

func applyDestRulesForHostnames(kctx krt.HandlerContext, wrappedDestRules krt.Collection[DestinationRuleWrapper],
	destinationRulesIndex destRuleIndex, workloadNs string, ep EndpointsForUpstream, podLabels map[string]string) *CLA {
	// host that would match the dest rule from the endpoints.
	// get the matching dest rule
	// get the lb info from the dest rules and call prioritize

	hostname := fromEndpoint(ep)
	key := nsWithHostname{ns: workloadNs, hostname: hostname}
	destrules := krt.Fetch(kctx, wrappedDestRules, krt.FilterIndex(destinationRulesIndex.idx, key), krt.FilterSelects(podLabels))

	lbInfo := getDestruleFor(destrules)
	var priorities *priorities
	if lbInfo != nil {
		priorities = newPriorities(lbInfo.failoverPriority, lbInfo.proxyLabels)
	}

	return prioritize(ep, lbInfo, priorities)
}

func fromEndpoint(ep EndpointsForUpstream) string {
	// get the upstream name and namespace
	// TODO: suppport other suffixes that are not cluster.local
	return fmt.Sprintf("%s.%s.svc.cluster.local", ep.UpstreamRef.Name, ep.UpstreamRef.Namespace)
}

func getDestruleFor(destrules []DestinationRuleWrapper) *LBInfo {

	// use oldest. TODO -  we need to merge them.
	oldestDestRule := slices.MinFunc(destrules, func(i DestinationRuleWrapper, j DestinationRuleWrapper) int {
		return i.CreationTimestamp.Time.Compare(j.CreationTimestamp.Time)
	})
	localityLb := oldestDestRule.Spec.GetTrafficPolicy().GetLoadBalancer().GetLocalityLbSetting()
	if localityLb == nil {
		return nil
	}

	if localityLb.GetFailoverPriority() != nil {
		failoverPriority := localityLb.GetFailoverPriority()
		return &LBInfo{
			failoverPriority: failoverPriority,
		}
	}

	panic("TODO: implement other locality lb")
}
