package vault

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	keyMetadata = "_metadata"
	//keyResourceName      = "_resource_name"
	//keyResourceNamespace = "_resource_namespace"
	//keyResourceVersion   = "_resource_version"
)

func toVaultSecret(resource thirdparty.ThirdPartyResource) map[string]interface{} {
	var data thirdparty.Data
	switch res := resource.(type) {
	case *thirdparty.Artifact:
		data = res.Data
	case *thirdparty.Secret:
		data = res.Data
	default:
		panic("invalid resource type " + reflect.TypeOf(resource).String())
	}
	values := make(map[string]interface{})
	for k, v := range data.Values {
		values[k] = v
	}
	metaBytes, err := json.Marshal(data.Metadata)
	if err != nil {
		panic("unexpected marshal err: " + err.Error())
	}
	values[keyMetadata] = string(metaBytes)
	return values
}

func fromVaultSecret(secret *api.Secret) (thirdparty.Data, error) {
	values := make(map[string]string)
	for k, v := range secret.Data {
		if k == keyMetadata {
			continue
		}
		values[k] = v.(string)
	}
	metaValue, ok := secret.Data[keyMetadata]
	if !ok {
		return thirdparty.Data{}, errors.Errorf("secret missing required key %v", keyMetadata)
	}
	metaString, ok := metaValue.(string)
	if !ok {
		return thirdparty.Data{}, errors.Errorf("key %v present but value was not string", keyMetadata)
	}
	var meta core.Metadata
	err := json.Unmarshal([]byte(metaString), &meta)
	if err != nil {
		return thirdparty.Data{}, errors.Wrapf(err, "key %v present but value was not string", keyMetadata)
	}
	return thirdparty.Data{
		Metadata: meta,
		Values:   values,
	}, nil
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr+1)
}

type thirdPartyResourceClient struct {
	vault                  *api.Client
	root                   string
	thirdPartyResourceType thirdparty.ThirdPartyResource
}

func NewThirdPartyResourceClient(client *api.Client, rootKey string, thirdPartyResourceType thirdparty.ThirdPartyResource) *thirdPartyResourceClient {
	return &thirdPartyResourceClient{
		vault: client,
		root:  rootKey,
		thirdPartyResourceType: thirdPartyResourceType,
	}
}
func (rc *thirdPartyResourceClient) kind() string {
	return reflect.TypeOf(rc.thirdPartyResourceType).String()
}

func (rc *thirdPartyResourceClient) newThirdPartyResource(data thirdparty.Data) thirdparty.ThirdPartyResource {
	switch rc.thirdPartyResourceType.(type) {
	case *thirdparty.Artifact:
		return &thirdparty.Artifact{
			Data: data,
		}
	case *thirdparty.Secret:
		return &thirdparty.Secret{
			Data: data,
		}
	}
	panic("unknown or unsupported third party resource type " + rc.kind())
}

func (rc *thirdPartyResourceClient) Read(namespace, name string, opts clients.ReadOpts) (thirdparty.ThirdPartyResource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	key := rc.thirdPartyResourceKey(namespace, name)

	secret, err := rc.vault.Logical().Read(key)
	if err != nil {
		return nil, errors.Wrapf(err, "performing vault KV get")
	}
	if secret == nil {
		return nil, errors.NewNotExistErr(namespace, name)
	}

	data, err := fromVaultSecret(secret)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing vault secret")
	}

	thirdPartyResource := rc.newThirdPartyResource(data)
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
	meta.ResourceVersion = newOrIncrementResourceVer(meta.ResourceVersion)
	clone.SetMetadata(meta)

	if _, err := rc.vault.Logical().Write(key, toVaultSecret(clone)); err != nil {
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
	_, err := rc.vault.Logical().Delete(key)
	if err != nil {
		return errors.Wrapf(err, "deleting thirdPartyResource %v", name)
	}
	return nil
}

func (rc *thirdPartyResourceClient) List(namespace string, opts clients.ListOpts) ([]thirdparty.ThirdPartyResource, error) {
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

	var thirdPartyResourceList []thirdparty.ThirdPartyResource
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
		data, err := fromVaultSecret(secret)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing vault secret")
		}

		thirdPartyResource := rc.newThirdPartyResource(data)
		if labels.SelectorFromSet(opts.Selector).Matches(labels.Set(thirdPartyResource.GetMetadata().Labels)) {
			thirdPartyResourceList = append(thirdPartyResourceList, thirdPartyResource)
		}
	}
	return thirdPartyResourceList, nil

	sort.SliceStable(thirdPartyResourceList, func(i, j int) bool {
		return thirdPartyResourceList[i].GetMetadata().Name < thirdPartyResourceList[j].GetMetadata().Name
	})

	return thirdPartyResourceList, nil
}

func (rc *thirdPartyResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []thirdparty.ThirdPartyResource, <-chan error, error) {
	opts = opts.WithDefaults()
	namespace = clients.DefaultNamespaceIfEmpty(namespace)
	thirdPartyResourcesChan := make(chan []thirdparty.ThirdPartyResource)
	errs := make(chan error)
	var cached []thirdparty.ThirdPartyResource
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
		thirdPartyResourcesChan <- list
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
					thirdPartyResourcesChan <- list
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()

	return thirdPartyResourcesChan, errs, nil
}

func (rc *thirdPartyResourceClient) thirdPartyResourceKey(namespace, name string) string {
	return filepath.Join(rc.root, namespace, name)
}
