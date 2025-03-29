package proxy_syncer

import (
	"github.com/solo-io/gloo/pkg/utils/envutils"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"istio.io/api/label"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/ptr"
)

type istioNetworkSingleton struct {
	krt.Singleton[string]
}

// istioNetworkSingleton as set by the label on the system namespace.
// This is the cluster-wide default for objects that don't declare their own network
// via the label on Pod or WorkloadEntry (or the field network on WorkloadEntry).
func newIstioNetworkSingleton(client kube.Client) istioNetworkSingleton {
	istioNamespace := envutils.GetOrDefault("ISTIO_NAMESPACE", "istio-system", false)
	filter := kclient.Filter{ObjectFilter: client.ObjectFilter(), FieldSelector: "metadata.name=" + istioNamespace}
	istioNamespaceCol := krt.NewInformerFiltered[*corev1.Namespace](client, filter, krt.WithName("IstioNamespace"))
	singleton := krt.NewSingleton(func(ctx krt.HandlerContext) *string {
		network := ""
		namespace := ptr.Flatten(krt.FetchOne(ctx, istioNamespaceCol, krt.FilterObjectName(types.NamespacedName{Name: istioNamespace})))
		if namespace != nil {
			if networkLabel, ok := namespace.GetLabels()[label.TopologyNetwork.Name]; ok {
				network = networkLabel
			}
		}
		return &network
	})
	return istioNetworkSingleton{Singleton: singleton}
}

func (i *istioNetworkSingleton) Fetch(ctx krt.HandlerContext) string {
	network := krt.FetchOne(ctx, i.AsCollection())
	return ptr.OrEmpty(network)
}

func istioNetworkFromLabels(labels map[string]string, defaultNetwork string) string {
	if net, ok := labels[label.TopologyCluster.Name]; ok {
		return net
	}
	return defaultNetwork
}
