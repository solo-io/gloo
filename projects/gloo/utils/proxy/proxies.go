package proxy

import (
	"context"
	"sort"

	apps_v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1/sets"
	core_v1_sets "github.com/solo-io/external-apis/pkg/api/k8s/core/v1/sets"
	gateway_defaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/skv2/contrib/pkg/sets"
	"github.com/solo-io/solo-projects/projects/gloo/utils/images"
	"github.com/solo-io/solo-projects/projects/gloo/utils/locality"
	"go.uber.org/zap"
	k8s_core_types "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type WorkloadController int32

const (
	DEPLOYMENT WorkloadController = 1
	DAEMON_SET WorkloadController = 2
)

func (w WorkloadController) String() string {
	switch w {
	case DEPLOYMENT:
		return "DEPLOYMENT"
	case DAEMON_SET:
		return "DAEMON_SET"
	default:
		return "UNDEFINED"
	}
}

type Proxy struct {
	Replicas                      int32
	AvailableReplicas             int32
	ReadyReplicas                 int32
	WasmEnabled                   bool
	ReadConfigMulticlusterEnabled bool
	Version                       string
	Name                          string
	Namespace                     string
	WorkloadControllerType        WorkloadController
	Zones                         []string
	IngressEndpoints              []*locality.IngressEndpoint
}

func (p *Proxy) GetName() string {
	return p.Name
}
func (p *Proxy) GetNamespace() string {
	return p.Namespace
}

// Get the proxies for a Gloo instance. To bypass the cluster check (e.g. for single-cluster use), pass in "" for the cluster.
func GetGlooInstanceProxies(
	ctx context.Context,
	cluster string,
	namespace string,
	deployments apps_v1_sets.DeploymentSet,
	daemonSets apps_v1_sets.DaemonSetSet,
	services core_v1_sets.ServiceSet,
	localityFinder locality.LocalityFinder,
	ipFinder locality.ExternalIpFinder,
) []*Proxy {

	var proxies []*Proxy
	for _, deployment := range deployments.List() {
		deployment := deployment
		if cluster != "" && deployment.GetClusterName() != cluster {
			continue
		}
		if deployment.GetNamespace() != namespace {
			continue
		}
		if tag, wasmEnabled, isProxy := images.GetGatewayProxyImage(&deployment.Spec.Template); isProxy {
			zones, err := localityFinder.ZonesForDeployment(ctx, deployment)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("failed to get zones for deployment", zap.Error(err), zap.Any("deployment", deployment))
			}

			proxy := &Proxy{
				Replicas:                      deployment.Status.Replicas,
				AvailableReplicas:             deployment.Status.AvailableReplicas,
				ReadyReplicas:                 deployment.Status.ReadyReplicas,
				WasmEnabled:                   wasmEnabled,
				ReadConfigMulticlusterEnabled: findProxyConfigDumpService(services, deployment.GetName(), deployment.GetNamespace()),
				Version:                       tag,
				Name:                          deployment.GetName(),
				Namespace:                     deployment.GetNamespace(),
				WorkloadControllerType:        DEPLOYMENT,
				Zones:                         zones,
			}

			ingressEndpoints, err := ipFinder.GetExternalIps(
				ctx,
				filterServices(services, deployment, deployment.Spec.Template.Labels),
			)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("error while determining ingress endpoints", zap.Error(err))
			}
			proxy.IngressEndpoints = ingressEndpoints

			proxies = append(proxies, proxy)
		}
	}

	for _, daemonSet := range daemonSets.List() {
		daemonSet := daemonSet
		if daemonSet.GetNamespace() != namespace {
			continue
		}
		if tag, wasmEnabled, isProxy := images.GetGatewayProxyImage(&daemonSet.Spec.Template); isProxy {
			zones, err := localityFinder.ZonesForDaemonSet(ctx, daemonSet)
			if err != nil {
				contextutils.LoggerFrom(ctx).Debugw("failed to get zones for daemonSet", zap.Error(err), zap.Any("daemonSet", daemonSet))
			}

			proxy := &Proxy{
				Replicas:               daemonSet.Status.DesiredNumberScheduled,
				AvailableReplicas:      daemonSet.Status.NumberAvailable,
				ReadyReplicas:          daemonSet.Status.NumberReady,
				WasmEnabled:            wasmEnabled,
				Version:                tag,
				Name:                   daemonSet.GetName(),
				Namespace:              daemonSet.GetNamespace(),
				WorkloadControllerType: DAEMON_SET,
				Zones:                  zones,
			}

			ingressEndpoints, err := ipFinder.GetExternalIps(
				ctx,
				filterServices(services, daemonSet, daemonSet.Spec.Template.Labels),
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

func findProxyConfigDumpService(services core_v1_sets.ServiceSet, name, namespace string) bool {
	for _, svc := range services.List() {
		if svc.Name == gateway_defaults.GatewayProxyConfigDumpServiceName(name) && svc.Namespace == namespace {
			return true
		}
	}
	return false
}

func filterServices(services core_v1_sets.ServiceSet, obj meta_v1.Object, matchLabels map[string]string) []*k8s_core_types.Service {
	var result []*k8s_core_types.Service
	for _, svc := range services.List() {
		if svc.GetClusterName() == obj.GetClusterName() &&
			svc.GetNamespace() == obj.GetNamespace() &&
			labels.SelectorFromSet(svc.Spec.Selector).Matches(labels.Set(matchLabels)) {
			result = append(result, svc)
		}
	}
	return result
}
