package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/setup"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/todo"
	"k8s.io/client-go/kubernetes"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	stats.StartStatsServer()
	port := flag.Int("p", 8082, "port to bind")
	dev := flag.Bool("dev", false, "use memory instead of connecting to real gloo storage")
	flag.Parse()
	debugMode := os.Getenv("DEBUG") == "1"
	settingsClient, glooOpts, err := DefaultKubernetesConstructOpts()
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
	if err := setup.Setup(ctx, *port, *dev, debugMode, settingsClient, glooOpts, gatewayOpts, sqoopOpts); err != nil {
		return err
	}
	return nil
}

// DEPRECATED: TODO(ilackarms): remove without breaking things, move to a test package
func DefaultKubernetesConstructOpts() (v1.SettingsClient, bootstrap.Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}
	ctx := contextutils.WithLogger(context.Background(), "gloo")
	cache := kube.NewKubeCache()

	// TODO(ilackarms): pass in settings configuration from an environment variable or CLI flag, rather than hard-coding to k8s
	settingsClient, err := v1.NewSettingsClient(&factory.KubeResourceClientFactory{
		Crd:         v1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: cache,
	})
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}

	return settingsClient, bootstrap.Opts{
		WriteNamespace: defaults.GlooSystem,
		Upstreams: &factory.KubeResourceClientFactory{
			Crd:         v1.UpstreamCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd:         v1.ProxyCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		Secrets: &factory.KubeSecretClientFactory{
			Clientset: clientset,
		},
		Artifacts: &factory.KubeConfigMapClientFactory{
			Clientset: clientset,
		},
		Schemas: &factory.KubeResourceClientFactory{
			Crd:         sqoopv1.SchemaCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		WatchNamespaces: []string{"default", defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		BindAddr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8080,
		},
		KubeClient: clientset,
		DevMode:    true,
	}, nil
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
