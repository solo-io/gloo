package devportal

import (
	"os"
	"time"

	"github.com/solo-io/dev-portal/pkg/admin"
	"github.com/solo-io/dev-portal/pkg/admin/grpc/apidoc"
	"github.com/solo-io/dev-portal/pkg/admin/grpc/group"
	"github.com/solo-io/dev-portal/pkg/admin/grpc/user"

	adminapi "github.com/solo-io/dev-portal/pkg/api/grpc/admin"
	devportalkubev1 "github.com/solo-io/dev-portal/pkg/api/kube/core/v1"
	"github.com/solo-io/dev-portal/pkg/assets"

	"github.com/google/wire"
	portalgrpc "github.com/solo-io/dev-portal/pkg/admin/grpc/portal"
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
	NewApiDocClient,
	NewUserClient,
	NewGroupClient,
	NewConfigMapClient,
	NewSecretClient,
	assets.NewConfigMapStorage,
	wire.Bind(new(assets.AssetStorage), new(*assets.ConfigMapStorage)),
	admin.NewResourceLabeler,
	wire.Bind(new(admin.ResourceLinker), new(*admin.ResourceLabeler)),
	portalgrpc.NewPortalGrpcService,
	wire.Bind(new(adminapi.PortalApiServer), new(*portalgrpc.GrpcService)),
	apidoc.NewApiDocGrpcService,
	wire.Bind(new(adminapi.ApiDocApiServer), new(*apidoc.GrpcService)),
	user.NewUserGrpcService,
	wire.Bind(new(adminapi.UserApiServer), new(*user.GrpcService)),
	group.NewGroupGrpcService,
	wire.Bind(new(adminapi.GroupApiServer), new(*group.GrpcService)),
	NewRegistrar,
)

// Starts a manager and binds it to the root context
// TODO(marco): consider abstracting the manager to provide more control over start/stop/error detection,
// 	taking this as an example: https://github.com/solo-io/mesh-projects/blob/0fd27533a40c8959269e7bd8beab6bea13182fc0/services/common/multicluster/manager/interfaces.go#L41
func NewManager(ctx context.Context, cfg *rest.Config, podNamespace string) (controllerruntime.Manager, error) {
	mgrCtx, cancel := context.WithCancel(ctx)

	watchNamespace := podNamespace
	if os.Getenv("RBAC_NAMESPACED") == "false" {
		// If we are not running in namespaced mode, watches all namespaces
		watchNamespace = ""
	}

	mgr, err := manager.New(cfg, manager.Options{
		Namespace: watchNamespace,
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

// We need these additional wrappers because wire expects providers to return the exact type that is needed
// for injection and we cannot use `Bind` as the structs returned by the New*Client functions are not exported.
func NewPortalClient(client client.Client) devportalv1.PortalClient {
	return devportalv1.NewPortalClient(client)
}
func NewApiDocClient(client client.Client) devportalv1.ApiDocClient {
	return devportalv1.NewApiDocClient(client)
}
func NewUserClient(client client.Client) devportalv1.UserClient {
	return devportalv1.NewUserClient(client)
}
func NewGroupClient(client client.Client) devportalv1.GroupClient {
	return devportalv1.NewGroupClient(client)
}

func NewConfigMapClient(client client.Client) devportalkubev1.ConfigMapClient {
	return devportalkubev1.NewConfigMapClient(client)
}
func NewSecretClient(client client.Client) devportalkubev1.SecretClient {
	return devportalkubev1.NewSecretClient(client)
}
