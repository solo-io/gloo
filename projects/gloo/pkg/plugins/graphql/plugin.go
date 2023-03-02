package graphql

import (
	"context"
	"strings"

	"github.com/solo-io/go-utils/contextutils"

	"github.com/solo-io/solo-projects/pkg/license"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/rotisserie/eris"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"
	"k8s.io/utils/lru"
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
	removeUnused              bool
	filterRequiredForListener map[*v1.HttpListener]struct{}
	// This cache prevents us from re-printing the same stitched schema on every translation loop, which can be expensive
	processedSchemaLruCache *lru.Cache
	// This cache prevents us from rerunning the stitching info JS script on every translation loop, which can be expensive
	stitchingInfoLruCache *lru.Cache
	// Is the graphQL plugin enabled
	graphQLFeatureState *license.FeatureState
	ctx                 context.Context
}

// NewPlugin creates the basic graphql plugin structure.
// context.Background() is used for ctx as we cant pass nil to contextutils.WithLogger() in Init
func NewPlugin(graphQLFeatureState *license.FeatureState) *plugin {
	return &plugin{
		processedSchemaLruCache: lru.New(1024),
		stitchingInfoLruCache:   lru.New(1024),
		graphQLFeatureState:     graphQLFeatureState,
		ctx:                     context.Background(),
	}
}

// Name returns the ExtensionName for overwriting purposes.
func (p *plugin) Name() string {
	return ExtensionName
}

// Init resets the plugin and creates the maps within the structure.
func (p *plugin) Init(params plugins.InitParams) {
	if params.Ctx != nil {
		p.ctx = params.Ctx
	}
	p.ctx = contextutils.WithLogger(p.ctx, "graphql-plugin")
	if p.graphQLFeatureState.Enabled {
		p.removeUnused = params.Settings.GetGloo().GetRemoveUnusedFilters().GetValue()
		p.filterRequiredForListener = make(map[*v1.HttpListener]struct{})
	}
	if p.graphQLFeatureState.Reason != "" {
		contextutils.LoggerFrom(p.ctx).Debug(p.graphQLFeatureState.Reason)
	}
}

// HttpFilters sets up the filters for envoy if it is needed.
func (p *plugin) HttpFilters(_ plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if !p.graphQLFeatureState.Enabled {
		return []plugins.StagedHttpFilter{}, nil
	}
	var filters []plugins.StagedHttpFilter

	_, ok := p.filterRequiredForListener[listener]
	if !ok && p.removeUnused {
		return filters, nil
	}

	emptyConf := &v2.GraphQLConfig{}
	stagedFilter, err := plugins.NewStagedFilter(FilterName, emptyConf, FilterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

// ProcessRoute applies any needed configurations related to graphql.
// If any configs are found then mark us needing this filter in our chain.
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	gqlRef := in.GetGraphqlApiRef()
	if gqlRef == nil {
		return nil
	}

	gql, err := params.Snapshot.GraphqlApis.Find(gqlRef.GetNamespace(), gqlRef.GetName())

	if !p.graphQLFeatureState.Enabled {
		contextutils.LoggerFrom(p.ctx).Errorf("Will not translate graphql api custom resource `%s` in namespace `%s`, no graphql features included with current license", gqlRef.GetName(), gqlRef.GetNamespace())
		return nil
	}
	if err != nil {
		ret := ""
		for _, api := range params.Snapshot.GraphqlApis {
			ret += " " + api.Metadata.Name
		}
		return eris.Wrapf(err, "unable to find graphql api custom resource `%s` in namespace `%s`, list of all graphqlapis found: %s", gqlRef.GetName(), gqlRef.GetNamespace(), ret)
	}

	p.filterRequiredForListener[params.HttpListener] = struct{}{} // Set here as user is at least attempting to use graphql at this point so might as well place it in the filterchain.
	routeConf, err := p.translateGraphQlApiToRouteConf(params, in, gql)

	if err != nil {
		return eris.Wrapf(err, "unable to translate graphql api control plane config to data plane config")
	}

	return pluginutils.SetRoutePerFilterConfig(out, FilterName, routeConf)
}

func (p *plugin) translateGraphQlApiToRouteConf(params plugins.RouteParams, in *v1.Route, api *v1beta1.GraphQLApi) (*v2.GraphQLRouteConfig, error) {
	execSchema, err := translation.CreateGraphQlApi(
		translation.CreateGraphQLApiParams{
			Artifacts:            params.Snapshot.Artifacts,
			Upstreams:            params.Snapshot.Upstreams,
			Graphqlapis:          params.Snapshot.GraphqlApis,
			Graphqlapi:           api,
			ProcessedSchemaCache: p.processedSchemaLruCache,
			StitchingInfoCache:   p.stitchingInfoLruCache,
		},
	)
	if err != nil {
		return nil, eris.Wrap(err, "error creating executable schema")
	}
	statsPrefix := in.GetGraphqlApiRef().Key()
	if sp := api.GetStatPrefix().GetValue(); sp != "" {
		statsPrefix = sp
	}
	statsPrefix = strings.TrimSuffix(statsPrefix, ".") + "."
	cacheConf := &v2.PersistedQueryCacheConfig{}
	if cc := api.GetPersistedQueryCacheConfig(); cc != nil {
		cacheConf.CacheSize = cc.CacheSize
	}
	return &v2.GraphQLRouteConfig{
		ExecutableSchema:          execSchema,
		StatPrefix:                statsPrefix,
		PersistedQueryCacheConfig: cacheConf,
		AllowedQueryHashes:        api.GetAllowedQueryHashes(),
	}, nil
}
