package proxy_syncer

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/hashutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	syncerstats "github.com/solo-io/gloo/projects/gloo/pkg/syncer/stats"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
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
	return proxyValidationReports
}
