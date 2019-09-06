package client

import (
	"context"
	"sync"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"

	consulapi "github.com/hashicorp/consul/api"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//go:generate mockgen -destination mocks/mock_ks8_client_cache.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client ClientCache

// Factory for producing k8s resource clients
// References returned by these methods may change from one call to the next, so they should not be kept around long-term
// All methods will attempt to acquire a mutex, so they may block the current goroutine
type ClientCache interface {
	GetVirtualServiceClient() gatewayv1.VirtualServiceClient
	GetGatewayClient() gatewayv2.GatewayClient
	GetUpstreamClient() gloov1.UpstreamClient
	GetSettingsClient() gloov1.SettingsClient
	GetSecretClient() gloov1.SecretClient
	GetArtifactClient() gloov1.ArtifactClient
	GetProxyClient() gloov1.ProxyClient

	SetCacheState(virtualServiceClient gatewayv1.VirtualServiceClient,
		gatewayClient gatewayv2.GatewayClient,
		upstreamClient gloov1.UpstreamClient,
		settingsClient gloov1.SettingsClient,
		secretClient gloov1.SecretClient,
		artifactClient gloov1.ArtifactClient,
		proxyClient gloov1.ProxyClient,
	)
}

type clientCache struct {
	virtualServiceClient gatewayv1.VirtualServiceClient
	gatewayClient        gatewayv2.GatewayClient
	upstreamClient       gloov1.UpstreamClient
	settingsClient       gloov1.SettingsClient
	secretClient         gloov1.SecretClient
	artifactClient       gloov1.ArtifactClient
	proxyClient          gloov1.ProxyClient

	// pointer updates in go are not guaranteed to be atomic, so be sure to lock and unlock when updating or getting a reference to a client
	cacheUpdateMutex *sync.Mutex
}

var _ ClientCache = &clientCache{}

func (c *clientCache) GetVirtualServiceClient() gatewayv1.VirtualServiceClient {
	c.cacheUpdateMutex.Lock()
	vsClient := c.virtualServiceClient
	c.cacheUpdateMutex.Unlock()

	return vsClient
}

func (c *clientCache) GetGatewayClient() gatewayv2.GatewayClient {
	c.cacheUpdateMutex.Lock()
	ret := c.gatewayClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) GetUpstreamClient() gloov1.UpstreamClient {
	c.cacheUpdateMutex.Lock()
	ret := c.upstreamClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) GetSettingsClient() gloov1.SettingsClient {
	c.cacheUpdateMutex.Lock()
	ret := c.settingsClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) GetSecretClient() gloov1.SecretClient {
	c.cacheUpdateMutex.Lock()
	ret := c.secretClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) GetArtifactClient() gloov1.ArtifactClient {
	c.cacheUpdateMutex.Lock()
	ret := c.artifactClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) GetProxyClient() gloov1.ProxyClient {
	c.cacheUpdateMutex.Lock()
	ret := c.proxyClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) SetCacheState(virtualServiceClient gatewayv1.VirtualServiceClient,
	gatewayClient gatewayv2.GatewayClient,
	upstreamClient gloov1.UpstreamClient,
	settingsClient gloov1.SettingsClient,
	secretClient gloov1.SecretClient,
	artifactClient gloov1.ArtifactClient,
	proxyClient gloov1.ProxyClient,
) {
	c.cacheUpdateMutex.Lock()

	c.virtualServiceClient = virtualServiceClient
	c.gatewayClient = gatewayClient
	c.upstreamClient = upstreamClient
	c.settingsClient = settingsClient
	c.secretClient = secretClient
	c.artifactClient = artifactClient
	c.proxyClient = proxyClient

	c.cacheUpdateMutex.Unlock()
}

// Returns a set of clients that use Kubernetes as storage
func NewClientCache(ctx context.Context, settings *gloov1.Settings, cfg *rest.Config, token setup.Token) (ClientCache, error) {

	// New shared cache
	k8sCache := kube.NewKubeCache(context.TODO())

	var clientset kubernetes.Interface
	memCache := memory.NewInMemoryResourceCache()
	opts, err := constructOpts(ctx, &clientset, k8sCache, nil, nil, memCache, settings)
	if err != nil {
		return nil, err
	}

	upstreamClient, err := gloov1.NewUpstreamClientWithToken(opts.Upstreams, *token)
	if err != nil {
		return nil, err
	}

	vsClient, err := gatewayv1.NewVirtualServiceClientWithToken(factoryFor(gatewayv1.VirtualServiceCrd, *cfg, k8sCache), *token)
	if err != nil {
		return nil, err
	}

	gatewayClient, err := gatewayv2.NewGatewayClientWithToken(factoryFor(gatewayv2.GatewayCrd, *cfg, k8sCache), *token)
	if err != nil {
		return nil, err
	}

	proxyClient, err := gloov1.NewProxyClientWithToken(factoryFor(gloov1.ProxyCrd, *cfg, k8sCache), *token)
	if err != nil {
		return nil, err
	}

	settingsClient, err := gloov1.NewSettingsClientWithToken(factoryFor(gloov1.SettingsCrd, *cfg, k8sCache), *token)
	if err != nil {
		return nil, err
	}

	// Needed only for the clients backed by the KubeResourceClientFactory
	// so that they register with the cache they share
	if err = registerAll(upstreamClient, vsClient, settingsClient, gatewayClient, proxyClient); err != nil {
		return nil, err
	}

	// replace this with the gloo factory
	secretClient, err := gloov1.NewSecretClientWithToken(opts.Secrets, *token)
	if err != nil {
		return nil, err
	}

	artifactClient, err := gloov1.NewArtifactClientWithToken(opts.Artifacts, *token)
	if err != nil {
		return nil, err
	}

	cache := &clientCache{
		cacheUpdateMutex: &sync.Mutex{},
	}

	cache.SetCacheState(
		vsClient,
		gatewayClient,
		upstreamClient,
		settingsClient,
		secretClient,
		artifactClient,
		proxyClient,
	)

	return cache, nil
}

type registrant interface {
	Register() error
}

func registerAll(clients ...registrant) error {
	for _, client := range clients {
		if err := client.Register(); err != nil {
			return err
		}
	}
	return nil
}

func factoryFor(crd crd.Crd, cfg rest.Config, cache kube.SharedCache) factory.ResourceClientFactory {
	return &factory.KubeResourceClientFactory{
		Crd:                crd,
		Cfg:                &cfg,
		SharedCache:        cache,
		NamespaceWhitelist: []string{v1.NamespaceAll},
		SkipCrdCreation:    settingsutil.GetSkipCrdCreation(),
	}
}

// TODO make this public in Gloo...
func constructOpts(ctx context.Context, clientset *kubernetes.Interface, kubeCache kube.SharedCache, consulClient *consulapi.Client, vaultClient *vaultapi.Client, memCache memory.InMemoryResourceCache, settings *gloov1.Settings) (bootstrap.Opts, error) {

	var (
		cfg           *rest.Config
		kubeCoreCache corecache.KubeCoreCache
	)

	params := bootstrap.NewConfigFactoryParams(
		settings,
		memCache,
		kubeCache,
		&cfg,
		consulClient,
	)

	upstreamFactory, err := bootstrap.ConfigFactoryForSettings(params, gloov1.UpstreamCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	kubeServiceClient, err := bootstrap.KubeServiceClientForSettings(
		ctx,
		settings,
		memCache,
		&cfg,
		clientset,
		&kubeCoreCache,
	)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	proxyFactory, err := bootstrap.ConfigFactoryForSettings(params, gloov1.ProxyCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	secretFactory, err := bootstrap.SecretFactoryForSettings(
		ctx,
		settings,
		memCache,
		&cfg,
		clientset,
		&kubeCoreCache,
		vaultClient,
		gloov1.SecretCrd.Plural,
	)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	upstreamGroupFactory, err := bootstrap.ConfigFactoryForSettings(params, gloov1.UpstreamGroupCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	artifactFactory, err := bootstrap.ArtifactFactoryForSettings(
		ctx,
		settings,
		memCache,
		&cfg,
		clientset,
		&kubeCoreCache,
		consulClient,
		gloov1.ArtifactCrd.Plural,
	)
	if err != nil {
		return bootstrap.Opts{}, err
	}
	return bootstrap.Opts{
		Upstreams:         upstreamFactory,
		KubeServiceClient: kubeServiceClient,
		Proxies:           proxyFactory,
		UpstreamGroups:    upstreamGroupFactory,
		Secrets:           secretFactory,
		Artifacts:         artifactFactory,
	}, nil
}
