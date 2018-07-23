package consul

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
)

type ResourceClient struct {
	consul       *api.Client
	root         string
	resourceType resources.Resource
}

func NewResourceClient(client *api.Client, rootKey string, resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		consul:       client,
		root:         rootKey,
		resourceType: resourceType,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(name string, opts clients.ReadOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	key := rc.resourceKey(opts.Namespace, name)

	kvPair, _, err := rc.consul.KV().Get(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "performing consul KV get")
	}
	if kvPair == nil {
		return nil, errors.NewNotExistErr(opts.Namespace, name)
	}
	resource := resources.Clone(rc.resourceType)
	if err := protoutils.UnmarshalBytes(kvPair.Value, resource); err != nil {
		return nil, errors.Wrapf(err, "reading KV into %v", reflect.TypeOf(rc.resourceType))
	}
	resources.UpdateMetadata(resource, func(meta *core.Metadata) {
		meta.ResourceVersion = fmt.Sprintf("%v", kvPair.ModifyIndex)
	})
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.Validate(resource); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	if meta.Namespace == "" {
		meta.Namespace = clients.DefaultNamespace
	}
	key := rc.resourceKey(meta.Namespace, meta.Name)

	if !opts.OverwriteExisting {
		kvPair, _, err := rc.consul.KV().Get(key, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "performing consul KV get")
		}
		if kvPair != nil {
			return nil, errors.NewExistErr(meta)
		}
	}

	// mutate and return clone
	clone := proto.Clone(resource).(resources.Resource)
	clone.SetMetadata(meta)

	data, err := protoutils.MarshalBytes(clone)
	if err != nil {
		panic(errors.Wrapf(err, "internal err: failed to marshal resource"))
	}
	var modifyIndex uint64
	if meta.GetResourceVersion() != "" {
		if i, err := strconv.Atoi(meta.GetResourceVersion()); err == nil {
			modifyIndex = uint64(i)
		}
	}
	kvPair := &api.KVPair{
		Key:         key,
		Value:       data,
		ModifyIndex: modifyIndex,
	}
	if _, err := rc.consul.KV().Put(kvPair, nil); err != nil {
		return nil, errors.Wrapf(err, "writing to KV")
	}
	// return a read object to update the modify index
	return rc.Read(meta.Name, clients.ReadOpts{Ctx: opts.Ctx, Namespace: meta.Namespace})
}

func (rc *ResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	key := rc.resourceKey(opts.Namespace, name)
	if !opts.IgnoreNotExist {
		if _, err := rc.Read(name, clients.ReadOpts{Namespace: opts.Namespace, Ctx: opts.Ctx}); err != nil {
			return errors.NewNotExistErr(opts.Namespace, name, err)
		}
	}
	_, err := rc.consul.KV().Delete(key, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()

	namespacePrefix := filepath.Join(rc.root, opts.Namespace)
	kvPairs, _, err := rc.consul.KV().List(namespacePrefix, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "reading namespace root")
	}

	var resourceList []resources.Resource
	for _, kvPair := range kvPairs {
		resource := resources.Clone(rc.resourceType)
		if err := protoutils.UnmarshalBytes(kvPair.Value, resource); err != nil {
			return nil, errors.Wrapf(err, "reading KV into %v", reflect.TypeOf(rc.resourceType))
		}
		resources.UpdateMetadata(resource, func(meta *core.Metadata) {
			meta.ResourceVersion = fmt.Sprintf("%v", kvPair.ModifyIndex)
		})
		resourceList = append(resourceList, resource)
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(opts clients.WatchOpts) (<-chan []resources.Resource, <-chan error, error) {
	opts = opts.WithDefaults()
	var lastIndex uint64
	namespacePrefix := filepath.Join(rc.root, opts.Namespace)
	resourcesChan := make(chan []resources.Resource)
	errs := make(chan error)
	go func() {
		// watch should open up with an initial read
		list, err := rc.List(clients.ListOpts{
			Ctx:       opts.Ctx,
			Selector:  opts.Selector,
			Namespace: opts.Namespace,
		})
		if err != nil {
			errs <- err
			return
		}
		resourcesChan <- list
	}()
	updatedResourceList := func() ([]resources.Resource, error) {
		kvPairs, meta, err := rc.consul.KV().List(namespacePrefix,
			&api.QueryOptions{
				RequireConsistent: true,
				WaitIndex:         lastIndex,
				WaitTime:          opts.RefreshRate,
			})
		if err != nil {
			return nil, errors.Wrapf(err, "getting kv-pairs list")
		}
		// no change since last poll
		if lastIndex == meta.LastIndex {
			return nil, nil
		}
		var resourceList []resources.Resource
		for _, kvPair := range kvPairs {
			resource := resources.Clone(rc.resourceType)
			if err := protoutils.UnmarshalBytes(kvPair.Value, resource); err != nil {
				return nil, errors.Wrapf(err, "reading KV into %v", reflect.TypeOf(rc.resourceType))
			}
			resources.UpdateMetadata(resource, func(meta *core.Metadata) {
				meta.ResourceVersion = fmt.Sprintf("%v", kvPair.ModifyIndex)
			})
			resourceList = append(resourceList, resource)
		}

		sort.SliceStable(resourceList, func(i, j int) bool {
			return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
		})

		// update index
		lastIndex = meta.LastIndex
		return resourceList, nil
	}

	go func() {
		for {
			list, err := updatedResourceList()
			if err != nil {
				errs <- err
			}
			if list != nil {
				resourcesChan <- list
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) resourceKey(namespace, name string) string {
	return filepath.Join(rc.root, namespace, name)
}
