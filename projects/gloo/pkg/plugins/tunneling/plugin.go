package tunneling

import (
	"errors"

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
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	_ plugins.Plugin                  = new(plugin)
	_ plugins.ResourceGeneratorPlugin = new(plugin)
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

// GeneratedResources scans Upstreams for a tunneling configuration and sets up
// clusters and listeners to forward traffic to an HTTP CONNECT supporting proxy.
//
// HTTP CONNECT tunneling is provided by an Envoy Listener filter. To send traffic to the Listener,
// we must generate a new forwarding Cluster. The generated Cluster sends traffic over
// a pipe to the Listener which forwards the traffic to the original Upstream.
//
// The SSL configuration for the original cluster is copied to the generated cluster and the
// supplied proxy SSL configuration is set on the original cluster.
func (p *plugin) GeneratedResources(params plugins.Params,
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
	error,
) {
	logger := contextutils.LoggerFrom(params.Ctx)

	var newClusters []*envoy_config_cluster_v3.Cluster
	var newListeners []*envoy_config_listener_v3.Listener

	// track the clusters we have transformed so we don't do it twice
	processedClusters := sets.Set[string]{}

	// find all upstreams with tunneling enabled
	for _, us := range params.Snapshot.Upstreams {
		clusterName := translator.UpstreamToClusterName(us.GetMetadata().Ref())

		// skip if the cluster has already been processed
		if processedClusters.Has(clusterName) {
			logger.Warnf("cluster %v already processed", clusterName)
			continue
		}

		// regardless of whether we successfully process the upstream, mark it as processed
		processedClusters.Insert(clusterName)

		newCluster, newListener, warnings, err := processUpstream(params, us, inClusters)
		// make sure to add warnings to the report before we handle the error
		if len(warnings) > 0 {
			for _, warning := range warnings {
				reports.AddWarnings(us, warning.Error())
			}
		}

		if err != nil {
			reports.AddError(us, err)
			continue
		}

		if newCluster != nil || newListener != nil {
			newClusters = append(newClusters, newCluster)
			newListeners = append(newListeners, newListener)
		}
	}

	return newClusters, nil, nil, newListeners, nil
}

func processUpstream(
	params plugins.Params,
	us *v1.Upstream,
	inClusters []*envoy_config_cluster_v3.Cluster,
) (*envoy_config_cluster_v3.Cluster, *envoy_config_listener_v3.Listener, []*translator.Warning, error) {
	var warningList []*translator.Warning

	// skip if the upstream does not have tunneling enabled
	httpProxyHostname := us.GetHttpProxyHostname().GetValue()
	if httpProxyHostname == "" {
		return nil, nil, nil, nil
	}

	clusterName := translator.UpstreamToClusterName(us.GetMetadata().Ref())

	// find the cluster to update
	cluster := findClusters(inClusters, clusterName)
	if cluster == nil {
		return nil, nil, nil, errors.New("cluster not found")
	}

	// change the original cluster name to avoid conflicts with the new cluster
	newOriginalClusterName := clusterName + OriginalClusterSuffix

	// use an in-memory pipe to ourselves (only works on linux)
	forwardingPipe := "@/" + clusterName

	var originalTransportSocket *envoy_config_core_v3.TransportSocket
	tunnelingHeaders := envoyHeadersFromHttpConnectHeaders(us)

	// create the listener with the tunneling configuration and point it the original clusters
	newListener, err := generateForwardingTcpListener(clusterName, newOriginalClusterName,
		forwardingPipe, httpProxyHostname, tunnelingHeaders)
	if err != nil {
		return nil, nil, nil, err
	}
	// we copy the transport socket to the generated cluster.
	// the generated cluster will use upstream TLS context to leverage TLS origination;
	if cluster.GetTransportSocket() != nil {
		tmp := *cluster.GetTransportSocket()
		originalTransportSocket = &tmp
	}

	// update the original cluster to use the new SSL configuration for the HTTP proxy, if provided
	if us.GetHttpConnectSslConfig() != nil {
		// user told us to configure ssl for the http connect proxy
		cfg, err := utils.NewSslConfigTranslator().ResolveUpstreamSslConfig(params.Snapshot.Secrets,
			us.GetHttpConnectSslConfig())
		if err != nil {
			// if we are configured to warn on missing tls secret and we match that error, add a
			// warning instead of error to the report.
			if params.Settings.GetGateway().GetValidation().GetWarnMissingTlsSecret().GetValue() &&
				errors.Is(err, utils.SslSecretNotFoundError) {
				warningList = append(warningList, &translator.Warning{
					Message: err.Error(),
				})
			} else {
				return nil, nil, nil, err
			}
		} else {
			typedConfig, err := utils.MessageToAny(cfg)
			if err != nil {
				return nil, nil, nil, err
			} else {
				cluster.TransportSocket = &envoy_config_core_v3.TransportSocket{
					Name:       wellknown.TransportSocketTls,
					ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{TypedConfig: typedConfig},
				}
			}
		}
	} else {
		cluster.TransportSocket = nil
	}

	// when we encapsulate in HTTP Connect the tcp data being proxied will
	// be encrypted (thus we don't need the original transport socket metadata here)
	cluster.TransportSocketMatches = nil

	// generate new cluster with original cluster's name and transport socket that points to
	// the new listener's pipe
	newCluster := generateForwardingCluster(clusterName, forwardingPipe, originalTransportSocket)

	// update the original route's cluster name
	cluster.Name = newOriginalClusterName

	return newCluster, newListener, warningList, nil
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
