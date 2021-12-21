package graphql

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes/any"
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
	resolutions, err := translateResolutions(params, schema.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())
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
				Specifier: &v3.DataSource_InlineString{InlineString: schema.GetExecutableSchema().GetSchemaDefinition()},
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

func translateResolutions(params plugins.RouteParams, resolvers []*v1alpha1.Resolution) ([]*v2.Resolution, error) {
	if len(resolvers) == 0 {
		return nil, nil
	}

	var converted []*v2.Resolution
	for _, r := range resolvers {
		matcher, err := translateQueryMatcher(r.Matcher)
		if err != nil {
			return nil, err
		}
		res, err := translateResolver(params, r)
		if err != nil {
			return nil, err
		}
		resolver := &v2.Resolution{
			Matcher:  matcher,
			Resolver: res,
		}
		converted = append(converted, resolver)
	}

	return converted, nil
}

func translateQueryMatcher(matcher *v1alpha1.QueryMatcher) (*v2.QueryMatcher, error) {
	qm := &v2.QueryMatcher{}
	switch m := matcher.Match.(type) {
	case *v1alpha1.QueryMatcher_FieldMatcher_:
		qm.Match = &v2.QueryMatcher_FieldMatcher_{
			FieldMatcher: &v2.QueryMatcher_FieldMatcher{
				Type:  m.FieldMatcher.Type,
				Field: m.FieldMatcher.Field,
			},
		}
	default:
		return nil, errors.Errorf("unimplemented matcher type: %T", m)
	}
	return qm, nil
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
