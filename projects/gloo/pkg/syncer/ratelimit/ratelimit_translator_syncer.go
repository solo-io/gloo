package ratelimit

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/hashstructure"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/hashutils"
	solo_api_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors"
	rate_limiter_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
)

// TODO(marco): generate these in solo-kit
//go:generate mockgen -package mocks -destination mocks/cache.go github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache SnapshotCache
//go:generate mockgen -package mocks -destination mocks/reporter.go github.com/solo-io/solo-kit/pkg/api/v2/reporter Reporter

// Compile-time assertion
var (
	_ syncer.TranslatorSyncerExtension            = new(translatorSyncerExtension)
	_ syncer.UpgradeableTranslatorSyncerExtension = new(translatorSyncerExtension)
)

var (
	rlConnectedStateDescription = "zero indicates gloo detected an error with the rate limit config and did not update its XDS snapshot, check the gloo logs for errors"
	rlConnectedState            = stats.Int64("glooe.ratelimit/connected_state", rlConnectedStateDescription, "1")

	rlConnectedStateView = &view.View{
		Name:        "glooe.ratelimit/connected_state",
		Measure:     rlConnectedState,
		Description: rlConnectedStateDescription,
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{},
	}
)

const (
	Name = "rate-limit"
)

func init() {
	_ = view.Register(rlConnectedStateView)
}

type translatorSyncerExtension struct {
	collectorFactory collectors.ConfigCollectorFactory
	domainGenerator  rate_limiter_shims.RateLimitDomainGenerator
}

// The rate limit server sends xDS discovery requests to Gloo to get its configuration from Gloo. This constant determines
// the value of the nodeInfo.Metadata.role field that the server sends along to retrieve its configuration snapshot,
// similarly to how the regular Gloo gateway-proxies do.
const RateLimitServerRole = "ratelimit"

func NewTranslatorSyncerExtension(_ context.Context, params syncer.TranslatorSyncerExtensionParams) (syncer.TranslatorSyncerExtension, error) {
	var settings ratelimit.ServiceSettings

	if params.RateLimitServiceSettings.GetDescriptors() != nil {
		settings.Descriptors = params.RateLimitServiceSettings.Descriptors
	}
	if params.RateLimitServiceSettings.GetSetDescriptors() != nil {
		settings.SetDescriptors = params.RateLimitServiceSettings.SetDescriptors
	}

	return NewTranslatorSyncer(
		collectors.NewCollectorFactory(
			&settings,
			rate_limiter_shims.NewGlobalRateLimitTranslator(),
			rate_limiter_shims.NewRateLimitConfigTranslator(),
			translation.NewBasicRateLimitTranslator(),
		),
		rate_limiter_shims.NewRateLimitDomainGenerator(),
	), nil
}

func NewTranslatorSyncer(
	collectorFactory collectors.ConfigCollectorFactory,
	domainGenerator rate_limiter_shims.RateLimitDomainGenerator,
) syncer.TranslatorSyncerExtension {
	return &translatorSyncerExtension{
		collectorFactory: collectorFactory,
		domainGenerator:  domainGenerator,
	}
}

