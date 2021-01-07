package translator

import (
	"context"
	"sort"

	k8s_core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"

	k8s_core_sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stringutils"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	v1sets "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/sets"
	enterprise_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/enterprise.gloo.solo.io/v1/check"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	gateway_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/gateway.solo.io/v1/check"
	gloo_check "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/gloo.solo.io/v1/check"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/images"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator/internal/checker"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator/internal/locality"
	"go.uber.org/zap"
	k8s_core_types "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
func (t *translator) FromSnapshot(ctx context.Context, snapshot input.Snapshot) []*v1.GlooInstance {
	var instances []*v1.GlooInstance
	for _, deployment := range snapshot.Deployments().List() {
		image, isEnterprise, ok := images.GetControlPlaneImage(deployment)
		if !ok {
			// There should be one Gloo Instance per cluster's control plane deployment, so skip all other deployments.
			continue
		}

		cluster := deployment.GetClusterName()
		client, err := t.mcClientSet.Cluster(cluster)
		if err != nil {
			contextutils.LoggerFrom(ctx).Errorf("couldn't get client for cluster %v", cluster)
			return nil
		}
		ipFinder := locality.NewExternalIpFinder(cluster, client.Pods(), client.Nodes())
		localityFinder := locality.NewLocalityFinder(client.Nodes(), client.Pods())
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
					Version:           images.GetVersion(image),
					Namespace:         deployment.Namespace,
					WatchedNamespaces: watchedNamespaces,
				},
				Region: region,
				Check:  &types.GlooInstanceSpec_Check{},
			},
		}

		// Proxy and admin info
		instance.Spec.Proxies = getProxiesFromSnapshot(ctx, cluster, localityFinder, ipFinder, snapshot)
		instance.Spec.Admin = getAdminInfo(instance.Spec.ControlPlane.GetNamespace(), watchedNamespaces, instance.Spec.Proxies)

		// Check -- workloads
		instance.Spec.Check.Pods = checker.GetPodsSummary(ctx, snapshot.Pods(), instance.Spec.ControlPlane.GetNamespace(), cluster)
		instance.Spec.Check.Deployments = checker.GetDeploymentsSummary(ctx, snapshot.Deployments(), instance.Spec.ControlPlane.GetNamespace(), cluster)

		// Check -- resources
		// Gloo
		instance.Spec.Check.Settings = gloo_check.GetSettingsSummary(ctx, snapshot.Settings(), watchedNamespaces, cluster)
		instance.Spec.Check.Upstreams = gloo_check.GetUpstreamSummary(ctx, snapshot.Upstreams(), watchedNamespaces, cluster)
		instance.Spec.Check.UpstreamGroups = gloo_check.GetUpstreamGroupSummary(ctx, snapshot.UpstreamGroups(), watchedNamespaces, cluster)
		instance.Spec.Check.Proxies = gloo_check.GetProxySummary(ctx, snapshot.Proxies(), watchedNamespaces, cluster)
		instance.Spec.Check.AuthConfigs = enterprise_check.GetAuthConfigSummary(ctx, snapshot.AuthConfigs(), watchedNamespaces, cluster)

		// Gateway
		instance.Spec.Check.Gateways = gateway_check.GetGatewaySummary(ctx, snapshot.Gateways(), watchedNamespaces, cluster)
		instance.Spec.Check.VirtualServices = gateway_check.GetVirtualServiceSummary(ctx, snapshot.VirtualServices(), watchedNamespaces, cluster)
		instance.Spec.Check.RouteTables = gateway_check.GetRouteTableSummary(ctx, snapshot.RouteTables(), watchedNamespaces, cluster)

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

