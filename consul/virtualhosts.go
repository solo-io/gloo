package consul

import (
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
)

// TODO: evaluate efficiency of LSing a whole dir on every op
// so far this is preferable to caring what files are named
type virtualHostsClient struct {
	base *baseClient
}

func (c *virtualHostsClient) Create(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	out, err := c.base.Create(item, configObjectTypeVirtualHost)
	if err != nil {
		return nil, err
	}
	return out.(*v1.VirtualHost), nil
}

func (c *virtualHostsClient) Update(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	out, err := c.base.Update(item, configObjectTypeVirtualHost)
	if err != nil {
		return nil, err
	}
	return out.(*v1.VirtualHost), nil
}

func (c *virtualHostsClient) Delete(name string) error {
	return c.base.Delete(name)
}

func (c *virtualHostsClient) Get(name string) (*v1.VirtualHost, error) {
	out, err := c.base.Get(name, configObjectTypeVirtualHost)
	if err != nil {
		return nil, err
	}
	return out.(*v1.VirtualHost), nil
}

func (c *virtualHostsClient) List() ([]*v1.VirtualHost, error) {
	list, err := c.base.List(configObjectTypeVirtualHost)
	if err != nil {
		return nil, err
	}
	var virtualhosts []*v1.VirtualHost
	for _, obj := range list {
		virtualhosts = append(virtualhosts, obj.(*v1.VirtualHost))
	}
	return virtualhosts, nil
}

func (c *virtualHostsClient) Watch(handlers ...storage.VirtualHostEventHandler) (*storage.Watcher, error) {
	var configObjHandlers []storage.ConfigObjectEventHandler
	for _, h := range handlers {
		configObjHandlers = append(configObjHandlers, storage.ConfigObjectEventHandler{VirtualHostEventHandler: h})
	}
	return c.base.Watch(configObjectTypeVirtualHost, configObjHandlers...)
}
