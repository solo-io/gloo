package envoy

import (
	"fmt"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"maps"
	"slices"
	"strconv"
	"strings"

	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"google.golang.org/protobuf/proto"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (o *Outputs) doRoutes(routes *adminv3.RoutesConfigDump, parentRef gwv1.ParentReference, ri RouteInfo) {
	for _, dr := range routes.DynamicRouteConfigs {

		var r envoy_config_route_v3.RouteConfiguration
		err := dr.GetRouteConfig().UnmarshalTo(&r)
		if err != nil {
			panic(err)
		}
		if r.Name == ri.Rds {
			for _, vh := range r.VirtualHosts {
				h := o.convertVH(vh, ri.FiltersOnChain)
				h.Spec.ParentRefs = []gwv1.ParentReference{parentRef}
				o.AddRoute(h)
			}
		}
	}
}

func (o *Outputs) convertVH(vh *envoy_config_route_v3.VirtualHost, filtersOnChain map[string][]proto.Message) *gwv1.HTTPRoute {
	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: "gateway.networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "gloo-system",
		},
	}
	for _, d := range vh.Domains {
		// TODO: handle *?
		hr.Spec.Hostnames = append(hr.Spec.Hostnames, gwv1.Hostname(d))
	}
	hr.Name = vh.Name
	// TODO: handle per filter config
	vho, err := o.convertVhostPolicy(vh, filtersOnChain)
	if err != nil {
		panic(err)
	}
	if vho != nil {
		o.AddVirtualHostOption(vho)
	}
	for _, r := range vh.Routes {
		rule := o.convertRoute(r)
		rp, filters, err := o.convertRoutePolicy(r, filtersOnChain)
		if rp != nil {
			rp.Name = fmt.Sprintf("%s-%s", vh.Name, r.Name)
			// TODO: optimize and see if we can use targetRef if all httpRoutes have the same config
			rule.Filters = []gwv1.HTTPRouteFilter{
				{
					ExtensionRef: &gwv1.LocalObjectReference{
						Name:  gwv1.ObjectName(r.Name),
						Group: "gateway.solo.io",
						Kind:  "RouteOption",
					},
					Type: gwv1.HTTPRouteFilterExtensionRef,
				},
			}
			o.AddRouteOption(rp)
		} else if err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("error converting route policy %w", err))
		}
		rule.Filters = append(rule.Filters, filters...)
		rule.Retry = o.convertRetries(r)
		// TODO: convert route options
		hr.Spec.Rules = append(hr.Spec.Rules, rule)

		if !isEmpty(r) {
			o.Errors = append(o.Errors, fmt.Errorf("unhandled route %s", r.Name))
		}
	}
	return hr
}

func (o *Outputs) convertRetries(rt *envoy_config_route_v3.Route) *gwv1.HTTPRouteRetry {
	ret := rt.GetRoute().GetRetryPolicy()
	if ret == nil {
		return nil
	}
	n := ret.GetNumRetries().GetValue()
	if n == 0 {
		return nil
	}
	var codes []gwv1.HTTPRouteRetryStatusCode
	var backoff *gwv1.Duration

	if baseInterval := ret.GetRetryBackOff().GetBaseInterval(); baseInterval != nil {
		backoff = ptr.To(convertDuration(baseInterval))
	}
	retriesOn := sets.New(strings.Split(ret.RetryOn, ",")...)
	for retryon := range retriesOn {
		switch retryon {
		case "5xx":
			// no op
		case "gateway-error":
			codes = append(codes, 502, 503, 504)
			delete(retriesOn, retryon)
		case "retriable-4xx":
			codes = append(codes, 409)
			delete(retriesOn, retryon)
		case "retriable-status-codes":
			for _, code := range ret.GetRetriableStatusCodes() {
				codes = append(codes, gwv1.HTTPRouteRetryStatusCode(code))
			}
			delete(retriesOn, retryon)
		case "connect-failure":
			fallthrough
		case "refused-stream":
			fallthrough
		case "unavailable":
			fallthrough
		case "cancelled":
			// These should be on by default according to gw api spec
			delete(retriesOn, retryon)
		}
	}

	if len(retriesOn) > 0 {
		ret.RetryOn = strings.Join(slices.Collect(maps.Keys(retriesOn)), ",")
	} else {
		ret.RetryOn = ""
	}

	// zero out the fields we extracted.
	ret.RetriableStatusCodes = nil
	ret.NumRetries = nil
	if ret.GetRetryBackOff() != nil {
		ret.GetRetryBackOff().BaseInterval = nil
	}
	// if not empty, there are fields where that can't be expressed in gw api
	if !isEmpty(ret) {
		o.Errors = append(o.Errors, fmt.Errorf("retry policy is not empty %v", ret))
	} else {
		rt.GetRoute().RetryPolicy = nil
	}

	return &gwv1.HTTPRouteRetry{
		Attempts: ptr.To(int(n)),
		Codes:    codes,
		Backoff:  backoff,
	}
}

