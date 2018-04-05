package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"

	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/dependencies"
	"github.com/solo-io/gloo/pkg/log"
)

type fileStorage struct {
	rootPath      string
	consul        *api.Client
	syncFrequency time.Duration
}

func NewFileStorage(cfg *api.Config, rootPath string, syncFrequency time.Duration) (dependencies.FileStorage, error) {
	// Get a new client
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "creating consul client")
	}

	return &fileStorage{
		consul:        client,
		rootPath:      rootPath + "/files",
		syncFrequency: syncFrequency,
	}, nil
}

func copyFile(file *dependencies.File) *dependencies.File {
	contents := make([]byte, len(file.Contents))
	copy(contents, file.Contents)
	return &dependencies.File{
		Ref:             file.Ref,
		Contents:        contents,
		ResourceVersion: file.ResourceVersion,
	}
}

func (s *fileStorage) Create(file *dependencies.File) (*dependencies.File, error) {
	if _, err := s.Get(file.Ref); err == nil {
		return nil, errors.Errorf("file %v already exists", file.Ref)
	}
	p := toKVPair(s.rootPath, file)

	// error if the key already exists
	existingP, _, err := s.consul.KV().Get(p.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query consul")
	}
	if existingP != nil {
		return nil, storage.NewAlreadyExistsErr(
			errors.Errorf("key found for file %s: %s", file.Ref, p.Key))
	}
	if _, err := s.consul.KV().Put(p, nil); err != nil {
		return nil, errors.Wrapf(err, "writing kv pair %s", p.Key)
	}
	// set the resource version from the CreateIndex of the created kv pair
	p, _, err = s.consul.KV().Get(p.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "getting newly created kv pair %s", p.Key)
	}
	copied := copyFile(file)
	setResourceVersion(copied, p)
	return copied, nil
}

func (s *fileStorage) Update(file *dependencies.File) (*dependencies.File, error) {
	updatedP := toKVPair(s.rootPath, file)

	// error if the key doesn't already exist
	exsitingP, _, err := s.consul.KV().Get(updatedP.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to query consul")
	}
	if exsitingP == nil {
		return nil, errors.Errorf("key not found for configObject %s: %s", file.Ref, updatedP.Key)
	}
	if success, _, err := s.consul.KV().CAS(updatedP, nil); err != nil {
		return nil, errors.Wrapf(err, "writing kv pair %s", updatedP.Key)
	} else if !success {
		return nil, errors.Errorf("resource version was invalid for configObject: %s", file.Ref)
	}

	// set the resource version from the CreateIndex of the created kv pair
	updatedP, _, err = s.consul.KV().Get(updatedP.Key, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "getting newly created kv pair %s", updatedP.Key)
	}
	copied := copyFile(file)
	setResourceVersion(copied, updatedP)
	return copied, nil
}

func (s *fileStorage) Delete(name string) error {
	key := key(s.rootPath, name)

	_, err := s.consul.KV().Delete(key, nil)
	if err != nil {
		return errors.Wrapf(err, "deleting %s", name)
	}
	return nil
}

func (s *fileStorage) Get(name string) (*dependencies.File, error) {
	key := key(s.rootPath, name)
	p, _, err := s.consul.KV().Get(key, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "getting pair for for key %v", key)
	}
	if p == nil {
		return nil, errors.Errorf("keypair %s not found for file %s", key, name)
	}
	file := fileFromKVPair(s.rootPath, p)
	return file, nil
}

func (s *fileStorage) List() ([]*dependencies.File, error) {
	pairs, _, err := s.consul.KV().List(s.rootPath, &api.QueryOptions{RequireConsistent: true})
	if err != nil {
		return nil, errors.Wrapf(err, "listing key-value pairs for root %s", s.rootPath)
	}
	var files []*dependencies.File
	for _, p := range pairs {
		files = append(files, fileFromKVPair(s.rootPath, p))
	}
	return files, nil
}

// TODO: be clear that watch for consul only calls update
func (s *fileStorage) Watch(handlers ...dependencies.FileEventHandler) (*storage.Watcher, error) {
	var lastIndex uint64
	sync := func() error {
		pairs, meta, err := s.consul.KV().List(s.rootPath, &api.QueryOptions{RequireConsistent: true})
		if err != nil {
			return errors.Wrap(err, "getting kv-pairs list")
		}
		// no change since last poll
		if lastIndex == meta.LastIndex {
			return nil
		}
		var (
			files []*dependencies.File
		)
		for _, p := range pairs {
			files = append(files, fileFromKVPair(s.rootPath, p))
		}
		lastIndex = meta.LastIndex
		for _, h := range handlers {
			h.OnUpdate(files, nil)
		}
		return nil
	}
	return storage.NewWatcher(func(stop <-chan struct{}, errs chan error) {
		for {
			select {
			case <-time.After(s.syncFrequency):
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