func (s *translatorSyncerExtension) Sync(
	ctx context.Context,
	snap *gloov1.ApiSnapshot,
	settings *gloov1.Settings,
	xdsCache envoycache.SnapshotCache,
	reports reporter.ResourceReports,
) (string, error) {
	ctx = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	snapHash := hashutils.MustHash(snap)
	logger.Infof("begin rate limit sync %v (%v proxies, %v rate limit configs)", snapHash, len(snap.Proxies), len(snap.Ratelimitconfigs))
	defer logger.Infof("end sync %v", snapHash)

	reports.Accept(snap.Proxies.AsInputResources()...)
	reports.Accept(snap.Ratelimitconfigs.AsInputResources()...)

	configCollectors, err := newCollectorSet(s.collectorFactory, snap, reports, logger)
	if err != nil {
		return syncerError(ctx, err)
	}

	for _, proxy := range snap.Proxies {
		for _, listener := range proxy.Listeners {
			httpListener, ok := listener.ListenerType.(*gloov1.Listener_HttpListener)
			if !ok {
				// Rate limiting is currently supported only on HTTP listeners
				continue
			}

			for _, virtualHost := range httpListener.HttpListener.VirtualHosts {

				// Sanitize the name of the virtual host
				virtualHostClone := proto.Clone(virtualHost).(*gloov1.VirtualHost)
				virtualHostClone.Name = glooutils.SanitizeForEnvoy(ctx, virtualHostClone.Name, "virtual host")

				// Process rate limit configuration on this virtual host
				configCollectors.processVirtualHost(virtualHostClone, proxy, reports)

				// Process rate limit configuration on routes
				for _, route := range virtualHost.GetRoutes() {
					configCollectors.processRoute(route, virtualHostClone, proxy, reports)
				}
			}
		}
	}

	var (
		errs              = &multierror.Error{}
		snapshotResources []envoycache.Resource
	)

	// Get all the collected rate limit config
	configs, err := configCollectors.toXdsConfigurations()
	errs = multierror.Append(errs, err)

	for _, cfg := range configs {
		// Verify xDS configuration can be translated to valid RL server configuration
		if _, err := s.domainGenerator.NewRateLimitDomain(ctx, cfg.Domain,
			&solo_api_rl.RateLimitConfigSpec_Raw{Descriptors: cfg.Descriptors, SetDescriptors: cfg.SetDescriptors}); err != nil {
			errs = multierror.Append(errs, err)
		}
		snapshotResources = append(snapshotResources, v1.NewRateLimitConfigXdsResourceWrapper(cfg))
	}
	if err := errs.ErrorOrNil(); err != nil {
		return syncerError(ctx, err)
	}

	hashedResources, err := hashstructure.Hash(snapshotResources, nil)
	if err != nil {
		contextutils.LoggerFrom(ctx).With(zap.Error(err)).DPanic("error hashing rate limit")
		return syncerError(ctx, err)
	}

	rateLimitSnapshot := envoycache.NewEasyGenericSnapshot(fmt.Sprintf("%d", hashedResources), snapshotResources)

	err = xdsCache.SetSnapshot(RateLimitServerRole, rateLimitSnapshot)
	if err != nil {
		return syncerError(ctx, err)
	}

	stats.Record(ctx, rlConnectedState.M(int64(1)))

	return RateLimitServerRole, nil
}

func (s *translatorSyncerExtension) ExtensionName() string {
	return Name
}

func (s *translatorSyncerExtension) IsUpgrade() bool {
	return true
}

func syncerError(ctx context.Context, err error) (string, error) {
	stats.Record(ctx, rlConnectedState.M(int64(0)))
	return RateLimitServerRole, err
}

// Helper object to reduce boilerplate.
type configCollectorSet struct {
	collectors []collectors.ConfigCollector
}

func newCollectorSet(
	collectorFactory collectors.ConfigCollectorFactory,
	snapshot *gloov1.ApiSnapshot,
	reports reporter.ResourceReports,
	logger *zap.SugaredLogger,
) (configCollectorSet, error) {
	set := configCollectorSet{}

	for _, collectorType := range []collectors.CollectorType{
		collectors.Global,
		collectors.Basic,
		collectors.Crd,
	} {
		collector, err := collectorFactory.MakeInstance(collectorType, snapshot, logger)
		if err != nil {
			return configCollectorSet{}, err
		}
		set.collectors = append(set.collectors, collector)
	}
	return set, nil
}

func (c configCollectorSet) processVirtualHost(
	virtualHost *gloov1.VirtualHost,
	proxy *gloov1.Proxy,
	reports reporter.ResourceReports,
) {
	for _, collector := range c.collectors {
		collector.ProcessVirtualHost(virtualHost, proxy, reports)
	}
}

func (c configCollectorSet) processRoute(
	route *gloov1.Route,
	virtualHost *gloov1.VirtualHost,
	proxy *gloov1.Proxy,
	reports reporter.ResourceReports,
) {
	for _, collector := range c.collectors {
		collector.ProcessRoute(route, virtualHost, proxy, reports)
	}
}

func (c configCollectorSet) toXdsConfigurations() ([]*v1.RateLimitConfig, error) {
	var result []*v1.RateLimitConfig
	var errs *multierror.Error
	for _, collector := range c.collectors {
		config, err := collector.ToXdsConfiguration()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
		result = append(result, config)
	}
	return result, errs.ErrorOrNil()
}
