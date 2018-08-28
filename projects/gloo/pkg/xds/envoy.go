package xds

import (
	"context"
	"fmt"
	"sync"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyv2 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	envoyserver "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"google.golang.org/grpc"
)

// TODO(ilackarms): create a new xds server on each sync if possible, need to know the names of the proxy configs

// used to let nodes know they have a bad config
// we assign a "fix me" snapshot for bad nodes
const fallbackNodeKey = "misconfigured-node"

const (
	// TODO(ilackarms): expose these as a configuration option (maybe)
	fallbackBindPort   = 80
	fallbackBindAddr   = "::"
	fallbackStatusCode = 500
)

// One hasher, like the cache, is meant to be shared between multiple event loops
// Allows us to serve XDS for configurations in multiple namespaces
type EnvoyInstanceHasher struct {
	ctx       context.Context
	validKeys map[string][]string
	lock      *sync.RWMutex
}

const errorString = `
Envoy proxies are assigned configuration by Gloo based on their Node ID.
Proxies must register to Gloo with their node ID in the format "NAMESPACE~NAME"
Where NAMESPACE and NAME are the namespace and name of the correlating Proxy resource.`

func (h *EnvoyInstanceHasher) ID(node *core.Node) string {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for namespace, keys := range h.validKeys {
		for _, key := range keys {
			// This is where Node ID is defined from namespace~name
			nodeId := fmt.Sprintf("%v~%v", namespace, key)
			if nodeId == node.Id {
				return nodeId
			}
		}
	}
	contextutils.LoggerFrom(h.ctx).Warnf("invalid id provided by Envoy: %v", node.Id)
	contextutils.LoggerFrom(h.ctx).Debugf(errorString)
	return fallbackNodeKey
}

// Called in translator
func (h *EnvoyInstanceHasher) SetValidKeys(namespace string, validKeys []string) {
	h.lock.Unlock()
	defer h.lock.Lock()
	h.validKeys[namespace] = validKeys
}

func newNodeHasher(ctx context.Context) *EnvoyInstanceHasher {
	return &EnvoyInstanceHasher{
		ctx:       ctx,
		validKeys: make(map[string][]string),
		lock:      &sync.RWMutex{},
	}
}

func SetupEnvoyXds(ctx context.Context, grpcServer *grpc.Server, callbacks envoyserver.Callbacks) (*EnvoyInstanceHasher, envoycache.SnapshotCache) {
	ctx = contextutils.WithLogger(ctx, "envoy-xds-server")
	hasher := newNodeHasher(ctx)
	envoyCache := envoycache.NewSnapshotCache(true, hasher, contextutils.LoggerFrom(ctx))
	xdsServer := envoyserver.NewServer(envoyCache, callbacks)
	envoyv2.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterEndpointDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterClusterDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterRouteDiscoveryServiceServer(grpcServer, xdsServer)
	v2.RegisterListenerDiscoveryServiceServer(grpcServer, xdsServer)
	envoyCache.SetSnapshot(fallbackNodeKey, fallbackSnapshot(fallbackBindAddr, fallbackBindPort, fallbackStatusCode))

	return hasher, envoyCache
}
