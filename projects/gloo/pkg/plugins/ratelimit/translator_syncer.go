package ratelimit

import (
	"context"
	"fmt"

	"github.com/mitchellh/hashstructure"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/cache"
)

type translatorSyncer struct {
	xdsCache envoycache.SnapshotCache
}

func NewTranslatorSyncer(xdsCache envoycache.SnapshotCache) *translatorSyncer {

	return &translatorSyncer{xdsCache: xdsCache}
}

func (t *translatorSyncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {
	for _, proxy := range snap.Proxies.List() {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*v1.Listener_HttpListener)
			if !ok {
				continue
			}
			virtualHosts := httpListener.HttpListener.VirtualHosts
			for _, virtualHost := range virtualHosts {
				if virtualHost.VirtualHostPlugins == nil {
					continue
				}
				if virtualHost.VirtualHostPlugins.RateLimits == nil {
					continue
				}
				cfg, err := translateUserConfigToRateLimitServerConfig(*virtualHost.VirtualHostPlugins.RateLimits)
				resource := v1.NewRateLimitConfigXdsResourceWrapper(cfg)
				resources := []envoycache.Resource{resource}
				h, err := hashstructure.Hash(resources, nil)
				if err != nil {
					panic(err)
				}
				rlsnap := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
				t.xdsCache.SetSnapshot("ratelimit", rlsnap)
				// TODO(yuval-k): For now we don't support more than one rate limit config, we need to solve this
				// very soon. either buy changing the plugin or potentially the filter.
				return nil

			}
		}
	}

	// find our plugin
	return nil
}
