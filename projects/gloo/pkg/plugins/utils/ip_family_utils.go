package utils

import (
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

// FIXME: (kasunt) it would be good to change this to V6_PREFERRED in the future when it becomes available in the API
func TranslateIpFamily(ipFamily v1.DnsIpFamily) envoy_config_cluster_v3.Cluster_DnsLookupFamily {
	switch ipFamily {
	case v1.DnsIpFamily_V4_ONLY:
		return envoy_config_cluster_v3.Cluster_V4_ONLY
	case v1.DnsIpFamily_V6_ONLY:
		return envoy_config_cluster_v3.Cluster_V6_ONLY
	case v1.DnsIpFamily_V4_PREFERRED:
		return envoy_config_cluster_v3.Cluster_V4_PREFERRED
	case v1.DnsIpFamily_ALL:
		return envoy_config_cluster_v3.Cluster_ALL
	default:
		return envoy_config_cluster_v3.Cluster_V4_ONLY
	}
}
