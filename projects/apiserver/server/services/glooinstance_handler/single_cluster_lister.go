package glooinstance_handler

import (
	"context"
	"sort"

	"github.com/rotisserie/eris"
	apps_v1 "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1"
	apps_v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/sets"
	core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	core_v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/ezkube"
	enterprise_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	ratelimit_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	enterprise_gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler"
	gateway_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/gateway.solo.io/v1/handler"
	gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/gloo.solo.io/v1/handler"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"github.com/solo-io/solo-projects/projects/gloo/utils/checker"
	"github.com/solo-io/solo-projects/projects/gloo/utils/images"
	"github.com/solo-io/solo-projects/projects/gloo/utils/locality"
	"go.uber.org/zap"
	k8s_apps_types "k8s.io/api/apps/v1"
	k8s_core_types "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Single-cluster Gloo Edge doesn't have a notion of clusters like Gloo Fed (which requires cluster registration
	// with user-provided cluster names), so we hardcode an arbitrary cluster name here for display purposes
	ClusterName = "cluster"

	// Note: currently we only show the GlooInstance in the current install namespace. In the future, if we need to
	// support listing all GlooInstances on a cluster, make this a configurable helm flag
	SingleInstanceOnly = true
)

//go:generate mockgen -source ./single_cluster_lister.go -destination mocks/single_cluster_lister.go

func NewSingleClusterGlooInstanceLister(
	coreClientset core_v1.Clientset,
	appsClientset apps_v1.Clientset,
	gatewayClientset gateway_v1.Clientset,
	glooClientset gloo_v1.Clientset,
	enterpriseGlooClientset enterprise_gloo_v1.Clientset,
	ratelimitClientset ratelimit_v1alpha1.Clientset,
) SingleClusterGlooInstanceLister {
	return &singleClusterGlooInstanceLister{
		coreClientset:           coreClientset,
		appsClientset:           appsClientset,
		gatewayClientset:        gatewayClientset,
		glooClientset:           glooClientset,
		enterpriseGlooClientset: enterpriseGlooClientset,
		ratelimitClientset:      ratelimitClientset,
	}
}

type SingleClusterGlooInstanceLister interface {
	ListGlooInstances(ctx context.Context) ([]*rpc_edge_v1.GlooInstance, error)
	GetGlooInstance(ctx context.Context, glooInstanceRef *skv2_v1.ObjectRef) (*rpc_edge_v1.GlooInstance, error)
}

type singleClusterGlooInstanceLister struct {
	coreClientset           core_v1.Clientset
	appsClientset           apps_v1.Clientset
	gatewayClientset        gateway_v1.Clientset
	glooClientset           gloo_v1.Clientset
	enterpriseGlooClientset enterprise_gloo_v1.Clientset
	ratelimitClientset      ratelimit_v1alpha1.Clientset
}

