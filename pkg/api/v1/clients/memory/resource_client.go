package memory

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

const separator = "~;~"

type InMemoryResourceCache interface {
	Get(key string) (resources.Resource, bool)
	Delete(key string)
	Set(key string, resource resources.Resource)
	List(prefix string) resources.ResourceList
	Subscribe(subscription chan struct{})
	Unsubscribe(subscription chan struct{})
}

type inMemoryResourceCache struct {
	store       map[string]resources.Resource
	lock        sync.RWMutex
	subscribers []chan struct{}
}

func (c *inMemoryResourceCache) signalUpdate() {
	for _, subscription := range c.subscribers {
		go func() {
			subscription <- struct{}{}
		}()
	}
}

func (c *inMemoryResourceCache) Get(key string) (resources.Resource, bool) {
	c.lock.RLock()
	resource, ok := c.store[key]
	c.lock.RUnlock()
	return resource, ok
}

func (c *inMemoryResourceCache) Delete(key string) {
	c.lock.Lock()
	delete(c.store, key)
	c.signalUpdate()
	c.lock.Unlock()
}

func (c *inMemoryResourceCache) Set(key string, resource resources.Resource) {
	c.lock.Lock()
	c.store[key] = resource
	c.signalUpdate()
	c.lock.Unlock()
}

func (c *inMemoryResourceCache) List(prefix string) resources.ResourceList {
	var ress resources.ResourceList
	c.lock.RLock()
	defer c.lock.RUnlock()
	for key, resource := range c.store {
		if !strings.HasPrefix(key, prefix) {
			continue
		}
		ress = append(ress, resource)
	}
	return ress
}

func (c *inMemoryResourceCache) Subscribe(subscription chan struct{}) {
	c.lock.Lock()
	c.subscribers = append(c.subscribers, subscription)
	c.lock.Unlock()
}

func (c *inMemoryResourceCache) Unsubscribe(subscription chan struct{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	for i, sub := range c.subscribers {
		if sub == subscription {
			c.subscribers = append(c.subscribers[:i], c.subscribers[i+1:]...)
			return
		}
	}
}

func NewInMemoryResourceCache() InMemoryResourceCache {
	return &inMemoryResourceCache{
		store: make(map[string]resources.Resource),
	}
}

type ResourceClient struct {
	resourceType resources.Resource
	cache        InMemoryResourceCache
}

func NewResourceClient(cache InMemoryResourceCache, resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		cache:        cache,
		resourceType: resourceType,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Kind() string {
	return resources.Kind(rc.resourceType)
}

func (rc *ResourceClient) NewResource() resources.Resource {
	return resources.Clone(rc.resourceType)
}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	resource, ok := rc.cache.Get(rc.key(namespace, name))
	if !ok {
		return nil, errors.NewNotExistErr(namespace, name)
	}
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)

	key := rc.key(meta.Namespace, meta.Name)

	original, err := rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{})
	if original != nil && err == nil {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if meta.ResourceVersion != original.GetMetadata().ResourceVersion {
			return nil, errors.Errorf("resource version error. must update new resource version to match current")
		}
	}

	// mutate and return clone
	clone := proto.Clone(resource).(resources.Resource)
	// initialize or increment resource version
	meta.ResourceVersion = newOrIncrementResourceVer(meta.ResourceVersion)
	clone.SetMetadata(meta)

	rc.cache.Set(key, clone)

	return clone, nil
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	key := rc.key(namespace, name)
	_, ok := rc.cache.Get(key)
	if !ok {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	rc.cache.Delete(key)
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) (resources.ResourceList, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	cachedResources := rc.cache.List(namespace + separator)
	var resourceList resources.ResourceList
	for _, resource := range cachedResources {
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(resource.GetMetadata().Labels)) {
			resourceList = append(resourceList, resource)
		}
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan resources.ResourceList, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	resourcesChan := make(chan resources.ResourceList)
	errs := make(chan error)
	updateResourceList := func() {
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list.FilterByKind(rc.Kind())
	}
	go updateResourceList()
	subscription := make(chan struct{})
	rc.cache.Subscribe(subscription)
	go func() {
		for {
			select {
			case <-time.After(opts.RefreshRate):
				updateResourceList()
			case <-subscription:
				updateResourceList()
			case <-opts.Ctx.Done():
				close(subscription)
				close(resourcesChan)
				close(errs)
				rc.cache.Unsubscribe(subscription)
				return
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) key(namespace, name string) string {
	return namespace + separator + name
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr+1)
}
