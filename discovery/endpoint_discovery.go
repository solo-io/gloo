package discovery

type Clusters map[string]Cluster

type Cluster struct {
	Hostname  string
	Endpoints []Endpoint
}

type Endpoint struct {
	Address string
	Port    uint32
}

type Discovery interface {
	// cluster ref is the discovery-specific identifier for the cluster
	// in kubernetes, this would be the namespace+service name
	DiscoveryFor(hostnames []string)

	// secrets are pushed here whenever they are read
	Clusters() <-chan Clusters

	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
