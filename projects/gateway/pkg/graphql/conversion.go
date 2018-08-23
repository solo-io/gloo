package graphql

import (
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-apiserver/pkg/api/userdata"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type Converter struct{}

func (c *Converter) ConvertInputUpstreams(upstream []InputUpstream) []*v1.Upstream {
	var result []*v1.Upstream
	for _, us := range upstream {
		result = append(result, c.ConvertInputUpstream(us))
	}
	return result
}

func (c *Converter) ConvertInputUpstream(upstream InputUpstream) *v1.Upstream {
	return &v1.Upstream{
		Name:              upstream.Name,
		Type:              upstream.Type,
		ConnectionTimeout: upstream.ConnectionTimeout.GetDuration(),
		Spec:              upstream.Spec.GetStruct(),
		Functions:         convertInputFunctions(upstream.Functions),
		ServiceInfo:       convertInputServiceInfo(upstream.ServiceInfo),
		Metadata:          convertInputMetadata(upstream.Metadata),
	}
}

func convertInputServiceInfo(svcInfo *InputServiceInfo) *v1.ServiceInfo {
	if svcInfo == nil {
		return nil
	}
	return &v1.ServiceInfo{
		Type:       svcInfo.Type,
		Properties: svcInfo.Properties.GetStruct(),
	}
}

func convertInputFunctions(functions []*InputFunction) []*v1.Function {
	var v1Funcs []*v1.Function
	for _, fn := range functions {
		v1Funcs = append(v1Funcs, convertInputFunction(fn))
	}
	return v1Funcs
}

func convertInputFunction(function *InputFunction) *v1.Function {
	return &v1.Function{
		Name: function.Name,
		Spec: function.Spec.GetStruct(),
	}
}

func (c *Converter) ConvertOutputUpstreams(upstreams []*v1.Upstream) []*Upstream {
	var result []*Upstream
	for _, us := range upstreams {
		result = append(result, c.ConvertOutputUpstream(us))
	}
	return result
}

func (c *Converter) ConvertOutputUpstream(upstream *v1.Upstream) *Upstream {
	return &Upstream{
		Name:              upstream.Name,
		Type:              upstream.Type,
		ConnectionTimeout: NewDuration(upstream.ConnectionTimeout),
		Spec:              NewStruct(upstream.Spec),
		Functions:         convertOutputFunctions(upstream.Functions),
		ServiceInfo:       convertOutputServiceInfo(upstream.ServiceInfo),
		Status:            convertOutputStatus(upstream.Status),
		Metadata:          convertOutputMetadata(upstream.Metadata),
	}
}

func convertOutputServiceInfo(svcInfo *v1.ServiceInfo) *ServiceInfo {
	if svcInfo == nil {
		return nil
	}
	return &ServiceInfo{
		Type:       svcInfo.Type,
		Properties: NewStruct(svcInfo.Properties),
	}
}

func convertOutputFunctions(functions []*v1.Function) []*Function {
	var v1Funcs []*Function
	for _, fn := range functions {
		v1Funcs = append(v1Funcs, convertOutputFunction(fn))
	}
	return v1Funcs
}

func convertOutputFunction(function *v1.Function) *Function {
	return &Function{
		Name: function.Name,
		Spec: NewStruct(function.Spec),
	}
}

func (c *Converter) ConvertInputVirtualServices(virtualService []InputVirtualService) ([]*v1.VirtualService, error) {
	var result []*v1.VirtualService
	for _, vs := range virtualService {
		converted, err := c.ConvertInputVirtualService(vs)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

func (c *Converter) ConvertInputVirtualService(virtualService InputVirtualService) (*v1.VirtualService, error) {
	routes, err := c.ConvertInputRoutes(virtualService.Routes)
	if err != nil {
		return nil, errors.Wrap(err, "validating input routes")
	}

	return &v1.VirtualService{
		Name:      virtualService.Name,
		Domains:   dePointerify(virtualService.Domains),
		Routes:    routes,
		SslConfig: convertInputSSLConfig(virtualService.SslConfig),
		Roles:     dePointerify(virtualService.Roles),
		Metadata:  convertInputMetadata(virtualService.Metadata),
	}, nil
}

func (c *Converter) ConvertUserConfigProperties(userData *userdata.UserDataResponse) ([]*UserConfigProperty, error) {
	output := []*UserConfigProperty{}
	for _, ud := range userData.UserData {
		output = append(output, &UserConfigProperty{
			Key:   ud.Key,
			Value: ud.Value,
		})
	}
	return output, nil
}

func (c *Converter) ConvertInputRoutes(routes []*InputRoute) ([]*v1.Route, error) {
	var v1Routes []*v1.Route
	for _, fn := range routes {
		converted, err := c.ConvertInputRoute(fn)
		if err != nil {
			return nil, err
		}
		v1Routes = append(v1Routes, converted)
	}
	return v1Routes, nil
}

func (c *Converter) ConvertInputRoute(route *InputRoute) (*v1.Route, error) {
	var prefixRewrite string
	if route.PrefixRewrite != nil {
		prefixRewrite = *route.PrefixRewrite
	}
	v1Route := &v1.Route{
		PrefixRewrite: prefixRewrite,
		Extensions:    route.Extensions.GetStruct(),
	}
	switch {
	case route.Destination.MultiDestinations != nil && route.Destination.SingleDestination == nil:
		dest, err := convertInputDestination(*route.Destination.SingleDestination)
		if err != nil {
			return nil, err
		}
		v1Route.SingleDestination = dest
	case route.Destination.SingleDestination != nil && route.Destination.MultiDestinations == nil:
		weightedDestinations, err := convertDestinations(route.Destination.MultiDestinations)
		if err != nil {
			return nil, err
		}
		v1Route.MultipleDestinations = weightedDestinations
	default:
		return nil, errors.Errorf("must specify exactly one of SingleDestination or MultiDestinations")
	}

	switch {
	case route.Matcher.EventMatcher != nil && route.Matcher.RequestMatcher == nil:
		v1Route.Matcher = &v1.Route_EventMatcher{
			EventMatcher: &v1.EventMatcher{
				EventType: route.Matcher.EventMatcher.EventType,
			},
		}
	case route.Matcher.RequestMatcher != nil && route.Matcher.EventMatcher == nil:
		match, err := convertInputRequestMatcher(*route.Matcher.RequestMatcher)
		if err != nil {
			return nil, err
		}
		v1Route.Matcher = match
	default:
		return nil, errors.Errorf("must specify exactly one of RequestMatcher or EventMatcher")
	}
	return v1Route, nil
}

func convertDestinations(inputDests []*InputWeightedDestination) ([]*v1.WeightedDestination, error) {
	var weightedDests []*v1.WeightedDestination
	for _, inDest := range inputDests {
		dest, err := convertInputDestination(inDest.Destination)
		if err != nil {
			return nil, err
		}
		weightedDests = append(weightedDests, &v1.WeightedDestination{
			Destination: dest,
			Weight:      uint32(inDest.Weight),
		})
	}
	return weightedDests, nil
}

func convertInputDestination(inputDest InputSingleDestination) (*v1.Destination, error) {
	switch {
	case inputDest.UpstreamDestination != nil && inputDest.FunctionDestination == nil:
		return &v1.Destination{
			DestinationType: &v1.Destination_Upstream{
				Upstream: &v1.UpstreamDestination{
					Name: inputDest.UpstreamDestination.Name,
				},
			},
		}, nil
	case inputDest.FunctionDestination != nil && inputDest.UpstreamDestination == nil:
		return &v1.Destination{
			DestinationType: &v1.Destination_Function{
				Function: &v1.FunctionDestination{
					UpstreamName: inputDest.FunctionDestination.UpstreamName,
					FunctionName: inputDest.FunctionDestination.FunctionName,
				},
			},
		}, nil
	default:
		return nil, errors.Errorf("must specify exactly one of UpstreamDestination or FunctionDestination")
	}
}

func convertInputRequestMatcher(match InputRequestMatcher) (*v1.Route_RequestMatcher, error) {
	v1Match := &v1.RequestMatcher{
		Headers:     match.Headers.GetMap(),
		QueryParams: match.QueryParams.GetMap(),
		Verbs:       dePointerify(match.Verbs),
	}
	switch {
	case match.PathExact != nil && match.PathPrefix == nil && match.PathRegex == nil:
		v1Match.Path = &v1.RequestMatcher_PathExact{
			PathExact: match.PathExact.Path,
		}
	case match.PathPrefix != nil && match.PathExact == nil && match.PathRegex == nil:
		v1Match.Path = &v1.RequestMatcher_PathPrefix{
			PathPrefix: match.PathPrefix.Path,
		}
	case match.PathRegex != nil && match.PathExact == nil && match.PathPrefix == nil:
		v1Match.Path = &v1.RequestMatcher_PathRegex{
			PathRegex: match.PathRegex.Path,
		}
	default:
		return nil, errors.Errorf("must specify exactly one of PathPrefix PathRegex or PathExact")
	}
	return &v1.Route_RequestMatcher{RequestMatcher: v1Match}, nil
}

func convertInputSSLConfig(inSSL *InputSSLConfig) *v1.SSLConfig {
	if inSSL == nil {
		return nil
	}
	return &v1.SSLConfig{
		SecretRef: inSSL.SecretRef,
	}
}

func (c *Converter) ConvertOutputVirtualServices(virtualServices []*v1.VirtualService) []*VirtualService {
	var result []*VirtualService
	for _, vs := range virtualServices {
		result = append(result, c.ConvertOutputVirtualService(vs))
	}
	return result
}

func (c *Converter) ConvertOutputVirtualService(virtualService *v1.VirtualService) *VirtualService {
	return &VirtualService{
		Name:      virtualService.Name,
		Domains:   pointerify(virtualService.Domains),
		Routes:    convertOutputRoutes(virtualService.Routes),
		SslConfig: convertOutputSSLConfig(virtualService.SslConfig),
		Roles:     pointerify(virtualService.Roles),
		Status:    convertOutputStatus(virtualService.Status),
		Metadata:  convertOutputMetadata(virtualService.Metadata),
	}
}

func convertOutputRoutes(routes []*v1.Route) []*Route {
	var outRoutes []*Route
	for _, route := range routes {
		var (
			prefixRewrite *string
			matcher       Matcher
		)
		switch v1Match := route.Matcher.(type) {
		case *v1.Route_RequestMatcher:
			var path Path
			switch v1Path := v1Match.RequestMatcher.Path.(type) {
			case *v1.RequestMatcher_PathPrefix:
				path = &PathPrefix{
					Path: v1Path.PathPrefix,
				}
			case *v1.RequestMatcher_PathRegex:
				path = &PathRegex{
					Path: v1Path.PathRegex,
				}
			case *v1.RequestMatcher_PathExact:
				path = &PathExact{
					Path: v1Path.PathExact,
				}
			}
			matcher = &RequestMatcher{
				Path:        path,
				Headers:     NewMapStringString(v1Match.RequestMatcher.Headers),
				QueryParams: NewMapStringString(v1Match.RequestMatcher.QueryParams),
				Verbs:       pointerify(v1Match.RequestMatcher.Verbs),
			}
		}
		var destination Destination
		switch {
		case route.MultipleDestinations != nil:
			var weightedDestinations []*WeightedDestination
			for _, dest := range route.MultipleDestinations {
				weightedDestinations = append(weightedDestinations, &WeightedDestination{
					Destination: *convertSingleDestination(dest.Destination),
					Weight:      int(dest.Weight),
				})
			}
			destination = &MultiDestination{
				WeighedDestinations: weightedDestinations,
			}
		case route.SingleDestination != nil:
			destination = convertSingleDestination(route.SingleDestination)
		}
		if route.PrefixRewrite != "" {
			prefixRewrite = &route.PrefixRewrite
		}
		outRoute := &Route{
			Matcher:       matcher,
			Destination:   destination,
			PrefixRewrite: prefixRewrite,
			Extensions:    NewStruct(route.Extensions),
		}
		outRoutes = append(outRoutes, outRoute)
	}
	return outRoutes
}

func convertSingleDestination(dest *v1.Destination) *SingleDestination {
	var destinationType SingleDestinationUnion
	switch destType := dest.DestinationType.(type) {
	case *v1.Destination_Upstream:
		destinationType = &UpstreamDestination{
			Name: destType.Upstream.Name,
		}
	case *v1.Destination_Function:
		destinationType = &FunctionDestination{
			UpstreamName: destType.Function.UpstreamName,
			FunctionName: destType.Function.FunctionName,
		}
	}
	return &SingleDestination{
		Destination: destinationType,
	}
}

func convertOutputSSLConfig(v1SSL *v1.SSLConfig) *SSLConfig {
	if v1SSL == nil {
		return nil
	}
	return &SSLConfig{
		SecretRef: v1SSL.SecretRef,
	}
}

// common

func convertInputMetadata(meta *InputMetadata) *v1.Metadata {
	if meta == nil {
		return nil
	}
	var namespace string
	if meta.Namespace != nil {
		namespace = *meta.Namespace
	}
	return &v1.Metadata{
		ResourceVersion: meta.ResourceVersion,
		Namespace:       namespace,
		Annotations:     meta.Annotations.GetMap(),
	}
}

func convertOutputMetadata(meta *v1.Metadata) *Metadata {
	if meta == nil {
		return nil
	}
	var namespace *string
	if meta.Namespace != "" {
		namespace = &meta.Namespace
	}
	return &Metadata{
		ResourceVersion: meta.ResourceVersion,
		Namespace:       namespace,
		Annotations:     NewMapStringString(meta.Annotations),
	}
}

func convertOutputStatus(status *v1.Status) *Status {
	if status == nil {
		return nil
	}
	var state State
	switch status.State {
	case v1.Status_Pending:
		state = StatePending
	case v1.Status_Accepted:
		state = StateAccepted
	case v1.Status_Rejected:
		state = StateRejected
	}
	var reason *string
	if status.Reason != "" {
		reason = &status.Reason
	}
	return &Status{
		State:  state,
		Reason: reason,
	}
}

func dePointerify(in []*string) []string {
	var out []string
	for _, p := range in {
		out = append(out, *p)
	}
	return out
}

func pointerify(in []string) []*string {
	var out []*string
	for _, p := range in {
		out = append(out, &p)
	}
	return out
}
