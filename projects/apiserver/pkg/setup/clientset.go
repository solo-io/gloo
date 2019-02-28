package setup

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/graph"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/auth"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql"
	vcsv1 "github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"
	v1 "k8s.io/api/core/v1"
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

func (ptc PerTokenClientsets) ClientsetForToken(token string) (*ClientSet, error) {
	ptc.lock.Lock()
	defer ptc.lock.Unlock()
	clientsetForToken, ok := ptc.clients[token]
	if ok {
		return clientsetForToken, nil
	}

	// TODO: temporary flag to switch to VCS clientset
	var clientset *ClientSet
	var err error
	if os.Getenv("VCS") == "1" {
		clientset, err = NewTempClientSet(token)
	} else {
		// Use regular clientset, to support current UI
		clientset, err = NewClientSet(token)
	}

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
func NewClientSet(token string) (*ClientSet, error) {

	// TODO: temporary solution to bypass authentication.
	if token == "" {
		if config.SkipAuth == "" {
			log.Panic("token is empty and auth is not bypassed. Should never happen.")
		}
		return newAdminClientSet()
	}

	// When running in-cluster, this configuration will hold a token associated with the pod service account
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	kubeClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// New shared cache
	cache := kube.NewKubeCache(context.TODO())

	kubeCoreCache, err := corecache.NewKubeCoreCache(context.TODO(), kubeClientset)
	if err != nil {
		return nil, err
	}

	upstreamClient, err := gloov1.NewUpstreamClientWithToken(factoryFor(gloov1.UpstreamCrd, *cfg, cache), token)
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

	secretClient, err := gloov1.NewSecretClientWithToken(&factory.KubeSecretClientFactory{
		Clientset: kubeClientset,
		Cache:     kubeCoreCache,
	}, token)
	if err != nil {
		return nil, err
	}

	artifactClient, err := gloov1.NewArtifactClientWithToken(&factory.KubeConfigMapClientFactory{
		Clientset: kubeClientset,
		Cache:     kubeCoreCache,
	}, token)
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

// Returns a set of clients that use a Changeset as storage
func NewChangesetClientSet(token string) (*ClientSet, error) {

	// When running in-cluster, this configuration will hold a token associated with the pod service account
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	kubeClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Validate bearer token and retrieve associated user information
	username, err := auth.GetUsername(kubeClientset.AuthenticationV1(), token)
	if err != nil {
		return nil, err
	}

	// Create a kubernetes client for changesets
	changesetClient, err := vcsv1.NewChangeSetClientWithToken(&factory.KubeResourceClientFactory{
		Crd:                vcsv1.ChangeSetCrd,
		Cfg:                cfg,
		SharedCache:        kube.NewKubeCache(context.TODO()),
		SkipCrdCreation:    true,
		NamespaceWhitelist: []string{defaults.GlooSystem},
	}, token)
	if err = changesetClient.Register(); err != nil {
		return nil, err
	}

	// Clients built on top of this factory will use the changeset with the given name as storage
	changesetClientFactory := &vcsv1.ChangesetResourceClientFactory{
		ChangesetClient: changesetClient,
		ChangesetName:   username,
	}

	vsClient, err := gatewayv1.NewVirtualServiceClient(changesetClientFactory)
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		VirtualServiceClient: vsClient,
	}, nil
}

// TODO: temporary hybrid clientset that writes only virtual services to a changeset
func NewTempClientSet(token string) (*ClientSet, error) {

	// When running in-cluster, this configuration will hold a token associated with the pod service account
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	kubeClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// New shared cache
	cache := kube.NewKubeCache(context.TODO())
	kubeCoreCache, err := corecache.NewKubeCoreCache(context.TODO(), kubeClientset)
	if err != nil {
		return nil, err
	}

	upstreamClient, err := gloov1.NewUpstreamClientWithToken(factoryFor(gloov1.UpstreamCrd, *cfg, cache), token)
	if err != nil {
		return nil, err
	}

	settingsClient, err := gloov1.NewSettingsClientWithToken(factoryFor(gloov1.SettingsCrd, *cfg, cache), token)
	if err != nil {
		return nil, err
	}

	// Needed only for the clients backed by a KubeResourceClientFactory to register with the cache they share
	if err = registerAll(upstreamClient, settingsClient); err != nil {
		return nil, err
	}

	secretClient, err := gloov1.NewSecretClientWithToken(&factory.KubeSecretClientFactory{
		Clientset: kubeClientset,
		Cache:     kubeCoreCache,
	}, token)
	if err != nil {
		return nil, err
	}

	artifactClient, err := gloov1.NewArtifactClientWithToken(&factory.KubeConfigMapClientFactory{
		Clientset: kubeClientset,
		Cache:     kubeCoreCache,
	}, token)
	if err != nil {
		return nil, err
	}

	// validate bearer token and retrieve associated user information
	username, err := auth.GetUsername(kubeClientset.AuthenticationV1(), token)
	if err != nil {
		return nil, err
	}

	// Create a kubernetes client for changesets
	changesetClient, err := vcsv1.NewChangeSetClientWithToken(&factory.KubeResourceClientFactory{
		Crd:                vcsv1.ChangeSetCrd,
		Cfg:                cfg,
		SharedCache:        kube.NewKubeCache(context.TODO()),
		SkipCrdCreation:    true,
		NamespaceWhitelist: []string{defaults.GlooSystem},
	}, token)
	if err = changesetClient.Register(); err != nil {
		return nil, err
	}

	// Clients built on top of this factory will use the changeset with the given name as storage
	changesetClientFactory := &vcsv1.ChangesetResourceClientFactory{
		ChangesetClient: changesetClient,
		ChangesetName:   username,
	}

	vsClient, err := gatewayv1.NewVirtualServiceClient(changesetClientFactory)

	return &ClientSet{
		UpstreamClient:       upstreamClient,
		VirtualServiceClient: vsClient,
		SettingsClient:       settingsClient,
		SecretClient:         secretClient,
		ArtifactClient:       artifactClient,
	}, nil
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

// Returns a set of clients that use Kubernetes as storage
// Uses the InClusterConfig k8s configuration
func newAdminClientSet() (*ClientSet, error) {

	// When running in-cluster, this configuration will hold a token associated with the pod service account
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	kubeClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	// New shared cache
	cache := kube.NewKubeCache(context.TODO())
	kubeCoreCache, err := corecache.NewKubeCoreCache(context.TODO(), kubeClientset)
	if err != nil {
		return nil, err
	}

	upstreamClient, err := gloov1.NewUpstreamClient(factoryFor(gloov1.UpstreamCrd, *cfg, cache))
	if err != nil {
		return nil, err
	}

	vsClient, err := gatewayv1.NewVirtualServiceClient(factoryFor(gatewayv1.VirtualServiceCrd, *cfg, cache))
	if err != nil {
		return nil, err
	}

	settingsClient, err := gloov1.NewSettingsClient(factoryFor(gloov1.SettingsCrd, *cfg, cache))
	if err != nil {
		return nil, err
	}

	// Needed only for the clients backed by the KubeResourceClientFactory
	// so that they register with the cache they share
	if err = registerAll(upstreamClient, vsClient, settingsClient); err != nil {
		return nil, err
	}

	secretClient, err := gloov1.NewSecretClient(&factory.KubeSecretClientFactory{
		Clientset: kubeClientset,
		Cache:     kubeCoreCache,
	})
	if err != nil {
		return nil, err
	}

	artifactClient, err := gloov1.NewArtifactClient(&factory.KubeConfigMapClientFactory{
		Clientset: kubeClientset,
		Cache:     kubeCoreCache,
	})
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
