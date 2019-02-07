package graphql

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/util"
	extauthv1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	ratelimitapi "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"

	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/azure"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/grpc"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/rest"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/transformation"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	. "github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

type Converter struct {
	r   *ApiResolver
	ctx context.Context
}

func NewConverter(r *ApiResolver, ctx context.Context) *Converter {
	return &Converter{r: r, ctx: ctx}
}

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

	case spec.Grpc != nil:
		serviceArray := []*grpc.ServiceSpec_GrpcService{}
		for _, grpcValue := range spec.Grpc.GrpcServices {
			serviceArray = append(serviceArray, &grpc.ServiceSpec_GrpcService{
				PackageName:   grpcValue.PackageName,
				ServiceName:   grpcValue.ServiceName,
				FunctionNames: grpcValue.FunctionNames,
			})
		}

		return &plugins.ServiceSpec{PluginType: &plugins.ServiceSpec_Grpc{
			Grpc: &grpc.ServiceSpec{
				GrpcServices: serviceArray,
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

func (c *Converter) ConvertOutputUpstreams(upstreams v1.UpstreamList) ([]*Upstream, error) {
	var result []*Upstream
	for _, us := range upstreams {
		gqlUpstream, err := c.ConvertOutputUpstream(us)
		if err != nil {
			return nil, err
		}
		result = append(result, gqlUpstream)
	}
	return result, nil
}

func (c *Converter) ConvertOutputUpstream(upstream *v1.Upstream) (*Upstream, error) {
	usSpec, err := c.convertOutputUpstreamSpec(upstream.UpstreamSpec)
	if err != nil {
		return nil, err
	}

	return &Upstream{
		Spec:     usSpec,
		Metadata: convertOutputMetadata(&v1.Upstream{}, upstream.Metadata),
		Status:   convertOutputStatus(upstream.Status),
	}, nil
}

func (c *Converter) convertOutputUpstreamSpec(spec *v1.UpstreamSpec) (UpstreamSpec, error) {
	switch specType := spec.UpstreamType.(type) {
	case *v1.UpstreamSpec_Aws:
		return &AwsUpstreamSpec{
			Region:    specType.Aws.Region,
			SecretRef: convertOutputRef(specType.Aws.SecretRef),
			Functions: convertOutputLambdaFunctions(specType.Aws.LambdaFunctions),
		}, nil
	case *v1.UpstreamSpec_Azure:
		return &AzureUpstreamSpec{
			FunctionAppName: specType.Azure.FunctionAppName,
			Functions:       convertOutputAzureFunctions(specType.Azure.Functions),
		}, nil
	case *v1.UpstreamSpec_Kube:
		serviceSpec, err := c.convertOutputServiceSpec(specType.Kube.ServiceSpec)
		if err != nil {
			return nil, err
		}
		return &KubeUpstreamSpec{
			ServicePort:      int(specType.Kube.ServicePort),
			ServiceNamespace: specType.Kube.ServiceNamespace,
			ServiceName:      specType.Kube.ServiceName,
			Selector:         NewMapStringString(specType.Kube.Selector),
			ServiceSpec:      serviceSpec,
		}, nil
	case *v1.UpstreamSpec_Static:
		var hosts []StaticHost
		for _, h := range specType.Static.Hosts {
			hosts = append(hosts, StaticHost{
				Addr: h.Addr,
				Port: int(h.Port),
			})
		}
		serviceSpec, err := c.convertOutputServiceSpec(specType.Static.ServiceSpec)
		if err != nil {
			return nil, err
		}
		return &StaticUpstreamSpec{
			Hosts:       hosts,
			UseTLS:      specType.Static.UseTls,
			ServiceSpec: serviceSpec,
		}, nil
	}
	log.Printf("unsupported upstream type %v", spec)
	return nil, nil
}

// TODO (ilackarms): finish these methods
func (c *Converter) convertOutputServiceSpec(spec *plugins.ServiceSpec) (ServiceSpec, error) {
	if spec == nil {
		return nil, nil
	}
	switch serviceSpec := spec.PluginType.(type) {
	case *plugins.ServiceSpec_Rest:
		return &RestServiceSpec{
			Functions: convertOutputTransformations(serviceSpec.Rest.Transformations),
		}, nil
	case *plugins.ServiceSpec_Grpc:
		return &GrpcServiceSpec{
			GrpcServices: convertOutputGrpcServices(serviceSpec.Grpc.GrpcServices),
		}, nil
	}
	panic("unsupported")
}

func convertOutputGrpcServices(grpcServices []*grpc.ServiceSpec_GrpcService) []*GrpcService {
	if len(grpcServices) == 0 {
		return nil
	}
	var convertedGrpcServices []*GrpcService
	for _, svc := range grpcServices {
		convertedGrpcServices = append(convertedGrpcServices, &GrpcService{
			PackageName:   svc.PackageName,
			ServiceName:   svc.ServiceName,
			FunctionNames: svc.FunctionNames,
		})
	}
	return convertedGrpcServices
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

	virtualServicePlugins, err := convertInputVirtualServicePlugins(virtualService.RateLimitConfig, virtualService.ExtAuthConfig)
	if err != nil {
		return nil, errors.Wrap(err, "converting virtual service plugins")
	}

	return &gatewayv1.VirtualService{
		VirtualHost: &v1.VirtualHost{
			Domains:            virtualService.Domains,
			Routes:             routes,
			VirtualHostPlugins: virtualServicePlugins,
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
	case spec.Rest != nil:
		var params *transformation.Parameters
		if spec.Rest.Parameters != nil {
			headers := spec.Rest.Parameters.Headers.GoType()
			if len(headers) > 0 {
				if params == nil {
					params = &transformation.Parameters{}
				}
				params.Headers = headers
			}
			if inPath := spec.Rest.Parameters.Path; inPath != nil && *inPath != "" {
				if params == nil {
					params = &transformation.Parameters{}
				}
				params.Path = &types.StringValue{Value: *inPath}
			}
		}
		return &v1.DestinationSpec{
			DestinationType: &v1.DestinationSpec_Rest{
				Rest: &rest.DestinationSpec{
					FunctionName: spec.Rest.FunctionName,
					Parameters:   params,
				},
			},
		}, nil
	case spec.Grpc != nil:
		return &v1.DestinationSpec{
			DestinationType: &v1.DestinationSpec_Grpc{
				Grpc: &grpc.DestinationSpec{
					Package:  spec.Grpc.Package,
					Service:  spec.Grpc.Service,
					Function: spec.Grpc.Function,
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

func updateInputRateLimitSpec(in *InputRateLimit, out *ratelimitapi.IngressRateLimit, authorized bool) error {
	if out == nil {
		// Must not pass a nil pointer
		fmt.Errorf("Unable to create rate limit specification.")
	}
	rl := &ratelimitapi.RateLimit{}
	if in == nil {
		rl = nil
	} else {
		unit, err := convertGQLTimeUnitEnum(in.Unit)
		if err != nil {
			return err
		}
		rl = &ratelimitapi.RateLimit{
			Unit:            unit,
			RequestsPerUnit: uint32(in.RequestsPerUnit),
		}
	}
	if authorized {
		out.AuthorizedLimits = rl
	} else {
		out.AnonymousLimits = rl
	}
	return nil
}

func convertInputRateLimitToProto(inputConfig *InputRateLimitConfig) (*ratelimitapi.IngressRateLimit, error) {
	if inputConfig == nil {
		return nil, nil
	}
	rlProto := &ratelimitapi.IngressRateLimit{}

	if err := updateInputRateLimitSpec(inputConfig.AuthorizedLimits, rlProto, true); err != nil {
		return nil, err
	}
	if err := updateInputRateLimitSpec(inputConfig.AnonymousLimits, rlProto, false); err != nil {
		return nil, err
	}
	return rlProto, nil
}

func convertInputRateLimitToStruct(inputConfig *InputRateLimitConfig) (*types.Struct, error) {
	if inputConfig == nil {
		return nil, nil
	}
	rlProto, err := convertInputRateLimitToProto(inputConfig)
	if err != nil {
		return nil, err
	}
	return util.MessageToStruct(rlProto)
}

func convertInputExtAuthConfigToProto(inputConfig *InputExtAuthConfig) (*extauthv1.VhostExtension, error) {
	if inputConfig == nil {
		return nil, nil
	}
	if inputConfig.BasicAuth != nil {
		return nil, errors.Errorf("BasicAuth resolvers are not implemented yet.")
	}
	if inputConfig.OAuth == nil {
		return nil, nil
	}
	oidc := inputConfig.OAuth
	if oidc.AppURL == "" {
		return nil, errors.Errorf("invalid app url specified: %v", oidc.AppURL)
	}
	if oidc.IssuerURL == "" {
		return nil, errors.Errorf("invalid issuer url specified: %v", oidc.IssuerURL)
	}
	if oidc.ClientID == "" {
		return nil, errors.Errorf("invalid client id specified: %v", oidc.ClientID)
	}
	if oidc.CallbackPath == "" {
		return nil, errors.Errorf("invalid callback path specified: %v", oidc.CallbackPath)
	}
	if oidc.ClientSecretRef.Name == "" || oidc.ClientSecretRef.Namespace == "" {
		return nil, errors.Errorf("invalid client secret ref specified: %v.%v", oidc.ClientSecretRef.Namespace, oidc.ClientSecretRef.Name)
	}
	secretRef := convertInputRef(oidc.ClientSecretRef)
	eaProto := &extauthv1.VhostExtension{
		AuthConfig: &extauthv1.VhostExtension_Oauth{
			Oauth: &extauthv1.OAuth{
				AppUrl:          oidc.AppURL,
				CallbackPath:    oidc.CallbackPath,
				ClientId:        oidc.ClientID,
				ClientSecretRef: &secretRef,
				IssuerUrl:       oidc.IssuerURL,
			},
		},
	}
	return eaProto, nil
}

func convertInputExtAuthConfigToStruct(inputConfig *InputExtAuthConfig) (*types.Struct, error) {
	if inputConfig == nil {
		return nil, nil
	}
	rlProto, err := convertInputExtAuthConfigToProto(inputConfig)
	if err != nil {
		return nil, err
	}
	return util.MessageToStruct(rlProto)
}

func convertInputVirtualServicePlugins(inputRateLimitConfig *InputRateLimitConfig, inputExtAuthConfig *InputExtAuthConfig) (*v1.VirtualHostPlugins, error) {

	configs := map[string]*types.Struct{}

	rlStruct, err := convertInputRateLimitToStruct(inputRateLimitConfig)
	if err != nil {
		return nil, err
	}
	if rlStruct != nil {
		configs[ratelimit.ExtensionName] = rlStruct

	}

	eaStruct, err := convertInputExtAuthConfigToStruct(inputExtAuthConfig)
	if err != nil {
		return nil, err
	}
	if eaStruct != nil {
		configs[extauth.ExtensionName] = eaStruct

	}

	result := &v1.VirtualHostPlugins{
		Extensions: &v1.Extensions{
			Configs: configs,
		},
	}

	return result, nil
}

func (c *Converter) ConvertOutputVirtualServices(virtualServices gatewayv1.VirtualServiceList) ([]*VirtualService, error) {
	var result []*VirtualService
	for _, vs := range virtualServices {
		gqlVs, err := c.ConvertOutputVirtualService(vs)
		if err != nil {
			return nil, err
		}
		result = append(result, gqlVs)
	}
	return result, nil
}

func (c *Converter) ConvertOutputVirtualService(virtualService *gatewayv1.VirtualService) (*VirtualService, error) {
	gqlRoutes, err := c.ConvertOutputRoutes(virtualService.VirtualHost.Routes)
	if err != nil {
		return nil, err
	}
	extAuthConfig, err := convertOutputExtAuthConfig(virtualService.VirtualHost.VirtualHostPlugins)
	if err != nil {
		return nil, err
	}
	rateLimitConfig, err := convertOutputRateLimitConfig(virtualService.VirtualHost.VirtualHostPlugins)
	if err != nil {
		return nil, err
	}
	return &VirtualService{
		Domains:         virtualService.VirtualHost.Domains,
		Routes:          gqlRoutes,
		SslConfig:       convertOutputSSLConfig(virtualService.SslConfig),
		Status:          convertOutputStatus(virtualService.Status),
		Metadata:        convertOutputMetadata(&gatewayv1.VirtualService{}, virtualService.Metadata),
		ExtAuthConfig:   extAuthConfig,
		RateLimitConfig: rateLimitConfig,
	}, nil
}

func (c *Converter) ConvertOutputRoutes(routes []*v1.Route) ([]Route, error) {
	var outRoutes []Route
	for _, r := range routes {
		route, err := c.ConvertOutputRoute(r)
		if err != nil {
			return nil, err
		}
		outRoutes = append(outRoutes, route)
	}
	return outRoutes, nil
}

func (c *Converter) ConvertOutputRoute(route *v1.Route) (Route, error) {
	action, ok := route.Action.(*v1.Route_RouteAction)
	if !ok {
		return Route{}, errors.Errorf("%v does not have a RouteAction", route)
	}
	gqlDest, err := c.convertOutputDestination(action.RouteAction)
	if err != nil {
		return Route{}, err
	}
	return Route{
		Matcher:     convertOutputMatcher(route.Matcher),
		Destination: gqlDest,
		Plugins:     convertOutputRoutePlugins(route.RoutePlugins),
	}, nil
}

func (c *Converter) convertOutputDestination(action *v1.RouteAction) (Destination, error) {
	var outDest Destination
	switch dest := action.Destination.(type) {
	case *v1.RouteAction_Single:
		gqlDest, err := c.convertOutputSingleDestination(dest.Single)
		if err != nil {
			return nil, err
		}
		outDest = gqlDest
	case *v1.RouteAction_Multi:
		gqlDest, err := c.convertOutputMultiDestination(dest.Multi.Destinations)
		if err != nil {
			return nil, err
		}
		outDest = gqlDest
	}
	return outDest, nil
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

func (c *Converter) convertOutputMultiDestination(dests []*v1.WeightedDestination) (*MultiDestination, error) {
	var weightedDests []WeightedDestination
	for _, v1Dest := range dests {
		gqlDest, err := c.convertOutputSingleDestination(v1Dest.Destination)
		if err != nil {
			return nil, err
		}
		weightedDests = append(weightedDests, WeightedDestination{
			Destination: gqlDest,
			Weight:      int(v1Dest.Weight),
		})
	}
	return &MultiDestination{Destinations: weightedDests}, nil
}

func (c *Converter) convertOutputSingleDestination(dest *v1.Destination) (SingleDestination, error) {
	if dest.Upstream.Namespace == "" || dest.Upstream.Name == "" {
		return SingleDestination{}, errors.Errorf("must provide destination upstream")
	}
	gqlUs, err := c.r.Namespace().Upstream(c.ctx, &customtypes.Namespace{Name: dest.Upstream.Namespace}, dest.Upstream.Name)
	if err != nil {
		return SingleDestination{}, err
	}
	ds, err := c.convertOutputDestinationSpec(dest.DestinationSpec)
	if err != nil {
		return SingleDestination{}, err
	}
	return SingleDestination{
		Upstream:        *gqlUs,
		DestinationSpec: ds,
	}, nil
}

func (c *Converter) convertOutputDestinationSpec(spec *v1.DestinationSpec) (DestinationSpec, error) {
	if spec == nil {
		return nil, nil
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
			LogicalName:            destSpec.Aws.LogicalName,
			InvocationStyle:        invocationStyle,
			ResponseTransformation: destSpec.Aws.ResponseTrasnformation,
		}, nil
	case *v1.DestinationSpec_Azure:
		return &AzureDestinationSpec{
			FunctionName: destSpec.Azure.FunctionName,
		}, nil
	case *v1.DestinationSpec_Rest:
		return &RestDestinationSpec{
			FunctionName: destSpec.Rest.FunctionName,
			Parameters:   convertOutputTransformation(destSpec.Rest.Parameters),
		}, nil
	case *v1.DestinationSpec_Grpc:
		return &GrpcDestinationSpec{
			Package:    destSpec.Grpc.Package,
			Service:    destSpec.Grpc.Service,
			Function:   destSpec.Grpc.Function,
			Parameters: convertOutputTransformation(destSpec.Grpc.Parameters),
		}, nil
	}
	return nil, errors.Errorf("unknown destination spec type: %v", spec)
}

func convertOutputTransformation(params *transformation.Parameters) *TransformationParameters {
	if params == nil {
		return nil
	}
	var headers *MapStringString
	if len(params.Headers) > 0 {
		headers = NewMapStringString(params.Headers)
	}
	var path *string
	if params.Path != nil {
		path = &params.Path.Value
	}
	if path != nil || headers != nil {
		return &TransformationParameters{
			Path:    path,
			Headers: headers,
		}
	}
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

func convertOutputExtAuthConfig(plugins *v1.VirtualHostPlugins) (*ExtAuthConfig, error) {
	if plugins == nil {
		return nil, nil
	}
	if plugins.Extensions == nil {
		return nil, nil
	}
	var extAuth extauthv1.VhostExtension
	err := utils.UnmarshalExtension(plugins, extauth.ExtensionName, &extAuth)
	if err != nil {
		if err == utils.NotFoundError {
			// plugin not present, just return nil
			return nil, nil
		}
		// plugin present and marshalling failed, return error
		return nil, errors.Wrapf(err, "failed to unmarshal proto message to %v plugin", extauth.ExtensionName)
	}

	result := &ExtAuthConfig{}
	switch auth := extAuth.AuthConfig.(type) {
	case *extauthv1.VhostExtension_BasicAuth:
		result.AuthType = &BasicAuth{
			Realm: auth.BasicAuth.Realm,
		}

	case *extauthv1.VhostExtension_Oauth:
		secretKey := auth.Oauth.ClientSecretRef.Key()
		result.AuthType = &OAuthConfig{
			ClientID:     &auth.Oauth.ClientId,
			ClientSecret: &secretKey,
			IssuerURL:    &auth.Oauth.IssuerUrl,
			AppURL:       &auth.Oauth.AppUrl,
			CallbackPath: &auth.Oauth.CallbackPath,
		}
	default:
		return nil, errors.Errorf("Unrecognized auth type %v", auth)
	}

	return result, nil
}

func convertOutputRateLimitConfig(plugins *v1.VirtualHostPlugins) (*RateLimitConfig, error) {
	if plugins == nil {
		return nil, nil
	}
	if plugins.Extensions == nil {
		return nil, nil
	}
	if rl, ok := plugins.Extensions.Configs[ratelimit.ExtensionName]; !ok || rl == nil {
		return nil, nil
	}

	var rateLimit ratelimitapi.IngressRateLimit
	err := utils.UnmarshalExtension(plugins, ratelimit.ExtensionName, &rateLimit)
	if err != nil {
		if err == utils.NotFoundError {
			// plugin not present, just return nil
			return nil, nil
		}
		// plugin present and marshalling failed, return error
		return nil, errors.Wrapf(err, "failed to unmarshal proto message to %v plugin", ratelimit.ExtensionName)
	}

	result := &RateLimitConfig{}
	result.AuthorizedLimits, err = convertOutputRateLimit(rateLimit.AuthorizedLimits)
	if err != nil {
		return nil, err
	}
	result.AnonymousLimits, err = convertOutputRateLimit(rateLimit.AnonymousLimits)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func convertOutputRateLimit(proto *ratelimitapi.RateLimit) (*RateLimit, error) {
	// If the proto representation is nil, we have nothing to do
	// RateLimit_UNKNOWN is the "zero value" for initialized RateLimits.
	// If present here, it indicates an inactive rate limit
	if proto == nil || proto.Unit == ratelimitapi.RateLimit_UNKNOWN {
		return nil, nil
	}
	authUnit, err := convertProtoTimeUnitEnum(proto.Unit)
	if err != nil {
		return nil, err
	}
	return &RateLimit{
		Unit:            authUnit,
		RequestsPerUnit: customtypes.UnsignedInt(proto.RequestsPerUnit),
	}, nil
}

func convertProtoTimeUnitEnum(protoEnum ratelimitapi.RateLimit_Unit) (TimeUnit, error) {
	switch protoEnum {
	case ratelimitapi.RateLimit_SECOND:
		return TimeUnitSecond, nil
	case ratelimitapi.RateLimit_MINUTE:
		return TimeUnitMinute, nil
	case ratelimitapi.RateLimit_HOUR:
		return TimeUnitHour, nil
	case ratelimitapi.RateLimit_DAY:
		return TimeUnitDay, nil
	default:
		return "", errors.Errorf("invalid rate limit unit: %v", protoEnum)
	}
}

func convertGQLTimeUnitEnum(gqlEnum TimeUnit) (ratelimitapi.RateLimit_Unit, error) {
	switch gqlEnum {
	case TimeUnitSecond:
		return ratelimitapi.RateLimit_SECOND, nil
	case TimeUnitMinute:
		return ratelimitapi.RateLimit_MINUTE, nil
	case TimeUnitHour:
		return ratelimitapi.RateLimit_HOUR, nil
	case TimeUnitDay:
		return ratelimitapi.RateLimit_DAY, nil
	default:
		return 0, errors.Errorf("invalid rate limit time unit: %v", gqlEnum)
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

func convertOutputMetadata(resource resources.Resource, meta core.Metadata) Metadata {
	resource = resources.Clone(resource)
	resource.SetMetadata(meta)
	return Metadata{
		GUID:            resources.Key(resource),
		Namespace:       meta.Namespace,
		Name:            meta.Name,
		ResourceVersion: meta.ResourceVersion,
		Labels:          NewMapStringString(meta.Labels),
		Annotations:     NewMapStringString(meta.Annotations),
	}
}

func (c *Converter) ConvertOutputSecrets(secrets v1.SecretList) []*Secret {
	var result []*Secret
	for _, us := range secrets {
		result = append(result, c.ConvertOutputSecret(us))
	}
	return result
}

func (c *Converter) ConvertOutputSecret(secret *v1.Secret) *Secret {
	return convertOutputSecret(secret)
}

// We are not actually returning the secret values, just their names
// Consider updating the graphql schema to only return Metadata
func convertOutputSecret(secret *v1.Secret) *Secret {
	out := &Secret{
		Metadata: convertOutputMetadata(&v1.Secret{}, secret.Metadata),
	}
	switch secret.Kind.(type) {
	case *v1.Secret_Aws:
		out.Kind = &AwsSecret{
			// AccessKey: sec.Aws.AccessKey,
			// SecretKey: sec.Aws.SecretKey,
			AccessKey: "(access key)",
			SecretKey: "(secret key)",
		}
	case *v1.Secret_Azure:
		out.Kind = &AzureSecret{
			// APIKeys: NewMapStringString(sec.Azure.ApiKeys),
			APIKeys: NewMapStringString(map[string]string{"secretKey": "secretValue"}),
		}
	case *v1.Secret_Tls:
		out.Kind = &TlsSecret{
			// CertChain:  sec.Tls.CertChain,
			// RootCa:     sec.Tls.RootCa,
			// PrivateKey: sec.Tls.PrivateKey,
			CertChain:  "(cert chain)",
			RootCa:     "(root ca)",
			PrivateKey: "(private key)",
		}
	case *v1.Secret_Extension:
		// Currently we only have one type of secret extension so we don't need to do a check
		out.Kind = &OauthSecret{
			ClientSecret: "n/a",
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
	case secret.Kind.Oauth != nil:
		secretStruct, err := util.MessageToStruct(&extauthv1.OauthSecret{
			ClientSecret: secret.Kind.Oauth.ClientSecret,
		})
		if err != nil {
			return nil, err
		}
		out.Kind = &v1.Secret_Extension{
			Extension: &v1.Extension{
				Config: secretStruct,
			},
		}
	default:
		return nil, errors.Errorf("invalid input secret:  requires one of Aws, Azure, TLS, or OAuth set")
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
		Metadata: convertOutputMetadata(&v1.Artifact{}, artifact.Metadata),
	}
}

func (c *Converter) ConvertOutputSettings(settings *v1.Settings) *Settings {
	refreshRate, err := types.DurationFromProto(settings.RefreshRate)
	if err != nil {
		log.Printf("weird error trying to convert duration from proto: %v", err)
	}
	dur := customtypes.Duration(refreshRate)
	return &Settings{
		WatchNamespaces: settings.WatchNamespaces,
		RefreshRate:     &dur,
		Metadata:        convertOutputMetadata(&v1.Settings{}, settings.Metadata),
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

func (c *Converter) ConvertInputSettings(settings InputSettings) (*v1.Settings, error) {
	var refreshRate *types.Duration
	if settings.RefreshRate != nil {
		refreshRate = types.DurationProto(time.Duration(*settings.RefreshRate))
	}
	return &v1.Settings{
		WatchNamespaces: settings.WatchNamespaces,
		RefreshRate:     refreshRate,
		Metadata:        convertInputMetadata(settings.Metadata),
	}, nil
}
