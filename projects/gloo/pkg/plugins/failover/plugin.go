package failover

import (
	"fmt"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_api_v2_endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/gogoutils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

//go:generate mockgen -destination mocks/mock_utils.go github.com/solo-io/gloo/projects/gloo/pkg/utils SslConfigTranslator

const (
	TransportSocketMatchKey = "envoy.transport_socket_match"
)

var (
	NoHealthCheckError = eris.New("No health checks found, at least one is required to enable failover")

	_ plugins.Plugin         = new(failoverPluginImpl)
	_ plugins.UpstreamPlugin = new(failoverPluginImpl)
	_ plugins.EndpointPlugin = new(failoverPluginImpl)
)

func NewFailoverPlugin(translator utils.SslConfigTranslator) plugins.Plugin {
	return &failoverPluginImpl{
		sslConfigTranslator: translator,
		endpoints:           map[core.ResourceRef][]*envoy_api_v2_endpoint.LocalityLbEndpoints{},
	}
}

type failoverPluginImpl struct {
	sslConfigTranslator utils.SslConfigTranslator
	endpoints           map[core.ResourceRef][]*envoy_api_v2_endpoint.LocalityLbEndpoints
}

func (f *failoverPluginImpl) Init(_ plugins.InitParams) error {
	return nil
}

func (f *failoverPluginImpl) ProcessUpstream(params plugins.Params, in *gloov1.Upstream, out *v2.Cluster) error {
	failoverCfg := in.GetFailover()
	if failoverCfg == nil {
		return nil
	}

	// If no health checks have been set on the Upstream, throw an error as this will cause failover to fail in envoy.
	if len(in.GetHealthChecks()) == 0 {
		return NoHealthCheckError
	}

	endpoints, matches, err := f.buildLocalityLBEndpoints(params, failoverCfg)
	if err != nil {
		return err
	}

	if out.GetType() == v2.Cluster_EDS {
		f.endpoints[core.ResourceRef{
			Name:      in.Metadata.Name,
			Namespace: in.Metadata.Namespace,
		}] = endpoints
	} else {
		// Otherwise add the endpoints directly to the LoadAssignment of the Cluster
		out.LoadAssignment.Endpoints = append(out.LoadAssignment.Endpoints, endpoints...)
	}
	// append all of the upstream ssl transport socket matches to the existing
	out.TransportSocketMatches = append(out.TransportSocketMatches, matches...)

	return nil
}

func (f *failoverPluginImpl) ProcessEndpoints(
	params plugins.Params,
	in *gloov1.Upstream,
	out *v2.ClusterLoadAssignment,
) error {
	failoverCfg := in.GetFailover()
	if failoverCfg == nil {
		return nil
	}

	// If this is an eds cluster save the endpoints to the EDS `ClusterLoadAssignment`
	if len(f.endpoints) > 0 {
		out.Endpoints = append(out.Endpoints, f.endpoints[core.ResourceRef{
			Name:      in.Metadata.Name,
			Namespace: in.Metadata.Namespace,
		}]...)
	}

	return nil
}

func (f *failoverPluginImpl) buildLocalityLBEndpoints(
	params plugins.Params,
	failoverCfg *gloov1.Failover,
) ([]*envoy_api_v2_endpoint.LocalityLbEndpoints, []*v2.Cluster_TransportSocketMatch, error) {
	var transportSocketMatches []*v2.Cluster_TransportSocketMatch
	var localityLbEndpoints []*envoy_api_v2_endpoint.LocalityLbEndpoints
	for idx, priority := range failoverCfg.GetPrioritizedLocalities() {
		for _, localityEndpoints := range priority.GetLocalityEndpoints() {
			// Use index+1 for the priority because the priority of the primary endpoints is automatically set to 0,
			// and each corresponding failover endpoint has 1 greater
			envoyEndpoints, socketMatches, err := GlooLocalityLbEndpointToEnvoyLocalityLbEndpoint(
				localityEndpoints,
				uint32(idx+1),
				f.sslConfigTranslator,
				params.Snapshot.Secrets,
			)
			if err != nil {
				return nil, nil, err
			}
			localityLbEndpoints = append(localityLbEndpoints, envoyEndpoints)
			transportSocketMatches = append(transportSocketMatches, socketMatches...)
		}
	}

	return localityLbEndpoints, transportSocketMatches, nil
}

/*
	Create a unique name based on the details available while creating `LbEndpoints`
	The priority comes from the index in the top most loop of the API.
	idx is derived from the index of the endpoint in the list of `LbEndpoint`s.
	These names are used for the `MetadataMatch` in the `Cluster_TransportSocketMatch`
	https://www.envoyproxy.io/docs/envoy/v1.14.1/api-v2/api/v2/cluster.proto#envoy-api-msg-cluster-transportsocketmatch
*/
func PrioritizedEndpointName(address string, port, priority uint32, idx int) string {
	return fmt.Sprintf("%s_%d_p%d_idx%d", address, port, priority, idx)
}

func GlooLocalityToEnvoyLocality(locality *gloov1.Locality) *envoy_api_v2_core.Locality {
	if locality == nil {
		return nil
	}
	return &envoy_api_v2_core.Locality{
		Region:  locality.GetRegion(),
		Zone:    locality.GetZone(),
		SubZone: locality.GetSubZone(),
	}
}

func GlooLocalityLbEndpointToEnvoyLocalityLbEndpoint(
	endpoints *gloov1.LocalityLbEndpoints,
	priority uint32,
	translator utils.SslConfigTranslator,
	secrets []*gloov1.Secret,
) (*envoy_api_v2_endpoint.LocalityLbEndpoints, []*v2.Cluster_TransportSocketMatch, error) {
	var lbEndpoints []*envoy_api_v2_endpoint.LbEndpoint
	var transportSocketMatches []*v2.Cluster_TransportSocketMatch
	// Generate an envoy `LbEndpoint` for each endpoint in the locality.
	for idx, v := range endpoints.GetLbEndpoints() {
		uniqueName := PrioritizedEndpointName(v.GetAddress(), v.GetPort(), priority, idx)
		// Create a unique metadata match for each endpoint to support unique Transport Sockets in the Cluster
		metadataMatch := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				uniqueName: {
					Kind: &structpb.Value_BoolValue{
						BoolValue: true,
					},
				},
			},
		}
		lbEndpoint := &envoy_api_v2_endpoint.LbEndpoint{
			HostIdentifier: &envoy_api_v2_endpoint.LbEndpoint_Endpoint{
				Endpoint: &envoy_api_v2_endpoint.Endpoint{
					Address: &envoy_api_v2_core.Address{
						Address: &envoy_api_v2_core.Address_SocketAddress{
							SocketAddress: &envoy_api_v2_core.SocketAddress{
								Address: v.GetAddress(),
								PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
									PortValue: v.GetPort(),
								},
							},
						},
					},
				},
			},
			LoadBalancingWeight: gogoutils.UInt32GogoToProto(v.GetLoadBalancingWeight()),
		}
		if v.GetHealthCheckConfig() != nil {
			lbEndpoint.GetEndpoint().HealthCheckConfig = &envoy_api_v2_endpoint.Endpoint_HealthCheckConfig{
				PortValue: v.GetHealthCheckConfig().GetPortValue(),
				Hostname:  v.GetHealthCheckConfig().GetHostname(),
			}
		}
		// If UpstreamSslConfig is set, resolve it.
		if v.GetUpstreamSslConfig() != nil {
			cfg, err := translator.ResolveUpstreamSslConfig(secrets, v.GetUpstreamSslConfig())
			if err != nil {
				return nil, nil, err
			}
			anyCfg, err := utils.MessageToAny(cfg)
			if err != nil {
				return nil, nil, err

			}
			transportSocketMatches = append(transportSocketMatches, &v2.Cluster_TransportSocketMatch{
				Name:  uniqueName,
				Match: metadataMatch,
				TransportSocket: &envoy_api_v2_core.TransportSocket{
					Name: wellknown.TransportSocketTls,
					ConfigType: &envoy_api_v2_core.TransportSocket_TypedConfig{
						TypedConfig: anyCfg,
					},
				},
			})
			// Set the match criteria for the transport socket match on the endpoint
			lbEndpoint.Metadata = &envoy_api_v2_core.Metadata{
				FilterMetadata: map[string]*structpb.Struct{
					TransportSocketMatchKey: metadataMatch,
				},
			}
		}
		lbEndpoints = append(lbEndpoints, lbEndpoint)
	}
	return &envoy_api_v2_endpoint.LocalityLbEndpoints{
		Locality:            GlooLocalityToEnvoyLocality(endpoints.GetLocality()),
		LbEndpoints:         lbEndpoints,
		LoadBalancingWeight: gogoutils.UInt32GogoToProto(endpoints.GetLoadBalancingWeight()),
		Priority:            priority,
	}, transportSocketMatches, nil
}
