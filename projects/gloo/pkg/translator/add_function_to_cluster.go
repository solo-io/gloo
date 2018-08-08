package translator

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/solo-io/gloo/pkg/coreplugins/common"

	"github.com/gogo/protobuf/types"
)

const (
	multiFunctionListDestinationKey = "functions"
)

// utility method for functional plugins
func AddFunctionToCluster(out *envoyapi.Cluster, funcName string, spec *types.Struct) {
	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	multiFunctionMetadata := getFunctionalFilterMetadata(multiFunctionListDestinationKey, out.Metadata)

	if multiFunctionMetadata.Fields[funcName] == nil {
		multiFunctionMetadata.Fields[funcName] = &types.Value{}
	}
	multiFunctionMetadata.Fields[funcName].Kind = &types.Value_StructValue{StructValue: spec}
}

func getFunctionalFilterMetadata(key string, meta *envoycore.Metadata) *types.Struct {
	initFunctionalFilterMetadata(key, meta)
	return meta.FilterMetadata[filterName].Fields[key].Kind.(*types.Value_StructValue).StructValue
}

// sets anything that might be nil so we don't get a nil pointer / map somewhere
func initFunctionalFilterMetadata(key string, meta *envoycore.Metadata) {
	filterMetadata := common.InitFilterMetadataField(filterName, key, meta)
	if filterMetadata.Kind == nil {
		filterMetadata.Kind = &types.Value_StructValue{}
	}
	_, isStructValue := filterMetadata.Kind.(*types.Value_StructValue)
	if !isStructValue {
		filterMetadata.Kind = &types.Value_StructValue{}
	}
	if filterMetadata.Kind.(*types.Value_StructValue).StructValue == nil {
		filterMetadata.Kind.(*types.Value_StructValue).StructValue = &types.Struct{}
	}
	if filterMetadata.Kind.(*types.Value_StructValue).StructValue.Fields == nil {
		filterMetadata.Kind.(*types.Value_StructValue).StructValue.Fields = make(map[string]*types.Value)
	}
}