func (o *Outputs) convertVhostPolicy(rt *envoy_config_route_v3.VirtualHost, filtersOnChain map[string][]proto.Message) (*gatewaykube.VirtualHostOption, error) {
	if rt.GetTypedPerFilterConfig() == nil {
		return nil, nil
	}

	keys := slices.Collect(maps.Keys(rt.GetTypedPerFilterConfig()))
	for _, k := range keys {
		v, err := convertAny(rt.GetTypedPerFilterConfig()[k])
		if err != nil {
			o.Errors = append(o.Errors, err)
			continue
		}

		switch v := v.(type) {
		case *jwt.StagedJwtAuthnPerRoute:
			o.convertJwtStaged(filtersOnChain[k], v)
		default:
			o.Errors = append(o.Errors, fmt.Errorf("vhost: unhandled per filter config %v", v))
		}
	}

	return nil, nil
}

func (o *Outputs) convertRoute(rt *envoy_config_route_v3.Route) gwv1.HTTPRouteRule {
	var hrr gwv1.HTTPRouteRule
	if rt.GetMatch() != nil {
		hrr.Matches = []gwv1.HTTPRouteMatch{convertMatcher(rt.GetMatch())}
	}
	if rt.GetRoute() != nil {
		for _, cluster := range getClusters(rt.GetRoute()) {
			br, err := o.convertBackendRef(cluster, rt)
			if err != nil {
				o.Errors = append(o.Errors, fmt.Errorf("error converting backend ref %w", err))
				continue
			}
			if br == nil {
				panic(fmt.Errorf("backend ref is nil for cluster %s on route %s", cluster.ClusterName, rt.Name))
			}
			hrr.BackendRefs = append(hrr.BackendRefs, *br)
		}
	}
	return hrr
}

func (o *Outputs) convertBackendRef(cluster clusterRef, route *envoy_config_route_v3.Route) (*gwv1.HTTPBackendRef, error) {

	backendRef := &gwv1.HTTPBackendRef{}
	if cluster.Weight > 0 {
		backendRef.Weight = ptr.To(int32(cluster.Weight))
	}
	clusterName := cluster.ClusterName

	// need to determine if the cluster is an upstream of k8s service
	pipeCount := strings.Count(clusterName, "|")
	// if there are 3 then its an istio generated cluster name
	if pipeCount == 3 {
		return o.parseIstioEnvoyCluster(cluster, backendRef)
	}
	// non istio based calls, we need to get the LB endpoints or lambda
	c := o.GetClusterByName(clusterName)
	if c == nil {
		// If the customer is using VD with TLS termination there will be no cluster, instead we need to see if the
		// x-gloo-mesh-federated-host header is being added, is that an upstream?
		return o.convertVDWithTLSTerminationHosts(route, cluster)
	}
	var err error
	backendRef, err = o.convertExternalServices(c)
	if err != nil {
		return nil, err
	}

	return backendRef, nil
}

