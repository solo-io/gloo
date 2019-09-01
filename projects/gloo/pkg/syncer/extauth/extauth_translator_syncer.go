package extauth

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
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
	ctx = contextutils.WithLogger(ctx, "extAuthTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, )", snap.Hash(),
		len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts))
	defer logger.Infof("end sync %v", snap.Hash())
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

				virtualHost = proto.Clone(virtualHost).(*gloov1.VirtualHost)
				virtualHost.Name = glooutils.SanitizeForEnvoy(ctx, virtualHost.Name, "virtual host")

				var extAuthVhost extauth.VhostExtension
				err := utils.UnmarshalExtension(virtualHost.VirtualHostPlugins, extAuthPlugin.ExtensionName, &extAuthVhost)
				if err != nil {
					if err == utils.NotFoundError {
						// no ext auth extension on this vhost - nothing else to do.
						continue
					}
					return errors.Wrapf(err, "Error converting proto to any ingress ext auth plugin")
				}

				extAuth, err := extAuthPlugin.TranslateUserConfigToExtAuthServerConfig(ctx, proxy, listener, virtualHost, snap, extAuthVhost)
				if err != nil {
					return err
				}
				if extAuth != nil {
					cfgs = append(cfgs, extAuth)
				}
			}
		}
	}

	var resources []envoycache.Resource
	for _, cfg := range cfgs {
		resource := extauth.NewExtAuthConfigXdsResourceWrapper(cfg)
		resources = append(resources, resource)
	}
	h, err := hashstructure.Hash(resources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing ext auth")
		return err
	}
	extAuthSnapshot := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
	_ = xdsCache.SetSnapshot("extauth", extAuthSnapshot)
	// find our plugin
	return nil
}
