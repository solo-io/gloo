package client

import (
	"context"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"
	"k8s.io/client-go/rest"
)

//go:generate mockgen -destination mocks/mock_client_updater.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client Updater

// Run a settings watch loop in a new goroutine, rebuilding all the k8s resource clients when the settings are updated
type Updater interface {
	StartWatch(ctx context.Context)
}

const clientUpdaterLogger = "clientUpdater"

type updater struct {
	clientCache  ClientCache
	token        setup.Token
	cfg          *rest.Config
	podNamespace string
}

var _ Updater = &updater{}

func NewClientUpdater(clientCache ClientCache, cfg *rest.Config, token setup.Token, podNamespace string) Updater {
	return &updater{
		clientCache:  clientCache,
		cfg:          cfg,
		token:        token,
		podNamespace: podNamespace,
	}
}

func (u *updater) StartWatch(ctx context.Context) {
	go func() {
		// setuputils.Main will block until the watch loop ends
		err := setuputils.Main(setuputils.SetupOpts{
			LoggerName:  clientUpdaterLogger,
			Version:     version.Version,
			SetupFunc:   u.buildReceiverFunc(),
			ExitOnError: false,
		})

		if err != nil {
			contextutils.LoggerFrom(ctx).Warnf("Settings watch loop has died, attempting to restart")

			// let this goroutine die and start another
			go u.StartWatch(ctx)
		}
	}()
}

func (u *updater) buildReceiverFunc() setuputils.SetupFunc {
	return func(ctx context.Context, kubeCache kube.SharedCache, inMemoryCache memory.InMemoryResourceCache, settings *gloov1.Settings) error {

		// build a new cache from scratch, and update the existing one from the new one's state
		newCache, err := NewClientCache(ctx, settings, u.cfg, u.token, u.podNamespace)
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
			newCache.GetUpstreamGroupClient(),
			newCache.GetRouteTableClient(),
		)

		return nil
	}
}
