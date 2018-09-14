package graphql

import (
	"log"
	"sort"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/rest"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/static"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/transformation"
)

type Converter struct{}

func (c *Converter) ConvertInputUpstreams(upstream []InputUpstream) (v1.UpstreamList, error) {
	var result v1.UpstreamList
	for _, us := range upstream {
		converted, err := c.ConvertInputUpstream(us)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

func convertInputRef(ref InputResourceRef) core.ResourceRef {
	return core.ResourceRef{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}
}

func convertOutputRef(ref core.ResourceRef) ResourceRef {
	return ResourceRef{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}
}

func (c *Converter) ConvertInputUpstream(upstream InputUpstream) (*v1.Upstream, error) {
	upstreamSpec, err := convertInputUpstreamSpec(upstream.Spec)
	if err != nil {
		return nil, err
	}
	return &v1.Upstream{
		Metadata:     convertInputMetadata(upstream.Metadata),
		UpstreamSpec: upstreamSpec,
	}, nil
}

func convertInputUpstreamSpec(spec InputUpstreamSpec) (*v1.UpstreamSpec, error) {
	out := &v1.UpstreamSpec{}
	switch {
	case spec.Aws != nil:
		out.UpstreamType = &v1.UpstreamSpec_Aws{
			Aws: &aws.UpstreamSpec{
				Region:          spec.Aws.Region,
				SecretRef:       convertInputRef(spec.Aws.SecretRef),
				LambdaFunctions: convertInputLambdaFunctions(spec.Aws.Functions),
			},
		}
	case spec.Azure != nil:
		var ref core.ResourceRef
		if spec.Azure.SecretRef != nil {
			ref = convertInputRef(*spec.Azure.SecretRef)
		}
		out.UpstreamType = &v1.UpstreamSpec_Azure{
			Azure: &azure.UpstreamSpec{
				FunctionAppName: spec.Azure.FunctionAppName,
				SecretRef:       ref,
				Functions:       convertInputAzureFunctions(spec.Azure.Functions),
			},
		}
	case spec.Kube != nil:
		if err := spec.Kube.Selector.Validate(); err != nil {
			return nil, errors.Wrapf(err, "invalid spec")
		}
		out.UpstreamType = &v1.UpstreamSpec_Kube{
			Kube: &kubernetes.UpstreamSpec{
				Selector:         spec.Kube.Selector.GoType(),
				ServiceName:      spec.Kube.ServiceName,
				ServiceNamespace: spec.Kube.ServiceNamespace,
				ServicePort:      uint32(spec.Kube.ServicePort),
			},
		}
	case spec.Static != nil:
		serviceSpec, err := convertInputServiceSpec(spec.Static.ServiceSpec)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid service spec")
		}
		var hosts []*static.Host
		for _, h := range spec.Static.Hosts {
			hosts = append(hosts, &static.Host{
				Addr: h.Addr,
				Port: uint32(h.Port),
			})
		}
		out.UpstreamType = &v1.UpstreamSpec_Static{
			Static: &static.UpstreamSpec{
				Hosts:       hosts,
				UseTls:      spec.Static.UseTLS,
				ServiceSpec: serviceSpec,
			},
		}
	default:
		log.Printf("invalid spec: %#v", spec)
	}
	return out, nil
}

func convertInputLambdaFunctions(inputFuncs []InputAwsLambdaFunction) []*aws.LambdaFunctionSpec {
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

func convertInputAzureFunctions(inputFuncs []InputAzureFunction) []*azure.UpstreamSpec_FunctionSpec {
	var funcs []*azure.UpstreamSpec_FunctionSpec
	for _, inFn := range inputFuncs {
		var authLevel azure.UpstreamSpec_FunctionSpec_AuthLevel
		switch AzureFnAuthLevel(inFn.AuthLevel) {
		case AzureFnAuthLevelAnonymous:
			authLevel = azure.UpstreamSpec_FunctionSpec_Anonymous
		case AzureFnAuthLevelAdmin:
			authLevel = azure.UpstreamSpec_FunctionSpec_Admin
		case AzureFnAuthLevelFunction:
			authLevel = azure.UpstreamSpec_FunctionSpec_Function
		}
		funcs = append(funcs, &azure.UpstreamSpec_FunctionSpec{
			FunctionName: inFn.FunctionName,
			AuthLevel:    authLevel,
		})
	}
	return funcs
}

// TODO (ilackarms): finish these methods
func convertInputServiceSpec(spec *InputServiceSpec) (*plugins.ServiceSpec, error) {
	if spec == nil {
		return nil, nil
	}
	switch {
	case spec.Rest != nil:
		var swaggerInfo *rest.ServiceSpec_SwaggerInfo
		if spec.Rest.InlineSwaggerDoc != nil {
			swaggerInfo = &rest.ServiceSpec_SwaggerInfo{
				SwaggerSpec: &rest.ServiceSpec_SwaggerInfo_Inline{Inline: *spec.Rest.InlineSwaggerDoc},
			}
		}
		return &plugins.ServiceSpec{PluginType: &plugins.ServiceSpec_Rest{
			Rest: &rest.ServiceSpec{
				Transformations: convertInputTransformations(spec.Rest.Functions),
				SwaggerInfo:     swaggerInfo,
			},
		}}, nil
	}
	return nil, errors.Errorf("unsupported spec: %v", spec)
}

func convertInputTransformations(in []InputTransformation) map[string]*transformation.TransformationTemplate {
	transforms := make(map[string]*transformation.TransformationTemplate)
	for _, trans := range in {
		glooTransformation := &transformation.TransformationTemplate{}
		if trans.Body != nil {
			glooTransformation.BodyTransformation = &transformation.TransformationTemplate_Body{
				Body: injaTemplateFromString(*trans.Body),
			}
		}
		if headers := trans.Headers.GoType(); len(headers) > 0 {
			glooHeaders := make(map[string]*transformation.InjaTemplate)
			for k, v := range headers {
				glooHeaders[k] = injaTemplateFromString(v)
			}
			glooTransformation.Headers = glooHeaders
		}
		transforms[trans.FunctionName] = glooTransformation
	}
	return transforms
}

func injaTemplateFromString(str string) *transformation.InjaTemplate {
	return &transformation.InjaTemplate{
		Text: str,
	}
}

func (c *Converter) ConvertOutputUpstreams(upstreams v1.UpstreamList) []*Upstream {
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
			SecretRef: convertOutputRef(specType.Aws.SecretRef),
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
			Selector:         NewMapStringString(specType.Kube.Selector),
			ServiceSpec:      convertOutputServiceSpec(specType.Kube.ServiceSpec),
		}
	case *v1.UpstreamSpec_Static:
		var hosts []StaticHost
		for _, h := range specType.Static.Hosts {
			hosts = append(hosts, StaticHost{
				Addr: h.Addr,
				Port: int(h.Port),
			})
		}
		return &StaticUpstreamSpec{
			Hosts:       hosts,
			UseTLS:      specType.Static.UseTls,
			ServiceSpec: convertOutputServiceSpec(specType.Static.ServiceSpec),
		}
	}
	log.Printf("unsupported upstream type %v", spec)
	return nil
}

