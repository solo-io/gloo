package file

import (
	"context"
	"fmt"
	"net"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gatewaysetup "github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	sqoopsetup "github.com/solo-io/solo-kit/projects/sqoop/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/todo"
)

// in the future, this may include other types of clients
// there should always be two, a file client for working with git
// and a client that matches the particular deployment
type DualClientSet struct {
	Kube ClientSet
	File ClientSet
}

type ClientSet struct {
	v1.UpstreamClient
	gatewayv1.VirtualServiceClient
	v1.SettingsClient
	v1.SecretClient
	v1.ArtifactClient
	sqoopv1.ResolverMapClient
	sqoopv1.SchemaClient
}

func NewDualClient(deploymentType string) (DualClientSet, error) {
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

	fileSettingsClient, fileGlooOpts, err := FileConstructOpts()
	if err != nil {
		return DualClientSet{}, err
	}
	fileGatewayOpts, err := FileGatewayOpts()
	if err != nil {
		return DualClientSet{}, err
	}
	fileSqoopOpts, err := FileSqoopOpts()
	if err != nil {
		return DualClientSet{}, err
	}
	fileClients, err := registerClients(fileSettingsClient, fileGlooOpts, fileGatewayOpts, fileSqoopOpts)
	if err != nil {
		return DualClientSet{}, err
	}
	// TODO -- placeholder
	// kubeClients, err := registerClients(kubeSettingsClient, kubeGlooOpts, kubeGatewayOpts, kubeSqoopOpts)
	// fileClients, err := registerClients(fileSettingsClient, fileGlooOpts, fileGatewayOpts, fileSqoopOpts)
	return DualClientSet{
		Kube: kubeClients,
		File: fileClients,
	}, nil
}

func registerClients(settings v1.SettingsClient, glooOpts bootstrap.Opts, gatewayOpts gatewaysetup.Opts, sqoopOpts sqoopsetup.Opts) (ClientSet, error) {
	// initial resource registration
	upstreams, err := v1.NewUpstreamClient(glooOpts.Upstreams)
	if err != nil {
		return ClientSet{}, err
	}
	if err := upstreams.Register(); err != nil {
		return ClientSet{}, err
	}
	secrets, err := v1.NewSecretClient(glooOpts.Secrets)
	if err != nil {
		return ClientSet{}, err
	}
	if err := secrets.Register(); err != nil {
		return ClientSet{}, err
	}
	artifacts, err := v1.NewArtifactClient(glooOpts.Artifacts)
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

func FileGatewayOpts() (gatewaysetup.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return gatewaysetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Gateways: &factory.FileResourceClientFactory{
			// TODO(mitchdraft) make these constants
			RootDir: "gloo/gateways",
		},
		VirtualServices: &factory.FileResourceClientFactory{
			RootDir: "gloo/virtualservices",
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: "gloo/proxies",
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

func FileSqoopOpts() (sqoopsetup.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return sqoopsetup.Opts{
		WriteNamespace: defaults.GlooSystem,
		Schemas: &factory.FileResourceClientFactory{
			RootDir: "gloo/schemas",
		},
		ResolverMaps: &factory.FileResourceClientFactory{
			RootDir: "gloo/resolvermaps",
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: "gloo/proxies",
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

func KubernetesConstructOpts() (v1.SettingsClient, bootstrap.Opts, error) {
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

func FileConstructOpts() (v1.SettingsClient, bootstrap.Opts, error) {
	ctx := contextutils.WithLogger(context.Background(), "gloo")

	// TODO(ilackarms): pass in settings configuration from an environment variable or CLI flag, rather than hard-coding to k8s
	settingsClient, err := v1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: "gloo/settings",
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
			RootDir: "gloo/upstreams",
		},
		Proxies: &factory.FileResourceClientFactory{
			RootDir: "gloo/proxies",
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
