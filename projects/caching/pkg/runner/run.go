package runner

import (
	"context"

	"github.com/solo-io/caching-service/pkg/runner"
	"github.com/solo-io/caching-service/pkg/settings"
	"github.com/solo-io/go-utils/log"
)

// Run the caching service
func Run(ctx context.Context, serverSettings settings.Settings) {

	err := runner.StartService(ctx, serverSettings)
	if err != nil {
		log.Fatalf("server stopped with unexpected error: %s", err.Error())
	}
}
