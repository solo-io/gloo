package server

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned/typed/solo.io/v1"
	solov1 "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilrt "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	// TODO (ashish) convert these into configurations
	maxRetries  = 5
	resyncDelay = 5 * time.Minute
)

// Handler represents the interface for anything that is
// interested in listening to Upstream events
// TODO(ashish) - rename to Subscriber
type Handler interface {
	Update(*solov1.Upstream)
	Remove(*solov1.Upstream)
}

type controller struct {
	repo     v1.UpstreamInterface
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	handlers []Handler
}

func newController(upstreamRepo v1.UpstreamInterface) *controller {
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return upstreamRepo.List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return upstreamRepo.Watch(options)
		},
	}
	informer := cache.NewSharedIndexInformer(lw, &solov1.Upstream{}, resyncDelay, cache.Indexers{})
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				queue.Add(key)
			}
		},
	})

	return &controller{repo: upstreamRepo, queue: queue, informer: informer}
}

func (c *controller) AddHandler(h Handler) {
	c.handlers = append(c.handlers, h)
}

func (c *controller) Run(stop <-chan struct{}) {
	defer utilrt.HandleCrash()
	defer c.queue.ShutDown()

	go c.informer.Run(stop)
	// wait for cache sync before starting
	if !cache.WaitForCacheSync(stop, c.informer.HasSynced) {
		utilrt.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	}
	go wait.Until(c.runWorker, time.Second, stop)
	log.Println("Controller started")
}

func (c *controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *controller) processNextItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	err := c.processItem(key.(string))
	if err == nil {
		c.queue.Forget(key)
	} else if c.queue.NumRequeues(key) < maxRetries {
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
		utilrt.HandleError(err)
	}
	return true
}

func (c *controller) processItem(key string) error {
	log.Println("processing: " + key)
	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return errors.Wrapf(err, "error fetching object with key %s from store", key)
	}
	if !exists {
		for _, h := range c.handlers {
			h.Remove(obj.(*solov1.Upstream))
		}
		return nil
	}
	for _, h := range c.handlers {
		h.Update(obj.(*solov1.Upstream))
	}
	return nil
}

func (c *controller) get(key string) (*solov1.Upstream, bool, error) {
	o, b, e := c.informer.GetIndexer().GetByKey(key)
	if e != nil {
		return nil, false, e
	}
	if !b {
		return nil, false, nil
	}
	return o.(*solov1.Upstream), b, nil
}

func (c *controller) set(u *solov1.Upstream) error {
	_, err := c.repo.Update(u)
	return err
}
