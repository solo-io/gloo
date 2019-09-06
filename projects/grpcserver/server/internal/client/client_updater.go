package client

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

//go:generate mockgen -destination mocks/mock_client_updater.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client Updater

// Run a settings watch loop in a new goroutine, rebuilding all the k8s resource clients when the settings are updated
type Updater interface {
	StartWatch() error
}

const clientUpdaterLogger = "clientUpdater"

type updater struct {
	clientCache ClientCache
	token       setup.Token
	cfg         *rest.Config
}

var _ Updater = &updater{}

func NewClientUpdater(clientCache ClientCache, cfg *rest.Config, token setup.Token) Updater {
	return &updater{
		clientCache: clientCache,
		cfg:         cfg,
		token:       token,
	}
}

func (u *updater) StartWatch() error {
	return setuputils.Main(setuputils.SetupOpts{
		LoggingPrefix: clientUpdaterLogger,
		SetupFunc:     u.buildReceiverFunc(),
		ExitOnError:   false,
	})
}

func (u *updater) buildReceiverFunc() setuputils.SetupFunc {
	return func(ctx context.Context, kubeCache kube.SharedCache, inMemoryCache memory.InMemoryResourceCache, settings *gloov1.Settings) error {

		// build a new cache from scratch, and update the existing one from the new one's state
		newCache, err := NewClientCache(ctx, settings, u.cfg, u.token)
		if err != nil {
			return err
		}

		u.clientCache.SetCacheState(
			newCache.GetVirtualServiceClient(),
			newCache.GetGatewayClient(),
			newCache.GetUpstreamClient(),
			newCache.GetSettingsClient(),
			newCache.GetSecretClient(),
			newCache.GetArtifactClient(),
			newCache.GetProxyClient(),
		)

		return nil
	}
}
