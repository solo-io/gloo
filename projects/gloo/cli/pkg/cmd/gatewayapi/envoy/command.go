package envoy

import (
	"encoding/json"
	"fmt"
	v8 "github.com/cncf/xds/go/udpa/type/v1"
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	v3_extensions "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/cors/v3"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	faultv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/fault/v3"
	rbacv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	http_connection_managerv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	tcp_proxyv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	api "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	v4 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	jwt2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	rbac2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	v5 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/cors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/faultinjection"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/headers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/retries"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
	"strconv"
	"strings"
)

var (
	output = &GatewayAPIOutput{
		HTTPRoutes:         make([]*gwv1.HTTPRoute, 0),
		RouteOptions:       make([]*gatewaykube.RouteOption, 0),
		VirtualHostOptions: make([]*gatewaykube.VirtualHostOption, 0),
		Upstreams:          make([]*glookube.Upstream, 0),
		AuthConfigs:        make([]*v1.AuthConfig, 0),
		Gateways:           make([]*gwv1.Gateway, 0),
	}
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
	//data, err := ioutil.ReadFile("nick/shipt/envoy.nick.json")
	data, err := ioutil.ReadFile("nick/grainger/config_dump.grainger.nick.json")
	//data, err := ioutil.ReadFile("nick/demo/demo-central.nick.json")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("Parsing envoy configuration")
	// Parse the configuration
	snapshot, err := parseEnvoyConfig(data)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Completed parsing envoy configuration")

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

	for _, listener := range snapshot.Listeners.DynamicListeners {
		var v3Listener v3.Listener
		if err := listener.ActiveState.Listener.UnmarshalTo(&v3Listener); err != nil {
			return err
		}
		if v3Listener.Address.GetSocketAddress().Address == "0.0.0.0" {
			log.Printf("Evaluating Listener %v", listener.Name)
			// this is a listener we want to generate output for each port
			gwGateway := &gwv1.Gateway{
				//TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:      fmt.Sprintf("ingress-%d", v3Listener.Address.GetSocketAddress().GetPortValue()),
					Namespace: "gloo-system",
				},
				Spec: gwv1.GatewaySpec{
					GatewayClassName: "gloo-gateway",
					Listeners:        []gwv1.Listener{},
				},
				//Addresses:      nil,
				//Infrastructure: nil,
				//BackendTLS:     nil,
			}

			for fi, fc := range v3Listener.FilterChains {
				listeners, err := processFilterChain(gwGateway.Name, gwGateway.Namespace, snapshot, fc, fi, &v3Listener)
				if err != nil {
					return err
				}
				gwGateway.Spec.Listeners = append(gwGateway.Spec.Listeners, listeners...)
				output.Gateways = append(output.Gateways, gwGateway)
			}

		}
	}

	// write all the outputs to their files
	//only write or
	txt, err := output.ToString()
	if err != nil {
		return err
	}
	filename := "gateway-api.rendered.nick.yaml"
	_, _ = fmt.Fprintf(os.Stdout, "Writing File: %s\n", filename)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintf(file, "%s", txt)
	//_, _ = fmt.Fprintf(os.Stdout, "%s\n", txt)

	return nil
}

