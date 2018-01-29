package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	"github.com/solo-io/glue/pkg/log"
)

// Controller is the controller implementation for Route resources
type Controller struct {
	name string

	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface

	syncFuncs []cache.InformerSynced
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder

	onAddFuncs    []func(namespace, name string, obj interface{})
	onUpdateFuncs []func(namespace, name string, obj interface{})
	onDeleteFuncs []func(namespace, name string, obj interface{})
}

// NewController returns a new sample controller
func NewController(
	controllerAgentName string,
	kubeclientset kubernetes.Interface,
	informers ...cache.SharedInformer) *Controller {

	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	var hasSyncedFuncs []cache.InformerSynced

	controller := &Controller{
		syncFuncs: hasSyncedFuncs,
		name:      controllerAgentName,
		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Routes"),
		recorder:  recorder,
	}

	glog.Info("Setting up event handlers")
	for _, informer := range informers {
		hasSyncedFuncs = append(hasSyncedFuncs, informer.HasSynced)

		// Set up an event handler for when glue resources change
		informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: controller.enqueueSync(Added),
			UpdateFunc: func(old, new interface{}) {
				controller.enqueueSync(Updated)(new)
			},
			DeleteFunc: controller.enqueueSync(Deleted),
		})
	}

	return controller
}

func (c *Controller) AddEventHandler(event Event, handler func(namespace, name string, obj interface{})) {
	switch event {
	case Added:
		c.onAddFuncs = append(c.onAddFuncs, handler)
	case Updated:
		c.onUpdateFuncs = append(c.onUpdateFuncs, handler)
	case Deleted:
		c.onDeleteFuncs = append(c.onDeleteFuncs, handler)
	}
}

type Event int

const (
	Added Event = iota
	Updated
	Deleted
)

type wrapped struct {
	event Event
	key   string
	obj   interface{}
}

func (c *Controller) enqueueSync(event Event) func(interface{}) {
	return func(obj interface{}) {
		var key string
		var err error
		if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
			runtime.HandleError(err)
			return
		}
		log.Printf("%s event: %s: %s", c.name, event, key)
		c.workqueue.AddRateLimited(wrapped{
			event: event,
			key:   key,
			obj:   obj,
		})
	}
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting %v controller", c.name)

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, []cache.InformerSynced(c.syncFuncs)...); !ok {
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
		var w wrapped
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if w, ok = obj.(wrapped); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected wrapped type in workqueue but got %#v", obj))
			return nil
		}
		namespace, name, _ := cache.SplitMetaNamespaceKey(w.key)
		switch w.event {
		case Added:
			for _, onAdd := range c.onAddFuncs {
				onAdd(namespace, name, w.obj)
			}
		case Updated:
			for _, onUpdate := range c.onUpdateFuncs {
				onUpdate(namespace, name, w.obj)
			}
		case Deleted:
			for _, onDelete := range c.onDeleteFuncs {
				onDelete(namespace, name, w.obj)
			}
		}

		c.workqueue.Forget(obj)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}
