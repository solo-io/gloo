package translation

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/grpc"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/mock"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/rest"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/directives"
	printer2 "github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/printer"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"
)

func CreateGraphQlApi(upstreams types.UpstreamList, graphqlapis types.GraphQLApiList, graphQLApi *v1beta1.GraphQLApi) (*v2.ExecutableSchema, error) {
	executableSchema, err := translateSchema(upstreams, graphqlapis, graphQLApi)
	if err != nil {
		return nil, err
	}
	executableSchema.LogRequestResponseInfo = graphQLApi.GetOptions().GetLogSensitiveInfo()
	return executableSchema, nil
}

func translateSchema(upstreams types.UpstreamList, graphqlapis types.GraphQLApiList, graphQLApi *v1beta1.GraphQLApi) (*v2.ExecutableSchema, error) {
	switch schema := graphQLApi.GetSchema().(type) {
	case *v1beta1.GraphQLApi_StitchedSchema:
		{
			return translateStitchedSchema(upstreams, graphqlapis, schema.StitchedSchema)
		}
	case *v1beta1.GraphQLApi_ExecutableSchema:
		{
			return translateExecutableSchema(upstreams, graphQLApi)
		}
	default:
		{
			return nil, eris.Errorf("unknown schema type %T", graphQLApi.GetSchema())
		}
	}
}

func translateExecutableSchema(upstreams types.UpstreamList, graphQLApi *v1beta1.GraphQLApi) (*v2.ExecutableSchema, error) {
	extensions, err := translateExtensions(graphQLApi)
	if err != nil {
		return nil, err
	}
	schemaStr := graphQLApi.GetExecutableSchema().GetSchemaDefinition()
	_, resolutions, processedSchema, err := processGraphqlSchema(upstreams, schemaStr, graphQLApi.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())
	if err != nil {
		return nil, err
	}

	return &v2.ExecutableSchema{
		Executor: &v2.Executor{
			Executor: &v2.Executor_Local_{
				Local: &v2.Executor_Local{
					Resolutions:         resolutions,
					EnableIntrospection: graphQLApi.GetExecutableSchema().GetExecutor().GetLocal().GetEnableIntrospection(),
				},
			},
		},
		SchemaDefinition: &v3.DataSource{
			Specifier: &v3.DataSource_InlineString{InlineString: printer2.PrettyPrintKubeString(processedSchema)},
		},
		Extensions: extensions,
	}, nil
}

func translateExtensions(api *v1beta1.GraphQLApi) (map[string]*any.Any, error) {
	extensions := map[string]*any.Any{}

	if reg := api.GetExecutableSchema().GetGrpcDescriptorRegistry(); reg != nil {

		grpcDescRegistry := &v2.GrpcDescriptorRegistry{
			ProtoDescriptors: &v3.DataSource{
				Specifier: nil, // filled in later
			},
		}

		switch regType := reg.DescriptorSet.(type) {
		case *v1beta1.GrpcDescriptorRegistry_ProtoDescriptor:
			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_Filename{
				Filename: reg.GetProtoDescriptor(),
			}
		case *v1beta1.GrpcDescriptorRegistry_ProtoDescriptorBin:
			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_InlineBytes{
				InlineBytes: reg.GetProtoDescriptorBin(),
			}
		default:
			return nil, eris.Errorf("unimplemented type %T for grpc resolver proto descriptor translation", regType)
		}
		extensions[grpc.GrpcRegistryExtensionName] = utils.MustMessageToAny(grpcDescRegistry)
	}

	if len(extensions) == 0 {
		return nil, nil
	}
	return extensions, nil
}

func processGraphqlSchema(upstreams types.UpstreamList, schema string, resolutions map[string]*v1beta1.Resolution) (*ast.Document, []*v2.Resolution, string, error) {
	doc, err := ParseGraphQLSchema(schema)
	if err != nil {
		return nil, nil, "", err
	}
	visitor := directives.NewGraphqlASTVisitor()
	var result []*v2.Resolution
	// Adds a directive visitor to the ast visitor which looks for `@resolve` directives in the schema. For each
	// resolve directive, it translates the resolver with the given name in the control plane resolutions map into
	// a data plane resolver.
	AddResolveDirectiveVisitor(visitor, upstreams, resolutions, &result)
	// Adds a directive visitor to the ast visitor which looks for `@cacheControl` directives in the schema. For each
	// cacheControl directive, it translates it to a data plane cache control and attaches it to the associated resolver.
	AddCacheControlDirectiveVisitor(visitor, &result)
	err = visitor.Visit(doc)
	if err != nil {
		return nil, nil, "", err
	}
	return doc, result, fmt.Sprintf("%s", printer.Print(doc)), nil
}

func ParseGraphQLSchema(schema string) (*ast.Document, error) {
	doc, err := parser.Parse(parser.ParseParams{Source: schema})
	if err != nil {
		return nil, eris.Wrapf(err, "unable to parse graphql schema %s", schema)
	}
	return doc, nil
}

