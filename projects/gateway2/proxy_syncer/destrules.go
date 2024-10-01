package proxy_syncer

import (
	"context"
	"fmt"
	"slices"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/types"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/solo-io/gloo/pkg/utils"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
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

const (
	FailoverPriorityLabelDefaultSeparator = '='
)

type plugin struct {
	destRuleIndex *DestruleIndex
}

func NewPlugin() *plugin {
	p := &plugin{}
	// TODO: this is very hacky, don't do this
	kubeRestConfig, err := kube.DefaultRestConfig("", "")
	if err != nil {
		panic(err)
	}
	cfg := kube.NewClientConfigForRestConfig(kubeRestConfig)
	client, err := kube.NewClient(cfg, "gloo-gateway")
	if err != nil {
		panic(err)
	}

	p.destRuleIndex = newDestruleIndex(client, "", nil)

	return p
}

var (
	_ plugins.Plugin         = new(plugin)
	_ plugins.UpstreamPlugin = new(plugin)
)

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
}

/*
func (p *plugin) ProcessEndpoints(params plugins.EndpointParams, in *v1.Upstream, endpoints []*v1.Endpoint, out *envoy_config_endpoint_v3.ClusterLoadAssignment) error {
	// re-order the endpoints by locality?!
	destRule := p.getDestRuleForUpstream(in)
	if destRule == nil {
		return nil
	}

	// go over the labels and figure out priorities

	return nil
}
*/

