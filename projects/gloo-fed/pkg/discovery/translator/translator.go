package translator

import (
	"context"

	checker2 "github.com/solo-io/solo-projects/projects/gloo/pkg/utils/checker"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/images"
	locality2 "github.com/solo-io/solo-projects/projects/gloo/pkg/utils/locality"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/proxy"

	k8s_core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stringutils"
	sk_sets "github.com/solo-io/skv2/contrib/pkg/sets/v2"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	enterprise_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/check"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	gateway_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/gateway.solo.io/v1/check"
	gloo_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/gloo.solo.io/v1/check"
	ratelimit_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/ratelimit.api.solo.io/v1alpha1/check"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:generate mockgen -source ./translator.go -destination ./mocks/mock_translator.go

type GlooInstanceTranslator interface {
	// FromSnapshot translates the given snapshot into a list of GlooInstances.
	FromSnapshot(ctx context.Context, snapshot input.Snapshot) []*fedv1.GlooInstance
}

type translator struct {
	writeNamespace string
	mcClientSet    k8s_core_v1.MulticlusterClientset
}

func NewTranslator(
	writeNamespace string,
	mcClientSet k8s_core_v1.MulticlusterClientset,
) GlooInstanceTranslator {
	return &translator{
		writeNamespace: writeNamespace,
		mcClientSet:    mcClientSet,
	}
}

// Detect all GlooInstances from a multicluster snapshot of all the deployments
func (t *translator) FromSnapshot(ctx context.Context, snapshot input.Snapshot) []*fedv1.GlooInstance {
	var instances []*fedv1.GlooInstance
	for _, deployment := range snapshot.Deployments().List() {
		tag, isEnterprise, isControlPlane := images.GetControlPlaneImage(deployment)
		if !isControlPlane {
			// There should be one Gloo Instance per cluster's control plane deployment, so skip all other deployments.
			continue
		}

		cluster := deployment.GetClusterName()
		client, err := t.mcClientSet.Cluster(cluster)
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorf("couldn't get client for cluster %v", cluster)
			return nil
		}
		ipFinder := locality2.NewExternalIpFinder(cluster, client.Pods(), client.Nodes())
		localityFinder := locality2.NewLocalityFinder(client.Nodes(), client.Pods())
		// Ignore error that may occur when listing nodes
		region, _ := localityFinder.GetRegion(ctx)

		settings := getSettings(snapshot.Settings(), deployment.Namespace, deployment.ClusterName)
		var watchedNamespaces []string
		if settings != nil {
			watchedNamespaces = settings.Spec.WatchNamespaces
		}

		instance := &fedv1.GlooInstance{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: t.writeNamespace,
				Name:      GetGlooInstanceName(cluster, deployment.Namespace),
			},
			Spec: types.GlooInstanceSpec{
				Cluster:      cluster,
				IsEnterprise: isEnterprise,
				ControlPlane: &types.GlooInstanceSpec_ControlPlane{
					Version:           tag,
					Namespace:         deployment.Namespace,
					WatchedNamespaces: watchedNamespaces,
				},
				Region: region,
				Check:  &types.GlooInstanceSpec_Check{},
			},
		}

		// Proxy and admin info
		instance.Spec.Proxies = convertProxies(proxy.GetGlooInstanceProxies(ctx, cluster, deployment.Namespace, snapshot.Deployments(), snapshot.DaemonSets(), snapshot.Services(), localityFinder, ipFinder))
		instance.Spec.Admin = getAdminInfo(instance.Spec.ControlPlane.GetNamespace(), watchedNamespaces, instance.Spec.Proxies)

		// Check -- workloads
		instance.Spec.Check.Pods = convertSummary(checker2.GetPodsSummary(ctx, snapshot.Pods(), instance.Spec.ControlPlane.GetNamespace(), cluster))
		instance.Spec.Check.Deployments = convertSummary(checker2.GetDeploymentsSummary(ctx, snapshot.Deployments(), instance.Spec.ControlPlane.GetNamespace(), cluster))

		// Check -- resources
		// Gloo
		instance.Spec.Check.Settings = gloo_check.GetSettingsSummary(ctx, snapshot.Settings(), watchedNamespaces, cluster)
		instance.Spec.Check.Upstreams = gloo_check.GetUpstreamSummary(ctx, snapshot.Upstreams(), watchedNamespaces, cluster)
		instance.Spec.Check.UpstreamGroups = gloo_check.GetUpstreamGroupSummary(ctx, snapshot.UpstreamGroups(), watchedNamespaces, cluster)
		instance.Spec.Check.Proxies = gloo_check.GetProxySummary(ctx, snapshot.Proxies(), watchedNamespaces, cluster)
		instance.Spec.Check.AuthConfigs = enterprise_check.GetAuthConfigSummary(ctx, snapshot.AuthConfigs(), watchedNamespaces, cluster)
		instance.Spec.Check.RateLimitConfigs = ratelimit_check.GetRateLimitConfigSummary(ctx, snapshot.RateLimitConfigs(), watchedNamespaces, cluster)

		// Gateway
		instance.Spec.Check.Gateways = gateway_check.GetGatewaySummary(ctx, snapshot.Gateways(), watchedNamespaces, cluster)
		instance.Spec.Check.VirtualServices = gateway_check.GetVirtualServiceSummary(ctx, snapshot.VirtualServices(), watchedNamespaces, cluster)
		instance.Spec.Check.RouteTables = gateway_check.GetRouteTableSummary(ctx, snapshot.RouteTables(), watchedNamespaces, cluster)
		instance.Spec.Check.MatchableHttpGateways = gateway_check.GetMatchableHttpGatewaySummary(ctx, snapshot.MatchableHttpGateways(), watchedNamespaces, cluster)

		instances = append(instances, instance)
	}

	return instances
}

