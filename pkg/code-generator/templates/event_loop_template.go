package templates

import (
	"text/template"
)

var ResourceGroupEventLoopTemplate = template.Must(template.New("resource_group_event_loop").Funcs(funcs).Parse(`package {{ .Project.PackageName }}

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

type {{ .GoName }}EventLoop interface {
	Run(namespaces []string, opts clients.WatchOpts) (<-chan error, error)
}

type {{ lower_camel .GoName }}EventLoop struct {
	cache  Cache
	syncer Syncer
}

func New{{ .GoName }}EventLoop(cache Cache, syncer Syncer) {{ .GoName }}EventLoop {
	return &{{ lower_camel .GoName }}EventLoop{
		cache:  cache,
		syncer: syncer,
	}
}

func (el *{{ lower_camel .GoName }}EventLoop) Run(namespaces []string, opts clients.WatchOpts) (<-chan error, error) {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "{{ .Project.PackageName }}.event_loop")
	logger := contextutils.LoggerFrom(opts.Ctx)
	logger.Infof("event loop started")

	errs := make(chan error)

	watch, cacheErrs, err := el.cache.Snapshots(namespaces, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "starting snapshot watch")
	}
	go errutils.AggregateErrs(opts.Ctx, errs, cacheErrs, "{{ .Project.PackageName }}.cache errors")
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
`))
