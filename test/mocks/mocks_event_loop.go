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
	Run(namespace string, opts clients.WatchOpts) error
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

func (el *eventLoop) Run(namespace string, opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "mocks.event_loop")
	logger := contextutils.LoggerFrom(opts.Ctx)
	logger.Infof("event loop started")
	errorHandler := contextutils.ErrorHandlerFrom(opts.Ctx)
	watch, errs, err := el.cache.Snapshots(namespace, opts)
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
