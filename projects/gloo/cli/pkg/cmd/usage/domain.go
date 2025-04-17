package usage

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	v1 "k8s.io/api/core/v1"
)

type K8sClusterInfo struct {
	Nodes    []v1.Node    `json:"nodes"`
	Pods     []v1.Pod     `json:"pods"`
	Services []v1.Service `json:"services"`
}
type UsageStats struct {
	GlooFeatureUsage map[API]*GlooFeatureUsage  `json:"glooFeatureUsage" yaml:"glooFeatureUsage"`
	GlooProxyStats   map[string]*GlooProxyStats `json:"glooProxyStats" yaml:"glooProxyStats"`
	KubernetesStats  *KubernetesStats           `json:"kubernetesStats" yaml:"kubernetesStats"`
}

type GlooFeatureUsage struct {
	APICounts            map[string]int                `json:"apiCounts" yaml:"apiCounts"`
	FeatureCountPerProxy map[string]*ProxyFeatureCount `json:"featureCountPerProxy" yaml:"featureCountPerProxy"`
}

type ProxyFeatureCount struct {
	FeatureCount map[FeatureType]int `json:"featureCount" yaml:"featureCount"`
}

type GlooProxyStats struct {
	GlooProxyMetrics *EnvoyMetrics   `json:"glooProxyMetrics" yaml:"glooProxyMetrics"`
	GlooProxyStats   []*UpstreamStat `json:"glooProxyStats" yaml:"glooProxyStats"`
}

type KubernetesStats struct {
	Pods          int            `json:"pods" yaml:"pods"`
	NodeResources *NodeResources `json:"nodeResources" yaml:"nodeResources"`
	Services      int            `json:"services" yaml:"services"`
}
type ProxyData struct {
	Name         string
	Namespace    string
	EnvoyMetrics *EnvoyMetrics
}

type Inputs struct {
	GlooEdgeConfigs *snapshot.Instance    `json:"glooEdgeConfigs"`
	K8sClusterInfo  *K8sClusterInfo       `json:"k8sClusterInfo"`
	ProxyStats      map[string]*ProxyInfo `json:"proxyStats"`
}

type NodeResources struct {
	TotalCPUCores int64 `json:"totalCPUCores" yaml:"totalCPUCores"`
	TotalMemoryGb int64 `json:"totalMemoryGb" yaml:"totalMemoryGb"`
	Nodes         int   `json:"nodes" yaml:"nodes"`
}

func (u *UsageStats) Print(format string) error {
	//featureCount := map[FeatureType]int{}

	if format == "json" {
		err := json.NewEncoder(os.Stdout).Encode(u)
		return err
	}

	if format == "yaml" {
		err := yaml.NewEncoder(os.Stdout).Encode(u)
		if err != nil {
			return err
		}
	}

	// text output format

	//fmt.Printf("\nK8s Resources:\n")
	//fmt.Printf("\tNodes: %d\n", u.KubernetesStats.NodeResources.Nodes)
	//fmt.Printf("\tPods: %d\n", u.KubernetesStats.Pods)
	//fmt.Printf("\tServices: %d\n", u.KubernetesStats.Services)
	//// Print node resource information
	//fmt.Printf("\nNode Resources:\n")
	//fmt.Printf("Total CPU Capacity: %.2f cores\n", float64(u.KubernetesStats.NodeResources.TotalCapacityCPU)/1000)
	//fmt.Printf("Total Memory Capacity: %.2f GB\n", float64(u.KubernetesStats.NodeResources.TotalCapacityMemory)/math.Pow(1024, 3))
	//
	//fmt.Printf("Gloo Edge APIs: \n")
	//fmt.Printf("\tGloo Gateways: %d\n", len(u.GlooEdgeConfigs.GlooGateways()))
	//fmt.Printf("\tVirtualService: %d\n", len(u.GlooEdgeConfigs.VirtualServices()))
	//fmt.Printf("\tRouteTables: %d\n", len(u.GlooEdgeConfigs.RouteTables()))
	//fmt.Printf("\tUpstreams: %d\n", len(u.GlooEdgeConfigs.Upstreams()))
	//fmt.Printf("\tVirtualHostOptions: %d\n", len(u.GlooEdgeConfigs.VirtualHostOptions()))
	//fmt.Printf("\tListenerOptions: %d\n", len(u.GlooEdgeConfigs.ListenerOptions()))
	//fmt.Printf("\tHTTPListenerOptions: %d\n", len(u.GlooEdgeConfigs.HTTPListenerOptions()))
	//fmt.Printf("\tAuthConfigs: %d\n", len(u.GlooEdgeConfigs.AuthConfigs()))
	//fmt.Printf("\tSettings: %d\n", len(u.GlooEdgeConfigs.Settings()))
	//
	//fmt.Printf("\nGateway API APIs: \n")
	//fmt.Printf("\tGateways: %d\n", len(u.GlooEdgeConfigs.Gateways()))
	//fmt.Printf("\tListenerSets: %d\n", len(u.GlooEdgeConfigs.ListenerSets()))
	//fmt.Printf("\tHTTPRoutes: %d\n", len(u.GlooEdgeConfigs.HTTPRoutes()))
	//fmt.Printf("\tDirectResponses: %d\n", len(u.GlooEdgeConfigs.DirectResponses()))
	//fmt.Printf("\tGatewayParameters: %d\n", len(u.GlooEdgeConfigs.GatewayParameters()))
	//
	//fmt.Printf("\nTotal Features Used Per API\n")
	//fmt.Printf("\tGloo Edge API: %d\n", len(u.GlooFeatureUsage[GlooEdgeAPI]))
	//fmt.Printf("\tGateway API: %d\n", len(u.GlooFeatureUsage[GatewayAPI]))
	//fmt.Printf("\tkGateway API: %d\n\n", len(u.GlooFeatureUsage[KGatewayAPI]))
	//
	//for api, features := range u.GlooFeatureUsage {
	//	fmt.Printf("API: %s", api)
	//
	//	// organize by category
	//	categories := map[Category]map[FeatureType]int{}
	//
	//	// group all the stats by their codes
	//	for _, feature := range features {
	//		//featureCount[stat.Type]++
	//		if categories[feature.Metadata.Category] == nil {
	//			categories[feature.Metadata.Category] = map[FeatureType]int{}
	//		}
	//		categories[feature.Metadata.Category][feature.Type]++
	//
	//	}
	//	for category, features := range categories {
	//		fmt.Printf("\n\tCategory: %s", category)
	//		//sort the features
	//		//add all the features to a list for naming
	//		featureList := []FeatureType{}
	//		for key := range features {
	//			featureList = append(featureList, key)
	//		}
	//
	//		sort.Slice(featureList, func(i, j int) bool {
	//			return featureList[i] < featureList[j]
	//		})
	//		for _, feature := range featureList {
	//			fmt.Printf("\n\t\t%s: %d", feature, features[feature])
	//		}
	//	}
	//}
	return nil
}
