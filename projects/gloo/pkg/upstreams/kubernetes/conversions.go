package kubernetes

import (
	"context"
	"fmt"
	"strings"

	"github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	corev1 "k8s.io/api/core/v1"
)

// FakeUpstreamNamePrefix is a prefix used to create/identify in-memory Upstreams for Kubernetes Services.
// It contains an invalid character so any accidental attempt to write to storage fails.
// Clusters created from these in-memory upstreams will also contain this prefix.
const FakeUpstreamNamePrefix = "kube-svc:"

// KubeUpstreamClusterPrefix is a prefix added to clusters created from "real" Upstreams of `kube` type.
// This prefix only appears in the cluster name and not in the actual name of the Upstream from which
// the cluster was created (since it contains an invalid character).
const KubeUpstreamClusterPrefix = "kube-upstream:"

// IsFakeKubeUpstream returns true if the given upstream is a "fake"/in-memory upstream representing
// a kubernetes service. An in-memory upstream is used when a routing destination is a Kubernetes
// Service (as opposed to an Upstream).
func IsFakeKubeUpstream(upstreamName string) bool {
	return strings.HasPrefix(upstreamName, FakeUpstreamNamePrefix)
}

// IsKubeCluster returns true if the given envoy cluster was created from an Upstream of `kube` type
// (which may be either "real" or "fake").
func IsKubeCluster(cluster string) bool {
	return strings.HasPrefix(cluster, FakeUpstreamNamePrefix) || strings.HasPrefix(cluster, KubeUpstreamClusterPrefix)
}

// DestinationToUpstreamRef converts a k8s service Destination to an in-memory upstream ref.
// This is called by the generic upstream DestinationToUpstreamRef.
func DestinationToUpstreamRef(svcDest *gloov1.KubernetesServiceDestination) *core.ResourceRef {
	return &core.ResourceRef{
		Namespace: svcDest.GetRef().GetNamespace(),
		Name:      fakeUpstreamName(svcDest.GetRef().GetName(), svcDest.GetRef().GetNamespace(), int32(svcDest.GetPort())),
	}
}

// UpstreamToClusterName converts a kube upstream (which may be either "real" or "fake") to a cluster name.
// This is called by the generic upstream UpstreamToClusterName.
// This should only be called if kube gateway is enabled.
func UpstreamToClusterName(upstreamName string, upstreamNamespace string, kubeSpec *kubernetes.UpstreamSpec) string {
	// Add an identifying prefix if it's a "real" upstream (fake upstreams already have such a prefix).
	if !IsFakeKubeUpstream(upstreamName) {
		upstreamName = fmt.Sprintf("%s%s", KubeUpstreamClusterPrefix, upstreamName)
	}
	return fmt.Sprintf("%s_%s_%s_%s_%v",
		upstreamName,
		upstreamNamespace,
		kubeSpec.GetServiceNamespace(),
		kubeSpec.GetServiceName(),
		kubeSpec.GetServicePort(),
	)
}

// ClusterToUpstreamRef converts an envoy cluster name to a kube upstream ref (the upstream may be either "real" or "fake").
// This does the inverse of UpstreamToClusterName.
// This is called by the generic ClusterToUpstreamRef.
// This should only be called if kube gateway is enabled.
func ClusterToUpstreamRef(cluster string) (*core.ResourceRef, error) {
	if !IsKubeCluster(cluster) {
		return nil, eris.Errorf("cluster %s does not refer to a kubernetes service", cluster)
	}

	// split cluster name into component parts (see UpstreamToClusterName for format)
	parts := strings.Split(cluster, "_")
	if len(parts) != 5 {
		return nil, eris.Errorf("unable to convert cluster %s back to upstream ref", cluster)
	}
	upstreamName := parts[0]
	upstreamNamespace := parts[1]

	// if it's a "real" upstream, remove the special cluster prefix, as that's not part of the actual upstream name
	if strings.HasPrefix(cluster, KubeUpstreamClusterPrefix) {
		upstreamName = upstreamName[len(KubeUpstreamClusterPrefix):]
	}

	return &core.ResourceRef{
		Name:      upstreamName,
		Namespace: upstreamNamespace,
	}, nil
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
