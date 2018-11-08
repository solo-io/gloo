package syncer

import (
	"context"
	"fmt"

	"github.com/mitchellh/hashstructure"

	envoycache "github.com/solo-io/solo-kit/pkg/control-plane/cache"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/ratelimit"
)

func (s *translatorSyncer) syncRateLimit(ctx context.Context, snap *v1.ApiSnapshot) error {
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
				cfg, err := ratelimit.TranslateUserConfigToRateLimitServerConfig(*virtualHost.VirtualHostPlugins.RateLimits)
				if err != nil {
					return err
				}
				resource := v1.NewRateLimitConfigXdsResourceWrapper(cfg)
				resources := []envoycache.Resource{resource}
				h, err := hashstructure.Hash(resources, nil)
				if err != nil {
					panic(err)
				}
				rlsnap := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
				s.xdsCache.SetSnapshot("ratelimit", rlsnap)
				// very soon. either buy changing the plugin or potentially the filter.
				return nil

			}
		}
	}

	// find our plugin
	return nil
}
