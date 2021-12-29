package graphql

import (
	"fmt"

	"github.com/graphql-go/graphql/language/kinds"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/graphql-go/graphql/gqlerrors"
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
)

const (
	FilterName    = "io.solo.filters.http.graphql"
	ExtensionName = "graphql"
)

var (
	_ plugins.Plugin           = new(Plugin)
	_ plugins.RoutePlugin      = new(Plugin)
	_ plugins.HttpFilterPlugin = new(Plugin)
	_ plugins.Upgradable       = new(Plugin)

	// This filter must be last as it is used to replace the router filter
	FilterStage = plugins.BeforeStage(plugins.RouteStage)
)

type Plugin struct {
}

var _ plugins.Plugin = new(Plugin)

func NewPlugin() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *Plugin) PluginName() string {
	return ExtensionName
}

func (p *Plugin) IsUpgrade() bool {
	return true
}

func (p *Plugin) HttpFilters(_ plugins.Params, _ *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	emptyConf := &v2.GraphQLConfig{}
	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, emptyConf, FilterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
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

	routeConf, err := translateGraphQlSchemaToRouteConf(params, gql)
	if err != nil {
		return eris.Wrapf(err, "unable to translate graphql schema control plane config to data plane config")
	}
	return pluginutils.SetRoutePerFilterConfig(out, FilterName, routeConf)
}

func translateGraphQlSchemaToRouteConf(params plugins.RouteParams, schema *v1alpha1.GraphQLSchema) (*v2.GraphQLRouteConfig, error) {
	schemaStr := schema.GetExecutableSchema().GetSchemaDefinition()
	_, resolutions, processedSchema, err := ProcessGraphqlSchema(params, schemaStr, schema.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())
	if err != nil {
		return nil, err
	}
	extensions, err := translateExtensions(schema)
	if err != nil {
		return nil, err
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

type Locatable interface {
	GetLoc() *ast.Location
}

func NewGraphqlSchemaError(l Locatable, description string, args ...interface{}) error {
	desc := fmt.Sprintf(description, args...)
	return gqlerrors.NewSyntaxError(l.GetLoc().Source, l.GetLoc().Start, desc)
}

func ProcessGraphqlSchema(params plugins.RouteParams, schema string, resolutions map[string]*v1alpha1.Resolution) (*ast.Document, []*v2.Resolution, string, error) {
	doc, err := parser.Parse(parser.ParseParams{Source: schema})
	if err != nil {
		return nil, nil, "", eris.Wrapf(err, "unable to parse graphql schema %s", schema)
	}
	visitor := NewGraphqlASTVisitor()
	var result []*v2.Resolution
	// Adds a directive visitor to the ast visitor which looks for `@resolve` direcitves and adds them to the resolution
	// map
	AddResolveDirectiveVisitor(visitor, params, resolutions, &result)
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

const (
	RESOLVER_DIRECTIVE     = "resolve"
	RESOLVER_NAME_ARGUMENT = "name"
)

func AddResolveDirectiveVisitor(visitor *GraphqlASTVisitor, params plugins.RouteParams, resolutions map[string]*v1alpha1.Resolution, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(RESOLVER_DIRECTIVE, func(directiveVisitorParams DirectiveVisitorParams) error {
		arguments := map[string]ast.Value{}
		directive := directiveVisitorParams.Directive
		for _, argument := range directive.Arguments {
			arguments[argument.Name.Value] = argument.Value
		}
		resolver_name, ok := arguments[RESOLVER_NAME_ARGUMENT]
		if !ok {
			return NewGraphqlSchemaError(directive, `the "resolve" directive must have a "name" argument to reference a resolver`)
		}
		if resolver_name.GetKind() != kinds.StringValue {
			return NewGraphqlSchemaError(resolver_name, `"name" argument must be a string value`)
		}
		name := resolver_name.GetValue().(string)
		// check if the resolver referenced here even exists
		resolution := resolutions[name]
		if resolution == nil {
			return NewGraphqlSchemaError(resolver_name, "resolver %s is not defined", name)
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
		*result = append(*result, &v2.Resolution{
			Matcher:  queryMatch,
			Resolver: res,
		})
		return nil
	})
}
