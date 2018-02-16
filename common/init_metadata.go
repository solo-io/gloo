package common

import (
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/types"
)

// sets anything that might be nil so we don't get a nil pointer / map somewhere
func InitFilterMetadata(filterName string, meta *envoycore.Metadata) {
	if meta.FilterMetadata == nil {
		meta.FilterMetadata = make(map[string]*types.Struct)
	}
	if meta.FilterMetadata[filterName] == nil {
		meta.FilterMetadata[filterName] = &types.Struct{
			Fields: make(map[string]*types.Value),
		}
	}
}

// sets anything that might be nil so we don't get a nil pointer / map somewhere
func InitFilterMetadataField(filterName, fieldName string, meta *envoycore.Metadata) {
	InitFilterMetadata(filterName, meta)
	if meta.FilterMetadata[filterName].Fields[fieldName] == nil {
		meta.FilterMetadata[filterName].Fields[fieldName] = &types.Value{}
	}
}
