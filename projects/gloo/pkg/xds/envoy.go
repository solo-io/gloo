package xds

import (
	"context"
	"fmt"
	"sync"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyv2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/cache"
	envoyserver "github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/server"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"google.golang.org/grpc"
)

// used to let nodes know they have a bad config
// we assign a "fix me" snapshot for bad nodes
const fallbackNodeKey = "misconfigured-node"

// TODO(ilackarms): expose these as a configuration option (maybe)
var fallbackBindPort = defaults.HttpPort

const (
	fallbackBindAddr   = "::"
	fallbackStatusCode = 500
)

type ProxyKeyHasher struct {
	ctx context.Context
	// (ilackarms) for the purpose of invalidation in the hasher
	validKeysLock sync.Mutex
	validKeys     []string
	lock          *sync.RWMutex
}

const errorString = `
Envoy proxies are assigned configuration by Gloo based on their Node ID.
Proxies must register to Gloo with their node ID in the format "NAMESPACE~NAME"
Where NAMESPACE and NAME are the namespace and name of the correlating Proxy resource.`

func (h *ProxyKeyHasher) ID(node *core.Node) string {

	role := ""
	if node.Metadata != nil {
		roleValue := node.Metadata.Fields["role"]
		if roleValue != nil {
			role = roleValue.GetStringValue()
		}
	}
	// TODO(yuval-k): once go-control-plane is implemented here we can implement default snapshot in it.
	return role
}

func SnapshotKey(proxy *v1.Proxy) string {
	namespace, name := proxy.GetMetadata().Ref().Strings()
	return fmt.Sprintf("%v~%v", namespace, name)
}

// Called in Syncer when a new set of proxies arrive
func (h *ProxyKeyHasher) SetKeysFromProxies(proxies v1.ProxyList) {
	var validKeys []string
	// This is where we correlate Node ID with proxy namespace~name
	for _, proxy := range proxies {
		validKeys = append(validKeys, SnapshotKey(proxy))
	}

	h.validKeysLock.Lock()
	h.validKeys = validKeys
	h.validKeysLock.Unlock()
}

func newNodeHasher(ctx context.Context) *ProxyKeyHasher {
	return &ProxyKeyHasher{
		ctx:  ctx,
		lock: &sync.RWMutex{},
	}
}

func SetupEnvoyXds(ctx context.Context, grpcServer *grpc.Server, callbacks envoyserver.Callbacks) (*ProxyKeyHasher, envoycache.SnapshotCache) {
	ctx = contextutils.WithLogger(ctx, "envoy-xds-server")
	hasher := newNodeHasher(ctx)
	envoyCache := envoycache.NewSnapshotCache(true, hasher, contextutils.LoggerFrom(ctx))
	xdsServer := envoyserver.NewServer(envoyCache, callbacks)
	// TODO(yuval-k): once we support generic cache, move this to a higher level.
	// we will probably invert dependencies and have this function received the cache, rather than produce it.
	envoyv2.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsServer)
	envoyServer := NewEnvoyServer(xdsServer)

	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, envoyServer)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, envoyServer)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, envoyServer)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, envoyServer)
	envoyCache.SetSnapshot(fallbackNodeKey, fallbackSnapshot(fallbackBindAddr, fallbackBindPort, fallbackStatusCode))

	return hasher, envoyCache
}
