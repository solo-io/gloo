package extensions

import (
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type RouteExtensionSpec struct {
	AddRequestHeaders     []HeaderValue `json:"add_request_headers,omitempty"`
	AddResponseHeaders    []HeaderValue `json:"add_response_headers,omitempty"`
	RemoveResponseHeaders []string      `json:"remove_response_headers,omitempty"`

	MaxRetries  uint32        `json:"max_retries,omitempty"`
	Timeout     time.Duration `json:"timeout,omitempty"`
	HostRewrite string        `json:"host_rewrite,omitempty"`

	//TODO: support RateLimit
	//TODO: support CORS policy
}

type HeaderValue struct {
	Key    string `json:"key,omitempty"`
	Value  string `json:"value,omitempty"`
	Append bool   `json:"append,omitempty"`
}

func DecodeRouteExtensions(generic *types.Struct) (RouteExtensionSpec, error) {
	var s RouteExtensionSpec
	err := protoutil.UnmarshalStruct(generic, &s)
	return s, err
}

func EncodeUpstreamSpec(spec RouteExtensionSpec) *types.Struct {
	v1Spec, err := protoutil.MarshalStruct(spec)
	if err != nil {
		panic(err)
	}
	return v1Spec
}
