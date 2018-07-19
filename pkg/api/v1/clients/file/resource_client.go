package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"

	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/fileutils"
	"v/github.com/radovskyb/watcher@v1.0.2"
)

type ResourceClient struct {
	dir          string
	resourceType resources.Resource
}

func NewResourceClient(dir string, resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		dir:          dir,
		resourceType: resourceType,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(name string, opts clients.GetOpts) (resources.Resource, error) {
	if err := resources.ValidateName(name); err != nil {
		return nil, errors.Wrapf(err, "validation error")
	}
	opts = opts.WithDefaults()
	path := rc.filename(opts.Namespace, name)
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil, errors.NewNotExistErr(opts.Namespace, name, err)
	}
	resource := rc.newResource()
	if err := fileutils.ReadFileInto(path, resource); err != nil {
		return nil, errors.Wrapf(err, "reading file into %v", reflect.TypeOf(rc.resourceType))
	}
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
	path := rc.filename(meta.Namespace, meta.Name)

	if !opts.OverwriteExisting {
		if _, err := os.Stat(path); err == nil {
			return nil, errors.NewExistErr(resource.GetMetadata())
		}
	}

	// mutate and return clone
	clone := proto.Clone(resource).(resources.Resource)
	// initialize or increment resource version
	meta.ResourceVersion = newOrIncrementResourceVer(meta.ResourceVersion)
	clone.SetMetadata(meta)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil && !os.IsExist(err) {
		return nil, errors.Wrapf(err, "creating directory")
	}
	if err := fileutils.WriteToFile(path, clone); err != nil {
		return nil, errors.Wrapf(err, "writing file")
	}
	return clone, nil
}

func (rc *ResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	path := rc.filename(opts.Namespace, name)
	err := os.Remove(path)
	switch {
	case err == nil:
		return nil
	case os.IsNotExist(err) && opts.IgnoreNotExist:
		return nil
	case os.IsNotExist(err) && !opts.IgnoreNotExist:
		return errors.NewNotExistErr(opts.Namespace, name, err)
	}
	return errors.Wrapf(err, "deleting resource %v", name)
}

func (rc *ResourceClient) List(opts clients.ListOpts) ([]resources.Resource, error) {
	opts = opts.WithDefaults()

	namespaceDir := filepath.Join(rc.dir, opts.Namespace)
	files, err := ioutil.ReadDir(namespaceDir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading namespace dir")
	}

	var resourceList []resources.Resource
	for _, file := range files {
		resource := rc.newResource()
		path := filepath.Join(namespaceDir, file.Name())
		if err := fileutils.ReadFileInto(path, resource); err != nil {
			return nil, errors.Wrapf(err, "reading file into %v", reflect.TypeOf(rc.resourceType))
		}
		resourceList = append(resourceList, resource)
	}

	sort.SliceStable(resourceList, func(i, j int) bool {
		return resourceList[i].GetMetadata().Name < resourceList[j].GetMetadata().Name
	})

	return resourceList, nil
}

func (rc *ResourceClient) Watch(opts clients.WatchOpts) (<-chan []resources.Resource, <-chan error, error) {
	opts = opts.WithDefaults()
	dir := filepath.Join(rc.dir, opts.Namespace)
	events, errs, err := rc.events(opts.Ctx, dir, opts.RefreshRate)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting watch on namespace dir")
	}
	resourceLists := make(chan []resources.Resource)
	go func() {
		select {
		case <-events:
			list, err := rc.List(clients.ListOpts{
				Ctx:       opts.Ctx,
				Selector:  opts.Selector,
				Namespace: opts.Namespace,
			})
			if err != nil {
				errs <- err
			}
			resourceLists <- list
		}
	}()

	return resourceLists, errs, nil
}

func (rc *ResourceClient) filename(namespace, name string) string {
	return filepath.Join(rc.dir, namespace, name) + ".yaml"
}

func (rc *ResourceClient) newResource() resources.Resource {
	return proto.Clone(rc.resourceType).(resources.Resource)
}

func (rc *ResourceClient) events(ctx context.Context, dir string, refreshRate time.Duration) (<-chan struct{}, chan error, error) {
	events := make(chan struct{})
	errs := make(chan error)
	w := watcher.New()
	w.SetMaxEvents(0)
	w.FilterOps(watcher.Create, watcher.Write, watcher.Remove, watcher.Rename, watcher.Move)
	if err := w.AddRecursive(dir); err != nil {
		return nil, nil, errors.Wrapf(err, "failed to watch directory %v", dir)
	}
	go func() {
		log.Printf("watching dir %v", dir)
		if err := w.Start(refreshRate); err != nil {
			errs <- err
		}
	}()
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Printf("event: %v", event.String())
				if event.IsDir() {
					continue
				}
				events <- struct{}{}
			case err := <-w.Error:
				errs <- errors.Wrapf(err, "file watcher error")
			case <-ctx.Done():
				w.Close()
				return
			}
		}
	}()
	return events, errs, nil
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr)
}
