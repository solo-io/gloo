package graphql

import (
	"strings"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/rotisserie/eris"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
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

// ProcessRoute applies any needed configurations related to graphql.
// If any configs are found then mark us needing this filter in our chain.
func (p *plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoy_config_route_v3.Route) error {
	gqlRef := in.GetGraphqlApiRef()
	if gqlRef == nil {
		return nil
	}

	gql, err := params.Snapshot.GraphqlApis.Find(gqlRef.GetNamespace(), gqlRef.GetName())
	if err != nil {
		ret := ""
		for _, api := range params.Snapshot.GraphqlApis {
			ret += " " + api.Metadata.Name
		}
		return eris.Wrapf(err, "unable to find graphql api custom resource `%s` in namespace `%s`, list of all graphqlapis found: %s", gqlRef.GetName(), gqlRef.GetNamespace(), ret)
	}

	p.filterNeeded = true // Set here as user is at least attempting to use graphql at this point so might as well place it in the filterchain.
	routeConf, err := translateGraphQlApiToRouteConf(params, in, gql)

	if err != nil {
		return eris.Wrapf(err, "unable to translate graphql api control plane config to data plane config")
	}

	return pluginutils.SetRoutePerFilterConfig(out, FilterName, routeConf)
}

func translateGraphQlApiToRouteConf(params plugins.RouteParams, in *v1.Route, api *v1beta1.GraphQLApi) (*v2.GraphQLRouteConfig, error) {
	execSchema, err := translation.CreateGraphQlApi(params.Snapshot.Artifacts, params.Snapshot.Upstreams, params.Snapshot.GraphqlApis, api)
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
