package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/fileutils"
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
	if opts.Namespace == "" {
		opts.Namespace = clients.DefaultNamespace
	}
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
	if opts.Namespace == "" {
		opts.Namespace = clients.DefaultNamespace
	}
	path := rc.filename(opts.Namespace, name)
	err := os.Remove(path)
	switch {
	case err == nil :
		return nil
	case os.IsNotExist(err) && opts.IgnoreNotExist:
		return nil
	case os.IsNotExist(err) && !opts.IgnoreNotExist:
		return errors.NewNotExistErr(opts.Namespace, name, err)
	}
	return errors.Wrapf(err, "deleting resource %v", name)
}

func (rc *ResourceClient) List(opts clients.ListOpts) ([]resources.Resource, error) {
	if opts.Namespace == "" {
		opts.Namespace = clients.DefaultNamespace
	}

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

func (rc *ResourceClient) Watch(opts clients.WatchOpts) (<-chan []resources.Resource, error) {
	panic("yay")
}

func (rc *ResourceClient) filename(namespace, name string) string {
	return filepath.Join(rc.dir, namespace, name) + ".yaml"
}

func (rc *ResourceClient) newResource() resources.Resource {
	return proto.Clone(rc.resourceType).(resources.Resource)
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr)
}
