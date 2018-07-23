package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

//go:generate protoc -I=./ -I=${GOPATH}/src/github.com/gogo/protobuf/ -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ -I=${GOPATH}/src --gogo_out=Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${GOPATH}/src --plugin=protoc-gen-solo-kit=${GOPATH}/bin/protoc-gen-solo-kit --solo-kit_out=. mock_resources.proto

func NewMockResource(namespace, name string) *MockResource {
	return &MockResource{
		Data: name,
		Metadata: core.Metadata{
			Name: name,
		},
	}
}