// TODO (ilackarms): finish these methods
func convertOutputServiceSpec(spec *plugins.ServiceSpec) ServiceSpec {
	if spec == nil {
		return nil
	}
	switch serviceSpec := spec.PluginType.(type) {
	case *plugins.ServiceSpec_Rest:
		return &RestServiceSpec{
			Functions: convertOutputTransformations(serviceSpec.Rest.Transformations),
		}
	}
	panic("unsupported")
}

func convertOutputTransformations(transformations map[string]*transformation.TransformationTemplate) []Transformation {
	var transforms []Transformation
	for fnName, trans := range transformations {
		var body *string
		if trans.BodyTransformation != nil {
			bodyTransform, ok := trans.BodyTransformation.(*transformation.TransformationTemplate_Body)
			if ok && bodyTransform.Body != nil {
				body = &bodyTransform.Body.Text
			}
		}
		var headers *MapStringString
		if len(trans.Headers) > 0 {
			h := make(map[string]string)
			for k, v := range trans.Headers {
				if v == nil {
					continue
				}
				h[k] = v.Text
			}
			headers = NewMapStringString(h)
		}
		transforms = append(transforms, Transformation{
			FunctionName: fnName,
			Body:         body,
			Headers:      headers,
		})
	}
	return transforms
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
		var authLevel AzureFnAuthLevel
		switch l.AuthLevel {
		case azure.UpstreamSpec_FunctionSpec_Anonymous:
			authLevel = AzureFnAuthLevelAnonymous
		case azure.UpstreamSpec_FunctionSpec_Admin:
			authLevel = AzureFnAuthLevelAdmin
		case azure.UpstreamSpec_FunctionSpec_Function:
			authLevel = AzureFnAuthLevelFunction
		}
		out = append(out, AzureFunction{
			FunctionName: l.FunctionName,
			AuthLevel:    authLevel,
		})
	}
	return out
}

