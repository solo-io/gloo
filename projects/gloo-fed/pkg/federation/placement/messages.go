package placement

import "fmt"

var (
	FailedToCreateClientForCluster = func(cluster string) string {
		return fmt.Sprintf("Failed to create client for cluster %s. Check that it has been registered.", cluster)
	}

	FailedToUpsertResource = func(kind string) string {
		return fmt.Sprintf("Failed to write %s to the given destination. Check that the namespace exists and that Gloo Fed has permission to write to it.", kind)
	}

	FailedToUpsertResourceDueToConflict = func(kind string) string {
		return fmt.Sprintf("Failed to write %s to the given destination due to a version conflict. Retrying.", kind)
	}

	FailedToListResource = func(kind, cluster string) string {
		return fmt.Sprintf("Failed to list %s on cluster %s.", kind, cluster)
	}

	FailedToDeleteResource = func(kind string) string {
		return fmt.Sprintf("Failed to delete %s. A stale %s may still be present in the given destination.", kind, kind)
	}

	ClusterNotRegistered = func(cluster string) string {
		return fmt.Sprintf("Cluster %s is not registered", cluster)
	}

	ClusterEventTriggered = func(cluster string) string {
		return fmt.Sprintf("Cluster event triggered for %s, resyncing", cluster)
	}

	SpecTemplateMissing = "Spec template must not be nil"

	MetaTemplateMissing = "Metadata template must not be nil"

	FailedToPlaceResource = "Resource is not federated to all destinations"

	PlacementMissing = "Placement must not be nil"
)
