package vault

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/vault/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	keyMetadata = "_metadata"
	//keyResourceName      = "_resource_name"
	//keyResourceNamespace = "_resource_namespace"
	//keyResourceVersion   = "_resource_version"
)

func toVaultSecret(resource resources.DataResource) map[string]interface{} {
	values := make(map[string]interface{})
	for k, v := range resource.GetData() {
		values[k] = v
	}
	metaBytes, err := json.Marshal(resource.GetMetadata())
	if err != nil {
		panic("unexpected marshal err: " + err.Error())
	}
	values[keyMetadata] = string(metaBytes)
	return values
}

func fromVaultSecret(secret *api.Secret, into resources.DataResource) error {
	values := make(map[string]string)
	for k, v := range secret.Data {
		if k == keyMetadata {
			continue
		}
		values[k] = v.(string)
	}
	metaValue, ok := secret.Data[keyMetadata]
	if !ok {
		return errors.Errorf("secret missing required key %v", keyMetadata)
	}
	metaString, ok := metaValue.(string)
	if !ok {
		return errors.Errorf("key %v present but value was not string", keyMetadata)
	}
	var meta core.Metadata
	err := json.Unmarshal([]byte(metaString), &meta)
	if err != nil {
		return errors.Wrapf(err, "key %v present but value was not string", keyMetadata)
	}
	into.SetData(values)
	into.SetMetadata(meta)
	return nil
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr+1)
}

type ResourceClient struct {
	vault        *api.Client
	root         string
	resourceType resources.DataResource
}

func NewResourceClient(client *api.Client, rootKey string, resourceType resources.DataResource) *ResourceClient {
	return &ResourceClient{
		vault:        client,
		root:         rootKey,
		resourceType: resourceType,
	}
}

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
	key := rc.resourceKey(namespace, name)

	secret, err := rc.vault.Logical().Read(key)
	if err != nil {
		return nil, errors.Wrapf(err, "performing vault KV get")
	}
	if secret == nil {
		return nil, errors.NewNotExistErr(namespace, name)
	}

	resource := rc.NewResource()
	if err = fromVaultSecret(secret, resource.(resources.DataResource)); err != nil {
		return nil, errors.Wrapf(err, "parsing vault secret")
	}
	return resource, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error) {
	opts = opts.WithDefaults()
	if err := resources.ValidateName(resource.GetMetadata().Name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	meta := resource.GetMetadata()
	meta.Namespace = clients.DefaultNamespaceIfEmpty(meta.Namespace)
	key := rc.resourceKey(meta.Namespace, meta.Name)

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
	clone := proto.Clone(resource).(resources.DataResource)
	meta.ResourceVersion = newOrIncrementResourceVer(meta.ResourceVersion)
	clone.SetMetadata(meta)

	if _, err := rc.vault.Logical().Write(key, toVaultSecret(clone)); err != nil {
		return nil, errors.Wrapf(err, "writing to KV")
	}
	// return a read object to update the modify index
	return rc.Read(meta.Namespace, meta.Name, clients.ReadOpts{Ctx: opts.Ctx})
}

func (rc *ResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	key := rc.resourceKey(namespace, name)
	if !opts.IgnoreNotExist {
		if _, err := rc.Read(namespace, name, clients.ReadOpts{Ctx: opts.Ctx}); err != nil {
			return errors.NewNotExistErr(namespace, name, err)
		}
	}
	_, err := rc.vault.Logical().Delete(key)
	if err != nil {
		return errors.Wrapf(err, "deleting resource %v", name)
	}
	return nil
}

func (rc *ResourceClient) List(namespace string, opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)

	namespacePrefix := filepath.Join(rc.root, namespace)
	secrets, err := rc.vault.Logical().List(namespacePrefix)
	if err != nil {
		return nil, errors.Wrapf(err, "reading namespace root")
	}
	val, ok := secrets.Data["keys"]
	if !ok {
		return nil, errors.Errorf("vault secret list at root %s did not contain key \"keys\"", namespacePrefix)
	}
	keys, ok := val.([]interface{})
	if !ok {
		return nil, errors.Errorf("expected secret list of type []interface{} but got %v", reflect.TypeOf(val))
	}

	var resourceList []resources.Resource
	for _, keyAsInterface := range keys {
		key, ok := keyAsInterface.(string)
		if !ok {
			return nil, errors.Errorf("expected key of type string but got %v", reflect.TypeOf(keyAsInterface))
		}
		secret, err := rc.vault.Logical().Read(namespacePrefix + "/" + key)
		if err != nil {
			return nil, errors.Wrapf(err, "getting secret %s", key)
		}
		if secret == nil {
			return nil, errors.Errorf("unexpected nil err on %v", key)
		}

		resource := rc.NewResource()
		if err = fromVaultSecret(secret, resource.(resources.DataResource)); err != nil {
			return nil, errors.Wrapf(err, "parsing vault secret")
		}
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(resource.GetMetadata().Labels)) {
			resourceList = append(resourceList, resource)
		}
	}
	return resourceList, nil

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
	var cached []resources.Resource
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
		cached = list
		resourcesChan <- list
	}()
	go func() {
		for {
			select {
			case <-time.After(opts.RefreshRate):
				list, err := rc.List(namespace, clients.ListOpts{
					Ctx: opts.Ctx,
				})
				if err != nil {
					errs <- err
				}
				if list != nil && !reflect.DeepEqual(list, cached) {
					cached = list
					resourcesChan <- list
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()

	return resourcesChan, errs, nil
}

func (rc *ResourceClient) resourceKey(namespace, name string) string {
	return filepath.Join(rc.root, namespace, name)
}
