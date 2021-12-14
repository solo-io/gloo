package graphql

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/graphql-go/graphql"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/go-utils/cliutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi/graphqlschematranslation/types"
)

type GetResolverParams struct {
	Operation      *types.Operation
	ArgsFromLink   map[string]interface{}
	PayloadName    string
	ResponseName   string
	Data           *types.PreprocessingData
	BaseUrl        string
	RequestOptions *RequestOptions
}

type RequestOptions struct {
}

// During the schema building, we stored the operation metadata (GetResolverParams) from the openapi spec in the
// `ResolveFn` thunk. We now get use the result of calling the thunk to build the resolver for each
// field.
func (t *OasToGqlTranslator) CreateResolverForField(parentTypeName string, field *graphql.Object, resolutions map[string]*Resolution) {
	// field does not return object type, do not dig into field
	if field == nil {
		return
	}
	for fName, f := range field.Fields() {
		if _, resolutionExists := resolutions[parentTypeName+"|"+fName]; resolutionExists {
			continue
		}
		if f.Resolve == nil {
			continue
		}
		getResolverParams, _ := f.Resolve(graphql.ResolveParams{})
		gRParams := getResolverParams.(GetResolverParams)
		resolvers := t.GetResolver(gRParams)
		for _, r := range resolvers {
			key := parentTypeName + "|" + fName
			resolutions[key] = &Resolution{
				Matcher: &QueryMatcher{
					Match: &QueryMatcher_FieldMatcher_{
						FieldMatcher: &QueryMatcher_FieldMatcher{
							Type:  parentTypeName,
							Field: fName,
						},
					},
				},
				Resolver: &Resolution_RestResolver{
					RestResolver: r,
				},
			}
		}
		if childField, ok := f.Type.(*graphql.Object); ok {
			t.CreateResolverForField(f.Type.String(), childField, resolutions)

		}
	}
}

func (t *OasToGqlTranslator) GetResolver(params GetResolverParams) []*RESTResolver {
	var resolvers []*RESTResolver
	baseUrl := params.BaseUrl
	if baseUrl == "" {
		var err error
		baseUrl, err = t.GetBaseUrlPath(params.Operation)
		if err != nil {
			t.handleWarningf(INVALID_SERVER_URL, "", Location{Operation: params.Operation},
				"Invalid server url in operation, skipping building resolver for operation")
			return nil
		}
	}
	resolver := &RESTResolver{
		UpstreamRef: &core.ResourceRef{
			Name:      t.Upstream.GetMetadata().GetName(),
			Namespace: t.Upstream.GetMetadata().GetNamespace(),
		},
	}
	operation := params.Operation
	// arg doesn't exit
	if operation != nil {

		// handle arguments provided by links
		for paramName, value := range params.ArgsFromLink {
			// value is an interface, cast it to string
			valueStr := fmt.Sprintf("%s", value)
			/**
			 * see if the link parameter contains constants that are appended to the link parameter
			 *
			 * e.g. instead of:
			 * $response.body#/employerId
			 *
			 * it could be:
			 * abc_{$response.body#/employerId}
			 */
			var valueProvider *ValueProvider
			if matches := regexp.MustCompile(`{([^}]*)}`).FindAllStringSubmatch(valueStr, 1); len(matches) > 0 {
				t.handleWarningf(LINK_PARAM_TEMPLATE_PRESENT, "", Location{Operation: operation},
					"Templating values with link parameters is not currently supported. Using only the first extraction and ignoring all else.")
				// replace link parameters with appropriate values

				valueStr = matches[0][1]
			}
			valueProvider = t.ResolveLinkParameter(valueStr)
			if valueProvider == nil {
				t.handleWarningf(LINK_UNSUPPORTED_EXTRACTION, "", Location{Operation: operation}, "Link %s uses an unsupported extraction %s. only $response.body extractions are currently supported",
					paramName, valueStr)
				continue
			}

			resolvers = append(resolvers, t.ExtractRequestDataFromParent(baseUrl, operation, paramName, valueProvider))
		}

		// Swallowing error here as it will never happen. URLJoin only returns error when baseUrl is invalid, which we confirmed to be valid in GetBaseUrlPath.
		extendedUrl := path.Join(baseUrl, operation.Path)
		resolverForArgs := ExtractRequestDataFromArgs(extendedUrl, operation, operation.Parameters, nil, params.ArgsFromLink)
		if resolverForArgs == nil {
			return resolvers
		}
		resolver.RequestTransform = resolverForArgs

		/**
		  Determine the possible payload type
		*/
		if params.Operation.PayloadDefinition != nil {
			if jsonVal := TraverseGraphqlSchema(params.Operation, params.Operation.PayloadDefinition.GraphQLInputObjectTypeName, params.Operation.PayloadDefinition.Schema, params.Data.SaneMap); jsonVal != nil {
				if resolver.RequestTransform == nil {
					resolver.RequestTransform = &RequestTemplate{}
				}
				resolver.RequestTransform.OutgoingBody = jsonVal
			}
		}
	}
	resolvers = append(resolvers, resolver)
	return resolvers
}

