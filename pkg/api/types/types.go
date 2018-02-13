package types

import "time"

// currently only top-level objects are storable
type StorableConfigObject interface {
	SetStorageRef(string)
	GetStorageRef() string
}

type storageRef struct {
	ref string
}

func (r *storageRef) SetStorageRef(ref string) {
	r.ref = ref
}

func (r *storageRef) GetStorageRef() string {
	if r.ref == "" {
		panic("storage ref not set")
	}
	return r.ref
}

type spec map[string]interface{}

type RoutePluginSpec spec
type UpstreamSpec spec
type FunctionSpec spec

type Config struct {
	Upstreams    []Upstream
	VirtualHosts []VirtualHost
}

type VirtualHost struct {
	storageRef

	Name      string   `json:"name"`
	Domains   []string `json:"domains"`
	Routes    []Route
	SSLConfig SSLConfig `json:"ssl_config,omitemtpy"`
	// TODO: global route rules that live on the virtualhost level
}

type Route struct {
	Matcher       Matcher         `json:"matcher"`
	Destination   Destination     `json:"destination"`
	RewritePrefix string          `json:"rewrite_prefix"`
	Plugins       RoutePluginSpec `json:"plugins"`
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
	Weight uint32
}

type SingleDestination struct {
	// A valid destination can only contain one of:
	// FunctionDestination
	// UpstreamDestination
	FunctionDestination *FunctionDestination `json:"function_destination,omitemtpy"`
	UpstreamDestination *UpstreamDestination `json:"upstream_destination,omitemtpy"`
}

type Matcher struct {
	Path        Path              `json:"path"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	Verbs       []string          `json:"verbs"`
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
	storageRef

	Name              string        `json:"name"`
	Type              UpstreamType  `json:"type"`
	Spec              UpstreamSpec  `json:"spec"`
	ConnectionTimeout time.Duration `json:"connection_timeout"`
	Functions         []Function    `json:"functions"`
}

type Function struct {
	Name string `json:"name"`
	// upstream ref?
	Spec FunctionSpec `json:"spec"`
}

type SSLConfig struct {
	CACertPath string `json:"ca_cert_path"`
	SecretRef  string `json:"secret_ref"`
}
