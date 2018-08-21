package errutils

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/errors"
)

func AggregateErrs(ctx context.Context, dest chan error, src <-chan error, srcInfo string) {
	for {
		select {
		case err, ok := <-src:
			if !ok {
				return
			}
			dest <- errors.Wrapf(err, srcInfo)
		case <-ctx.Done():
			return
		}
	}
}
