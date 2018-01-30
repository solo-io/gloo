package configwatcher

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/sample-controller/pkg/signals"

	clientset "github.com/solo-io/glue/implemented_modules/kube/configwatcher/crd/client/clientset/versioned"
	informers "github.com/solo-io/glue/implemented_modules/kube/configwatcher/crd/client/informers/externalversions"
	listers "github.com/solo-io/glue/implemented_modules/kube/configwatcher/crd/client/listers/solo.io/v1"
	"github.com/solo-io/glue/implemented_modules/kube/pkg/controller"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
)

type crdController struct {
	configs chan *v1.Config
	errors  chan error
}

func newCrdController(cfg *rest.Config, resyncDuration time.Duration) (*crdController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	glueClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create glue clientset: %v", err)
	}

	glueInformerFactory := informers.NewSharedInformerFactory(glueClient, resyncDuration)

	routeInformer := glueInformerFactory.Glue().V1().Routes()
	upstreamInformer := glueInformerFactory.Glue().V1().Upstreams()
	virtualHostInformer := glueInformerFactory.Glue().V1().VirtualHosts()

	kubeController := controller.NewController("glue-crd-controller", kubeClient,
		routeInformer.Informer(),
		upstreamInformer.Informer(),
		virtualHostInformer.Informer())

	ctrl := &crdController{
		configs: make(chan *v1.Config),
		errors:  make(chan error),
	}

	kubeController.AddEventHandler(controller.Added, func(namespace, name string, _ interface{}) {
		ctrl.syncConfig(namespace, name, routeInformer.Lister(), upstreamInformer.Lister(), virtualHostInformer.Lister())
	})
	kubeController.AddEventHandler(controller.Updated, func(namespace, name string, _ interface{}) {
		ctrl.syncConfig(namespace, name, routeInformer.Lister(), upstreamInformer.Lister(), virtualHostInformer.Lister())
	})
	kubeController.AddEventHandler(controller.Deleted, func(namespace, name string, _ interface{}) {
		ctrl.syncConfig(namespace, name, routeInformer.Lister(), upstreamInformer.Lister(), virtualHostInformer.Lister())
	})

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	go glueInformerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

func (c *crdController) syncConfig(namespace, name string,
	routesLister listers.RouteLister,
	upstreamsLister listers.UpstreamLister,
	virtualHostsLister listers.VirtualHostLister) error {
	if err := func() error {
		log.Printf("syncing config after item %v/%v changed", namespace, name)

		routeList, err := routesLister.List(labels.Everything())
		if err != nil {
			return fmt.Errorf("error retrieving routes: %v", err)
		}
		upstreamList, err := upstreamsLister.List(labels.Everything())
		if err != nil {
			return fmt.Errorf("error retrieving upstreams: %v", err)
		}
		vHostList, err := virtualHostsLister.List(labels.Everything())
		if err != nil {
			return fmt.Errorf("error retrieving virtualhosts: %v", err)
		}
		var routes []v1.Route
		for _, route := range routeList {
			r := v1.Route(route.Spec)
			routes = append(routes, r)
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
			Routes:       routes,
			Upstreams:    upstreams,
			VirtualHosts: vHosts,
		}
		return nil
	}(); err != nil {
		c.errors <- err
	}
	return nil
}

func (c *crdController) Config() <-chan *v1.Config {
	return c.configs
}

func (c *crdController) Error() <-chan error {
	return c.errors
}
