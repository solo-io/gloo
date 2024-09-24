package proxy_syncer

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	syncerstats "github.com/solo-io/gloo/projects/gloo/pkg/syncer/stats"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	envoySnapshotOut   = stats.Int64("api.gloo.solo.io/translator/resources", "The number of resources in the snapshot in", "1")
	resourceNameKey, _ = tag.NewKey("resource")

	envoySnapshotOutView = &view.View{
		Name:        "api.gloo.solo.io/translator/resources",
		Measure:     envoySnapshotOut,
		Description: "The number of resources in the snapshot for envoy",
		Aggregation: view.LastValue(),
		TagKeys:     []tag.Key{syncerstats.ProxyNameKey, resourceNameKey},
	}
)

func init() {
	_ = view.Register(envoySnapshotOutView)
}

func measureResource(ctx context.Context, resource string, length int) {
	if ctxWithTags, err := tag.New(ctx, tag.Insert(resourceNameKey, resource)); err == nil {
		stats.Record(ctxWithTags, envoySnapshotOut.M(int64(length)))
	}
}

func (s *ProxyTranslator) glooSync(ctx context.Context, snap *v1snap.ApiSnapshot) []translatorutils.ProxyWithReports {
	// Reports used to aggregate results from xds and extension translation.
	// Will contain reports only `Gloo` components (i.e. Proxies, Upstreams, AuthConfigs, etc.)
	reports := make(reporter.ResourceReports)

	contextutils.LoggerFrom(ctx).Info("before gw sync envoy")
	// Execute the EnvoySyncer
	// This will update the xDS SnapshotCache for each entry that corresponds to a Proxy in the API Snapshot
	// TODO: need to pass in ggv2 proxies now
	proxyReports := s.syncEnvoy(ctx, snap, reports)
	contextutils.LoggerFrom(ctx).Info("after gw sync envoy")

	// Execute the SyncerExtensions
	// Each of these are responsible for updating a single entry in the SnapshotCache
	s.syncExtensions(ctx, snap, reports)

	// reports now has been merged from the envoy and extension translation/syncs
	// it also contains reports for all Gloo resources (Upstreams, Proxies, AuthConfigs, RLCs, etc.)
	// so let's filter out non-Proxy reports
	// TODO: we actually don't want to do this, we do need to report status
	filteredReports := reports.FilterByKind("Proxy")

	// build object used by status plugins
	var proxiesWithReports []translatorutils.ProxyWithReports
	for i, proxy := range snap.Proxies {
		proxy := proxy // still need pike?

		// build ResourceReports struct containing only this Proxy
		r := make(reporter.ResourceReports)
		r[proxy] = filteredReports[proxy]

		proxiesWithReports = append(proxiesWithReports, translatorutils.ProxyWithReports{
			Proxy: proxy,
			Reports: translatorutils.TranslationReports{
				ProxyReport:     proxyReports[i],
				ResourceReports: r,
			},
		})
	}

	// TODO(Law): confirm not needed; metrics can be derived from k8s conditions, may be needed for Policy GE-style status?
	// // Update resource status metrics
	// for resource, report := range reports {
	// 	status := s.reporter.StatusFromReport(report, nil)
	// 	s.statusMetrics.SetResourceStatus(ctx, resource, status)
	// }

	// need to write proxy reports

	contextutils.LoggerFrom(ctx).Info("LAW got proxieswithreports")
	return proxiesWithReports
}

