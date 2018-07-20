package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

//go:generate protoc -I=./ -I=${GOPATH}/src/github.com/gogo/protobuf/ -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ -I=${GOPATH}/src --gogo_out=Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${GOPATH}/src mock_resource.proto

func NewMockResource(name string) *MockResource {
	return &MockResource{
		Data: name,
		Metadata: core.Metadata{
			Name: name,
		},
	}
}

type MockCrdObject struct {
	resources.Resource
}

func (m *MockCrdObject) GetObjectKind() schema.ObjectKind {
	t := MockCrd.TypeMeta()
	return &t
}

func (m *MockCrdObject) DeepCopyObject() runtime.Object {
	return &MockCrdObject{
		Resource: resources.Clone(m.Resource),
	}
}

var MockCrd = crd.NewCrd("testing.solo.io",
	"mocks",
	"testing.solo.io",
	"v1",
	"MockCrdObject",
	"mk",
	&MockCrdObject{})
