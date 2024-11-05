package kubernetes

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes/serviceconverter"
	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"

	sanitizer "github.com/solo-io/k8s-utils/kubeutils"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
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
	// KubeSourceResourceLabel indicates the original name of the resource that
	// the synthetic resource is based on.
	KubeNameLabel = "~internal.solo.io/kubernetes-name"
	// KubeSourceResourceLabel indicates the original namespace of the resource
	// that the synthetic resource is based on.
	KubeNamespaceLabel = "~internal.solo.io/kubernetes-namespace"
	// KubeSourceResourceLabel indicates the service port when applicable.
	KubeServicePortLabel = "~internal.solo.io/kubernetes-service-port"
)

// ClusterNameForKube builds the cluster name based on _internal_ labels.
// All of the kind, name, namespace and port must be provided.
func ClusterNameForKube(us *v1.Upstream) (string, bool) {
	labels := us.GetMetadata().GetLabels()
	kind, kok := labels[KubeSourceResourceLabel]
	name, nok := labels[KubeNameLabel]
	ns, nsok := labels[KubeNamespaceLabel]
	port, pok := labels[KubeServicePortLabel]
	if !(kok && nok && nsok && pok) {
		return "", false
	}
	return fmt.Sprintf("%s_%s_%s_%s", kind, name, ns, port), true
}

type UpstreamConverter interface {
	UpstreamsForService(ctx context.Context, svc *corev1.Service) v1.UpstreamList
}

func DefaultUpstreamConverter() *KubeUpstreamConverter {
	kuc := new(KubeUpstreamConverter)
	kuc.serviceConverters = serviceconverter.DefaultServiceConverters
	return kuc
}

type KubeUpstreamConverter struct {
	serviceConverters []serviceconverter.ServiceConverter
}

func (uc *KubeUpstreamConverter) UpstreamsForService(ctx context.Context, svc *corev1.Service) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, port := range svc.Spec.Ports {
		upstreams = append(upstreams, uc.CreateUpstream(ctx, svc, port))
	}
	return upstreams
}

func (uc *KubeUpstreamConverter) CreateUpstream(ctx context.Context, svc *corev1.Service, port corev1.ServicePort) *v1.Upstream {
	meta := svc.ObjectMeta
	coremeta := kubeutils.FromKubeMeta(meta, false)
	coremeta.ResourceVersion = ""
	coremeta.Name = UpstreamName(meta.Namespace, meta.Name, port.Port)
	labels := coremeta.GetLabels()
	coremeta.Labels = map[string]string{
		// preserve parts of the source service in a structured way
		// so we don't rely on string parsing to recover these
		// this is more extensible than relying on casting Spec to Upstream_Kube
		KubeSourceResourceLabel: "kube-svc",
		KubeNameLabel:           meta.Name,
		KubeNamespaceLabel:      meta.Namespace,
		KubeServicePortLabel:    strconv.Itoa(int(port.Port)),
	}

	us := &v1.Upstream{
		Metadata: coremeta,
		UpstreamType: &v1.Upstream_Kube{
			Kube: &kubeplugin.UpstreamSpec{
				ServiceName:      meta.Name,
				ServiceNamespace: meta.Namespace,
				ServicePort:      uint32(port.Port),
				Selector:         svc.Spec.Selector,
			},
		},
		DiscoveryMetadata: &v1.DiscoveryMetadata{
			Labels: labels,
		},
	}

	for _, sc := range uc.serviceConverters {
		if err := sc.ConvertService(ctx, svc, port, us); err != nil {
			contextutils.LoggerFrom(ctx).Errorf("error: failed to process service options with err %v", err)
		}
	}

	return us
}

func UpstreamName(serviceNamespace, serviceName string, servicePort int32) string {
	return sanitizer.SanitizeNameV2(fmt.Sprintf("%s-%s-%v", serviceNamespace, serviceName, servicePort))
}

func skip(svc *corev1.Service, opts discovery.Opts) bool {
	// ilackarms: allow user to override the skip with an annotation
	// force discovery for a service with no selector
	if svc.ObjectMeta.Annotations[discoveryAnnotationKey] == discoveryAnnotationTrue {
		return false
	}
	// note: ilackarms: IgnoredServices is not set anywhere
	for _, name := range opts.KubeOpts.IgnoredServices {
		if svc.Name == name {
			return true
		}
	}
	return false
}

func (p *plugin) UpdateUpstream(original, desired *v1.Upstream) (bool, error) {
	return UpdateUpstream(original, desired)
}

func UpdateUpstream(original, desired *v1.Upstream) (didChange bool, err error) {
	originalSpec, ok := original.GetUpstreamType().(*v1.Upstream_Kube)
	if !ok {
		return false, errors.Errorf("internal error: expected *v1.Upstream_Kube, got %v", reflect.TypeOf(original.GetUpstreamType()).Name())
	}
	desiredSpec, ok := desired.GetUpstreamType().(*v1.Upstream_Kube)
	if !ok {
		return false, errors.Errorf("internal error: expected *v1.Upstream_Kube, got %v", reflect.TypeOf(original.GetUpstreamType()).Name())
	}
	// copy service spec, we don't want to overwrite that
	desiredSpec.Kube.ServiceSpec = originalSpec.Kube.GetServiceSpec()
	// copy labels; user may have written them over. cannot be auto-discovered
	desiredSpec.Kube.Selector = originalSpec.Kube.GetSelector()

	utils.UpdateUpstream(original, desired)

	return !upstreamsEqual(original, desired), nil
}

// we want to know if the upstreams are equal apart from their Status and Metadata
func upstreamsEqual(original, desired *v1.Upstream) bool {
	copyOriginal := *original
	copyDesired := *desired

	copyOriginal.Metadata = copyDesired.GetMetadata()
	copyOriginal.SetNamespacedStatuses(copyDesired.GetNamespacedStatuses())

	return copyOriginal.Equal(copyDesired)
}
