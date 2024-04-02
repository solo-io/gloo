package types

import (
	istiov1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	kindv1alpha4 "sigs.k8s.io/kind/pkg/apis/config/v1alpha4"

	"github.com/solo-io/gloo/pkg/utils/maputils"
)

type Cluster struct {
	// Name of the cluster.
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// IstioOperator configuration
	IstioOperators []*istiov1alpha1.IstioOperator `yaml:"istioOperators,omitempty" json:"istioOperators,omitempty"`

	// (Optional) KinD configuration if supplied will be used to create the cluster.
	KindConfig *kindv1alpha4.Cluster `yaml:"kindConfig,omitempty" json:"kindConfig,omitempty"`

	// Will be populated from values data in chart section of cluster.
	GlooEdge *GlooEdge `yaml:"glooEdge,omitempty" json:"glooEdge,omitempty"`

	// Namespaces to create on the cluster.
	Namespaces []*corev1.Namespace `yaml:"namespaces,omitempty" json:"namespaces,omitempty"`

	// List of helm charts to install
	Charts []*Chart `yaml:"charts,omitempty" json:"charts,omitempty"`

	// Test apps to install on the cluster.
	Apps []*App `yaml:"apps,omitempty" json:"apps,omitempty"`
}

func (c *Cluster) GetImages() []string {
	if c.GlooEdge == nil {
		return nil
	}
	return c.GlooEdge.Images()
}

func (c *Cluster) UpdateChart(name string, values ...map[string]any) {
	for i, chart := range c.Charts {
		if chart.Name != name {
			continue
		}

		base := &c.Charts[i].Values

		for _, value := range values {
			*base = maputils.MergeValueMaps(*base, value)
		}

		(*c.Charts[i]).Values = *base
	}
}

func (c *Cluster) GetChart(name string) *Chart {
	for _, chart := range c.Charts {
		if chart.Name != name {
			continue
		}
		return chart
	}
	return nil
}

// GetPrioritizedCharts returns the charts in the order they should be installed.
func (c *Cluster) GetPrioritizedCharts() []*Chart {
	var prioritized []*Chart
	for _, chart := range c.Charts {
		if chart.Prioritize {
			prioritized = append(prioritized, chart)
		}
	}

	return prioritized
}

func (c *Cluster) GetUnprioritizedCharts() []*Chart {
	var unprioritized []*Chart
	for _, chart := range c.Charts {
		if !chart.Prioritize {
			unprioritized = append(unprioritized, chart)
		}
	}

	return unprioritized
}
