package settings

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"go.uber.org/zap"
)

const DefaultRefreshRate = 10 * time.Minute

// A client for accessing specific values from a user's settings.
// Returns reasonable defaults if valid values aren't found in the settings CRD.
type ValuesClient interface {
	GetRefreshRate() time.Duration
}

type client struct {
	ctx            context.Context
	settingsClient v1.SettingsClient
}

func (c *client) GetRefreshRate() time.Duration {
	var refreshRate time.Duration
	// TODO cleanup hardcoded resource namespace / name
	settings, err := c.settingsClient.Read("gloo-system", "default", clients.ReadOpts{Ctx: c.ctx})
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

func NewSettingsValuesClient(ctx context.Context, settingsClient v1.SettingsClient) ValuesClient {
	return &client{
		ctx:            ctx,
		settingsClient: settingsClient,
	}
}
