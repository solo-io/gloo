package graphql

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"

	structpb "github.com/golang/protobuf/ptypes/struct"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/graphql-go/graphql"
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
			var valueProvider string
			if matches := regexp.MustCompile(`{([^}]*)}`).FindAllStringSubmatch(valueStr, 1); len(matches) > 0 {
				t.handleWarningf(LINK_PARAM_TEMPLATE_PRESENT, "", Location{Operation: operation},
					"Templating values with link parameters is not currently supported. Using only the first extraction and ignoring all else.")
				// replace link parameters with appropriate values

				valueStr = matches[0][1]
			}
			valueProvider = t.ResolveLinkParameter(valueStr)
			if valueProvider == "" {
				t.handleWarningf(LINK_UNSUPPORTED_EXTRACTION, "", Location{Operation: operation}, "Link %s uses an unsupported extraction %s. only $response.body extractions are currently supported",
					paramName, valueStr)
				continue
			}

			resolvers = append(resolvers, t.ExtractRequestDataFromParent(baseUrl, operation, paramName, valueProvider))
		}

		// Swallowing error here as it will never happen. URLJoin only returns error when baseUrl is invalid, which we confirmed to be valid in GetBaseUrlPath.
		extendedUrl := path.Join(baseUrl, operation.Path)
		resolverForArgs := ExtractRequestDataFrom(extendedUrl, operation, operation.Parameters, "", params.ArgsFromLink)
		if resolverForArgs == nil {
			return resolvers
		}
		resolver.Request = resolverForArgs

		/**
		  Determine the possible payload type
		*/
		if params.Operation.PayloadDefinition != nil {
			if jsonVal := TraverseGraphqlSchema(params.Operation, params.Operation.PayloadDefinition.GraphQLInputObjectTypeName, params.Operation.PayloadDefinition.Schema, params.Data.SaneMap); jsonVal != nil {
				if resolver.Request == nil {
					resolver.Request = &RequestTemplate{}
				}
				resolver.Request.Body = jsonVal
			}
		}
	}
	resolvers = append(resolvers, resolver)
	return resolvers
}

func TraverseGraphqlSchema(operation *types.Operation, inputTypeName string, schema *openapi.Schema, saneMap map[string]string) *structpb.Value {

	inputSaneName := sanitizeString(inputTypeName, CaseStyle_camelCase)
	if schema.Type == "object" {
		ret := &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		}
		for propName, _ := range schema.Properties {
			p := sanitizeString(propName, CaseStyle_camelCase)
			ret.Fields[p] = &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: fmt.Sprintf("{$args.%s.%s}", inputSaneName, saneMap[p]),
				},
			}
		}
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: ret,
			},
		}
	} else if schema.Type == "array" {
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("{$args.%s}", inputSaneName)}}
	} else {
		return nil
	}
}

func ExtractRequestDataFrom(path string, operation *types.Operation, parameters openapi.Parameters, providerString string, argsFromLink map[string]interface{}) *RequestTemplate {
	method := operation.Method
	requestTemplate := &RequestTemplate{
		Headers: map[string]string{
			":path":   path,
			":method": method,
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
		if providerString == "" {
			providerString = fmt.Sprintf("{$args.%s}", sanitizedParamName)
		}
		switch p.In {
		case "path":
			// replace /pet/{petid} from openapispec to /pet/{$args.ARG_NAME} template string
			// for extraction
			// todo - support multiple name parameters here.
			pString := "{" + p.Name + "}"
			requestTemplate.Headers[":path"] = strings.ReplaceAll(path, pString, providerString)

		case "query":
			if requestTemplate.QueryParams == nil {
				requestTemplate.QueryParams = map[string]string{}
			}
			requestTemplate.QueryParams[p.Name] = providerString
		case "header":
			requestTemplate.Headers[p.Name] = providerString
		case "cookie":
			requestTemplate.Headers["cookie"] = providerString
		}
	}
	return requestTemplate
}
