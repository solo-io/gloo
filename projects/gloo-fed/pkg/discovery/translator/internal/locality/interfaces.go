package locality

import (
	"context"

	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	appsv1 "k8s.io/api/apps/v1"
	k8s_core_types "k8s.io/api/core/v1"
)

/*
	LocalityFinder aids in the discovery of the locality associated with kubernetes workloads.
	It does this by checking well-known node labels https://kubernetes.io/docs/reference/kubernetes-api/labels-annotations-taints/.
	If no locality can be found on the nodes, than an empty string will be returned for region, and an empty list for zones.
*/
type LocalityFinder interface {
	// GetRegion attempts to get the region for a given kubernetes cluster. If no region can be found, will return ""
	GetRegion(ctx context.Context) (string, error)
	// ZonesForDeployment gets the zones in which a deployment's replicas reside in, if none are found will return nil
	ZonesForDeployment(ctx context.Context, deployment *appsv1.Deployment) ([]string, error)
	// ZonesForDaemonSet gets the zones in which a daemonset's pods reside in, if none are found will return nil
	ZonesForDaemonSet(ctx context.Context, set *appsv1.DaemonSet) ([]string, error)
}

/*
	ExternalIpFinder takes a list of k8s services and the cluster they are located on, it can then find
	the string address to communicate with the endpoints via an external network
*/
type ExternalIpFinder interface {
	GetExternalIps(
		ctx context.Context,
		svcs []*k8s_core_types.Service,
	) ([]*types.GlooInstanceSpec_Proxy_IngressEndpoint, error)
}
