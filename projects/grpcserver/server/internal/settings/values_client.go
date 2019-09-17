package settings

import (
	"context"
	"time"

	"github.com/solo-io/gloo/pkg/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	clientscache "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"
	"k8s.io/kubernetes/pkg/apis/core"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/go-utils/contextutils"
	v1clients "github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"go.uber.org/zap"
)

const DefaultRefreshRate = 10 * time.Minute

//go:generate mockgen -destination mocks/mock_values_client.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings ValuesClient

// A client for accessing specific values from a user's settings.
// Returns reasonable defaults if valid values aren't found in the settings CRD.
type ValuesClient interface {
	GetRefreshRate() time.Duration
	GetWatchNamespaces() []string
}

type client struct {
	ctx          context.Context
	clientCache  clientscache.ClientCache
	podNamespace string
}

func (c *client) GetRefreshRate() time.Duration {
	var refreshRate time.Duration
	var err error

	settings := c.getSettings()
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

func (c *client) GetWatchNamespaces() []string {
	settings := c.getSettings()
	if settings.GetWatchNamespaces() != nil {
		return utils.ProcessWatchNamespaces(settings.GetWatchNamespaces(), settings.GetDiscoveryNamespace())
	}

	// Return all namespaces by default, as this is the default when no watch namespaces are specified.
	return []string{core.NamespaceAll}
}

func (c *client) getSettings() *v1.Settings {
	settings, err := c.clientCache.GetSettingsClient().Read(c.podNamespace, defaults.SettingsName, v1clients.ReadOpts{Ctx: c.ctx})
	if err != nil {
		contextutils.LoggerFrom(c.ctx).Errorw("Failed to read settings",
			zap.Error(err))
	}
	return settings
}

func NewSettingsValuesClient(ctx context.Context, clientCache clientscache.ClientCache, podNamespace string) ValuesClient {
	return &client{
		ctx:          ctx,
		clientCache:  clientCache,
		podNamespace: podNamespace,
	}
}