// Gets a list of Gloo instances on this cluster.
// If SingleInstanceOnly is true, only returns the Gloo instance in the same install namespace as
// the currently running apiserver instance (rather than all Gloo instances on the cluster).
func (l *singleClusterGlooInstanceLister) ListGlooInstances(ctx context.Context) ([]*rpc_edge_v1.GlooInstance, error) {
	var instances []*rpc_edge_v1.GlooInstance
	installNamespace := apiserverutils.GetInstallNamespace()

	// if we only want the Gloo instance in the current install namespace, we can get it directly
	if SingleInstanceOnly {
		instance, err := l.GetGlooInstance(ctx, &skv2_v1.ObjectRef{
			Name:      images.ControlPlaneDeploymentName,
			Namespace: installNamespace,
		})
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
		return instances, nil
	}

	deployments, err := l.appsClientset.Deployments().ListDeployment(ctx)
	if err != nil {
		return nil, err
	}
	daemonSets, err := l.appsClientset.DaemonSets().ListDaemonSet(ctx)
	if err != nil {
		return nil, err
	}
	services, err := l.coreClientset.Services().ListService(ctx)
	if err != nil {
		return nil, err
	}
	pods, err := l.coreClientset.Pods().ListPod(ctx)
	if err != nil {
		return nil, err
	}

	// go through all the deployments to find the Gloo control planes
	for _, deployment := range deployments.Items {
		deployment := deployment

		// get info about the control plane
		tag, isEnterprise, isControlPlane := images.GetControlPlaneImage(&deployment)
		if !isControlPlane {
			continue
		}

		instance, err := l.buildGlooInstanceFromDeployment(ctx, &deployment, tag, isEnterprise, deployments, daemonSets, services, pods)
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	sortGlooInstances(instances)
	return instances, nil
}

func (l *singleClusterGlooInstanceLister) GetGlooInstance(ctx context.Context, glooInstanceRef *skv2_v1.ObjectRef) (*rpc_edge_v1.GlooInstance, error) {
	// get the deployment with the same name/namespace as the glooInstanceRef
	deployment, err := l.appsClientset.Deployments().GetDeployment(ctx, ezkube.MakeClientObjectKey(glooInstanceRef))
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw("Failed to get deployment for GlooInstance", zap.Error(err), zap.Any("glooInstanceRef", glooInstanceRef))
		return nil, err
	}
	tag, isEnterprise, isControlPlane := images.GetControlPlaneImage(deployment)
	if !isControlPlane {
		return nil, eris.Errorf("Ref does not refer to a GlooInstance: %v", glooInstanceRef)
	}

	deployments, err := l.appsClientset.Deployments().ListDeployment(ctx)
	if err != nil {
		return nil, err
	}
	daemonSets, err := l.appsClientset.DaemonSets().ListDaemonSet(ctx)
	if err != nil {
		return nil, err
	}
	services, err := l.coreClientset.Services().ListService(ctx)
	if err != nil {
		return nil, err
	}
	pods, err := l.coreClientset.Pods().ListPod(ctx)
	if err != nil {
		return nil, err
	}
	glooInstance, err := l.buildGlooInstanceFromDeployment(ctx, deployment, tag, isEnterprise, deployments, daemonSets, services, pods)
	if err != nil {
		return nil, err
	}
	return glooInstance, nil
}

