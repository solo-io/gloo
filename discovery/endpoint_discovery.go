package discovery

import "github.com/solo-io/glue/pkg/api/types/v1"

// cluster of endpoints for an upstream
type Clusters map[string][]Endpoint

type Endpoint struct {
	Address string
	Port    uint32
}

type Discovery interface {
	// cluster ref is the discovery-specific identifier for the cluster
	// in kubernetes, this would be the namespace+service name

	/*
		upstream:
		- name: my_upstream_v1
		  type: kubernetes
		  spec:
			kubernetes_name: service
		    labels:
		    - version: v1
	*/

	TrackUpstreams(upstreams []v1.Upstream)

	// secrets are pushed here whenever they are read
	Clusters() <-chan Clusters

	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
