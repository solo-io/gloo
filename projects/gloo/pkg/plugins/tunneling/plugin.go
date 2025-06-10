package tunneling

import (
	"errors"
	"fmt"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoytcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	_ plugins.Plugin                           = new(plugin)
	_ plugins.ResourceGeneratorPlugin          = new(plugin)
	_ plugins.UpstreamGeneratedResourcesPlugin = new(plugin)
)

const (
	ExtensionName                 = "tunneling"
	OriginalClusterSuffix         = "_original"
	forwardingListenerPrefix      = "solo_io_generated_self_listener_"
	forwardingListenerStatsPrefix = "soloioTcpStats"
)

type plugin struct{}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return ExtensionName
}

func (p *plugin) Init(_ plugins.InitParams) {
}

// UpstreamGeneratedResources checks the Upstream for a tunneling configuration
// and sets up clusters and listeners to forward traffic to an HTTP CONNECT supporting proxy.
//
// HTTP CONNECT tunneling is provided by an Envoy Listener filter. To send traffic to the Listener,
// we must generate a new forwarding Cluster. The generated Cluster sends traffic over
// a pipe to the Listener which forwards the traffic to the original Upstream.
//
// The SSL configuration for the original cluster is copied to the generated cluster and the
// supplied proxy SSL configuration is set on the original cluster.
//
// It's important that this method not modify the original cluster UNTIL AFTER all possible
// error cases have been checked. We don't want to partially transform the cluster and
// the additional cluster/listener to be generated if we encounter an error.
func (p *plugin) UpstreamGeneratedResources(
	params plugins.Params,
	in *v1.Upstream,
	out *envoy_config_cluster_v3.Cluster,
	reports reporter.ResourceReports,
) ([]*envoy_config_cluster_v3.Cluster, []*envoy_config_listener_v3.Listener, error) {
	var newClusters []*envoy_config_cluster_v3.Cluster
	var newListeners []*envoy_config_listener_v3.Listener
	var warning error

	// skip if the upstream does not have tunneling enabled
	httpProxyHostname := in.GetHttpProxyHostname().GetValue()
	if httpProxyHostname == "" {
		return nil, nil, nil
	}

	clusterName := out.GetName()

	// change the original cluster name to avoid conflicts with the new cluster
	newInClusterName := clusterName + OriginalClusterSuffix

	// use an in-memory pipe to ourselves (only works on linux)
	forwardingPipe := "@/" + clusterName

	var originalTransportSocket *envoy_config_core_v3.TransportSocket
	tunnelingHeaders := envoyHeadersFromHttpConnectHeaders(in)

	newListener, err := generateForwardingTcpListener(clusterName, newInClusterName,
		forwardingPipe, httpProxyHostname, tunnelingHeaders)
	if err != nil {
		return nil, nil, err
	}

	// we copy the transport socket to the generated cluster.
	// the generated cluster will use upstream TLS context to leverage TLS origination;
	if out.GetTransportSocket() != nil {
		tmp := *out.GetTransportSocket()
		originalTransportSocket = &tmp
	}

	// update the original cluster to use the new SSL configuration for the HTTP proxy, if provided
	if in.GetHttpConnectSslConfig() != nil {
		// user told us to configure ssl for the http connect proxy
		cfg, err := utils.NewSslConfigTranslator().ResolveUpstreamSslConfig(params.Snapshot.Secrets,
			in.GetHttpConnectSslConfig())
		if err != nil {
			// if we are configured to warn on missing tls secret and we match that error, add a
			// warning instead of error to the report.
			if params.Settings.GetGateway().GetValidation().GetWarnMissingTlsSecret().GetValue() &&
				errors.Is(err, utils.SslSecretNotFoundError) {
				reports.AddWarning(in, err.Error())
			} else {
				return nil, nil, err
			}
		} else {
			typedConfig, err := utils.MessageToAny(cfg)
			if err != nil {
				return nil, nil, err
			} else {
				out.TransportSocket = &envoy_config_core_v3.TransportSocket{
					Name:       wellknown.TransportSocketTls,
					ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{TypedConfig: typedConfig},
				}
			}
		}
	} else {
		// if no ssl config is provided, we need to remove the transport socket
		out.TransportSocket = nil
	}

	// when we encapsulate in HTTP Connect the tcp data being proxied will
	// be encrypted (thus we don't need the original transport socket metadata here)
	out.TransportSocketMatches = nil

	// generate new cluster with original cluster's name and transport socket that points to
	// the new listener's pipe
	newCluster := generateForwardingCluster(clusterName, forwardingPipe, originalTransportSocket)

	// to avoid having to change the parent's reference to the new cluster, we change the name
	// of the original cluster and use the original cluster name for the new cluster
	// this saves this plugin having to know about every place a parent may reference a cluster
	out.Name = newInClusterName

	newClusters = append(newClusters, newCluster)
	newListeners = append(newListeners, newListener)

	return newClusters, newListeners, warning
}

