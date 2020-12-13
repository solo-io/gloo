package ratelimit

import (
	"context"
	"fmt"

	solo_api_rl "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"

	"github.com/hashicorp/go-multierror"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit/collectors"
	rate_limiter_shims "github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation"

	"github.com/golang/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/go-utils/hashutils"
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
)

// TODO(marco): generate these in solo-kit
//go:generate mockgen -package mocks -destination mocks/cache.go github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache SnapshotCache
//go:generate mockgen -package mocks -destination mocks/reporter.go github.com/solo-io/solo-kit/pkg/api/v2/reporter Reporter

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

func init() {
	_ = view.Register(rlConnectedStateView)
}

type translatorSyncerExtension struct {
	collectorFactory collectors.ConfigCollectorFactory
	domainGenerator  rate_limiter_shims.RateLimitDomainGenerator
	reporter         reporter.Reporter
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
		params.Reporter,
	), nil
}

func NewTranslatorSyncer(
	collectorFactory collectors.ConfigCollectorFactory,
	domainGenerator rate_limiter_shims.RateLimitDomainGenerator,
	reporter reporter.Reporter,
) syncer.TranslatorSyncerExtension {
	return &translatorSyncerExtension{
		collectorFactory: collectorFactory,
		domainGenerator:  domainGenerator,
		reporter:         reporter,
	}
}

func (s *translatorSyncerExtension) Sync(ctx context.Context, snap *gloov1.ApiSnapshot, xdsCache envoycache.SnapshotCache) (string, error) {
	ctx = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	snapHash := hashutils.MustHash(snap)
	logger.Infof("begin rate limit sync %v (%v proxies, %v rate limit configs)", snapHash, len(snap.Proxies), len(snap.Ratelimitconfigs))
	defer logger.Infof("end sync %v", snapHash)

	reports := make(reporter.ResourceReports)
	reports.Accept(snap.Proxies.AsInputResources()...)
	reports.Accept(snap.Ratelimitconfigs.AsInputResources()...)

	// Make sure we always write reports
	defer func() {
		if err := s.reporter.WriteReports(ctx, reports, nil); err != nil {
			contextutils.LoggerFrom(ctx).Warnf("Failed writing report for rate limit configs: %v", err)
		}
	}()

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
				configCollectors.processVirtualHost(virtualHostClone, proxy)

				// Process rate limit configuration on routes
				for _, route := range virtualHost.GetRoutes() {
					configCollectors.processRoute(route, virtualHostClone, proxy)
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
		collector, err := collectorFactory.MakeInstance(collectorType, snapshot, reports, logger)
		if err != nil {
			return configCollectorSet{}, err
		}
		set.collectors = append(set.collectors, collector)
	}
	return set, nil
}

func (c configCollectorSet) processVirtualHost(virtualHost *gloov1.VirtualHost, proxy *gloov1.Proxy) {
	for _, collector := range c.collectors {
		collector.ProcessVirtualHost(virtualHost, proxy)
	}
}

func (c configCollectorSet) processRoute(route *gloov1.Route, virtualHost *gloov1.VirtualHost, proxy *gloov1.Proxy) {
	for _, collector := range c.collectors {
		collector.ProcessRoute(route, virtualHost, proxy)
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
