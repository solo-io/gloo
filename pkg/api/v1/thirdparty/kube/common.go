package kube

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func fromKubeMeta(meta v1.ObjectMeta) core.Metadata {
	return core.Metadata{
		Name:            meta.Name,
		Namespace:       meta.Namespace,
		ResourceVersion: meta.ResourceVersion,
		Labels:          meta.Labels,
		Annotations:     meta.Annotations,
	}
}

func toKubeMeta(meta core.Metadata) v1.ObjectMeta {
	return v1.ObjectMeta{
		Name:            meta.Name,
		Namespace:       clients.DefaultNamespaceIfEmpty(meta.Namespace),
		ResourceVersion: meta.ResourceVersion,
		Labels:          meta.Labels,
		Annotations:     meta.Annotations,
	}
}
