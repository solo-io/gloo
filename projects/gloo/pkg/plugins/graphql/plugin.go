package graphql

import (
	"time"

	"github.com/pkg/errors"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"google.golang.org/protobuf/types/known/durationpb"
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
	resolutions, err := translateResolutions(params, schema.Resolutions)
	if err != nil {
		return nil, err
	}
	return &v2.GraphQLRouteConfig{
		Schema: &v3.DataSource{
			Specifier: &v3.DataSource_InlineString{InlineString: schema.GetSchema()},
		},
		EnableIntrospection: schema.EnableIntrospection,
		Resolutions:         resolutions,
	}, nil
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
	typedExtensionConf := &v3.TypedExtensionConfig{}
	switch r := resolver.Resolver.(type) {
	case *v1alpha1.Resolution_RestResolver:
		t, err := translateRequestTransform(r.RestResolver.RequestTransform)
		if err != nil {
			return nil, err
		}
		us, err := params.Snapshot.Upstreams.Find(r.RestResolver.UpstreamRef.GetNamespace(), r.RestResolver.UpstreamRef.GetName())
		if err != nil {
			return nil, eris.Wrapf(err, "unable to find upstream `%s` in namespace `%s` to resolve schema", r.RestResolver.UpstreamRef, r.RestResolver.UpstreamRef)
		}
		restResolver := &v2.RESTResolver{
			ServerUri: &v3.HttpUri{
				Uri: "ignored", // ignored by graphql filter
				HttpUpstreamType: &v3.HttpUri_Cluster{
					Cluster: translator.UpstreamToClusterName(us.GetMetadata().Ref()),
				},
				Timeout: durationpb.New(1 * time.Second),
			},
			RequestTransform: t,
			SpanName:         r.RestResolver.SpanName,
		}
		typedExtensionConf = &v3.TypedExtensionConfig{
			Name:        "io.solo.graphql.resolver.rest",
			TypedConfig: utils.MustMessageToAny(restResolver),
		}
	default:
		return nil, errors.Errorf("unimplemented resolver type: %T", r)
	}
	return typedExtensionConf, nil
}

func translateRequestTransform(transform *v1alpha1.RequestTemplate) (*v2.RequestTemplate, error) {
	if transform == nil {
		return nil, nil
	}
	headersMap, err := translateStringValueProviderMap(transform.GetHeaders())
	if err != nil {
		return nil, err
	}
	queryParamsMap, err := translateStringValueProviderMap(transform.GetQueryParams())
	if err != nil {
		return nil, err
	}
	rt := &v2.RequestTemplate{
		Headers:      headersMap,
		QueryParams:  queryParamsMap,
		OutgoingBody: nil, // filled in later
	}

	jv, err := translateJsonValue(transform.OutgoingBody)
	if err != nil {
		return nil, err
	}
	rt.OutgoingBody = jv
	return rt, nil
}

func translateJsonNode(jn *v1alpha1.JsonNode) (*v2.JsonNode, error) {
	if jn == nil || len(jn.KeyValues) == 0 {
		return nil, nil
	}
	var convertedKvs []*v2.JsonKeyValue
	for _, kv := range jn.KeyValues {
		newVal, err := translateJsonValue(kv.Value)
		if err != nil {
			return nil, err
		}

		newkv := &v2.JsonKeyValue{
			Key: kv.Key,
			Value: &v2.JsonValue{
				JsonVal: newVal.JsonVal,
			},
		}
		convertedKvs = append(convertedKvs, newkv)
	}
	return &v2.JsonNode{KeyValues: convertedKvs}, nil
}

func translateJsonValue(kv *v1alpha1.JsonValue) (*v2.JsonValue, error) {
	if kv == nil || kv.JsonVal == nil {
		return nil, nil
	}
	newkv := &v2.JsonValue{
		JsonVal: nil, // filled in later
	}

	switch jv := kv.GetJsonVal().(type) {
	case *v1alpha1.JsonValue_Node:
		recurseNode, err := translateJsonNode(jv.Node)
		if err != nil {
			return nil, err
		}
		node := &v2.JsonValue_Node{Node: recurseNode}
		newkv.JsonVal = node
	case *v1alpha1.JsonValue_List:
		var convertedList []*v2.JsonValue
		for _, val := range jv.List.Values {
			newVal, err := translateJsonValue(val)
			if err != nil {
				return nil, err
			}
			convertedList = append(convertedList, newVal)
		}

		list := &v2.JsonValue_List{
			List: &v2.JsonValueList{
				Values: convertedList,
			},
		}
		newkv.JsonVal = list
	case *v1alpha1.JsonValue_ValueProvider:
		convertedVp, err := translateValueProvider(jv.ValueProvider)
		if err != nil {
			return nil, err
		}
		vp := &v2.JsonValue_ValueProvider{
			ValueProvider: convertedVp,
		}
		newkv.JsonVal = vp
	default:
		return nil, errors.Errorf("unimplemented json value type: %T", jv)
	}
	return newkv, nil
}