// For each filter chain we need to gather a bunch of information then process the routes.
// Need to get the HCM name, SNIs, and Filters that need to be available to reference from the routes
func processFilterChain(gwName string, gwNamespace string, snapshot *EnvoySnapshot, fc *v3.FilterChain, fi int, v3Listener *v3.Listener) ([]gwv1.Listener, error) {
	var snis []string
	var gwListeners []gwv1.Listener
	if fc.FilterChainMatch != nil {
		var err error
		var tlsContext *gwv1.GatewayTLSConfig
		//grabs the tls information if it exists
		tlsContext, snis, err = findTLSContext(fc)
		if err != nil {
			return nil, err
		}

		var jwtProviders map[string]interface{}
		jwtProviders, _, err = findJWTProviders(fc)

		if err != nil {
			return nil, err
		}

		//No SNIs exist so we pull it from the HTTP connection manager
		var routeName string // HCM Route Name for reference
		for i, filter := range fc.Filters {
			if filter.Name == "envoy.filters.network.tcp_proxy" {
				tcpListeners, err := generateTCPListeners(snis, fi, tlsContext, i, v3Listener.Address.GetSocketAddress().GetPortValue())
				if err != nil {
					return nil, err
				}
				for _, listener := range tcpListeners {
					gwListeners = append(gwListeners, *listener)
				}

				var tcpp tcp_proxyv3.TcpProxy
				if err := filter.GetTypedConfig().UnmarshalTo(&tcpp); err != nil {
					return nil, err
				}
				//TODO need to generate the TCP Route to the backend
				//tcpp.

			}
			if filter.Name == "envoy.filters.network.http_connection_manager" {
				httpListeners, err := generateHTTPListeners(filter, snis, fi, tlsContext, i, v3Listener.Address.GetSocketAddress().GetPortValue())
				if err != nil {
					return nil, err
				}
				gwListeners = append(gwListeners, httpListeners...)

				var hcm http_connection_managerv3.HttpConnectionManager
				if err := filter.GetTypedConfig().UnmarshalTo(&hcm); err != nil {
					return nil, err
				}
				if hcm.GetRds() != nil {
					//// pull in the route
					routeName = hcm.GetRds().RouteConfigName
					log.Printf("grabbing route: %s", routeName)
					//
					rt, err := snapshot.GetRouteByName(routeName)
					if err != nil {
						return nil, err
					}
					if rt == nil {
						log.Printf("Route not found: %s", routeName)
					}
					routes, upstreams, routeOptions, err := generateHTTPRoutes(gwName, gwNamespace, rt, snapshot, jwtProviders)
					if err != nil {
						return nil, err
					}
					for _, ro := range routeOptions {
						output.RouteOptions = append(output.RouteOptions, ro)
					}
					for _, upstream := range upstreams {
						output.Upstreams = append(output.Upstreams, upstream)
					}
					for _, rt := range routes {
						output.HTTPRoutes = append(output.HTTPRoutes, rt)
					}
				}
			}
		}
	}
	return gwListeners, nil
}

// TODO will need to find the other providers too
func findJWTProviders(fc *v3.FilterChain) (map[string]interface{}, map[string][]string, error) {
	jwtProviders := make(map[string]interface{})
	filterStateRules := make(map[string][]string)
	for _, filter := range fc.Filters {
		// extract all the JWT Policies to feed the RouteOptions
		if filter.Name == "envoy.filters.network.http_connection_manager" {
			var hcm http_connection_managerv3.HttpConnectionManager
			if err := filter.GetTypedConfig().UnmarshalTo(&hcm); err != nil {
				return nil, nil, err
			}

			for _, ftlr := range hcm.GetHttpFilters() {
				if ftlr.Name == "io.solo.filters.http.solo_jwt_authn_staged" {
					var ts v8.TypedStruct
					if err := ftlr.GetTypedConfig().UnmarshalTo(&ts); err != nil {
						return nil, nil, err
					}

					jsonData, err := ts.Value.MarshalJSON()
					if err != nil {
						return nil, nil, err
					}

					// Convert JSON to map interface
					var result map[string]interface{}
					err = json.Unmarshal(jsonData, &result)
					if err != nil {
						return nil, nil, err
					}

					// filter state rules
					fsr := result["jwt_authn"].(map[string]interface{})
					authConfigMapping := NewAuthGenerator()
					ros, err := authConfigMapping.TransformJWT(fsr)
					if err != nil {
						return nil, nil, err
					}
					for _, ro := range ros {
						output.RouteOptions = append(output.RouteOptions, ro)
					}
				}
			}
		}
	}
	return jwtProviders, filterStateRules, nil
}

func findTLSContext(fc *v3.FilterChain) (*gwv1.GatewayTLSConfig, []string, error) {
	// there is a TLS listener?
	var tlsContext *gwv1.GatewayTLSConfig

	snis := []string{}
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
			return nil, nil, err
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
	return tlsContext, snis, nil
}