func (o *Outputs) convertExternalServices(c *envoy_config_cluster_v3.Cluster) (*gwv1.HTTPBackendRef, error) {

	if c == nil {
		panic("cluster is nil")
	}
	backendRef := &gwv1.HTTPBackendRef{}
	// check to see if upstream already exists
	up := o.GetUpstream(c.Name)
	if up != nil {
		backendRef.Kind = ptr.To(gwv1.Kind("Upstream"))
		backendRef.Group = ptr.To(gwv1.Group("gloo.solo.io"))
		backendRef.Port = ptr.To(gwv1.PortNumber(up.Spec.GetStatic().GetHosts()[0].Port))
		backendRef.Name = gwv1.ObjectName(up.Name)
		backendRef.Namespace = ptr.To(gwv1.Namespace(up.Namespace))
		return backendRef, nil
	}

	//if cluster.Weight > 0 {
	//	backendRef.Weight = ptr.To(int32(cluster.Weight))
	//}
	staticUpstream := &v1.Upstream_Static{
		Static: &static.UpstreamSpec{},
	}
	var sslConfig *ssl.UpstreamSslConfig

	staticUpstreamSpec, sslConfigSpec, err := o.generateTLSContext(c)
	if err != nil {
		return nil, err
	}
	sslConfig = sslConfigSpec
	if staticUpstreamSpec != nil {
		staticUpstream.Static = staticUpstreamSpec
	}

	var servicePort uint32
	//var host string

	if c.LoadAssignment != nil && len(c.LoadAssignment.Endpoints) > 0 {
		if len(c.LoadAssignment.Endpoints) > 1 {
			o.Errors = append(o.Errors, fmt.Errorf("multiple endpoints assigned to the same load-assignment for cluster %s", c.Name))
		}

		for i, ep := range c.LoadAssignment.Endpoints[0].LbEndpoints {
			port := ep.GetEndpoint().GetAddress().GetSocketAddress().GetPortValue()
			if staticUpstream.Static.Hosts == nil {
				staticUpstream.Static.Hosts = []*static.Host{}
			}
			staticUpstream.Static.Hosts = append(staticUpstream.Static.Hosts, &static.Host{
				Addr: ep.GetEndpoint().GetAddress().GetSocketAddress().GetAddress(),
				Port: port,
			})
			if i == 0 {
				//currently just grabbing the first port
				servicePort = port
			}
		}
	}
	upstream := &glookube.Upstream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Upstream",
			APIVersion: "gloo.solo.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: "gloo-system",
		},
		Spec: v1.Upstream{
			UpstreamType: staticUpstream,
			SslConfig:    sslConfig,
		},
	}

	o.AddUpstream(upstream)
	backendRef.Kind = ptr.To(gwv1.Kind("Upstream"))
	backendRef.Group = ptr.To(gwv1.Group("gloo.solo.io"))
	backendRef.Port = ptr.To(gwv1.PortNumber(servicePort))
	backendRef.Name = gwv1.ObjectName(upstream.Name)
	backendRef.Namespace = ptr.To(gwv1.Namespace(upstream.Namespace))
	return backendRef, nil
}

