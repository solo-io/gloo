package consul

import (
	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
)

func ToKVPair(rootPath string, item v1.ConfigObject) (*api.KVPair, error) {
	data, err := proto.Marshal(item)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling cfg object to proto")
	}
	return &api.KVPair{
		Key:   rootPath + "/" + item.GetName(),
		Value: data,
	}, nil
}
