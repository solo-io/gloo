package thirdparty

import (
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

const DefaultNamespace = "default"

var DefaultRefreshRate = time.Second * 30

func DefaultNamespaceIfEmpty(namespace string) string {
	if namespace == "" {
		return DefaultNamespace
	}
	return namespace
}

type Secret struct{}

type SecretClient interface {
	Kind() string
	NewResource() resources.Resource
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (resources.Resource, error)
	Write(resource resources.Resource, opts clients.WriteOpts) (resources.Resource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]resources.Resource, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []resources.Resource, <-chan error, error)
}

/*
thirdarty resource can be either secret or plaintext
same object same interface
*/