func generateHTTPRoutes(gwName string, gwNamespace string, route *route.RouteConfiguration, snapshot *EnvoySnapshot, jwtProviders map[string]interface{}) ([]*gwv1.HTTPRoute, []*glookube.Upstream, []*gatewaykube.RouteOption, error) {
	httpRoutes := make([]*gwv1.HTTPRoute, 0)
	upstreams := make([]*glookube.Upstream, 0)
	routeOptions := make([]*gatewaykube.RouteOption, 0)

	for _, virtualHost := range route.VirtualHosts {
		// Generate a route per virtualhost
		gwRoute := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      virtualHost.Name,
				Namespace: "gloo-system",
			},
			Spec: gwv1.HTTPRouteSpec{
				CommonRouteSpec: gwv1.CommonRouteSpec{
					ParentRefs: []gwv1.ParentReference{
						{
							Namespace: ptr.To(gwv1.Namespace(gwNamespace)),
							Name:      gwv1.ObjectName(gwName),
						},
					},
				},
				Hostnames: make([]gwv1.Hostname, 0),
				Rules:     make([]gwv1.HTTPRouteRule, 0),
			},
		}
		for _, vhRoute := range virtualHost.Routes {
			// matches
			gwrr := &gwv1.HTTPRouteRule{
				Matches:     make([]gwv1.HTTPRouteMatch, 0),
				Filters:     make([]gwv1.HTTPRouteFilter, 0),
				BackendRefs: make([]gwv1.HTTPBackendRef, 0),
			}
			match, err := convertMatcher(vhRoute.Match)
			if err != nil {
				return nil, nil, nil, err
			}
			gwrr.Matches = append(gwrr.Matches, match)

			// Add the filters and the upstreams
			routeOption, err := generateRouteOption(route.Name, vhRoute, jwtProviders)
			if err != nil {
				return nil, nil, nil, err
			}
			if routeOption != nil {
				gwrr.Filters = append(gwrr.Filters, gwv1.HTTPRouteFilter{
					Type: "ExtensionRef",
					ExtensionRef: &gwv1.LocalObjectReference{
						Group: "gateway.solo.io",
						Kind:  "RouteOption",
						Name:  gwv1.ObjectName(routeOption.Name),
					},
				})
				routeOptions = append(routeOptions, routeOption)
			}

			//cluster lookup
			//"outbound|8080||campaign-service-webserver.mesh.internal"
			cluster, err := snapshot.GetClusterByName(vhRoute.GetRoute().GetCluster())
			if err != nil {
				return nil, nil, nil, err
			}
			if cluster == nil {
				log.Printf("cluster not found " + vhRoute.GetRoute().GetCluster())
			}
			if cluster != nil {
				backendRef, upstream, err := generateBackendRef(vhRoute, cluster)
				if err != nil {
					return nil, nil, nil, err
				}
				if upstream != nil {
					upstreams = append(upstreams, upstream)
				}
				if backendRef != nil {
					gwrr.BackendRefs = append(gwrr.BackendRefs, *backendRef)
				}
			}

			gwRoute.Spec.Rules = append(gwRoute.Spec.Rules, *gwrr)
		}
		httpRoutes = append(httpRoutes, gwRoute)
	}

	return httpRoutes, upstreams, routeOptions, nil
}

