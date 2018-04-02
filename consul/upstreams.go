package consul

import (
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"

	"time"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
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

//func (c *upstreamsClient) Get(name string) (*v1.Upstream, error) {
//	upstreamFiles, err := c.pathsToUpstreams()
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to read upstream dir")
//	}
//	// error if exists already
//	for _, existingUps := range upstreamFiles {
//		if existingUps.Name == name {
//			return existingUps, nil
//		}
//	}
//	return nil, errors.Errorf("file not found for upstream %v", name)
//}
//
//func (c *upstreamsClient) List() ([]*v1.Upstream, error) {
//	upstreamPaths, err := c.pathsToUpstreams()
//	if err != nil {
//		return nil, err
//	}
//	var upstreams []*v1.Upstream
//	for _, up := range upstreamPaths {
//		upstreams = append(upstreams, up)
//	}
//	return upstreams, nil
//}
//
//func (c *upstreamsClient) pathsToUpstreams() (map[string]*v1.Upstream, error) {
//	files, err := ioutil.ReadDir(c.dir)
//	if err != nil {
//		return nil, errors.Wrap(err, "could not read dir")
//	}
//	upstreams := make(map[string]*v1.Upstream)
//	for _, f := range files {
//		path := filepath.Join(c.dir, f.Name())
//		if !strings.HasSuffix(path, ".yml") && !strings.HasSuffix(path, ".yaml") {
//			continue
//		}
//		var upstream v1.Upstream
//		err := ReadFileInto(path, &upstream)
//		if err != nil {
//			return nil, errors.Wrap(err, "unable to parse .yml file as upstream")
//		}
//		upstreams[path] = &upstream
//	}
//	return upstreams, nil
//}
//
//func (u *upstreamsClient) Watch(handlers ...storage.UpstreamEventHandler) (*storage.Watcher, error) {
//	w := watcher.New()
//	w.SetMaxEvents(0)
//	w.FilterOps(watcher.Create, watcher.Write, watcher.Remove)
//	if err := w.AddRecursive(u.dir); err != nil {
//		return nil, errors.Wrapf(err, "failed to add directory %v", u.dir)
//	}
//
//	return storage.NewWatcher(func(stop <-chan struct{}, errs chan error) {
//		go func() {
//			if err := w.Start(u.syncFrequency); err != nil {
//				errs <- err
//			}
//		}()
//		for {
//			select {
//			case event := <-w.Event:
//				if err := u.onEvent(event, handlers...); err != nil {
//					runtime.HandleError(err)
//				}
//			case err := <-w.Error:
//				runtime.HandleError(fmt.Errorf("watcher encoutnered error: %v", err))
//				return
//			case err := <-errs:
//				runtime.HandleError(fmt.Errorf("failed to start watcher to: %v", err))
//				return
//			case <-stop:
//				w.Close()
//				return
//			}
//		}
//	}), nil
//}
//
//func (u *upstreamsClient) onEvent(event watcher.Event, handlers ...storage.UpstreamEventHandler) error {
//	log.Debugf("file event: %v [%v]", event.Path, event.Op)
//	current, err := u.List()
//	if err != nil {
//		return err
//	}
//	if event.IsDir() {
//		return nil
//	}
//	switch event.Op {
//	case watcher.Create:
//		for _, h := range handlers {
//			var created v1.Upstream
//			err := ReadFileInto(event.Path, &created)
//			if err != nil {
//				return err
//			}
//			h.OnAdd(current, &created)
//		}
//	case watcher.Write:
//		for _, h := range handlers {
//			var updated v1.Upstream
//			err := ReadFileInto(event.Path, &updated)
//			if err != nil {
//				return err
//			}
//			h.OnUpdate(current, &updated)
//		}
//	case watcher.Remove:
//		for _, h := range handlers {
//			// can't read the deleted object
//			// callers beware
//			h.OnDelete(current, nil)
//		}
//	}
//	return nil
//}
