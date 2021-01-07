package placement

import fed_core_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_interfaces.go

// StatusBuilder facilitates the construction of PlacementStatuses.
type StatusBuilder interface {
	// AddDestination adds a namespaced placement status for the given cluster and namespace.
	AddDestination(cluster, namespace string, namespaceStatus fed_core_v1.PlacementStatus_Namespace) StatusBuilder
	// AddDestinations adds a namespaced placement status for all the given cluster and namespace.
	AddDestinations(clusters, namespaces []string, namespaceStatus fed_core_v1.PlacementStatus_Namespace) StatusBuilder
	// UpdateUnprocessed updates an existing PlacementStatus to indicate a change in status which indicates that
	// the specific cluster/namespace placement of a resource has not been updated.
	UpdateUnprocessed(status *fed_core_v1.PlacementStatus, reason string, state fed_core_v1.PlacementStatus_State) StatusBuilder
	// Build computes a top-level resource status from per-destination statuses and returns a complete PlacementStatus.
	Build(generation int64) *fed_core_v1.PlacementStatus
	// Eject returns the placement status stored in the builder without recalculating the top-level resource status.
	Eject(generation int64) *fed_core_v1.PlacementStatus
}

// StatusBuilderFactory is a factory for StatusBuilders.
type StatusBuilderFactory interface {
	// GetBuilder returns a new instance of StatusBuilder.
	GetBuilder() StatusBuilder
}
