package extauth

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"go.uber.org/zap"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	extAuthPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

type ExtAuthTranslatorSyncerExtension struct {
}

var _ syncer.TranslatorSyncerExtension = NewTranslatorSyncerExtension()

func NewTranslatorSyncerExtension() *ExtAuthTranslatorSyncerExtension {
	return &ExtAuthTranslatorSyncerExtension{}
}

// TODO: move this into solo-kit
type SnapshotSetter interface {
	SetSnapshot(node string, snapshot envoycache.Snapshot) error
}

func (s *ExtAuthTranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) error {
	return s.SyncAndSet(ctx, snap, xdsCache)
}

func (s *ExtAuthTranslatorSyncerExtension) SyncAndSet(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache SnapshotSetter) error {
	var cfgs []*extauth.ExtAuthConfig

	for _, proxy := range snap.Proxies {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*gloov1.Listener_HttpListener)
			if !ok {
				// not an http listener - skip it as currently ext auth is only supported for http
				continue
			}

			virtualHosts := httpListener.HttpListener.VirtualHosts
			for _, virtualHost := range virtualHosts {

				var extAuthVhost extauth.VhostExtension
				err := utils.UnmarshalExtension(virtualHost.VirtualHostPlugins, extAuthPlugin.ExtensionName, &extAuthVhost)
				if err != nil {
					if err == utils.NotFoundError {
						// no ext auth extension on this vhost - nothing else to do.
						continue
					}
					return errors.Wrapf(err, "Error converting proto any to ingress ext auth plugin")
				}

				extath, err := extAuthPlugin.TranslateUserConfigToExtAuthServerConfig(proxy, listener, virtualHost, snap, extAuthVhost)
				if err != nil {
					return err
				}
				if extath != nil {
					cfgs = append(cfgs, extath)
				}
			}
		}
	}

	resources := []envoycache.Resource{}
	for _, cfg := range cfgs {
		resource := extauth.NewExtAuthConfigXdsResourceWrapper(cfg)
		resources = append(resources, resource)
	}
	h, err := hashstructure.Hash(resources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing ext auth")
		return err
	}
	rlsnap := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
	xdsCache.SetSnapshot("extauth", rlsnap)
	// find our plugin
	return nil
}
