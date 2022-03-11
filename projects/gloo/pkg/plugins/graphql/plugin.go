package graphql

import (
	"strings"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	directive_utils "github.com/solo-io/solo-projects/projects/gloo/utils/graphql/directives"
)

var (
	_ plugins.Plugin           = new(plugin)
	_ plugins.RoutePlugin      = new(plugin)
	_ plugins.HttpFilterPlugin = new(plugin)
)

const (
	FilterName    = "io.solo.filters.http.graphql"
	ExtensionName = "graphql"
)

var (
	// This filter must be last as it is used to replace the router filter
	FilterStage = plugins.BeforeStage(plugins.RouteStage)
)

type plugin struct {
	filterNeeded bool
}

// NewPlugin creates the basic graphql plugin structure.
func NewPlugin() *plugin {
	return &plugin{}
}

// Name returns the ExtensionName for overwriting purposes.
func (p *plugin) Name() string {
	return ExtensionName
}

// Init resets the plugin and creates the maps within the structure.
func (p *plugin) Init(params plugins.InitParams) error {
	p.filterNeeded = !params.Settings.GetGloo().GetRemoveUnusedFilters().GetValue()
	return nil
}

// HttpFilters sets up the filters for envoy if it is needed.
func (p *plugin) HttpFilters(_ plugins.Params, _ *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	if !p.filterNeeded {
		return filters, nil
	}

	emptyConf := &v2.GraphQLConfig{}
	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, emptyConf, FilterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

// ProcessRoute aplying any needed configurations related to grapql.
// If any configs are found then mark us needing this filter in our chain.
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	gqlRef := in.GetGraphqlSchemaRef()
	if gqlRef == nil {
		return nil
	}

	gql, err := params.Snapshot.GraphqlSchemas.Find(gqlRef.GetNamespace(), gqlRef.GetName())
	if err != nil {
		ret := ""
		for _, schema := range params.Snapshot.GraphqlSchemas {
			ret += " " + schema.Metadata.Name
		}
		return eris.Wrapf(err, "unable to find graphql schema custom resource `%s` in namespace `%s`, list of all graphqlschemas found: %s", gqlRef.GetName(), gqlRef.GetNamespace(), ret)
	}

	p.filterNeeded = true // Set here as user is at least attempting to use graphql at this point so might as well place it in the filterchain.
	routeConf, err := translateGraphQlSchemaToRouteConf(params, in, gql)

	if err != nil {
		return eris.Wrapf(err, "unable to translate graphql schema control plane config to data plane config")
	}

	return pluginutils.SetRoutePerFilterConfig(out, FilterName, routeConf)
}

func translateGraphQlSchemaToRouteConf(params plugins.RouteParams, in *v1.Route, schema *v1alpha1.GraphQLSchema) (*v2.GraphQLRouteConfig, error) {
	schemaStr := schema.GetExecutableSchema().GetSchemaDefinition()
	_, resolutions, processedSchema, err := processGraphqlSchema(params, schemaStr, schema.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())
	if err != nil {
		return nil, err
	}
	extensions, err := translateExtensions(schema)
	if err != nil {
		return nil, err
	}
	statsPrefix := in.GetGraphqlSchemaRef().Key()
	if sp := schema.GetStatPrefix().GetValue(); sp != "" {
		statsPrefix = sp
	}
	statsPrefix = strings.TrimSuffix(statsPrefix, ".") + "."
	cacheConf := &v2.PersistedQueryCacheConfig{}
	if cc := schema.GetPersistedQueryCacheConfig(); cc != nil {
		cacheConf.CacheSize = cc.CacheSize
	}
	return &v2.GraphQLRouteConfig{
		ExecutableSchema: &v2.ExecutableSchema{
			Executor: &v2.Executor{
				Executor: &v2.Executor_Local_{
					Local: &v2.Executor_Local{
						Resolutions:         resolutions,
						EnableIntrospection: schema.GetExecutableSchema().GetExecutor().GetLocal().GetEnableIntrospection(),
					},
				},
			},
			SchemaDefinition: &v3.DataSource{
				Specifier: &v3.DataSource_InlineString{InlineString: processedSchema},
			},
			Extensions: extensions,
		},
		StatPrefix:                statsPrefix,
		PersistedQueryCacheConfig: cacheConf,
		AllowedQueryHashes:        schema.GetAllowedQueryHashes(),
	}, nil
}