// GeneratedResources scans Upstreams for a tunneling configuration and sets up
// clusters and listeners to forward traffic to an HTTP CONNECT supporting proxy.
//
// Deprecated: This methods is for the Edge (non-krt) translator and is not safe
// to use with krt. Use UpstreamGeneratedResources instead.
//
// Before this can be removed, the stubbed error/warning handling in the krt path
// must be resolved.
func (p *plugin) GeneratedResources(
	params plugins.Params,
	proxy *v1.Proxy,
	inClusters []*envoy_config_cluster_v3.Cluster,
	inEndpoints []*envoy_config_endpoint_v3.ClusterLoadAssignment,
	inRouteConfigurations []*envoy_config_route_v3.RouteConfiguration,
	inListeners []*envoy_config_listener_v3.Listener,
	reports reporter.ResourceReports,
) (
	[]*envoy_config_cluster_v3.Cluster,
	[]*envoy_config_endpoint_v3.ClusterLoadAssignment,
	[]*envoy_config_route_v3.RouteConfiguration,
	[]*envoy_config_listener_v3.Listener,
) {
	var newClusters []*envoy_config_cluster_v3.Cluster
	var newListeners []*envoy_config_listener_v3.Listener

	// track the clusters we have transformed so we don't do it twice
	processedClusters := sets.Set[string]{}

	// find all upstreams with tunneling enabled
	for _, us := range params.Snapshot.Upstreams {
		clusterName := translator.UpstreamToClusterName(us.GetMetadata().Ref())

		// skip if the cluster has already been processed
		if processedClusters.Has(clusterName) {
			continue
		}

		// regardless of whether we successfully process the upstream, mark it as processed
		processedClusters.Insert(clusterName)

		// find the cluster to update
		cluster := findClusters(inClusters, clusterName)
		if cluster == nil {
			reports.AddError(us, fmt.Errorf("The cluster for the %s Upstream was not found."+
				"Check the status of the Upstream resources for errors.", clusterName))
			continue
		}

		upstreamsClusters, upstreamsListeners, err := p.UpstreamGeneratedResources(params, us, cluster, reports)
		if err != nil {
			reports.AddError(us, err)
			continue
		}

		newClusters = append(newClusters, upstreamsClusters...)
		newListeners = append(newListeners, upstreamsListeners...)
	}

	return newClusters, nil, nil, newListeners
}

// generateForwardingCluster generates a cluster will replace the original cluster and send
// the traffic to the generated listener
func generateForwardingCluster(
	clusterName,
	pipePath string,
	originalTransportSocket *envoy_config_core_v3.TransportSocket,
) *envoy_config_cluster_v3.Cluster {
	return &envoy_config_cluster_v3.Cluster{
		ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_STATIC,
		},
		ConnectTimeout:  &duration.Duration{Seconds: 5},
		Name:            clusterName,
		TransportSocket: originalTransportSocket,
		LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
				{
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_Pipe{
											Pipe: &envoy_config_core_v3.Pipe{
												Path: pipePath,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// generateForwardingTcpListener generates a listener that will forwards traffic to the
// HTTP Connect proxy
func generateForwardingTcpListener(
	cluster,
	originalCluster,
	pipePath,
	tunnelingHostname string,
	tunnelingHeadersToAdd []*envoy_config_core_v3.HeaderValueOption,
) (*envoy_config_listener_v3.Listener, error) {
	cfg := &envoytcp.TcpProxy{
		StatPrefix: forwardingListenerStatsPrefix + cluster,
		TunnelingConfig: &envoytcp.TcpProxy_TunnelingConfig{Hostname: tunnelingHostname,
			HeadersToAdd: tunnelingHeadersToAdd},
		ClusterSpecifier: &envoytcp.TcpProxy_Cluster{Cluster: originalCluster}, // route to original target
	}
	typedConfig, err := utils.MessageToAny(cfg)
	if err != nil {
		return nil, err
	}

	// FUTURE: considering using an internal listener
	// https://www.envoyproxy.io/docs/envoy/latest/configuration/other_features/internal_listener
	return &envoy_config_listener_v3.Listener{
		Name: forwardingListenerPrefix + cluster,
		Address: &envoy_config_core_v3.Address{
			Address: &envoy_config_core_v3.Address_Pipe{
				Pipe: &envoy_config_core_v3.Pipe{
					Path: pipePath,
				},
			},
		},
		FilterChains: []*envoy_config_listener_v3.FilterChain{
			{
				Filters: []*envoy_config_listener_v3.Filter{
					{
						Name: "tcp",
						ConfigType: &envoy_config_listener_v3.Filter_TypedConfig{
							TypedConfig: typedConfig,
						},
					},
				},
			},
		},
	}, nil
}

// envoyHeadersFromHttpConnectHeaders converts the http connect headers to envoy headers
func envoyHeadersFromHttpConnectHeaders(us *v1.Upstream) []*envoy_config_core_v3.HeaderValueOption {
	var tunnelingHeaders []*envoy_config_core_v3.HeaderValueOption
	for _, header := range us.GetHttpConnectHeaders() {
		tunnelingHeaders = append(tunnelingHeaders, &envoy_config_core_v3.HeaderValueOption{
			Header: &envoy_config_core_v3.HeaderValue{
				Key:   header.GetKey(),
				Value: header.GetValue(),
			},
			Append: &wrappers.BoolValue{Value: false},
		})
	}
	return tunnelingHeaders
}

// findClusters returns the cluster in the slice of clusters
func findClusters(clusters []*envoy_config_cluster_v3.Cluster, clusterName string) *envoy_config_cluster_v3.Cluster {
	for _, cluster := range clusters {
		if cluster.GetName() == clusterName {
			return cluster
		}
	}
	return nil
}
