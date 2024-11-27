package upstreams

import (
	"strings"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// isRealUpstream returns true if the given upstream name refers to a "real" upstream that is
// stored in etcd (as opposed a "fake"/in-memory upstream)
func isRealUpstream(upstreamName string) bool {
	return !(kubernetes.IsFakeKubeUpstream(upstreamName) || consul.IsConsulUpstream(upstreamName))
}

// DestinationToUpstreamRef converts a route destination to an upstream ref. The upstream ref may
// refer to either a real or in-memory upstream, depending on the destination type.
func DestinationToUpstreamRef(dest *v1.Destination) (*core.ResourceRef, error) {
	var ref *core.ResourceRef
	switch d := dest.GetDestinationType().(type) {
	case *v1.Destination_Upstream:
		ref = d.Upstream
	case *v1.Destination_Kube:
		ref = kubernetes.DestinationToUpstreamRef(d.Kube)
	case *v1.Destination_Consul:
		ref = consul.DestinationToUpstreamRef(d.Consul)
	default:
		return nil, eris.Errorf("no destination type specified")
	}
	return ref, nil
}

// UpstreamToClusterName converts an Upstream to a cluster name to be used in envoy.
func UpstreamToClusterName(upstream *v1.Upstream) string {
	// when kube gateway is enabled, we use a new cluster name format that is more easily parseable
	if envutils.IsEnvTruthy(constants.GlooGatewayEnableK8sGwControllerEnv) {
		// currently only kube-type upstreams are handled, and the rest will fall back to the old cluster format
		switch upstreamType := upstream.GetUpstreamType().(type) {
		case *v1.Upstream_Kube:
			return kubernetes.UpstreamToClusterName(
				upstream.GetMetadata().GetName(),
				upstream.GetMetadata().GetNamespace(),
				upstreamType.Kube,
			)

			// TODO: other upstream types can be handled here later
		}
	}

	// fall back to "legacy" format
	return utils.ResourceRefToKey(upstream.GetMetadata().Ref())
}

// ClusterToUpstreamRef converts an envoy cluster name back to an upstream ref.
// (this is currently only used in the tunneling plugin and the old UI)
// This does the inverse of UpstreamToClusterName
func ClusterToUpstreamRef(cluster string) (*core.ResourceRef, error) {
	// if kube gateway is enabled and it's a kube-type cluster name, use the new parsing logic
	if envutils.IsEnvTruthy(constants.GlooGatewayEnableK8sGwControllerEnv) &&
		kubernetes.IsKubeCluster(cluster) {
		return kubernetes.ClusterToUpstreamRef(cluster)
	}

	// TODO: if we add support for more cluster name formats based on upstream type
	// (in UpstreamToClusterName), add the reverse conversion logic here

	// otherwise fall back to old logic:
	// legacy and non-kube cluster names consist of `upstreamName_upstreamNamespace`
	split := strings.Split(cluster, "_")
	if len(split) > 2 || len(split) < 1 {
		return nil, eris.Errorf("unable to convert cluster %s back to upstream ref", cluster)
	}

	ref := &core.ResourceRef{
		Name: split[0],
	}

	if len(split) == 2 {
		ref.Namespace = split[1]
	}
	return ref, nil
}
