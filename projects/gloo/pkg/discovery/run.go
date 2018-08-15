package discovery

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

// run once then ya done
func RunUds(disc *Discovery, opts clients.WatchOpts, discOpts Opts) (chan error, error) {
	errs, err := disc.StartUds(opts, discOpts)
	if err != nil {
		return nil, err
	}
	return errs, nil
}

// run once, watch upstreams
func RunEds(upstreamClient v1.UpstreamClient, disc *Discovery, watchNamespace string, opts clients.WatchOpts) (chan error, error) {
	errs := make(chan error)
	upstreams, upstreamErrs, err := upstreamClient.Watch(watchNamespace, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "beginning upstream watch")
	}
	var cancel context.CancelFunc
	ctx := opts.Ctx
	go func() {
		for {
			select {
			case err := <-upstreamErrs:
				errs <- err
			case upstreamList := <-upstreams:
				if cancel != nil {
					cancel()
				}
				opts.Ctx, cancel = context.WithCancel(ctx)

				edsErrs, err := disc.StartEds(upstreamList, opts)
				if err != nil {
					errs <- err
					continue
				}
				go errutils.AggregateErrs(opts.Ctx, errs, edsErrs)

			}
		}
	}()
	return errs, nil
}
