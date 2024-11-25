package kubernetes

import (
	"context"
	"strings"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	corev1 "k8s.io/api/core/v1"
)

func IsKubeUpstream(upstreamName string) bool {
	return strings.HasPrefix(upstreamName, UpstreamNamePrefix)
}

// DestinationToUpstreamRef converts a Destination of type k8s service to an in-memory upstream ref
// This is called by the generic upstream DestinationToUpstreamRef which is used by many plugins
// to convert a route destination to an upstream ref.
func DestinationToUpstreamRef(svcDest *v1.KubernetesServiceDestination) *core.ResourceRef {
	return &core.ResourceRef{
		Namespace: svcDest.GetRef().GetNamespace(),
		Name:      fakeUpstreamName(svcDest.GetRef().GetName(), svcDest.GetRef().GetNamespace(), int32(svcDest.GetPort())),
	}
}

func fakeUpstreamName(serviceName, serviceNamespace string, port int32) string {
	regularServiceName := kubeplugin.UpstreamName(serviceNamespace, serviceName, port, false)
	return UpstreamNamePrefix + regularServiceName
}

// KubeServicesToUpstreams converts a list of k8s Services to a list of in-memory Upstreams.
// Public because it's needed in the translator test
// Also used by k8s upstream client (List/Watch)
func KubeServicesToUpstreams(ctx context.Context, services skkube.ServiceList) v1.UpstreamList {
	var result v1.UpstreamList
	for _, svc := range services {
		for _, port := range svc.Spec.Ports {
			kubeSvc := svc.Service.GetKubeService()
			result = append(result, ServiceToUpstream(ctx, &kubeSvc, port))
		}
	}
	return result
}

// ServiceToUpstream converts a k8s Service to an in-memory Upstream (with the kube-svc prefix).
// Called by KubeServicesToUpstreams (above) and kube gateway proxy syncer when initializing
// in-memory upstreams from k8s services.
func ServiceToUpstream(ctx context.Context, svc *corev1.Service, port corev1.ServicePort) *gloov1.Upstream {
	us := kubeplugin.DefaultUpstreamConverter().CreateUpstream(ctx, svc, port)

	us.GetMetadata().Name = fakeUpstreamName(svc.GetName(), svc.GetNamespace(), port.Port)
	us.GetMetadata().Namespace = svc.GetNamespace()
	us.GetMetadata().ResourceVersion = ""

	return us
}
