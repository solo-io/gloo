package bootstrap

import (
	"bytes"
	"context"
	"os"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap/validation"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

const defaultEnvoyPath = "/usr/local/bin/envoy"

func getEnvoyPath() string {
	ep := os.Getenv("ENVOY_BINARY_PATH")
	if len(ep) == 0 {
		ep = defaultEnvoyPath
	}
	return ep
}

func ValidateBootstrap(
	ctx context.Context,
	filterName string,
	msg proto.Message,
) error {
	bootstrapYaml, err := buildPerFilterBootstrapYaml(filterName, msg)
	if err != nil {
		return err
	}

	return validation.ValidateBootstrap(ctx, bootstrapYaml)
}

func buildPerFilterBootstrapYaml(filterName string, msg proto.Message) (string, error) {

	typedFilter, err := utils.MessageToAny(msg)
	if err != nil {
		return "", err
	}
	vhosts := []*envoy_config_route_v3.VirtualHost{
		{
			Name:    "placeholder_host",
			Domains: []string{"*"},
			TypedPerFilterConfig: map[string]*any.Any{
				filterName: {
					TypeUrl: typedFilter.GetTypeUrl(),
					Value:   typedFilter.GetValue(),
				},
			},
		},
	}

	rc := &envoy_config_route_v3.RouteConfiguration{VirtualHosts: vhosts}

	hcm := &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix:     "placeholder",
		RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{RouteConfig: rc},
	}

	hcmAny, err := utils.MessageToAny(hcm)
	if err != nil {
		return "", err
	}
	bootstrap := &envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Id:      "imspecial",
			Cluster: "doesntmatter",
		},
		StaticResources: &envoy_config_bootstrap_v3.Bootstrap_StaticResources{
			Listeners: []*envoy_config_listener_v3.Listener{
				{
					Name: "placeholder_listener",
					Address: &envoy_config_core_v3.Address{
						Address: &envoy_config_core_v3.Address_SocketAddress{SocketAddress: &envoy_config_core_v3.SocketAddress{
							Address:       "0.0.0.0",
							PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{PortValue: 8081},
						}},
					},
					FilterChains: []*envoy_config_listener_v3.FilterChain{
						{
							Name: "placeholder_filter_chain",
							Filters: []*envoy_config_listener_v3.Filter{
								{
									ConfigType: &envoy_config_listener_v3.Filter_TypedConfig{
										TypedConfig: hcmAny,
									},
									Name: wellknown.HTTPConnectionManager,
								},
							},
						},
					},
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	marshaler := &jsonpb.Marshaler{
		OrigName: true,
	}
	marshaler.Marshal(buf, bootstrap)
	json := string(buf.Bytes())
	return json, nil // returns a json, but json is valid yaml
}
