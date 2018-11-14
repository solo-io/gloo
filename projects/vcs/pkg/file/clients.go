package file

import (
	"context"
	"fmt"
	"net"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	gatewayv1 "github.com/solo-io/solo-projects/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-projects/projects/gateway/pkg/syncer"
	gloov1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
	sqoopv1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-projects/projects/sqoop/pkg/syncer"
	"github.com/solo-io/solo-projects/projects/sqoop/pkg/todo"
	"k8s.io/client-go/kubernetes"
)

const (
	gatewayRootDir        = "/gateways"
	virtualServiceRootDir = "/virtual-services"
	proxyRootDir          = "/proxies"
	schemaRootDir         = "/schemas"
	resolverMapRootDir    = "/resolver-maps"
	upstreamRootDir       = "/upstreams"
	settingsRootDir       = "/settings"
)

type ClientSet struct {
	gloov1.UpstreamClient
	gatewayv1.VirtualServiceClient
	gatewayv1.GatewayClient
	gloov1.ProxyClient
	gloov1.SettingsClient
	gloov1.SecretClient
	gloov1.ArtifactClient
	sqoopv1.ResolverMapClient
	sqoopv1.SchemaClient
}

func NewFileClient(rootDir string) (ClientSet, error) {

	fileSettingsClient, fileGlooOpts, err := FileConstructOpts(rootDir)
	if err != nil {
		return ClientSet{}, err
	}
	fileGatewayOpts, err := FileGatewayOpts(rootDir)
	if err != nil {
		return ClientSet{}, err
	}
	fileSqoopOpts, err := FileSqoopOpts(rootDir)
	if err != nil {
		return ClientSet{}, err
	}
	return registerClients(fileSettingsClient, fileGlooOpts, fileGatewayOpts, fileSqoopOpts)
}

func NewKubeClient() (ClientSet, error) {
	kubeSettingsClient, kubeGlooOpts, err := KubernetesConstructOpts()
	if err != nil {
		return ClientSet{}, err
	}
	kubeGatewayOpts, err := KubeGatewayOpts()
	if err != nil {
		return ClientSet{}, err
	}
	kubeSqoopOpts, err := KubeSqoopOpts()
	if err != nil {
		return ClientSet{}, err
	}
	return registerClients(kubeSettingsClient, kubeGlooOpts, kubeGatewayOpts, kubeSqoopOpts)
}

func registerClients(settings gloov1.SettingsClient, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) (ClientSet, error) {
	// initial resource registration
	upstreams, err := gloov1.NewUpstreamClient(glooOpts.Upstreams)
	if err != nil {
		return ClientSet{}, err
	}
	if err := upstreams.Register(); err != nil {
		return ClientSet{}, err
	}
	secrets, err := gloov1.NewSecretClient(glooOpts.Secrets)
	if err != nil {
		return ClientSet{}, err
	}
	if err := secrets.Register(); err != nil {
		return ClientSet{}, err
	}
	artifacts, err := gloov1.NewArtifactClient(glooOpts.Artifacts)
	if err != nil {
		return ClientSet{}, err
	}
	if err := artifacts.Register(); err != nil {
		return ClientSet{}, err
	}
	proxies, err := gloov1.NewProxyClient(glooOpts.Proxies)
	if err != nil {
		return ClientSet{}, err
	}
	if err := proxies.Register(); err != nil {
		return ClientSet{}, err
	}
	virtualServices, err := gatewayv1.NewVirtualServiceClient(gatewayOpts.VirtualServices)
	if err != nil {
		return ClientSet{}, err
	}
	if err := virtualServices.Register(); err != nil {
		return ClientSet{}, err
	}
	gateways, err := gatewayv1.NewGatewayClient(gatewayOpts.Gateways)
	if err != nil {
		return ClientSet{}, err
	}
	if err := gateways.Register(); err != nil {
		return ClientSet{}, err
	}
	resolverMaps, err := sqoopv1.NewResolverMapClient(sqoopOpts.ResolverMaps)
	if err != nil {
		return ClientSet{}, err
	}
	if err := resolverMaps.Register(); err != nil {
		return ClientSet{}, err
	}
	schemas, err := sqoopv1.NewSchemaClient(sqoopOpts.Schemas)
	if err != nil {
		return ClientSet{}, err
	}
	if err := schemas.Register(); err != nil {
		return ClientSet{}, err
	}

	// ( we already created a settingsClient )
	if err := settings.Register(); err != nil {
		return ClientSet{}, err
	}

	return ClientSet{
		UpstreamClient:       upstreams,
		VirtualServiceClient: virtualServices,
		GatewayClient:        gateways,
		ProxyClient:          proxies,
		SettingsClient:       settings,
		SecretClient:         secrets,
		ArtifactClient:       artifacts,
		ResolverMapClient:    resolverMaps,
		SchemaClient:         schemas,
	}, nil
}

