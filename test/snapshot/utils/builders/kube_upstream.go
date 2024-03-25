package builders

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeUpstreamBuilder contains options for building Upstreams
type KubeUpstreamBuilder struct {
	name      string
	namespace string
	labels    map[string]string

	discoveryMetadata *gloov1.DiscoveryMetadata
	upstreamType      *gloov1.UpstreamSpec_Kube

	sniPattern   sniPattern
	healthChecks []*core.HealthCheck
}

type sniPattern int

const (
	noSni sniPattern = iota
	uniqueSni
	consistentSni
)

func NewUpstreamBuilder() *KubeUpstreamBuilder {
	return &KubeUpstreamBuilder{}
}

func (b *KubeUpstreamBuilder) WithName(name string) *KubeUpstreamBuilder {
	b.name = name
	return b
}

func (b *KubeUpstreamBuilder) WithNamespace(namespace string) *KubeUpstreamBuilder {
	b.namespace = namespace
	return b
}

func (b *KubeUpstreamBuilder) WithLabel(key, value string) *KubeUpstreamBuilder {
	b.labels[key] = value
	return b
}

func (b *KubeUpstreamBuilder) WithUniqueSni() *KubeUpstreamBuilder {
	b.sniPattern = uniqueSni
	return b
}

func (b *KubeUpstreamBuilder) WithConsistentSni() *KubeUpstreamBuilder {
	b.sniPattern = consistentSni
	return b
}

func (b *KubeUpstreamBuilder) WithHealthChecks(healthChecks []*core.HealthCheck) *KubeUpstreamBuilder {
	b.healthChecks = healthChecks
	return b
}

func (b *KubeUpstreamBuilder) WithDiscoveryMetadata(discoveryMeta *gloov1.DiscoveryMetadata) *KubeUpstreamBuilder {
	b.discoveryMetadata = discoveryMeta
	return b
}

func (b *KubeUpstreamBuilder) WithKubeUpstream(kubeUpstream *gloov1.UpstreamSpec_Kube) *KubeUpstreamBuilder {
	b.upstreamType = kubeUpstream
	return b
}

func (b *KubeUpstreamBuilder) Build() *gloov1.Upstream {
	upstream := &gloov1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
			Labels:    b.labels,
		},
		Spec: gloov1.UpstreamSpec{
			DiscoveryMetadata: b.discoveryMetadata,
			UpstreamType:      b.upstreamType,
		},
	}
	return upstream
}
