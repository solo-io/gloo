package graphql

import (
	"log"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/customtypes"
	. "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/models"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
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
		Metadata:     convertInputMetadata(upstream.Metadata),
		UpstreamSpec: convertInputUpstreamSpec(upstream.Spec),
	}
}

func convertInputUpstreamSpec(spec InputUpstreamSpec) *v1.UpstreamSpec {
	out := &v1.UpstreamSpec{}
	switch {
	case spec.Aws != nil:
		out.UpstreamType = &v1.UpstreamSpec_Aws{
			Aws: &aws.UpstreamSpec{
				Region:          spec.Aws.Region,
				SecretRef:       spec.Aws.SecretRef,
				LambdaFunctions: convertLambdaFunctions(spec.Aws.Functions),
			},
		}
	case spec.Azure != nil:
		var ref string
		if spec.Azure.SecretRef != nil {
			ref = *spec.Azure.SecretRef
		}
		out.UpstreamType = &v1.UpstreamSpec_Azure{
			Azure: &azure.UpstreamSpec{
				FunctionAppName: spec.Azure.FunctionAppName,
				SecretRef:       ref,
				Functions:       convertAzureFunctions(spec.Azure.Functions),
			},
		}
	case spec.Kube != nil:
		out.UpstreamType = &v1.UpstreamSpec_Kube{
			Kube: &kubernetes.UpstreamSpec{
				Selector:         spec.Kube.Selector.GetMap(),
				ServiceName:      spec.Kube.ServiceName,
				ServiceNamespace: spec.Kube.ServiceNamespace,
				ServicePort:      uint32(spec.Kube.ServicePort),
			},
		}
	default:
		log.Printf("invalid spec: %#v", spec)
	}
	return out
}

func convertLambdaFunctions(inputFuncs []InputAwsLambdaFunction) []*aws.LambdaFunctionSpec {
	var funcs []*aws.LambdaFunctionSpec
	for _, inFn := range inputFuncs {
		funcs = append(funcs, &aws.LambdaFunctionSpec{
			LogicalName:        inFn.LogicalName,
			LambdaFunctionName: inFn.FunctionName,
			Qualifier:          inFn.Qualifier,
		})
	}
	return funcs
}

func convertAzureFunctions(inputFuncs []InputAzureFunction) []*azure.UpstreamSpec_FunctionSpec {
	var funcs []*azure.UpstreamSpec_FunctionSpec
	for _, inFn := range inputFuncs {
		funcs = append(funcs, &azure.UpstreamSpec_FunctionSpec{
			FunctionName: inFn.FunctionName,
			AuthLevel:    inFn.AuthLevel,
		})
	}
	return funcs
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
		Spec:     convertOutputUpstreamSpec(upstream.UpstreamSpec),
		Metadata: convertOutputMetadata(upstream.Metadata),
		Status:   convertOutputStatus(upstream.Status),
	}
}

func convertOutputUpstreamSpec(spec *v1.UpstreamSpec) UpstreamSpec {
	switch specType := spec.UpstreamType.(type) {
	case *v1.UpstreamSpec_Aws:
		return &AwsUpstreamSpec{
			Region:    specType.Aws.Region,
			SecretRef: specType.Aws.SecretRef,
			Functions: convertOutputLambdaFunctions(specType.Aws.LambdaFunctions),
		}
	case *v1.UpstreamSpec_Azure:
		return &AzureUpstreamSpec{
			FunctionAppName: specType.Azure.FunctionAppName,
			Functions:       convertOutputAzureFunctions(specType.Azure.Functions),
		}
	case *v1.UpstreamSpec_Kube:
		return &KubeUpstreamSpec{
			ServicePort:      int(specType.Kube.ServicePort),
			ServiceNamespace: specType.Kube.ServiceNamespace,
			ServiceName:      specType.Kube.ServiceName,
			Selector:         customtypes.NewMapStringString(specType.Kube.Selector),
		}
	}
	log.Printf("unsupported upstream type %v", spec)
	return nil
}

func convertOutputLambdaFunctions(lambdas []*aws.LambdaFunctionSpec) []AwsLambdaFunction {
	var out []AwsLambdaFunction
	for _, l := range lambdas {
		out = append(out, AwsLambdaFunction{
			LogicalName:  l.LogicalName,
			FunctionName: l.LambdaFunctionName,
			Qualifier:    l.Qualifier,
		})
	}
	return out
}

