package upstreams

import (
	"github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
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
