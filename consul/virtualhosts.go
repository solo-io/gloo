package consul

import (
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/internal/base"
)

type virtualHostsClient struct {
	base *base.ConsulStorageClient
}

func (c *virtualHostsClient) Create(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	out, err := c.base.Create(&base.StorableItem{VirtualHost: item})
	if err != nil {
		return nil, err
	}
	return out.VirtualHost, nil
}

func (c *virtualHostsClient) Update(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	out, err := c.base.Update(&base.StorableItem{VirtualHost: item})
	if err != nil {
		return nil, err
	}
	return out.VirtualHost, nil
}

func (c *virtualHostsClient) Delete(name string) error {
	return c.base.Delete(name)
}

func (c *virtualHostsClient) Get(name string) (*v1.VirtualHost, error) {
	out, err := c.base.Get(name)
	if err != nil {
		return nil, err
	}
	return out.VirtualHost, nil
}

func (c *virtualHostsClient) List() ([]*v1.VirtualHost, error) {
	list, err := c.base.List()
	if err != nil {
		return nil, err
	}
	var virtualHosts []*v1.VirtualHost
	for _, obj := range list {
		virtualHosts = append(virtualHosts, obj.VirtualHost)
	}
	return virtualHosts, nil
}

func (c *virtualHostsClient) Watch(handlers ...storage.VirtualHostEventHandler) (*storage.Watcher, error) {
	var baseHandlers []base.StorableItemEventHandler
	for _, h := range handlers {
		baseHandlers = append(baseHandlers, base.StorableItemEventHandler{VirtualHostEventHandler: h})
	}
	return c.base.Watch(baseHandlers...)
}
