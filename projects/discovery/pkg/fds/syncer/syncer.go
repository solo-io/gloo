package syncer

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"time"
)

type syncer struct {
	disc        *discovery.UpstreamDiscovery
	refreshRate time.Duration
	discOpts    discovery.Opts
}

func NewSyncer(disc *discovery.UpstreamDiscovery, discOpts discovery.Opts, refreshRate time.Duration) v1.DiscoverySyncer {
	s := &syncer{
		disc:        disc,
		refreshRate: refreshRate,
		discOpts:    discOpts,
	}
	return s
}

func (s *syncer) Sync(ctx context.Context, snap *v1.DiscoverySnapshot) error {
	ctx = contextutils.WithLogger(ctx, "syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v upstreams)", snap.Hash(), len(snap.Upstreams.List()))
	defer logger.Infof("end sync %v", snap.Hash())

	logger.Debugf("%v", snap)
	allResourceErrs := make(reporter.ResourceErrors)
	allResourceErrs.Accept(snap.Upstreams.List().AsInputResources()...)

	opts := clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: s.refreshRate,
	}

	udsErrs, err := s.disc.StartUds(opts, s.discOpts)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case err := <-udsErrs:
				contextutils.LoggerFrom(ctx).Errorf("error in UDS: %v", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