// outbound|443||sts-qa.grainger-development.auth0app.com
func (o *Outputs) generateTLSContext(cluster *envoy_config_cluster_v3.Cluster) (*static.UpstreamSpec, *ssl.UpstreamSslConfig, error) {

	if cluster.GetTransportSocket() != nil {
		if cluster.GetTransportSocket().Name == "envoy.transport_sockets.tls" {
			// we should create an upstream for this service as it's not local to k8s and is not a known VD
			var sslConfig *ssl.UpstreamSslConfig
			staticUpstream := &static.UpstreamSpec{
				UseTls: wrapperspb.Bool(true),
			}

			transportSocket, err := convertAny(cluster.GetTransportSocket().GetTypedConfig())
			if err != nil {
				return nil, nil, err
			}

			if transportSocket == nil || isEmpty(transportSocket) {
				return nil, nil, fmt.Errorf("transport socket is empty %s", cluster.GetTransportSocket().Name)
			}
			switch ts := transportSocket.(type) {
			case *tlsv3.UpstreamTlsContext:
				staticUpstream.AutoSniRewrite = wrapperspb.Bool(ts.AutoHostSni)

				sslConfig = &ssl.UpstreamSslConfig{
					Sni:       ts.Sni,
					OneWayTls: wrapperspb.Bool(true),
				}
				//sts.amazonaws.com
				//launchpoint.internal.graingercloud.com
				if ts.CommonTlsContext != nil {

					if ts.CommonTlsContext.GetCombinedValidationContext() != nil {
						defaultConfig := ts.CommonTlsContext.GetCombinedValidationContext().GetDefaultValidationContext()
						sdsConfig := ts.CommonTlsContext.GetCombinedValidationContext().GetValidationContextSdsSecretConfig()
						if sdsConfig != nil {
							sslFiles := &ssl.SSLFiles{}
							if strings.HasPrefix(sdsConfig.Name, "file-root:") {
								splitName := strings.Split(sdsConfig.Name, ":")
								sslFiles.RootCa = splitName[1]

							}
							if strings.HasPrefix(sdsConfig.Name, "file-cert:") {
								split1 := strings.Split(sdsConfig.Name, ":")
								split2 := strings.Split(split1[1], "~")
								sslFiles.TlsCert = split2[0]
								sslFiles.TlsKey = split2[1]
							}
							sslConfig.SslSecrets = &ssl.UpstreamSslConfig_SslFiles{
								SslFiles: sslFiles,
							}
						}
						if defaultConfig != nil {
							sslConfig.VerifySubjectAltName = []string{}
							for _, san := range defaultConfig.MatchTypedSubjectAltNames {
								//TODO do we need to find more types of matches?
								if san.GetMatcher() != nil && san.GetMatcher().GetExact() != "" {
									sslConfig.VerifySubjectAltName = append(sslConfig.VerifySubjectAltName, san.GetMatcher().GetExact())
									sslConfig.OneWayTls = wrapperspb.Bool(false)
								}
							}
						}
					}
				}

			default:
				//just looking for upstream tls context
				return nil, nil, fmt.Errorf("unknown transport socket %s", cluster.TransportSocket.GetConfigType())
			}

			return staticUpstream, sslConfig, nil
		} else {
			return nil, nil, fmt.Errorf("unknown transport socket %s", cluster.TransportSocket.Name)
		}
	}
	return nil, nil, nil
}

func (o *Outputs) convertVDWithTLSTerminationHosts(route *envoy_config_route_v3.Route, cluster clusterRef) (*gwv1.HTTPBackendRef, error) {
	var err error
	backendRef := &gwv1.HTTPBackendRef{}
	// check to see if upstream already exists
	up := o.GetUpstream(cluster.ClusterName)
	if up != nil {
		backendRef.Kind = ptr.To(gwv1.Kind("Upstream"))
		backendRef.Group = ptr.To(gwv1.Group("gloo.solo.io"))
		backendRef.Port = ptr.To(gwv1.PortNumber(up.Spec.GetStatic().GetHosts()[0].Port))
		backendRef.Name = gwv1.ObjectName(up.Name)
		backendRef.Namespace = ptr.To(gwv1.Namespace(up.Namespace))
		return backendRef, nil
	}

	if cluster.Weight > 0 {
		backendRef.Weight = ptr.To(int32(cluster.Weight))
	}
	upstream := &glookube.Upstream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Upstream",
			APIVersion: "gloo.solo.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cluster.ClusterName,
			Namespace: "gloo-system",
		},
		Spec: v1.Upstream{},
	}
	var host string
	var port int
	for _, v := range route.GetRequestHeadersToAdd() {
		if v.GetHeader() != nil {
			if v.GetHeader().Key == "x-gloo-mesh-federated-host" {
				// this is a VD route with TLS Termination
				host = v.GetHeader().GetValue()
			}
			if v.GetHeader().Key == "x-gloo-mesh-federated-port" {
				// this is a VD route with TLS Termination
				port, err = strconv.Atoi(v.GetHeader().GetValue())
				if err != nil {
					return nil, fmt.Errorf("error converting port to int: %w for cluster %s", err, cluster.ClusterName)
				}
			}
		}
	}
	if host == "" {
		return backendRef, fmt.Errorf("cluster %s not found but referenced by route", cluster.ClusterName)
	}
	staticUpstream := &static.UpstreamSpec{
		Hosts: []*static.Host{{Addr: host}},
	}
	upstream.Spec.UpstreamType = &v1.Upstream_Static{
		Static: staticUpstream,
	}
	o.AddUpstream(upstream)
	backendRef.Kind = ptr.To(gwv1.Kind("Upstream"))
	backendRef.Group = ptr.To(gwv1.Group("gloo.solo.io"))
	backendRef.Port = ptr.To(gwv1.PortNumber(port))
	backendRef.Name = gwv1.ObjectName(cluster.ClusterName)
	backendRef.Namespace = ptr.To(gwv1.Namespace(upstream.Namespace))
	return backendRef, nil
}