func (p *plugin) ProcessUpstream(params plugins.Params, in *v1.Upstream, out *envoy_config_cluster_v3.Cluster) error {
	destRule := p.getDestRuleForUpstream(in)
	if destRule == nil {
		return nil
	}

	// apply settings

	// apply mtls settings

	applyPolicy(params.Ctx, out, destRule.Spec.TrafficPolicy)

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

func (p *plugin) getDestRuleForUpstream(in *v1.Upstream) *networkingclient.DestinationRule {
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
	return &oldestDestRule
}

/*
func applyLocalityLoadBalancer(
	loadAssignment *envoy_config_endpoint_v3.ClusterLoadAssignment,
	wrappedLocalityLbEndpoints []*v1.Endpoint,
	locality *envoy_config_core_v3.Locality,
	proxyLabels map[string]string,
	localityLB *networkingv1.LocalityLoadBalancerSetting,
	enableFailover bool,
) {
	// before calling this function localityLB.enabled field has been checked.
	if localityLB == nil || loadAssignment == nil {
		return
	}

	// one of Distribute or Failover settings can be applied.
	if localityLB.GetDistribute() != nil {
		applyLocalityWeights(locality, loadAssignment, localityLB.GetDistribute())
		// Failover needs outlier detection, otherwise Envoy will never drop down to a lower priority.
		// Do not apply default failover when locality LB is disabled.
	} else if enableFailover {
		if len(localityLB.FailoverPriority) > 0 {
			// Apply user defined priority failover settings.
			applyFailoverPriorities(loadAssignment, wrappedLocalityLbEndpoints, proxyLabels, localityLB.FailoverPriority)
			// If failover is expliciltly configured with failover priority, apply failover settings also.
			if len(localityLB.Failover) != 0 {
				applyLocalityFailover(locality, loadAssignment, localityLB.Failover)
			}
		} else {
			// Apply default failover settings or user defined region failover settings.
			applyLocalityFailover(locality, loadAssignment, localityLB.Failover)
		}
	}
}

// set locality loadbalancing weight based on user defined weights.
func applyLocalityWeights(
	locality *envoy_config_core_v3.Locality,
	loadAssignment *envoy_config_endpoint_v3.ClusterLoadAssignment,
	distribute []*networkingv1.LocalityLoadBalancerSetting_Distribute,
) {
	if distribute == nil {
		return
	}

	// Support Locality weighted load balancing
	// (https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/upstream/load_balancing/locality_weight#locality-weighted-load-balancing)
	// by providing weights in LocalityLbEndpoints via load_balancing_weight.
	// By setting weights across different localities, it can allow
	// Envoy to do weighted load balancing across different zones and geographical locations.
	for _, localityWeightSetting := range distribute {
		if localityWeightSetting != nil &&
			util.LocalityMatch(locality, localityWeightSetting.From) {
			misMatched := sets.Set[int]{}
			for i := range loadAssignment.Endpoints {
				misMatched.Insert(i)
			}
			for locality, weight := range localityWeightSetting.To {
				// index -> original weight
				destLocMap := map[int]uint32{}
				totalWeight := uint32(0)
				for i, ep := range loadAssignment.Endpoints {
					if misMatched.Contains(i) {
						if util.LocalityMatch(ep.Locality, locality) {
							delete(misMatched, i)
							if ep.LoadBalancingWeight != nil {
								destLocMap[i] = ep.LoadBalancingWeight.Value
							} else {
								destLocMap[i] = 1
							}
							totalWeight += destLocMap[i]
						}
					}
				}
				// in case wildcard dest matching multi groups of endpoints
				// the load balancing weight for a locality is divided by the sum of the weights of all localities
				for index, originalWeight := range destLocMap {
					destWeight := float64(originalWeight*weight) / float64(totalWeight)
					if destWeight > 0 {
						loadAssignment.Endpoints[index].LoadBalancingWeight = &wrappers.UInt32Value{
							Value: uint32(math.Ceil(destWeight)),
						}
					}
				}
			}

			// remove groups of endpoints in a locality that miss matched
			for i := range misMatched {
				if loadAssignment.Endpoints[i] != nil {
					loadAssignment.Endpoints[i].LbEndpoints = nil
				}
			}
			break
		}
	}
}

// set locality loadbalancing priority - This is based on Region/Zone/SubZone matching.
func applyLocalityFailover(
	locality *envoy_config_core_v3.Locality,
	loadAssignment *envoy_config_endpoint_v3.ClusterLoadAssignment,
	failover []*networkingv1.LocalityLoadBalancerSetting_Failover,
) {
	// key is priority, value is the index of the LocalityLbEndpoints in ClusterLoadAssignment
	priorityMap := map[int][]int{}

	// 1. calculate the LocalityLbEndpoints.Priority compared with proxy locality
	for i, localityEndpoint := range loadAssignment.Endpoints {
		// if region/zone/subZone all match, the priority is 0.
		// if region/zone match, the priority is 1.
		// if region matches, the priority is 2.
		// if locality not match, the priority is 3.
		priority := util.LbPriority(locality, localityEndpoint.Locality)
		// region not match, apply failover settings when specified
		// update localityLbEndpoints' priority to 4 if failover not match
		if priority == 3 {
			for _, failoverSetting := range failover {
				if failoverSetting.From == locality.Region {
					if localityEndpoint.Locality == nil || localityEndpoint.Locality.Region != failoverSetting.To {
						priority = 4
					}
					break
				}
			}
		}
		// priority is calculated using the already assigned priority using failoverPriority.
		// Since there are at most 5 priorities can be assigned using locality failover(0-4),
		// we multiply the priority by 5 for maintaining the priorities already assigned.
		// Afterwards the final priorities can be calculted from 0 (highest) to N (lowest) without skipping.
		priorityInt := int(loadAssignment.Endpoints[i].Priority*5) + priority
		loadAssignment.Endpoints[i].Priority = uint32(priorityInt)
		priorityMap[priorityInt] = append(priorityMap[priorityInt], i)
	}

	// since Priorities should range from 0 (highest) to N (lowest) without skipping.
	// 2. adjust the priorities in order
	// 2.1 sort all priorities in increasing order.
	priorities := []int{}
	for priority := range priorityMap {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)
	// 2.2 adjust LocalityLbEndpoints priority
	// if the index and value of priorities array is not equal.
	for i, priority := range priorities {
		if i != priority {
			// the LocalityLbEndpoints index in ClusterLoadAssignment.Endpoints
			for _, index := range priorityMap[priority] {
				loadAssignment.Endpoints[index].Priority = uint32(i)
			}
		}
	}
}

// set loadbalancing priority by failover priority label.
func applyFailoverPriorities(
	loadAssignment *envoy_config_endpoint_v3.ClusterLoadAssignment,
	wrappedLocalityLbEndpoints []*v1.Endpoint,
	proxyLabels map[string]string,
	failoverPriorities []string,
) {
	if len(proxyLabels) == 0 || len(wrappedLocalityLbEndpoints) == 0 {
		return
	}
	priorityMap := make(map[int][]int, len(failoverPriorities))
	localityLbEndpoints := []*envoy_config_endpoint_v3.LocalityLbEndpoints{}
	for _, wrappedLbEndpoint := range wrappedLocalityLbEndpoints {
		localityLbEndpointsPerLocality := applyFailoverPriorityPerLocality(proxyLabels, wrappedLbEndpoint, failoverPriorities)
		localityLbEndpoints = append(localityLbEndpoints, localityLbEndpointsPerLocality...)
	}
	for i, ep := range localityLbEndpoints {
		priorityMap[int(ep.Priority)] = append(priorityMap[int(ep.Priority)], i)
	}
	// since Priorities should range from 0 (highest) to N (lowest) without skipping.
	// adjust the priorities in order
	// 1. sort all priorities in increasing order.
	priorities := []int{}
	for priority := range priorityMap {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)
	// 2. adjust LocalityLbEndpoints priority
	// if the index and value of priorities array is not equal.
	for i, priority := range priorities {
		if i != priority {
			// the LocalityLbEndpoints index in ClusterLoadAssignment.Endpoints
			for _, index := range priorityMap[priority] {
				localityLbEndpoints[index].Priority = uint32(i)
			}
		}
	}
	loadAssignment.Endpoints = localityLbEndpoints
}

// Returning the label names in a separate array as the iteration of map is not ordered.
func priorityLabelOverrides(labels []string) ([]string, map[string]string) {
	priorityLabels := make([]string, 0, len(labels))
	overriddenValueByLabel := make(map[string]string, len(labels))
	var tempStrings []string
	for _, labelWithValue := range labels {
		tempStrings = strings.Split(labelWithValue, string(FailoverPriorityLabelDefaultSeparator))
		priorityLabels = append(priorityLabels, tempStrings[0])
		if len(tempStrings) == 2 {
			overriddenValueByLabel[tempStrings[0]] = tempStrings[1]
			continue
		}
	}
	return priorityLabels, overriddenValueByLabel
}

// set loadbalancing priority by failover priority label.
// split one LocalityLbEndpoints to multiple LocalityLbEndpoints based on failover priorities.
func applyFailoverPriorityPerLocality(
	proxyLabels map[string]string,
	ep []*v1.Endpoint,
	failoverPriorities []string,
) []*envoy_config_endpoint_v3.LocalityLbEndpoints {
	lowestPriority := len(failoverPriorities)
	// key is priority, value is the index of LocalityLbEndpoints.LbEndpoints
	priorityMap := map[int][]int{}
	priorityLabels, priorityLabelOverrides := priorityLabelOverrides(failoverPriorities)
	for i, istioEndpoint := range ep {
		var priority int
		// failoverPriority labels match
		for j, label := range priorityLabels {
			valueForProxy, ok := priorityLabelOverrides[label]
			if !ok {
				valueForProxy = proxyLabels[label]
			}
			if valueForProxy != istioEndpoint.Metadata.Labels[label] {
				priority = lowestPriority - j
				break
			}
		}
		priorityMap[priority] = append(priorityMap[priority], i)
	}

	// sort all priorities in increasing order.
	priorities := []int{}
	for priority := range priorityMap {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)

	out := make([]*envoy_config_endpoint_v3.LocalityLbEndpoints, len(priorityMap))
	for i, priority := range priorities {
		out[i] = nil // util.CloneLocalityLbEndpoint(ep.LocalityLbEndpoints)
		out[i].LbEndpoints = nil
		out[i].Priority = uint32(priority)
		var weight uint32
		for _, index := range priorityMap[priority] {
			out[i].LbEndpoints = append(out[i].LbEndpoints, ep.LocalityLbEndpoints.LbEndpoints[index])
			weight += 1 //ep.LocalityLbEndpoints.LbEndpoints[index].GetLoadBalancingWeight().GetValue()
		}
		// reset weight
		if weight > 0 {
			out[i].LoadBalancingWeight = &wrappers.UInt32Value{
				Value: weight,
			}
		} else {
			out[i].LoadBalancingWeight = nil
		}
	}

	return out
}
*/

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
	OurDestRulesByHostName krt.Index[string, networkingclient.DestinationRule]
}

