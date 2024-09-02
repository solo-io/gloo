package proxy_syncer

import (
	"context"
	"sync"
	"time"

	"github.com/rotisserie/eris"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
)

type statusSyncer struct {
	// shared with translator syncer; no data race because we own the reporter.
	// if translator syncer starts doing writes with the reporter, we should add locks
	reporter reporter.StatusReporter

	syncNeeded          chan struct{}
	identity            leaderelector.Identity
	leaderStartupAction *leaderelector.LeaderStartupAction
	reportsLock         sync.RWMutex
	latestReports       reporter.ResourceReports
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

// syncExtensions executes each of the TranslatorSyncerExtensions
// These are responsible for updating xDS cache entries
func (s *ProxyTranslator) syncExtensions(
	ctx context.Context,
	snap *v1snap.ApiSnapshot,
	reports reporter.ResourceReports,
) {
	for _, syncerExtension := range s.syncerExtensions {
		intermediateReports := make(reporter.ResourceReports)
		syncerExtension.Sync(ctx, snap, s.settings, s.xdsCache, intermediateReports)
		reports.Merge(intermediateReports)
	}
}

func (s *statusSyncer) syncStatusOnEmit(ctx context.Context) {
	var retryChan <-chan time.Time

	doSync := func() {
		err := s.syncStatus(ctx)
		if err != nil {
			contextutils.LoggerFrom(ctx).Debugw("failed to sync status; will try again shortly.", "error", err)
			retryChan = time.After(time.Second)
		} else {
			retryChan = nil
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-retryChan:
			doSync()
		case <-s.syncNeeded:
			doSync()
		}
	}
}

func (s *statusSyncer) forceSync() {
	if len(s.syncNeeded) > 0 {
		// sync is already needed; no reason to block on send
		return
	}
	s.syncNeeded <- struct{}{}
}

func (s *statusSyncer) syncStatus(ctx context.Context) error {
	s.reportsLock.RLock()
	// deep copy the reports so we can release the lock
	reports := make(reporter.ResourceReports, len(s.latestReports))
	for k, v := range s.latestReports {
		reports[k] = v
	}
	s.reportsLock.RUnlock()

	if len(reports) == 0 {
		return nil
	}

	logger := contextutils.LoggerFrom(ctx)
	if s.identity.IsLeader() {
		// Only leaders will write reports
		//
		// while tempting to write statuses in parallel to increase performance, we should actually first consider recommending the user tunes k8s qps/burst:
		// https://github.com/solo-io/gloo/blob/a083522af0a4ce22f4d2adf3a02470f782d5a865/projects/gloo/api/v1/settings.proto#L337-L350
		//
		// add TEMPORARY wrap to our WriteReports error that we should remove in Gloo Edge ~v1.16.0+.
		// to get the status performance improvements, we need to make the assumption that the user has the latest CRDs installed.
		// if a user forgets the error message is very confusing (invalid request during kubectl patch);
		// this should help them understand what's going on in case they did not read the changelog.
		if err := s.reporter.WriteReports(ctx, reports, nil); err != nil {
			logger.Debugf("Failed writing report for proxies: %v", err)

			wrappedErr := eris.Wrapf(err, "failed to write reports. "+
				"did you make sure your CRDs have been updated since v1.13.0-beta14 of open-source? (i.e. `status` and `status.statuses` fields exist on your CR)")
			return wrappedErr
		}
	} else {
		logger.Debugf("Not a leader, skipping reports writing")
		s.leaderStartupAction.SetAction(func() error {
			// Store the closure in the StartupAction so that it is invoked if this component becomes the new leader
			// That way we can be sure that statuses are updated even if no changes occur after election completes
			// https://github.com/solo-io/gloo/issues/7148
			return s.reporter.WriteReports(ctx, reports, nil)
		})
	}
	return nil
}
