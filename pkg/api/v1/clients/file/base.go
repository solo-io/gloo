package file

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"time"
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
func (rc *ResourceClient) Create(resource resources.Resource, opts *clients.CreateOptions) (resources.Resource, error) {panic("yay")}
func (rc *ResourceClient) Update(resource resources.Resource, opts *clients.UpdateOptions) (resources.Resource, error) {panic("yay")}
func (rc *ResourceClient) Delete(name string, opts *clients.DeleteOptions) error {panic("yay")}
func (rc *ResourceClient) List(opts *clients.ListOptions) ([]resources.Resource, error) {panic("yay")}
func (rc *ResourceClient) Watch(opts *clients.WatchOptions) (<-chan []resources.Resource, error) {panic("yay")}
