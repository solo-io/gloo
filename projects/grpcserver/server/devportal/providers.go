package devportal

import (
	"time"

	"github.com/google/wire"
	devportalgrpc "github.com/solo-io/dev-portal/pkg/services/grpc/portal"
	"go.uber.org/zap"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"context"

	devportalv1 "github.com/solo-io/dev-portal/pkg/api/devportal.solo.io/v1"
)

var ProviderSet = wire.NewSet(
	NewManager,
	NewDynamicClient,
	NewPortalClient,
	devportalgrpc.NewPortalGrpcService,
	NewRegistrar,
)

// Starts a manager and binds it to the root context
// TODO(marco): consider abstracting the manager to provide more control over start/stop/error detection,
// 	taking this as an example: https://github.com/solo-io/mesh-projects/blob/0fd27533a40c8959269e7bd8beab6bea13182fc0/services/common/multicluster/manager/interfaces.go#L41
func NewManager(ctx context.Context, cfg *rest.Config, podNamespace string) (controllerruntime.Manager, error) {
	mgrCtx, cancel := context.WithCancel(ctx)

	mgr, err := manager.New(cfg, manager.Options{
		// Dev portal resources always reside in the same namespace as Gloo
		Namespace: podNamespace,
		// TODO(marco): should conditionally enable
		MetricsBindAddress: "0",
	})
	if err != nil {
		return nil, err
	}

	// Add dev portal types to manager scheme
	if err := devportalv1.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, err
	}

	go func() {
		// Blocks until ctx is done
		if err := mgr.Start(mgrCtx.Done()); err != nil {
			contextutils.LoggerFrom(mgrCtx).Errorw("error while starting manager", zap.Error(err))
		}
		contextutils.LoggerFrom(mgrCtx).Info("manager has been shut down")
	}()

	withDeadline, _ := context.WithTimeout(mgrCtx, 5*time.Second)
	if synced := mgr.GetCache().WaitForCacheSync(withDeadline.Done()); !synced {
		cancel()
		return nil, errors.New("timed out waiting for manager caches to sync")
	}
	return mgr, nil
}

func NewDynamicClient(manager controllerruntime.Manager) client.Client {
	return manager.GetClient()
}

// We need this additional wrapper because wire expects providers
// to return the exact type that is needed for injection.
func NewPortalClient(client client.Client) devportalv1.PortalClient {
	return devportalv1.NewPortalClient(client)
}
