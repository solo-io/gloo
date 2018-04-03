package consul

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"time"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/log"
)

// TODO: evaluate efficiency of LSing a whole dir on every op
// so far this is preferable to caring what files are named
type baseClient struct {
	rootPath      string
	consul        *api.Client
	syncFrequency time.Duration
}

type ConfigObjectType string

const (
	configObjectTypeUpstream    = "Upstream"
	configObjectTypeVirtualHost = "VirtualHost"
)

func (c *baseClient) Create(item v1.ConfigObject, t ConfigObjectType) (v1.ConfigObject, error) {
	p, err := toKVPair(c.rootPath, item)
	if err != nil {
		return nil, errors.Wrapf(err, "converting %s to kv pair", item.GetName())
	}

	// error if the key already exists
	existingP, _, err := c.consul.KV().Get(p.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query consul")
	}
	if existingP != nil {
		return nil, storage.NewAlreadyExistsErr(
			errors.Errorf("key found for configObject %s: %s", item.GetName(), p.Key))
	}
	if _, err := c.consul.KV().Put(p, nil); err != nil {
		return nil, errors.Wrapf(err, "writing kv pair %s", p.Key)
	}
	// set the resource version from the CreateIndex of the created kv pair
	p, _, err = c.consul.KV().Get(p.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "getting newly created kv pair %s", p.Key)
	}
	// set resourceversion on clone
	configObjectClone := proto.Clone(item).(v1.ConfigObject)
	setResourceVersion(configObjectClone, p, t)
	return configObjectClone, nil
}

func (c *baseClient) Update(item v1.ConfigObject, t ConfigObjectType) (v1.ConfigObject, error) {
	updatedP, err := toKVPair(c.rootPath, item)
	if err != nil {
		return nil, errors.Wrapf(err, "converting %s to kv pair", item.GetName())
	}

	// error if the key doesn't already exist
	exsitingP, _, err := c.consul.KV().Get(updatedP.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query consul")
	}
	if exsitingP == nil {
		return nil, errors.Errorf("key not found for configObject %s: %s", item.GetName(), updatedP.Key)
	}
	if success, _, err := c.consul.KV().CAS(updatedP, nil); err != nil {
		return nil, errors.Wrapf(err, "writing kv pair %s", updatedP.Key)
	} else if !success {
		return nil, errors.Errorf("resource version was invalid for configObject: %s", item.GetName())
	}

	// set the resource version from the CreateIndex of the created kv pair
	updatedP, _, err = c.consul.KV().Get(updatedP.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "getting newly created kv pair %s", updatedP.Key)
	}
	// set resourceversion on clone
	configObjectClone, ok := proto.Clone(item).(v1.ConfigObject)
	if !ok {
		return nil, errors.New("internal error: output of proto.Clone was not expected type")
	}
	setResourceVersion(configObjectClone, updatedP, t)
	return configObjectClone, nil
}

func (c *baseClient) Delete(name string) error {
	key := key(c.rootPath, name)

	_, err := c.consul.KV().Delete(key, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting %s", name)
	}
	return nil
}

func (c *baseClient) Get(name string, t ConfigObjectType) (v1.ConfigObject, error) {
	key := key(c.rootPath, name)
	p, _, err := c.consul.KV().Get(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "getting pair for for key %v", key)
	}
	if p == nil {
		return nil, errors.Errorf("keypair %s not found for configObject %s", key, name)
	}
	obj, err := configObjectFromKVPair(p, t)
	if err != nil {
		return nil, errors.Wrap(err, "converting consul kv-pair to configObject")
	}
	return obj, nil
}

func (c *baseClient) List(t ConfigObjectType) ([]v1.ConfigObject, error) {
	pairs, _, err := c.consul.KV().List(c.rootPath, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "listing key-value pairs for root %s", c.rootPath)
	}
	var configObjects []v1.ConfigObject
	for _, p := range pairs {
		obj, err := configObjectFromKVPair(p, t)
		if err != nil {
			return nil, errors.Wrapf(err, "converting %s to configObject", p.Key)
		}
		configObjects = append(configObjects, obj)
	}
	return configObjects, nil
}

// TODO: be clear that watch for consul only calls update
func (c *baseClient) Watch(t ConfigObjectType, handlers ...storage.ConfigObjectEventHandler) (*storage.Watcher, error) {
	var lastIndex uint64
	sync := func() error {
		pairs, meta, err := c.consul.KV().List(c.rootPath, &api.QueryOptions{RequireConsistent: true})
		if err != nil {
			return errors.Wrap(err, "getting kv-pairs list")
		}
		// no change since last poll
		if lastIndex == meta.LastIndex {
			return nil
		}
		var (
			virtualHosts []*v1.VirtualHost
			upstreams    []*v1.Upstream
		)
		for _, p := range pairs {
			obj, err := configObjectFromKVPair(p, t)
			if err != nil {
				return errors.Wrapf(err, "converting %s to configObject", p.Key)
			}

			switch t {
			case configObjectTypeUpstream:
				upstreams = append(upstreams, obj.(*v1.Upstream))
			case configObjectTypeVirtualHost:
				virtualHosts = append(virtualHosts, obj.(*v1.VirtualHost))
			}
		}
		lastIndex = meta.LastIndex
		for _, h := range handlers {
			switch t {
			case configObjectTypeUpstream:
				h.UpstreamEventHandler.OnUpdate(upstreams, nil)
			case configObjectTypeVirtualHost:
				h.VirtualHostEventHandler.OnUpdate(virtualHosts, nil)
			}
		}
		return nil
	}
	return storage.NewWatcher(func(stop <-chan struct{}, errs chan error) {
		for {
			select {
			case <-time.After(c.syncFrequency):
				if err := sync(); err != nil {
					log.Warnf("error syncing with consul kv-pairs: %v", err)
				}
			case err := <-errs:
				log.Warnf("failed to start watcher to: %v", err)
				return
			case <-stop:
				return
			}
		}
	}), nil
}
