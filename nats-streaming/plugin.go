package gfunc

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-plugins/common/annotations"
	"github.com/solo-io/gloo/pkg/coreplugins/common"
	"github.com/solo-io/gloo/pkg/plugin"
)

func init() {
	plugin.Register(&Plugin{}, nil)
}

type Plugin struct{}

const (
	ServiceTypeNatsStreaming = "nats-streaming"

	// generic plugin info
	filterName  = "io.solo.nats-streaming"
	pluginStage = plugin.OutAuth

	clusterId                 = "cluster_id"
	clusterIdAnnotations      = "gloo.solo.io/cluster_id"
	defaultClusterId          = "test-cluster"
	discoverPrefix            = "discover_prefix"
	discoverPrefixAnnotations = "gloo.solo.io/discover_prefix"
	defaultDiscoverPrefix     = "_STAN.discover"
)

func (p *Plugin) GetDependencies(cfg *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) HttpFilter(params *plugin.FilterPluginParams) (*envoyhttp.HttpFilter, plugin.Stage) {
	return &envoyhttp.HttpFilter{Name: filterName}, pluginStage
}

func (p *Plugin) ProcessRoute(_ *plugin.RoutePluginParams, in *v1.Route, out *envoyroute.Route) error {
	// nothing to do here
	return nil
}

func (p *Plugin) ProcessUpstream(params *plugin.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if in.Metadata.Annotations[annotations.ServiceType] != ServiceTypeNatsStreaming {
		return nil
	}
	//    string nats_streaming_cluster_id = 3;
	//    string discover_prefix = 4;
	// in.Metadata

	cid := in.Metadata.Annotations[clusterIdAnnotations]
	if cid == "" {
		cid = defaultClusterId
	}
	dp := in.Metadata.Annotations[clusterIdAnnotations]
	if dp == "" {
		dp = defaultDiscoverPrefix
	}
	if out.Metadata == nil {
		out.Metadata = &envoycore.Metadata{}
	}
	common.InitFilterMetadataField(filterName, clusterId, out.Metadata).Kind = &types.Value_StringValue{StringValue: defaultClusterId}
	common.InitFilterMetadataField(filterName, discoverPrefix, out.Metadata).Kind = &types.Value_StringValue{StringValue: dp}

	return nil
}

func (p *Plugin) ParseFunctionSpec(params *plugin.FunctionPluginParams, in v1.FunctionSpec) (*types.Struct, error) {

}