func (o *Outputs) parseIstioEnvoyCluster(cluster clusterRef, backendRef *gwv1.HTTPBackendRef) (*gwv1.HTTPBackendRef, error) {
	clusterName := cluster.ClusterName
	parsed := strings.Split(clusterName, "|")
	i, err := strconv.Atoi(parsed[1])
	if err != nil {
		return nil, err
	}
	servicePort := i
	serviceName := parsed[3]
	// Istio services are in the format of <inbound/outbound>|port||service-name
	if strings.HasSuffix(clusterName, "svc.cluster.local") {
		//its a k8s service
		serviceSplit := strings.Split(serviceName, ".")
		backendRef.Name = gwv1.ObjectName(serviceSplit[0])
		backendRef.Namespace = ptr.To(gwv1.Namespace(serviceSplit[1]))
		backendRef.Port = ptr.To(gwv1.PortNumber(servicePort))
		return backendRef, nil
	}
	// we expect a cluster to be generated for all istio
	c := o.GetClusterByName(clusterName)

	if c == nil {
		return nil, fmt.Errorf("istio cluster %s not found", clusterName)
	}
	// non local service (service entry or VD), lets look for mTLS
	foundIstioMTLS := false
	if len(c.GetTransportSocketMatches()) > 0 {
		for _, match := range c.GetTransportSocketMatches() {
			if match.Name == "tlsMode-istio" {
				foundIstioMTLS = true
			}
		}
	}
	if foundIstioMTLS {
		// if the cluster is not svc.cluster.local, probably a VirtualDestination
		backendRef.Name = gwv1.ObjectName(serviceName)
		backendRef.Kind = ptr.To(gwv1.Kind("Hostname"))
		backendRef.Group = ptr.To(gwv1.Group("networking.istio.io"))
		backendRef.Port = ptr.To(gwv1.PortNumber(servicePort))
		// these all need to live in the gateway ns
		backendRef.Namespace = ptr.To(gwv1.Namespace("gloo-system"))
		return backendRef, nil
	}

	// its an ExternalService
	return o.convertExternalServices(c)
}

type clusterRef struct {
	ClusterName string
	Weight      uint32
}

func getClusters(rt *envoy_config_route_v3.RouteAction) []clusterRef {
	var clusters []clusterRef
	if rt.GetCluster() != "" {
		clusters = append(clusters, clusterRef{ClusterName: rt.GetCluster()})
		rt.ClusterSpecifier = nil
	}
	if rt.GetWeightedClusters() != nil {
		for _, wc := range rt.GetWeightedClusters().Clusters {
			clusters = append(clusters, clusterRef{ClusterName: wc.Name, Weight: wc.Weight.GetValue()})
		}
		rt.ClusterSpecifier = nil
	}
	return clusters
}
