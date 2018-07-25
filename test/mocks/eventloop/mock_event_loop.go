package eventloop

import (
	"github.com/solo-io/solo-kit/test/mocks"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type Syncer interface {
	Sync(*mocks.Snapshot) error
}

type EventLoop interface {
	Run(opts clients.WatchOpts) error
}

type eventLoop struct {
	cache  mocks.Cache
	syncer Syncer
}

func NewEventLoop(cache mocks.Cache, syncer Syncer) EventLoop {
	return &eventLoop{
		cache:  cache,
		syncer: syncer,
	}
}

func (el *eventLoop) Run(opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	logger := contextutils.GetLogger(opts.Ctx)
	logger.Printf(contextutils.LogLevelInfo, "mock event loop started")
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
