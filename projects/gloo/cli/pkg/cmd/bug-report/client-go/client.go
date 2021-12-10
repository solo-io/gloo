package client_go

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"io"
	"k8s.io/client-go/metadata"
	"k8s.io/client-go/metadata/metadatainformer"
	"net/http"
	"sync"

	kubeExtClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/api/meta"
	kubeVersion "k8s.io/apimachinery/pkg/version"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/openapi"
	"k8s.io/kubectl/pkg/validation"
)

// Client is a helper for common Kubernetes client operations. This contains various different kubernetes
// clients using a shared config. It is expected that all of Istiod can share the same set of clients and
// informers. Sharing informers is especially important for load on the API server/Istiod itself.
type Client interface {
	// TODO - stop embedding this, it will conflict with future additions. Use Kube() instead is preferred
	kubernetes.Interface
	// RESTConfig returns the Kubernetes rest.Config used to configure the clients.
	RESTConfig() *rest.Config

	// Ext returns the API extensions client.
	Ext() kubeExtClient.Interface

	// EnvoyDo makes an http request to the Envoy in the specified pod.
	EnvoyDo(ctx context.Context, podName, podNamespace, method, path string) ([]byte, error)

	// Kube returns the core kube client
	Kube() kubernetes.Interface

	KubeClientSet() *kubernetes.Clientset

	// Dynamic client.
	Dynamic() dynamic.Interface

	// Metadata returns the Metadata kube client.
	Metadata() metadata.Interface

	// KubeInformer returns an informer for core kube client
	KubeInformer() informers.SharedInformerFactory

	// DynamicInformer returns an informer for dynamic client
	DynamicInformer() dynamicinformer.DynamicSharedInformerFactory

	// MetadataInformer returns an informer for metadata client
	MetadataInformer() metadatainformer.SharedInformerFactory

	// RunAndWait starts all informers and waits for their caches to sync.
	// Warning: this must be called AFTER .Informer() is called, which will register the informer.
	RunAndWait(stop <-chan struct{})

	// GetKubernetesVersion returns the Kubernetes server version
	GetKubernetesVersion() (*kubeVersion.Info, error)
}

// NewClient creates a Kubernetes client from the given rest config.
func NewClient(clientConfig clientcmd.ClientConfig) (Client, error) {
	return newClientInternal(newClientFactory(clientConfig), "")
}

const resyncInterval = 0

// newClientInternal creates a Kubernetes client from the given factory.
func newClientInternal(clientFactory util.Factory, revision string) (*client, error) {
	var c client
	var err error

	c.clientFactory = clientFactory

	c.config, err = clientFactory.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	c.kubeClient, err = clientFactory.KubernetesClientSet()
	if err != nil {
		return nil, err
	}
	c.restClient, err = clientFactory.RESTClient()
	if err != nil {
		return nil, err
	}

	c.discoveryClient, err = clientFactory.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}
	c.mapper = restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(c.discoveryClient))

	c.Interface, err = kubernetes.NewForConfig(c.config)
	c.kube = c.Interface
	if err != nil {
		return nil, err
	}
	c.kubeInformer = informers.NewSharedInformerFactory(c.Interface, resyncInterval)

	c.metadata, err = metadata.NewForConfig(c.config)
	if err != nil {
		return nil, err
	}
	c.metadataInformer = metadatainformer.NewSharedInformerFactory(c.metadata, resyncInterval)

	c.dynamic, err = dynamic.NewForConfig(c.config)
	if err != nil {
		return nil, err
	}
	c.dynamicInformer = dynamicinformer.NewDynamicSharedInformerFactory(c.dynamic, resyncInterval)

	c.extSet, err = kubeExtClient.NewForConfig(c.config)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Client is a helper wrapper around the Kube RESTClient for istioctl -> Pilot/Envoy/Mesh related things
type client struct {
	kubernetes.Interface

	clientFactory util.Factory
	config        *rest.Config

	extSet kubeExtClient.Interface

	kube         kubernetes.Interface
	kubeInformer informers.SharedInformerFactory

	dynamic         dynamic.Interface
	dynamicInformer dynamicinformer.DynamicSharedInformerFactory

	metadata         metadata.Interface
	metadataInformer metadatainformer.SharedInformerFactory

	// If enable, will wait for cache syncs with extremely short delay. This should be used only for tests
	fastSync               bool
	informerWatchesPending *atomic.Int32

	kubeClient *kubernetes.Clientset

	// These may be set only when creating an extended client.
	restClient      *rest.RESTClient
	discoveryClient discovery.CachedDiscoveryInterface
	mapper          meta.RESTMapper

	versionOnce sync.Once
	version     *kubeVersion.Info
	versionErr  error
}

