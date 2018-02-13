package v1

import types "github.com/golang/protobuf/ptypes/struct"

type Config struct {
	Upstreams    []*Upstream
	VirtualHosts []*VirtualHost
}

type FunctionSpec *types.Struct
