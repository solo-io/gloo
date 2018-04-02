package consul

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/runtime"

	"time"

	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo/pkg/log"
)

// TODO: evaluate efficiency of LSing a whole dir on every op
// so far this is preferable to caring what files are named
type upstreamsClient struct {
	rootPath      string
	consul        *api.Client
	syncFrequency time.Duration
}

func (c *upstreamsClient) Create(item *v1.Upstream) (*v1.Upstream, error) {
	p, err := toKVPair(c.rootPath, item)
	if err != nil {
		return nil, errors.Wrapf(err, "converting %s to kv pair", item.Name)
	}

	// error if the key already exists
	existingP, _, err := c.consul.KV().Get(p.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query consul")
	}
	if existingP != nil {
		return nil, storage.NewAlreadyExistsErr(
			errors.Errorf("key found for upstream %s: %s", item.Name, p.Key))
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
	upstreamClone := proto.Clone(item).(*v1.Upstream)
	setResourceVersion(upstreamClone, p)
	return upstreamClone, nil
}

func (c *upstreamsClient) Update(item *v1.Upstream) (*v1.Upstream, error) {
	updatedP, err := toKVPair(c.rootPath, item)
	if err != nil {
		return nil, errors.Wrapf(err, "converting %s to kv pair", item.Name)
	}

	// error if the key doesn't already exist
	exsitingP, _, err := c.consul.KV().Get(updatedP.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query consul")
	}
	if exsitingP == nil {
		return nil, errors.Errorf("key not found for upstream %s: %s", item.Name, updatedP.Key)
	}
	if success, _, err := c.consul.KV().CAS(updatedP, nil); err != nil {
		return nil, errors.Wrapf(err, "writing kv pair %s", updatedP.Key)
	} else if !success {
		return nil, errors.Errorf("resource version was invalid for upstream: %s", item.Name)
	}

	// set the resource version from the CreateIndex of the created kv pair
	updatedP, _, err = c.consul.KV().Get(updatedP.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "getting newly created kv pair %s", updatedP.Key)
	}
	// set resourceversion on clone
	upstreamClone, ok := proto.Clone(item).(*v1.Upstream)
	if !ok {
		return nil, errors.New("internal error: output of proto.Clone was not expected type")
	}
	setResourceVersion(upstreamClone, updatedP)
	return upstreamClone, nil
}

func (c *upstreamsClient) Delete(name string) error {
	key := key(c.rootPath, name)

	_, err := c.consul.KV().Delete(key, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting %s", name)
	}
	return nil
}

func (c *upstreamsClient) Get(name string) (*v1.Upstream, error) {
	key := key(c.rootPath, name)
	p, _, err := c.consul.KV().Get(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "getting pair for for key %v", key)
	}
	if p == nil {
		return nil, errors.Errorf("keypair %s not found for upstream %s", key, name)
	}
	us, err := upstreamFromKVPair(p)
	if err != nil {
		return nil, errors.Wrap(err, "converting consul kv-pair to upstream")
	}
	return us, nil
}

func (c *upstreamsClient) List() ([]*v1.Upstream, error) {
	pairs, _, err := c.consul.KV().List(c.rootPath, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "listing key-value pairs for root %s", c.rootPath)
	}
	var upstreams []*v1.Upstream
	for _, p := range pairs {
		us, err := upstreamFromKVPair(p)
		if err != nil {
			return nil, errors.Wrapf(err, "converting %s to upstream", p.Key)
		}
		upstreams = append(upstreams, us)
	}
	return upstreams, nil
}

func (c *upstreamsClient) Watch(handlers ...storage.UpstreamEventHandler) (*storage.Watcher, error) {
	_, meta, err := c.consul.KV().List(c.rootPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "getting initial list")
	}
	lastIndex := meta.LastIndex
	sync := func() error {
		pairs, meta, err := c.consul.KV().List(c.rootPath, nil)
		if err != nil {
			return errors.Wrap(err, "getting kv-pairs list")
		}
		// no change since last poll
		if lastIndex == meta.LastIndex {
			return nil
		}
		var upstreams []*v1.Upstream
		for _, p := range pairs {
			us, err := upstreamFromKVPair(p)
			if err != nil {
				return errors.Wrapf(err, "converting %s to upstream", p.Key)
			}
			upstreams = append(upstreams, us)
		}
		lastIndex = meta.LastIndex
		for _, h := range handlers {
			// TODO: be clear that watch for consul only calls update
			h.OnUpdate(upstreams, nil)
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
				runtime.HandleError(fmt.Errorf("failed to start watcher to: %v", err))
				return
			case <-stop:
				return
			}
		}
	}), nil
}