func getAdminInfo(installNamespace string, watchedNamespaces []string, proxies []*types.GlooInstanceSpec_Proxy) *types.GlooInstanceSpec_Admin {
	proxy := getProxyForGlooInstance(installNamespace, proxies)
	proxyId := &skv2v1.ObjectRef{
		Name:      proxy.GetName(),
		Namespace: proxy.GetNamespace(),
	}

	return &types.GlooInstanceSpec_Admin{
		WriteNamespace: getWriteNamespaceForGlooInstance(installNamespace, watchedNamespaces),
		ProxyId:        proxyId,
	}
}

/*
	GetWriteNamespaceForGlooInstance returns the write namespace for a gloo instance based on the
	following heuristics:

		* Watch Namespaces is not empty
			1. Installation Namespace (if it is watched)
			2. gloo-system (if it is watched)
			3. First item in the watch namespaces list
		* Watch Namespaces is empty
			1. Installation Namespace
*/
func getWriteNamespaceForGlooInstance(installNamespace string, watchedNamespaces []string) string {
	if len(watchedNamespaces) > 0 {
		if stringutils.ContainsString(installNamespace, watchedNamespaces) {
			return installNamespace
		} else if stringutils.ContainsString(defaults.GlooSystem, watchedNamespaces) {
			return defaults.GlooSystem
		}
		return watchedNamespaces[0]
	}
	return installNamespace
}

/*
	GetProxyForGlooInstance returns the proxy which should be used by gloo-fed for routing in failover and
	other related scenarios. Select based on the following heuristics:

		1. Proxy found with default gateway name "gateway-proxy" in install namespace
		2. First proxy found in install namespace
		3. First proxy in list

*/
func getProxyForGlooInstance(installNamespace string, proxies []*types.GlooInstanceSpec_Proxy) *types.GlooInstanceSpec_Proxy {
	if len(proxies) == 0 {
		return nil
	}

	// heuristics
	var glooSystemProxyIdx, inNamespaceProxyIdx = -1, -1
	for idx, proxy := range proxies {
		if proxy.GetNamespace() == installNamespace {
			if proxy.GetName() == "gateway-proxy" {
				// First heuristic satisfied, can break now
				glooSystemProxyIdx = idx
				break
			}
			inNamespaceProxyIdx = idx
		}
	}

	switch {
	case glooSystemProxyIdx != -1:
		return proxies[glooSystemProxyIdx]
	case inNamespaceProxyIdx != -1:
		return proxies[inNamespaceProxyIdx]
	default:
		return proxies[0]
	}
}

