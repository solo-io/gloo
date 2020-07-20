package runner

import (
	"context"

	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/shims"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/xds"
)

func Run(ctx context.Context, serverSettings server.Settings, xdsSettings xds.Settings) {
	xDSRunner := xds.NewConfigSource(xdsSettings, shims.NewRateLimitDomainGenerator())
	rateLimitServer := server.NewRateLimitServer(serverSettings, xDSRunner)

	// Blocks until shutdown signals are received. Normal shutdown will return nil error.
	if err := rateLimitServer.Start(ctx); err != nil {
		log.Fatalf("server stopped with unexpected error: %s", err.Error())
	}
}
