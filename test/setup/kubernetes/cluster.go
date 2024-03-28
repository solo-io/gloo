package kubernetes

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/go-utils/contextutils"
)

type ClusterMapping struct {
	clusters sync.Map
}

func NewClusterMapping() *ClusterMapping {
	return &ClusterMapping{}
}

func (c *ClusterMapping) UpdateClusterMapping(
	ctx context.Context,
	name, kubecontext, kubeconfig string,
	client kubernetes.Interface,
	controller client.Client,
) error {
	logger := contextutils.LoggerFrom(ctx)

	logger.Info("updating cluster mapping", zap.String("cluster", name), zap.String("config", kubeconfig), zap.String("context", kubecontext))

	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	cluster := &Cluster{
		name:        name,
		client:      client,
		controller:  controller,
		kubeconfig:  kubeconfig,
		kubecontext: kubecontext,
		logger:      *logger.With(zap.String("cluster", name)),
	}

	c.clusters.Store(name, cluster)
	return nil
}

func (c *ClusterMapping) GetClusterInfo(name string) *Cluster {
	if v, ok := c.clusters.Load(name); ok {
		return v.(*Cluster)
	}
	return nil
}

type Cluster struct {
	name        string
	kubeconfig  string // Path to kubeconfig file
	kubecontext string
	client      kubernetes.Interface
	controller  client.Client
	logger      zap.SugaredLogger
}

func (c *Cluster) GetName() string {
	return c.name
}

func (c *Cluster) GetKubeConfig() string {
	return c.kubeconfig
}

func (c *Cluster) GetKubeContext() string {
	return c.kubecontext
}

func (c *Cluster) GetController() client.Client {
	return c.controller
}

func (c *Cluster) GetKubernetes() kubernetes.Interface {
	return c.client
}

func (c *Cluster) GetLogger() zap.SugaredLogger {
	return c.logger
}
