package consul

import (
	"strconv"

	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
)

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
		Key:         rootPath + "/" + item.GetName(),
		Value:       data,
		ModifyIndex: modifyIndex,
	}, nil
}

func setResourceVersion(item *v1.Upstream, p *api.KVPair) {
	resourceVersion := fmt.Sprintf("%v", p.ModifyIndex)
	if item.Metadata == nil {
		item.Metadata = &v1.Metadata{}
	}
	item.Metadata.ResourceVersion = resourceVersion
}

func upstreamFromKVPair(p *api.KVPair) (*v1.Upstream, error) {
	var us v1.Upstream
	err := proto.Unmarshal(p.Value, &us)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling cfg object from proto")
	}
	setResourceVersion(&us, p)
	return &us, nil
}
