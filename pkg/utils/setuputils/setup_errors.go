package setuputils

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
)

// The event loop does not close its error channel on cancellation. Returning
// errors lets the process supervisor cancel the other controller before exit.
func waitForSetupErrors(ctx context.Context, errs <-chan error, exitOnError bool) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-errs:
			if !ok {
				return nil
			}
			if exitOnError {
				return err
			}
			contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
		}
	}
}