func getProxiesFromSnapshot(
	ctx context.Context,
	cluster string,
	localityFinder locality.LocalityFinder,
	ipFinder locality.ExternalIpFinder,
	snapshot input.Snapshot,
) []*types.GlooInstanceSpec_Proxy {

	var proxies []*types.GlooInstanceSpec_Proxy
	for _, deploymentIter := range snapshot.Deployments().List() {
		deployment := deploymentIter
		if deployment.GetClusterName() != cluster {
			continue
		}
		if image, wasmEnabled, ok := images.GetGatewayProxyImage(&deployment.Spec.Template); ok {
			zones, err := localityFinder.ZonesForDeployment(ctx, deployment)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("failed to get zones for deployment", zap.Error(err), zap.Any("deployment", deployment))
			}

			proxy := &types.GlooInstanceSpec_Proxy{
				Replicas:               deployment.Status.Replicas,
				AvailableReplicas:      deployment.Status.AvailableReplicas,
				ReadyReplicas:          deployment.Status.ReadyReplicas,
				WasmEnabled:            wasmEnabled,
				Version:                images.GetVersion(image),
				Name:                   deployment.GetName(),
				Namespace:              deployment.GetNamespace(),
				WorkloadControllerType: types.GlooInstanceSpec_Proxy_DEPLOYMENT,
				Zones:                  zones,
			}

			ingressEndpoints, err := ipFinder.GetExternalIps(
				ctx,
				filterServices(snapshot.Services(), deployment, deployment.Spec.Template.Labels),
			)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("error while determining ingress endpoints", zap.Error(err))
			}
			proxy.IngressEndpoints = ingressEndpoints

			proxies = append(proxies, proxy)
		}
	}

	for _, daemonset := range snapshot.DaemonSets().List() {
		if image, wasmEnabled, ok := images.GetGatewayProxyImage(&daemonset.Spec.Template); ok {
			zones, err := localityFinder.ZonesForDaemonSet(ctx, daemonset)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("failed to get zones for daemonset", zap.Error(err), zap.Any("daemonset", daemonset))
			}

			proxy := &types.GlooInstanceSpec_Proxy{
				Replicas:               daemonset.Status.DesiredNumberScheduled,
				AvailableReplicas:      daemonset.Status.NumberAvailable,
				ReadyReplicas:          daemonset.Status.NumberReady,
				WasmEnabled:            wasmEnabled,
				Version:                images.GetVersion(image),
				Name:                   daemonset.GetName(),
				Namespace:              daemonset.GetNamespace(),
				WorkloadControllerType: types.GlooInstanceSpec_Proxy_DAEMON_SET,
				Zones:                  zones,
			}

			ingressEndpoints, err := ipFinder.GetExternalIps(
				ctx,
				filterServices(snapshot.Services(), daemonset, daemonset.Spec.Template.Labels),
			)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("error while determining ingress endpoints", zap.Error(err))
			}
			proxy.IngressEndpoints = ingressEndpoints

			proxies = append(proxies, proxy)
		}
	}

	// Sort all proxies for idempotence, this is especially important when selecting admin info.
	sort.SliceStable(proxies, func(i, j int) bool {
		return proxies[i].WorkloadControllerType.String()+sets.Key(proxies[i]) < proxies[j].WorkloadControllerType.String()+sets.Key(proxies[j])
	})

	return proxies

}

func filterServices(svcs k8s_core_sets.ServiceSet, obj metav1.Object, matchLabels map[string]string) []*k8s_core_types.Service {
	var result []*k8s_core_types.Service
	for _, svc := range svcs.List() {
		if svc.GetClusterName() == obj.GetClusterName() &&
			svc.GetNamespace() == obj.GetNamespace() &&
			labels.SelectorFromSet(svc.Spec.Selector).Matches(labels.Set(matchLabels)) {
			result = append(result, svc)
		}
	}
	return result
}

// GetSettings returns the first settings object in the given namespace and cluster.
// In all known real-world cases, this should be the "default"-named settings object.
func getSettings(set v1sets.SettingsSet, namespace, cluster string) *gloov1.Settings {
	for _, settingsIter := range set.List() {
		settings := settingsIter
		if settings.Namespace == namespace && settings.ClusterName == cluster {
			return settings
		}
	}
	return nil
}
