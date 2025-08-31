package utils

import (
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

// GetPreferredIpFamily utility to manage the preferred ip family
// - if local override is set always set it appropriately
// - if a global value is present and no local override is done then get the global value
// - otherwise if forceToV4Family is set it forces to V4 family. This flag is mainly introduced to restore the original behavior prior to introducing Ip Family management
func GetPreferredIpFamily(out *envoy_config_cluster_v3.Cluster, global *v1.UpstreamOptions, local v1.DnsIpFamily, forceToV4Family bool) {
	if local != v1.DnsIpFamily_DEFAULT {
		out.DnsLookupFamily = translateIpFamily(local)
	} else if global != nil && global.GetDnsLookupIpFamily() != v1.DnsIpFamily_DEFAULT {
		out.DnsLookupFamily = translateIpFamily(global.GetDnsLookupIpFamily())
	} else if forceToV4Family {
		out.DnsLookupFamily = envoy_config_cluster_v3.Cluster_V4_ONLY
	}
}

// translateIpFamily maps the given upstream DnsIpFamily to the equivalent cluster family in envoy
func translateIpFamily(ipFamily v1.DnsIpFamily) envoy_config_cluster_v3.Cluster_DnsLookupFamily {
	switch ipFamily {
	case v1.DnsIpFamily_V4_ONLY:
		return envoy_config_cluster_v3.Cluster_V4_ONLY
	case v1.DnsIpFamily_V6_ONLY:
		return envoy_config_cluster_v3.Cluster_V6_ONLY
	case v1.DnsIpFamily_V4_PREFERRED:
		return envoy_config_cluster_v3.Cluster_V4_PREFERRED
	case v1.DnsIpFamily_V6_PREFERRED:
		// FIXME: (kasunt) it would be good to change this to V6_PREFERRED in the future when it becomes available in the API
		return envoy_config_cluster_v3.Cluster_AUTO
	case v1.DnsIpFamily_DUAL_IP_FAMILY:
		return envoy_config_cluster_v3.Cluster_ALL
	default:
		return envoy_config_cluster_v3.Cluster_V4_ONLY
	}
}