func (l *singleClusterGlooInstanceLister) buildGlooInstanceFromDeployment(
	ctx context.Context,
	deployment *k8s_apps_types.Deployment,
	tag string,
	isEnterprise bool,
	deployments *k8s_apps_types.DeploymentList,
	daemonSets *k8s_apps_types.DaemonSetList,
	services *k8s_core_types.ServiceList,
	pods *k8s_core_types.PodList) (*rpc_edge_v1.GlooInstance, error) {

	// get locality and ingress info
	podClient := l.coreClientset.Pods()
	nodeClient := l.coreClientset.Nodes()
	ipFinder := locality.NewExternalIpFinder(ClusterName, podClient, nodeClient)
	localityFinder := locality.NewLocalityFinder(nodeClient, podClient)
	region, err := localityFinder.GetRegion(ctx)
	if err != nil {
		// log the error but keep going; we just won't have region info
		contextutils.LoggerFrom(ctx).Warnw("Failed to get region of control plane", zap.Error(err), zap.Any("deployment", deployment))
	}

	// get the watched namespaces from this Gloo instance's settings
	settingsClient := l.glooClientset.Settings()
	settings, err := settingsClient.GetSettings(ctx, client.ObjectKey{
		Namespace: deployment.GetNamespace(),
		Name:      defaults.SettingsName,
	})
	if err != nil {
		return nil, err
	}
	watchedNamespaces := settings.Spec.GetWatchNamespaces()

	instance := &rpc_edge_v1.GlooInstance{
		Metadata: &rpc_edge_v1.ObjectMeta{
			// use the gloo deployment name and namespace
			Name:      deployment.GetName(),
			Namespace: deployment.GetNamespace(),
		},
		Spec: &rpc_edge_v1.GlooInstance_GlooInstanceSpec{
			Cluster:      ClusterName,
			IsEnterprise: isEnterprise,
			ControlPlane: &rpc_edge_v1.GlooInstance_GlooInstanceSpec_ControlPlane{
				Version:           tag,
				Namespace:         deployment.GetNamespace(),
				WatchedNamespaces: watchedNamespaces,
			},
			Region: region,
			Check:  &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check{},
		},
	}

	// proxy info
	instance.Spec.Proxies = l.getProxies(ctx, deployment.GetNamespace(), deployments, daemonSets, services, localityFinder, ipFinder)

	// pod and deployment summaries
	podSet := core_v1_sets.NewPodSetFromList(pods)
	deploymentSet := apps_v1_sets.NewDeploymentSetFromList(deployments)
	instance.Spec.Check.Pods = convertSummary(checker.GetPodsSummary(ctx, podSet, instance.Spec.ControlPlane.GetNamespace(), ""))
	instance.Spec.Check.Deployments = convertSummary(checker.GetDeploymentsSummary(ctx, deploymentSet, instance.Spec.ControlPlane.GetNamespace(), ""))

	// gloo resource summaries
	instance.Spec.Check.Settings = gloo_resource_handler.GetSettingsSummary(ctx, settingsClient, watchedNamespaces)
	instance.Spec.Check.Upstreams = gloo_resource_handler.GetUpstreamSummary(ctx, l.glooClientset.Upstreams(), watchedNamespaces)
	instance.Spec.Check.UpstreamGroups = gloo_resource_handler.GetUpstreamGroupSummary(ctx, l.glooClientset.UpstreamGroups(), watchedNamespaces)
	instance.Spec.Check.Proxies = gloo_resource_handler.GetProxySummary(ctx, l.glooClientset.Proxies(), watchedNamespaces)
	instance.Spec.Check.AuthConfigs = enterprise_gloo_resource_handler.GetAuthConfigSummary(ctx, l.enterpriseGlooClientset.AuthConfigs(), watchedNamespaces)

	// gateway resource summaries
	instance.Spec.Check.Gateways = gateway_resource_handler.GetGatewaySummary(ctx, l.gatewayClientset.Gateways(), watchedNamespaces)
	instance.Spec.Check.VirtualServices = gateway_resource_handler.GetVirtualServiceSummary(ctx, l.gatewayClientset.VirtualServices(), watchedNamespaces)
	instance.Spec.Check.RouteTables = gateway_resource_handler.GetRouteTableSummary(ctx, l.gatewayClientset.RouteTables(), watchedNamespaces)

	return instance, nil
}

// This follows the same logic as the Gloo Fed translator projects/gloo-fed/pkg/discovery/translator/translator.go
// but for a single cluster
func (l *singleClusterGlooInstanceLister) getProxies(
	ctx context.Context,
	namespace string,
	deployments *k8s_apps_types.DeploymentList,
	daemonSets *k8s_apps_types.DaemonSetList,
	services *k8s_core_types.ServiceList,
	localityFinder locality.LocalityFinder,
	ipFinder locality.ExternalIpFinder,
) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy {

	var proxies []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy
	for _, deployment := range deployments.Items {
		deployment := deployment
		if deployment.GetNamespace() != namespace {
			continue
		}
		if tag, wasmEnabled, isProxy := images.GetGatewayProxyImage(&deployment.Spec.Template); isProxy {
			zones, err := localityFinder.ZonesForDeployment(ctx, &deployment)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("failed to get zones for deployment", zap.Error(err), zap.Any("deployment", deployment))
			}

			proxy := &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
				Replicas:                      deployment.Status.Replicas,
				AvailableReplicas:             deployment.Status.AvailableReplicas,
				ReadyReplicas:                 deployment.Status.ReadyReplicas,
				WasmEnabled:                   wasmEnabled,
				ReadConfigMulticlusterEnabled: findProxyConfigDumpService(services, deployment.GetName(), deployment.GetNamespace()),
				Version:                       tag,
				Name:                          deployment.GetName(),
				Namespace:                     deployment.GetNamespace(),
				WorkloadControllerType:        rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DEPLOYMENT,
				Zones:                         zones,
			}

			ingressEndpoints, err := ipFinder.GetExternalIps(
				ctx,
				filterServices(services, &deployment, deployment.Spec.Template.Labels),
			)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("error while determining ingress endpoints", zap.Error(err))
			}
			proxy.IngressEndpoints = convertIngressEndpoints(ingressEndpoints)

			proxies = append(proxies, proxy)
		}
	}

	for _, daemonSet := range daemonSets.Items {
		daemonSet := daemonSet
		if daemonSet.GetNamespace() != namespace {
			continue
		}
		if tag, wasmEnabled, isProxy := images.GetGatewayProxyImage(&daemonSet.Spec.Template); isProxy {
			zones, err := localityFinder.ZonesForDaemonSet(ctx, &daemonSet)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("failed to get zones for daemonSet", zap.Error(err), zap.Any("daemonSet", daemonSet))
			}

			proxy := &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy{
				Replicas:               daemonSet.Status.DesiredNumberScheduled,
				AvailableReplicas:      daemonSet.Status.NumberAvailable,
				ReadyReplicas:          daemonSet.Status.NumberReady,
				WasmEnabled:            wasmEnabled,
				Version:                tag,
				Name:                   daemonSet.GetName(),
				Namespace:              daemonSet.GetNamespace(),
				WorkloadControllerType: rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_DAEMON_SET,
				Zones:                  zones,
			}

			ingressEndpoints, err := ipFinder.GetExternalIps(
				ctx,
				filterServices(services, &daemonSet, daemonSet.Spec.Template.Labels),
			)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("error while determining ingress endpoints", zap.Error(err))
			}
			proxy.IngressEndpoints = convertIngressEndpoints(ingressEndpoints)

			proxies = append(proxies, proxy)
		}
	}

	// Sort all proxies for idempotence, this is especially important when selecting admin info.
	sort.SliceStable(proxies, func(i, j int) bool {
		return proxies[i].WorkloadControllerType.String()+sets.Key(proxies[i]) < proxies[j].WorkloadControllerType.String()+sets.Key(proxies[j])
	})

	return proxies

}

