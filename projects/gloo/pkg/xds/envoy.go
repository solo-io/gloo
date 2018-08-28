package xds

import (
	"context"
	"strings"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// TODO(ilackarms): create a new xds server on each sync if possible, need to know the names of the proxy configs

// used to let nodes know they have a bad config
// we assign a "fix me" virtualhost for bad nodes
const badNodeKey = "misconfigured-node"

type hasher struct {
	ctx context.Context
}

const errorString = `
Envoy proxies are assigned configuration by Gloo based on their Node ID.
Proxies must register to Gloo with their node ID in the format "NAMESPACE~NAME"
Where NAMESPACE and NAME are the namespace and name of the correlating Proxy resource.`

func (h hasher) ID(node *core.Node) string {
	parts := strings.SplitN(node.Id, "~", 2)
	if len(parts) != 2 {
		contextutils.LoggerFrom(h.ctx).Warnf("invalid id provided by Envoy: %v", node.Id)
		contextutils.LoggerFrom(h.ctx).Debugf(errorString)
		return badNodeKey
	}
	contextutils.LoggerFrom(h.ctx).Debugf("node %v registered with role %v", parts[1], parts[0])
	return parts[0]
}

func NewEnvoyXds(ctx context.Context) {
	ctx = contextutils.WithLogger(ctx, "envoy-xds-server")
	envoyCache := envoycache.NewSnapshotCache(true, hasher{}, contextutils.LoggerFrom(ctx))
}