// GetSettings returns the first settings object in the given namespace and cluster.
// In all known real-world cases, this should be the "default"-named settings object.
func getSettings(set sk_sets.ResourceSet[*gloov1.Settings], namespace, cluster string) *gloov1.Settings {
	for _, settingsIter := range set.List() {
		settings := settingsIter
		if settings.Namespace == namespace && settings.ClusterName == cluster {
			return settings
		}
	}
	return nil
}

// Convert from the generic structs to the Gloo Fed types

func convertProxies(proxies []*proxy.Proxy) []*types.GlooInstanceSpec_Proxy {
	convertedProxies := make([]*types.GlooInstanceSpec_Proxy, 0, len(proxies))
	for _, p := range proxies {
		convertedProxies = append(convertedProxies, &types.GlooInstanceSpec_Proxy{
			Replicas:                      p.Replicas,
			AvailableReplicas:             p.AvailableReplicas,
			ReadyReplicas:                 p.ReadyReplicas,
			WasmEnabled:                   p.WasmEnabled,
			ReadConfigMulticlusterEnabled: p.ReadConfigMulticlusterEnabled,
			Version:                       p.Version,
			Name:                          p.Name,
			Namespace:                     p.Namespace,
			WorkloadControllerType:        convertWorkloadControllerType(p.WorkloadControllerType),
			Zones:                         p.Zones,
			IngressEndpoints:              convertIngressEndpoints(p.IngressEndpoints),
		})
	}
	return convertedProxies
}

func convertWorkloadControllerType(workloadControllerType proxy.WorkloadController) types.GlooInstanceSpec_Proxy_WorkloadController {
	switch workloadControllerType {
	case proxy.DEPLOYMENT:
		return types.GlooInstanceSpec_Proxy_DEPLOYMENT
	case proxy.DAEMON_SET:
		return types.GlooInstanceSpec_Proxy_DAEMON_SET
	default:
		return types.GlooInstanceSpec_Proxy_UNDEFINED
	}
}

func convertIngressEndpoints(ingressEndpoints []*locality2.IngressEndpoint) []*types.GlooInstanceSpec_Proxy_IngressEndpoint {
	convertedEndpoints := make([]*types.GlooInstanceSpec_Proxy_IngressEndpoint, 0, len(ingressEndpoints))
	for _, endpoint := range ingressEndpoints {
		convertedEndpoints = append(convertedEndpoints, &types.GlooInstanceSpec_Proxy_IngressEndpoint{
			Address:     endpoint.Address,
			Ports:       convertPorts(endpoint.Ports),
			ServiceName: endpoint.ServiceName,
		})
	}
	return convertedEndpoints
}

func convertPorts(ports []*locality2.Port) []*types.GlooInstanceSpec_Proxy_IngressEndpoint_Port {
	convertedPorts := make([]*types.GlooInstanceSpec_Proxy_IngressEndpoint_Port, 0, len(ports))
	for _, port := range ports {
		convertedPorts = append(convertedPorts, &types.GlooInstanceSpec_Proxy_IngressEndpoint_Port{
			Port: port.Port,
			Name: port.Name,
		})
	}
	return convertedPorts
}

func convertSummary(summary *checker2.Summary) *types.GlooInstanceSpec_Check_Summary {
	return &types.GlooInstanceSpec_Check_Summary{
		Total:    summary.Total,
		Errors:   convertResourceReports(summary.Errors),
		Warnings: convertResourceReports(summary.Warnings),
	}
}

func convertResourceReports(resourceReports []*checker2.ResourceReport) []*types.GlooInstanceSpec_Check_Summary_ResourceReport {
	if len(resourceReports) == 0 {
		return nil
	}
	convertedResourceReports := make([]*types.GlooInstanceSpec_Check_Summary_ResourceReport, 0, len(resourceReports))
	for _, resourceReport := range resourceReports {
		convertedResourceReports = append(convertedResourceReports, &types.GlooInstanceSpec_Check_Summary_ResourceReport{
			Ref:     resourceReport.Ref,
			Message: resourceReport.Message,
		})
	}
	return convertedResourceReports
}
