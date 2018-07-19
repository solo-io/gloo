package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:generate protoc -I=./ -I=${GOPATH}/src/github.com/gogo/protobuf/ -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf/ -I=${GOPATH}/src --gogo_out=Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${GOPATH}/src mock_resource.proto

func (r *MockResource) SetStatus(status core.Status) {
	r.Status = status
}

func (r *MockResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &MockResource{}

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