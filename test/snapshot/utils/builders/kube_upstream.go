package builders

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KubeUpstreamBuilder contains options for building Upstreams
type KubeUpstreamBuilder struct {
	name      string
	namespace string
	labels    map[string]string

	DiscoveryMetadata *gloov1.DiscoveryMetadata
	UpstreamType      *gloov1.UpstreamSpec_Kube

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
	b.DiscoveryMetadata = discoveryMeta
	return b
}

func (b *KubeUpstreamBuilder) WithKubeUpstream(kubeUpstream *gloov1.UpstreamSpec_Kube) *KubeUpstreamBuilder {
	b.UpstreamType = kubeUpstream
	return b
}

func (b *KubeUpstreamBuilder) Build() *gloov1.Upstream {
	upstream := &gloov1.Upstream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "httpbin-htppbin-8000",
			Namespace: gloodefaults.GlooSystem,
		},
		Spec: gloov1.UpstreamSpec{
			DiscoveryMetadata: b.DiscoveryMetadata,
			UpstreamType:      b.UpstreamType,
		},
	}
	return upstream
}