func findProxyConfigDumpService(services *k8s_core_types.ServiceList, name, namespace string) bool {
	for _, svc := range services.Items {
		if svc.Name == name+"-config-dump-service" && svc.Namespace == namespace {
			return true
		}
	}
	return false
}

func filterServices(services *k8s_core_types.ServiceList, obj meta_v1.Object, matchLabels map[string]string) []*k8s_core_types.Service {
	var result []*k8s_core_types.Service
	for _, svc := range services.Items {
		if svc.GetClusterName() == obj.GetClusterName() &&
			svc.GetNamespace() == obj.GetNamespace() &&
			labels.SelectorFromSet(svc.Spec.Selector).Matches(labels.Set(matchLabels)) {
			result = append(result, &svc)
		}
	}
	return result
}

// Convert from the generic structs to the GlooEE apiserver types
func convertIngressEndpoints(ingressEndpoints []*locality.IngressEndpoint) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint {
	convertedEndpoints := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint, 0, len(ingressEndpoints))
	for _, endpoint := range ingressEndpoints {
		convertedEndpoints = append(convertedEndpoints, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint{
			Address:     endpoint.Address,
			Ports:       convertPorts(endpoint.Ports),
			ServiceName: endpoint.ServiceName,
		})
	}
	return convertedEndpoints
}

func convertPorts(ports []*locality.Port) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port {
	convertedPorts := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port, 0, len(ports))
	for _, port := range ports {
		convertedPorts = append(convertedPorts, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Proxy_IngressEndpoint_Port{
			Port: port.Port,
			Name: port.Name,
		})
	}
	return convertedPorts
}

func convertSummary(summary *checker.Summary) *rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary {
	return &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary{
		Total:    summary.Total,
		Errors:   convertResourceReports(summary.Errors),
		Warnings: convertResourceReports(summary.Warnings),
	}
}

func convertResourceReports(resourceReports []*checker.ResourceReport) []*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport {
	if len(resourceReports) == 0 {
		return nil
	}
	convertedResourceReports := make([]*rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport, 0, len(resourceReports))
	for _, resourceReport := range resourceReports {
		convertedResourceReports = append(convertedResourceReports, &rpc_edge_v1.GlooInstance_GlooInstanceSpec_Check_Summary_ResourceReport{
			Ref:     resourceReport.Ref,
			Message: resourceReport.Message,
		})
	}
	return convertedResourceReports
}
