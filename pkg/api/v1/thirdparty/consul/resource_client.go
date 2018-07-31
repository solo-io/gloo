package consul

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"

	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
	"github.com/solo-io/solo-kit/pkg/errors"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/labels"
)

type thirdPartyResourceClient struct {
	consul                 *api.Client
	root                   string
	thirdPartyResourceType thirdparty.ThirdPartyResource
}

func NewThirdPartyResourceClient(client *api.Client, rootKey string, thirdPartyResourceType thirdparty.ThirdPartyResource) *thirdPartyResourceClient {
	return &thirdPartyResourceClient{
		consul: client,
		root:   rootKey,
		thirdPartyResourceType: thirdPartyResourceType,
	}
}
func (rc *thirdPartyResourceClient) kind() string {
	return reflect.TypeOf(rc.thirdPartyResourceType).String()
}

func (rc *thirdPartyResourceClient) newThirdPartyResource() thirdparty.ThirdPartyResource {
	switch rc.thirdPartyResourceType.(type) {
	case *thirdparty.Artifact:
		return &thirdparty.Artifact{}
	case *thirdparty.Secret:
		return &thirdparty.Secret{}
	}
	panic("unknown or unsupported third party resource type " + rc.kind())
}

func (rc *thirdPartyResourceClient) Read(namespace, name string, opts clients.ReadOpts) (thirdparty.ThirdPartyResource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	key := rc.thirdPartyResourceKey(namespace, name)

	kvPair, _, err := rc.consul.KV().Get(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "performing consul KV get")
	}
	if kvPair == nil {
		return nil, errors.NewNotExistErr(namespace, name)
	}
	thirdPartyResource := rc.newThirdPartyResource()
	if err := yaml.Unmarshal(kvPair.Value, thirdPartyResource); err != nil {
		return nil, errors.Wrapf(err, "reading KV into %v", rc.kind())
	}
	resources.UpdateMetadata(thirdPartyResource, func(meta *core.Metadata) {
		meta.ResourceVersion = fmt.Sprintf("%v", kvPair.ModifyIndex)
	})
	return thirdPartyResource, nil
}

func (rc *thirdPartyResourceClient) Write(thirdPartyResource thirdparty.ThirdPartyResource, opts clients.WriteOpts) (thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	if err := resources.ValidateName(thirdPartyResource.GetMetadata().Name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := thirdPartyResource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)
	key := rc.thirdPartyResourceKey(meta.Namespace, meta.Name)

	original, err := rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{})
	if original != nil && err == nil {
		if !opts.OverwriteExisting {
			return nil, errors.NewExistErr(meta)
		}
		if meta.ResourceVersion != original.GetMetadata().ResourceVersion {
			return nil, errors.Errorf("thirdPartyResource version error. must update new thirdPartyResource version to match current")
		}
	}

	// mutate and return clone
	clone := thirdPartyResource.DeepCopy()
	clone.SetMetadata(meta)

	data, err := yaml.Marshal(clone)
	if err != nil {
		panic(errors.Wrapf(err, "internal err: failed to marshal thirdPartyResource"))
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
	return rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *thirdPartyResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	key := rc.thirdPartyResourceKey(namespace, name)
	if !opts.IgnoreNotExist {
		if _, err := rc.Read(namespace, name, clients.ReadOpts{Ctx: opts.Ctx}); err != nil {
			return errors.NewNotExistErr(namespace, name, err)
		}
	}
	_, err := rc.consul.KV().Delete(key, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting thirdPartyResource %v", name)
	}
	return nil
}

func (rc *thirdPartyResourceClient) List(namespace string, opts clients.ListOpts) ([]thirdparty.ThirdPartyResource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	namespacePrefix := filepath.Join(rc.root, namespace)
	kvPairs, _, err := rc.consul.KV().List(namespacePrefix, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "reading namespace root")
	}

	var thirdPartyResourceList []thirdparty.ThirdPartyResource
	for _, kvPair := range kvPairs {
		thirdPartyResource := rc.newThirdPartyResource()
		if err := yaml.Unmarshal(kvPair.Value, thirdPartyResource); err != nil {
			return nil, errors.Wrapf(err, "reading KV into %v", rc.kind())
		}
		resources.UpdateMetadata(thirdPartyResource, func(meta *core.Metadata) {
			meta.ResourceVersion = fmt.Sprintf("%v", kvPair.ModifyIndex)
		})
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(thirdPartyResource.GetMetadata().Labels)) {
			thirdPartyResourceList = append(thirdPartyResourceList, thirdPartyResource)
		}
	}

	sort.SliceStable(thirdPartyResourceList, func(i, j int) bool {
		return thirdPartyResourceList[i].GetMetadata().Name < thirdPartyResourceList[j].GetMetadata().Name
	})

	return thirdPartyResourceList, nil
}

func (rc *thirdPartyResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []thirdparty.ThirdPartyResource, <-chan error, error) {
	opts = opts.WithDefaults()
	var lastIndex uint64
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	namespacePrefix := filepath.Join(rc.root, namespace)
	thirdPartyResourcesChan := make(chan []thirdparty.ThirdPartyResource)
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
		thirdPartyResourcesChan <- list
	}()
	updatedThirdPartyResourceList := func() ([]thirdparty.ThirdPartyResource, error) {
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
		var thirdPartyResourceList []thirdparty.ThirdPartyResource
		for _, kvPair := range kvPairs {
			thirdPartyResource := rc.newThirdPartyResource()
			if err := yaml.Unmarshal(kvPair.Value, thirdPartyResource); err != nil {
				return nil, errors.Wrapf(err, "reading KV into %v", rc.kind())
			}
			resources.UpdateMetadata(thirdPartyResource, func(meta *core.Metadata) {
				meta.ResourceVersion = fmt.Sprintf("%v", kvPair.ModifyIndex)
			})
			if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(thirdPartyResource.GetMetadata().Labels)) {
				thirdPartyResourceList = append(thirdPartyResourceList, thirdPartyResource)
			}
		}

		sort.SliceStable(thirdPartyResourceList, func(i, j int) bool {
			return thirdPartyResourceList[i].GetMetadata().Name < thirdPartyResourceList[j].GetMetadata().Name
		})

		// update index
		lastIndex = meta.LastIndex
		return thirdPartyResourceList, nil
	}

	go func() {
		for {
			list, err := updatedThirdPartyResourceList()
			if err != nil {
				errs <- err
			}
			if list != nil {
				thirdPartyResourcesChan <- list
			}
		}
	}()

	return thirdPartyResourcesChan, errs, nil
}

func (rc *thirdPartyResourceClient) thirdPartyResourceKey(namespace, name string) string {
	return filepath.Join(rc.root, namespace, name)
}
