package ratelimit

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"go.uber.org/zap"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	rateLimitPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

type RateLimitTranslatorSyncerExtension struct {
}

func NewTranslatorSyncerExtension() syncer.TranslatorSyncerExtension {
	return &RateLimitTranslatorSyncerExtension{}
}

func (s *RateLimitTranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) error {

	rl := &v1.RateLimitConfig{
		Domain: rateLimitPlugin.IngressDomain,
	}

	for _, proxy := range snap.Proxies.List() {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*gloov1.Listener_HttpListener)
			if !ok {
				continue
			}

			virtualHosts := httpListener.HttpListener.VirtualHosts
			for _, virtualHost := range virtualHosts {
				var rateLimit ratelimit.IngressRateLimit
				err := utils.UnmarshalExtension(virtualHost.VirtualHostPlugins, rateLimitPlugin.ExtensionName, &rateLimit)
				if err != nil {
					if err == utils.NotFoundError {
						continue
					}
					return errors.Wrapf(err, "Error converting proto any to ingress rate limit plugin")
				}

				vhostConstraint, err := rateLimitPlugin.TranslateUserConfigToRateLimitServerConfig(virtualHost.Name, rateLimit)
				if err != nil {
					return err
				}
				rl.Constraints = append(rl.Constraints, vhostConstraint)
			}
		}
	}

	// TODO(yuval-k): we should make sure that we add the proxy name and listener name to the descriptors
	var resources []envoycache.Resource

	resource := v1.NewRateLimitConfigXdsResourceWrapper(rl)
	resources = append(resources, resource)
	h, err := hashstructure.Hash(resources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing rate limit")
		return err
	}
	rlsnap := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
	xdsCache.SetSnapshot("ratelimit", rlsnap)

	// find our plugin
	return nil
}