func translateStringValueProviderMap(headers map[string]*v1alpha1.ValueProvider) (map[string]*v2.ValueProvider, error) {
	if len(headers) == 0 {
		return nil, nil
	}
	converted := map[string]*v2.ValueProvider{}
	for header, provider := range headers {
		vp, err := translateValueProvider(provider)
		if err != nil {
			return nil, err
		}
		converted[header] = vp
	}
	return converted, nil
}

func translateValueProvider(vp *v1alpha1.ValueProvider) (*v2.ValueProvider, error) {
	if vp == nil {
		return nil, nil
	}

	converted := &v2.ValueProvider{
		ProviderTemplate: vp.GetProviderTemplate(),
	}
	switch p := vp.Provider.(type) {
	case *v1alpha1.ValueProvider_GraphqlArg:
		ps, err := translatePath(p.GraphqlArg.Path)
		if err != nil {
			return nil, err
		}
		graphqlArg := &v2.ValueProvider_GraphqlArg{
			GraphqlArg: &v2.ValueProvider_GraphQLArgExtraction{
				ArgName:  p.GraphqlArg.ArgName,
				Path:     ps,
				Required: p.GraphqlArg.Required,
			},
		}
		converted.Provider = graphqlArg
	case *v1alpha1.ValueProvider_GraphqlParent:
		ps, err := translatePath(p.GraphqlParent.Path)
		if err != nil {
			return nil, err
		}
		graphqlParent := &v2.ValueProvider_GraphqlParent{
			GraphqlParent: &v2.ValueProvider_GraphQLParentExtraction{Path: ps},
		}
		converted.Provider = graphqlParent
	case *v1alpha1.ValueProvider_TypedProvider:
		t, err := translateType(p.TypedProvider.Type)
		if err != nil {
			return nil, err
		}
		tp := &v2.ValueProvider_TypedProvider{
			TypedProvider: &v2.ValueProvider_TypedValueProvider{
				Type:        t,
				ValProvider: nil, // filled in later
			},
		}
		switch vp := p.TypedProvider.ValProvider.(type) {
		case *v1alpha1.ValueProvider_TypedValueProvider_Header:
			convertedValProvider := &v2.ValueProvider_TypedValueProvider_Header{
				Header: vp.Header,
			}
			tp.TypedProvider.ValProvider = convertedValProvider
		case *v1alpha1.ValueProvider_TypedValueProvider_Value:
			convertedValProvider := &v2.ValueProvider_TypedValueProvider_Value{
				Value: vp.Value,
			}
			tp.TypedProvider.ValProvider = convertedValProvider
		default:
			return nil, errors.Errorf("unimplemented val provider type: %T", vp)
		}
		converted.Provider = tp
	default:
		return nil, errors.Errorf("unimplemented value provider type: %T", p)
	}
	return converted, nil
}

func translateType(t v1alpha1.ValueProvider_TypedValueProvider_Type) (v2.ValueProvider_TypedValueProvider_Type, error) {
	switch t.Enum().Number() {
	case v1alpha1.ValueProvider_TypedValueProvider_STRING.Number():
		return v2.ValueProvider_TypedValueProvider_STRING, nil
	case v1alpha1.ValueProvider_TypedValueProvider_INT.Number():
		return v2.ValueProvider_TypedValueProvider_INT, nil
	case v1alpha1.ValueProvider_TypedValueProvider_FLOAT.Number():
		return v2.ValueProvider_TypedValueProvider_FLOAT, nil
	case v1alpha1.ValueProvider_TypedValueProvider_BOOLEAN.Number():
		return v2.ValueProvider_TypedValueProvider_BOOLEAN, nil
	default:
		return v2.ValueProvider_TypedValueProvider_STRING, errors.Errorf("unimplemented typed value provider type: %T", t)
	}
}

func translatePath(path []*v1alpha1.PathSegment) ([]*v2.PathSegment, error) {
	if len(path) == 0 {
		return nil, nil
	}
	var converted []*v2.PathSegment
	for _, pathSegment := range path {
		ps := &v2.PathSegment{}
		switch p := pathSegment.Segment.(type) {
		case *v1alpha1.PathSegment_Key:
			key := &v2.PathSegment_Key{
				Key: p.Key,
			}
			ps.Segment = key
		case *v1alpha1.PathSegment_Index:
			index := &v2.PathSegment_Index{
				Index: p.Index,
			}
			ps.Segment = index
		case *v1alpha1.PathSegment_All:
			all := &v2.PathSegment_All{
				All: p.All,
			}
			ps.Segment = all
		default:
			return nil, errors.Errorf("unimplemented path segment type: %T", p)
		}
		converted = append(converted, ps)
	}
	return converted, nil
}
