package serviceentry

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	networkingclient "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
)

// All krt collection types must provide a key through Resourcenamer or being a controllers.Object
// Optionally, they can be a LabelSelectorer to allow fetching using selected labels (i.e. pod selected by ServiceEntry).

var _ krt.ResourceNamer = endpoint{}

type endpoints []endpoint

func (eps endpoints) Unwrap() v1.EndpointList {
	out := make(v1.EndpointList, 0, len(eps))
	for _, ep := range eps {
		out = append(out, ep.Endpoint)
	}
	return out
}

type endpoint struct {
	*v1.Endpoint
}

func (ep endpoint) ResourceName() string {
	return ep.Metadata.GetName() + "/" + ep.Metadata.GetNamespace()
}

var _ krt.ResourceNamer = upstream{}

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

type serviceEntryUpstreams struct {
	*networkingclient.ServiceEntry
	// list of upstreams for each of the service ports
	upstreams map[uint32][]upstream
	// mapping from service port to target port
	targetPorts map[uint32]uint32
}

func (s serviceEntryUpstreams) GetLabelSelector() map[string]string {
	return s.Spec.GetWorkloadSelector().GetLabels()
}
