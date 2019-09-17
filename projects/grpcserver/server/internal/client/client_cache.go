package client

import (
	"context"
	"sync"

	"github.com/solo-io/gloo/pkg/utils"
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
	GetUpstreamGroupClient() gloov1.UpstreamGroupClient
	GetRouteTableClient() gatewayv1.RouteTableClient

	SetCacheState(virtualServiceClient gatewayv1.VirtualServiceClient,
		gatewayClient gatewayv2.GatewayClient,
		upstreamClient gloov1.UpstreamClient,
		settingsClient gloov1.SettingsClient,
		secretClient gloov1.SecretClient,
		artifactClient gloov1.ArtifactClient,
		proxyClient gloov1.ProxyClient,
		upstreamGroupClient gloov1.UpstreamGroupClient,
		routeTableClient gatewayv1.RouteTableClient,
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
	upstreamGroupClient  gloov1.UpstreamGroupClient
	routeTableClient     gatewayv1.RouteTableClient

	// pointer updates in go are not guaranteed to be atomic, so be sure to lock and unlock when updating or getting a reference to a client
	cacheUpdateMutex sync.RWMutex
}

var _ ClientCache = &clientCache{}

func (c *clientCache) GetVirtualServiceClient() gatewayv1.VirtualServiceClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	vsClient := c.virtualServiceClient

	return vsClient
}

func (c *clientCache) GetGatewayClient() gatewayv2.GatewayClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	ret := c.gatewayClient

	return ret
}

func (c *clientCache) GetUpstreamClient() gloov1.UpstreamClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	ret := c.upstreamClient

	return ret
}

func (c *clientCache) GetSettingsClient() gloov1.SettingsClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	ret := c.settingsClient

	return ret
}

func (c *clientCache) GetSecretClient() gloov1.SecretClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	ret := c.secretClient

	return ret
}

func (c *clientCache) GetArtifactClient() gloov1.ArtifactClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	ret := c.artifactClient

	return ret
}

func (c *clientCache) GetProxyClient() gloov1.ProxyClient {
	c.cacheUpdateMutex.RLock()
	defer c.cacheUpdateMutex.RUnlock()
	ret := c.proxyClient

	return ret
}

func (c *clientCache) GetRouteTableClient() gatewayv1.RouteTableClient {
	c.cacheUpdateMutex.Lock()
	ret := c.routeTableClient
	c.cacheUpdateMutex.Unlock()

	return ret
}

func (c *clientCache) GetUpstreamGroupClient() gloov1.UpstreamGroupClient {
	c.cacheUpdateMutex.Lock()
	ret := c.upstreamGroupClient
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
	upstreamGroupClient gloov1.UpstreamGroupClient,
	routeTableClient gatewayv1.RouteTableClient,
) {

	c.cacheUpdateMutex.Lock()

	c.virtualServiceClient = virtualServiceClient
	c.gatewayClient = gatewayClient
	c.upstreamClient = upstreamClient
	c.settingsClient = settingsClient
	c.secretClient = secretClient
	c.artifactClient = artifactClient
	c.proxyClient = proxyClient
	c.upstreamGroupClient = upstreamGroupClient
	c.routeTableClient = routeTableClient

	c.cacheUpdateMutex.Unlock()
}

// Returns a set of clients that use Kubernetes as storage
func NewClientCache(ctx context.Context, settings *gloov1.Settings, cfg *rest.Config, token setup.Token, podNamespace string) (ClientCache, error) {

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

	vsClient, err := gatewayv1.NewVirtualServiceClientWithToken(factoryFor(gatewayv1.VirtualServiceCrd, *cfg, k8sCache, settings, podNamespace), *token)
	if err != nil {
		return nil, err
	}

	upstreamGroupClient, err := gloov1.NewUpstreamGroupClientWithToken(opts.UpstreamGroups, *token)
	if err != nil {
		return nil, err
	}

	gatewayClient, err := gatewayv2.NewGatewayClientWithToken(factoryFor(gatewayv2.GatewayCrd, *cfg, k8sCache, settings, podNamespace), *token)
	if err != nil {
		return nil, err
	}

	proxyClient, err := gloov1.NewProxyClientWithToken(factoryFor(gloov1.ProxyCrd, *cfg, k8sCache, settings, podNamespace), *token)
	if err != nil {
		return nil, err
	}

	settingsClient, err := gloov1.NewSettingsClientWithToken(factoryFor(gloov1.SettingsCrd, *cfg, k8sCache, settings, podNamespace), *token)
	if err != nil {
		return nil, err
	}

	routeTableClient, err := gatewayv1.NewRouteTableClientWithToken(factoryFor(gatewayv1.RouteTableCrd, *cfg, k8sCache, settings, podNamespace), *token)
	if err != nil {
		return nil, err
	}

	// Needed only for the clients backed by the KubeResourceClientFactory
	// so that they register with the cache they share
	if err = registerAll(upstreamClient, vsClient, upstreamGroupClient, gatewayClient, proxyClient, settingsClient, routeTableClient); err != nil {
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
		cacheUpdateMutex: sync.RWMutex{},
	}

	cache.SetCacheState(
		vsClient,
		gatewayClient,
		upstreamClient,
		settingsClient,
		secretClient,
		artifactClient,
		proxyClient,
		upstreamGroupClient,
		routeTableClient,
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

func factoryFor(crd crd.Crd, cfg rest.Config, cache kube.SharedCache, settings *gloov1.Settings, podNamespace string) factory.ResourceClientFactory {
	// TODO refactor this into shareable setup utility for all gloo components.
	writeNamespace := settings.GetDiscoveryNamespace()
	if writeNamespace == "" {
		writeNamespace = podNamespace
	}
	watchNamespaces := utils.ProcessWatchNamespaces(settings.GetWatchNamespaces(), settings.GetDiscoveryNamespace())

	return &factory.KubeResourceClientFactory{
		Crd:                crd,
		Cfg:                &cfg,
		SharedCache:        cache,
		NamespaceWhitelist: watchNamespaces,
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
