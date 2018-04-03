package consul

import (
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
)

// TODO: evaluate efficiency of LSing a whole dir on every op
// so far this is preferable to caring what files are named
type upstreamsClient struct {
	base *baseClient
}

func (c *upstreamsClient) Create(item *v1.Upstream) (*v1.Upstream, error) {
	out, err := c.base.Create(item, configObjectTypeUpstream)
	if err != nil {
		return nil, err
	}
	return out.(*v1.Upstream), nil
}

func (c *upstreamsClient) Update(item *v1.Upstream) (*v1.Upstream, error) {
	out, err := c.base.Update(item, configObjectTypeUpstream)
	if err != nil {
		return nil, err
	}
	return out.(*v1.Upstream), nil
}

func (c *upstreamsClient) Delete(name string) error {
	return c.base.Delete(name)
}

func (c *upstreamsClient) Get(name string) (*v1.Upstream, error) {
	out, err := c.base.Get(name, configObjectTypeUpstream)
	if err != nil {
		return nil, err
	}
	return out.(*v1.Upstream), nil
}

func (c *upstreamsClient) List() ([]*v1.Upstream, error) {
	list, err := c.base.List(configObjectTypeUpstream)
	if err != nil {
		return nil, err
	}
	var upstreams []*v1.Upstream
	for _, obj := range list {
		upstreams = append(upstreams, obj.(*v1.Upstream))
	}
	return upstreams, nil
}

func (c *upstreamsClient) Watch(handlers ...storage.UpstreamEventHandler) (*storage.Watcher, error) {
	var configObjHandlers []storage.ConfigObjectEventHandler
	for _, h := range handlers {
		configObjHandlers = append(configObjHandlers, storage.ConfigObjectEventHandler{UpstreamEventHandler: h})
	}
	return c.base.Watch(configObjectTypeUpstream, configObjHandlers...)
}
