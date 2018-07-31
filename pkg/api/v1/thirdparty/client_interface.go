package thirdparty

import (
	"encoding/json"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ThirdPartyResource interface {
	resources.BaseResource
	GetData() map[string]string
	SetData(map[string]string)
	IsSecret() bool
	DeepCopy() ThirdPartyResource
}

type Data struct {
	core.Metadata
	Values map[string]string
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

func (a *Data) GetData() map[string]string {
	return a.Values
}

func (a *Data) SetData(values map[string]string) {
	a.Values = values
}

func (a *Artifact) DeepCopy() ThirdPartyResource {
	b, err := json.Marshal(a)
	if err != nil {
		panic("unexpected error marshalling " + err.Error())
	}
	var deepcopy Artifact
	err = json.Unmarshal(b, &deepcopy)
	if err != nil {
		panic("unexpected error unmarshalling " + err.Error())
	}
	return &deepcopy
}

func (a *Artifact) IsSecret() bool {
	return false
}

func (a *Secret) DeepCopy() ThirdPartyResource {
	b, err := json.Marshal(a)
	if err != nil {
		panic("unexpected error marshalling " + err.Error())
	}
	var deepcopy Artifact
	err = json.Unmarshal(b, &deepcopy)
	if err != nil {
		panic("unexpected error unmarshalling " + err.Error())
	}
	return &deepcopy
}

func (a *Secret) IsSecret() bool {
	return true
}

var _ ThirdPartyResource = &Secret{}
var _ ThirdPartyResource = &Artifact{}

type ThirdPartyResourceClient interface {
	Read(namespace, name string, opts clients.ReadOpts) (ThirdPartyResource, error)
	Write(resource ThirdPartyResource, opts clients.WriteOpts) (ThirdPartyResource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]ThirdPartyResource, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []ThirdPartyResource, <-chan error, error)
}
