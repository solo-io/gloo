package transformation

import (
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/protoutil"
)

type TransformationSpec struct {
	Input  Input  `json:"input"`
	Output Output `json:"output"`
}

type Input struct {
	HeaderParameters map[string]string `json:"header_parameters"`
	// we will parse this
	// get X capture groups
	// generate regex string from it
	// and header == :path
	PathParameter string `json:"path_parameter,omitempty"`
}

type Output struct {
	HeaderTemplates map[string]string `json:"header_templates"`
	BodyTemplate    string            `json:"body_template"`
}

func DecodeTransformationSpec(generic *types.Struct) (TransformationSpec, error) {
	var s TransformationSpec
	err := protoutil.UnmarshalStruct(generic, &s)
	return s, err
}

func EncodeUpstreamSpec(spec TransformationSpec) *types.Struct {
	v1Spec, err := protoutil.MarshalStruct(spec)
	if err != nil {
		panic(err)
	}
	return v1Spec
}