func generateRouteOption(routeName string, r *route.Route, jwtProviders map[string]interface{}) (*gatewaykube.RouteOption, error) {

	if r == nil {
		return nil, nil
	}
	ro := &gatewaykube.RouteOption{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RouteOption",
			APIVersion: gatewaykube.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      RandStringRunes(8),
			Namespace: "gloo-system",
		},
		Spec: api.RouteOption{
			Options: &v2.RouteOptions{},
		},
	}
	if r.GetRoute() != nil {
		if r.GetRoute().GetPrefixRewrite() != "" {
			ro.Spec.Options.PrefixRewrite = wrapperspb.String(r.GetRoute().GetPrefixRewrite())
		}
		if r.GetRoute().GetRegexRewrite() != nil && r.GetRoute().GetRegexRewrite().GetPattern() != nil {
			ro.Spec.Options.RegexRewrite = &v4.RegexMatchAndSubstitute{
				Pattern: &v4.RegexMatcher{
					EngineType: &v4.RegexMatcher_GoogleRe2{},
					Regex:      r.GetRoute().GetRegexRewrite().GetPattern().GetRegex(),
				},
				Substitution: r.GetRoute().GetRegexRewrite().GetSubstitution(),
			}
			r.GetRoute().GetRegexRewrite().String()
		}
		if r.GetRoute().RetryPolicy != nil {
			rtp := r.GetRoute().RetryPolicy
			roRTP := &retries.RetryPolicy{
				RetryOn:       rtp.RetryOn,
				NumRetries:    rtp.NumRetries.Value,
				PerTryTimeout: rtp.PerTryTimeout,
				//PriorityPredicate:    nil,
				RetriableStatusCodes: rtp.RetriableStatusCodes,
			}
			if rtp.RetryBackOff != nil {
				roRTP.RetryBackOff = &retries.RetryBackOff{
					BaseInterval: rtp.RetryBackOff.GetBaseInterval(),
					MaxInterval:  rtp.RetryBackOff.GetMaxInterval(),
				}
			}
			ro.Spec.Options.Retries = roRTP
		}
	}
	// filters
	if len(r.GetTypedPerFilterConfig()) > 0 {
		for filterName, filterConfig := range r.GetTypedPerFilterConfig() {
			if filterName == "envoy.filters.http.cors" {
				//cors
				var corsPolicy v3_extensions.CorsPolicy
				if err := filterConfig.UnmarshalTo(&corsPolicy); err != nil {
					return nil, err
				}
				roCors := generateCorsPolicy(&corsPolicy)
				ro.Spec.Options.Cors = roCors
			}
			if filterName == "envoy.filters.http.fault" {
				var fault faultv3.HTTPFault
				if err := filterConfig.UnmarshalTo(&fault); err != nil {
					return nil, err
				}
				roFault := generateRouteFaults(&fault)
				ro.Spec.Options.Faults = roFault
			}
			if filterName == "envoy.filters.http.rbac" {
				//var ts v8.TypedStruct
				//if err := filterConfig.UnmarshalTo(&ts); err != nil {
				//	return nil, err
				//}
				//
				//jsonData, err := ts.Value.MarshalJSON()
				//if err != nil {
				//	return nil, err
				//}
				//data := ts.Value.Fields
				//log.Fatalf("%v", data)
				//
				//// Convert JSON to RBACPerRoute
				//var rbacPerRoute rbacv3.RBACPerRoute
				//err = json.Unmarshal(jsonData, &rbacPerRoute)
				//if err != nil {
				//	return nil, err
				//}
				//
				//roRbac := generateJWTRBAC(&rbacPerRoute)
				//if roRbac != nil {
				//	ro.Spec.Options.Rbac = roRbac
				//}
			}
			if filterName == "io.solo.filters.http.solo_jwt_authn_staged" {
				// TODO for inline JWT policies we need to generate a new RouteOption. If its just a reference we need to look htat up
				//var ts v8.TypedStruct
				//if err := filterConfig.UnmarshalTo(&ts); err != nil {
				//	return nil, err
				//}
				//
				//jsonData, err := ts.Value.MarshalJSON()
				//
				//if err != nil {
				//	return nil, err
				//}
				//
				//// Convert JSON to RBACPerRoute
				//var jwtPerRoute jwt.SoloJwtAuthnPerRoute
				//var result map[string]interface{}
				//
				//err = json.Unmarshal(jsonData, &result)
				//if err != nil {
				//	return nil, err
				//}
				//
				//perRouteJson, err := json.Marshal(result["jwt_configs"].(map[string]interface{})["0"])
				//if err != nil {
				//	return nil, err
				//}

				//if err := json.Unmarshal(perRouteJson, &jwtPerRoute); err != nil {
				//	return nil, err
				//}
				//
				//jwtProviderName := jwtPerRoute.Requirement
				////jwtSpecName := strings.Split(jwtProviderName, ".")[2]
				//if jwtProviders[routeName].(map[string]interface{})[jwtProviderName] == nil {
				//	log.Fatalf("JWT provider %s not found", jwtProviderName)
				//}
				//jwtProvider := jwtProviders[jwtProviderName].(map[string]interface{})
				//
				//log.Printf("%v", jwtProvider)
				//
				//roProvider := &jwt2.Provider{
				//	Jwks:             nil,
				//	Audiences:        jwtProvider["audiences"].([]string),
				//	Issuer:           jwtProvider["issuer"].(string),
				//	TokenSource:      nil,
				//	KeepToken:        jwtProvider["forward"].(bool),
				//	ClaimsToHeaders:  nil,
				//	ClockSkewSeconds: wrapperspb.UInt32(jwtProvider["clock_skew_seconds"].(uint32)),
				//}
				//if jwtProvider["remote_jwks"] != nil {
				//remote_jwks:
				//                                      http_uri:
				//                                        uri: https://member-auth-poc.shipt.com/.well-known/jwks.json
				//                                        cluster: outbound|80||member-auth-poc.shipt.com
				//                                        timeout: 5s
				//                                      async_fetch:
				//                                        fast_listener: true
				//                                    forward: true
				//                                    from_headers:
				//                                      - name: Authorization
				//                                        value_prefix: "Bearer "
				//                                    from_params:
				//                                      - access_token
				//                                    payload_in_metadata: principal
				//                                    clock_skew_seconds: 60
				//roProvider.Jwks = roJWKS
				//}

				roJWT := &v2.RouteOptions_JwtProvidersStaged{
					JwtProvidersStaged: &jwt2.JwtStagedRouteProvidersExtension{
						AfterExtAuth: &jwt2.VhostExtension{
							Providers: map[string]*jwt2.Provider{
								//providerName: roProvider,
							},
							//TODO need to figure out AllowMissingOrFailed (filter_state_rules)
							//AllowMissingOrFailedJwt: false,
							//ValidationPolicy:        0,
						},
					},
				}
				ro.Spec.Options.JwtConfig = roJWT
			}
		}
	}

	if len(r.GetRequestHeadersToRemove()) > 0 {
		if ro.Spec.Options.HeaderManipulation == nil {
			ro.Spec.Options.HeaderManipulation = &headers.HeaderManipulation{}
		}
		ro.Spec.Options.HeaderManipulation.RequestHeadersToRemove = r.GetRequestHeadersToRemove()
	}
	if len(r.GetRequestHeadersToAdd()) > 0 {
		if ro.Spec.Options.HeaderManipulation == nil {
			ro.Spec.Options.HeaderManipulation = &headers.HeaderManipulation{
				RequestHeadersToAdd: []*core.HeaderValueOption{},
			}
		}
		for _, a := range r.GetRequestHeadersToAdd() {
			addRequestHeader := &core.HeaderValueOption{
				HeaderOption: &core.HeaderValueOption_Header{
					Header: &core.HeaderValue{
						Key:   a.Header.Key,
						Value: a.Header.Value,
					},
				},
			}
			if a.AppendAction == corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD {
				addRequestHeader.Append = wrapperspb.Bool(false)
			}
			if a.AppendAction == corev3.HeaderValueOption_APPEND_IF_EXISTS_OR_ADD {
				addRequestHeader.Append = wrapperspb.Bool(true)
			}
			ro.Spec.Options.HeaderManipulation.RequestHeadersToAdd = append(ro.Spec.Options.HeaderManipulation.RequestHeadersToAdd, addRequestHeader)
		}
	}
	if len(r.GetResponseHeadersToRemove()) > 0 {
		if ro.Spec.Options.HeaderManipulation == nil {
			ro.Spec.Options.HeaderManipulation = &headers.HeaderManipulation{}
		}
		ro.Spec.Options.HeaderManipulation.ResponseHeadersToRemove = r.GetResponseHeadersToRemove()
	}

	if len(r.GetResponseHeadersToAdd()) > 0 {
		if ro.Spec.Options.HeaderManipulation == nil {
			ro.Spec.Options.HeaderManipulation = &headers.HeaderManipulation{
				ResponseHeadersToAdd: []*headers.HeaderValueOption{},
			}
		}
		for _, a := range r.GetResponseHeadersToAdd() {
			responseHeaderToAdd := &headers.HeaderValueOption{
				Header: &headers.HeaderValue{
					Key:   a.Header.Key,
					Value: a.Header.Value,
				},
			}
			if a.AppendAction == corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD {
				responseHeaderToAdd.Append = wrapperspb.Bool(false)
			}
			if a.AppendAction == corev3.HeaderValueOption_APPEND_IF_EXISTS_OR_ADD {
				responseHeaderToAdd.Append = wrapperspb.Bool(true)
			}
			ro.Spec.Options.HeaderManipulation.ResponseHeadersToAdd = append(ro.Spec.Options.HeaderManipulation.ResponseHeadersToAdd, responseHeaderToAdd)
		}
	}
	return ro, nil
}

