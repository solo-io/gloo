package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	clientset "github.com/solo-io/glue/config/watcher/crd/client/clientset/versioned"
	gluescheme "github.com/solo-io/glue/config/watcher/crd/client/clientset/versioned/scheme"
	informers "github.com/solo-io/glue/config/watcher/crd/client/informers/externalversions"
	listers "github.com/solo-io/glue/config/watcher/crd/client/listers/solo.io/v1"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/log"
)

const controllerAgentName = "glue-crd-controller"

// Controller is the controller implementation for Route resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// glueclientset is a clientset for our own API group
	glueclientset clientset.Interface

	routesLister listers.RouteLister
	routesSynced cache.InformerSynced

	upstreamsLister listers.UpstreamLister
	upstreamsSynced cache.InformerSynced

	virtualHostsLister listers.VirtualHostLister
	virtualHostsSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	configs chan *v1.Config
}

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	glueclientset clientset.Interface,
	glueInformerFactory informers.SharedInformerFactory) *Controller {

	// obtain references to shared index informers for the Deployment and Route
	// types.
	routeInformer := glueInformerFactory.Glue().V1().Routes()
	upstreamInformer := glueInformerFactory.Glue().V1().Upstreams()
	virtualHostInformer := glueInformerFactory.Glue().V1().VirtualHosts()

	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	gluescheme.AddToScheme(scheme.Scheme)
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:      kubeclientset,
		glueclientset:      glueclientset,
		routesLister:       routeInformer.Lister(),
		routesSynced:       routeInformer.Informer().HasSynced,
		upstreamsLister:    upstreamInformer.Lister(),
		upstreamsSynced:    upstreamInformer.Informer().HasSynced,
		virtualHostsLister: virtualHostInformer.Lister(),
		virtualHostsSynced: virtualHostInformer.Informer().HasSynced,
		workqueue:          workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Routes"),
		recorder:           recorder,
		configs:            make(chan *v1.Config),
	}

	glog.Info("Setting up event handlers")
	for _, informer := range []cache.SharedInformer{routeInformer.Informer(), upstreamInformer.Informer(), virtualHostInformer.Informer()} {
		// Set up an event handler for when glue resources change
		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueSync("added"),
			UpdateFunc: func(old, new interface{}) {
				controller.enqueueSync("updated")(new)
			},
			DeleteFunc: controller.enqueueSync("deleted"),
		})
	}

	return controller
}

// registers crds to the kube apiserver
func (c *Controller) RegisterCRDs() {

}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Glue controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.routesSynced, c.upstreamsSynced, c.virtualHostsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process glue resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Route resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// enqueueSync takes a Glue resource and converts it into a namespace/name
// string which is then put onto the work queue.
func (c *Controller) enqueueSync(event string) func(interface{}) {
	return func(obj interface{}) {
		var key string
		var err error
		if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
			runtime.HandleError(err)
			return
		}
		log.Printf("crd event: %s: %s", event, key)
		c.workqueue.AddRateLimited(key)
	}
}

// syncHandler retrieves all crds from kubernetes
// and constructs an updated version of the config
// we really don't care what resource type we were passed.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	log.Printf("updating config after item %v changed", name)

	c.virtualHostsLister.List(labels.Everything())
	routeList, err := c.routesLister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("error retrieving routes: %v", err)
	}
	upstreamList, err := c.upstreamsLister.List(labels.Everything())
	if err != nil {
		return fmt.Errorf("error retrieving upstreams: %v", err)
	}
	vHostList, err := c.virtualHostsLister.List(labels.Everything())
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
}

func (c *Controller) Configs() <-chan *v1.Config {
	return c.configs
}
