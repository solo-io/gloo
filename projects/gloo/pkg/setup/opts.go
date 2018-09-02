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
	"github.com/solo-io/solo-kit/pkg/namespacing"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

type Opts struct {
	writeNamespace string
	upstreams      factory.ResourceClientFactoryOpts
	proxies        factory.ResourceClientFactoryOpts
	secrets        factory.ResourceClientFactoryOpts
	artifacts      factory.ResourceClientFactoryOpts
	namespacer     namespacing.Namespacer
	bindAddr       net.Addr
	grpcServer     *grpc.Server
	watchOpts      clients.WatchOpts
}

//  ilackarms: We can just put any hacky stuff we need here

func DefaultKubernetesConstructOpts() (Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return Opts{}, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return Opts{}, err
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
	return Opts{
		writeNamespace: "gloo-system",
		upstreams: &factory.KubeResourceClientOpts{
			Crd: v1.UpstreamCrd,
			Cfg: cfg,
		},
		proxies: &factory.KubeResourceClientOpts{
			Crd: v1.ProxyCrd,
			Cfg: cfg,
		},
		secrets: &factory.KubeSecretClientOpts{
			Clientset: clientset,
		},
		artifacts: &factory.KubeConfigMapClientOpts{
			Clientset: clientset,
		},
		namespacer: static.NewNamespacer([]string{"default", "gloo-system"}),
		watchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		},
		bindAddr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8080,
		},
		grpcServer: grpcServer,
	}, nil
}
