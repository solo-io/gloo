package als

import (
	"context"
	"fmt"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var (
	_ plugins.Plugin                      = new(plugin)
	_ plugins.HttpConnectionManagerPlugin = new(plugin)
	_ plugins.ListenerPlugin              = new(plugin)
	_ plugins.ResourceGeneratorPlugin     = new(plugin)
)

const (
	ExtensionName = "als"
	ClusterName   = "access_log_cluster"
)

// Access logging plugin can contain a context for logging
type plugin struct {
	ctx context.Context
}

// NewPlugin creates an empty als plugin with no extra data
func NewPlugin() *plugin {
	return &plugin{}
}

// Name returns "als"
func (p *plugin) Name() string {
	return ExtensionName
}

// Init grabs the context for logging
func (p *plugin) Init(params plugins.InitParams) {
	p.ctx = params.Ctx
}

// ProcessHcmNetworkFilter will configure access logging for the hcm.
// This delegates most of its logic to ProcessAccessLogPlugins, which is also used by
// the TCP plugin and the listener level configuration.
func (p *plugin) ProcessHcmNetworkFilter(params plugins.Params, parentListener *v1.Listener,
	_ *v1.HttpListener, out *envoyhttp.HttpConnectionManager) error {
	if out == nil {
		return nil
	}

	// AccessLog settings are defined on the root listener, and applied to each HCM instance
	alsSettings := parentListener.GetOptions().GetAccessLoggingService()
	if alsSettings == nil {
		return nil
	}

	fmt.Printf("ProcessHcmNetworkFilter alsSettings: %v\n", alsSettings)

	var err error
	out.AccessLog, err = ProcessAccessLogPlugins(params, alsSettings, out.GetAccessLog())
	if err := DetectUnusefulCmds(Hcm, out.GetAccessLog()); err != nil {
		contextutils.LoggerFrom(p.ctx).Warnf("warning non-useful access log operator on %s's hcm: %s",
			parentListener.GetName(), err.Error())
	}
	return err
}

// ProcessListener will configure access logging at the listener level.
func (p *plugin) ProcessListener(params plugins.Params, parentListener *v1.Listener,
	out *envoy_config_listener_v3.Listener) error {

	alsSettings := parentListener.GetOptions().GetListenerAccessLoggingService()
	if alsSettings == nil {
		return nil
	}

	fmt.Printf("ProcessListener alsSettings: %v\n", alsSettings)

	var err error
	out.AccessLog, err = ProcessAccessLogPlugins(params, alsSettings, out.GetAccessLog())
	if err := DetectUnusefulCmds(HttpListener, out.GetAccessLog()); err != nil {
		contextutils.LoggerFrom(p.ctx).Warnf("non-useful access log operator configured on %s: %s",
			parentListener.GetName(), err.Error())
	}

	return err
}

// GeneratedResources scans the proxy for listeners with access logging settings
// and generates clusters if needed. This is mostly for the OTEL collector as
// we decided to not make customers configure upstreams.
func (p *plugin) GeneratedResources(
	params plugins.Params,
	proxy *v1.Proxy,
	_ []*envoy_config_cluster_v3.Cluster,
	_ []*envoy_config_endpoint_v3.ClusterLoadAssignment,
	_ []*envoy_config_route_v3.RouteConfiguration,
	_ []*envoy_config_listener_v3.Listener,
	reports reporter.ResourceReports) (
	[]*envoy_config_cluster_v3.Cluster,
	[]*envoy_config_endpoint_v3.ClusterLoadAssignment,
	[]*envoy_config_route_v3.RouteConfiguration,
	[]*envoy_config_listener_v3.Listener) {

	generatedClusters := []*envoy_config_cluster_v3.Cluster{}
	for _, listener := range proxy.GetListeners() {
		listenerOpts := listener.GetOptions()
		if listenerOpts == nil {
			continue
		}

		alsSettings := listenerOpts.GetAccessLoggingService()
		if alsSettings != nil {
			clusters := getClustersForAccessLogs(params, proxy, reports, alsSettings)
			generatedClusters = append(generatedClusters, clusters...)
		}

		alsSettings = listenerOpts.GetListenerAccessLoggingService()
		if alsSettings != nil {
			clusters := getClustersForAccessLogs(params, proxy, reports, alsSettings)
			generatedClusters = append(generatedClusters, clusters...)
		}
	}

	return generatedClusters, nil, nil, nil
}