// newClientFactory creates a new util.Factory from the given clientcmd.ClientConfig.
func newClientFactory(clientConfig clientcmd.ClientConfig) util.Factory {
	out := &clientFactory{
		clientConfig: clientConfig,
	}

	out.factory = util.NewFactory(out)
	return out
}

// clientFactory implements the kubectl util.Factory, which is provides access to various k8s clients.
type clientFactory struct {
	clientConfig clientcmd.ClientConfig
	factory      util.Factory
}

func (c *clientFactory) ToRESTConfig() (*rest.Config, error) {
	restConfig, err := c.clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}
	return SetRestDefaults(restConfig), nil
}

func (c *clientFactory) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	restConfig, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	d, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return memory.NewMemCacheClient(d), nil
}

func (c *clientFactory) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := c.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func (c *clientFactory) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return c.clientConfig
}

func (c *clientFactory) DynamicClient() (dynamic.Interface, error) {
	restConfig, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(restConfig)
}

func (c *clientFactory) KubernetesClientSet() (*kubernetes.Clientset, error) {
	restConfig, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(restConfig)
}

func (c *clientFactory) RESTClient() (*rest.RESTClient, error) {
	return c.factory.RESTClient()
}

func (c *clientFactory) NewBuilder() *resource.Builder {
	return c.factory.NewBuilder()
}

func (c *clientFactory) ClientForMapping(mapping *meta.RESTMapping) (resource.RESTClient, error) {
	return c.factory.ClientForMapping(mapping)
}

func (c *clientFactory) UnstructuredClientForMapping(mapping *meta.RESTMapping) (resource.RESTClient, error) {
	return c.factory.UnstructuredClientForMapping(mapping)
}

func (c *clientFactory) Validator(validate bool) (validation.Schema, error) {
	return c.factory.Validator(validate)
}

func (c *clientFactory) OpenAPISchema() (openapi.Resources, error) {
	return c.factory.OpenAPISchema()
}

func (c *client) RESTConfig() *rest.Config {
	if c.config == nil {
		return nil
	}
	cpy := *c.config
	return &cpy
}

func (c *client) Ext() kubeExtClient.Interface {
	return c.extSet
}

func (c *client) Dynamic() dynamic.Interface {
	return c.dynamic
}

func (c *client) Kube() kubernetes.Interface {
	return c.kube
}

func (c *client) Metadata() metadata.Interface {
	return c.metadata
}

func (c *client) KubeInformer() informers.SharedInformerFactory {
	return c.kubeInformer
}

func (c *client) DynamicInformer() dynamicinformer.DynamicSharedInformerFactory {
	return c.dynamicInformer
}

func (c *client) KubeClientSet() *kubernetes.Clientset {
	//p, _ := c.kubeClient.CoreV1().Pods("hi").
	return c.kubeClient
}

func (c *client) EnvoyDo(ctx context.Context, podName, podNamespace, method, path string) ([]byte, error) {
	return c.portForwardRequest(ctx, podName, podNamespace, method, path, 19000)
}

func (c *client) portForwardRequest(ctx context.Context, podName, podNamespace, method, path string, port int) ([]byte, error) {
	formatError := func(err error) error {
		return fmt.Errorf("failure running port forward process: %v", err)
	}

	fw, err := c.NewPortForwarder(podName, podNamespace, "127.0.0.1", 0, port)
	if err != nil {
		return nil, err
	}
	if err = fw.Start(); err != nil {
		return nil, formatError(err)
	}
	defer fw.Close()
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s/%s", fw.Address(), path), nil)
	if err != nil {
		return nil, formatError(err)
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, formatError(err)
	}
	defer func() { _ = resp.Body.Close() }()
	out, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, formatError(err)
	}

	return out, nil
}

func (c *client) NewPortForwarder(podName, ns, localAddress string, localPort int, podPort int) (PortForwarder, error) {
	return newPortForwarder(c.config, podName, ns, localAddress, localPort, podPort)
}

func (c *client) MetadataInformer() metadatainformer.SharedInformerFactory {
	return c.metadataInformer
}

// RunAndWait starts all informers and waits for their caches to sync.
// Warning: this must be called AFTER .Informer() is called, which will register the informer.
func (c *client) RunAndWait(stop <-chan struct{}) {
	panic("run and wait not impl")
}

func (c *client) GetKubernetesVersion() (*kubeVersion.Info, error) {
	c.versionOnce.Do(func() {
		v, err := c.Discovery().ServerVersion()
		c.version = v
		c.versionErr = err
	})
	return c.version, c.versionErr
}
