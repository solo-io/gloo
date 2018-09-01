package syncer

import (
	"context"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/xds"
)

type syncer struct {
	translator translator.Translator
	xdsCache   envoycache.SnapshotCache
	xdsHasher  *xds.ProxyKeyHasher
	reporter   reporter.Reporter
}

func NewSyncer(translator translator.Translator, xdsCache envoycache.SnapshotCache, xdsHasher *xds.ProxyKeyHasher, reporter reporter.Reporter) v1.Syncer {
	return &syncer{
		translator: translator,
		xdsCache:   xdsCache,
		xdsHasher:  xdsHasher,
		reporter:   reporter,
	}
}

func (s *syncer) Sync(ctx context.Context, snap *v1.Snapshot) error {
	ctx = contextutils.WithLogger(ctx, "gloo.syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("Beginning translation loop for snapshot %v", snap.Hash())
	logger.Debugf("%v", snap)
	allResourceErrs := make(reporter.ResourceErrors)
	allResourceErrs.Initialize(snap.Upstreams.List().AsInputResources()...)
	allResourceErrs.Initialize(snap.Proxies.List().AsInputResources()...)

	params := plugins.Params{
		Ctx:      ctx,
		Snapshot: snap,
	}

	s.xdsHasher.SetKeysFromProxies(snap.Proxies.List())

	for _, proxy := range snap.Proxies.List() {
		xdsSnapshot, resourceErrs := s.translator.Translate(params, proxy)

		allResourceErrs.Merge(resourceErrs)

		if err := resourceErrs.Validate(); err != nil {
			logger.Warnf("proxy %v was rejected due to invalid config: %v\nxDS cache will not be updated.", err)
			continue
		}
		if err := s.xdsCache.SetSnapshot(proxy.Metadata.Name, xdsSnapshot); err != nil {
			return errors.Wrapf(err, "failed while updating xds snapshot cache")
		}
		logger.Infof("Setting xDS Snapshot for Proxy %v: %v clusters, %v listeners, %v route configs, %v endpoints",
			proxy.Metadata.Name, len(xdsSnapshot.Clusters.Items), len(xdsSnapshot.Listeners.Items),
			len(xdsSnapshot.Routes.Items), len(xdsSnapshot.Endpoints.Items))

		logger.Debugf("Full snapshot for proxy %v: %v", proxy.Metadata.Name, xdsSnapshot)
	}
	if err := s.reporter.WriteReports(ctx, allResourceErrs); err != nil {
		return errors.Wrapf(err, "writing reports")
	}
	logger.Infof("finished translation loop %v", snap.Hash())
	return nil
}
