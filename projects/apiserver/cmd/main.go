package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/setup"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/todo"
	"github.com/solo-io/solo-kit/test/config"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	port := flag.Int("p", 8082, "port to bind")
	dev := flag.Bool("dev", false, "use memory instead of connecting to real gloo storage")
	flag.Parse()
	glooOpts, err := config.DefaultKubernetesConstructOpts()
	if err != nil {
		return err
	}
	gatewayOpts, err := DefaultGatewayOpts()
	if err != nil {
		return err
	}
	sqoopOpts, err := DefaultSqoopOpts()
	if err != nil {
		return err
	}

	ctx := contextutils.WithLogger(context.Background(), "apiserver")

	contextutils.LoggerFrom(ctx).Infof("listening on :%v", *port)
	if err := setup.Setup(*port, *dev, glooOpts, gatewayOpts, sqoopOpts); err != nil {
		return err
	}
	return nil
}

func DefaultGatewayOpts() (gatewaysetup.Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return gatewaysetup.Opts{}, err
	}
	cache := kube.NewKubeCache()
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return gatewaysetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Gateways: &factory.KubeResourceClientFactory{
			Crd:         gatewayv1.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		VirtualServices: &factory.KubeResourceClientFactory{
			Crd:         gatewayv1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd:         gloov1.ProxyCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		WatchNamespaces: []string{"default", defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode: true,
	}, nil
}

func DefaultSqoopOpts() (sqoopsetup.Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return sqoopsetup.Opts{}, err
	}
	cache := kube.NewKubeCache()

	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return sqoopsetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Schemas: &factory.KubeResourceClientFactory{
			Crd:         sqoopv1.SchemaCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		ResolverMaps: &factory.KubeResourceClientFactory{
			Crd:         sqoopv1.ResolverMapCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd:         gloov1.ProxyCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		WatchNamespaces: []string{"default", defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode:     false,
		SidecarAddr: fmt.Sprintf("%v:%v", "127.0.0.1", TODO.SqoopSidecarBindPort),
	}, nil
}
