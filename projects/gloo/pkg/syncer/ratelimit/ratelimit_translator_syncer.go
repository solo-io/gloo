package ratelimit

import (
	"context"
	"fmt"

	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"

	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/go-multierror"
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
	_ syncer.TranslatorSyncerExtension = new(translatorSyncerExtension)
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

var (
	emptyTypedResources = map[string]envoycache.Resources{
		v1.RateLimitConfigType: {
			Version: emptyVersionKey,
			Items:   map[string]envoycache.Resource{},
		},
	}
)

const (
	emptyVersionKey = "empty"
)

func init() {
	_ = view.Register(rlConnectedStateView)
}

type translatorSyncerExtension struct {
	collectorFactory collectors.ConfigCollectorFactory
	domainGenerator  rate_limiter_shims.RateLimitDomainGenerator
	hasher           func(resources []envoycache.Resource) uint64
}

func NewTranslatorSyncerExtension(_ context.Context, params syncer.TranslatorSyncerExtensionParams) syncer.TranslatorSyncerExtension {
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
		params.Hasher,
	)
}

func NewTranslatorSyncer(
	collectorFactory collectors.ConfigCollectorFactory,
	domainGenerator rate_limiter_shims.RateLimitDomainGenerator,
	hasher func(resources []envoycache.Resource) uint64,
) syncer.TranslatorSyncerExtension {
	return &translatorSyncerExtension{
		collectorFactory: collectorFactory,
		domainGenerator:  domainGenerator,
		hasher:           hasher,
	}
}

func (s *translatorSyncerExtension) ID() string {
	return xds.ServerRole
}

func (s *translatorSyncerExtension) Sync(
	ctx context.Context,
	snap *gloov1snap.ApiSnapshot,
	_ *gloov1.Settings,
	snapshotSetter syncer.SnapshotSetter,
	reports reporter.ResourceReports,
) {
	ctx = contextutils.WithLogger(ctx, "rateLimitTranslatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	snapHash := hashutils.MustHash(snap)
	logger.Infof("begin rate limit sync %v (%v proxies, %v rate limit configs)", snapHash, len(snap.Proxies), len(snap.Ratelimitconfigs))
	defer logger.Infof("end sync %v", snapHash)

	reports.Accept(snap.Proxies.AsInputResources()...)
	reports.Accept(snap.Ratelimitconfigs.AsInputResources()...)

	configCollectors := newCollectorSet(s.collectorFactory, snap, reports, logger)

	for _, proxy := range snap.Proxies {
		for _, listener := range proxy.Listeners {
			virtualHosts := glooutils.GetVirtualHostsForListener(listener)

			for _, virtualHost := range virtualHosts {

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
		if _, err := s.domainGenerator.NewRateLimitDomain(ctx, cfg.Domain, cfg.Domain,
			&solo_api_rl.RateLimitConfigSpec_Raw{Descriptors: cfg.Descriptors, SetDescriptors: cfg.SetDescriptors}); err != nil {
			errs = multierror.Append(errs, err)
		}
		snapshotResources = append(snapshotResources, v1.NewRateLimitConfigXdsResourceWrapper(cfg))
	}
	if err := errs.ErrorOrNil(); err != nil {
		// This means that one or more rate limit configs could not be translated into the appropriate xDS format
		// Historically, we would error the syncer here
		// In the future, we should assign these errors to the status of the individual resource
		// For now, we will loudly log the error and continue
		logger.Warnf("Encountered errors during sync: %+v", err)
	}

	var rateLimitSnapshot envoycache.Snapshot
	if snapshotResources == nil {
		// If there are no rate limit configs, use an empty configuration
		//
		// The SnapshotCache can now differentiate between nil and empty resources in a snapshot.
		// This was introduced with: https://github.com/solo-io/solo-kit/pull/410
		// A nil resource is not updated, whereas an empty resource is intended to be modified.
		//
		// The ratelimit service only becomes healthy after it has received configuration
		// from Gloo via xDS. Therefore, we must set the ratelimit config resource to empty in the snapshot
		// so that ratelimiit picks up the empty config, and becomes healthy
		rateLimitSnapshot = envoycache.NewGenericSnapshot(emptyTypedResources)
	} else {
		snapshotVersion := fmt.Sprintf("%d", s.hasher(snapshotResources))
		rateLimitSnapshot = envoycache.NewEasyGenericSnapshot(snapshotVersion, snapshotResources)
	}

	snapshotSetter.SetSnapshot(s.ID(), rateLimitSnapshot)
	stats.Record(ctx, rlConnectedState.M(int64(1)))
}

// Helper object to reduce boilerplate.
type configCollectorSet struct {
	collectors []collectors.ConfigCollector
}

func newCollectorSet(
	collectorFactory collectors.ConfigCollectorFactory,
	snapshot *gloov1snap.ApiSnapshot,
	reports reporter.ResourceReports,
	logger *zap.SugaredLogger,
) configCollectorSet {
	set := configCollectorSet{}

	for _, collectorType := range []collectors.CollectorType{
		collectors.Global,
		collectors.Basic,
		collectors.Crd,
	} {
		collector := collectorFactory.MakeInstance(collectorType, snapshot, logger)
		set.collectors = append(set.collectors, collector)
	}
	return set
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
