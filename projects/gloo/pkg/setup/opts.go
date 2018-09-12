package setup

import (
	"context"
	"net"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
)

//  ilackarms: We can just put any hacky stuff we need here

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
				logger.Debugf("gRPC call: %v", info.FullMethod)
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
			RefreshRate: time.Minute,
		},
		BindAddr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8080,
		},
		GrpcServer: grpcServer,
		KubeClient: clientset,
		DevMode:    false,
	}, nil
}