func generateJWTRBAC(rbac *rbacv3.RBACPerRoute) *rbac2.ExtensionSettings {
	// JWT RBAC
	roRbac := &rbac2.ExtensionSettings{
		Disable:  false,
		Policies: make(map[string]*rbac2.Policy),
	}
	if rbac.Rbac != nil && rbac.Rbac.Rules != nil && rbac.Rbac.Rules.Policies != nil {
		for name, policy := range rbac.Rbac.Rules.Policies {
			roPolicy := &rbac2.Policy{
				Principals: make([]*rbac2.Principal, 0),
			}
			//
			//                  cds-account-segmentation-dev-acct-seg-read-jwt-dev:
			//                    permissions:
			//                    - any: true
			//                    principals:
			//                    - metadata:
			//                        filter: envoy.filters.http.jwt_authn
			//                        path:
			//                        - key: principal
			//                        - key: scope
			//                        value:
			//                          list_match:
			//                            one_of:
			//                              string_match:
			//                                exact: read:SegOnboarding

			// add a principal for both scope and principal?
			roPolicy.Principals = append(roPolicy.Principals, &rbac2.Principal{
				JwtPrincipal: &rbac2.JWTPrincipal{
					Claims: map[string]string{
						"scope": policy.Principals[0].GetMetadata().Value.GetListMatch().GetOneOf().GetStringMatch().GetExact(),
					},
				},
			})
			roPolicy.Principals = append(roPolicy.Principals, &rbac2.Principal{
				JwtPrincipal: &rbac2.JWTPrincipal{
					Claims: map[string]string{
						"principal": policy.Principals[0].GetMetadata().Value.GetListMatch().GetOneOf().GetStringMatch().GetExact(),
					},
				},
			})

			roRbac.Policies[name] = roPolicy
		}
		return roRbac
	}
	return nil
}

