package usage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/fatih/color"
	"github.com/rodaine/table"
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
	APICounts            map[string]int                        `json:"apiCounts" yaml:"apiCounts"`
	FeatureCountPerProxy map[string]*ProxyFeatureCountCategory `json:"featureCountPerProxy" yaml:"featureCountPerProxy"`
}
type ProxyFeatureCountCategory struct {
	Categories map[Category]ProxyFeatureCount `json:"categories" yaml:"categories"`
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
	Name         string        `json:"name" yaml:"name"`
	Namespace    string        `json:"namespace" yaml:"namespace"`
	EnvoyMetrics *EnvoyMetrics `json:"envoyMetrics" yaml:"envoyMetrics"`
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
		if err != nil {
			return err
		}
		return nil
	}

	if format == "yaml" {
		err := yaml.NewEncoder(os.Stdout).Encode(u)
		if err != nil {
			return err
		}
		return nil
	}

	// text output format
	if u.KubernetesStats != nil {
		fmt.Printf("\nK8s Resources:")
		fmt.Printf("\n\tNodes: %d", u.KubernetesStats.NodeResources.Nodes)
		fmt.Printf("\n\tPods: %d", u.KubernetesStats.Pods)
		fmt.Printf("\n\tServices: %d", u.KubernetesStats.Services)
		// Print node resource information
		fmt.Printf("\n\n\tNode Resources:")
		fmt.Printf("\n\t\tTotal CPU Capacity: %.2f cores", float64(u.KubernetesStats.NodeResources.TotalCPUCores))
		fmt.Printf("\n\t\tTotal Memory Capacity: %.2f GB", float64(u.KubernetesStats.NodeResources.TotalMemoryGb))
	}

	if u.GlooFeatureUsage[GlooEdgeAPI] != nil {
		fmt.Printf("\n\nGloo Edge APIs:")
		fmt.Printf("\n\tGloo Gateways: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["Gloo Gateways"])
		fmt.Printf("\n\tVirtualService: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["VirtualService"])
		fmt.Printf("\n\tRouteTables: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["RouteTables"])
		fmt.Printf("\n\tUpstreams: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["Upstreams"])
		fmt.Printf("\n\tVirtualHostOptions: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["VirtualHostOptions"])
		fmt.Printf("\n\tListenerOptions: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["ListenerOptions"])
		fmt.Printf("\n\tHTTPListenerOptions: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["HTTPListenerOptions"])
		fmt.Printf("\n\tAuthConfigs: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["AuthConfigs"])
		fmt.Printf("\n\tSettings: %d", u.GlooFeatureUsage[GlooEdgeAPI].APICounts["Settings"])
	}

	if u.GlooFeatureUsage[GatewayAPI] != nil {
		fmt.Printf("\n\nGateway API APIs:")
		fmt.Printf("\n\tGateways: %d", u.GlooFeatureUsage[GatewayAPI].APICounts["Gateways"])
		fmt.Printf("\n\tListenerSets: %d", u.GlooFeatureUsage[GatewayAPI].APICounts["ListenerSets"])
		fmt.Printf("\n\tHTTPRoutes: %d", u.GlooFeatureUsage[GatewayAPI].APICounts["HTTPRoutes"])
		fmt.Printf("\n\tDirectResponses: %d", u.GlooFeatureUsage[GatewayAPI].APICounts["DirectResponses"])
		fmt.Printf("\n\tGatewayParameters: %d", u.GlooFeatureUsage[GatewayAPI].APICounts["GatewayParameters"])
	}

	for api, features := range u.GlooFeatureUsage {
		fmt.Printf("\n\nAPI: %s", api)
		for proxyName, feature := range features.FeatureCountPerProxy {
			fmt.Printf("\n\t%s:", proxyName)
			for category, featureCount := range feature.Categories {
				fmt.Printf("\n\t\t%s:", category)
				for featureType, count := range featureCount.FeatureCount {
					fmt.Printf("\n\t\t\t%s: %d", featureType, count)
				}
			}
		}
	}

	if u.GlooProxyStats != nil {
		fmt.Printf("\n\nProxy Stats:")
		for proxyName, stats := range u.GlooProxyStats {
			fmt.Printf("\n\tProxy: %s", proxyName)
			if stats.GlooProxyMetrics != nil {
				fmt.Printf("\n\t\tGloo 2xx Responses: %v", stats.GlooProxyMetrics.Total2xxResponses)
				fmt.Printf("\n\t\tGloo 3xx Responses: %v", stats.GlooProxyMetrics.Total3xxResponses)
				fmt.Printf("\n\t\tGloo 4xx Responses: %v", stats.GlooProxyMetrics.Total4xxResponses)
				fmt.Printf("\n\t\tGloo 5xx Responses: %v", stats.GlooProxyMetrics.Total5xxResponses)
				fmt.Printf("\n\t\tUptime: %v", time.Duration(stats.GlooProxyMetrics.UptimeSeconds)*time.Second)
				fmt.Printf("\n\t\tAverage Responses Per Second: %v", stats.GlooProxyMetrics.AverageResponsesPerSecond)
			}
			if len(stats.GlooProxyStats) > 0 {
				headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
				columnFmt := color.New(color.FgYellow).SprintfFunc()
				tbl := table.New("Upstream", "IP Address", "Port", "Rq Success", "Rq Error", "Cx Active", "Cx Connect Fail")
				tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
				for _, stat := range stats.GlooProxyStats {
					tbl.AddRow(stat.Upstream, stat.IPAddress, stat.Port, stat.RqSuccess, stat.RqError, stat.CxActive, stat.CxConnectFail)
				}
				fmt.Printf("\n------------%s Upstream Calls ------------\n", proxyName)
				tbl.Print()
				fmt.Printf("--------------------------------------\n")
			}
		}
	}

	return nil
}
