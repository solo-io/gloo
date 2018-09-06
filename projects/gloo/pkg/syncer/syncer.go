package syncer

import (
	"context"

	"fmt"
	"net/http"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gorilla/mux"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/log"
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
	// used for debugging purposes only
	latestSnap *v1.ApiSnapshot
}

func NewSyncer(translator translator.Translator, xdsCache envoycache.SnapshotCache, xdsHasher *xds.ProxyKeyHasher, reporter reporter.Reporter) v1.ApiSyncer {
	s := &syncer{
		translator: translator,
		xdsCache:   xdsCache,
		xdsHasher:  xdsHasher,
		reporter:   reporter,
	}
	// TODO(ilackarms): move this somewhere else, make it part of dev-mode
	go s.ServeXdsSnapshots()
	return s
}

func (s *syncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {
	// TODO(yuval-k): @ilackarms, why do we get nil here sometimes?
	if snap == nil {
		return nil
	}
	s.latestSnap = snap
	ctx = contextutils.WithLogger(ctx, "syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v resources)", snap.Hash(),
		len(snap.Upstreams)+len(snap.Secrets)+len(snap.Artifacts)+len(snap.Proxies)+len(snap.Endpoints))
	defer logger.Infof("end sync %v", snap.Hash())

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
		xdsSnapshot, resourceErrs, err := s.translator.Translate(params, proxy)
		if err != nil {
			return errors.Wrapf(err, "translation loop failed")
		}

		allResourceErrs.Merge(resourceErrs)

		if err := resourceErrs.Validate(); err != nil {
			logger.Warnf("proxy %v was rejected due to invalid config: %v\nxDS cache will not be updated.", err)
			continue
		}
		key := xds.SnapshotKey(proxy)
		if err := s.xdsCache.SetSnapshot(key, xdsSnapshot); err != nil {
			return errors.Wrapf(err, "failed while updating xds snapshot cache")
		}
		logger.Infof("Setting xDS Snapshot for Key %v: %v clusters, %v listeners, %v route configs, %v endpoints",
			key, len(xdsSnapshot.Clusters.Items), len(xdsSnapshot.Listeners.Items),
			len(xdsSnapshot.Routes.Items), len(xdsSnapshot.Endpoints.Items))

		logger.Debugf("Full snapshot for proxy %v: %v", proxy.Metadata.Name, xdsSnapshot)
	}
	if err := s.reporter.WriteReports(ctx, allResourceErrs); err != nil {
		return errors.Wrapf(err, "writing reports")
	}
	return nil
}

// TODO(ilackarms): move this somewhere else, make it part of dev-mode
func (s *syncer) ServeXdsSnapshots() error {
	r := mux.NewRouter()
	r.HandleFunc("/xds", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, log.Sprintf("%v", s.xdsCache))
	})
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, log.Sprintf("%v", s.latestSnap))
	})
	return http.ListenAndServe(":9090", r)
}
