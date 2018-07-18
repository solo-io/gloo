package file

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"time"
	"github.com/gogo/protobuf/proto"
	"strconv"
	"fmt"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/fileutils"
	"path/filepath"
	"os"
	"reflect"
)

type ResourceClient struct {
	dir         string
	refreshRate time.Duration
	resourceType resources.Resource
}

func NewResourceClient(dir string, refreshRate time.Duration, resourceType resources.Resource) *ResourceClient {
	return &ResourceClient{
		dir:         dir,
		refreshRate: refreshRate,
		resourceType: resourceType,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Register() error {
	return nil
}

func (rc *ResourceClient) Read(name string, opts clients.GetOptions) (resources.Resource, error) {
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
	into := proto.Clone(rc.resourceType).(resources.Resource)
	if err := fileutils.ReadFileInto(path, into); err != nil {
		return nil, errors.Wrapf(err, "reading file into %v", reflect.TypeOf(rc.resourceType).Name())
	}
	return into, nil
}

func (rc *ResourceClient) Write(resource resources.Resource, opts clients.WriteOptions) (resources.Resource, error) {
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

func (rc *ResourceClient) Delete(name string, opts clients.DeleteOptions) error { panic("yay") }

func (rc *ResourceClient) List(opts clients.ListOptions) ([]resources.Resource, error) {
	panic("yay")
}

func (rc *ResourceClient) Watch(opts clients.WatchOptions) (<-chan []resources.Resource, error) { panic("yay") }

func (rc *ResourceClient) filename(namespace, name string) string {
	return filepath.Join(rc.dir, namespace, name)
}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr)
}