func KubeGatewayOpts() (gatewaysetup.Opts, error) {
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
		WatchNamespaces: []string{clients.DefaultNamespace, defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode: true,
	}, nil
}

func FileGatewayOpts(path string) (gatewaysetup.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return gatewaysetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Gateways: &factory.FileResourceClientFactory{
			RootDir: path + gatewayRootDir,
		},
		VirtualServices: &factory.FileResourceClientFactory{
			RootDir: path + virtualServiceRootDir,
		},
		WatchNamespaces: []string{defaults.GlooSystem, defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode: true,
	}, nil
}

func KubeSqoopOpts() (sqoopsetup.Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return sqoopsetup.Opts{}, err
	}
	cache := kube.NewKubeCache()

	ctx := contextutils.WithLogger(context.Background(), "sqoop")
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
		WatchNamespaces: []string{clients.DefaultNamespace, defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode:     false,
		SidecarAddr: fmt.Sprintf("%v:%v", "127.0.0.1", TODO.SqoopSidecarBindPort),
	}, nil
}

func FileSqoopOpts(path string) (sqoopsetup.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "sqoop")
	return sqoopsetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Schemas: &factory.FileResourceClientFactory{
			RootDir: path + schemaRootDir,
		},
		ResolverMaps: &factory.FileResourceClientFactory{
			RootDir: path + resolverMapRootDir,
		},
		WatchNamespaces: []string{clients.DefaultNamespace, defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode:     false,
		SidecarAddr: fmt.Sprintf("%v:%v", "127.0.0.1", TODO.SqoopSidecarBindPort),
	}, nil
}

func KubernetesConstructOpts() (gloov1.SettingsClient, bootstrap.Opts, error) {
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
	settingsClient, err := gloov1.NewSettingsClient(&factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: cache,
	})
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}

	return settingsClient, bootstrap.Opts{
		WriteNamespace: defaults.GlooSystem,
		Upstreams: &factory.KubeResourceClientFactory{
			Crd:         gloov1.UpstreamCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd:         gloov1.ProxyCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		Secrets: &factory.KubeSecretClientFactory{
			Clientset: clientset,
		},
		Artifacts: &factory.KubeConfigMapClientFactory{
			Clientset: clientset,
		},
		WatchNamespaces: []string{clients.DefaultNamespace, defaults.GlooSystem},
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

func FileConstructOpts(path string) (gloov1.SettingsClient, bootstrap.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gloo")

	// TODO(ilackarms): pass in settings configuration from an environment variable or CLI flag, rather than hard-coding to k8s
	settingsClient, err := gloov1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: path + settingsRootDir,
	})
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}

	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, bootstrap.Opts{}, err
	}

	return settingsClient, bootstrap.Opts{
		WriteNamespace: defaults.GlooSystem,
		Upstreams: &factory.FileResourceClientFactory{
			RootDir: path + upstreamRootDir,
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: path + proxyRootDir,
		},
		// TODO - make less sketchy (fileClients don't store secrets so we're using kube for that - find a good way to fill in the git/fileClient missing pieces)
		Secrets: &factory.KubeSecretClientFactory{
			Clientset: clientset,
		},
		// TODO - add to fileClient?
		Artifacts: &factory.KubeConfigMapClientFactory{
			Clientset: clientset,
		},
		WatchNamespaces: []string{clients.DefaultNamespace, defaults.GlooSystem},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		BindAddr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8080,
		},
		// KubeClient: clientset,
		DevMode: true,
	}, nil
}
