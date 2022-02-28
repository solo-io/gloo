package graphql

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/dot_notation"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	ARBITRARY_PROVIDER_NAME              = "ARBITRARY_PROVIDER_NAME"
	restResolverTypedExtensionConfigName = "io.solo.graphql.resolver.rest"

	// allowed setter keywords; for now just $body but may add $headers in the future
	BODY_SETTER = "body"
)

func translateRestResolver(params plugins.RouteParams, r *RESTResolver) (*v3.TypedExtensionConfig, error) {
	requestTransform, err := translateRequestTransform(r.Request)
	if err != nil {
		return nil, err
	}
	responseTransform, err := translateResponseTransform(r.Response)
	if err != nil {
		return nil, err
	}
	us, err := params.Snapshot.Upstreams.Find(r.UpstreamRef.GetNamespace(), r.UpstreamRef.GetName())
	if err != nil {
		return nil, eris.Wrapf(err, "unable to find upstream `%s` in namespace `%s` to resolve schema", r.UpstreamRef.GetName(), r.UpstreamRef.GetNamespace())
	}
	restResolver := &v2.RESTResolver{
		ServerUri: &v3.HttpUri{
			Uri: "ignored", // ignored by graphql filter
			HttpUpstreamType: &v3.HttpUri_Cluster{
				Cluster: translator.UpstreamToClusterName(us.GetMetadata().Ref()),
			},
			Timeout: durationpb.New(1 * time.Second),
		},
		RequestTransform:      requestTransform,
		PreExecutionTransform: responseTransform,
		SpanName:              r.SpanName,
	}
	return &v3.TypedExtensionConfig{
		Name:        restResolverTypedExtensionConfigName,
		TypedConfig: utils.MustMessageToAny(restResolver),
	}, nil
}

func translateResponseTransform(transform *ResponseTemplate) (*v2.ResponseTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	resultRoot, err := dot_notation.DotNotationToPathSegments(transform.ResultRoot)
	if err != nil {
		return nil, eris.Wrapf(err, "error translating result root path %s", transform.ResultRoot)
	}
	setters := map[string]*v2.TemplatedPath{}
	for k, path := range transform.Setters {
		templatedPath, err := TranslateSetter(path)
		if err != nil {
			return nil, eris.Wrapf(err, "error, unable to translate setter string to templated path")
		}
		setters[k] = templatedPath
	}
	result := &v2.ResponseTemplate{
		ResultRoot: resultRoot,
		Setters:    setters,
	}
	return result, nil
}

func translateRequestTransform(transform *RequestTemplate) (*v2.RequestTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	headersMap, err := TranslateStringValueProviderMap(transform.GetHeaders())
	if err != nil {
		return nil, err
	}
	queryParamsMap, err := TranslateStringValueProviderMap(transform.GetQueryParams())
	if err != nil {
		return nil, err
	}
	jv, err := TranslateJsonValue(transform.GetBody())
	if err != nil {
		return nil, err
	}
	rt := &v2.RequestTemplate{
		Headers:      headersMap,
		QueryParams:  queryParamsMap,
		OutgoingBody: jv, // filled in later
	}

	return rt, nil
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
	providerTemplateRe = regexp.MustCompile(`{\$([[a-zA-Z0-9.\[\]]+)}`)
)

func translateValueProvider(vpString string) (*v2.ValueProvider, error) {
	ret := &v2.ValueProvider{
		Providers:        map[string]*v2.ValueProvider_Provider{},
		ProviderTemplate: vpString,
	}
	subMatches := providerTemplateRe.FindAllStringSubmatch(vpString, -1)
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

func TranslateSetter(templatedPathString string) (*v2.TemplatedPath, error) {
	ret := &v2.TemplatedPath{
		PathTemplate: templatedPathString,
		NamedPaths:   map[string]*v2.Path{},
	}
	subMatches := providerTemplateRe.FindAllStringSubmatch(templatedPathString, -1)
	if len(subMatches) < 1 {
		return nil, eris.Errorf("malformed setter %s, needs to match template string regex %s", templatedPathString, providerTemplateRe.String())
	}

	for _, matches := range subMatches {
		match := matches[1]
		pathSegment, err := dot_notation.DotNotationToPathSegments(match)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to parse graphql extraction string %s", templatedPathString)
		}
		if len(pathSegment) <= 1 {
			return nil, eris.New("invalid extraction string " + templatedPathString)
		}

		if pathSegment[0].GetKey() != BODY_SETTER {
			return nil, eris.New("currently only support for grabbing from the request body with {$body.}")
		}
		// remove $body from path segments, this was a special keyword
		// we didn't remove it from `match` (NamedPaths` key) so we could leave room in our api for new keywords, e.g. $headers
		pathSegment = pathSegment[1:]

		tp := &v2.Path{
			Segments: pathSegment,
		}
		matchSanitized := strings.ReplaceAll(match, ".", "")
		ret.NamedPaths[matchSanitized] = tp
		ret.PathTemplate = strings.ReplaceAll(ret.PathTemplate, matches[0], "{"+matchSanitized+"}")
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
