package swagger

import (
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type UpstreamSpec struct {
	SwaggerURI string `json:"swagger_uri"`
}

func DecodeUpstreamSpec(generic v1.UpstreamSpec) (*UpstreamSpec, error) {
	s := new(UpstreamSpec)
	if err := protoutil.UnmarshalStruct(generic, s); err != nil {
		return nil, err
	}
	return s, nil
}

func EncodeUpstreamSpec(spec UpstreamSpec) *types.Struct {
	pb, err := protoutil.MarshalStruct(spec)
	if err != nil {
		panic(err)
	}
	return pb
}

func IsSwagger(us *v1.Upstream) bool {
	spec, err := DecodeUpstreamSpec(us.Spec)
	if err != nil {
		return false
	}
	return spec != nil
}
