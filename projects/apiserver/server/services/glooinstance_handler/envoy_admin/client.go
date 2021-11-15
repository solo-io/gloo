package envoy_admin

import (
	"context"
	"fmt"

	envoy_admin_v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	"github.com/golang/protobuf/jsonpb"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/go-utils/contextutils"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"k8s.io/client-go/rest"
)

//go:generate mockgen -source ./client.go -destination mocks/mock_client.go

type EnvoyAdminClient interface {
	GetConfigs(ctx context.Context, glooInstance *rpc_edge_v1.GlooInstance, restClient rest.Interface) ([]*rpc_edge_v1.ConfigDump, error)
	GetHostsByUpstream(ctx context.Context, glooInstance *rpc_edge_v1.GlooInstance, restClient rest.Interface) (map[string]*rpc_edge_v1.HostList, error)
}

type envoyAdminClient struct {
}

func NewEnvoyAdminClient() EnvoyAdminClient {
	return &envoyAdminClient{}
}

func (c *envoyAdminClient) GetConfigs(ctx context.Context, glooInstance *rpc_edge_v1.GlooInstance, restClient rest.Interface) ([]*rpc_edge_v1.ConfigDump, error) {
	var configDumps []*rpc_edge_v1.ConfigDump

	for _, proxy := range glooInstance.Spec.Proxies { // Get a config dump for each proxy.
		if !proxy.GetReadConfigMulticlusterEnabled() {
			configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
				Name:  proxy.GetName(),
				Error: "No config dumps available! When installing Gloo Edge, the readConfig and readConfigMulticluster values must be set to true.",
			})
			continue
		}
		if proxy.GetAvailableReplicas() == 0 {
			configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
				Name:  proxy.GetName(),
				Error: "No replicas available for this proxy yet!",
			})
			continue
		}
		configProxyServicePath := getConfigDumpPath(proxy.GetName(), glooInstance.Spec.GetControlPlane().GetNamespace())
		result := restClient.Get().AbsPath(configProxyServicePath).Do(ctx)
		configDump, err := result.Raw()
		if err != nil {
			configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
				Name:  proxy.GetName(),
				Error: "No config dumps found!",
			})
		} else {
			configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
				Name: proxy.GetName(),
				Raw:  string(configDump),
			})
		}
	}
	return configDumps, nil
}

func (c *envoyAdminClient) GetHostsByUpstream(ctx context.Context, glooInstance *rpc_edge_v1.GlooInstance, restClient rest.Interface) (map[string]*rpc_edge_v1.HostList, error) {
	upstreamToHosts := make(map[string]*rpc_edge_v1.HostList)

	// get the clusters for each proxy
	for _, proxy := range glooInstance.Spec.Proxies {
		if !proxy.GetReadConfigMulticlusterEnabled() {
			contextutils.LoggerFrom(ctx).Debugf("readConfig and/or readConfigMulticluster is not enabled on proxy %s.%s", proxy.GetNamespace(), proxy.GetName())
			continue
		}
		if proxy.GetAvailableReplicas() == 0 {
			contextutils.LoggerFrom(ctx).Warnf("No replicas available for proxy %s.%s yet", proxy.GetNamespace(), proxy.GetName())
			continue
		}

		// make a request to the envoy admin clusters endpoint.
		// we need to request json in order to unmarshal it correctly
		clustersPath := getClustersPath(proxy.GetName(), glooInstance.Spec.GetControlPlane().GetNamespace())
		result := restClient.Get().AbsPath(clustersPath).Param("format", "json").Do(ctx)
		clustersJsonBytes, err := result.Raw()
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorf("Unable to get envoy clusters: %v", err)
			continue
		}

		// unmarshal the result into an envoy Clusters struct
		var clusters envoy_admin_v3.Clusters
		if err = jsonpb.UnmarshalString(string(clustersJsonBytes), &clusters); err != nil {
			contextutils.LoggerFrom(ctx).Errorf("Error unmarshalling clusters: %v", err)
			continue
		}

		// collect the host info for each cluster
		for _, clusterStatus := range clusters.GetClusterStatuses() {
			// for each cluster, translate its name back into an upstream key (namespace.name)
			// note, not all clusters will correspond to an actual upstream, but we don't verify it here
			clusterName := clusterStatus.GetName()
			upstreamRef, err := translator.ClusterToUpstreamRef(clusterName)
			if err != nil {
				contextutils.LoggerFrom(ctx).Errorf("Error getting upstream ref from cluster %s: %v", clusterName, err)
				continue
			}
			key := upstreamRef.Key()

			// make sure a map entry exists for this upstream key
			hostList, ok := upstreamToHosts[key]
			if !ok {
				hostList = &rpc_edge_v1.HostList{}
				upstreamToHosts[key] = hostList
			}

			// add each host to the list for this upstream
			for _, hostStatus := range clusterStatus.GetHostStatuses() {
				hostList.Hosts = append(hostList.Hosts, &rpc_edge_v1.Host{
					Address: hostStatus.GetAddress().GetSocketAddress().GetAddress(),
					Port:    hostStatus.GetAddress().GetSocketAddress().GetPortValue(),
					Weight:  hostStatus.GetWeight(),
					ProxyRef: &skv2_v1.ObjectRef{
						Name:      proxy.GetName(),
						Namespace: proxy.GetNamespace(),
					},
				})
			}
		}
	}

	return upstreamToHosts, nil
}

// Note that the proxy name should match the name of the service created for gloo
// here: https://github.com/solo-io/gloo/blob/9923fa18c792bf5f71fe58f5fad9ebce2a9a5dec/install/helm/gloo/templates/8-gateway-proxy-service.yaml#L89
func getConfigDumpPath(proxyName, namespace string) string {
	return fmt.Sprintf("/api/v1/namespaces/%s/services/%s/proxy/config_dump",
		namespace, defaults.GatewayProxyConfigDumpServiceName(proxyName))
}

func getClustersPath(proxyName, namespace string) string {
	return fmt.Sprintf("/api/v1/namespaces/%s/services/%s/proxy/clusters",
		namespace, defaults.GatewayProxyConfigDumpServiceName(proxyName))
}