func generateRouteFaults(fault *faultv3.HTTPFault) *faultinjection.RouteFaults {
	roFault := &faultinjection.RouteFaults{}

	if fault.Abort != nil && fault.Abort.GetHttpStatus() != 0 {
		roFault.Abort = &faultinjection.RouteAbort{
			Percentage: float32(100.0),
			HttpStatus: fault.Abort.GetHttpStatus(),
		}
	}
	if fault.Delay != nil && fault.Delay.GetFixedDelay().Seconds != 0 {
		roFault.Delay = &faultinjection.RouteDelay{
			Percentage: float32(100.0),
			FixedDelay: fault.Delay.GetFixedDelay(),
		}
	}
	return roFault
}

func generateCorsPolicy(corsPolicy *v3_extensions.CorsPolicy) *v5.CorsPolicy {
	roCors := &v5.CorsPolicy{
		AllowOrigin:      make([]string, 0),
		AllowOriginRegex: make([]string, 0),
		AllowMethods:     strings.Split(strings.ReplaceAll(corsPolicy.AllowMethods, " ", ""), ","),
		AllowHeaders:     strings.Split(strings.ReplaceAll(corsPolicy.AllowHeaders, " ", ""), ","),
		ExposeHeaders:    strings.Split(strings.ReplaceAll(corsPolicy.ExposeHeaders, " ", ""), ","),
		MaxAge:           corsPolicy.MaxAge,
		AllowCredentials: corsPolicy.AllowCredentials.Value,
		DisableForRoute:  false,
	}
	if corsPolicy.ShadowEnabled != nil && corsPolicy.ShadowEnabled.DefaultValue != nil && corsPolicy.ShadowEnabled.DefaultValue.Numerator != 100 {
		roCors.DisableForRoute = true
	}
	if len(corsPolicy.AllowOriginStringMatch) > 0 {
		for _, origin := range corsPolicy.AllowOriginStringMatch {
			if origin.GetExact() != "" {
				roCors.AllowOrigin = append(roCors.AllowOrigin, origin.GetExact())
			}
			if origin.GetSafeRegex() != nil && origin.GetSafeRegex().Regex != "" {
				roCors.AllowOriginRegex = append(roCors.AllowOriginRegex, origin.GetSafeRegex().Regex)
			}
		}
	}
	return roCors
}

