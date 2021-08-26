package config_getter

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"k8s.io/client-go/discovery"
)

//go:generate mockgen -source ./getter.go -destination mocks/mock_getter.go

type EnvoyConfigDumpGetter interface {
	GetConfigs(ctx context.Context, glooInstance *rpc_edge_v1.GlooInstance, discoveryClient discovery.DiscoveryClient) ([]*rpc_edge_v1.ConfigDump, error)
}

type configDumpGetter struct {
}

func NewEnvoyConfigDumpGetter() EnvoyConfigDumpGetter {
	return &configDumpGetter{}
}

func (c *configDumpGetter) GetConfigs(ctx context.Context, glooInstance *rpc_edge_v1.GlooInstance, discoveryClient discovery.DiscoveryClient) ([]*rpc_edge_v1.ConfigDump, error) {
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
		configProxyServicePath := kubeAPIConfigProxyServicePath(proxy.GetName(), glooInstance.Spec.GetControlPlane().GetNamespace())
		result := discoveryClient.RESTClient().Get().AbsPath(configProxyServicePath).Do(ctx)
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

// Note that the proxy name should match the name of the service created for gloo
// here: https://github.com/solo-io/gloo/blob/9923fa18c792bf5f71fe58f5fad9ebce2a9a5dec/install/helm/gloo/templates/8-gateway-proxy-service.yaml#L89
func kubeAPIConfigProxyServicePath(proxyName, namespace string) string {
	return fmt.Sprintf("/api/v1/namespaces/%s/services/%s/proxy/config_dump",
		namespace, defaults.GatewayProxyConfigDumpServiceName(proxyName))
}
