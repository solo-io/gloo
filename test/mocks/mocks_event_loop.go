package mocks

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Syncer interface {
	Sync(context.Context, *Snapshot) error
}

type EventLoop interface {
	Run(namespaces []string, opts clients.WatchOpts) (<-chan error, error)
}

type eventLoop struct {
	cache  Cache
	syncer Syncer
}

func NewEventLoop(cache Cache, syncer Syncer) EventLoop {
	return &eventLoop{
		cache:  cache,
		syncer: syncer,
	}
}

func (el *eventLoop) Run(namespaces []string, opts clients.WatchOpts) (<-chan error, error) {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "mocks.event_loop")
	logger := contextutils.LoggerFrom(opts.Ctx)
	logger.Infof("event loop started")

	errs := make(chan error)

	watch, cacheErrs, err := el.cache.Snapshots(namespaces, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "starting snapshot watch")
	}
	go errutils.AggregateErrs(opts.Ctx, errs, cacheErrs, "mocks.cache errors")
	go func() {
		// create a new context for each loop, cancel it before each loop
		var cancel context.CancelFunc = func() {}
		defer cancel()
		for {
			select {
			case snapshot := <-watch:
				// cancel any open watches from previous loop
				cancel()
				ctx, canc := context.WithCancel(opts.Ctx)
				cancel = canc
				err := el.syncer.Sync(ctx, snapshot)
				if err != nil {
					errs <- err
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()
	return errs, nil
}
