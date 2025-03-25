package envoy

import (
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
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
				o.Routes = append(o.Routes, h)
			}
		}
	}
}

func (o *Outputs) convertVH(vh *envoy_config_route_v3.VirtualHost, filtersOnChain map[string][]proto.Message) gwv1.HTTPRoute {
	hr := gwv1.HTTPRoute{
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
		o.VhostOptions = append(o.VhostOptions, vho)
	}
	for _, r := range vh.Routes {
		rule := o.convertRoute(r)
		rp, filters, err := o.convertRoutePolicy(r, filtersOnChain)
		if rp != nil {
			rp.Name = fmt.Sprintf("%s-%s", vh.Name, r.Name)
			// TODO: optimize and see if we can use targetRef if all routes have the same config
			rule.Filters = []gwv1.HTTPRouteFilter{
				{
					ExtensionRef: &gwv1.LocalObjectReference{
						Name:  gwv1.ObjectName(r.Name),
						Group: "gateway.solo.io",
						Kind:  "RouteOption",
					},
				},
			}
			o.RouteOptions = append(o.RouteOptions, rp)
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
			convertJwtStaged(filtersOnChain[k], v)
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
			br, err := convertBackendRef(cluster)
			if err != nil {
				o.Errors = append(o.Errors, fmt.Errorf("error converting backend ref %w", err))
				continue
			}
			hrr.BackendRefs = append(hrr.BackendRefs, *br)
		}
	}
	return hrr
}

func convertBackendRef(cluster clusterRef) (*gwv1.HTTPBackendRef, error) {

	backendRef := &gwv1.HTTPBackendRef{}
	if cluster.Weight > 0 {
		backendRef.Weight = ptr.To(int32(cluster.Weight))
	}

	// need to determine if the cluster is an upstream of k8s service

	serviceName := cluster.ClusterName
	parsed := strings.Split(serviceName, "|")

	if len(parsed) < 4 {
		return nil, fmt.Errorf("invalid service name %s", serviceName)
	}

	i, err := strconv.Atoi(parsed[1])
	if err != nil {
		return nil, err
	}
	backendRef.Port = ptr.To(gwv1.PortNumber(i))

	if strings.HasSuffix(parsed[3], "svc.cluster.local") {
		//its a k8s service
		serviceSplit := strings.Split(parsed[3], ".")
		backendRef.Name = gwv1.ObjectName(serviceSplit[0])
		backendRef.Namespace = ptr.To(gwv1.Namespace(serviceSplit[1]))
	} else {
		// if the cluster is not svc.cluster.local, probably a VirtualDestination
		backendRef.Name = gwv1.ObjectName(parsed[3])
		backendRef.Kind = ptr.To(gwv1.Kind("Hostname"))
		backendRef.Group = ptr.To(gwv1.Group("networking.istio.io"))
	}

	//TODO non k8s services
	return backendRef, nil
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
