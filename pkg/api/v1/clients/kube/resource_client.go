package kube

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	clientFactory = NewResourceClientSharedInformerFactory()

	cacheUpdatedWatchers      []chan struct{}
	cacheUpdatedWatchersMutex sync.Mutex

	factoryStarter sync.Once
)

func addWatch() <-chan struct{} {
	cacheUpdatedWatchersMutex.Lock()
	defer cacheUpdatedWatchersMutex.Unlock()
	c := make(chan struct{}, 1)
	cacheUpdatedWatchers = append(cacheUpdatedWatchers, c)
	return c
}
func removeWatch(c <-chan struct{}) {
	cacheUpdatedWatchersMutex.Lock()
	defer cacheUpdatedWatchersMutex.Unlock()
	for i, cacheUpdated := range cacheUpdatedWatchers {
		if cacheUpdated == c {
			cacheUpdatedWatchers = append(cacheUpdatedWatchers[:i], cacheUpdatedWatchers[i+1:]...)
			return
		}
	}
}

func startFactory(ctx context.Context, client kubernetes.Interface) {
	factoryStarter.Do(func() {
		clientFactory.Start(ctx, client, updatedOccured)
	})
}

// TODO(yuval-k): See if we can get more fine grained updates here, about which resources was udpated
func updatedOccured() {
	cacheUpdatedWatchersMutex.Lock()
	defer cacheUpdatedWatchersMutex.Unlock()
	for _, cacheUpdated := range cacheUpdatedWatchers {
		select {
		case cacheUpdated <- struct{}{}:
		default:
		}
	}
}

// lazy start in list & watch
// register informers in register

type ResourceClient struct {
	crd          crd.Crd
	apiexts      apiexts.Interface
	kube         versioned.Interface
	kubeClient   kubernetes.Interface
	ownerLabel   string
	resourceName string
	resourceType resources.InputResource
}

func NewResourceClient(crd crd.Crd, cfg *rest.Config, resourceType resources.InputResource) (*ResourceClient, error) {
	apiExts, err := apiexts.NewForConfig(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "creating api extensions client")
	}
	crdClient, err := versioned.NewForConfig(cfg, crd)
	if err != nil {
		return nil, errors.Wrapf(err, "creating crd client")
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	typeof := reflect.TypeOf(resourceType)
	resourceName := strings.Replace(typeof.String(), "*", "", -1)
	resourceName = strings.Replace(resourceName, ".", "", -1)

	return &ResourceClient{
		crd:          crd,
		apiexts:      apiExts,
		kube:         crdClient,
		kubeClient:   kubeClient,
		resourceName: resourceName,
		resourceType: resourceType,
	}, nil
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Kind() string {
	return resources.Kind(rc.resourceType)
}

func (rc *ResourceClient) NewResource() resources.Resource {
	return resources.Clone(rc.resourceType)
}

func (rc *ResourceClient) Register() error {
	/*
		shared informer factory for all namespaces; and then filter desired namespace in list and watch!
		zbam!

	*/
	clientFactory.Register(rc)
	return rc.crd.Register(rc.apiexts)
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	resourceCrd, err := rc.kube.ResourcesV1().Resources(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.NewNotExistErr(namespace, name, err)
		}
		return nil, errors.Wrapf(err, "reading resource from kubernetes")
	}
	resource := rc.NewResource()
	if resourceCrd.Spec != nil {
		if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
			return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
		}
	}
	resources.UpdateMetadata(resource, func(meta *core.Metadata) {
		meta.ResourceVersion = resourceCrd.ResourceVersion
	})
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)

	// mutate and return clone
	clone := proto.Clone(resource).(resources.InputResource)
	clone.SetMetadata(meta)
	resourceCrd := rc.crd.KubeResource(clone)

	if rc.exist(meta.Namespace, meta.Name) {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if _, updateerr := rc.kube.ResourcesV1().Resources(meta.Namespace).Update(resourceCrd); updateerr != nil {
			original, err := rc.kube.ResourcesV1().Resources(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			if err == nil {
				return nil, errors.Wrapf(updateerr, "updating kube resource %v:%v (want %v)", resourceCrd.Name, resourceCrd.ResourceVersion, original.ResourceVersion)
			}
			return nil, errors.Wrapf(updateerr, "updating kube resource %v", resourceCrd.Name)
		}
	} else {
		if _, err := rc.kube.ResourcesV1().Resources(meta.Namespace).Create(resourceCrd); err != nil {
			return nil, errors.Wrapf(err, "creating kube resource %v", resourceCrd.Name)
		}
	}

	// return a read object to update the resource version
	return rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	if !rc.exist(namespace, name) {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	if err := rc.kube.ResourcesV1().Resources(namespace).Delete(name, nil); err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {
	startFactory(context.TODO(), rc.kubeClient)

	lister := clientFactory.GetLister(rc.crd.Type)
	allResources, err := lister.List(labels.SelectorFromSet(opts.Selector))
	if err != nil {
		return nil, errors.Wrapf(err, "listing resources in %v", namespace)
	}
	var listedResources []*v1.Resource
	if namespace != "" {
		for _, r := range allResources {
			if r.ObjectMeta.Namespace == namespace {
				listedResources = append(listedResources, r)
			}
		}
	} else {
		listedResources = allResources
	}

	var resourceList resources.ResourceList
	for _, resourceCrd := range listedResources {
		resource := rc.NewResource()
		if resourceCrd.Spec != nil {
			if err := protoutils.UnmarshalMap(*resourceCrd.Spec, resource); err != nil {
				return nil, errors.Wrapf(err, "reading crd spec into %v", rc.resourceName)
			}
		}
		resources.UpdateMetadata(resource, func(meta *core.Metadata) {
			meta.Namespace = resourceCrd.Namespace
			meta.Name = resourceCrd.Name
			meta.ResourceVersion = resourceCrd.ResourceVersion
		})
		resourceList = append(resourceList, resource)
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan resources.ResourceList, <-chan error, error) {
	startFactory(context.TODO(), rc.kubeClient)

	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	resourcesChan := make(chan resources.ResourceList)
	errs := make(chan error)
	ctx := opts.Ctx

	updateResourceList := func() {
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}
	// watch should open up with an initial read
	cacheUpdated := addWatch()

	go func() {
		defer removeWatch(cacheUpdated)
		defer close(resourcesChan)
		defer close(errs)

		updateResourceList()

		for {
			select {
			case <-time.After(opts.RefreshRate): // TODO(yuval-k): can we remove this? is the factory takes care of that...
				updateResourceList()
			case <-cacheUpdated:
				updateResourceList()
			case <-ctx.Done():
				return
			}
		}

	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) exist(namespace, name string) bool {
	_, err := rc.kube.ResourcesV1().Resources(namespace).Get(name, metav1.GetOptions{}) // TODO(yuval-k): check error for real
	return err == nil

}
