package ratelimit

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"go.uber.org/zap"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	rateLimitPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

type RateLimitTranslatorSyncerExtension struct {
	settings *ratelimit.ServiceSettings
}

// TODO(kdorosh) delete once we stop supporting opaque rate-limiting config
type extensionsContainer struct {
	params syncer.TranslatorSyncerExtensionParams
}

// TODO(kdorosh) delete once we stop supporting opaque rate-limiting config
func (t *extensionsContainer) GetExtensions() *gloov1.Extensions {
	return t.params.SettingExtensions
}

// TODO(kdorosh) delete once we stop supporting opaque rate-limiting config
func getDeprecatedSettings(params syncer.TranslatorSyncerExtensionParams) (ratelimit.ServiceSettings, error) {
	var settings ratelimit.EnvoySettings
	err := utils.UnmarshalExtension(&extensionsContainer{params}, rateLimitPlugin.EnvoyExtensionName, &settings)
	if err != nil {
		if err == utils.NotFoundError {
			return ratelimit.ServiceSettings{}, nil
		}
		return ratelimit.ServiceSettings{}, err
	}
	return ratelimit.ServiceSettings{Descriptors: settings.GetCustomConfig().GetDescriptors()}, nil
}

func NewTranslatorSyncerExtension(ctx context.Context, params syncer.TranslatorSyncerExtensionParams) (syncer.TranslatorSyncerExtension, error) {
	var settings ratelimit.ServiceSettings
	settings, err := getDeprecatedSettings(params)
	if err != nil {
		return nil, err
	}

	if params.RateLimitDescriptorSettings.GetDescriptors() != nil {
		settings.Descriptors = params.RateLimitDescriptorSettings.Descriptors
	}

	return &RateLimitTranslatorSyncerExtension{
		settings: &settings,
	}, nil
}

func (s *RateLimitTranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) error {
	ctx = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, )", snap.Hash(),
		len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts))
	defer logger.Infof("end sync %v", snap.Hash())

	var customrl *v1.RateLimitConfig
	if s.settings.GetDescriptors() != nil {
		customrl = &v1.RateLimitConfig{
			Domain:      rateLimitPlugin.CustomDomain,
			Descriptors: s.settings.Descriptors,
		}
	}

	rl := &v1.RateLimitConfig{
		Domain: rateLimitPlugin.IngressDomain,
	}

	for _, proxy := range snap.Proxies {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*gloov1.Listener_HttpListener)
			if !ok {
				continue
			}

			virtualHosts := httpListener.HttpListener.VirtualHosts
			for _, virtualHost := range virtualHosts {
				var rateLimitDeprecated ratelimit.IngressRateLimit
				rateLimit := virtualHost.GetVirtualHostPlugins().GetRatelimitBasic()
				err := utils.UnmarshalExtension(virtualHost.VirtualHostPlugins, rateLimitPlugin.ExtensionName, &rateLimitDeprecated)
				if err != nil {
					if err == utils.NotFoundError && rateLimit == nil {
						// no rate limit virtual host config found, nothing to do here
						continue
					}
					return errors.Wrapf(err, "Error converting proto any to ingress rate limit plugin")
				}
				virtualHost = proto.Clone(virtualHost).(*gloov1.VirtualHost)
				virtualHost.Name = glooutils.SanitizeForEnvoy(ctx, virtualHost.Name, "virtual host")

				if rateLimit == nil {
					rateLimit = &rateLimitDeprecated
				}

				vhostConstraint, err := rateLimitPlugin.TranslateUserConfigToRateLimitServerConfig(virtualHost.Name, *rateLimit)
				if err != nil {
					return err
				}
				rl.Descriptors = append(rl.Descriptors, vhostConstraint)
			}
		}
	}

	// TODO(yuval-k): we should make sure that we add the proxy name and listener name to the descriptors
	var resources []envoycache.Resource

	resource := v1.NewRateLimitConfigXdsResourceWrapper(rl)
	resources = append(resources, resource)
	if customrl != nil {
		resource := v1.NewRateLimitConfigXdsResourceWrapper(customrl)
		resources = append(resources, resource)
	}
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