func generateBackendRef(r *route.Route, cluster *envoy_config_cluster_v3.Cluster) (*gwv1.HTTPBackendRef, *glookube.Upstream, error) {

	backendRef := &gwv1.HTTPBackendRef{}
	if cluster == nil {
		return nil, nil, nil
	}
	// need to determine if the cluster is an upstream of k8s service
	if cluster.GetType() == envoy_config_cluster_v3.Cluster_EDS {
		if cluster.GetEdsClusterConfig() != nil && cluster.GetEdsClusterConfig().GetServiceName() != "" {
			serviceName := cluster.GetEdsClusterConfig().GetServiceName()
			parsed := strings.Split(serviceName, "|")

			if strings.HasSuffix(parsed[3], "svc.cluster.local") {
				//its a k8s service
				serviceSplit := strings.Split(parsed[3], ".")
				backendRef.Name = gwv1.ObjectName(serviceSplit[0])
				backendRef.Namespace = ptr.To(gwv1.Namespace(serviceSplit[1]))
				i, err := strconv.Atoi(parsed[1])
				if err != nil {
					return nil, nil, err
				}
				backendRef.Port = ptr.To(gwv1.PortNumber(i))
			} else {
				foundIstioMTLS := false
				if len(cluster.GetTransportSocketMatches()) > 0 {
					for _, match := range cluster.GetTransportSocketMatches() {
						if match.Name == "tlsMode-istio" {
							foundIstioMTLS = true
						}
					}
				}
				if foundIstioMTLS {
					// if the cluster is type EDS, has mTLS enabled, and is not svc.cluster.local, probably a VirtualDestination
					backendRef.Name = gwv1.ObjectName(parsed[3])
					backendRef.Kind = ptr.To(gwv1.Kind("Hostname"))
					backendRef.Group = ptr.To(gwv1.Group("networking.istio.io"))
					i, err := strconv.Atoi(parsed[1])
					if err != nil {
						return nil, nil, err
					}
					backendRef.Port = ptr.To(gwv1.PortNumber(i))
				} else {
					log.Printf("unknown cluster type, cant convert to backendRef %v", cluster.Name)
				}

			}
		}

		return backendRef, nil, nil
	}

	//TODO non k8s services
	return nil, nil, nil
}

func convertMatcher(match *route.RouteMatch) (gwv1.HTTPRouteMatch, error) {
	gwMatch := gwv1.HTTPRouteMatch{}

	if match.GetPrefix() != "" {
		gwMatch.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchPathPrefix),
			Value: ptr.To(match.GetPrefix()),
		}
	}
	if match.GetPath() != "" {
		gwMatch.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchExact),
			Value: ptr.To(match.GetPath()),
		}
	}
	if match.GetSafeRegex() != nil && match.GetSafeRegex().Regex != "" {
		gwMatch.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchRegularExpression),
			Value: ptr.To(match.GetSafeRegex().Regex),
		}
	}

	if len(match.GetHeaders()) > 0 {

		for _, header := range match.GetHeaders() {
			gwHM := gwv1.HTTPHeaderMatch{}
			gwHM.Name = gwv1.HTTPHeaderName(header.Name)

			if header.GetStringMatch() != nil {
				//TODO GWAPI does nto support prefix header matching
				//if header.GetStringMatch().GetPrefix() != "" {
				//	gwHM.Type = ptr.To(gwv1.HeaderMatchExact)
				//}
				if header.GetStringMatch().GetExact() != "" {
					gwHM.Type = ptr.To(gwv1.HeaderMatchExact)
					gwHM.Value = header.GetStringMatch().GetExact()
				}
			}
			if gwMatch.Headers == nil {
				gwMatch.Headers = make([]gwv1.HTTPHeaderMatch, 0)
			}
			gwMatch.Headers = append(gwMatch.Headers, gwHM)
		}
	}
	return gwMatch, nil
}

