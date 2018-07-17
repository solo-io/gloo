package file

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"time"
	"github.com/gogo/protobuf/proto"
	"strconv"
	"fmt"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type ResourceClient struct {
	dir         string
	refreshRate time.Duration
}

func NewResourceClient(dir string, refreshRate time.Duration) *ResourceClient {
	return &ResourceClient{
		dir:         dir,
		refreshRate: refreshRate,
	}
}

var _ clients.ResourceClient = &ResourceClient{}

func (rc *ResourceClient) Register() error {panic("yay")}

func (rc *ResourceClient) Get(name string, opts *clients.GetOptions) (resources.Resource, error) {panic("yay")}

func (rc *ResourceClient) Write(resource resources.Resource, opts *clients.WriteOptions) (resources.Resource, error) {
	if !opts.OverwriteExisting {
		getOpts := &clients.GetOptions{}
		if opts != nil {
			getOpts.Ctx = opts.Ctx
		}
		if _, err := rc.Get(resource.GetMetadata().Name, getOpts); err == nil {
			return nil, errors.NewAlreadyExistsErr(resource)
		}
	}
	// only mutate clone
	clone, ok := proto.Clone(resource).(resources.Resource)
	if !ok {
		panic("internal error: output of proto.Clone was not expected type")
	}
	meta := clone.GetMetadata()
	// initialize or increment resource version
	meta.ResourceVersion = newOrIncrementResourceVer(meta.ResourceVersion)
	clone.SetMetadata(meta)

}

func (rc *ResourceClient) Delete(name string, opts *clients.DeleteOptions) error {panic("yay")}

func (rc *ResourceClient) List(opts *clients.ListOptions) ([]resources.Resource, error) {panic("yay")}

func (rc *ResourceClient) Watch(opts *clients.WatchOptions) (<-chan []resources.Resource, error) {panic("yay")}

// util methods
func newOrIncrementResourceVer(resourceVersion string) string {
	curr, err := strconv.Atoi(resourceVersion)
	if err != nil {
		curr = 1
	}
	return fmt.Sprintf("%v", curr)
}
