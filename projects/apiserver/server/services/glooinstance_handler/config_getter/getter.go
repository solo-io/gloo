package config_getter

import (
	"context"

	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"

	skv2_multicluster "github.com/solo-io/skv2/pkg/multicluster"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"k8s.io/client-go/discovery"
)

//go:generate mockgen -source ./getter.go -destination mocks/mock_getter.go

type EnvoyConfigDumpGetter interface {
	GetConfigs(ctx context.Context, glooInstance fedv1.GlooInstance) ([]*rpc_edge_v1.ConfigDump, error)
}

type configDumpGetter struct {
	clusterWatcher skv2_multicluster.ManagerSet
}

func (c *configDumpGetter) GetConfigs(ctx context.Context, glooInstance fedv1.GlooInstance) ([]*rpc_edge_v1.ConfigDump, error) {
	var configDumps []*rpc_edge_v1.ConfigDump
	mgr, err := c.clusterWatcher.Cluster(glooInstance.Spec.GetCluster())
	if err != nil {
		return nil, err
	}
	clientset, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	for _, proxy := range glooInstance.Spec.Proxies { // Get a config dump for each proxy.
		if !proxy.GetReadConfigMulticlusterEnabled() {
			configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
				Name:  proxy.GetName(),
				Error: "No config dumps available! When installing Gloo, the readConfig and readConfigMulticluster values must be set to true.",
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
		configProxyServicePath := kubeAPIConfigProxyServicePath(proxy.GetName(), glooInstance.Spec.GetControlPlane().GetNamespace())
		result := clientset.RESTClient().Get().AbsPath(configProxyServicePath).Do(ctx)
		configDump, err := result.Raw()
		if err != nil {
			configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
				Name:  proxy.GetName(),
				Error: "No config dumps found!",
			})
		}
		configDumps = append(configDumps, &rpc_edge_v1.ConfigDump{
			Name: proxy.GetName(),
			Raw:  string(configDump),
		})
	}
	return configDumps, nil
}

func NewEnvoyConfigDumpGetter(clusterWatcher skv2_multicluster.ManagerSet) EnvoyConfigDumpGetter {
	return &configDumpGetter{
		clusterWatcher: clusterWatcher,
	}
}

// Note that the proxy name should match the name of the service created for gloo
// here: https://github.com/solo-io/gloo/blob/9923fa18c792bf5f71fe58f5fad9ebce2a9a5dec/install/helm/gloo/templates/8-gateway-proxy-service.yaml#L89
func kubeAPIConfigProxyServicePath(proxyName, namespace string) string {
	return "/api/v1/namespaces/" + namespace + "/services/" + proxyName + "-config-dump-service/proxy/config_dump"
}
