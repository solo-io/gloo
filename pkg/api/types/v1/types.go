package v1

type ConfigObject interface {
	IsConfigObject()
}

func (c *Route) IsConfigObject()       {}
func (c *Upstream) IsConfigObject()    {}
func (c *VirtualHost) IsConfigObject() {}
func (c *Function) IsConfigObject()    {}

type spec map[string]interface{}

type RoutePluginSpec spec
type UpstreamSpec spec
type FunctionSpec spec

type Config struct {
	Upstreams    []Upstream
	VirtualHosts []VirtualHost
}

type Route struct {
	Matcher       Matcher                    `json:"matcher"`
	Destination   Destination                `json:"destination"`
	RewritePrefix string                     `json:"rewrite_prefix"`
	Plugins       map[string]RoutePluginSpec `json:"plugins"`
}

type Destination struct {
	// A valid destination can only contain one of:
	// SingleDestination
	// Destinations
	SingleDestination
	Destinations []WeightedDestination
}

type WeightedDestination struct {
	SingleDestination
	Weight uint
}

type SingleDestination struct {
	// A valid destination can only contain one of:
	// FunctionDestination
	// UpstreamDestination
	FunctionDestination *FunctionDestination `json:"function_destination,omitemtpy"`
	UpstreamDestination *UpstreamDestination `json:"upstream_destination,omitemtpy"`
}

type Matcher struct {
	Path    Path              `json:"path"`
	Headers map[string]string `json:"headers"`
	Verbs   []string          `json:"verbs"`
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
	UpstreamName string `json:"upstream_name"` /// Move to function object?
	FunctionName string `json:"function_name"`
}

type UpstreamDestination struct {
	UpstreamName string `json:"upstream_name"`
}

type UpstreamType string

type Upstream struct {
	Name      string                 `json:"name"`
	Type      UpstreamType           `json:"type"`
	Spec      map[string]interface{} `json:"spec"`
	Functions []Function             `json:"functions"`
}

type Function struct {
	Name string `json:"name"`
	// upstream ref?
	Spec map[string]interface{} `json:"spec"`
}

type VirtualHost struct {
	Name      string   `json:"name"`
	Domains   []string `json:"domains"`
	Routes    []Route
	SSLConfig SSLConfig `json:"ssl_config,omitemtpy"`
	// ^ secret ref | or file
	// should route rules live here?
}

type SSLConfig struct {
	CACertPath string `json:"ca_cert_path"`
	SecretRef  string `json:"secret_ref"`
}