// syncEnvoy will translate, sanitize, and set the xds snapshot for each of the proxies in the provided api snapshot.
// Reports from translation attempts on every Proxy will be merged into allReports.
func (s *ProxyTranslator) syncEnvoy(
	ctx context.Context,
	snap *v1snap.ApiSnapshot,
	allReports reporter.ResourceReports,
) []*validation.ProxyReport {
	ctx, span := trace.StartSpan(ctx, "gloo.kube.syncer.Sync")
	defer span.End()

	ctx = contextutils.WithLogger(ctx, "kube-gateway-envoy-translatorSyncer")
	logger := contextutils.LoggerFrom(ctx)
	snapHash := hashutils.MustHash(snap)
	logger.Infof("begin kube gw sync %v (%v proxies, %v upstreams, %v endpoints, %v secrets, %v artifacts, %v auth configs, %v rate limit configs, %v graphql apis)", snapHash,
		len(snap.Proxies), len(snap.Upstreams), len(snap.Endpoints), len(snap.Secrets), len(snap.Artifacts), len(snap.AuthConfigs), len(snap.Ratelimitconfigs), len(snap.GraphqlApis))
	defer logger.Infof("end sync %v", snapHash)

	// stringifying the snapshot may be an expensive operation, so we'd like to avoid building the large
	// string if we're not even going to log it anyway
	if contextutils.GetLogLevel() == zapcore.DebugLevel {
		// logger.Debug(syncutil.StringifySnapshot(snap))
	}

	allReports.Accept(snap.Proxies.AsInputResources()...)
	// accept Upstream[Group]s as they may be reported on during xds translation, but we will drop them.
	// the main GE translator_syncer will manage them, it is not our responsibility
	allReports.Accept(snap.Upstreams.AsInputResources()...)
	allReports.Accept(snap.UpstreamGroups.AsInputResources()...)

	// parallel slice to snap.Proxies containing corresponding proxyReport to translation
	var proxyValidationReports []*validation.ProxyReport
	for _, proxy := range snap.Proxies {
		proxyCtx := ctx
		metaKey := xds.SnapshotCacheKey(proxy)
		if ctxWithTags, err := tag.New(proxyCtx, tag.Insert(syncerstats.ProxyNameKey, metaKey)); err == nil {
			proxyCtx = ctxWithTags
		}

		params := plugins.Params{
			Ctx:      proxyCtx,
			Settings: s.settings,
			Snapshot: snap,
			Messages: map[*core.ResourceRef][]string{},
		}

		xdsSnapshot, reports, proxyReport := s.translator.Translate(params, proxy)
		proxyValidationReports = append(proxyValidationReports, proxyReport)

		// Messages are aggregated during translation, and need to be added to reports
		for _, messages := range params.Messages {
			reports.AddMessages(proxy, messages...)
		}

		if validateErr := reports.ValidateStrict(); validateErr != nil {
			logger.Warnw("Proxy had invalid config", zap.Any("proxy", proxy.GetMetadata().Ref()), zap.Error(validateErr))
		}

		allReports.Merge(reports)

		key := xds.SnapshotCacheKey(proxy)
		// if the snapshot is not consistent, make it so
		xdsSnapshot.MakeConsistent()
		s.xdsCache.SetSnapshot(key, xdsSnapshot)

		// Record some metrics
		clustersLen := len(xdsSnapshot.GetResources(types.ClusterTypeV3).Items)
		listenersLen := len(xdsSnapshot.GetResources(types.ListenerTypeV3).Items)
		routesLen := len(xdsSnapshot.GetResources(types.RouteTypeV3).Items)
		endpointsLen := len(xdsSnapshot.GetResources(types.EndpointTypeV3).Items)

		measureResource(proxyCtx, "clusters", clustersLen)
		measureResource(proxyCtx, "listeners", listenersLen)
		measureResource(proxyCtx, "routes", routesLen)
		measureResource(proxyCtx, "endpoints", endpointsLen)

		logger.Infow("Setting xDS Snapshot", "key", key,
			"clusters", clustersLen,
			"listeners", listenersLen,
			"routes", routesLen,
			"endpoints", endpointsLen)

		logger.Debugf("Full snapshot for proxy %v: %+v", proxy.GetMetadata().GetName(), xdsSnapshot)
	}

	logger.Debugf("gloo reports to be written: %v", allReports)
	s.syncStatus(ctx, allReports)
	return proxyValidationReports
}
func (s *ProxyTranslator) syncStatus(ctx context.Context, reports reporter.ResourceReports) error {
	// leftover from translator_syncer's statusSyncer
	// analyze our plan for concurrency, data ownership, do we need locks, etc.?

	// s.reportsLock.RLock()
	// // deep copy the reports so we can release the lock
	// reports := make(reporter.ResourceReports, len(s.latestReports))
	// for k, v := range s.latestReports {
	// 	reports[k] = v
	// }
	// s.reportsLock.RUnlock()

	logger := contextutils.LoggerFrom(ctx)

	// if s.identity.IsLeader() {
	if err := s.glooReporter.WriteReports(ctx, reports, nil); err != nil {
		logger.Debugf("Failed writing report for proxies: %v", err)
		return err
	}
	return nil
}

// syncExtensions executes each of the TranslatorSyncerExtensions
// These are responsible for updating xDS cache entries
func (s *ProxyTranslator) syncExtensions(
	ctx context.Context,
	snap *v1snap.ApiSnapshot,
	reports reporter.ResourceReports,
) {
	for _, syncerExtension := range s.syncerExtensions {
		intermediateReports := make(reporter.ResourceReports)
		// we use the no-op setter here as we don't actually sync the extensions here,
		// that is classic edge syncer's job [see: projects/gloo/pkg/syncer/translator_syncer.go#Sync(...)]
		// all we care about is getting the reports, as our `Proxies` will get reports for errors/warns
		// related to the extension processing
		syncerExtension.Sync(ctx, snap, s.settings, s.noopSnapSetter, intermediateReports)
		reports.Merge(intermediateReports)
	}
}
