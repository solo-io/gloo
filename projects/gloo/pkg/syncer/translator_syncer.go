package syncer

import (
	"context"

	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/projects/gloo/pkg/control-plane/cache"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/xds"
)

var (
	proxyNameKey, _ = tag.NewKey("proxyname")
)

type translatorSyncer struct {
	translator translator.Translator
	xdsCache   envoycache.SnapshotCache
	xdsHasher  *xds.ProxyKeyHasher
	reporter   reporter.Reporter
	// used for debugging purposes only
	latestSnap *v1.ApiSnapshot
}

func NewTranslatorSyncer(translator translator.Translator, xdsCache envoycache.SnapshotCache, xdsHasher *xds.ProxyKeyHasher, reporter reporter.Reporter, devMode bool) v1.ApiSyncer {
	s := &translatorSyncer{
		translator: translator,
		xdsCache:   xdsCache,
		xdsHasher:  xdsHasher,
		reporter:   reporter,
	}
	if devMode {
		// TODO(ilackarms): move this somewhere else?
		go s.ServeXdsSnapshots()
	}
	return s
}

func (s *translatorSyncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {
	err := s.syncEnvoy(ctx, snap)
	if err != nil {
		return err
	}

	err = s.syncRateLimit(ctx, snap)
	if err != nil {
		return err
	}
	return nil
}
