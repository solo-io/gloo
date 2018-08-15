package errutils

import (
	"context"
)

func AggregateErrs(ctx context.Context, dest chan error, src <-chan error) {
	for {
		select {
		case err := <-src:
			dest <- err
		case <-ctx.Done():
			return
		}
	}
}
