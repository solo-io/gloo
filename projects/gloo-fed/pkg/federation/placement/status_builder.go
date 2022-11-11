package placement

import (
	fed_core_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	// https://go.dev/doc/effective_go#blank_import
	// Important to implicitly import this so that it gets imported during code-gen.
	// This code is only imported in generated code which is deleted during generation.
	_ "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
)

var _ StatusBuilder = new(statusBuilder)

type statusBuilder struct {
	podName string
	status  *fed_core_v1.PlacementStatus
}

func NewStatusBuilder(podName string) StatusBuilder {
	return &statusBuilder{
		podName: podName,
	}
}

func (b *statusBuilder) init() {
	if b.status == nil {
		b.status = &fed_core_v1.PlacementStatus{}
	}
}

func (b *statusBuilder) AddDestination(cluster, namespace string, namespaceStatus fed_core_v1.PlacementStatus_Namespace) StatusBuilder {
	b.init()

	if b.status.Clusters == nil {
		b.status.Clusters = make(map[string]*fed_core_v1.PlacementStatus_Cluster)
	}

	if _, ok := b.status.Clusters[cluster]; !ok {
		b.status.Clusters[cluster] = &fed_core_v1.PlacementStatus_Cluster{
			Namespaces: make(map[string]*fed_core_v1.PlacementStatus_Namespace),
		}
	}

	b.status.Clusters[cluster].Namespaces[namespace] = &namespaceStatus

	return b
}

func (b *statusBuilder) AddDestinations(clusters, namespaces []string, namespaceStatus fed_core_v1.PlacementStatus_Namespace) StatusBuilder {
	b.init()

	for _, cluster := range clusters {
		for _, namespace := range namespaces {
			b.AddDestination(cluster, namespace, namespaceStatus)
		}
	}

	return b
}

func (b *statusBuilder) UpdateUnprocessed(status *fed_core_v1.PlacementStatus, reason string, state fed_core_v1.PlacementStatus_State) StatusBuilder {
	if status == nil {
		status = &fed_core_v1.PlacementStatus{}
	}
	b.status = status
	b.status.State = state
	b.status.Message = reason
	return b
}

func (b *statusBuilder) Build(generation int64) *fed_core_v1.PlacementStatus {
	b.init()

	// Add resource and pod specific metadata
	b.status.ObservedGeneration = generation
	b.status.WrittenBy = b.podName

	// Assume success unless a namespace is not placed
	b.status.State = fed_core_v1.PlacementStatus_PLACED

	for _, clusterStatus := range b.status.Clusters {
		for _, namespaceStatus := range clusterStatus.GetNamespaces() {
			if namespaceStatus.State != fed_core_v1.PlacementStatus_PLACED {
				b.status.State = fed_core_v1.PlacementStatus_FAILED
				b.status.Message = FailedToPlaceResource
				return b.status
			}
		}
	}

	return b.status
}

func (b *statusBuilder) Eject(generation int64) *fed_core_v1.PlacementStatus {
	b.init()

	// Add resource and pod specific metadata
	b.status.ObservedGeneration = generation
	b.status.WrittenBy = b.podName

	return b.status
}