func TraverseGraphqlSchema(operation *types.Operation, inputTypeName string, schema *openapi.Schema, saneMap map[string]string) *JsonValue {

	inputSaneName := sanitizeString(inputTypeName, CaseStyle_camelCase)
	if schema.Type == "object" {
		ret := &JsonNode{
			KeyValues: []*JsonKeyValue{},
		}
		for propName, _ := range schema.Properties {
			p := sanitizeString(propName, CaseStyle_camelCase)
			ret.KeyValues = append(ret.KeyValues, &JsonKeyValue{
				Key: p,
				Value: &JsonValue{
					JsonVal: &JsonValue_ValueProvider{
						ValueProvider: &ValueProvider{
							Providers: map[string]*ValueProvider_Provider{"namedProvider": {
								Provider: &ValueProvider_Provider_GraphqlArg{
									GraphqlArg: &ValueProvider_GraphQLArgExtraction{
										ArgName:  inputSaneName,
										Required: cliutils.Contains(schema.Required, propName),
										Path: []*PathSegment{
											{
												Segment: &PathSegment_Key{
													Key: saneMap[p],
												}}}}}}}}}}})
		}
		return &JsonValue{
			JsonVal: &JsonValue_Node{
				Node: ret,
			},
		}
	} else if schema.Type == "array" {
		return &JsonValue{
			JsonVal: &JsonValue_ValueProvider{
				ValueProvider: &ValueProvider{
					Providers: map[string]*ValueProvider_Provider{"namedProvider": {
						Provider: &ValueProvider_Provider_GraphqlArg{
							GraphqlArg: &ValueProvider_GraphQLArgExtraction{
								ArgName:  inputSaneName,
								Required: operation.PayloadRequired,
								Path: []*PathSegment{{
									Segment: &PathSegment_All{
										All: true,
									}}}}}}}}},
		}
	} else {
		return nil
	}
}

func ExtractRequestDataFromArgs(path string, operation *types.Operation, parameters openapi.Parameters, provider *ValueProvider, argsFromLink map[string]interface{}) *RequestTemplate {
	method := operation.Method
	requestTemplate := &RequestTemplate{
		Headers: map[string]*ValueProvider{
			":path": {
				Providers: map[string]*ValueProvider_Provider{
					"namedProvider": {
						Provider: &ValueProvider_Provider_TypedProvider{
							TypedProvider: &ValueProvider_TypedValueProvider{
								ValProvider: &ValueProvider_TypedValueProvider_Value{
									Value: path,
								},
							},
						},
					},
				},
			},
			":method": {
				Providers: map[string]*ValueProvider_Provider{
					"namedProvider": {
						Provider: &ValueProvider_Provider_TypedProvider{
							TypedProvider: &ValueProvider_TypedValueProvider{
								ValProvider: &ValueProvider_TypedValueProvider_Value{
									Value: method,
								},
							},
						},
					},
				},
			},
		},
	}
	if len(parameters) == 0 {
		return requestTemplate
	}
	var parametersNotFromLink openapi.Parameters
	for _, param := range parameters {
		if _, ok := argsFromLink[param.Value.Name]; !ok {
			parametersNotFromLink = append(parametersNotFromLink, param)
		}
	}

	if len(parametersNotFromLink) == 0 {
		// all parameters were provided by link and we should not create a separate resolver
		return nil
	}
	for _, param := range parameters {

		p := param.Value
		sanitizedParamName := sanitizeString(p.Name, CaseStyle_camelCase)
		required := false
		if operation.PayloadDefinition != nil {
			required = cliutils.Contains(operation.PayloadDefinition.Required, sanitizedParamName)
		}
		if provider == nil {
			provider = &ValueProvider{
				Providers: map[string]*ValueProvider_Provider{
					"namedProvider": {
						Provider: &ValueProvider_Provider_GraphqlArg{
							GraphqlArg: &ValueProvider_GraphQLArgExtraction{
								ArgName:  sanitizedParamName,
								Required: required,
							},
						},
					},
				},
			}
		}
		switch p.In {
		case "path":
			// replace /pet/{petid} from openapispec to /pet/{} template string
			// for extraction
			// todo - support multiple name parameters here.
			pString := "{" + p.Name + "}"
			requestTemplate.Headers[":path"] = &ValueProvider{
				ProviderTemplate: strings.ReplaceAll(path, pString, "{namedProvider}"), // TODO(sai) fixme
				Providers:        provider.Providers,
			}
		case "query":
			if requestTemplate.QueryParams == nil {
				requestTemplate.QueryParams = map[string]*ValueProvider{}
			}
			requestTemplate.QueryParams[p.Name] = &ValueProvider{
				Providers: provider.Providers,
			}
		case "header":
			requestTemplate.Headers[p.Name] = &ValueProvider{
				Providers: provider.Providers,
			}
		case "cookie":
			requestTemplate.Headers["cookie"] = &ValueProvider{
				Providers: provider.Providers,
			}
		}
	}
	return requestTemplate
}
