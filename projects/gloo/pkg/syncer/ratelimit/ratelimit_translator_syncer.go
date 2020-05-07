package ratelimit

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	rlPlugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/go-utils/hashutils"
	configproto "github.com/solo-io/solo-projects/projects/rate-limit/pkg/config"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"

	"github.com/mitchellh/hashstructure"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	rlIngressPlugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

var (
	rlConnectedStateDescription = "0 indicates gloo detected an error with the rate limit config and did not update its XDS snapshot, check the gloo logs for errors"
	rlConnectedState            = stats.Int64("glooe.ratelimit/connected_state", rlConnectedStateDescription, "1")

	rlConnectedStateView = &view.View{
		Name:        "glooe.ratelimit/connected_state",
		Measure:     rlConnectedState,
		Description: rlConnectedStateDescription,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	_ = view.Register(rlConnectedStateView)
}

type RateLimitTranslatorSyncerExtension struct {
	settings *ratelimit.ServiceSettings
}

// The ratelimit server sends xDS discovery requests to Gloo to get its configuration from Gloo. This constant determines
// the value of the nodeInfo.Metadata.role field that the server sends along to retrieve its configuration snapshot,
// similarly to how the regular Gloo gateway-proxies do.
const RateLimitServerRole = "ratelimit"

func NewTranslatorSyncerExtension(ctx context.Context, params syncer.TranslatorSyncerExtensionParams) (syncer.TranslatorSyncerExtension, error) {
	var settings ratelimit.ServiceSettings

	if params.RateLimitServiceSettings.GetDescriptors() != nil {
		settings.Descriptors = params.RateLimitServiceSettings.Descriptors
	}

	return &RateLimitTranslatorSyncerExtension{
		settings: &settings,
	}, nil
}

func (s *RateLimitTranslatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) (string, error) {
	ctx = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	snapHash := hashutils.MustHash(snap)
	logger.Infof("begin sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, %v auth configs)", snapHash,
		len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts), len(snap.AuthConfigs))
	defer logger.Infof("end sync %v", snapHash)

	customrl := &v1.RateLimitConfig{
		Domain: rlPlugin.CustomDomain,
	}
	if s.settings.GetDescriptors() != nil {
		customrl.Descriptors = s.settings.GetDescriptors()
	}

	rl := &v1.RateLimitConfig{
		Domain: rlIngressPlugin.IngressDomain,
	}

	for _, proxy := range snap.Proxies {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*gloov1.Listener_HttpListener)
			if !ok {
				continue
			}

			virtualHosts := httpListener.HttpListener.VirtualHosts
			for _, virtualHost := range virtualHosts {
				rateLimit := virtualHost.GetOptions().GetRatelimitBasic()
				if rateLimit == nil {
					// no rate limit virtual host config found, nothing to do here
					continue
				}
				virtualHost = proto.Clone(virtualHost).(*gloov1.VirtualHost)
				virtualHost.Name = glooutils.SanitizeForEnvoy(ctx, virtualHost.Name, "virtual host")

				vhostConstraint, err := rlIngressPlugin.TranslateUserConfigToRateLimitServerConfig(virtualHost.Name, *rateLimit)
				if err != nil {
					return syncerError(ctx, err)
				}
				rl.Descriptors = append(rl.Descriptors, vhostConstraint)
			}
		}
	}

	// TODO(yuval-k): we should make sure that we add the proxy name and listener name to the descriptors
	resources := []envoycache.Resource{
		v1.NewRateLimitConfigXdsResourceWrapper(customrl),
		v1.NewRateLimitConfigXdsResourceWrapper(rl),
	}

	// Verify settings can be translated to valid RL config
	generator := configproto.NewConfigGenerator(contextutils.LoggerFrom(ctx))
	if _, err := generator.GenerateConfig([]*v1.RateLimitConfig{customrl, rl}); err != nil {
		return syncerError(ctx, err)
	}

	h, err := hashstructure.Hash(resources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing rate limit")
		return syncerError(ctx, err)
	}
	rlsnap := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", h), resources)
	xdsCache.SetSnapshot(RateLimitServerRole, rlsnap)

	stats.Record(ctx, rlConnectedState.M(int64(1)))

	// find our plugin
	return RateLimitServerRole, nil
}

func syncerError(ctx context.Context, err error) (string, error) {
	stats.Record(ctx, rlConnectedState.M(int64(0)))
	return RateLimitServerRole, err
}
