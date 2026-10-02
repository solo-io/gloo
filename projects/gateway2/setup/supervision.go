package setup

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

// RunWithGateway runs the Gateway API producer alongside the legacy setup loop.
// An unexpected producer exit is fatal to this process: otherwise setup can
// wait forever for a Proxy snapshot that no controller will ever publish.
// Both functions must stop when their context is canceled. A nil gateway keeps
// the legacy-only startup path synchronous.
func RunWithGateway(ctx context.Context, gateway, setup func(context.Context) error) error {
	if gateway == nil {
		return setup(ctx)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	group, runCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		err := gateway(runCtx)
		if runCtx.Err() != nil {
			return nil
		}
		if err == nil {
			return fmt.Errorf("Gateway API controller exited unexpectedly")
		}
		return fmt.Errorf("Gateway API controller failed: %w", err)
	})
	group.Go(func() error {
		err := setup(runCtx)
		if err == nil {
			cancel()
		}
		return err
	})
	return group.Wait()
}
