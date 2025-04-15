package usage

type Stats struct {
	Stats []Stat `json:"stats"`
}

type Stat struct {
	Name  string `json:"name"`
	Value interface{}    `json:"value"`
}

type SocketAddress struct {
	Address   string `json:"address"`
	PortValue int    `json:"port_value"`
}

type Address struct {
	SocketAddress SocketAddress `json:"socket_address,omitempty"`
	Pipe          struct {
		Path string `json:"path"`
	} `json:"pipe,omitempty"`
}

type EndpointStat struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
}

type HealthStatus struct {
	EdsHealthStatus string `json:"eds_health_status"`
}

type HostStatus struct {
	Address      Address       `json:"address"`
	Stats        []EndpointStat `json:"stats"`
	HealthStatus HealthStatus `json:"health_status"`
	Weight       int          `json:"weight"`
	Hostname     string       `json:"hostname,omitempty"`
	Locality     Locality     `json:"locality"`
}

type Locality struct {
	Region  string `json:"region"`
	Zone    string `json:"zone"`
	SubZone string `json:"sub_zone"`
}

type Threshold struct {
	MaxConnections     int    `json:"max_connections"`
	MaxPendingRequests int    `json:"max_pending_requests"`
	MaxRequests        int    `json:"max_requests"`
	MaxRetries         int    `json:"max_retries"`
	Priority           string `json:"priority,omitempty"`
}

type CircuitBreakers struct {
	Thresholds []Threshold `json:"thresholds"`
}

type ClusterStatus struct {
	Name              string          `json:"name"`
	AddedViaAPI       bool            `json:"added_via_api"`
	HostStatuses      []HostStatus    `json:"host_statuses"`
	CircuitBreakers   CircuitBreakers `json:"circuit_breakers"`
	ObservabilityName string          `json:"observability_name"`
	EdsServiceName    string          `json:"eds_service_name,omitempty"`
}

type Clusters struct {
	ClusterStatuses []ClusterStatus `json:"cluster_statuses"`
}
