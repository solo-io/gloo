package apiserverutils

import (
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ToMetadata(objectMeta metav1.ObjectMeta) *rpc_edge_v1.ObjectMeta {
	return &rpc_edge_v1.ObjectMeta{
		Name:            objectMeta.Name,
		Namespace:       objectMeta.Namespace,
		Labels:          objectMeta.Labels,
		Annotations:     objectMeta.Annotations,
		ResourceVersion: objectMeta.ResourceVersion,
		Uid:             string(objectMeta.UID),
	}
}

func RefToObjectMeta(ref v1.ClusterObjectRef) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      ref.GetName(),
		Namespace: ref.GetNamespace(),
	}
}

func ToObjectRef(name, namespace string) v1.ObjectRef {
	return v1.ObjectRef{
		Name:      name,
		Namespace: namespace,
	}
}

func ToClusterObjectRef(name, namespace, cluster string) v1.ClusterObjectRef {
	return v1.ClusterObjectRef{
		Name:        name,
		Namespace:   namespace,
		ClusterName: cluster,
	}
}
