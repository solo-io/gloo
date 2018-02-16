package v1

import (
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
)

type UpstreamSpec *types.Struct
type FunctionSpec *types.Struct

type GlueObject interface {
	proto.Message
	GetMetadata() *Metadata
}
