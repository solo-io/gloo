package settings

import (
	"context"
	"time"

	clientscache "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/go-utils/contextutils"
	v1clients "github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"go.uber.org/zap"
)

const DefaultRefreshRate = 10 * time.Minute

// A client for accessing specific values from a user's settings.
// Returns reasonable defaults if valid values aren't found in the settings CRD.
type ValuesClient interface {
	GetRefreshRate() time.Duration
}

type client struct {
	ctx         context.Context
	clientCache clientscache.ClientCache
}

func (c *client) GetRefreshRate() time.Duration {
	var refreshRate time.Duration
	// TODO cleanup hardcoded resource namespace / name

	settings, err := c.clientCache.GetSettingsClient().Read("gloo-system", "default", v1clients.ReadOpts{Ctx: c.ctx})
	if err != nil {
		contextutils.LoggerFrom(c.ctx).Errorw("Failed to read settings",
			zap.Error(err))
	}

	if settings.GetRefreshRate() != nil {
		refreshRate, err = types.DurationFromProto(settings.GetRefreshRate())
		if err != nil {
			contextutils.LoggerFrom(c.ctx).Errorw("Invalid refresh rate stored in settings", zap.Error(err))
		}
	}

	if refreshRate == 0 {
		contextutils.LoggerFrom(c.ctx).Infow("Falling back to default refresh rate")
		return DefaultRefreshRate
	}

	return refreshRate
}

func NewSettingsValuesClient(ctx context.Context, clientCache clientscache.ClientCache) ValuesClient {
	return &client{
		ctx:         ctx,
		clientCache: clientCache,
	}
}
