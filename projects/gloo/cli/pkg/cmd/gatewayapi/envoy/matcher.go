package envoy

import (
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func convertMatcher(match *envoy_config_route_v3.RouteMatch) gwv1.HTTPRouteMatch {
	gwMatch := gwv1.HTTPRouteMatch{}

	if match.GetPrefix() != "" {
		gwMatch.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchPathPrefix),
			Value: ptr.To(match.GetPrefix()),
		}
		match.PathSpecifier = nil
	}
	if match.GetPath() != "" {
		gwMatch.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchExact),
			Value: ptr.To(match.GetPath()),
		}
		match.PathSpecifier = nil
	}
	if match.GetSafeRegex() != nil && match.GetSafeRegex().Regex != "" {
		gwMatch.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchRegularExpression),
			Value: ptr.To(match.GetSafeRegex().Regex),
		}
		match.PathSpecifier = nil
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
					header.HeaderMatchSpecifier = nil
				}
			}
			if gwMatch.Headers == nil {
				gwMatch.Headers = make([]gwv1.HTTPHeaderMatch, 0)
			}
			gwMatch.Headers = append(gwMatch.Headers, gwHM)
		}

		var headersLeft []*envoy_config_route_v3.HeaderMatcher
		for _, header := range match.GetHeaders() {
			if isEmpty(header) {
				continue
			}
			headersLeft = append(headersLeft, header)
		}
		match.Headers = headersLeft
	}
	return gwMatch
}
