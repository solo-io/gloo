package v1

import (
	"github.com/gogo/protobuf/proto"
)

type Resource interface {
	proto.Message
	GetName() string
	GetStatus() Status
	GetMetadata() Metadata
	SetName(name string)
	SetStatus(status Status)
	SetMetadata(meta Metadata)
}