func (c *Converter) ConvertInputVirtualServices(virtualService []InputVirtualService) (gatewayv1.VirtualServiceList, error) {
	var result gatewayv1.VirtualServiceList
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
	action, err := convertInputDestinationToAction(route.Destination)
	if err != nil {
		return nil, err
	}
	return &v1.Route{
		Matcher:      match,
		RoutePlugins: convertInputRoutePlugins(route.Plugins),
		Action: &v1.Route_RouteAction{
			RouteAction: action,
		},
	}, nil
}

func convertInputDestinationToAction(dest InputDestination) (*v1.RouteAction, error) {
	action := &v1.RouteAction{}
	switch {
	case dest.SingleDestination != nil:
		dest, err := convertInputSingleDestination(*dest.SingleDestination)
		if err != nil {
			return nil, err
		}
		action.Destination = &v1.RouteAction_Single{
			Single: dest,
		}
	case dest.MultiDestination != nil:
		weightedDestinations, err := convertInputDestinations(dest.MultiDestination.Destinations)
		if err != nil {
			return nil, err
		}
		action.Destination = &v1.RouteAction_Multi{
			Multi: &v1.MultiDestination{
				Destinations: weightedDestinations,
			},
		}
	default:
		return nil, errors.Errorf("must specify exactly one of SingleDestination or MultiDestinations")
	}
	return action, nil
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
	// TODO(ilackaitems): convert route plugins when there are any
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
	var invocationstyle aws.DestinationSpec_InvocationStyle
	switch {
	case spec.Rest != nil:
		return &v1.DestinationSpec{DestinationType: &v1.DestinationSpec_Rest{
			Rest: &rest.DestinationSpec{
				FunctionName: spec.Rest.FunctionName,
			},
		}}, nil
	case spec.Aws != nil:
		switch spec.Aws.InvocationStyle {
		case AwsLambdaInvocationStyleAsync:
			invocationstyle = aws.DestinationSpec_ASYNC
		case AwsLambdaInvocationStyleSync:
			invocationstyle = aws.DestinationSpec_SYNC
		}
		return &v1.DestinationSpec{
			DestinationType: &v1.DestinationSpec_Aws{
				Aws: &aws.DestinationSpec{
					LogicalName:            spec.Aws.LogicalName,
					InvocationStyle:        invocationstyle,
					ResponseTrasnformation: spec.Aws.ResponseTransformation,
				},
			},
		}, nil
	case spec.Azure != nil:
		return &v1.DestinationSpec{
			DestinationType: &v1.DestinationSpec_Azure{
				Azure: &azure.DestinationSpec{
					FunctionName: spec.Azure.FunctionName,
				},
			},
		}, nil
	}
	return nil, nil
}

func convertInputSingleDestination(inputDest InputSingleDestination) (*v1.Destination, error) {
	destSpec, err := convertInputDestinationSpec(inputDest.DestinationSpec)
	if err != nil {
		return nil, err
	}
	return &v1.Destination{
		Upstream:        convertInputRef(inputDest.Upstream),
		DestinationSpec: destSpec,
	}, nil
}

func convertInputSSLConfig(ssl *InputSslConfig) *v1.SslConfig {
	if ssl == nil {
		return nil
	}
	ref := convertInputRef(ssl.SecretRef)
	return &v1.SslConfig{
		SslSecrets: &v1.SslConfig_SecretRef{
			SecretRef: &ref,
		},
	}
}

func (c *Converter) ConvertOutputVirtualServices(virtualServices gatewayv1.VirtualServiceList) []*VirtualService {
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
	return Route{
		Matcher:     convertOutputMatcher(route.Matcher),
		Destination: convertOutputDestination(action.RouteAction),
		Plugins:     convertOutputRoutePlugins(route.RoutePlugins),
	}, true
}

