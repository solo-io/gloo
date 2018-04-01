package grpc

import (
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type ServiceProperties struct {
	// the name of the gRPC services defined in the descriptors (to route to)
	GRPCServiceNames []string `json:"service_names"`
	// file ref for the proto descriptors generated for is gRPC service
	DescriptorsFileRef string `json:"descriptors_file_ref"`
}

func DecodeServiceProperties(generic *types.Struct) (ServiceProperties, error) {
	var p ServiceProperties
	err := protoutil.UnmarshalStruct(generic, &p)
	return p, err
}

func EncodeServiceProperties(properties ServiceProperties) *types.Struct {
	v1Properties, err := protoutil.MarshalStruct(properties)
	if err != nil {
		panic(err)
	}
	return v1Properties
}
