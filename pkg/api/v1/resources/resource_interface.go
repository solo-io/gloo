package resources

import (
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type Resource interface {
	proto.Message
	GetName() string
	GetStatus() core.Status
	GetMetadata() core.Metadata
	SetName(name string)
	SetStatus(status core.Status)
	SetMetadata(meta core.Metadata)
}
