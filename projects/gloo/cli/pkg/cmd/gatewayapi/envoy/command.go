package envoy

import (
	"fmt"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	http_connection_managerv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcp_proxyv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	_ "istio.io/api/envoy/config/filter/network/metadata_exchange"
	_ "istio.io/api/envoy/extensions/stats"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/utils/ptr"
	"log"
	"os"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func init() {
	runtimeScheme = runtime.NewScheme()

	if err := SchemeBuilder.AddToScheme(runtimeScheme); err != nil {
		log.Fatal(err)
	}

	codecs = serializer.NewCodecFactory(runtimeScheme)
	decoder = codecs.UniversalDeserializer()
}

const (
	RandomSuffix = 4
	RandomSeed   = 1
)

var runtimeScheme *runtime.Scheme
var codecs serializer.CodecFactory
var decoder runtime.Decoder

func RootCmd() *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:   "envoy",
		Short: "Convert Envoy Config to Gateway API",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}
	opts.addToFlags(cmd.PersistentFlags())
	cmd.SilenceUsage = true
	return cmd
}

func run(opts *Options) error {
	// Read the Envoy configuration file
	//data, err := ioutil.ReadFile("envoy.nick.json")
	data, err := ioutil.ReadFile("config_dump.grainger.nick.json")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Parse the configuration
	snapshot, err := parseEnvoyConfig(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = generateGatwayAPIConfig(snapshot)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return nil
}

// Function to parse Envoy configuration
func parseEnvoyConfig(data []byte) (*EnvoySnapshot, error) {
	// Unmarshal the JSON into the ConfigDump struct
	var configDump adminv3.ConfigDump
	if err := protojson.Unmarshal(data, &configDump); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	envoysnapshot := &EnvoySnapshot{}

	for _, config := range configDump.Configs {
		if config.GetTypeUrl() == "type.googleapis.com/envoy.admin.v3.ListenersConfigDump" {
			if err := config.UnmarshalTo(&envoysnapshot.Listeners); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
		}
		if config.GetTypeUrl() == "type.googleapis.com/envoy.admin.v3.RoutesConfigDump" {
			if err := config.UnmarshalTo(&envoysnapshot.Routes); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
		}
		if config.GetTypeUrl() == "type.googleapis.com/envoy.admin.v3.ClustersConfigDump" {
			if err := config.UnmarshalTo(&envoysnapshot.Clusters); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
		}
	}

	return envoysnapshot, nil
}

func generateGatwayAPIConfig(snapshot *EnvoySnapshot) error {

	output := &GatewayAPIOutput{
		HTTPRoutes:         make([]*gwv1.HTTPRoute, 0),
		RouteOptions:       make([]*gatewaykube.RouteOption, 0),
		VirtualHostOptions: make([]*gatewaykube.VirtualHostOption, 0),
		Upstreams:          make([]*glookube.Upstream, 0),
		AuthConfigs:        make([]*v1.AuthConfig, 0),
		Gateways:           make([]*gwv1.Gateway, 0),
	}

	for _, listener := range snapshot.Listeners.DynamicListeners {
		var v3Listener v3.Listener
		if err := listener.ActiveState.Listener.UnmarshalTo(&v3Listener); err != nil {
			return err
		}
		if v3Listener.Address.GetSocketAddress().Address == "0.0.0.0" {
			// this is a listener we want to generate output for each port
			gwGateway := &gwv1.Gateway{
				//TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("ingress-%d", v3Listener.Address.GetSocketAddress().GetPortValue()),
					Namespace: "gloo-system",
				},
				Spec: gwv1.GatewaySpec{
					GatewayClassName: "gloo-gateway",
				},
				//Addresses:      nil,
				//Infrastructure: nil,
				//BackendTLS:     nil,
			}

			//TODO what to do about multiple filter chains?!!?
			var snis []string
			for fi, fc := range v3Listener.FilterChains {
				if fc.FilterChainMatch != nil {
					// there is a TLS listener?
					var tlsContext *gwv1.GatewayTLSConfig

					if fc.TransportSocket != nil && fc.TransportSocket.Name == "envoy.transport_sockets.tls" {
						//we need to generate a listener per SNI if they exist
						//TODO MTLS
						snis = fc.FilterChainMatch.ServerNames
						tlsContext = &gwv1.GatewayTLSConfig{
							Mode:            ptr.To(gwv1.TLSModeTerminate),
							CertificateRefs: []gwv1.SecretObjectReference{},
							// TODO CIPHERS
							//Options:            nil,
						}
						var downstreamTLSContext tlsv3.DownstreamTlsContext
						if err := fc.TransportSocket.GetTypedConfig().UnmarshalTo(&downstreamTLSContext); err != nil {
							return err
						}
						if len(downstreamTLSContext.CommonTlsContext.TlsCertificateSdsSecretConfigs) > 0 {
							for _, secret := range downstreamTLSContext.CommonTlsContext.TlsCertificateSdsSecretConfigs {
								//TODO no namespace support kubernetes://prod-wildcard-shipt-com-tls

								tlsContext.CertificateRefs = append(tlsContext.CertificateRefs, gwv1.SecretObjectReference{
									Name: gwv1.ObjectName(secret.Name[13:]), //remove kubernetes://
								})
							}
						}
					}

					//No SNIs exist so we pull it from the HTTP connection manager
					for i, filter := range fc.Filters {
						if filter.Name == "envoy.filters.network.tcp_proxy" {
							if len(snis) > 0 {
								for i, sni := range snis {
									listener := gwv1.Listener{
										Port:     gwv1.PortNumber(v3Listener.Address.GetSocketAddress().GetPortValue()),
										Hostname: ptr.To(gwv1.Hostname(sni)),
										Protocol: gwv1.TLSProtocolType,
										Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d", fi, i)),
										AllowedRoutes: &gwv1.AllowedRoutes{
											Kinds: []gwv1.RouteGroupKind{
												{
													Kind: "TCPRoute",
												},
											},
										},
									}
									if tlsContext != nil {
										listener.TLS = tlsContext
									}
									gwGateway.Spec.Listeners = append(gwGateway.Spec.Listeners, listener)
								}
							} else {
								listener := gwv1.Listener{
									Port:     gwv1.PortNumber(v3Listener.Address.GetSocketAddress().GetPortValue()),
									Protocol: gwv1.TCPProtocolType,
									Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d", fi, i)),
									AllowedRoutes: &gwv1.AllowedRoutes{
										Kinds: []gwv1.RouteGroupKind{
											{
												Kind: "TCPRoute",
											},
										},
									},
								}
								if tlsContext != nil {
									listener.TLS = tlsContext
								}
								gwGateway.Spec.Listeners = append(gwGateway.Spec.Listeners, listener)
							}

							var tcpp tcp_proxyv3.TcpProxy
							if err := filter.GetTypedConfig().UnmarshalTo(&tcpp); err != nil {
								return err
							}
							//TODO need to generate the TCP Route to the backend
							//tcpp.
							//	listener := gwv1.Listener{
							//	Port:     gwv1.PortNumber(v3Listener.Address.GetSocketAddress().GetPortValue()),
							//	Hostname: ptr.To(gwv1.Hostname(domain)),
							//	Protocol: gwv1.HTTPProtocolType,
							//	Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d-%d", fi, i, j)),

						}
						if filter.Name == "envoy.filters.network.http_connection_manager" {
							//SNIs exist so we use those as the listener domains
							if len(snis) > 0 {
								for i, sni := range snis {
									listener := gwv1.Listener{
										Port:     gwv1.PortNumber(v3Listener.Address.GetSocketAddress().GetPortValue()),
										Hostname: ptr.To(gwv1.Hostname(sni)),
										Protocol: gwv1.HTTPSProtocolType,
										Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d", fi, i)),
									}
									if tlsContext != nil {
										listener.TLS = tlsContext
									}
									gwGateway.Spec.Listeners = append(gwGateway.Spec.Listeners, listener)
								}
							}
							var hcm http_connection_managerv3.HttpConnectionManager
							if err := filter.GetTypedConfig().UnmarshalTo(&hcm); err != nil {
								return err
							}
							if hcm.GetRouteConfig() != nil {
								for j, vh := range hcm.GetRouteConfig().VirtualHosts {
									for _, domain := range vh.Domains {
										listener := gwv1.Listener{
											Port:     gwv1.PortNumber(v3Listener.Address.GetSocketAddress().GetPortValue()),
											Hostname: ptr.To(gwv1.Hostname(domain)),
											Protocol: gwv1.HTTPProtocolType,
											Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d-%d", fi, i, j)),
										}
										if tlsContext != nil {
											listener.TLS = tlsContext
											listener.Protocol = gwv1.HTTPSProtocolType
										}
										gwGateway.Spec.Listeners = append(gwGateway.Spec.Listeners, listener)
									}
								}
							}
							//// pull in the route
							//routeName := hcm.Config.GetRds().RouteConfigName
							////
							//rt, err := snapshot.GetRouteByName(routeName)
							//if err != nil {
							//	return err
							//}
							//listener.Name = gwv1.SectionName(routeName)

						}
					}
				}
			}

			output.Gateways = append(output.Gateways, gwGateway)
		}

	}

	// write all the outputs to their files
	//only write or
	txt, err := output.ToString()
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", txt)

	return nil
}
