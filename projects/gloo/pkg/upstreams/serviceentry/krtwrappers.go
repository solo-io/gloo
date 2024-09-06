package serviceentry

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
)

// All krt collection types must provide a key through Resourcenamer or being a controllers.Object
// Optionally, they can be a LabelSelectorer to allow fetching using selected labels (i.e. pod selected by ServiceEntry).

// endpoints is a convenience type for convering lists of our krt
// wrapped `endpoint` back into Gloo's `v1.EndpointList`
type endpoints []glooEndpoint

func (eps endpoints) Unwrap() v1.EndpointList {
	out := make(v1.EndpointList, 0, len(eps))
	for _, ep := range eps {
		out = append(out, ep.Endpoint)
	}
	return out
}

var _ krt.ResourceNamer = glooEndpoint{}

// glooEndpoint provides a krt keying function for Gloo's `v1.Endpoint`
type glooEndpoint struct {
	*v1.Endpoint
}

func (ep glooEndpoint) ResourceName() string {
	return ep.Metadata.GetName() + "/" + ep.Metadata.GetNamespace()
}

var _ krt.ResourceNamer = upstream{}

// upstream provides a keying function for Gloo's `v1.Upstream`
type upstream struct {
	*v1.Upstream
}

func (us upstream) ResourceName() string {
	return us.Metadata.GetName() + "/" + us.Metadata.GetNamespace()
}

var (
	_ krt.LabelSelectorer = serviceEntryUpstreams{}
	_ controllers.Object  = serviceEntryUpstreams{}
)

// serviceEntryUpstreams contains back-references to Upstreams that link to
// this ServiceEntry by hostname. We also implement LabelSelectorer to support
// krt selection queries.
type serviceEntryUpstreams struct {
	*networkingclient.ServiceEntry
	// list of upstreams for each of the service ports
	upstreams map[uint32][]upstream
	// mapping from service port to target port
	targetPorts map[uint32]uint32
	// mapping of svc port to name
	portNames map[uint32]string
}

func (s serviceEntryUpstreams) GetLabelSelector() map[string]string {
	return s.Spec.GetWorkloadSelector().GetLabels()
}
