package file

import (
	"context"
	"fmt"
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
	"net"
)

// in the future, this may include other types of clients
// there should always be two, a file client for working with git
// and a client that matches the particular deployment
type DualClientSet struct {
	Kube ClientSet
	File ClientSet
}

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

func NewDualClient(deploymentType, fileClientRootDir string) (DualClientSet, error) {
	if deploymentType != "kube" {
		panic("we only support kubernetes clients at this time")
	}

	kubeSettingsClient, kubeGlooOpts, err := KubernetesConstructOpts()
	if err != nil {
		return DualClientSet{}, err
	}
	kubeGatewayOpts, err := KubeGatewayOpts()
	if err != nil {
		return DualClientSet{}, err
	}
	kubeSqoopOpts, err := KubeSqoopOpts()
	if err != nil {
		return DualClientSet{}, err
	}
	kubeClients, err := registerClients(kubeSettingsClient, kubeGlooOpts, kubeGatewayOpts, kubeSqoopOpts)
	if err != nil {
		return DualClientSet{}, err
	}

	fileSettingsClient, fileGlooOpts, err := FileConstructOpts(fileClientRootDir)
	if err != nil {
		return DualClientSet{}, err
	}
	fileGatewayOpts, err := FileGatewayOpts(fileClientRootDir)
	if err != nil {
		return DualClientSet{}, err
	}
	fileSqoopOpts, err := FileSqoopOpts(fileClientRootDir)
	if err != nil {
		return DualClientSet{}, err
	}
	fileClients, err := registerClients(fileSettingsClient, fileGlooOpts, fileGatewayOpts, fileSqoopOpts)
	if err != nil {
		return DualClientSet{}, err
	}
	return DualClientSet{
		Kube: kubeClients,
		File: fileClients,
	}, nil
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
	proxies, err := gloov1.NewProxyClient(glooOpts.Proxies)
	if err != nil {
		return ClientSet{}, err
	}
	if err := proxies.Register(); err != nil {
		return ClientSet{}, err
	}
	// ( we already created a settingsClient )
	if err := settings.Register(); err != nil {
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

func FileGatewayOpts(path string) (gatewaysetup.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return gatewaysetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Gateways: &factory.FileResourceClientFactory{
			// TODO(mitchdraft) make these constants
			RootDir: path + "/gateways",
		},
		VirtualServices: &factory.FileResourceClientFactory{
			RootDir: path + "/virtualservices",
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: path + "/proxies",
		},
		WatchNamespaces: []string{"default", defaults.GlooSystem},
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

func FileSqoopOpts(path string) (sqoopsetup.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return sqoopsetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Schemas: &factory.FileResourceClientFactory{
			RootDir: path + "/schemas",
		},
		ResolverMaps: &factory.FileResourceClientFactory{
			RootDir: path + "/resolvermaps",
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: path + "/proxies",
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

func FileConstructOpts(path string) (gloov1.SettingsClient, bootstrap.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gloo")

	// TODO(ilackarms): pass in settings configuration from an environment variable or CLI flag, rather than hard-coding to k8s
	settingsClient, err := gloov1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: path + "/settings",
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
			RootDir: path + "/upstreams",
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: path + "/proxies",
		},
		// TODO - make less sketchy (fileClients don't store secrets so we're using kube for that - find a good way to fill in the git/fileClient missing pieces)
		Secrets: &factory.KubeSecretClientFactory{
			Clientset: clientset,
		},
		// TODO - add to fileClient?
		Artifacts: &factory.KubeConfigMapClientFactory{
			Clientset: clientset,
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
		// KubeClient: clientset,
		DevMode: true,
	}, nil
}
