package syncer

import (
	"context"
	"sync"
	"time"

	"github.com/solo-io/gloo/pkg/utils/statsutils/metrics"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/iosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	gwsyncer "github.com/solo-io/gloo/projects/gateway/pkg/syncer"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/sanitizer"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
)

type translatorSyncer struct {
	translator translator.Translator
	sanitizer  sanitizer.XdsSanitizer
	xdsCache   envoycache.SnapshotCache
	reporter   reporter.StatusReporter // shared with status syncer; no data race because StatusFromReport() only reads values set at reporter initialization

	syncerExtensions []TranslatorSyncerExtension
	settings         *v1.Settings
	statusMetrics    metrics.ConfigStatusMetrics
	gatewaySyncer    *gwsyncer.TranslatorSyncer
	proxyClient      v1.ProxyClient
	writeNamespace   string

	// used for debugging purposes only
	// Deprecated: https://github.com/solo-io/gloo/issues/6494
	// Prefer to use the iosnapshot.History
	latestSnap *v1snap.ApiSnapshot

	// snapshotHistory is used for debugging purposes
	// The syncer updates the History each time it runs, and the History is then used by the Admin Server
	snapshotHistory iosnapshot.History

	statusSyncer *statusSyncer
}

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

func NewTranslatorSyncer(
	ctx context.Context,
	translator translator.Translator,
	xdsCache envoycache.SnapshotCache,
	sanitizer sanitizer.XdsSanitizer,
	reporter reporter.StatusReporter,
	devMode bool,
	extensions []TranslatorSyncerExtension,
	settings *v1.Settings,
	statusMetrics metrics.ConfigStatusMetrics,
	gatewaySyncer *gwsyncer.TranslatorSyncer,
	proxyClient v1.ProxyClient,
	writeNamespace string,
	identity leaderelector.Identity,
	snapshotHistory iosnapshot.History,
) v1snap.ApiSyncer {
	s := &translatorSyncer{
		translator:       translator,
		xdsCache:         xdsCache,
		reporter:         reporter,
		syncerExtensions: extensions,
		sanitizer:        sanitizer,
		settings:         settings,
		statusMetrics:    statusMetrics,
		gatewaySyncer:    gatewaySyncer,
		proxyClient:      proxyClient,
		writeNamespace:   writeNamespace,
		statusSyncer: &statusSyncer{
			reporter:            reporter,
			syncNeeded:          make(chan struct{}, 1),
			identity:            identity,
			leaderStartupAction: leaderelector.NewLeaderStartupAction(identity),
			reportsLock:         sync.RWMutex{},
		},
		snapshotHistory: snapshotHistory,
	}
	if devMode {
		// TODO(ilackarms): move this somewhere else?
		go func() {
			_ = s.ContextuallyServeXdsSnapshots(ctx)
		}()
	}
	go s.statusSyncer.syncStatusOnEmit(ctx)
	s.statusSyncer.leaderStartupAction.WatchElectionResults(ctx)
	return s
}

