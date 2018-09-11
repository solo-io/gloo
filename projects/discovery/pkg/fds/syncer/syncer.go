package syncer

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/discovery/pkg/fds"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type syncer struct {
	fd *fds.FunctionDiscovery
}

func NewSyncer(fd *fds.FunctionDiscovery) v1.DiscoverySyncer {
	s := &syncer{
		fd: fd,
	}
	return s
}

func (s *syncer) Sync(ctx context.Context, snap *v1.DiscoverySnapshot) error {
	ctx = contextutils.WithLogger(ctx, "syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v upstreams)", snap.Hash(), len(snap.Upstreams.List()))
	defer logger.Infof("end sync %v", snap.Hash())

	logger.Debugf("%v", snap)

	return s.fd.Update(snap.Upstreams.List(), snap.Secrets.List())

	//opts := clients.WatchOpts{
	//	Ctx:         ctx,
	//	RefreshRate: s.refreshRate,
	//}
	//
	//udsErrs, err := s.disc.StartUds(opts, s.discOpts)
	//if err != nil {
	//	return err
	//}
	//
	//go func() {
	//	for {
	//		select {
	//		case err := <-udsErrs:
	//			contextutils.LoggerFrom(ctx).Errorf("error in UDS: %v", err)
	//		case <-ctx.Done():
	//			return
	//		}
	//	}
	//}()
	//return nil
}
