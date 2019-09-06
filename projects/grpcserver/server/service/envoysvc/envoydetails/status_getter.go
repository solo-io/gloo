package envoydetails

import (
	"context"
	"fmt"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"go.uber.org/zap"
	kubev1 "k8s.io/api/core/v1"
)

//go:generate mockgen -destination mocks/status_getter_mock.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails ProxyStatusGetter

type ProxyStatusGetter interface {
	GetProxyStatus(ctx context.Context, pod kubev1.Pod) *v1.Status
}

type proxyStatusGetter struct {
	proxyClient gloov1.ProxyClient
}

var _ ProxyStatusGetter = &proxyStatusGetter{}

func NewProxyStatusGetter(proxyClient gloov1.ProxyClient) ProxyStatusGetter {
	return &proxyStatusGetter{proxyClient: proxyClient}
}

func (g *proxyStatusGetter) GetProxyStatus(ctx context.Context, pod kubev1.Pod) *v1.Status {
	if pod.Status.Phase != kubev1.PodRunning {
		return &v1.Status{
			Code:    v1.Status_ERROR,
			Message: GatewayProxyPodIsNotRunning(pod.Namespace, pod.Name, string(pod.Status.Phase)),
		}
	}

	// TODO joekelley is there a stronger association among gateways, proxies, and gateway-proxies?
	proxy, err := g.proxyClient.Read(pod.Namespace, pod.Labels[GatewayProxyIdLabel], clients.ReadOpts{Ctx: ctx})
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
	case core.Status_Accepted:
		return &v1.Status{Code: v1.Status_OK}
	case core.Status_Pending:
		return &v1.Status{
			Code:    v1.Status_WARNING,
			Message: ProxyResourcePending(proxy.GetMetadata().Namespace, proxy.GetMetadata().Name),
		}
	default:
		return &v1.Status{
			Code:    v1.Status_ERROR,
			Message: ProxyResourceRejected(proxy.GetMetadata().Namespace, proxy.GetMetadata().Name, proxy.GetStatus().Reason),
		}
	}
}
