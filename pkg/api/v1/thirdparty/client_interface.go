package thirdparty

import (
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ThirdPartyResource interface {
	GetMetadata() core.Metadata
	SetMetadata(meta core.Metadata)
	Data() string
	IsSecret() bool
}

type Data struct {
	core.Metadata
	data string
}

type Secret struct {
	Data
}

type Artifact struct {
	Data
}

func (a *Data) GetMetadata() core.Metadata {
	return a.Metadata
}

func (a *Data) SetMetadata(meta core.Metadata) {
	a.Metadata = meta
}

func (a *Data) Data() string {
	return a.data
}

func (a *Artifact) IsSecret() bool {
	return false
}

func (a *Secret) IsSecret() bool {
	return true
}

var _ ThirdPartyResource = &Secret{}
var _ ThirdPartyResource = &Artifact{}

const DefaultNamespace = "default"

var DefaultRefreshRate = time.Second * 30

func DefaultNamespaceIfEmpty(namespace string) string {
	if namespace == "" {
		return DefaultNamespace
	}
	return namespace
}

type ThirdPartyResourceClient interface {
	Kind() string
	NewResource() ThirdPartyResource
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (ThirdPartyResource, error)
	Write(resource ThirdPartyResource, opts clients.WriteOpts) (ThirdPartyResource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]ThirdPartyResource, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []ThirdPartyResource, <-chan error, error)
}
