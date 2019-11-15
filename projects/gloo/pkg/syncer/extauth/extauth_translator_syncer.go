package extauth

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/projects/extauth/pkg/runner"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"go.uber.org/zap"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"github.com/solo-io/go-utils/contextutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	extAuthPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

type TranslatorSyncerExtension struct{}

var _ syncer.TranslatorSyncerExtension = NewTranslatorSyncerExtension()

func NewTranslatorSyncerExtension() *TranslatorSyncerExtension {
	return &TranslatorSyncerExtension{}
}

// TODO(marco): report errors on auth config resources once we have the strongly typed API. Currently it is not possible
//  to do this consistently, since we need to parse the raw extension to get to the auth config, an operation that might itself fail.
func (s *TranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) error {
	ctx = contextutils.WithLogger(ctx, "extAuthTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin auth sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, %v auth configs)", snap.Hash(),
		len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts), len(snap.AuthConfigs))
	defer logger.Infof("end auth sync %v", snap.Hash())

	return s.SyncAndSet(ctx, snap, xdsCache)
}

type SnapshotSetter interface {
	SetSnapshot(node string, snapshot envoycache.Snapshot) error
}

func (s *TranslatorSyncerExtension) SyncAndSet(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache SnapshotSetter) error {
	helper := newHelper()

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

				if err := helper.processAuthExtension(ctx, snap, virtualHost.GetVirtualHostPlugins().GetExtauth()); err != nil {
					return err
				}

				for _, route := range virtualHost.Routes {

					if err := helper.processAuthExtension(ctx, snap, route.GetRoutePlugins().GetExtauth()); err != nil {
						return err
					}

					for _, weightedDestination := range route.GetRouteAction().GetMulti().GetDestinations() {
						if err := helper.processAuthExtension(ctx, snap, weightedDestination.GetWeightedDestinationPlugins().GetExtauth()); err != nil {
							return err
						}
					}
				}
			}
		}
	}

	var resources []envoycache.Resource
	for _, cfg := range helper.translatedConfigs {
		resource := extauth.NewExtAuthConfigXdsResourceWrapper(cfg)
		resources = append(resources, resource)
	}
	h, err := hashstructure.Hash(resources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing ext auth")
		return err
	}
	extAuthSnapshot := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
	_ = xdsCache.SetSnapshot(runner.ExtAuthServerRole, extAuthSnapshot)
	return nil
}

// This translation helper contains a map where each key is the unique identifier of an AuthConfig and the corresponding
// value is the translated config. We use it avoid translating the same configuration multiple times.
type helper struct {
	translatedConfigs map[string]*extauth.ExtAuthConfig
}

func newHelper() *helper {
	return &helper{
		translatedConfigs: make(map[string]*extauth.ExtAuthConfig),
	}
}

func (h *helper) processAuthExtension(ctx context.Context, snap *gloov1.ApiSnapshot, config *extauth.ExtAuthExtension) error {
	configRef := config.GetConfigRef()
	if configRef == nil {
		// Just return if there is nothing to translate
		return nil
	}

	// Don't perform duplicate work if we already have translated this resource
	if _, ok := h.translatedConfigs[configRef.Key()]; ok {
		return nil
	}

	translatedConfig, err := extAuthPlugin.TranslateExtAuthConfig(ctx, snap, configRef)
	if err != nil {
		return err
	} else if translatedConfig == nil {
		// Do nothing if config is empty or consists only of custom auth configs
		return nil
	}

	h.translatedConfigs[configRef.Key()] = translatedConfig
	return nil
}