func translateExtensions(schema *v1alpha1.GraphQLSchema) (map[string]*any.Any, error) {
	extensions := map[string]*any.Any{}

	if reg := schema.GetExecutableSchema().GetGrpcDescriptorRegistry(); reg != nil {

		grpcDescRegistry := &v2.GrpcDescriptorRegistry{
			ProtoDescriptors: &v3.DataSource{
				Specifier: nil, // filled in later
			},
		}

		switch regType := reg.DescriptorSet.(type) {
		case *v1alpha1.GrpcDescriptorRegistry_ProtoDescriptor:
			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_Filename{
				Filename: reg.GetProtoDescriptor(),
			}
		case *v1alpha1.GrpcDescriptorRegistry_ProtoDescriptorBin:
			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_InlineBytes{
				InlineBytes: reg.GetProtoDescriptorBin(),
			}
		default:
			return nil, eris.Errorf("unimplemented type %T for grpc resolver proto descriptor translation", regType)
		}
		extensions[grpcRegistryExtensionName] = utils.MustMessageToAny(grpcDescRegistry)
	}

	if len(extensions) == 0 {
		return nil, nil
	}
	return extensions, nil
}

func processGraphqlSchema(params plugins.RouteParams, schema string, resolutions map[string]*v1alpha1.Resolution) (*ast.Document, []*v2.Resolution, string, error) {
	doc, err := parser.Parse(parser.ParseParams{Source: schema})
	if err != nil {
		return nil, nil, "", eris.Wrapf(err, "unable to parse graphql schema %s", schema)
	}
	visitor := directive_utils.NewGraphqlASTVisitor()
	var result []*v2.Resolution
	// Adds a directive visitor to the ast visitor which looks for `@resolve` directives in the schema. For each
	// resolve directive, it translates the resolver with the given name in the control plane resolutions map into
	// a data plane resolver.
	addResolveDirectiveVisitor(visitor, params, resolutions, &result)
	// Adds a directive visitor to the ast visitor which looks for `@cacheControl` directives in the schema. For each
	// cacheControl directive, it translates it to a data plane cache control and attaches it to the associated resolver.
	addCacheControlDirectiveVisitor(visitor, &result)
	err = visitor.Visit(doc)
	if err != nil {
		return nil, nil, "", err
	}
	return doc, result, schema, nil
}

func translateResolver(params plugins.RouteParams, resolver *v1alpha1.Resolution) (*v3.TypedExtensionConfig, error) {
	switch r := resolver.Resolver.(type) {
	case *v1alpha1.Resolution_RestResolver:
		return translateRestResolver(params, r.RestResolver)
	case *v1alpha1.Resolution_GrpcResolver:
		return translateGrpcResolver(params, r.GrpcResolver)
	default:
		return nil, errors.Errorf("unimplemented resolver type: %T", r)
	}
}

func addResolveDirectiveVisitor(visitor *directive_utils.GraphqlASTVisitor, params plugins.RouteParams, resolutions map[string]*v1alpha1.Resolution, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(directive_utils.RESOLVER_DIRECTIVE, func(directiveVisitorParams directive_utils.DirectiveVisitorParams) error {
		// validate correct usage of the resolve directive
		resolveDirective := directive_utils.NewResolveDirective()
		err := resolveDirective.Validate(directiveVisitorParams)
		if err != nil {
			return err
		}

		// check if the resolver referenced here even exists
		resolution := resolutions[resolveDirective.ResolverName]
		if resolution == nil {
			return directive_utils.NewGraphqlSchemaError(resolveDirective.ResolverNameAstValue, "resolver %s is not defined",
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
		res, err := translateResolver(params, resolution)
		if err != nil {
			return err
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
					return nil
				}
			}
		}
		*result = append(*result, &v2.Resolution{
			Matcher:    queryMatch,
			Resolver:   res,
			StatPrefix: statsPrefix,
		})
		return nil
	})
}

func addCacheControlDirectiveVisitor(visitor *directive_utils.GraphqlASTVisitor, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(directive_utils.CACHE_CONTROL_DIRECTIVE, func(directiveVisitorParams directive_utils.DirectiveVisitorParams) error {
		// validate correct usage of the cacheControl directive
		cacheControlDirective := directive_utils.NewCacheControlDirective()
		err := cacheControlDirective.Validate(directiveVisitorParams)
		if err != nil {
			return err
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
			return eris.Errorf("logic error: no query match generated but `@cacheControl` directive was found")
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
		return nil
	})
}
