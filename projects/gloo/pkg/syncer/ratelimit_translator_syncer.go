package syncer

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/utils/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	rateLimitPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

type RateLimitTranslatorSyncerExtension struct {
}

func NewTranslatorSyncerExtension() syncer.TranslatorSyncerExtension {
	return &RateLimitTranslatorSyncerExtension{}
}

func (s *RateLimitTranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) error {
	for _, proxy := range snap.Proxies.List() {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*gloov1.Listener_HttpListener)
			if !ok {
				continue
			}
			virtualHosts := httpListener.HttpListener.VirtualHosts
			for _, virtualHost := range virtualHosts {
				if virtualHost.VirtualHostPlugins == nil {
					continue
				}
				if virtualHost.VirtualHostPlugins.Plugins == nil {
					continue
				}
				plugins := virtualHost.VirtualHostPlugins.Plugins
				var rateLimit ratelimit.IngressRateLimit
				err := proto.UnmarshalAnyFromMap(plugins, rateLimitPlugin.PluginName, &rateLimit)
				if err != nil {
					if err == proto.NotFoundError {
						continue
					}
					return errors.Wrapf(err, "Error converting proto any to ingress rate limit plugin")
				}

				cfg, err := rateLimitPlugin.TranslateUserConfigToRateLimitServerConfig(rateLimit)
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
				xdsCache.SetSnapshot("ratelimit", rlsnap)
				// very soon. either buy changing the plugin or potentially the filter.
				return nil

			}
		}
	}

	// find our plugin
	return nil
}
