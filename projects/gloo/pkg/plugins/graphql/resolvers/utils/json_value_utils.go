package utils

import (
	"fmt"
	"regexp"
	"strings"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/rotisserie/eris"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/dot_notation"
)

func GetTypedValueProvider(value string, t v2.ValueProvider_TypedValueProvider_Type) *v2.JsonValue {
	return &v2.JsonValue{
		JsonVal: &v2.JsonValue_ValueProvider{
			ValueProvider: &v2.ValueProvider{
				Providers: map[string]*v2.ValueProvider_Provider{
					ARBITRARY_PROVIDER_NAME: {
						Provider: &v2.ValueProvider_Provider_TypedProvider{
							TypedProvider: &v2.ValueProvider_TypedValueProvider{
								Type: t,
								ValProvider: &v2.ValueProvider_TypedValueProvider_Value{
									Value: value,
								},
							},
						},
					},
				},
			},
		},
	}
}

func TranslateJsonValue(body *structpb.Value) (*v2.JsonValue, error) {
	if body == nil {
		return nil, nil
	}
	switch b := body.GetKind().(type) {
	case *structpb.Value_StructValue:
		structNode, err := translateJsonStruct(b.StructValue)
		if err != nil {
			return nil, err
		}
		return &v2.JsonValue{
			JsonVal: &v2.JsonValue_Node{
				Node: structNode,
			},
		}, nil
	case *structpb.Value_ListValue:
		structpbList, err := translateJsonList(b.ListValue)
		if err != nil {
			return nil, err
		}
		return &v2.JsonValue{
			JsonVal: &v2.JsonValue_List{
				List: &v2.JsonValueList{
					Values: structpbList,
				},
			},
		}, nil
	case *structpb.Value_BoolValue:
		value := "false"
		if b.BoolValue {
			value = "true"
		}
		return GetTypedValueProvider(value, v2.ValueProvider_TypedValueProvider_BOOLEAN), nil
	case *structpb.Value_NumberValue:
		return GetTypedValueProvider(fmt.Sprintf("%f", b.NumberValue), v2.ValueProvider_TypedValueProvider_FLOAT), nil
	case *structpb.Value_StringValue:
		vp, err := translateValueProvider(b.StringValue)
		if err != nil {
			return nil, err
		}
		return &v2.JsonValue{
			JsonVal: &v2.JsonValue_ValueProvider{
				ValueProvider: vp,
			},
		}, nil

	}
	return nil, nil
}

func translateJsonList(value *structpb.ListValue) ([]*v2.JsonValue, error) {
	var convertedList []*v2.JsonValue
	for _, v := range value.GetValues() {
		newVal, err := TranslateJsonValue(v)
		if err != nil {
			return nil, err
		}
		convertedList = append(convertedList, newVal)
	}
	return convertedList, nil
}

func translateJsonStruct(p *structpb.Struct) (*v2.JsonNode, error) {
	var convertedKvs []*v2.JsonKeyValue
	for k, v := range p.GetFields() {
		newVal, err := TranslateJsonValue(v)
		if err != nil {
			return nil, err
		}
		newkv := &v2.JsonKeyValue{
			Key: k,
			Value: &v2.JsonValue{
				JsonVal: newVal.JsonVal,
			},
		}
		convertedKvs = append(convertedKvs, newkv)

	}
	return &v2.JsonNode{KeyValues: convertedKvs}, nil
}

func TranslateStringValueProviderMap(headers map[string]string) (map[string]*v2.ValueProvider, error) {
	if len(headers) == 0 {
		return nil, nil
	}
	converted := map[string]*v2.ValueProvider{}
	for header, providerString := range headers {
		vp, err := translateValueProvider(providerString)
		if err != nil {
			return nil, err
		}
		converted[header] = vp
	}
	return converted, nil
}

var (
	ProviderTemplateRegex   = regexp.MustCompile(`{\$([[a-zA-Z0-9\[\]*.\-_]+)}`)
	ARBITRARY_PROVIDER_NAME = "ARBITRARY_PROVIDER_NAME"
)

func translateValueProvider(vpString string) (*v2.ValueProvider, error) {
	ret := &v2.ValueProvider{
		Providers:        map[string]*v2.ValueProvider_Provider{},
		ProviderTemplate: vpString,
	}
	subMatches := ProviderTemplateRegex.FindAllStringSubmatch(vpString, -1)
	if len(subMatches) < 1 {
		// no templated value providers are needed here
		ret.Providers[ARBITRARY_PROVIDER_NAME] = &v2.ValueProvider_Provider{
			Provider: &v2.ValueProvider_Provider_TypedProvider{
				TypedProvider: &v2.ValueProvider_TypedValueProvider{
					ValProvider: &v2.ValueProvider_TypedValueProvider_Value{
						Value: vpString,
					},
				},
			},
		}
		return ret, nil
	}

	for _, matches := range subMatches {
		match := matches[1]
		pathSegment, err := dot_notation.DotNotationToPathSegments(match)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to parse graphql extraction string %s", vpString)
		}
		if len(pathSegment) <= 1 {
			return nil, eris.New("invalid extraction string " + vpString)
		}
		vp, err := getExtraction(pathSegment, NewExtractionCb(vpString))
		if err != nil {
			return nil, err
		}
		matchSanitized := strings.ReplaceAll(match, ".", "")
		ret.Providers[matchSanitized] = vp
		ret.ProviderTemplate = strings.ReplaceAll(ret.ProviderTemplate, matches[0], "{"+matchSanitized+"}")
	}
	return ret, nil
}

type ExtractionCb func(string, ...interface{}) error

func NewExtractionCb(extractionString string) ExtractionCb {
	return func(fmtString string, args ...interface{}) error {
		return eris.Wrapf(eris.New(fmt.Sprintf("Error in extraction %s", extractionString)), fmtString, args...)
	}
}

const (
	HEADERS = "headers"
	ARGS    = "args"
	PARENT  = "parent"
)

func getExtraction(pathSegment []*v2.PathSegment, errorCb ExtractionCb) (*v2.ValueProvider_Provider, error) {
	// here we assume that pathSegment has a length of at least 2 elements
	switch pathSegment[0].GetKey() {
	case HEADERS:
		header := pathSegment[1].GetKey()
		if len(pathSegment) > 2 || header == "" {
			return nil, errorCb("you may only specify one header to extract and it must be a string")
		}
		return &v2.ValueProvider_Provider{
			Provider: &v2.ValueProvider_Provider_TypedProvider{
				TypedProvider: &v2.ValueProvider_TypedValueProvider{
					ValProvider: &v2.ValueProvider_TypedValueProvider_Header{
						Header: header,
					},
				},
			},
		}, nil
	case ARGS:
		argName := pathSegment[1].GetKey()
		if argName == "" {
			return nil, errorCb("second element in path must be the arg name")
		}
		return &v2.ValueProvider_Provider{
			Provider: &v2.ValueProvider_Provider_GraphqlArg{
				GraphqlArg: &v2.ValueProvider_GraphQLArgExtraction{
					ArgName: argName,
					Path:    pathSegment[2:],
					//todo(sai) Required:
				},
			},
		}, nil
	case PARENT:
		return &v2.ValueProvider_Provider{
			Provider: &v2.ValueProvider_Provider_GraphqlParent{
				GraphqlParent: &v2.ValueProvider_GraphQLParentExtraction{
					Path: pathSegment[1:],
				},
			},
		}, nil
	default:
		return nil, errorCb("Invalid first key, must be one of '%s'", strings.Join([]string{HEADERS, ARGS, PARENT}, ","))
	}
}
