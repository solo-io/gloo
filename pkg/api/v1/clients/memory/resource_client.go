package memory

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

const separator = "~;~"

type ResourceClient struct {
	lock         sync.RWMutex
	cache        map[string]resources.Resource
	updates      chan struct{}
	resourceType resources.Resource
}

func NewResourceClient(resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		cache:        make(map[string]resources.Resource),
		updates:      make(chan struct{}),
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
	rc.lock.RLock()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	resource, ok := rc.cache[rc.key(namespace, name)]
	rc.lock.RUnlock()
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

	rc.lock.Lock()
	rc.cache[key] = clone
	rc.signalUpdate()
	rc.lock.Unlock()

	return clone, nil
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	key := rc.key(namespace, name)
	rc.lock.RLock()
	_, ok := rc.cache[key]
	rc.lock.RUnlock()
	if !ok {
		if !opts.IgnoreNotExist {
			return errors.NewNotExistErr(namespace, name)
		}
		return nil
	}

	rc.lock.Lock()
	delete(rc.cache, key)
	rc.signalUpdate()
	rc.lock.Unlock()
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	var resourceList []resources.Resource
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	for key, resource := range rc.cache {
		if !strings.HasPrefix(key, namespace+separator) {
			continue
		}
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(resource.GetMetadata().Labels)) {
			resourceList = append(resourceList, resource)
		}
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []resources.Resource, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	resourcesChan := make(chan []resources.Resource)
	errs := make(chan error)
	go func() {
		// watch should open up with an initial read
		list, err := rc.List(namespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}()
	go func() {
		for {
			select {
			case <-rc.updates:
				list, err := rc.List(namespace, clients.ListOpts{
					Ctx:      opts.Ctx,
					Selector: opts.Selector,
				})
				if err != nil {
					errs <- err
					continue
				}
				resourcesChan <- list
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) key(namespace, name string) string {
	return namespace + separator + name
}

func (rc *ResourceClient) signalUpdate() {
	go func() {
		rc.updates <- struct{}{}
	}()
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr+1)
}
