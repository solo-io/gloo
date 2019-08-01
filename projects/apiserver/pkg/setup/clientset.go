package setup

import (
	"context"
	"log"
	"sync"

	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/graph"
	v1 "k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	consulapi "github.com/hashicorp/consul/api"
	vaultapi "github.com/hashicorp/vault/api"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// PerTokenClientsets contains global settings and user-specific resource clients
// clients is a map from user token to resource clients
// the token is also used for authorizing actions on the resource clients
type PerTokenClientsets struct {
	lock    sync.RWMutex
	clients map[string]*ClientSet
}

func NewPerTokenClientsets() PerTokenClientsets {
	return PerTokenClientsets{
		clients: make(map[string]*ClientSet),
	}
}

func (ptc PerTokenClientsets) ClientsetForToken(ctx context.Context, settings *gloov1.Settings, token string) (*ClientSet, error) {
	ptc.lock.Lock()
	defer ptc.lock.Unlock()
	clientsetForToken, ok := ptc.clients[token]
	if ok {
		return clientsetForToken, nil
	}

	clientset, err := NewClientSet(ctx, settings, token)
	if err != nil {
		return nil, err
	}
	ptc.clients[token] = clientset
	return clientset, nil
}

// ClientSet is a collection of all the exposed resource clients
type ClientSet struct {
	gatewayv1.VirtualServiceClient
	gloov1.UpstreamClient
	gloov1.SettingsClient
	gloov1.SecretClient
	gloov1.ArtifactClient
	corev1.CoreV1Interface
}

func (c ClientSet) NewResolvers() graph.ResolverRoot {
	return graphql.NewResolvers(
		c.UpstreamClient,
		c.ArtifactClient,
		c.SettingsClient,
		c.SecretClient,
		c.VirtualServiceClient,
		c.CoreV1Interface)
}

// Returns a set of clients that use Kubernetes as storage
func NewClientSet(ctx context.Context, settings *gloov1.Settings, token string) (*ClientSet, error) {

	// When running in-cluster, this configuration will hold a token associated with the pod service account
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	// TODO: temporary solution to bypass authentication.
	if token == "" {
		if config.SkipAuth == "" {
			log.Panic("token is empty and auth is not bypassed. Should never happen.")
		}
		// When we want to skip auth, we use the token associated with the pod service account
		token = cfg.BearerToken
	}

	kubeClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// New shared cache
	cache := kube.NewKubeCache(context.TODO())

	var clientset kubernetes.Interface
	memCache := memory.NewInMemoryResourceCache()
	opts, err := constructOpts(ctx, &clientset, cache, nil, nil, memCache, settings)
	if err != nil {
		return nil, err
	}

	upstreamClient, err := gloov1.NewUpstreamClientWithToken(opts.Upstreams, token)
	if err != nil {
		return nil, err
	}

	vsClient, err := gatewayv1.NewVirtualServiceClientWithToken(factoryFor(gatewayv1.VirtualServiceCrd, *cfg, cache), token)
	if err != nil {
		return nil, err
	}

	settingsClient, err := gloov1.NewSettingsClientWithToken(factoryFor(gloov1.SettingsCrd, *cfg, cache), token)
	if err != nil {
		return nil, err
	}

	// Needed only for the clients backed by the KubeResourceClientFactory
	// so that they register with the cache they share
	if err = registerAll(upstreamClient, vsClient, settingsClient); err != nil {
		return nil, err
	}

	// replace this with the gloo factory
	secretClient, err := gloov1.NewSecretClientWithToken(opts.Secrets, token)
	if err != nil {
		return nil, err
	}

	artifactClient, err := gloov1.NewArtifactClientWithToken(opts.Artifacts, token)
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		UpstreamClient:       upstreamClient,
		VirtualServiceClient: vsClient,
		SettingsClient:       settingsClient,
		SecretClient:         secretClient,
		ArtifactClient:       artifactClient,
		CoreV1Interface:      kubeClientset.CoreV1(),
	}, nil
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
		SkipCrdCreation:    true,
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