func newDestruleIndex(c kube.Client, workloadNs string, workloadLabels map[string]string) *DestruleIndex {
	var d DestruleIndex

	ourNs := utils.GetPodNamespace()

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
		if labels.Instance(i.Spec.WorkloadSelector.MatchLabels).SubsetOf(workloadLabels) {
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

type DestinationRuleWrapper struct {
	*networkingclient.DestinationRule
}

func (s DestinationRuleWrapper) GetLabelSelector() map[string]string {
	return s.Spec.WorkloadSelector.MatchLabels
}

/*
func getDestruleFor(destinationRules krt.Collection[*networkingclient.DestinationRule], augmentedPods krt.Collection[augmentedPod]) {

	wrappedDestRules := krt.NewCollection(destinationRules, func(_ krt.HandlerContext, i *networkingclient.DestinationRule) *DestinationRuleWrapper {
		return &DestinationRuleWrapper{i}
	})
	destinationRulesByNamespace := krt.NewNamespaceIndex(wrappedDestRules)

	ourDestRules := krt.NewCollection(augmentedPods, func(kctx krt.HandlerContext, pod augmentedPod) *DestinationRuleWrapper {
		// make sure this either doesn't have selectors, or they select us:

		destrules := krt.Fetch(kctx, wrappedDestRules, krt.FilterIndex(destinationRulesByNamespace, pod.Namespace), krt.FilterSelects(pod.podLabels))
		if len(destrules) == 0 {
			return nil
		}
		oldestDestRule := slices.MinFunc(destrules, func(i DestinationRuleWrapper, j DestinationRuleWrapper) int {
			return i.CreationTimestamp.Time.Compare(j.CreationTimestamp.Time)
		})
		return &oldestDestRule
	})

}
*/

func snapshotPerClient(ucc krt.Collection[UniqlyConnectedClient],
	glooEndpoints krt.Collection[EndpointsForUpstream], mostXdsSnapshots krt.Collection[xdsSnapWrapper],
	mostXdsSnapshotsIndex krt.Index[types.NamespacedName, xdsSnapWrapper], wrappedDestRules krt.Collection[DestinationRuleWrapper],
	destinationRulesByNamespace krt.Index[string, DestinationRuleWrapper]) krt.Collection[xdsSnapWrapper] {

	xdsSnapshotsForUcc := krt.NewCollection(ucc, func(kctx krt.HandlerContext, ucc UniqlyConnectedClient) *xdsSnapWrapper {

		// get the proxy and the labels
		gwName, gwNamespace := roleToNameNamespace(ucc.Role)
		podLabels := ucc.Labels

		mostlySnaps := krt.Fetch(kctx, mostXdsSnapshots, krt.FilterIndex(mostXdsSnapshotsIndex, types.NamespacedName{Namespace: gwNamespace, Name: gwName}))
		if len(mostlySnaps) != 1 {
			return nil
		}
		mostlySnap := mostlySnaps[0]

		destrules := krt.Fetch(kctx, wrappedDestRules, krt.FilterIndex(destinationRulesByNamespace, gwNamespace), krt.FilterSelects(podLabels))
		endpoints := krt.Fetch(kctx, glooEndpoints)
		var endpointsProto []envoycache.Resource
		for _, ep := range endpoints {
			// host that would match the dest rule from the endpoints.
			// get the matching dest rule
			// get the lb info from the dest rules and call prioritize
			hostname := fromEndpoint(ep)
			lbInfo := getDestruleFor(destrules, hostname)
			var priorities *priorities
			if lbInfo != nil {
				priorities = newPriorities(lbInfo.failoverPriority, lbInfo.proxyLabels)
			}

			cla := prioritize(ep, lbInfo, priorities)
			endpointsProto = append(endpointsProto, resource.NewEnvoyResource(cla))
		}

		clustersVersion := mostlySnap.snap.Clusters.Version
		endpointsVersion := EnvoyCacheResourcesSetToFnvHash(endpointsProto)
		// fetch destrules with index and see if we have dest rules for us. if so modify the proxy cache key
		// if clusters are updated, provider a new version of the endpoints,
		// so the clusters are warm
		genericSnap := mostlySnap.snap
		mostlySnap.proxyKey = "TODO" // compute key that includes the labels
		mostlySnap.snap = &xds.EnvoySnapshot{
			Clusters:  genericSnap.Clusters,
			Endpoints: envoycache.NewResources(fmt.Sprintf("%v-%v", clustersVersion, endpointsVersion), endpointsProto),
			Routes:    genericSnap.Routes,
			Listeners: genericSnap.Listeners,
		}

		return &mostlySnap
	})
	return xdsSnapshotsForUcc

}

func roleToNameNamespace(role string) (name, namespace string) {
	panic("implement me")
}

func fromEndpoint(ep EndpointsForUpstream) string {
	panic("implement me")
}

func getDestruleFor(destrules []DestinationRuleWrapper, hostname string) *LBInfo {
	panic("implement me")
}
