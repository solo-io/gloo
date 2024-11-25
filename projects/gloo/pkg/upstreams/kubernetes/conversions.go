package kubernetes

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	corev1 "k8s.io/api/core/v1"
)

// these labels are used to propagate internal data
// on synthetic Gloo resources generated from other Kubernetes
// resources (generally Service).
// The `~` is an invalid character that prevents these labels from ending up
// on actual Kubernetes resources.
const (
	// KubeSourceResourceLabel indicates the kind of resource that the synthetic
	// resource is based on.
	KubeSourceResourceLabel = "~internal.solo.io/kubernetes-source-resource"
	// KubeNameLabel indicates the original name of the resource that
	// the synthetic resource is based on.
	KubeNameLabel = "~internal.solo.io/kubernetes-name"
	// KubeNamespaceLabel indicates the original namespace of the resource
	// that the synthetic resource is based on.
	KubeNamespaceLabel = "~internal.solo.io/kubernetes-namespace"
	// KubeServicePortLabel indicates the service port when applicable.
	KubeServicePortLabel = "~internal.solo.io/kubernetes-service-port"
)

func IsKubeUpstream(upstreamName string) bool {
	return strings.HasPrefix(upstreamName, UpstreamNamePrefix)
}

func DestinationToUpstreamRef(svcDest *v1.KubernetesServiceDestination) *core.ResourceRef {
	return &core.ResourceRef{
		Namespace: svcDest.GetRef().GetNamespace(),
		Name:      fakeUpstreamName(svcDest.GetRef().GetName(), svcDest.GetRef().GetNamespace(), int32(svcDest.GetPort())),
	}
}

func fakeUpstreamName(serviceName, serviceNamespace string, port int32) string {
	regularServiceName := kubeplugin.UpstreamName(serviceNamespace, serviceName, port)
	return UpstreamNamePrefix + regularServiceName
}

// KubeServicesToUpstreams converts a list of k8s Services to a list of in-memory Upstreams.
// Public because it's needed in the translator test
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
// called by KubeServicesToUpstreams and
func ServiceToUpstream(ctx context.Context, svc *corev1.Service, port corev1.ServicePort) *gloov1.Upstream {
	us := kubeplugin.DefaultUpstreamConverter().CreateUpstream(ctx, svc, port)

	us.GetMetadata().Name = fakeUpstreamName(svc.GetName(), svc.GetNamespace(), port.Port)
	us.GetMetadata().Namespace = svc.GetNamespace()
	us.GetMetadata().ResourceVersion = ""

	additionalLabels := map[string]string{
		// preserve parts of the source service in a structured way
		// so we don't rely on string parsing to recover these
		// this is more extensible than relying on casting Spec to Upstream_Kube
		KubeSourceResourceLabel: UpstreamNamePrefixNoSeparator,
		KubeNameLabel:           svc.GetName(),
		KubeNamespaceLabel:      svc.GetNamespace(),
		KubeServicePortLabel:    strconv.Itoa(int(port.Port)),
	}
	if us.GetMetadata().GetLabels() == nil {
		us.GetMetadata().Labels = map[string]string{}
	}
	for k, v := range additionalLabels {
		us.GetMetadata().Labels[k] = v
	}

	fmt.Printf("xxxxxxx ServiceToUpstream: %s\n", us.GetMetadata().GetName())
	return us
}
