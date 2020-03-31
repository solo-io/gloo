package envoydetails

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	clientscache "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"go.uber.org/zap"
	kubev1 "k8s.io/api/core/v1"
)

//go:generate mockgen -destination mocks/status_getter_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails EnvoyStatusGetter

type EnvoyStatusGetter interface {
	GetEnvoyStatus(ctx context.Context, pod kubev1.Pod) *v1.Status
}

type envoyStatusGetter struct {
	clientCache clientscache.ClientCache
}

var _ EnvoyStatusGetter = &envoyStatusGetter{}

func NewEnvoyStatusGetter(clientCache clientscache.ClientCache) EnvoyStatusGetter {
	return &envoyStatusGetter{clientCache: clientCache}
}

func (g *envoyStatusGetter) GetEnvoyStatus(ctx context.Context, pod kubev1.Pod) *v1.Status {
	if pod.Status.Phase != kubev1.PodRunning {
		return &v1.Status{
			Code:    v1.Status_ERROR,
			Message: GatewayProxyPodIsNotRunning(pod.Namespace, pod.Name, string(pod.Status.Phase)),
		}
	}

	// TODO joekelley is there a stronger association among gateways, proxies, and gateway-proxies?
	proxy, err := g.clientCache.GetProxyClient().Read(pod.Namespace, pod.Labels[GatewayProxyIdLabel], clients.ReadOpts{Ctx: ctx})
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw(
			fmt.Sprintf("Failed to load proxy resource for gateway proxy pod %v.%v", pod.Namespace, pod.Name),
			zap.String("namespace", pod.Namespace),
			zap.String("name", pod.Name))
	}

	if proxy == nil {
		return &v1.Status{
			Code:    v1.Status_ERROR,
			Message: ProxyResourceNotFound(getName(pod)),
		}
	}

	switch proxy.GetStatus().State {
	case core.Status_Pending:
		return &v1.Status{
			Code:    v1.Status_WARNING,
			Message: ProxyResourcePending(proxy.GetMetadata().Namespace, proxy.GetMetadata().Name),
		}
	case core.Status_Rejected:
		return &v1.Status{
			Code:    v1.Status_ERROR,
			Message: ProxyResourceRejected(proxy.GetMetadata().Namespace, proxy.GetMetadata().Name, proxy.GetStatus().Reason),
		}
	}

	// By default, check gateways in the same namespace as the proxy controller.
	// If ReadGatewaysFromAllNamespaces=true, then check gateways in all namespaces.
	gatewayNs := pod.Namespace
	settings, err := g.clientCache.GetSettingsClient().Read(pod.Namespace, defaults.SettingsName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnw(
			fmt.Sprintf("Failed to read settings for 'readGatewaysFromAllNamespaces' value, so only looking"+
				"for gateways in the %v namespace", pod.Namespace))
	} else {
		if settings.GetGateway().GetReadGatewaysFromAllNamespaces() {
			gatewayNs = kubev1.NamespaceAll // empty string means look in all namespaces
		}
	}
	// If gateways are rejected, then we need to propagate that error to the overall Envoy status
	// because no changes will make it into the proxy.
	gateways, err := g.clientCache.GetGatewayClient().List(gatewayNs, clients.ListOpts{Ctx: ctx})
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorw(
			fmt.Sprintf("Failed to load gateway resources in namespace: %v", pod.Namespace),
			zap.String("namespace", pod.Namespace),
			zap.String("name", pod.Name))
	}
	if len(gateways) == 0 {
		return &v1.Status{
			Code:    v1.Status_ERROR,
			Message: GatewayResourcesNotFound(gatewayNs),
		}
	}
	for _, gw := range gateways {
		if !(gw.GetStatus().State == core.Status_Accepted || gw.GetStatus().State == core.Status_Pending) {
			return &v1.Status{
				Code:    v1.Status_ERROR,
				Message: "Gateways are in a bad state",
			}
		}
	}
	return &v1.Status{Code: v1.Status_OK}
}
