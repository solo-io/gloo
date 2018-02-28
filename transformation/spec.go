package transformation

import (
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/protoutil"
)

// this goes on the route extension
type RouteExtension struct {
	TransformationParameters Parameters `json:"transformation_parameters"`
}

type Parameters struct {
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
	Header    map[string]string `json:"header"`
	Path      string            `json:"path,omitempty"`
	Authority string            `json:"authority"`
	//TODO: support query params
	//TODO: support form params
}

// SPECIAL VARIABLE NAMES:
// these variables will be provided even if there is no parameter source
// path
// method
// scheme
// authority
// any field from the request body, assuming it's json

// this goes on the function spec
type FunctionSpec struct {
	Path   string            `json:"path"`
	Header map[string]string `json:"headers"`
	Body   string            `json:"body"`
	//TODO: support query template
	//TODO: support form template
}

func DecodeRouteExtension(generic *types.Struct) (RouteExtension, error) {
	var s RouteExtension
	err := protoutil.UnmarshalStruct(generic, &s)
	return s, err
}

func EncodeRouteExtension(spec RouteExtension) *types.Struct {
	v1Spec, err := protoutil.MarshalStruct(spec)
	if err != nil {
		panic(err)
	}
	return v1Spec
}

func DecodeFunctionSpec(generic *types.Struct) (FunctionSpec, error) {
	var s FunctionSpec
	err := protoutil.UnmarshalStruct(generic, &s)
	return s, err
}

func EncodeFunctionSpec(spec FunctionSpec) *types.Struct {
	v1Spec, err := protoutil.MarshalStruct(spec)
	if err != nil {
		panic(err)
	}
	return v1Spec
}
