package base

import (
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/dependencies"
)

type StorableItem struct {
	Upstream    *v1.Upstream
	VirtualHost *v1.VirtualHost
	File        *dependencies.File
}

func (item *StorableItem) GetName() string {
	switch {
	case item.Upstream != nil:
		return item.Upstream.GetName()
	case item.VirtualHost != nil:
		return item.VirtualHost.GetName()
	case item.File != nil:
		return item.File.Ref
	default:
		panic("virtual host, file or upstream must be set")
	}
}

func (item *StorableItem) GetResourceVersion() string {
	switch {
	case item.Upstream != nil:
		if item.Upstream.GetMetadata() == nil {
			return ""
		}
		return item.Upstream.GetMetadata().GetResourceVersion()
	case item.VirtualHost != nil:
		if item.VirtualHost.GetMetadata() == nil {
			return ""
		}
		return item.VirtualHost.GetMetadata().GetResourceVersion()
	case item.File != nil:
		return item.File.ResourceVersion
	default:
		panic("virtual host, file or upstream must be set")
	}
}

func (item *StorableItem) SetResourceVersion(rv string) {
	switch {
	case item.Upstream != nil:
		if item.Upstream.GetMetadata() == nil {
			item.Upstream.Metadata = &v1.Metadata{}
		}
		item.Upstream.Metadata.ResourceVersion = rv
	case item.VirtualHost != nil:
		if item.VirtualHost.GetMetadata() == nil {
			item.VirtualHost.Metadata = &v1.Metadata{}
		}
		item.VirtualHost.Metadata.ResourceVersion = rv
	case item.File != nil:
		item.File.ResourceVersion = rv
	default:
		panic("virtual host, file or upstream must be set")
	}
}

func (item *StorableItem) GetBytes() ([]byte, error) {
	switch {
	case item.Upstream != nil:
		return proto.Marshal(item.Upstream)
	case item.VirtualHost != nil:
		return proto.Marshal(item.VirtualHost)
	case item.File != nil:
		return item.File.Contents, nil
	default:
		panic("virtual host, file or upstream must be set")
	}
}

func (item *StorableItem) GetTypeFlag() StorableItemType {
	switch {
	case item.Upstream != nil:
		return StorableItemTypeUpstream
	case item.VirtualHost != nil:
		return StorableItemTypeVirtualHost
	case item.File != nil:
		return StorableItemTypeFile
	default:
		panic("virtual host, file or upstream must be set")
	}
}

type StorableItemType uint64

const (
	StorableItemTypeUpstream StorableItemType = iota
	StorableItemTypeVirtualHost
	StorableItemTypeFile
)

type StorableItemEventHandler struct {
	UpstreamEventHandler    storage.UpstreamEventHandler
	VirtualHostEventHandler storage.VirtualHostEventHandler
	FileEventHandler        dependencies.FileEventHandler
}