func translateResolver(upstreams types.UpstreamList, resolver *v1beta1.Resolution) (*v3.TypedExtensionConfig, error) {
	switch r := resolver.Resolver.(type) {
	case *v1beta1.Resolution_RestResolver:
		return rest.TranslateRestResolver(upstreams, r.RestResolver)
	case *v1beta1.Resolution_GrpcResolver:
		return grpc.TranslateGrpcResolver(upstreams, r.GrpcResolver)
	case *v1beta1.Resolution_MockResolver:
		return mock.TranslateMockResolver(r.MockResolver)
	default:
		return nil, errors.Errorf("unimplemented resolver type: %T", r)
	}
}

func AddResolveDirectiveVisitor(visitor *directives.GraphqlASTVisitor, upstreams types.UpstreamList, resolutions map[string]*v1beta1.Resolution, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(directives.RESOLVER_DIRECTIVE, func(directiveVisitorParams directives.DirectiveVisitorParams) (bool, error) {
		// validate correct usage of the resolve directive
		resolveDirective := directives.NewResolveDirective()
		err := resolveDirective.Validate(directiveVisitorParams)
		if err != nil {
			return false, err
		}

		// check if the resolver referenced here even exists
		resolution := resolutions[resolveDirective.ResolverName]
		if resolution == nil {
			return false, directives.NewGraphqlSchemaError(resolveDirective.ResolverNameAstValue, "resolver %s is not defined",
				resolveDirective.ResolverName)
		}

		queryMatch := &v2.QueryMatcher{
			Match: &v2.QueryMatcher_FieldMatcher_{
				FieldMatcher: &v2.QueryMatcher_FieldMatcher{
					Type:  directiveVisitorParams.Type.Name.Value,
					Field: directiveVisitorParams.DirectiveField.Name.Value,
				},
			},
		}
		res, err := translateResolver(upstreams, resolution)
		if err != nil {
			return false, err
		}
		statsPrefix := resolveDirective.ResolverName
		if sp := resolution.StatPrefix; sp != nil {
			statsPrefix = sp.Value
		}
		statsPrefix = strings.TrimSuffix(statsPrefix, ".") + "."
		if result != nil {
			for i, resul := range *result {
				if proto.Equal(resul.Matcher, queryMatch) {
					(*result)[i].Resolver = res
					(*result)[i].StatPrefix = statsPrefix
					return true, nil
				}
			}
		}
		*result = append(*result, &v2.Resolution{
			Matcher:    queryMatch,
			Resolver:   res,
			StatPrefix: statsPrefix,
		})
		return true, nil
	})
}

func AddCacheControlDirectiveVisitor(visitor *directives.GraphqlASTVisitor, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(directives.CACHE_CONTROL_DIRECTIVE, func(directiveVisitorParams directives.DirectiveVisitorParams) (bool, error) {
		// validate correct usage of the cacheControl directive
		cacheControlDirective := directives.NewCacheControlDirective()
		_, err := cacheControlDirective.Validate(directiveVisitorParams)
		if err != nil {
			return false, err
		}

		cacheControl := cacheControlDirective.CacheControl

		var queryMatchList []*v2.QueryMatcher
		for _, df := range directiveVisitorParams.DirectiveFields {
			// this is a type-level directive. we will set on all fields
			queryMatch := &v2.QueryMatcher{
				Match: &v2.QueryMatcher_FieldMatcher_{
					FieldMatcher: &v2.QueryMatcher_FieldMatcher{
						Type:  directiveVisitorParams.Type.Name.Value,
						Field: df.Name.Value,
					},
				},
			}
			queryMatchList = append(queryMatchList, queryMatch)
		}

		if directiveVisitorParams.DirectiveField != nil {
			// this is a field-directive
			queryMatch := &v2.QueryMatcher{
				Match: &v2.QueryMatcher_FieldMatcher_{
					FieldMatcher: &v2.QueryMatcher_FieldMatcher{
						Type:  directiveVisitorParams.Type.Name.Value,
						Field: directiveVisitorParams.DirectiveField.Name.Value,
					},
				},
			}
			queryMatchList = append(queryMatchList, queryMatch)
		}

		if len(queryMatchList) == 0 {
			return false, eris.Errorf("logic error: no query match generated but `@cacheControl` directive was found")
		}

		for _, queryMatch := range queryMatchList {
			found := false
			if result != nil {
				for i, res := range *result {
					if proto.Equal(res.Matcher, queryMatch) {
						(*result)[i].CacheControl = cacheControl
						found = true
						break
					}
				}
			}
			if !found {
				*result = append(*result, &v2.Resolution{
					Matcher:      queryMatch,
					CacheControl: cacheControl,
				})
			}
		}
		return true, nil
	})
}
