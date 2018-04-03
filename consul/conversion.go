package consul

import (
	"strconv"

	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
)

func key(rootPath, itemName string) string {
	return rootPath + "/" + itemName
}

func toKVPair(rootPath string, item v1.ConfigObject) (*api.KVPair, error) {
	data, err := proto.Marshal(item)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling cfg object to proto")
	}
	var modifyIndex uint64
	if meta := item.GetMetadata(); meta != nil {
		if meta.ResourceVersion != "" {
			if i, err := strconv.Atoi(meta.ResourceVersion); err == nil {
				modifyIndex = uint64(i)
			}
		}
	}
	return &api.KVPair{
		Key:         key(rootPath, item.GetName()),
		Value:       data,
		ModifyIndex: modifyIndex,
	}, nil
}

func setResourceVersion(item v1.ConfigObject, p *api.KVPair, t ConfigObjectType) {
	resourceVersion := fmt.Sprintf("%v", p.ModifyIndex)
	switch t {
	case configObjectTypeUpstream:
		if item.(*v1.Upstream).Metadata == nil {
			item.(*v1.Upstream).Metadata = &v1.Metadata{}
		}
		item.(*v1.Upstream).Metadata.ResourceVersion = resourceVersion
	case configObjectTypeVirtualHost:
		if item.(*v1.VirtualHost).Metadata == nil {
			item.(*v1.VirtualHost).Metadata = &v1.Metadata{}
		}
		item.(*v1.VirtualHost).Metadata.ResourceVersion = resourceVersion
	default:
		panic("invalid type: " + t)
	}
}

func configObjectFromKVPair(p *api.KVPair, t ConfigObjectType) (v1.ConfigObject, error) {
	var item v1.ConfigObject
	switch t {
	case configObjectTypeUpstream:
		item = &v1.Upstream{}
	case configObjectTypeVirtualHost:
		item = &v1.VirtualHost{}
	default:
		panic("invalid type: " + t)
	}
	err := proto.Unmarshal(p.Value, item)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling cfg object from proto")
	}
	setResourceVersion(item, p, t)
	return item, nil
}