func (s *translatorSyncer) Sync(ctx context.Context, snap *v1snap.ApiSnapshot) error {
	logger := contextutils.LoggerFrom(ctx)
	var multiErr *multierror.Error

	// Log sync start with snapshot summary
	logger.Infow("Starting translator sync",
		"issue", "8539",
		"proxies", len(snap.Proxies),
		"upstreams", len(snap.Upstreams),
		"endpoints", len(snap.Endpoints),
		"secrets", len(snap.Secrets),
		"authConfigs", len(snap.AuthConfigs),
		"rateLimitConfigs", len(snap.Ratelimitconfigs),
		"virtualServices", len(snap.VirtualServices),
		"gateways", len(snap.Gateways),
		"routeTables", len(snap.RouteTables))

	// If gateway controller is enabled, run the gateway translation to generate proxies.
	// Use the ProxyClient interface to persist them either to an in-memory store or etcd as configured at startup.
	if s.gatewaySyncer != nil {
		logger.Infow("Starting gateway proxy translation",
			"issue", "8539",
			"gateways", len(snap.Gateways))
		if err := s.translateProxies(ctx, snap); err != nil {
			logger.Infow("Gateway proxy translation failed",
				"issue", "8539",
				"error", err.Error())
			multiErr = multierror.Append(multiErr, eris.Wrapf(err, "translating proxies"))
		} else {
			logger.Infow("Gateway proxy translation completed",
				"issue", "8539",
				"proxies", len(snap.Proxies))
		}
	}

	// Reports used to aggregate results from xds and extension translation.
	// Will contain reports for `Gloo` components (i.e. Proxies, Upstreams, AuthConfigs, etc.)
	reports := make(reporter.ResourceReports)

	// Execute the EnvoySyncer
	// This will update the xDS SnapshotCache for each entry that corresponds to a non-kube gw Proxy in the API Snapshot
	logger.Infow("Starting Envoy xDS sync",
		"issue", "8539",
		"proxies", len(snap.Proxies))
	s.syncEnvoy(ctx, snap, reports)
	logger.Infow("Envoy xDS sync completed",
		"issue", "8539")

	// Execute the SyncerExtensions
	// Each of these are responsible for updating a single entry in the SnapshotCache
	logger.Infow("Starting syncer extensions",
		"issue", "8539",
		"extensionCount", len(s.syncerExtensions))
	s.syncExtensions(ctx, snap, reports)
	logger.Infow("Syncer extensions completed",
		"issue", "8539")

	// Update resource status metrics and filter out kube gateway proxies
	filteredReports := make(reporter.ResourceReports)
	for resource, report := range reports {
		if proxy, ok := resource.(*v1.Proxy); ok {
			// if this is a proxy report for kube gw, skip it
			if utils.GetTranslatorValue(proxy.GetMetadata()) == utils.GatewayApiProxyValue {
				continue
			}
		}
		filteredReports[resource] = report
		status := s.reporter.StatusFromReport(report, nil)
		s.statusMetrics.SetResourceStatus(ctx, resource, status)
	}

	// After reports are written for proxies, save in gateway syncer (previously gw watched for status changes to proxies)
	if s.gatewaySyncer != nil {
		logger.Infow("Updating gateway proxy status",
			"issue", "8539")
		s.gatewaySyncer.UpdateStatusForAllProxies(ctx)
	}

	// Update status syncer with latest reports
	logger.Infow("Updating status reports",
		"issue", "8539",
		"reportCount", len(filteredReports))
	s.statusSyncer.reportsLock.Lock()
	s.statusSyncer.latestReports = filteredReports
	s.statusSyncer.reportsLock.Unlock()
	s.statusSyncer.forceSync()

	if multiErr.ErrorOrNil() != nil {
		logger.Infow("Translator sync completed with errors",
			"issue", "8539",
			"errorCount", len(multiErr.Errors))
	} else {
		logger.Infow("Translator sync completed successfully",
			"issue", "8539")
	}

	return multiErr.ErrorOrNil()
}

// syncExtensions executes each of the TranslatorSyncerExtensions
// These are responsible for updating xDS cache entries
func (s *translatorSyncer) syncExtensions(ctx context.Context, snap *v1snap.ApiSnapshot, reports reporter.ResourceReports) {
	for _, syncerExtension := range s.syncerExtensions {
		intermediateReports := make(reporter.ResourceReports)
		syncerExtension.Sync(ctx, snap, s.settings, s.xdsCache, intermediateReports)
		reports.Merge(intermediateReports)
	}
}

// translateProxies will call the gatewaySyncer to translate Proxies for the Gateways in the provided snapshot.
// It will then use the proxyClient to List() Proxies and *mutate the snapshot* to add those Proxies.
func (s *translatorSyncer) translateProxies(ctx context.Context, snap *v1snap.ApiSnapshot) error {
	logger := contextutils.LoggerFrom(ctx)
	var multiErr *multierror.Error

	logger.Infow("Starting gateway syncer translation",
		"issue", "8539",
		"gateways", len(snap.Gateways),
		"writeNamespace", s.writeNamespace)

	err := s.gatewaySyncer.Sync(ctx, snap)
	if err != nil {
		logger.Infow("Gateway syncer translation failed",
			"issue", "8539",
			"error", err.Error())
		multiErr = multierror.Append(multiErr, err)
	} else {
		logger.Infow("Gateway syncer translation completed successfully",
			"issue", "8539")
	}

	logger.Infow("Listing proxies from client",
		"issue", "8539",
		"namespace", s.writeNamespace)

	proxyList, err := s.proxyClient.List(s.writeNamespace, clients.ListOpts{})
	if err != nil {
		logger.Infow("Failed to list proxies from client",
			"issue", "8539",
			"error", err.Error(),
			"namespace", s.writeNamespace)
		multiErr = multierror.Append(multiErr, err)
	} else {
		logger.Infow("Successfully retrieved proxies from client",
			"issue", "8539",
			"proxyCount", len(proxyList),
			"namespace", s.writeNamespace)
	}

	snap.Proxies = proxyList
	return multiErr.ErrorOrNil()
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
	logger.Infow("Starting status sync",
		"issue", "8539",
		"reportCount", len(reports),
		"isLeader", s.identity.IsLeader())

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
			logger.Infow("Failed to write status reports",
				"issue", "8539",
				"error", err.Error(),
				"reportCount", len(reports))

			wrappedErr := eris.Wrapf(err, "failed to write reports. "+
				"did you make sure your CRDs have been updated since v1.13.0-beta14 of open-source? (i.e. `status` and `status.statuses` fields exist on your CR)")
			return wrappedErr
		} else {
			logger.Infow("Successfully wrote status reports",
				"issue", "8539",
				"reportCount", len(reports))
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
