package setup

import (
	"context"
	"fmt"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/todo"
)

func Main(settingsDir string) error {
	settingsClient, err := gloov1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: settingsDir,
	})
	if err != nil {
		return err
	}
	cache := gloov1.NewSetupEmitter(settingsClient)
	ctx := contextutils.WithLogger(context.Background(), "sqoop")
	eventLoop := gloov1.NewSetupEventLoop(cache, NewSetupSyncer())
	errs, err := eventLoop.Run([]string{"settings"}, clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Second,
	})
	if err != nil {
		return err
	}
	for err := range errs {
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}
	return nil
}

type Opts struct {
	WriteNamespace  string
	WatchNamespaces []string
	Schemas         factory.ResourceClientFactory
	ResolverMaps    factory.ResourceClientFactory
	Proxies         factory.ResourceClientFactory
	WatchOpts       clients.WatchOpts
	DevMode         bool
	SidecarAddr     string
}

func DefaultKubernetesConstructOpts() (Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return Opts{}, err
	}
	// clientset, err := kubernetes.NewForConfig(cfg)
	// if err != nil {
	// 	return Opts{}, err
	// }
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return Opts{
		WriteNamespace: defaults.GlooSystem,
		Schemas: &factory.KubeResourceClientFactory{
			Crd: v1.SchemaCrd,
			Cfg: cfg,
		},
		ResolverMaps: &factory.KubeResourceClientFactory{
			Crd: v1.ResolverMapCrd,
			Cfg: cfg,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd: gloov1.ProxyCrd,
			Cfg: cfg,
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