func convertOutputAzureFunctions(azureFns []*azure.UpstreamSpec_FunctionSpec) []AzureFunction {
	var out []AzureFunction
	for _, l := range azureFns {
		out = append(out, AzureFunction{
			FunctionName: l.FunctionName,
			AuthLevel:    l.AuthLevel,
		})
	}
	return out
}

func (c *Converter) ConvertInputVirtualServices(virtualService []InputVirtualService) ([]*gatewayv1.VirtualService, error) {
	var result []*gatewayv1.VirtualService
	for _, vs := range virtualService {
		converted, err := c.ConvertInputVirtualService(vs)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

func (c *Converter) ConvertInputVirtualService(virtualService InputVirtualService) (*gatewayv1.VirtualService, error) {
	routes, err := c.ConvertInputRoutes(virtualService.Routes)
	if err != nil {
		return nil, errors.Wrap(err, "validating input routes")
	}

	return &gatewayv1.VirtualService{
		VirtualHost: &v1.VirtualHost{
			Domains: virtualService.Domains,
			Routes:  routes,
		},
		SslConfig: convertInputSSLConfig(virtualService.SslConfig),
		Metadata:  convertInputMetadata(virtualService.Metadata),
	}, nil
}

func (c *Converter) ConvertInputRoutes(routes []InputRoute) ([]*v1.Route, error) {
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

func (c *Converter) ConvertInputRoute(route InputRoute) (*v1.Route, error) {
	match, err := convertInputMatcher(route.Matcher)
	if err != nil {
		return nil, err
	}
	action := &v1.Route_RouteAction{
		RouteAction: &v1.RouteAction{},
	}
	v1Route := &v1.Route{
		Matcher:      match,
		RoutePlugins: convertInputRoutePlugins(route.Plugins),
		Action:       action,
	}
	switch {
	case route.Destination.SingleDestination != nil:
		dest, err := convertInputSingleDestination(*route.Destination.SingleDestination)
		if err != nil {
			return nil, err
		}
		action.RouteAction.Destination = &v1.RouteAction_Single{
			Single: dest,
		}
	case route.Destination.MultiDestination != nil:
		weightedDestinations, err := convertInputDestinations(route.Destination.MultiDestination.Destinations)
		if err != nil {
			return nil, err
		}
		action.RouteAction.Destination = &v1.RouteAction_Multi{
			Multi: &v1.MultiDestination{
				Destinations: weightedDestinations,
			},
		}
	default:
		return nil, errors.Errorf("must specify exactly one of SingleDestination or MultiDestinations")
	}
	return v1Route, nil
}

func convertInputMatcher(match InputMatcher) (*v1.Matcher, error) {
	v1Match := &v1.Matcher{
		Headers:         convertInputHeaderMatcher(match.Headers),
		QueryParameters: convertInputQueryMatcher(match.QueryParameters),
		Methods:         match.Methods,
	}
	switch match.PathMatchType {
	case PathMatchTypeRegex:
		v1Match.PathSpecifier = &v1.Matcher_Regex{
			Regex: match.PathMatch,
		}
	case PathMatchTypeExact:
		v1Match.PathSpecifier = &v1.Matcher_Exact{
			Exact: match.PathMatch,
		}
	case PathMatchTypePrefix:
		v1Match.PathSpecifier = &v1.Matcher_Prefix{
			Prefix: match.PathMatch,
		}
	default:
		return nil, errors.Errorf("must specify one of PathPrefix PathRegex or PathExact")
	}
	return v1Match, nil
}

func convertInputHeaderMatcher(headers []InputKeyValueMatcher) []*v1.HeaderMatcher {
	var v1Headers []*v1.HeaderMatcher
	for _, h := range headers {
		v1Headers = append(v1Headers, &v1.HeaderMatcher{
			Name:  h.Name,
			Value: h.Value,
			Regex: h.IsRegex,
		})
	}
	return v1Headers
}

func convertInputQueryMatcher(queryM []InputKeyValueMatcher) []*v1.QueryParameterMatcher {
	var v1Query []*v1.QueryParameterMatcher
	for _, h := range queryM {
		v1Query = append(v1Query, &v1.QueryParameterMatcher{
			Name:  h.Name,
			Value: h.Value,
			Regex: h.IsRegex,
		})
	}
	return v1Query
}

func convertInputRoutePlugins(plugs *InputRoutePlugins) *v1.RoutePlugins {
	// TODO(ilackarms): convert route plugins when there are any
	return nil
}

func convertInputDestinations(inputDests []InputWeightedDestination) ([]*v1.WeightedDestination, error) {
	var weightedDests []*v1.WeightedDestination
	for _, inDest := range inputDests {
		dest, err := convertInputSingleDestination(inDest.Destination)
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

func convertInputDestinationSpec(spec *InputDestinationSpec) (*v1.DestinationSpec, error) {
	if spec == nil {
		return nil, nil
	}
	switch {
	case spec.Aws != nil:

	}
	return nil, errors.Errorf("unknown destination spec type: %#v", spec)
}

func convertInputSingleDestination(inputDest InputSingleDestination) (*v1.Destination, error) {
	destSpec, err := convertInputDestinationSpec(inputDest.DestinationSpec)
	if err != nil {
		return nil, err
	}
	return &v1.Destination{
		UpstreamName:    inputDest.UpstreamName,
		DestinationSpec: destSpec,
	}, nil
}

func convertInputSSLConfig(ssl *InputSslConfig) *v1.SslConfig {
	if ssl == nil {
		return nil
	}
	return &v1.SslConfig{
		SslSecrets: &v1.SslConfig_SecretRef{
			SecretRef: ssl.SecretRef,
		},
	}
}

func (c *Converter) ConvertOutputVirtualServices(virtualServices []*gatewayv1.VirtualService) []*VirtualService {
	var result []*VirtualService
	for _, vs := range virtualServices {
		result = append(result, c.ConvertOutputVirtualService(vs))
	}
	return result
}

func (c *Converter) ConvertOutputVirtualService(virtualService *gatewayv1.VirtualService) *VirtualService {
	return &VirtualService{
		Domains:   virtualService.VirtualHost.Domains,
		Routes:    convertOutputRoutes(virtualService.VirtualHost.Routes),
		SslConfig: convertOutputSSLConfig(virtualService.SslConfig),
		Status:    convertOutputStatus(virtualService.Status),
		Metadata:  convertOutputMetadata(virtualService.Metadata),
	}
}

func convertOutputRoutes(routes []*v1.Route) []Route {
	var outRoutes []Route
	for _, r := range routes {
		route, ok := convertOutputRoute(r)
		if !ok {
			continue
		}
		outRoutes = append(outRoutes, route)
	}
	return outRoutes
}

func convertOutputRoute(route *v1.Route) (Route, bool) {
	action, ok := route.Action.(*v1.Route_RouteAction)
	if !ok {
		log.Printf("warning: %v does not have a RouteAction, ignoring", route)
		return Route{}, false
	}
	var outDest Destination
	switch dest := action.RouteAction.Destination.(type) {
	case *v1.RouteAction_Single:
		outDest = convertOutputSingleDestination(dest.Single)
	case *v1.RouteAction_Multi:
		outDest = convertOutputMultiDestination(dest.Multi.Destinations)
	}
	return Route{
		Matcher:     convertOutputMatcher(route.Matcher),
		Destination: outDest,
		Plugins:     convertOutputRoutePlugins(route.RoutePlugins),
	}, true
}

func convertOutputMatcher(match *v1.Matcher) Matcher {
	var (
		path     string
		pathType PathMatchType
	)
	switch p := match.PathSpecifier.(type) {
	case *v1.Matcher_Exact:
		path = p.Exact
		pathType = PathMatchTypeExact
	case *v1.Matcher_Regex:
		path = p.Regex
		pathType = PathMatchTypeRegex
	case *v1.Matcher_Prefix:
		path = p.Prefix
		pathType = PathMatchTypePrefix
	}
	return Matcher{
		Headers:         convertOutputHeaderMatcher(match.Headers),
		QueryParameters: convertOutputQueryMatcher(match.QueryParameters),
		Methods:         match.Methods,
		PathMatch:       path,
		PathMatchType:   pathType,
	}
}

func convertOutputHeaderMatcher(headers []*v1.HeaderMatcher) []KeyValueMatcher {
	var v1Headers []KeyValueMatcher
	for _, h := range headers {
		v1Headers = append(v1Headers, KeyValueMatcher{
			Name:    h.Name,
			Value:   h.Value,
			IsRegex: h.Regex,
		})
	}
	return v1Headers
}

func convertOutputQueryMatcher(headers []*v1.QueryParameterMatcher) []KeyValueMatcher {
	var v1Headers []KeyValueMatcher
	for _, h := range headers {
		v1Headers = append(v1Headers, KeyValueMatcher{
			Name:    h.Name,
			Value:   h.Value,
			IsRegex: h.Regex,
		})
	}
	return v1Headers
}

func convertOutputRoutePlugins(plugs *v1.RoutePlugins) *RoutePlugins {
	// TODO(ilackarms): convert route plugins when there are any
	return nil
}

func convertOutputMultiDestination(dests []*v1.WeightedDestination) *MultiDestination {
	var weightedDests []WeightedDestination
	for _, v1Dest := range dests {
		weightedDests = append(weightedDests, WeightedDestination{
			Destination: convertOutputSingleDestination(v1Dest.Destination),
			Weight:      int(v1Dest.Weight),
		})
	}
	return &MultiDestination{Destinations: weightedDests}
}

func convertOutputSingleDestination(dest *v1.Destination) SingleDestination {
	return SingleDestination{
		UpstreamName:    dest.UpstreamName,
		DestinationSpec: convertOutputDestinationSpec(dest.DestinationSpec),
	}
}

func convertOutputDestinationSpec(spec *v1.DestinationSpec) DestinationSpec {
	if spec == nil {
		return nil
	}
	switch destSpec := spec.DestinationType.(type) {
	case *v1.DestinationSpec_Aws:
		var invocationStyle AwsLambdaInvocationStyle
		switch destSpec.Aws.InvocationStyle {
		case aws.DestinationSpec_ASYNC:
			invocationStyle = AwsLambdaInvocationStyleAsync
		case aws.DestinationSpec_SYNC:
			invocationStyle = AwsLambdaInvocationStyleSync
		}
		return &AwsDestinationSpec{
			LogicalName:     destSpec.Aws.LogicalName,
			InvocationStyle: invocationStyle,
		}
	}
	log.Printf("unknown destination spec type: %#v", spec)
	return nil
}

func convertOutputSSLConfig(ssl *v1.SslConfig) *SslConfig {
	if ssl == nil {
		return nil
	}
	secret, ok := ssl.SslSecrets.(*v1.SslConfig_SecretRef)
	if !ok {
		// file not supported atm
		return nil
	}

	return &SslConfig{
		SecretRef: secret.SecretRef,
	}
}

func (c *Converter) ConvertOutputResolverMaps(resolverMaps []*sqoopv1.ResolverMap) []*ResolverMap {
	var result []*ResolverMap
	for _, us := range resolverMaps {
		result = append(result, c.ConvertOutputResolverMap(us))
	}
	return result
}

func (c *Converter) ConvertOutputResolverMap(resolverMap *sqoopv1.ResolverMap) *ResolverMap {
	return &ResolverMap{
		Status:   convertOutputStatus(resolverMap.Status),
		Metadata: convertOutputMetadata(resolverMap.Metadata),
	}
}

// common
func convertInputMetadata(inMeta InputMetadata) core.Metadata {
	return core.Metadata{
		Namespace:       inMeta.Namespace,
		Name:            inMeta.Name,
		ResourceVersion: inMeta.ResourceVersion,
		Labels:          inMeta.Labels.GetMap(),
		Annotations:     inMeta.Annotations.GetMap(),
	}
}

func convertOutputMetadata(meta core.Metadata) Metadata {
	return Metadata{
		Namespace:       meta.Namespace,
		Name:            meta.Name,
		ResourceVersion: meta.ResourceVersion,
		Labels:          customtypes.NewMapStringString(meta.Labels),
		Annotations:     customtypes.NewMapStringString(meta.Annotations),
	}
}

func convertOutputStatus(status core.Status) Status {
	status = status.Flatten()
	var state State
	switch status.State {
	case core.Status_Pending:
		state = StatePending
	case core.Status_Accepted:
		state = StateAccepted
	case core.Status_Rejected:
		state = StateRejected
	}
	var reason *string
	if status.Reason != "" {
		reason = &status.Reason
	}
	return Status{
		State:  state,
		Reason: reason,
	}
}
