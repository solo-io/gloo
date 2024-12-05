package kubernetes

import (
	"context"
	"strings"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	corev1 "k8s.io/api/core/v1"
)

// FakeUpstreamNamePrefix is a prefix used to create/identify in-memory Upstreams for Kubernetes Services.
// It contains an invalid character so any accidental attempt to write to storage fails.
// Clusters created from these in-memory upstreams, as well as the cluster stats, will also contain this prefix.
// Note that `:` will be replaced by `_` in stats names.
const FakeUpstreamNamePrefix = "kube-svc:"

// KubeUpstreamStatsPrefix is a prefix added to stats to identify clusters created from "real" Upstreams of `kube` type.
// This prefix only appears in stats, and not in the cluster or Upstream names.
const KubeUpstreamStatsPrefix = "kube-upstream_"

// IsFakeKubeUpstream returns true if the given upstream is a "fake"/in-memory upstream representing
// a kubernetes service. An in-memory upstream is used when a routing destination is a Kubernetes
// Service (as opposed to an Upstream).
func IsFakeKubeUpstream(upstreamName string) bool {
	return strings.HasPrefix(upstreamName, FakeUpstreamNamePrefix)
}

// DestinationToUpstreamRef converts a k8s service Destination to an in-memory upstream ref.
// This is called by the generic upstream DestinationToUpstreamRef.
func DestinationToUpstreamRef(svcDest *gloov1.KubernetesServiceDestination) *core.ResourceRef {
	return &core.ResourceRef{
		Namespace: svcDest.GetRef().GetNamespace(),
		Name:      fakeUpstreamName(svcDest.GetRef().GetName(), svcDest.GetRef().GetNamespace(), int32(svcDest.GetPort())),
	}
}

// fakeUpstreamName generates a name for an in-memory upstream, based on the given service spec
func fakeUpstreamName(serviceName, serviceNamespace string, port int32) string {
	regularServiceName := kubeplugin.UpstreamName(serviceNamespace, serviceName, port)
	return FakeUpstreamNamePrefix + regularServiceName
}

// KubeServicesToUpstreams converts a list of k8s Services to a list of in-memory Upstreams.
// Public because it's needed in the translator test
// Also used by k8s upstream client (List/Watch)
func KubeServicesToUpstreams(ctx context.Context, services skkube.ServiceList) gloov1.UpstreamList {
	var result gloov1.UpstreamList
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
