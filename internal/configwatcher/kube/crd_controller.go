package kube

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	clientset "github.com/solo-io/glue/internal/configwatcher/kube/crd/client/clientset/versioned"
	informers "github.com/solo-io/glue/internal/configwatcher/kube/crd/client/informers/externalversions"
	listers "github.com/solo-io/glue/internal/configwatcher/kube/crd/client/listers/solo.io/v1"
	"github.com/solo-io/glue/internal/pkg/kube/controller"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
)

type crdController struct {
	configs chan *v1.Config
	errors  chan error

	upstreamsLister    listers.UpstreamLister
	virtualHostsLister listers.VirtualHostLister
}

func newCrdController(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) (*crdController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	glueClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create glue clientset: %v", err)
	}

	glueInformerFactory := informers.NewSharedInformerFactory(glueClient, resyncDuration)

	upstreamInformer := glueInformerFactory.Glue().V1().Upstreams()
	virtualHostInformer := glueInformerFactory.Glue().V1().VirtualHosts()

	ctrl := &crdController{
		configs:            make(chan *v1.Config),
		errors:             make(chan error),
		upstreamsLister:    upstreamInformer.Lister(),
		virtualHostsLister: virtualHostInformer.Lister(),
	}

	kubeController := controller.NewController("glue-crd-controller", kubeClient,
		ctrl.syncConfig,
		upstreamInformer.Informer(),
		virtualHostInformer.Informer())

	go glueInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

func (c *crdController) syncConfig(namespace, name string, _ interface{}) {
	if err := func() error {
		log.Debugf("syncing config after item %v/%v changed", namespace, name)

		upstreamList, err := c.upstreamsLister.List(labels.Everything())
		if err != nil {
			return fmt.Errorf("error retrieving upstreams: %v", err)
		}
		vHostList, err := c.virtualHostsLister.List(labels.Everything())
		if err != nil {
			return fmt.Errorf("error retrieving virtualhosts: %v", err)
		}
		var upstreams []v1.Upstream
		for _, upstream := range upstreamList {
			u := v1.Upstream(upstream.Spec)
			upstreams = append(upstreams, u)
		}
		var vHosts []v1.VirtualHost
		for _, vHost := range vHostList {
			v := v1.VirtualHost(vHost.Spec)
			vHosts = append(vHosts, v)
		}
		log.Debugf("config updated")
		c.configs <- &v1.Config{
			Upstreams:    upstreams,
			VirtualHosts: vHosts,
		}
		return nil
	}(); err != nil {
		c.errors <- err
	}
}

func (c *crdController) Config() <-chan *v1.Config {
	return c.configs
}

func (c *crdController) Error() <-chan error {
	return c.errors
}
