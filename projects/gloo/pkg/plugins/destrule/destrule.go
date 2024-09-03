package destrule

import (
	"context"
	"slices"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"

	"github.com/solo-io/gloo/pkg/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
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

func NewPlugin(ctx context.Context) (*plugin, error) {
	p := &plugin{}

	return p, nil
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(params plugins.InitParams) {
}

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

	// use oldest
	oldestDestRule := slices.MinFunc(destrules, func(i networkingclient.DestinationRule, j networkingclient.DestinationRule) int {
		return i.CreationTimestamp.Time.Compare(j.CreationTimestamp.Time)
	})

	// apply settings from the oldest one
	// TODO !!

	return nil
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

func (d *DestruleIndex) destRuleCollection(c kube.Client) {
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
		// do not test if len(selectors) == 0, because that means no selectors.
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

}
