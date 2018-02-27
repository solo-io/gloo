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
	// headers that will be used to derive the data for processing the output templates
	// if no syntax containing {variables} are detected in the header value,
	// the whole value will be substituted by its name into the template
	// for example:
	/*
		input:
			header_parmeters:
		      x-header-foo: bar
		output:
			body_template: "{\"path\": {{ path }}}"
	*/
	HeaderParameters   map[string]string `json:"header_parameters"`
	PathParameter      string            `json:"path_parameter,omitempty"`
	MethodParameter    string            `json:"method_parameter"`
	SchemeParameter    string            `json:"scheme_parameter"`
	AuthorityParameter string            `json:"authority_parameter"`
	//TODO: support query params
}

// SPECIAL VARIABLE NAMES:
// these variables will be provided even if there is no parameter source
// path
// method
// scheme
// authority
type Output struct {
	HeaderTemplates map[string]string `json:"header_templates"`
	BodyTemplate    string            `json:"body_template"`
	//TODO: support query template
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