func convertOutputDestination(action *v1.RouteAction) Destination {
	var outDest Destination
	switch dest := action.Destination.(type) {
	case *v1.RouteAction_Single:
		outDest = convertOutputSingleDestination(dest.Single)
	case *v1.RouteAction_Multi:
		outDest = convertOutputMultiDestination(dest.Multi.Destinations)
	}
	return outDest
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
	// TODO(ilackaitems): convert route plugins when there are any
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
		Upstream:        convertOutputRef(dest.Upstream),
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
	case *v1.DestinationSpec_Azure:
		return &AzureDestinationSpec{
			FunctionName: destSpec.Azure.FunctionName,
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

	var ref ResourceRef
	if secret.SecretRef != nil {
		ref = convertOutputRef(*secret.SecretRef)
	}

	return &SslConfig{
		SecretRef: ref,
	}
}

func (c *Converter) ConvertOutputResolverMaps(resolverMaps sqoopv1.ResolverMapList) []*ResolverMap {
	var result []*ResolverMap
	for _, us := range resolverMaps {
		result = append(result, c.ConvertOutputResolverMap(us))
	}
	return result
}

func (c *Converter) ConvertOutputResolverMap(resolverMap *sqoopv1.ResolverMap) *ResolverMap {
	var typeResolvers []TypeResolver
	for typeName, typeResolver := range resolverMap.Types {
		typeResolvers = append(typeResolvers, convertOutputTypeResolver(typeName, typeResolver))
	}
	sort.SliceStable(typeResolvers, func(i, j int) bool {
		return typeResolvers[i].TypeName < typeResolvers[j].TypeName
	})
	return &ResolverMap{
		Types:    typeResolvers,
		Status:   convertOutputStatus(resolverMap.Status),
		Metadata: convertOutputMetadata(resolverMap.Metadata),
	}
}

func convertOutputTypeResolver(typeName string, typeResolver *sqoopv1.TypeResolver) TypeResolver {
	var fieldResolvers []FieldResolver
	for fieldName, fieldResolver := range typeResolver.Fields {
		fieldResolvers = append(fieldResolvers, FieldResolver{
			FieldName: fieldName,
			Resolver:  convertOutputResolver(fieldResolver),
		})
	}
	sort.SliceStable(fieldResolvers, func(i, j int) bool {
		return fieldResolvers[i].FieldName < fieldResolvers[j].FieldName
	})
	return TypeResolver{
		TypeName: typeName,
		Fields:   fieldResolvers,
	}
}

func convertOutputResolver(resolver *sqoopv1.FieldResolver) Resolver {
	switch res := resolver.Resolver.(type) {
	case *sqoopv1.FieldResolver_GlooResolver:
		// Until implemented - bypass. TODO -implement
		if res.GlooResolver == nil {
			return nil
		}
		return &GlooResolver{
			RequestTemplate:  convertOutputRequestTemplate(res.GlooResolver.RequestTemplate),
			ResponseTemplate: convertOutputResponseTemplate(res.GlooResolver.ResponseTemplate),
			Destination:      convertOutputDestination(res.GlooResolver.Action),
		}
	case *sqoopv1.FieldResolver_TemplateResolver:
		return &TemplateResolver{}
	case *sqoopv1.FieldResolver_NodejsResolver:
		return &NodeJSResolver{}
	}
	log.Printf("invalid resolver type: %v", resolver)
	return nil
}

func convertOutputRequestTemplate(t *sqoopv1.RequestTemplate) *RequestTemplate {
	if t == nil {
		return nil
	}
	return &RequestTemplate{
		Verb:    t.Verb,
		Path:    t.Path,
		Body:    t.Body,
		Headers: NewMapStringString(t.Headers),
	}
}

func convertOutputResponseTemplate(t *sqoopv1.ResponseTemplate) *ResponseTemplate {
	if t == nil {
		return nil
	}
	return &ResponseTemplate{
		Body:    t.Body,
		Headers: NewMapStringString(t.Headers),
	}
}

func (c *Converter) ConvertInputResolverMaps(resolverMaps []*InputResolverMap) (sqoopv1.ResolverMapList, error) {
	var result sqoopv1.ResolverMapList
	for _, item := range resolverMaps {
		in, err := c.ConvertInputResolverMap(*item)
		if err != nil {
			return nil, err
		}
		result = append(result, in)
	}
	return result, nil
}

func (c *Converter) ConvertInputResolverMap(resolverMap InputResolverMap) (*sqoopv1.ResolverMap, error) {
	typeResolvers := make(map[string]*sqoopv1.TypeResolver)
	for _, typeResolver := range resolverMap.Types {
		res, err := convertInputTypeResolver(typeResolver)
		if err != nil {
			return nil, err
		}
		typeResolvers[typeResolver.TypeName] = res
	}
	return &sqoopv1.ResolverMap{
		Metadata: convertInputMetadata(resolverMap.Metadata),
		Types:    typeResolvers,
	}, nil
}

func convertInputTypeResolver(typeResolver InputTypeResolver) (*sqoopv1.TypeResolver, error) {
	fieldResolvers := make(map[string]*sqoopv1.FieldResolver)
	for _, fieldResolver := range typeResolver.Fields {
		resolver, err := ConvertInputResolver(fieldResolver.Resolver)
		if err != nil {
			return nil, err
		}
		fieldResolvers[fieldResolver.FieldName] = resolver
	}
	return &sqoopv1.TypeResolver{
		Fields: fieldResolvers,
	}, nil
}

// TODO(ilacakitems): implement these
func ConvertInputResolver(resolver InputResolver) (*sqoopv1.FieldResolver, error) {
	switch {
	case resolver.GlooResolver != nil:
		action, err := convertInputDestinationToAction(resolver.GlooResolver.Destination)
		if err != nil {
			return nil, err
		}
		return &sqoopv1.FieldResolver{
			Resolver: &sqoopv1.FieldResolver_GlooResolver{
				GlooResolver: &sqoopv1.GlooResolver{
					RequestTemplate:  convertInputRequestTemplate(resolver.GlooResolver.RequestTemplate),
					ResponseTemplate: convertInputResponseTemplate(resolver.GlooResolver.ResponseTemplate),
					Action:           action,
				},
			},
		}, nil
	case resolver.TemplateResolver != nil:
		return &sqoopv1.FieldResolver{
			Resolver: &sqoopv1.FieldResolver_TemplateResolver{
				TemplateResolver: &sqoopv1.TemplateResolver{},
			},
		}, nil
	case resolver.NodeResolver != nil:
		return &sqoopv1.FieldResolver{
			Resolver: &sqoopv1.FieldResolver_NodejsResolver{
				NodejsResolver: &sqoopv1.NodeJSResolver{},
			},
		}, nil
	}
	return nil, errors.Errorf("invalid input resolver: %#v", resolver)
}

func convertInputRequestTemplate(t *InputRequestTemplate) *sqoopv1.RequestTemplate {
	if t == nil {
		return nil
	}
	return &sqoopv1.RequestTemplate{
		Verb:    t.Verb,
		Path:    t.Path,
		Body:    t.Body,
		Headers: t.Headers.GoType(),
	}
}

func convertInputResponseTemplate(t *InputResponseTemplate) *sqoopv1.ResponseTemplate {
	if t == nil {
		return nil
	}
	return &sqoopv1.ResponseTemplate{
		Body:    t.Body,
		Headers: t.Headers.GoType(),
	}
}

// common
func convertInputMetadata(inMeta InputMetadata) core.Metadata {
	return core.Metadata{
		Namespace:       inMeta.Namespace,
		Name:            inMeta.Name,
		ResourceVersion: inMeta.ResourceVersion,
		Labels:          inMeta.Labels.GoType(),
		Annotations:     inMeta.Annotations.GoType(),
	}
}

func convertOutputMetadata(meta core.Metadata) Metadata {
	return Metadata{
		Namespace:       meta.Namespace,
		Name:            meta.Name,
		ResourceVersion: meta.ResourceVersion,
		Labels:          NewMapStringString(meta.Labels),
		Annotations:     NewMapStringString(meta.Annotations),
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

func (c *Converter) ConvertOutputSchemas(schemas sqoopv1.SchemaList) []*Schema {
	var result []*Schema
	for _, us := range schemas {
		result = append(result, c.ConvertOutputSchema(us))
	}
	return result
}

func (c *Converter) ConvertOutputSchema(schema *sqoopv1.Schema) *Schema {
	return &Schema{
		InlineSchema: schema.InlineSchema,
		Status:       convertOutputStatus(schema.Status),
		Metadata:     convertOutputMetadata(schema.Metadata),
	}
}

func (c *Converter) ConvertInputSchemas(schemas []*InputSchema) (sqoopv1.SchemaList, error) {
	var result sqoopv1.SchemaList
	for _, item := range schemas {
		in, err := c.ConvertInputSchema(*item)
		if err != nil {
			return nil, err
		}
		result = append(result, in)
	}
	return result, nil
}

func (c *Converter) ConvertInputSchema(schema InputSchema) (*sqoopv1.Schema, error) {
	return &sqoopv1.Schema{
		Metadata:     convertInputMetadata(schema.Metadata),
		InlineSchema: schema.InlineSchema,
	}, nil
}

func (c *Converter) ConvertOutputSecrets(secrets v1.SecretList) []*Secret {
	var result []*Secret
	for _, us := range secrets {
		result = append(result, c.ConvertOutputSecret(us))
	}
	return result
}

func (c *Converter) ConvertOutputSecret(secret *v1.Secret) *Secret {
	out := &Secret{
		Metadata: convertOutputMetadata(secret.Metadata),
	}
	switch sec := secret.Kind.(type) {
	case *v1.Secret_Aws:
		out.Kind = &AwsSecret{
			AccessKey: sec.Aws.AccessKey,
			SecretKey: sec.Aws.SecretKey,
		}
	case *v1.Secret_Azure:
		out.Kind = &AzureSecret{
			APIKeys: NewMapStringString(sec.Azure.ApiKeys),
		}
	case *v1.Secret_Tls:
		out.Kind = &TlsSecret{
			CertChain:  sec.Tls.CertChain,
			RootCa:     sec.Tls.RootCa,
			PrivateKey: sec.Tls.PrivateKey,
		}
	}
	return out
}

func (c *Converter) ConvertInputSecrets(secrets []*InputSecret) (v1.SecretList, error) {
	var result v1.SecretList
	for _, item := range secrets {
		in, err := c.ConvertInputSecret(*item)
		if err != nil {
			return nil, err
		}
		result = append(result, in)
	}
	return result, nil
}

func (c *Converter) ConvertInputSecret(secret InputSecret) (*v1.Secret, error) {
	out := &v1.Secret{
		Metadata: convertInputMetadata(secret.Metadata),
	}
	switch {
	case secret.Kind.Aws != nil:
		out.Kind = &v1.Secret_Aws{
			Aws: &v1.AwsSecret{
				AccessKey: secret.Kind.Aws.AccessKey,
				SecretKey: secret.Kind.Aws.SecretKey,
			},
		}
	case secret.Kind.Azure != nil:
		out.Kind = &v1.Secret_Azure{
			Azure: &v1.AzureSecret{
				ApiKeys: secret.Kind.Azure.APIKeys.GoType(),
			},
		}
	case secret.Kind.TLS != nil:
		out.Kind = &v1.Secret_Tls{
			Tls: &v1.TlsSecret{
				PrivateKey: secret.Kind.TLS.PrivateKey,
				RootCa:     secret.Kind.TLS.RootCa,
				CertChain:  secret.Kind.TLS.CertChain,
			},
		}
	default:
		return nil, errors.Errorf("invalid input secret:  requires one of Aws, Azure, or TLS set")
	}
	return out, nil
}

func (c *Converter) ConvertOutputArtifacts(artifacts v1.ArtifactList) []*Artifact {
	var result []*Artifact
	for _, us := range artifacts {
		result = append(result, c.ConvertOutputArtifact(us))
	}
	return result
}

func (c *Converter) ConvertOutputArtifact(artifact *v1.Artifact) *Artifact {
	return &Artifact{
		Metadata: convertOutputMetadata(artifact.Metadata),
	}
}

func (c *Converter) ConvertInputArtifacts(artifacts []*InputArtifact) (v1.ArtifactList, error) {
	var result v1.ArtifactList
	for _, item := range artifacts {
		in, err := c.ConvertInputArtifact(*item)
		if err != nil {
			return nil, err
		}
		result = append(result, in)
	}
	return result, nil
}

func (c *Converter) ConvertInputArtifact(artifact InputArtifact) (*v1.Artifact, error) {
	return &v1.Artifact{
		Metadata: convertInputMetadata(artifact.Metadata),
		Data:     artifact.Data,
	}, nil
}
