package channelutils

import (
	"context"
	"time"

	"github.com/solo-io/go-utils/contextutils"
)

func WaitForReady(ctx context.Context, timeout time.Duration, readychans ...<-chan struct{}) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	logger := contextutils.LoggerFrom(ctx)
	logger.Debugw("Starting WaitForReady",
		"issue", "8539",
		"timeout", timeout.String(),
		"channelCount", len(readychans))

	// wait for everything to finish warming up:
	for i, chans := range readychans {
		logger.Debugw("Waiting for ready channel",
			"issue", "8539",
			"channelIndex", i,
			"totalChannels", len(readychans))

		select {
		case <-chans:
			logger.Debugw("Ready channel signaled",
				"issue", "8539",
				"channelIndex", i,
				"totalChannels", len(readychans))
			// wait for component to become ready
		case <-ctx.Done():
			logger.Warnw("WaitForReady context cancelled or timed out",
				"issue", "8539",
				"channelIndex", i,
				"totalChannels", len(readychans),
				"error", ctx.Err().Error())
			return ctx.Err()
		}
	}

	logger.Debugw("All ready channels signaled successfully",
		"issue", "8539",
		"channelCount", len(readychans))
	return nil
}
