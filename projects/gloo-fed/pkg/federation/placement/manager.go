package placement

import (
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
)

// ResourceWithPlacement defines the minimum API that resources must implement to have their
// PlacementStatuses maintained by the placement Manager
type ResourceWithPlacement interface {
	// GetPlacementStatus returns the PlacementStatus for a resource
	//
	// Deprecated: Prefer GetNamespacedPlacementStatuses
	GetPlacementStatus() *v1.PlacementStatus

	// SetPlacementStatus assigns the PlacementStatus for a resource
	//
	// Deprecated: Prefer SetNamespacedPlacementStatuses
	SetPlacementStatus(*v1.PlacementStatus)

	// GetNamespacedPlacementStatuses returns the set of PlacementStatuses for a resource,
	// keyed by the name of the namespace of the operator that updated the placement status
	GetNamespacedPlacementStatuses() map[string]*v1.PlacementStatus

	// SetNamespacedPlacementStatuses assigns the set of PlacementStatuses for a resource,
	// keyed by the name of the namespace of the operator that updated the placement status
	SetNamespacedPlacementStatuses(map[string]*v1.PlacementStatus)
}

var _ Manager = new(managerImpl)

// Manager is responsible for maintaining the PlacementStatus for a Federated resource
type Manager interface {
	StatusBuilderFactory
	SetPlacementStatus(resource ResourceWithPlacement, placementStatus *v1.PlacementStatus)
	GetPlacementStatus(resource ResourceWithPlacement) *v1.PlacementStatus
}

type managerImpl struct {
	podName      string
	podNamespace string
}

func NewManager(podNamespace, podName string) Manager {
	return &managerImpl{
		podName:      podName,
		podNamespace: podNamespace,
	}
}

func (m managerImpl) GetBuilder() StatusBuilder {
	return NewStatusBuilder(m.podName)
}

func (m managerImpl) SetPlacementStatus(resource ResourceWithPlacement, placementStatus *v1.PlacementStatus) {
	statuses := resource.GetNamespacedPlacementStatuses()
	if statuses == nil {
		resource.SetNamespacedPlacementStatuses(
			map[string]*v1.PlacementStatus{
				m.podNamespace: placementStatus,
			})
		return
	}
	statuses[m.podNamespace] = placementStatus
}

func (m managerImpl) GetPlacementStatus(resource ResourceWithPlacement) *v1.PlacementStatus {
	statuses := resource.GetNamespacedPlacementStatuses()
	if statuses == nil {
		return nil
	}
	return statuses[m.podNamespace]
}
