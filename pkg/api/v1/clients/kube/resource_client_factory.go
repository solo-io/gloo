package kube

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	"github.com/solo-io/kubecontroller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	kubewatch "k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	"github.com/solo-io/solo-kit/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	MLists   = stats.Int64("kube/lists", "The number of lists", "1")
	MWatches = stats.Int64("kube/lists", "The  number of watches", "1")

	KeyKind, _ = tag.NewKey("kind")

	ListCountView = &view.View{
		Name:        "kube/lists-count",
		Measure:     MLists,
		Description: "The number of list calls",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			KeyKind,
		},
	}
	WatchCountView = &view.View{
		Name:        "kube/watches-count",
		Measure:     MWatches,
		Description: "The number of list calls",
		Aggregation: view.Count(),
		TagKeys: []tag.Key{
			KeyKind,
		},
	}
)

func init() {
	view.Register(ListCountView)
}

type ResourceLister interface {
	List(selector labels.Selector) (ret []*v1.Resource, err error)
}

type ResourceClientSharedInformerFactory struct {
	initError error

	lock          sync.Mutex
	defaultResync time.Duration

	informers map[reflect.Type]cache.SharedIndexInformer
	// startedInformers is used for tracking which informers have been started.
	// This allows Start() to be called multiple times safely.
	startedInformers map[reflect.Type]bool
	started          bool
}

func NewResourceClientSharedInformerFactory() *ResourceClientSharedInformerFactory {
	return &ResourceClientSharedInformerFactory{
		defaultResync:    12 * time.Hour,
		informers:        make(map[reflect.Type]cache.SharedIndexInformer),
		startedInformers: make(map[reflect.Type]bool),
	}
}
func (f *ResourceClientSharedInformerFactory) InitErr() error {
	return f.initError
}

func (f *ResourceClientSharedInformerFactory) Register(rc *ResourceClient) {
	ctx := context.TODO()
	if ctxWithTags, err := tag.New(ctx, tag.Insert(KeyKind, rc.resourceName)); err == nil {
		ctx = ctxWithTags
	}

	list := rc.kube.ResourcesV1().Resources(metav1.NamespaceAll).List
	watch := rc.kube.ResourcesV1().Resources(metav1.NamespaceAll).Watch
	sharedInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				//if tweakListOptions != nil {
				//	tweakListOptions(&options)
				//}

				if ctxWithTags, err := tag.New(ctx, tag.Insert(KeyOpKind, "list")); err == nil {
					ctx = ctxWithTags
				}
				stats.Record(ctx, MLists.M(1), MInFlight.M(1))
				defer stats.Record(ctx, MInFlight.M(-1))
				return list(options)
			},
			WatchFunc: func(options metav1.ListOptions) (kubewatch.Interface, error) {
				// if tweakListOptions != nil {
				// 	tweakListOptions(&options)
				// }

				if ctxWithTags, err := tag.New(ctx, tag.Insert(KeyOpKind, "watch")); err == nil {
					ctx = ctxWithTags
				}

				stats.Record(ctx, MWatches.M(1), MInFlight.M(1))
				defer stats.Record(ctx, MInFlight.M(-1))
				return watch(options)
			},
		},
		&v1.Resource{},  // TODO(yuval-k): can we make this rc.crd.Type ?
		f.defaultResync, // TODO(yuval-k): make this configurable!
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	f.RegisterInformer(ctx, rc.crd.Type, sharedInformer)
}

func (f *ResourceClientSharedInformerFactory) RegisterInformer(ctx context.Context, obj runtime.Object, newInformer cache.SharedIndexInformer) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.started {
		contextutils.LoggerFrom(ctx).DPanic("can't register informer after factory has started. this may change in the future")
	}

	informerType := reflect.TypeOf(obj)
	informer, exists := f.informers[informerType]
	if exists {
		return informer
	}
	informer = newInformer
	f.informers[informerType] = informer

	return informer
}
func (f *ResourceClientSharedInformerFactory) GetInformer(obj runtime.Object) cache.SharedIndexInformer {
	f.lock.Lock()
	defer f.lock.Unlock()

	informerType := reflect.TypeOf(obj)
	return f.informers[informerType]
}

func (f *ResourceClientSharedInformerFactory) Start(ctx context.Context, kubeClient kubernetes.Interface, updatecallback func()) {
	stop := ctx.Done()
	var sharedInformers []cache.SharedInformer
	func() {

		f.lock.Lock()
		defer f.lock.Unlock()

		for informerType, informer := range f.informers {
			if !f.startedInformers[informerType] {
				go informer.Run(stop)
				f.startedInformers[informerType] = true
				sharedInformers = append(sharedInformers, informer)
			}
		}

		f.started = true
	}()

	kubeController := kubecontroller.NewController("solo-resource-controller", kubeClient,
		kubecontroller.NewLockingSyncHandler(updatecallback),
		sharedInformers...)
	go kubeController.Run(2, stop)

	for _, informer := range sharedInformers {
		ok := cache.WaitForCacheSync(stop, informer.HasSynced)
		if !ok {
			// if initError is non-nil, the kube resource client will panic
			f.initError = errors.Errorf("waiting for initial cache sync failed")
			return
		}
	}
}

func (f *ResourceClientSharedInformerFactory) GetLister(obj runtime.Object) (ResourceLister, error) {
	informer := f.GetInformer(obj)
	if informer == nil {
		return nil, errors.Errorf("no lister has been registered for ObjectKind %v", obj.GetObjectKind())
	}
	return &resourceLister{indexer: informer.GetIndexer()}, nil

}

type resourceLister struct {
	indexer cache.Indexer
}

func (s *resourceLister) List(selector labels.Selector) (ret []*v1.Resource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Resource))
	})
	return ret, err

}
