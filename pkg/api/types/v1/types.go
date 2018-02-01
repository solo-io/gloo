package v1

type Config struct {
	Routes       []Route
	Upstreams    []Upstream
	VirtualHosts []VirtualHost
}

type Route struct {
	Matcher     Matcher                `json:"matcher"`
	Destination Destination            `json:"destination"`
	Weight      int                    `json:"weight"`
	Plugins     map[string]interface{} `json:"plugins"`
}

type Destination struct {
	// A valid destination can only contain one of:
	// FunctionDestination
	// UpstreamDestination
	FunctionDestionation *FunctionDestination `json:"function_destination,omitemtpy"`
	UpstreamDestination  *UpstreamDestination `json:"upstream_destination,omitemtpy"`
}

type Matcher struct {
	Path        Path              `json:"path"`
	Headers     map[string]string `json:"headers"`
	Verbs       []string          `json:"verbs"`
	VirtualHost string            `json:"virtual_host"`
}

type Path struct {
	// a valid path can only contain one of:
	// prefix
	// regex
	// exact
	Prefix string `json:"prefix,omitemtpy"`
	Regex  string `json:"regex,omitemtpy"`
	Exact  string `json:"exact,omitemtpy"`
}

type FunctionDestination struct {
	UpstreamName string `json:"upstream_name"`
	FunctionName string `json:"function_name"`
}

type UpstreamDestination struct {
	UpstreamName  string `json:"upstream_name"`
	RewritePrefix string `json:"rewrite_prefix"`
}

type UpstreamType string

type Upstream struct {
	Name      string                 `json:"name"`
	Type      UpstreamType           `json:"type"`
	Spec      map[string]interface{} `json:"spec"`
	Functions []Function             `json:"functions"`
}

type Function struct {
	Name string                 `json:"name"`
	Spec map[string]interface{} `json:"spec"`
}

type VirtualHost struct {
	Domains   []string  `json:"domains"`
	SSLConfig SSLConfig `json:"ssl_config,omitemtpy"`
}

type SSLConfig struct {
	CACertPath string `json:"ca_cert_path"`
}