func generateHTTPListeners(filter *v3.Filter, snis []string, fi int, tlsContext *gwv1.GatewayTLSConfig, i int, port uint32) ([]gwv1.Listener, error) {
	httpListeners := make([]gwv1.Listener, 0)
	var hcm http_connection_managerv3.HttpConnectionManager
	if err := filter.GetTypedConfig().UnmarshalTo(&hcm); err != nil {
		return nil, err
	}
	//SNIs exist so we use those as the listener domains
	if len(snis) > 0 {
		for i, sni := range snis {
			listener := gwv1.Listener{
				Port:     gwv1.PortNumber(port),
				Hostname: ptr.To(gwv1.Hostname(sni)),
				Protocol: gwv1.HTTPSProtocolType,
				Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d", fi, i)),
				AllowedRoutes: &gwv1.AllowedRoutes{
					Namespaces: &gwv1.RouteNamespaces{
						From: ptr.To(gwv1.FromNamespaces("gloo-system")),
					},
					Kinds: []gwv1.RouteGroupKind{
						{
							Kind: "HTTPRoute",
						},
					},
				},
			}
			if tlsContext != nil {
				listener.TLS = tlsContext
			}
			httpListeners = append(httpListeners, listener)
		}
	} else {
		//TODO hcm http_filters become listener filters
		// there are no SNIs so we should look at the VirtualHosts
		if hcm.GetRouteConfig() != nil {
			for j, vh := range hcm.GetRouteConfig().VirtualHosts {
				for _, domain := range vh.Domains {
					listener := gwv1.Listener{
						Port:     gwv1.PortNumber(port),
						Hostname: ptr.To(gwv1.Hostname(domain)),
						Protocol: gwv1.HTTPProtocolType,
						Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d-%d", fi, i, j)),
						AllowedRoutes: &gwv1.AllowedRoutes{
							Namespaces: &gwv1.RouteNamespaces{
								From: ptr.To(gwv1.FromNamespaces("gloo-system")),
							},
							Kinds: []gwv1.RouteGroupKind{
								{
									Kind: "HTTPRoute",
								},
							},
						},
					}
					if tlsContext != nil {
						listener.TLS = tlsContext
						listener.Protocol = gwv1.HTTPSProtocolType
					}
					httpListeners = append(httpListeners, listener)
				}
			}

		} else { // hcm.GetRouteConfig
			// wild card listener
			listener := gwv1.Listener{
				Port:     gwv1.PortNumber(port),
				Hostname: ptr.To(gwv1.Hostname("*")),
				Protocol: gwv1.HTTPProtocolType,
				Name:     gwv1.SectionName(fmt.Sprintf("listener-%d-%d", fi, i)),
				AllowedRoutes: &gwv1.AllowedRoutes{
					Namespaces: &gwv1.RouteNamespaces{
						From: ptr.To(gwv1.FromNamespaces("gloo-system")),
					},
					Kinds: []gwv1.RouteGroupKind{
						{
							Kind: "HTTPRoute",
						},
					},
				},
			}
			if tlsContext != nil {
				listener.TLS = tlsContext
				listener.Protocol = gwv1.HTTPSProtocolType
			}
			httpListeners = append(httpListeners, listener)
		}
	}
	return httpListeners, nil
}

func generateTCPListeners(snis []string, fi int, tlsContext *gwv1.GatewayTLSConfig, i int, port uint32) ([]*gwv1.Listener, error) {
	tcpListeners := make([]*gwv1.Listener, 0)
	if len(snis) > 0 {
		for i, sni := range snis {
			listener := gwv1.Listener{
				Port:     gwv1.PortNumber(port),
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
			tcpListeners = append(tcpListeners, &listener)
		}
	} else {
		listener := gwv1.Listener{
			Port:     gwv1.PortNumber(port),
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
		tcpListeners = append(tcpListeners, &listener)
	}
	return tcpListeners, nil
}
