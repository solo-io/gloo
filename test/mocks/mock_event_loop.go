package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

type Syncer interface {
	Sync(*Snapshot) error
}

type EventLoop interface {
	Run(opts clients.WatchOpts) error
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

func (el *eventLoop) Run(opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	logger := contextutils.GetLogger(opts.Ctx)
	logger = logger.WithPrefix("mocks.event_loop")
	opts.Ctx = contextutils.WithLogger(opts.Ctx, logger)
	logger.Printf(contextutils.LogLevelInfo, "mocks: event loop started")
	errorHandler := contextutils.GetErrorHandler(opts.Ctx)
	watch, errs, err := el.cache.Snapshots(opts)
	if err != nil {
		return errors.Wrapf(err, "starting snapshot watch")
	}
	for {
		select {
		case snapshot := <-watch:
			err := el.syncer.Sync(snapshot)
			if err != nil {
				errorHandler.HandleErr(err)
			}
		case err := <-errs:
			errorHandler.HandleErr(err)
		}
	}
}
