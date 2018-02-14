package file

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/solo-io/glue-storage"
	"github.com/solo-io/glue-storage/common"
	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	yaml "gopkg.in/yaml.v2"
)

const (
	upstreamPath     = "upstream"
	vhostPath        = "vhost"
	folderAtrributes = 0777
	fileAttributes   = 0644
)

var evtTypeToOperation = map[fsnotify.Op]storage.WatchOperation{
	fsnotify.Create: storage.CreateOp,
	fsnotify.Write:  storage.UpdateOp,
	fsnotify.Remove: storage.DeleteOp,
	fsnotify.Chmod:  storage.UpdateOp,
	fsnotify.Rename: storage.ErrorOp,
}

type FileStorage struct {
	namespacePath   string
	upstreamWatcher *fsnotify.Watcher
	vhostWatcher    *fsnotify.Watcher
}

func NewFileStorage(root, namespace string) (*FileStorage, error) {
	fullpath := path.Join(root, namespace)
	return &FileStorage{namespacePath: fullpath}, nil
}

func (c *FileStorage) Register(item storage.Item) error {
	switch item.(type) {
	case *gluev1.Upstream:
		err := createFolder(path.Join(c.namespacePath, upstreamPath))
		if err != nil {
			return err
		}
	case *gluev1.VirtualHost:
		err := createFolder(path.Join(c.namespacePath, vhostPath))
		if err != nil {
			return err
		}
	default:
		return common.UnknownType(item)
	}
	return nil
}

func (c *FileStorage) Create(item storage.Item) (storage.Item, error) {
	d, err := yaml.Marshal(item)
	if err != nil {
		return nil, err
	}
	n, err := getObjName(item)
	if err != nil {
		return nil, err
	}
	switch item.(type) {
	case *gluev1.Upstream:
		err = ioutil.WriteFile(path.Join(c.namespacePath, *n), d, fileAttributes)
	case *gluev1.VirtualHost:
		err = ioutil.WriteFile(path.Join(c.namespacePath, *n), d, fileAttributes)
	default:
		return nil, common.UnknownType(item)
	}
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *FileStorage) Update(item storage.Item) (storage.Item, error) {
	return c.Create(item)
}

func (c *FileStorage) Delete(item storage.Item) error {
	n, err := getObjName(item)
	if err != nil {
		return err
	}
	return os.Remove(path.Join(c.namespacePath, *n))
}

func (c *FileStorage) Get(item storage.Item, getOptions *storage.GetOptions) (storage.Item, error) {
	n, err := getObjName(item)
	if err != nil {
		return nil, err
	}
	return c.read(item, *n, c.namespacePath, getOptions)
}

func (c *FileStorage) List(item storage.Item, listOptions *storage.ListOptions) ([]storage.Item, error) {
	res := make([]storage.Item, 0)
	switch item.(type) {
	case *gluev1.Upstream:
		flist, err := ioutil.ReadDir(path.Join(c.namespacePath, upstreamPath))
		if err != nil {
			return nil, err
		}
		for _, finfo := range flist {
			if !finfo.IsDir() && finfo.Size() > 0 {
				obj, err := c.read(item, finfo.Name(), path.Join(c.namespacePath, upstreamPath), nil)
				if err != nil {
					return nil, err
				}
				res = append(res, obj)
			}
		}
		return res, nil
	case *gluev1.VirtualHost:
		flist, err := ioutil.ReadDir(path.Join(c.namespacePath, vhostPath))
		if err != nil {
			return nil, err
		}
		for _, finfo := range flist {
			if !finfo.IsDir() && finfo.Size() > 0 {
				obj, err := c.read(item, finfo.Name(), path.Join(c.namespacePath, vhostPath), nil)
				if err != nil {
					return nil, err
				}
				res = append(res, obj)
			}
		}
		return res, nil
	default:
		return nil, common.UnknownType(item)
	}
}
func (c *FileStorage) Watch(item storage.Item, watchOptions *storage.WatchOptions, callback func(item storage.Item, operation storage.WatchOperation)) error {

	var err error
	switch item.(type) {
	case *gluev1.Upstream:
		c.upstreamWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case event, ok := <-c.upstreamWatcher.Events:
					if !ok {
						return
					}
					if event.Op == fsnotify.Remove {
						callback(&gluev1.Upstream{Name: path.Base(event.Name)}, storage.DeleteOp)
					} else {
						obj, err := c.read(item, event.Name, "", nil)
						if err != nil {
							log.Println(err)
							continue
						}
						callback(obj, evtTypeToOperation[event.Op])
					}
				case err = <-c.upstreamWatcher.Errors:
					log.Println(err)
				}
			}
		}()
		c.upstreamWatcher.Add(path.Join(c.namespacePath, upstreamPath))

	case *gluev1.VirtualHost:
		c.vhostWatcher, err = fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		go func() {
			for {
				select {
				case event, ok := <-c.vhostWatcher.Events:
					if !ok {
						return
					}
					if event.Op == fsnotify.Remove {
						callback(&gluev1.VirtualHost{Name: path.Base(event.Name)}, storage.DeleteOp)
					} else {
						obj, err := c.read(item, event.Name, "", nil)
						if err != nil {
							log.Println(err)
							continue
						}
						callback(obj, evtTypeToOperation[event.Op])
					}
				case err = <-c.vhostWatcher.Errors:
					log.Println(err)
				}
			}
		}()
		c.vhostWatcher.Add(path.Join(c.namespacePath, vhostPath))

	default:
		return common.UnknownType(item)
	}
	return nil
}

func (c *FileStorage) WatchStop() {
	if c.upstreamWatcher != nil {
		c.upstreamWatcher.Close()
	}

	if c.vhostWatcher != nil {
		c.vhostWatcher.Close()
	}
}

func (c *FileStorage) read(item storage.Item, name, fpath string, getOptions *storage.GetOptions) (storage.Item, error) {

	d, err := ioutil.ReadFile(path.Join(fpath, name))
	if err != nil {
		return nil, err
	}
	if obj, ok := item.(*gluev1.Upstream); ok {
		err = yaml.Unmarshal(d, obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		err = yaml.Unmarshal(d, obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return nil, common.UnknownType(item)
}

func getObjName(item storage.Item) (*string, error) {
	if obj, ok := item.(*gluev1.Upstream); ok {
		p := path.Join(upstreamPath, obj.Name)
		return &p, nil
	} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		p := path.Join(vhostPath, obj.Name)
		return &p, nil
	}
	return nil, common.UnknownType(item)
}

func createFolder(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return os.MkdirAll(name, folderAtrributes)
	} else if err != nil {
		return err
	}
	return nil
}
