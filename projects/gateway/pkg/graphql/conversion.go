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

// Extra

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
	match, err := convertInputRequestMatcher(route.Matcher)
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
	case route.Destination.MultiDestination != nil:
		dest, err := convertInputSingleDestination(*route.Destination.SingleDestination)
		if err != nil {
			return nil, err
		}
		action.RouteAction.Destination = dest
	case route.Destination.SingleDestination != nil:
		weightedDestinations, err := convertDestinations(route.Destination.MultiDestination)
		if err != nil {
			return nil, err
		}
		action.RouteAction.Destination = weightedDestinations
	default:
		return nil, errors.Errorf("must specify exactly one of SingleDestination or MultiDestinations")
	}
	return v1Route, nil
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
		Name:      virtualService.Name,
		Domains:   pointerify(virtualService.Domains),
		Routes:    convertOutputRoutes(virtualService.Routes),
		SslConfig: convertOutputSSLConfig(virtualService.SslConfig),
		Roles:     pointerify(virtualService.Roles),
		Status:    convertOutputStatus(virtualService.Status),
		Metadata:  convertOutputMetadata(virtualService.Metadata),
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

func convertDestinations(inputDests []InputWeightedDestination) ([]*v1.WeightedDestination, error) {
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
		UpstreamName: inputDest.UpstreamName,
		DestinationSpec: destSpec,
	}, nil
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

func convertInputSSLConfig(inSSL *InputSSLConfig) *v1.SSLConfig {
	if inSSL == nil {
		return nil
	}
	return &v1.SSLConfig{
		SecretRef: inSSL.SecretRef,
	}
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
