package kube

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/sample-controller/pkg/signals"

	"strings"

	"github.com/solo-io/glue/adapters/kube/controller"
	"github.com/solo-io/glue/discovery"
	"github.com/solo-io/glue/pkg/log"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/listers/core/v1"
)

type discoveryController struct {
	clusters        chan discovery.Clusters
	errors          chan error
	serviceLister   v1.ServiceLister
	endpointsLister v1.EndpointsLister
	serviceRefs     []serviceRef
}

type serviceRef struct {
	namespace, name string
}

func newDiscoveryController(cfg *rest.Config, resyncDuration time.Duration) (*discoveryController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	serviceInformer := informerFactory.Core().V1().Services()
	endpointsLister := informerFactory.Core().V1().Endpoints()

	kubeController := controller.NewController("glue-endpoints-controller", kubeClient,
		serviceInformer.Informer(), endpointsLister.Informer())

	ctrl := &discoveryController{
		clusters:        make(chan discovery.Clusters),
		errors:          make(chan error),
		serviceLister:   serviceInformer.Lister(),
		endpointsLister: endpointsLister.Lister(),
	}

	kubeController.AddEventHandler(controller.Added, func(_, _ string, _ interface{}) {
		ctrl.getClusters()
	})
	kubeController.AddEventHandler(controller.Updated, func(namespace, name string, _ interface{}) {
		ctrl.getClusters()
	})
	kubeController.AddEventHandler(controller.Deleted, func(namespace, name string, _ interface{}) {
		ctrl.getClusters()
	})

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	go informerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

// triggers an update
func (c *discoveryController) UpdateClusterRefs(hostnames []string) {
	var serviceRefs []serviceRef
	for _, hostname := range hostnames {
		namespace, name, err := parseHostname(hostname)
		if err != nil {
			runtime.HandleError(fmt.Errorf("failed to parse hostname %v", err))
			continue
		}
		serviceRefs = append(serviceRefs, serviceRef{
			namespace: namespace,
			name:      name,
		})
	}
	c.syncClusters()
}

func (c *discoveryController) Clusters() <-chan discovery.Clusters {
	return c.clusters
}

func (c *discoveryController) Error() <-chan error {
	return c.errors
}

// pushes clusters or error to channel
func (c *discoveryController) syncClusters() {
	clusters, err := c.getClusters()
	if err != nil {
		c.errors <- err
		return
	}
	// ignore empty configs / no clusters to watch
	if len(clusters) == 0 {
		return
	}
	c.clusters <- clusters
}

// retrieves clusters from kubernetes
func (c *discoveryController) getClusters() (discovery.Clusters, error) {
	serviceList, err := c.serviceLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving routes: %v", err)
	}
	clusters := make(discovery.Clusters)
	for _, ref := range c.serviceRefs {
		var clusterFound bool
		for _, service := range serviceList {
			if service.Name == ref.name && service.Namespace == ref.namespace {
				log.Printf("updated cluster %s", ref)
				service.Spec
				clusters[ref] = service.Data
				break
			}
		}
		if !clusterFound {
			runtime.HandleError(fmt.Errorf("cluster for service %v not found", ref))
		}
	}
	return clusters, nil
}

// parseHostname extracts service name and namespace from the service hostname
func parseHostname(hostname string) (name string, namespace string, err error) {
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		err = fmt.Errorf("missing service name and namespace from the service hostname %q", hostname)
		return
	}
	name = parts[0]
	namespace = parts[1]
	return
}
