package v1

import (
	"context"
	"errors"
	"fmt"

	discovery "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/cache"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/client"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/server"
)

// Type Definitions:

const RateLimitConfigType = cache.TypePrefix + "/solo.api." + "RateLimitConfig"

/* Defined a resource - to be used by snapshot */
type RateLimitConfigResource struct {
	resourceProto *RateLimitConfig
}

// Make sure the Resource interface is implemented
var _ cache.Resource = &RateLimitConfigResource{}

func NewRateLimitConfigResource(resourceProto *RateLimitConfig) *RateLimitConfigResource {
	return &RateLimitConfigResource{
		resourceProto: resourceProto,
	}
}

func (e *RateLimitConfigResource) Self() cache.ResourceReference {
	return cache.ResourceReference{Name: e.resourceProto.Domain, Type: RateLimitConfigType}
}

func (e *RateLimitConfigResource) ResourceProto() cache.ResourceProto {
	return e.resourceProto
}
func (e *RateLimitConfigResource) References() []cache.ResourceReference {
	return nil
}

// Define a type record. This is used by the generic client library.
var RateLimitConfigTypeRecord = client.NewTypeRecord(
	RateLimitConfigType,

	// Return an empty message, that can be used to deserialize bytes into it.
	func() cache.ResourceProto { return &RateLimitConfig{} },

	// Covert the message to a resource suitable for use for protobuf's Any.
	func(r cache.ResourceProto) cache.Resource {
		return &RateLimitConfigResource{resourceProto: r.(*RateLimitConfig)}
	},
)

// Server Implementation:

// Wrap the generic server and implement the type sepcific methods:
type rateLimitDiscoveryServiceServer struct {
	server.Server
}

func NewServer(genericServer server.Server) RateLimitDiscoveryServiceServer {
	return &rateLimitDiscoveryServiceServer{Server: genericServer}
}

func (s *rateLimitDiscoveryServiceServer) StreamRateLimitConfig(stream RateLimitDiscoveryService_StreamRateLimitConfigServer) error {
	return s.Server.Stream(stream, RateLimitConfigType)
}

func (s *rateLimitDiscoveryServiceServer) FetchRateLimitConfig(ctx context.Context, req *discovery.DiscoveryRequest) (*discovery.DiscoveryResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.Unavailable, "empty request")
	}
	req.TypeUrl = RateLimitConfigType
	return s.Server.Fetch(ctx, req)
}

func (s *rateLimitDiscoveryServiceServer) IncrementalRateLimitConfig(_ RateLimitDiscoveryService_IncrementalRateLimitConfigServer) error {
	return errors.New("not implemented")
}

// Client Implementation: Generate a strongly typed client over the generic client

// The apply functions receives resources and returns an error if they were applied correctly.
// In theory the configuration can become valid in the future (i.e. eventually consistent), but I don't think we need to worry about that now
// As our current use cases only have one configuration resource, so no interactions are expected.
type ApplyRateLimitConfig func(version string, resources []*RateLimitConfig) error

// Convert the strongly typed apply to a generic apply.
func apply(rlapply ApplyRateLimitConfig) func(cache.Resources) error {
	return func(resources cache.Resources) error {

		var configs []*RateLimitConfig
		for _, r := range resources.Items {
			if proto, ok := r.ResourceProto().(*RateLimitConfig); !ok {
				return fmt.Errorf("resource %s of type %s incorrect", r.Self().Name, r.Self().Type)
			} else {
				configs = append(configs, proto)
			}
		}

		return rlapply(resources.Version, configs)
	}
}

func NewRateLimitConfigClient(nodeinfo *core.Node, rlapply ApplyRateLimitConfig) client.Client {
	return client.NewClient(nodeinfo, RateLimitConfigTypeRecord, apply(rlapply))
}
