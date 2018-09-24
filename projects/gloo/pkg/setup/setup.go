package setup

import (
	"context"
	"net"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

// TODO(ilackarms): remove this or move it to a test package, only use settings watch for prodution gloo
func writeSettings(settingsDir string) error {
	cli, err := v1.NewSettingsClient(&factory.FileResourceClientFactory{RootDir: settingsDir})
	if err != nil {
		return err
	}
	settings := &v1.Settings{
		ConfigSource: &v1.Settings_KubernetesConfigSource{
			KubernetesConfigSource: &v1.Settings_KubernetesCrds{},
		},
		ArtifactSource: &v1.Settings_KubernetesArtifactSource{
			KubernetesArtifactSource: &v1.Settings_KubernetesConfigmaps{},
		},
		SecretSource: &v1.Settings_KubernetesSecretSource{
			KubernetesSecretSource: &v1.Settings_KubernetesSecrets{},
		},
		BindAddr:    "0.0.0.0:9977",
		RefreshRate: types.DurationProto(time.Minute),
		DevMode:     true,
		Metadata: core.Metadata{
			Namespace: "settings",
			Name:      "gloo",
		},
	}
	cli.Delete(settings.Metadata.Namespace, settings.Metadata.Name, clients.DeleteOpts{})
	_, err = cli.Write(settings, clients.WriteOpts{})
	return err
}

func Main(settingsDir string) error {
	if err := writeSettings(settingsDir); err != nil && !errors.IsExist(err) {
		return err
	}
	settingsClient, err := v1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: settingsDir,
	})
	if err != nil {
		return err
	}
	cache := v1.NewSetupEmitter(settingsClient)
	ctx := contextutils.WithLogger(context.Background(), "gloo")
	eventLoop := v1.NewSetupEventLoop(cache, syncer.NewSetupSyncer())
	errs, err := eventLoop.Run([]string{"settings"}, clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Minute,
	})
	if err != nil {
		return err
	}
	for err := range errs {
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}
	return nil
}

func DefaultKubernetesConstructOpts() (bootstrap.Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return bootstrap.Opts{}, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return bootstrap.Opts{}, err
	}
	ctx := contextutils.WithLogger(context.Background(), "gloo")
	logger := contextutils.LoggerFrom(ctx)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop()),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				logger.Infof("gRPC call: %v", info.FullMethod)
				return handler(srv, ss)
			},
		)),
	)
	return bootstrap.Opts{
		WriteNamespace: defaults.GlooSystem,
		Upstreams: &factory.KubeResourceClientFactory{
			Crd: v1.UpstreamCrd,
			Cfg: cfg,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd: v1.ProxyCrd,
			Cfg: cfg,
		},
		Secrets: &factory.KubeSecretClientFactory{
			Clientset: clientset,
		},
		Artifacts: &factory.KubeConfigMapClientFactory{
			Clientset: clientset,
		},
		Namespacer: static.NewNamespacer([]string{"default", defaults.GlooSystem}),
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		BindAddr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8080,
		},
		GrpcServer: grpcServer,
		KubeClient: clientset,
		DevMode:    true,
	}, nil
}
