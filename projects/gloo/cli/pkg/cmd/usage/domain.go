package usage

import (
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	v1 "k8s.io/api/core/v1"
	"sort"
)

type K8sClusterInfo struct {
	Nodes    []v1.Node
	Pods     []v1.Pod
	Services []v1.Service
}
type UsageStats struct {
	ProxyData        *ProxyData
	GlooFeatureUsage *UsageStats
	GlooProxyMetrics *EnvoyMetrics
	GlooProxyStats   []*UpstreamStat
	GlooEdgeConfigs  *snapshot.Instance
}

type UsageInputs struct {
	GlooEdgeConfigs *snapshot.Instance
	K8sClusterInfo  *K8sClusterInfo
	ProxyStats      map[string]*ProxyInfo
}

type NodeResources struct {
	TotalCapacityCPU    int64
	TotalCapacityMemory int64
}

func (u *UsageStats) Print() {
	//featureCount := map[FeatureType]int{}

	fmt.Printf("Gloo Edge APIs: \n")
	fmt.Printf("\tGloo Gateways: %d\n", len(u.GlooEdgeConfigs.GlooGateways()))
	fmt.Printf("\tVirtualService: %d\n", len(u.GlooEdgeConfigs.VirtualServices()))
	fmt.Printf("\tRouteTables: %d\n", len(u.GlooEdgeConfigs.RouteTables()))
	fmt.Printf("\tUpstreams: %d\n", len(u.GlooEdgeConfigs.Upstreams()))
	fmt.Printf("\tVirtualHostOptions: %d\n", len(u.GlooEdgeConfigs.VirtualHostOptions()))
	fmt.Printf("\tListenerOptions: %d\n", len(u.GlooEdgeConfigs.ListenerOptions()))
	fmt.Printf("\tHTTPListenerOptions: %d\n", len(u.GlooEdgeConfigs.HTTPListenerOptions()))
	fmt.Printf("\tAuthConfigs: %d\n", len(u.GlooEdgeConfigs.AuthConfigs()))
	fmt.Printf("\tSettings: %d\n", len(u.GlooEdgeConfigs.Settings()))

	fmt.Printf("\nGateway API APIs: \n")
	fmt.Printf("\tGateways: %d\n", len(u.GlooEdgeConfigs.Gateways()))
	fmt.Printf("\tListenerSets: %d\n", len(u.GlooEdgeConfigs.ListenerSets()))
	fmt.Printf("\tHTTPRoutes: %d\n", len(u.GlooEdgeConfigs.HTTPRoutes()))
	fmt.Printf("\tDirectResponses: %d\n", len(u.GlooEdgeConfigs.DirectResponses()))
	fmt.Printf("\tGatewayParameters: %d\n", len(u.GlooEdgeConfigs.GatewayParameters()))

	fmt.Printf("\nTotal Features Used Per API\n")
	fmt.Printf("\tGloo Edge API: %d\n", len(u.stats[GlooEdgeAPI]))
	fmt.Printf("\tGateway API: %d\n", len(u.stats[GatewayAPI]))
	fmt.Printf("\tkGateway API: %d\n\n", len(u.stats[KGatewayAPI]))

	for api, stats := range u.GlooProxyStats {
		fmt.Printf("API: %s", api)

		// organize by category
		categories := map[Category]map[FeatureType]int{}

		// group all the stats by their codes
		for _, stat := range stats {
			//featureCount[stat.Type]++
			if categories[stat.Metadata.Category] == nil {
				categories[stat.Metadata.Category] = map[FeatureType]int{}
			}
			categories[stat.Metadata.Category][stat.Type]++

		}
		for category, features := range categories {
			fmt.Printf("\n\tCategory: %s", category)
			//sort the features
			//add all the features to a list for naming
			featureList := []FeatureType{}
			for key := range features {
				featureList = append(featureList, key)
			}

			sort.Slice(featureList, func(i, j int) bool {
				return featureList[i] < featureList[j]
			})
			for _, feature := range featureList {
				fmt.Printf("\n\t\t%s: %d", feature, features[feature])
			}
		}
	}

	//// Calculate node resources
	//nodeResources, err := calculateNodeResources(inputs.K8sClusterInfo.Nodes)
	//if err != nil {
	//	return err
	//}
	//fmt.Printf("\nK8s Resources:\n")
	//fmt.Printf("\tNodes: %d\n", len(inputs.K8sClusterInfo.Nodes))
	//fmt.Printf("\tPods: %d\n", len(inputs.K8sClusterInfo.Pods))
	//fmt.Printf("\tServices: %d\n", len(inputs.K8sClusterInfo.Services))
	//// Print node resource information
	//fmt.Printf("\nNode Resources:\n")
	//fmt.Printf("Total CPU Capacity: %.2f cores\n", float64(nodeResources.TotalCapacityCPU)/1000)
	//fmt.Printf("Total Memory Capacity: %.2f GB\n", float64(nodeResources.TotalCapacityMemory)/math.Pow(1024, 3))

}

type UsageMetadata struct {
	Name      string
	Namespace string
	Kind      string
	Category  Category
	API       API
}

type UsageStat struct {
	Type     FeatureType
	Metadata UsageMetadata
}

func (u *UsageStats) AddUsageStat(stat *UsageStat) {
	if stat.Metadata.API == "" {
		stat.Metadata.API = GlooEdgeAPI
	}
	if u.stats == nil {
		u.stats = make(map[API][]*UsageStat)
	}
	if u.stats[stat.Metadata.API] == nil {
		u.stats[stat.Metadata.API] = []*UsageStat{}
	}
	u.stats[stat.Metadata.API] = append(u.stats[stat.Metadata.API], stat)
}
