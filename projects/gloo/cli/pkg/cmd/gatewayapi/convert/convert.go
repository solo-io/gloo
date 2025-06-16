package convert

import (
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ai"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/types"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"

	"encoding/json"

	kgateway "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	gloogateway "github.com/solo-io/gloo-gateway/api/v1alpha1"
	gloogwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

func (o *GatewayAPIOutput) Convert() error {

	for _, gateway := range o.edgeCache.GlooGateways() {
		// We only translate virtual services for ones that match a gateway selector
		// TODO in the future we could blindly convert VS and not attach it to anything
		err := o.convertGatewayAndVirtualServices(gateway)
		if err != nil {
			return err
		}
	}

	for _, routeTable := range o.edgeCache.RouteTables() {
		err := o.convertRouteTableToHTTPRoute(routeTable)
		if err != nil {
			return err
		}
	}

	for _, upstream := range o.edgeCache.Upstreams() {
		// Add all existing upstreams except for kube services which will be referenced directly
		if upstream.Spec.GetKube() == nil {
			o.gatewayAPICache.AddBackend(convertUpstreamToBackend(upstream))
		}
	}

	//TODO(nick): what do we do with settings?
	// for _, settings := range o.edgeCache.Settings() {
	// 	o.gatewayAPICache.AddSettings(settings)
	// }

	// copy over any existing options
	return nil
}

// TODO(nick): this is a placeholder for now, we need to figure out how to convert the upstream to a backend
func convertUpstreamToBackend(upstream *snapshot.UpstreamWrapper) *snapshot.BackendWrapper {
	backend := &snapshot.BackendWrapper{
		Backend: &kgateway.Backend{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Backend",
				APIVersion: kgateway.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      upstream.GetName(),
				Namespace: upstream.GetNamespace(),
			},
			Spec: kgateway.BackendSpec{
				// Static: &kgateway.StaticBackend{
				// },
			},
		},
	}
	return backend
}

func (o *GatewayAPIOutput) convertGatewayAndVirtualServices(glooGateway *snapshot.GlooGatewayWrapper) error {

	// we first need to generate Gateway objects with the correct names based on proxy Names
	// spec.proxyNames
	o.generateGatewaysFromProxyNames(glooGateway)

	gatewayVs, err := o.edgeCache.GlooGatewayVirtualServices(glooGateway)
	if err != nil {
		return err
	}
	if len(gatewayVs) == 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NO_REFERENCES, glooGateway, "gateway does not contain virtual services")
	}
	for _, vs := range gatewayVs {
		proxyNames := glooGateway.Spec.GetProxyNames()
		if len(proxyNames) == 0 {
			proxyNames = append(proxyNames, "gateway-proxy")
		}
		for _, proxyName := range proxyNames {
			listenerName := fmt.Sprintf("%s-%d-%s-%s", proxyName, glooGateway.Spec.GetBindPort(), vs.Name, vs.Namespace)
			// convert the listener portion of the virtual service
			if err := o.convertVirtualServiceListener(vs, glooGateway, listenerName, proxyName); err != nil {
				return err
			}
			// convert the routing portion of the virtual service
			err := o.convertVirtualServiceHTTPRoutes(vs, glooGateway, listenerName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *GatewayAPIOutput) convertVirtualServiceListener(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string, gatewayName string) error {

	// for each VirtualService generate a listener set given the gateway port
	listenerSet := &apixv1a1.XListenerSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "XListenerSet",
			APIVersion: apixv1a1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      listenerName,
			Namespace: vs.GetNamespace(),
			Labels:    vs.GetLabels(),
		},
		Spec: apixv1a1.ListenerSetSpec{
			ParentRef: apixv1a1.ParentGatewayReference{
				Group:     ptr.To(gwv1.Group(wellknown.GatewayGroup)),
				Kind:      ptr.To(gwv1.Kind(wellknown.GatewayKind)),
				Namespace: ptr.To(gwv1.Namespace(glooGateway.GetNamespace())),
				Name:      gwv1.ObjectName(gatewayName),
			},
			Listeners: make([]apixv1a1.ListenerEntry, 0),
		},
	}

	// we only create the listener part, not the http matchers
	for _, hostname := range vs.Spec.GetVirtualHost().GetDomains() {
		if strings.Contains(hostname, ":") {
			o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, vs, "contains port in hostname %s, its being ignored for ListenerSet %s/%s", hostname, listenerSet.Namespace, listenerSet.Name)
			continue
		}

		// listener entry does not support wildcard
		listenerEntryName := strings.ReplaceAll(fmt.Sprintf("%s-%s", vs.Name, hostname), "*", "star")
		entry := apixv1a1.ListenerEntry{
			Name:     gwv1.SectionName(listenerEntryName),
			Hostname: ptr.To(gwv1.Hostname(hostname)),
			Port:     apixv1a1.PortNumber(glooGateway.Spec.GetBindPort()),
			Protocol: gwv1.HTTPProtocolType,
		}
		if vs.Spec.GetSslConfig() != nil {
			tlsConfig := o.generateTLSConfiguration(vs)
			if tlsConfig != nil {
				entry.TLS = tlsConfig
				entry.Protocol = gwv1.HTTPSProtocolType
			}
		}

		// allowed routes
		entry.AllowedRoutes = &gwv1.AllowedRoutes{
			Namespaces: &gwv1.RouteNamespaces{
				From: ptr.To(gwv1.NamespacesFromAll),
			},
		}
		listenerSet.Spec.Listeners = append(listenerSet.Spec.Listeners, entry)
	}

	if vs.Spec.GetVirtualHost().GetOptionsConfigRefs() != nil && len(vs.Spec.GetVirtualHost().GetOptionsConfigRefs().GetDelegateOptions()) > 0 {
		delegateOptions := vs.Spec.GetVirtualHost().GetOptionsConfigRefs().GetDelegateOptions()
		for _, delegateOption := range delegateOptions {
			// check to see if this already exists in gatewayAPI cache, if not move it over from edge cache
			gtp, exists := o.gatewayAPICache.GlooTrafficPolicies[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
			if !exists {
				vho, exists := o.edgeCache.VirtualHostOptions()[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
				if !exists {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, vs, "references VirtualHostOption %s that does not exist", types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()})
					continue
				}
				gtp = o.convertVirtualHostOptionToGlooTrafficPolicy(vho)
			}
			if listenerSet.Namespace != gtp.GlooTrafficPolicy.GetNamespace() {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "VirtualHostOption %s references a listener set in a different namespace %s which is not supported", types.NamespacedName{Name: vs.GetName(), Namespace: vs.GetNamespace()}, types.NamespacedName{Name: listenerSet.GetName(), Namespace: listenerSet.GetNamespace()})
			}
			// add the target ref to the listener
			gtp.GlooTrafficPolicy.Spec.TargetRefs = append(gtp.GlooTrafficPolicy.Spec.TargetRefs, kgateway.LocalPolicyTargetReferenceWithSectionName{
				LocalPolicyTargetReference: kgateway.LocalPolicyTargetReference{
					Group: apixv1a1.GroupName,
					Kind:  "XListenerSet",
					Name:  gwv1.ObjectName(listenerSet.Name),
				},
			})
			o.gatewayAPICache.AddGlooTrafficPolicy(gtp)
		}
	}

	// we need to get the virtualhostoptions and update their references
	if vs.Spec.GetVirtualHost().GetOptions() != nil {
		// create a separate virtualhost option and link it
		gtp := &gloogateway.GlooTrafficPolicy{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GlooTrafficPolicy",
				APIVersion: gloogateway.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      listenerSet.Name,
				Namespace: listenerSet.Namespace,
			},
			Spec: gloogateway.GlooTrafficPolicySpec{
				TrafficPolicySpec: kgateway.TrafficPolicySpec{
					TargetRefs: []kgateway.LocalPolicyTargetReferenceWithSectionName{
						{
							LocalPolicyTargetReference: kgateway.LocalPolicyTargetReference{
								Group: apixv1a1.GroupName,
								Kind:  "XListenerSet",
								Name:  gwv1.ObjectName(listenerSet.Name),
							},
							SectionName: nil,
						},
					},
					AI:             nil,
					Transformation: nil,
					ExtProc:        nil,
					ExtAuth:        nil,
					RateLimit:      nil,
				},
				Waf:                   nil,
				Retry:                 nil,
				Timeouts:              nil,
				RateLimitEnterprise:   nil,
				ExtAuthEnterprise:     nil,
				StagedTransformations: nil,
			},
		}
		//Options:
		//	vs.Spec.GetVirtualHost().GetOptions(),
		// go through each option and add it to traffic policy

		gtp.Spec = o.convertVHOOptionsToTrafficPolicySpec(vs.Spec.GetVirtualHost().GetOptions(), vs)

		o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(gtp, vs.FileOrigin()))
	}
	o.gatewayAPICache.AddListenerSet(snapshot.NewListenerSetWrapper(listenerSet, vs.FileOrigin()))

	return nil
}

func (o *GatewayAPIOutput) convertVirtualHostOptionToGlooTrafficPolicy(vho *snapshot.VirtualHostOptionWrapper) *snapshot.GlooTrafficPolicyWrapper {
	policy := &gloogateway.GlooTrafficPolicy{
		Spec: gloogateway.GlooTrafficPolicySpec{},
	}
	if vho != nil {
		policy.Spec = o.convertVHOOptionsToTrafficPolicySpec(vho.VirtualHostOption.Spec.Options, vho)
	}

	wrapper := snapshot.NewGlooTrafficPolicyWrapper(policy, vho.FileOrigin())
	return wrapper
}

func (o *GatewayAPIOutput) convertVHOOptionsToTrafficPolicySpec(vho *gloov1.VirtualHostOptions, wrapper snapshot.Wrapper) gloogateway.GlooTrafficPolicySpec {

	spec := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs:      nil,
			TargetSelectors: nil,
			AI:              nil,
			Transformation:  nil,
			ExtProc:         nil,
			ExtAuth:         nil,
			RateLimit:       nil,
		},
		Waf:                   nil,
		Retry:                 nil,
		Timeouts:              nil,
		RateLimitEnterprise:   nil,
		ExtAuthEnterprise:     nil,
		StagedTransformations: nil,
	}
	if vho != nil {
		if vho.GetCors() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "CORS is not currently supported in GlooTrafficPolicy")
		}
		if vho.GetBufferPerRoute() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "BufferPerRoute is not currently supported in GlooTrafficPolicy")
		}
		if vho.GetCsrf() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "CSRF is not currently supported in GlooTrafficPolicy")
		}
		if vho.GetExtauth() != nil {

			authConfigWrapper, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{
				Namespace: vho.GetExtauth().GetConfigRef().GetNamespace(),
				Name:      vho.GetExtauth().GetConfigRef().GetName(),
			}]
			if !exists {
				o.AddErrorFromWrapper(ERROR_TYPE_NO_REFERENCES, wrapper, "Unable to find referenced AuthConfig %s", types.NamespacedName{Namespace: vho.GetExtauth().GetConfigRef().GetNamespace(), Name: vho.GetExtauth().GetConfigRef().GetName()}.String())
			} else {
				//copy the auth config to gateway api cache if it doesnt already exist
				_, exists = o.gatewayAPICache.AuthConfigs[types.NamespacedName{
					Namespace: vho.GetExtauth().GetConfigRef().GetNamespace(),
					Name:      vho.GetExtauth().GetConfigRef().GetName(),
				}]
				if !exists {
					o.gatewayAPICache.AddAuthConfig(authConfigWrapper)
				}
			}

			// TODO need to copy auth config over and reference it
			spec.ExtAuthEnterprise = &gloogateway.ExtAuthEnterprise{
				//TODO(nick): need to get the server reference but for now use the default
				ExtensionRef: nil,
				AuthConfigRef: gloogateway.AuthConfigRef{
					Name:      vho.GetExtauth().GetConfigRef().GetName(),
					Namespace: vho.GetExtauth().GetConfigRef().GetNamespace(),
				},
			}
		}
		if vho.GetDlp() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "DLP is not currently supported in GlooTrafficPolicy")
		}
		if vho.GetExtProc() != nil {
			// TODO may need to get settings to create the extension ref, or at least look it up
			extProc := &kgateway.ExtProcPolicy{
				ExtensionRef:   nil,
				ProcessingMode: nil,
			}

			if vho.GetExtProc().GetOverrides() != nil {

			}

			spec.TrafficPolicySpec.ExtProc = extProc
		}
	}

	return spec
}

func (o *GatewayAPIOutput) generateTLSConfiguration(vs *snapshot.VirtualServiceWrapper) *gwv1.GatewayTLSConfig {
	tlsConfig := &gwv1.GatewayTLSConfig{
		Mode: ptr.To(gwv1.TLSModeTerminate),
		//FrontendValidation: nil, // TODO do we need to set this?
		//Options:            nil, // TODO do we need to set this?
	}
	if vs.Spec.GetSslConfig().GetSecretRef() != nil {
		tlsConfig.CertificateRefs = []gwv1.SecretObjectReference{
			{
				Group:     ptr.To(gwv1.Group("")),
				Kind:      ptr.To(gwv1.Kind("Secret")),
				Name:      gwv1.ObjectName(vs.Spec.GetSslConfig().GetSecretRef().GetName()),
				Namespace: ptr.To(gwv1.Namespace(vs.Spec.GetSslConfig().GetSecretRef().GetNamespace())),
			},
		}
	}
	// TODO There is a situation where a SSLSecret contains a ca.crt which triggers mTLS in Gloo Edge we have no way to determine this

	if vs.Spec.GetSslConfig().GetSslFiles() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has SSLFiles but its not supported in Gateway API")
	}
	if vs.Spec.GetSslConfig().GetSds() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has SDS Certificates but its not supported in Gateway API")
	}
	if len(vs.Spec.GetSslConfig().GetVerifySubjectAltName()) > 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has VerifySubjectAltName but its not supported in Gateway API")
	}
	if len(vs.Spec.GetSslConfig().GetAlpnProtocols()) > 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has AlpnProtocols but its not supported in Gateway API")
	}
	if vs.Spec.GetSslConfig().GetOcspStaplePolicy() > 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has OcspStaplePolicy %d but its not supported in Gateway API", vs.Spec.GetSslConfig().GetOcspStaplePolicy())
	}
	return tlsConfig
}

func (o *GatewayAPIOutput) generateGatewaysFromProxyNames(glooGateway *snapshot.GlooGatewayWrapper) {

	proxyNames := glooGateway.Gateway.Spec.GetProxyNames()

	if len(proxyNames) == 0 {
		proxyNames = append(proxyNames, "gateway-proxy")
	}

	for _, proxyName := range glooGateway.Gateway.Spec.GetProxyNames() {
		// check to see if we already created the Gateway, if we did then just move on
		existingGw := o.gatewayAPICache.GetGateway(types.NamespacedName{Name: proxyName, Namespace: glooGateway.Gateway.Namespace})
		if existingGw == nil {
			// create a new gateway
			gwGateway := &gwv1.Gateway{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Gateway",
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      proxyName,
					Namespace: glooGateway.Gateway.Namespace,
					Labels:    glooGateway.Gateway.Labels,
				},
				Spec: gwv1.GatewaySpec{
					AllowedListeners: &gwv1.AllowedListeners{
						Namespaces: &gwv1.ListenerNamespaces{
							From: ptr.To(gwv1.NamespacesFromAll),
						},
					},
					GatewayClassName: "gloo-gateway",
				},
			}
			o.gatewayAPICache.AddGateway(snapshot.NewGatewayWrapper(gwGateway, glooGateway.FileOrigin()))
		}
		if glooGateway.Spec.GetHttpGateway() != nil && glooGateway.Spec.GetHttpGateway().GetOptions() != nil {
			o.convertHTTPListenerOptions(glooGateway.Spec.GetHttpGateway().Options, glooGateway, proxyName)
		}
		if glooGateway.Spec.GetOptions() != nil && glooGateway.Spec.GetOptions() != nil {
			o.convertListenerOptions(glooGateway, proxyName)
		}
	}
}

func (o *GatewayAPIOutput) convertListenerOptions(glooGateway *snapshot.GlooGatewayWrapper, proxyName string) {
	options := glooGateway.Spec.GetOptions()
	if options == nil {
		return
	}
	listenerPolicy := &kgateway.HTTPListenerPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ListenerPolicy",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      glooGateway.GetName(),
			Namespace: glooGateway.GetNamespace(),
			Labels:    glooGateway.Gateway.Labels,
		},
		Spec: kgateway.HTTPListenerPolicySpec{
			TargetRefs: []kgateway.LocalPolicyTargetReference{
				{
					Group: gwv1.Group(gwv1.GroupVersion.Group),
					Kind:  "Gateway",
					Name:  gwv1.ObjectName(proxyName),
				},
			},
		},
	}
	if options.GetExtensions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option extensions are not supported for HTTPTrafficPolicy")
	}
	if options.GetSocketOptions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option socket options are not supported for HTTPTrafficPolicy")
	}
	if options.GetAccessLoggingService() != nil {
		o.convertListenerOptionAccessLogging(glooGateway, options, listenerPolicy)
	}
	if options.GetListenerAccessLoggingService() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option listenerAccessLoggingService is not supported for HTTPTrafficPolicy")
	}
	if options.GetConnectionBalanceConfig() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option connectionBalanceConfig is not supported for HTTPTrafficPolicy")
	}
	if options.GetPerConnectionBufferLimitBytes() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option perConnectionBufferLimitBytes is not supported for HTTPTrafficPolicy")
	}
	if options.GetProxyProtocol() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option proxyProtocol is not supported for HTTPTrafficPolicy")
	}
	if options.GetTcpStats() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option tcpStats is not supported for HTTPTrafficPolicy")
	}

	o.gatewayAPICache.AddHTTPListenerPolicy(snapshot.NewHTTPListenerPolicyWrapper(listenerPolicy, glooGateway.FileOrigin()))
}

func (o *GatewayAPIOutput) convertListenerOptionAccessLogging(glooGateway *snapshot.GlooGatewayWrapper, options *gloov1.ListenerOptions, listenerPolicy *kgateway.HTTPListenerPolicy) {
	accessLoggingService := options.GetAccessLoggingService()

	for _, edgeAccessLog := range accessLoggingService.GetAccessLog() {
		if listenerPolicy.Spec.AccessLog == nil {
			listenerPolicy.Spec.AccessLog = []kgateway.AccessLog{}
		}
		accessLog := kgateway.AccessLog{
			FileSink:    nil,
			GrpcService: nil,
			Filter:      nil,
		}
		if edgeAccessLog.GetFileSink() != nil {
			fileSink := &kgateway.FileSink{
				Path: edgeAccessLog.GetFileSink().Path,
			}
			if jsonFormat := edgeAccessLog.GetFileSink().GetJsonFormat(); jsonFormat != nil {
				jsonBytes, err := json.Marshal(jsonFormat.AsMap())
				if err != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, glooGateway, "unable to marshal json format for accessLoggingService %v", err)
				} else {
					fileSink.JsonFormat = &runtime.RawExtension{Raw: jsonBytes}
				}
			}
			if edgeAccessLog.GetFileSink().GetStringFormat() != "" {
				fileSink.StringFormat = edgeAccessLog.GetFileSink().GetStringFormat()
			}
			accessLog.FileSink = fileSink
		}
		if edgeAccessLog.GetGrpcService() != nil {
			accessLog.GrpcService = &kgateway.GrpcService{
				LogName:                         edgeAccessLog.GetGrpcService().LogName,
				AdditionalRequestHeadersToLog:   edgeAccessLog.GetGrpcService().AdditionalRequestHeadersToLog,
				AdditionalResponseHeadersToLog:  edgeAccessLog.GetGrpcService().AdditionalResponseHeadersToLog,
				AdditionalResponseTrailersToLog: edgeAccessLog.GetGrpcService().AdditionalResponseTrailersToLog,
			}

			// backend Ref
			switch edgeAccessLog.GetGrpcService().GetServiceRef().(type) {
			case *als.GrpcService_StaticClusterName:
				accessLog.GrpcService.BackendRef = &gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      gwv1.ObjectName(edgeAccessLog.GetGrpcService().GetStaticClusterName()),
						Namespace: ptr.To(gwv1.Namespace("UNKNOWN")),
						Port:      ptr.To(gwv1.PortNumber(0)),
					},
				}
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, glooGateway, "", edgeAccessLog.GetGrpcService().GetStaticClusterName())
			}
			if edgeAccessLog.GetGrpcService().GetFilterStateObjectsToLog() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "accessLoggingService grpcService has filterStateObjectsToLog but its not supported in kgateway")
			}
		}
		if edgeAccessLog.GetFilter() != nil {
			if accessLog.Filter == nil {
				accessLog.Filter = &kgateway.AccessLogFilter{}
			}
			if edgeAccessLog.GetFilter().GetOrFilter() != nil {

				if accessLog.Filter.OrFilter == nil {
					accessLog.Filter.OrFilter = []kgateway.FilterType{}
				}
				for _, filter := range edgeAccessLog.GetFilter().GetOrFilter().GetFilters() {
					accessLog.Filter.AndFilter = append(accessLog.Filter.AndFilter, *o.convertAccessLogFitler(filter, glooGateway))
				}
			} else if edgeAccessLog.GetFilter().GetAndFilter() != nil {
				if accessLog.Filter.AndFilter == nil {
					accessLog.Filter.AndFilter = []kgateway.FilterType{}
				}
				for _, filter := range edgeAccessLog.GetFilter().GetAndFilter().GetFilters() {
					accessLog.Filter.AndFilter = append(accessLog.Filter.AndFilter, *o.convertAccessLogFitler(filter, glooGateway))
				}
			} else {
				// just and inline filter
				accessLog.Filter.FilterType = o.convertAccessLogFitler(edgeAccessLog.GetFilter(), glooGateway)
			}
		}
		listenerPolicy.Spec.AccessLog = append(listenerPolicy.Spec.AccessLog, accessLog)
	}
}

func (o *GatewayAPIOutput) convertAccessLogFitler(filter *als.AccessLogFilter, wrapper snapshot.Wrapper) *kgateway.FilterType {

	filterType := &kgateway.FilterType{}

	if filter.GetDurationFilter() != nil {
		filterType.DurationFilter = &kgateway.DurationFilter{
			Op:    kgateway.Op(filter.GetDurationFilter().GetComparison().Op),
			Value: filter.GetDurationFilter().GetComparison().GetValue().GetDefaultValue(),
		}
	}
	if filter.GetHeaderFilter() != nil && filter.GetHeaderFilter().GetHeader() != nil {
		headerMatch := gwv1.HTTPHeaderMatch{
			Name: gwv1.HTTPHeaderName(filter.GetHeaderFilter().GetHeader().Name),
		}

		if filter.GetHeaderFilter().GetHeader().GetExactMatch() != "" {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService grpcService has filterStateObjectsToLog but its not supported in kgateway")
		}

		if filter.GetHeaderFilter().GetHeader().GetInvertMatch() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header invert match is not supported in kgateway")
		}

		if filter.GetHeaderFilter().GetHeader().GetPresentMatch() == true {
			// TODO(nick): is this supported in Gateway API?
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header range match match is not supported in kgateway")
		}

		if filter.GetHeaderFilter().GetHeader().GetPrefixMatch() != "" {
			//	HeaderMatchExact             HeaderMatchType = "Exact"
			//	HeaderMatchRegularExpression HeaderMatchType = "RegularExpression"
			//TODO(nick): can someone verify this is the equivalent?
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = filter.GetHeaderFilter().GetHeader().GetPrefixMatch() + ".*"
		}

		if filter.GetHeaderFilter().GetHeader().GetRangeMatch() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header range match match is not supported in kgateway")
		}

		if filter.GetHeaderFilter().GetHeader().GetSafeRegexMatch() != nil {
			// Edge only supported Googles Regex (RE2) which might not be compatible with Gateway API regex
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = filter.GetHeaderFilter().GetHeader().GetSafeRegexMatch().Regex
		}

		if filter.GetHeaderFilter().GetHeader().GetSuffixMatch() != "" {
			//TODO(nick): can someone verify this is the equivalent?
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = ".*" + filter.GetHeaderFilter().GetHeader().GetPrefixMatch()
		}

		filterType.HeaderFilter = &kgateway.HeaderFilter{
			Header: headerMatch,
		}
	}

	if filter.GetGrpcStatusFilter() != nil {
		grpcFilter := &kgateway.GrpcStatusFilter{
			Statuses: []kgateway.GrpcStatus{},
			Exclude:  filter.GetGrpcStatusFilter().Exclude,
		}
		for _, status := range filter.GetGrpcStatusFilter().Statuses {
			grpcFilter.Statuses = append(grpcFilter.Statuses, kgateway.GrpcStatus(status))
		}
		filterType.GrpcStatusFilter = grpcFilter
	}

	if filter.GetNotHealthCheckFilter() != nil {
		//unsure if this is correct. it appears this just needs to exist to function?
		filterType.NotHealthCheckFilter = true
	}

	if filter.GetResponseFlagFilter() != nil {
		filterType.ResponseFlagFilter = &kgateway.ResponseFlagFilter{
			Flags: filter.GetResponseFlagFilter().GetFlags(),
		}
	}

	if filter.GetRuntimeFilter() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService runtimeFilter is not supported in kgateway")
	}

	if filter.GetTraceableFilter() != nil {
		//unsure if this is correct. it appears this just needs to exist to function?
		filterType.TraceableFilter = true
	}

	if filter.GetStatusCodeFilter() != nil {
		filterType.StatusCodeFilter = &kgateway.StatusCodeFilter{
			Op:    kgateway.Op(filter.GetStatusCodeFilter().GetComparison().GetOp()),
			Value: filter.GetStatusCodeFilter().GetComparison().GetValue().GetDefaultValue(),
		}
	}

	return nil
}

// convertHTTPListenerOptions - generates GlooTrafficPolicy applied to the Gateway
func (o *GatewayAPIOutput) convertHTTPListenerOptions(options *gloov1.HttpListenerOptions, wrapper snapshot.Wrapper, proxyName string) {
	if options == nil {
		return
	}

	trafficPolicy := &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "TrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      wrapper.GetName(),
			Namespace: wrapper.GetNamespace(),
			Labels:    wrapper.GetLabels(),
		},
	}

	tps := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs: []kgateway.LocalPolicyTargetReferenceWithSectionName{
				{
					LocalPolicyTargetReference: kgateway.LocalPolicyTargetReference{
						Group: gwv1.Group(gwv1.GroupVersion.Group),
						Kind:  "Gateway",
						Name:  gwv1.ObjectName(proxyName),
					},
				},
			},
			TargetSelectors: nil,
			AI:              nil,
			Transformation:  nil,
			ExtProc:         nil,
			ExtAuth:         nil,
			RateLimit:       nil,
		},
		Waf:                   nil,
		Retry:                 nil,
		Timeouts:              nil,
		RateLimitEnterprise:   nil,
		ExtAuthEnterprise:     nil,
		StagedTransformations: nil,
	}

	// go through each option in Gateway Options and convert to listener policy
	if options.GetExtauth() != nil {
		// These are global extAuthSettings that are also on the Settings Object.
		// TODO(nick): Implement auth settings at GlooTrafficPolicy spec

		// If this exists we need to generate a GatewayExtensionObject for this
		gatewayExtensions := o.generateGatewayExtension(options.Extauth, wrapper.GetName(), wrapper)
		o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	if options.GetExtProc() != nil {
		// TODO(nick): Implement ext proc settings at GlooTrafficPolicy spec
	}

	if options.GetHttpLocalRatelimit() != nil {
		if options.GetHttpLocalRatelimit().GetEnableXRatelimitHeaders() != nil && options.GetHttpLocalRatelimit().GetEnableXRatelimitHeaders().GetValue() {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpLocalRateLimit enableXRateLimitHeaders is not supported in kgateway")
		}
		if options.GetHttpLocalRatelimit().GetLocalRateLimitPerDownstreamConnection() != nil && options.GetHttpLocalRatelimit().GetLocalRateLimitPerDownstreamConnection().GetValue() {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpLocalRateLimit localRateLimitPerDownstreamConnection is not supported in kgateway")
		}
		if options.GetHttpLocalRatelimit().GetDefaultLimit() != nil {
			rl := &kgateway.RateLimit{
				Local: &kgateway.LocalRateLimitPolicy{
					TokenBucket: &kgateway.TokenBucket{
						MaxTokens: options.GetHttpLocalRatelimit().GetDefaultLimit().GetMaxTokens(),
					},
				},
			}
			if options.GetHttpLocalRatelimit().GetDefaultLimit().GetTokensPerFill() != nil {
				rl.Local.TokenBucket.TokensPerFill = ptr.To(options.GetHttpLocalRatelimit().GetDefaultLimit().GetTokensPerFill().GetValue())
			}
			if options.GetHttpLocalRatelimit().GetDefaultLimit().GetFillInterval() != nil {
				rl.Local.TokenBucket.FillInterval = gwv1.Duration(options.GetHttpLocalRatelimit().GetDefaultLimit().GetFillInterval().AsDuration().String())
			}
			tps.TrafficPolicySpec.RateLimit = rl
		}
	}

	if options.GetWaf() != nil {
		waf := &gloogateway.Waf{
			Disabled:      ptr.To(options.GetWaf().Disabled),
			CustomMessage: options.GetWaf().CustomInterventionMessage,
			Rules:         make([]gloogateway.WafRule, len(options.GetWaf().RuleSets)),
		}
		for _, r := range options.GetWaf().RuleSets {
			waf.Rules = append(waf.Rules, gloogateway.WafRule{
				RuleStr: r.RuleStr,
			})
			if r.Files != nil && len(r.Files) > 0 {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF files is not supported in kgateway")
			}
		}
		tps.Waf = waf
	}
	if options.GetDisableExtProc() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "disableExtProc is not supported in kgateway")
	}
	if options.GetNetworkLocalRatelimit() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "networkLocalRateLimit is not supported in kgateway")
	}
	if options.GetDlp() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported in kgateway")
	}
	if options.GetCsrf() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "csrf is not supported in kgateway")
	}
	if options.GetBuffer() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "buffer is not supported in kgateway")
	}
	if options.GetCaching() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "caching is not supported in kgateway")
	}
	if options.GetConnectionLimit() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "connectionlimit is not supported in kgateway")
	}
	if options.GetDynamicForwardProxy() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dynamicForwardProxy (DFP) is not supported in kgateway")
	}
	if options.GetExtensions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported in kgateway")
	}
	if options.GetGrpcJsonTranscoder() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "grpcToJson is not supported in kgateway")
	}
	if options.GetGrpcWeb() != nil {
		//TODO(nick) : GRPCWeb is enabled by default in edge. we need to verify the same in kgateway.
		//o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "grpcWeb is not supported in kgateway")
	}
	if options.GetGzip() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gzip is not supported in kgateway")
	}
	if options.GetHeaderValidationSettings() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "header validation is not supported in kgateway")
	}
	if options.GetHealthCheck() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "health check is not supported in kgateway")
	}
	if options.GetHttpConnectionManagerSettings() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManagerSettings is not supported in kgateway")
	}
	if options.GetProxyLatency() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "proxy latency is not supported in kgateway")
	}
	if options.GetRouter() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "router (envoy filter maps) is not supported in kgateway")
	}
	if options.GetSanitizeClusterHeader() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "sanitize cluster header is not supported in kgateway")
	}
	if options.GetStatefulSession() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "statefulSession is not supported in kgateway")
	}
	if options.GetTap() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "Tap filter is not supported in kgateway")
	}
	if options.GetWasm() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WASM is not supported in kgateway")
	}

	trafficPolicy.Spec = tps

	o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(trafficPolicy, wrapper.FileOrigin()))
}

func (o *GatewayAPIOutput) generateGatewayExtension(extauth *v1.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {

	if extauth == nil {
		return nil
	}

	var grpcService *kgateway.ExtGrpcService

	if extauth.GetExtauthzServerRef() != nil {
		grpcService = &kgateway.ExtGrpcService{
			BackendRef: &gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(extauth.GetExtauthzServerRef().GetName()),
					Namespace: ptr.To(gwv1.Namespace(extauth.GetExtauthzServerRef().GetNamespace())),
					// using the default port here
					Port: ptr.To(gwv1.PortNumber(8083)),
				},
			},
		}

		if extauth.GetGrpcService() != nil {
			grpcService.Authority = ptr.To(extauth.GetGrpcService().GetAuthority())
		}
	}

	if extauth.GetClearRouteCache() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings clearRouteCache is not supported in kgateway")
	}
	if extauth.GetFailureModeAllow() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings clearRouteCache is not supported in kgateway")
	}
	if extauth.GetHttpService() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings httpService is not supported in kgateway")
	}
	if extauth.GetRequestBody() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings requestBody is not supported in kgateway")
	}
	if extauth.GetRequestTimeout() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings requestTimeout is not supported in kgateway")
	}
	if extauth.GetStatPrefix() != "" {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings statPrefix is not supported in kgateway")
	}
	if extauth.GetStatusOnError() != 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings statusOnError is not supported in kgateway")
	}
	if extauth.GetTransportApiVersion() != 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings transportApiVersion is not supported in kgateway")
	}
	if extauth.GetUserIdHeader() != "" {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings userIdHeader is not supported in kgateway")
	}

	gatewayExtension := &kgateway.GatewayExtension{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GatewayExtension",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: wrapper.GetNamespace(),
			Labels:    wrapper.GetLabels(),
		},
		Spec: kgateway.GatewayExtensionSpec{
			Type: kgateway.GatewayExtensionTypeExtAuth,
			ExtAuth: &kgateway.ExtAuthProvider{
				GrpcService: grpcService,
			},
		},
	}
	return gatewayExtension
}

func (o *GatewayAPIOutput) convertVirtualServiceHTTPRoutes(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string) error {

	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gwv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vs.GetName(),
			Namespace: vs.GetNamespace(),
			Labels:    vs.GetLabels(),
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{
				ParentRefs: []gwv1.ParentReference{
					{
						Name:      gwv1.ObjectName(listenerName),
						Namespace: ptr.To(gwv1.Namespace(glooGateway.GetNamespace())),
						Kind:      ptr.To(gwv1.Kind("XListenerSet")),
						Group:     ptr.To(gwv1.Group(apixv1a1.GroupVersion.Group)),
					},
				},
			},
			Hostnames: convertDomains(vs.Spec.GetVirtualHost().GetDomains()),
			Rules:     []gwv1.HTTPRouteRule{},
		},
	}

	for _, route := range vs.Spec.GetVirtualHost().GetRoutes() {
		rule, err := o.convertRouteToRule(route, vs)
		if err != nil {
			return err
		}
		hr.Spec.Rules = append(hr.Spec.Rules, rule)
	}

	o.gatewayAPICache.AddHTTPRoute(snapshot.NewHTTPRouteWrapper(hr, vs.FileOrigin()))

	return nil
}

func (o *GatewayAPIOutput) convertRouteOptions(
	options *gloov1.RouteOptions,
	routeName string,
	wrapper snapshot.Wrapper,
) (*gloogateway.GlooTrafficPolicy, *gwv1.HTTPRouteFilter) {

	var trafficPolicy *gloogateway.GlooTrafficPolicy
	var filter *gwv1.HTTPRouteFilter
	associationID := RandStringRunes(RandomSuffix)
	if routeName == "" {
		routeName = "route-association"
	}
	associationName := fmt.Sprintf("%s-%s", routeName, associationID)

	// converts options to RouteOptions but we need to this for everything except prefixrewrite and a few others now
	if isRouteOptionsSet(options) {
		gtpSpec := &gloogateway.GlooTrafficPolicySpec{
			TrafficPolicySpec: kgateway.TrafficPolicySpec{
				TargetRefs:      nil,
				TargetSelectors: nil,
				AI:              nil,
				Transformation:  nil,
				ExtProc:         nil,
				ExtAuth:         nil,
				RateLimit:       nil,
			},
			Waf:                   nil,
			Retry:                 nil,
			Timeouts:              nil,
			RateLimitEnterprise:   nil,
			ExtAuthEnterprise:     nil,
			StagedTransformations: nil,
		}

		//// Because we move rewrites to a filter we need to remove it from RouteOptions
		// TODO: delete this because this was for RouteOption and not needed for GlooTrafficPolicy
		//if options.GetPrefixRewrite() != nil {
		//	trafficPolicy.Spec.GetOptions().PrefixRewrite = nil
		//}

		filter = &gwv1.HTTPRouteFilter{
			Type: gwv1.HTTPRouteFilterExtensionRef,
			ExtensionRef: &gwv1.LocalObjectReference{
				Group: glookube.GroupName,
				Kind:  "GlooTrafficPolicy",
				Name:  gwv1.ObjectName(associationName),
			},
		}
		if options.GetExtauth() != nil && options.GetExtauth().GetConfigRef() != nil {
			// we need to copy over the auth config ref if it exists
			ref := options.GetExtauth().GetConfigRef()
			ac, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
			if !exists {
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
			}
			o.gatewayAPICache.AddAuthConfig(ac)
		}
		if options.GetAi() != nil {
			aip := &kgateway.AIPolicy{
				PromptEnrichment: nil,
				PromptGuard:      nil,
				Defaults:         make([]kgateway.FieldDefault, len(options.GetAi().GetDefaults())),
			}
			switch options.GetAi().GetRouteType() {
			case ai.RouteSettings_CHAT:
				aip.RouteType = ptr.To(kgateway.CHAT)
			case ai.RouteSettings_CHAT_STREAMING:
				aip.RouteType = ptr.To(kgateway.CHAT_STREAMING)
			}
			for _, d := range options.GetAi().GetDefaults() {
				aip.Defaults = append(aip.Defaults, kgateway.FieldDefault{
					Field:    d.Field,
					Value:    d.Value.String(),
					Override: ptr.To(d.Override),
				})
			}
			if options.GetAi().GetPromptEnrichment() != nil {
				enrichment := &kgateway.AIPromptEnrichment{}

				for _, prepend := range options.GetAi().GetPromptEnrichment().GetPrepend() {
					enrichment.Prepend = append(enrichment.Prepend, kgateway.Message{
						Role:    prepend.GetRole(),
						Content: prepend.GetContent(),
					})
				}

				aip.PromptEnrichment = enrichment
			}
			if options.GetAi().GetRag() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai RAG is not supported in kgateway")
			}
			if options.GetAi().GetPromptGuard() != nil {
				guard := o.generateAIPromptGuard(options, wrapper)
				aip.PromptGuard = guard
			}
			if options.GetAi().GetSemanticCache() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai SemanticCache is not supported in kgateway")
			}
			gtpSpec.AI = aip
		}
	}
	if options.GetWaf() != nil {

	}
	if options.GetDlp() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported in kgateway")
	}
	if options.GetCsrf() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "csrf is not supported in kgateway")
	}
	if options.GetExtensions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported in kgateway")
	}
	if options.GetBufferPerRoute() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "bufferPerRoute is not supported in kgateway")
	}
	if options.GetCors() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "cors is not supported in kgateway")
	}
	if options.GetAppendXForwardedHost() != nil && options.GetAppendXForwardedHost().GetValue() == true {

	}
	if options.GetAutoHostRewrite() != nil {

	}
	if options.GetEnvoyMetadata() != nil {

	}
	if options.GetFaults() != nil {

	}
	if options.GetHeaderManipulation() != nil {

	}
	if options.GetHostRewrite() != "" {

	}
	if options.GetHostRewriteHeader() != nil {

	}
	if options.GetHostRewritePathRegex() != nil {

	}
	if options.GetIdleTimeout() != nil {

	}
	if options.GetJwtProvidersStaged() != nil {

	}
	if options.GetJwtStaged() != nil {

	}
	if options.GetLbHash() != nil {

	}
	if options.GetMaxStreamDuration() != nil {

	}
	if options.GetPrefixRewrite() != nil {

	}
	if options.GetRatelimit() != nil {

	}
	if options.GetRatelimitBasic() != nil {

	}
	if options.GetRegexRewrite() != nil {

	}
	if options.GetRbac() != nil {

	}
	if options.GetRetries() != nil {

	}
	if options.GetShadowing() != nil {

	}
	if options.GetUpgrades() != nil {

	}
	trafficPolicy = &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      associationName,
			Namespace: wrapper.GetNamespace(),
		},
		Spec: gtpSpec,
	}

	return trafficPolicy, filter
}

func (o *GatewayAPIOutput) generateAIPromptGuard(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.AIPromptGuard {
	guard := &kgateway.AIPromptGuard{
		Request:  nil,
		Response: nil,
	}
	if options.GetAi().GetPromptGuard().GetRequest() != nil {
		guard.Request = &kgateway.PromptguardRequest{
			CustomResponse: &kgateway.CustomResponse{
				Message:    ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetMessage()),
				StatusCode: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetStatusCode()),
			},
		}
		if options.GetAi().GetPromptGuard().GetRequest().GetModeration() != nil && options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai() != nil {
			guard.Request.Moderation = &kgateway.Moderation{
				OpenAIModeration: &kgateway.OpenAIConfig{
					Model: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetModel()),
				},
			}
			if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken() != nil {
				authToken := kgateway.SingleAuthToken{}
				if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetInline() != "" {
					authToken.Kind = kgateway.Inline
					authToken.Inline = ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetInline())
				}
				if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef() != nil {
					authToken.Kind = kgateway.SecretRef
					authToken.SecretRef = &corev1.LocalObjectReference{
						Name: options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef().GetName(),
					}
					if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef().GetNamespace() != wrapper.GetNamespace() {
						o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "AI AuthToken secretRef may be referencing secret outside configuration namespace")
					}
				}
				guard.Request.Moderation.OpenAIModeration.AuthToken = authToken
			}
		}

		if options.GetAi().GetPromptGuard().GetRequest().GetWebhook() != nil {
			webhook := &kgateway.Webhook{
				Host: kgateway.Host{
					Host: options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetHost(),
					Port: gwv1.PortNumber(options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetPort()),
					//InsecureSkipVerify: nil,
				},
				ForwardHeaders: make([]gwv1.HTTPHeaderMatch, len(options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetForwardHeaders())),
			}
			for _, h := range options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetForwardHeaders() {
				match := gwv1.HTTPHeaderMatch{
					Name: gwv1.HTTPHeaderName(h.GetKey()),
					//Value: nil,
				}
				// TODO(nick) - We have a lot of options but gateway API only has exact or regex....
				switch h.GetMatchType() {
				case ai.AIPromptGuard_Webhook_HeaderMatch_CONTAINS:
					match.Type = ptr.To(gwv1.HeaderMatchExact)
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'contains' is not supported in kgateway")
				case ai.AIPromptGuard_Webhook_HeaderMatch_EXACT:
					match.Type = ptr.To(gwv1.HeaderMatchExact)
				case ai.AIPromptGuard_Webhook_HeaderMatch_PREFIX:
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'prefix' is not supported in kgateway")
					match.Type = ptr.To(gwv1.HeaderMatchExact)
				case ai.AIPromptGuard_Webhook_HeaderMatch_REGEX:
					match.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
				case ai.AIPromptGuard_Webhook_HeaderMatch_SUFFIX:
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'suffix' is not supported in kgateway")
					match.Type = ptr.To(gwv1.HeaderMatchExact)
				}
				webhook.ForwardHeaders = append(webhook.ForwardHeaders, match)
			}
			guard.Request.Webhook = webhook
		}

		if options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse() != nil {
			guard.Request.CustomResponse = &kgateway.CustomResponse{
				Message:    ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetMessage()),
				StatusCode: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetStatusCode()),
			}
		}
		if options.GetAi().GetPromptGuard().GetRequest().GetRegex() != nil {
			guard.Request.Regex = &kgateway.Regex{
				Matches:  make([]kgateway.RegexMatch, len(options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetMatches())),
				Builtins: make([]kgateway.BuiltIn, len(options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetBuiltins())),
			}
			switch options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetAction() {
			case ai.AIPromptGuard_Regex_MASK:
				guard.Request.Regex.Action = ptr.To(kgateway.MASK)
			case ai.AIPromptGuard_Regex_REJECT:
				guard.Request.Regex.Action = ptr.To(kgateway.REJECT)
			}

			for _, match := range options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetMatches() {
				guard.Request.Regex.Matches = append(guard.Request.Regex.Matches, kgateway.RegexMatch{
					Pattern: ptr.To(match.GetPattern()),
					Name:    ptr.To(match.GetName()),
				})
			}
			guard.Request.Regex.Builtins = make([]kgateway.BuiltIn, len(options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetBuiltins()))
			for _, builtIns := range options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetBuiltins() {
				switch builtIns {
				case ai.AIPromptGuard_Regex_SSN:
					guard.Request.Regex.Builtins = append(guard.Request.Regex.Builtins, kgateway.SSN)
				case ai.AIPromptGuard_Regex_CREDIT_CARD:
					guard.Request.Regex.Builtins = append(guard.Request.Regex.Builtins, kgateway.CREDIT_CARD)
				case ai.AIPromptGuard_Regex_PHONE_NUMBER:
					guard.Request.Regex.Builtins = append(guard.Request.Regex.Builtins, kgateway.PHONE_NUMBER)
				}
			}
		}
	}
	return guard
}

func (o *GatewayAPIOutput) convertRouteToRule(r *gloogwv1.Route, wrapper snapshot.Wrapper) (gwv1.HTTPRouteRule, error) {

	rr := gwv1.HTTPRouteRule{
		Matches:     []gwv1.HTTPRouteMatch{},
		Filters:     []gwv1.HTTPRouteFilter{},
		BackendRefs: []gwv1.HTTPBackendRef{},
	}

	// unused fields
	if r.GetInheritablePathMatchers() != nil && r.GetInheritablePathMatchers().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has inheritable path matchers but there is not equivalent in Gateway API")
	}
	if r.GetInheritableMatchers() != nil && r.GetInheritableMatchers().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has inheritable matchers but there is not equivalent in Gateway API")
	}

	for _, m := range r.GetMatchers() {
		match, err := o.convertMatch(m, wrapper)
		if err != nil {
			return rr, err
		}
		rr.Matches = append(rr.Matches, match)
	}
	if r.GetOptions() != nil {
		options := r.GetOptions()

		// prefix rewrite, sets it on HTTPRoute
		if options.GetPrefixRewrite() != nil {
			rf := o.generateFilterForURLRewrite(r, wrapper)
			if rf != nil {
				rr.Filters = append(rr.Filters, *rf)
			}
		}

		glooTrafficPolicy, filter := o.convertRouteOptions(options, r.GetName(), wrapper)
		if filter != nil {
			rr.Filters = append(rr.Filters, *filter)
		}
		if glooTrafficPolicy != nil {
			o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(glooTrafficPolicy, wrapper.FileOrigin()))
		}
	}
	// Process Route_Actions
	if r.GetRouteAction() != nil {
		// Route_Route Action
		if r.GetRouteAction().GetClusterHeader() != "" {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has cluster header action set but there is not equivalent in Gateway API")
		}
		if r.GetRouteAction().GetDynamicForwardProxy() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has dynamic forward proxy action set but there is not equivalent in Gateway API")
		}
		if r.GetRouteAction().GetMulti() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multi detination action set but there is not equivalent in Gateway API")
		}

		if r.GetRouteAction().GetSingle() != nil {
			// single static upstream
			if r.GetRouteAction().GetSingle().GetUpstream() != nil {
				backendRef := o.generateBackendRefForSingleUpstream(r, wrapper)

				rr.BackendRefs = append(rr.BackendRefs, backendRef)
			}
		}
		if r.GetRouteAction().GetUpstreamGroup() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has upstream group action set but there is not equivalent in Gateway API")
		}

	} else if r.GetRedirectAction() != nil {
		rdf := o.convertRedirect(r, wrapper)

		rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
			Type:            "RequestRedirect",
			RequestRedirect: rdf,
		})

	} else if r.GetDirectResponseAction() != nil {

		dr := convertDirectResponse(r.GetDirectResponseAction())
		if dr != nil {
			// TODO(nick): what if route name is nil?
			rName := r.GetName()
			if rName == "" {
				rName = RandStringRunes(6)
			}
			drName := fmt.Sprintf("directresponse-%s-%s", wrapper.GetName(), rName)
			dr.Name = drName
			dr.Namespace = wrapper.GetNamespace()
			o.gatewayAPICache.AddDirectResponse(snapshot.NewDirectResponseWrapper(dr, wrapper.FileOrigin()))

			rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
				Type: "ExtensionRef",
				ExtensionRef: &gwv1.LocalObjectReference{
					Group: v1alpha1.Group,
					Kind:  "DirectResponse",
					Name:  gwv1.ObjectName(drName),
				},
			})
		}

	} else if r.GetDelegateAction() != nil {
		// delegate action
		// intermediate delegation step. This is a placeholder for the next path to do delegation
		backendRef := o.generateBackendRefForDelegateAction(r, wrapper)

		if len(backendRef) > 0 {
			for _, b := range backendRef {
				rr.BackendRefs = append(rr.BackendRefs, *b)
			}
		}
	}

	if r.GetOptionsConfigRefs() != nil && len(r.GetOptionsConfigRefs().GetDelegateOptions()) > 0 {
		// these are references to other RouteOptions, we need to add them
		for _, delegateOptions := range r.GetOptionsConfigRefs().GetDelegateOptions() {
			if delegateOptions.GetNamespace() != wrapper.GetNamespace() {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "delegates to route options not in same namespace (this does not work in Gateway API)")
			}

			// grab that route option and convert it to GlooTrafficPolicy
			ro, exists := o.edgeCache.RouteOptions()[types.NamespacedName{Name: delegateOptions.GetName(), Namespace: delegateOptions.GetNamespace()}]
			if !exists {
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find RouteOption %s/%s for delegated route option reference", delegateOptions.GetNamespace(), delegateOptions.GetName())
			}

			if ro.Spec.GetOptions() != nil && ro.Spec.GetOptions().GetExtauth() != nil && ro.Spec.GetOptions().GetExtauth().GetConfigRef() != nil {
				// we need to copy over the auth config ref if it exists
				ref := ro.Spec.GetOptions().GetExtauth().GetConfigRef()
				ac, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
				if !exists {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, ro, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
				}
				o.gatewayAPICache.AddAuthConfig(ac)
			}

			gtp, filter := o.convertRouteOptions(ro.RouteOption.Spec.GetOptions(), delegateOptions.GetName(), ro)
			o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(gtp, ro.FileOrigin()))
			rr.Filters = append(rr.Filters, *filter)
		}
	}

	return rr, nil
}

func (o *GatewayAPIOutput) convertRedirect(r *gloogwv1.Route, wrapper snapshot.Wrapper) *gwv1.HTTPRequestRedirectFilter {
	rdf := &gwv1.HTTPRequestRedirectFilter{}

	action := r.GetRedirectAction()
	if action.GetHttpsRedirect() {
		rdf.Scheme = ptr.To("https")
	}
	if action.GetStripQuery() {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has stripQuery redirect action but there is not equivalent in Gateway API")
	}
	if action.GetRegexRewrite() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has regexRewrite redirect action but there is not equivalent in Gateway API")
	}
	if action.GetPrefixRewrite() != "" {
		match, err := isPrefixMatch(r.GetMatchers())
		if err != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
		}

		if match {
			// full path rewrite
			rdf.Path = &gwv1.HTTPPathModifier{
				Type:               gwv1.PrefixMatchHTTPPathModifier,
				ReplacePrefixMatch: ptr.To(action.GetPrefixRewrite()),
			}
		}

	}
	if action.GetHostRedirect() != "" {
		rdf.Hostname = ptr.To(gwv1.PreciseHostname(action.GetHostRedirect()))
	}
	if action.GetPathRedirect() != "" {
		match, err := isExactMatch(r.GetMatchers())
		if err != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
		}
		if match {
			// full path rewrite
			rdf.Path = &gwv1.HTTPPathModifier{
				Type:            gwv1.FullPathHTTPPathModifier,
				ReplaceFullPath: ptr.To(action.GetPathRedirect()),
			}
		}
	}

	if action.GetPortRedirect() != nil {
		rdf.Port = ptr.To(gwv1.PortNumber(action.GetPortRedirect().GetValue()))
	}

	switch action.GetResponseCode() {
	case gloov1.RedirectAction_MOVED_PERMANENTLY:
		rdf.StatusCode = ptr.To(301)
	case gloov1.RedirectAction_FOUND:
		rdf.StatusCode = ptr.To(302)
	case gloov1.RedirectAction_SEE_OTHER:
		rdf.StatusCode = ptr.To(303)
	case gloov1.RedirectAction_TEMPORARY_REDIRECT:
		rdf.StatusCode = ptr.To(307)
	case gloov1.RedirectAction_PERMANENT_REDIRECT:
		rdf.StatusCode = ptr.To(308)
	default:
		rdf.StatusCode = ptr.To(301)
	}
	return rdf
}
func convertDirectResponse(action *gloov1.DirectResponseAction) *kgateway.DirectResponse {
	if action == nil {
		return nil
	}
	dr := &kgateway.DirectResponse{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DirectResponse",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: kgateway.DirectResponseSpec{
			StatusCode: action.GetStatus(),
			Body:       action.GetBody(),
		},
	}

	return dr
}

func (o *GatewayAPIOutput) generateBackendRefForDelegateAction(
	r *gloogwv1.Route,
	wrapper snapshot.Wrapper,
) []*gwv1.HTTPBackendRef {
	var backends []*gwv1.HTTPBackendRef
	if r.GetDelegateAction().GetRef() != nil {
		delegate := r.GetDelegateAction().GetRef()
		backendRef := &gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(delegate.GetName()),
					Namespace: (*gwv1.Namespace)(ptr.To(delegate.GetNamespace())),
					Kind:      (*gwv1.Kind)(ptr.To("HTTPRoute")),
					Group:     (*gwv1.Group)(ptr.To(gwv1.GroupVersion.Group)),
				},
			},
		}
		backends = append(backends, backendRef)
	} else if r.GetDelegateAction().GetSelector() != nil {

		selector := r.GetDelegateAction().GetSelector()
		namespaces := selector.GetNamespaces()
		if namespaces != nil || len(selector.GetNamespaces()) == 0 {
			// default namespace is gloo-system
			namespaces = []string{"gloo-system"}
		}

		for _, namespace := range selector.GetNamespaces() {
			if namespace == "*" {
				namespace = "all"
			}

			if len(selector.GetLabels()) > 1 {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has delegate action with more than one label selector which is not supported in Gateway API")
			}
			// create a backend ref for each label
			for _, v := range selector.GetLabels() {
				// just grab the first label
				backendRef := &gwv1.HTTPBackendRef{
					BackendRef: gwv1.BackendRef{
						BackendObjectReference: gwv1.BackendObjectReference{
							Name:      gwv1.ObjectName(v),                               // label value
							Namespace: ptr.To(gwv1.Namespace(namespace)),                // defaults to parent namespace if unset
							Kind:      ptr.To(gwv1.Kind("label")),                       // label is the only value
							Group:     ptr.To(gwv1.Group("delegation.gateway.solo.io")), // custom group for delegation
						},
					},
				}
				backends = append(backends, backendRef)
				break
			}
		}
	}

	return backends
}

func (o *GatewayAPIOutput) generateFilterForURLRewrite(r *gloogwv1.Route, wrapper snapshot.Wrapper) *gwv1.HTTPRouteFilter {

	rf := &gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterURLRewrite,
		URLRewrite: &gwv1.HTTPURLRewriteFilter{
			Path: &gwv1.HTTPPathModifier{},
		},
	}
	match, err := isExactMatch(r.GetMatchers())
	if err != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers with different types in same route")
	}
	if match {
		rf.URLRewrite.Path.Type = gwv1.FullPathHTTPPathModifier
		rf.URLRewrite.Path.ReplaceFullPath = ptr.To(r.GetOptions().GetPrefixRewrite().GetValue())
		rf.URLRewrite.Path.ReplacePrefixMatch = nil
	}
	match, err = isPrefixMatch(r.GetMatchers())
	if err != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
	}

	if match {
		rf.URLRewrite.Path.Type = gwv1.PrefixMatchHTTPPathModifier
		rf.URLRewrite.Path.ReplacePrefixMatch = ptr.To(r.GetOptions().GetPrefixRewrite().GetValue())
		rf.URLRewrite.Path.ReplaceFullPath = nil
	}

	match, err = isRegexMatch(r.GetMatchers())
	if err != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers with different types in same route")
	}
	if match {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has regex matchers and cannot be used with path rewrites in Gateway API")
		return nil
	}
	// regex rewrite, NOT SUPPORTED IN GATEWAY API
	if r.GetOptions().GetRegexRewrite() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "regex rewrite not supported in Gateway API")
	}

	return rf
}

// Converts a single upstream to a GatewayAPI backend ref
func (o *GatewayAPIOutput) generateBackendRefForSingleUpstream(r *gloogwv1.Route, wrapper snapshot.Wrapper) gwv1.HTTPBackendRef {
	upstream := r.GetRouteAction().GetSingle().GetUpstream()
	var backendRef gwv1.HTTPBackendRef

	//TODO we need to lookup the upstream to see if its kube and then just reference kube directly
	var up *snapshot.UpstreamWrapper
	//if it is not a kube service or does not need http2
	var upstreamNs = upstream.GetNamespace()
	if upstreamNs == "" {
		upstreamNs = wrapper.GetNamespace()
	}

	up = o.edgeCache.GetUpstream(types.NamespacedName{Name: upstream.GetName(), Namespace: upstreamNs})

	if up == nil {
		// unknown reference to backend
		o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "upstream %s not found, referencing unknown upstream backend ref", upstream.GetName())

		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(upstream.GetName()),
					Namespace: (*gwv1.Namespace)(ptr.To(upstreamNs)),
					Kind:      (*gwv1.Kind)(ptr.To("Backend")),
					Group:     (*gwv1.Group)(ptr.To(glookube.GroupName)),
				},
			},
		}
	} else if up.Spec.GetKube() == nil {
		// non kubernetes upstream
		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(upstream.GetName()),
					Namespace: (*gwv1.Namespace)(ptr.To(upstreamNs)),
					Kind:      (*gwv1.Kind)(ptr.To("Backend")),
					Group:     (*gwv1.Group)(ptr.To(glookube.GroupName)),
				},
			},
		}
	} else if up.Spec.GetKube() != nil && up.Spec.GetUseHttp2() != nil && up.Spec.GetUseHttp2().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_UPDATE_OBJECT, wrapper, "service %s/%s uses http2, update its k8s service appProtocol=http2", up.Spec.GetKube().GetServiceNamespace(), up.Spec.GetKube().GetServiceName())
		// normal backend ref but let the user know htey need to annotate their service
		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(up.Spec.GetKube().GetServiceName()),
					Namespace: (*gwv1.Namespace)(ptr.To(up.Spec.GetKube().GetServiceNamespace())),
					Port:      ptr.To(gwv1.PortNumber(int32(up.Spec.GetKube().GetServicePort()))),
				},
			},
		}
	} else {
		//use kube backendref
		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(up.Spec.GetKube().GetServiceName()),
					Namespace: (*gwv1.Namespace)(ptr.To(up.Spec.GetKube().GetServiceNamespace())),
					Port:      ptr.To(gwv1.PortNumber(int32(up.Spec.GetKube().GetServicePort()))),
				},
			},
		}
	}

	// AWS lambda integration
	if r.GetRouteAction().GetSingle().GetDestinationSpec() != nil && r.GetRouteAction().GetSingle().GetDestinationSpec().GetAws() != nil {
		// we need to add a parameter for the lambda name reference
		backendRef.Filters = append(backendRef.Filters, gwv1.HTTPRouteFilter{
			Type: gwv1.HTTPRouteFilterExtensionRef,
			ExtensionRef: &gwv1.LocalObjectReference{
				Kind:  "Parameter",
				Group: glookube.GroupName,
				Name:  (gwv1.ObjectName)(r.GetRouteAction().GetSingle().GetDestinationSpec().GetAws().GetLogicalName()),
			},
		})
	}
	return backendRef
}

func (o *GatewayAPIOutput) convertMatch(m *matchers.Matcher, wrapper snapshot.Wrapper) (gwv1.HTTPRouteMatch, error) {
	hrm := gwv1.HTTPRouteMatch{
		QueryParams: []gwv1.HTTPQueryParamMatch{},
	}

	// header matching
	if len(m.GetHeaders()) > 0 {
		hrm.Headers = []gwv1.HTTPHeaderMatch{}
		for _, h := range m.GetHeaders() {
			// support invert header match https://github.com/solo-io/gloo/blob/main/projects/gateway2/translator/httproute/gateway_http_route_translator.go#L274
			if h.GetInvertMatch() == true {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "invert match not currently supported")
			}
			if h.GetRegex() {
				hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
					Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
					Value: h.GetValue(),
					Name:  gwv1.HTTPHeaderName(h.GetName()),
				})
			} else {
				if h.GetValue() == "" {
					// no header value set so any value is good
					hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Value: "*",
						Name:  gwv1.HTTPHeaderName(h.GetName()),
					})
				} else {
					hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Value: h.GetValue(),
						Name:  gwv1.HTTPHeaderName(h.GetName()),
					})
				}
			}
		}

	}

	// method matching
	if len(m.GetMethods()) > 0 {
		if len(m.GetMethods()) > 1 {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gateway API only supports 1 method match per rule and %d were detected", len(m.GetMethods()))
		}
		hrm.Method = (*gwv1.HTTPMethod)(ptr.To(strings.ToUpper(m.GetMethods()[0])))
	}

	// query param matching
	if len(m.GetQueryParameters()) > 0 {
		for _, m := range m.GetQueryParameters() {
			if m.GetRegex() {
				hrm.QueryParams = append(hrm.QueryParams, gwv1.HTTPQueryParamMatch{
					Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
					Name:  (gwv1.HTTPHeaderName)(m.GetName()),
					Value: m.GetValue(),
				})
			} else {
				hrm.QueryParams = append(hrm.QueryParams, gwv1.HTTPQueryParamMatch{
					Type:  ptr.To(gwv1.QueryParamMatchExact),
					Name:  (gwv1.HTTPHeaderName)(m.GetName()),
					Value: m.GetValue(),
				})
			}
		}
	}

	// Path matching
	if m.GetPathSpecifier() != nil {
		if m.GetPrefix() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchPathPrefix),
				Value: ptr.To(m.GetPrefix()),
			}
		}
		if m.GetExact() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchExact),
				Value: ptr.To(m.GetExact()),
			}
		}
		if m.GetRegex() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchRegularExpression),
				Value: ptr.To(m.GetRegex()),
			}
		}
	}
	return hrm, nil
}

func (o *GatewayAPIOutput) convertRouteTableToHTTPRoute(rt *snapshot.RouteTableWrapper) error {

	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gwv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rt.Name,
			Namespace: rt.Namespace,
			Labels:    rt.Labels,
		},
		Spec: gwv1.HTTPRouteSpec{
			// CommonRouteSpec: gwv1.CommonRouteSpec{},
			// Hostnames: [],
			Rules: []gwv1.HTTPRouteRule{},
		},
	}
	if rt.Spec.GetWeight() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, rt, "has weight set but there is no equivalent in Gateway API")
	}

	for _, route := range rt.Spec.GetRoutes() {
		rule, err := o.convertRouteToRule(route, rt)
		if err != nil {
			return err
		}
		hr.Spec.Rules = append(hr.Spec.Rules, rule)
	}
	o.gatewayAPICache.AddHTTPRoute(snapshot.NewHTTPRouteWrapper(hr, rt.FileOrigin()))

	return nil
}

// This function validates that the RouteRable matchers are the same match type prefix or exact
// The reason being is that if you are doing a rewrite you can only have one type of filter applied
func validateMatchersAreSame(matches []*matchers.Matcher) error {

	var foundExact, foundPrefix, foundRegex bool
	for _, m := range matches {
		if m.GetExact() != "" {
			if foundPrefix || foundRegex {
				return fmt.Errorf("multiple matchers found")
			}
			foundExact = true
		}
		if m.GetPrefix() != "" {
			if foundExact || foundRegex {
				return fmt.Errorf("multiple matchers found")
			}
			foundPrefix = true
		}
		if m.GetRegex() != "" {
			if foundExact || foundPrefix {
				return fmt.Errorf("multiple matchers found")
			}
			foundRegex = true
		}
	}
	return nil
}

// tests to see if all matchers are exact
func isExactMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetExact() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are exact
func isPrefixMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetPrefix() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are regex
func isRegexMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetRegex() != "" {
			return true, nil
		}
	}
	return false, nil
}

func doHttpRouteLabelsMatch(matches map[string]string, labels map[string]string) bool {
	for k, v := range matches {
		if labels[k] != v {
			return false
		}
	}
	return true
}

// TODO this is a mess
func isRouteOptionsSet(options *gloov1.RouteOptions) bool {
	return options.GetExtProc() != nil || options.GetCors() != nil || options.GetRetries() != nil || options.GetTimeout() != nil ||
		options.GetStagedTransformations() != nil || options.GetAutoHostRewrite() != nil ||
		options.GetFaults() != nil || options.GetExtensions() != nil || options.GetTracing() != nil || options.GetShadowing() != nil ||
		options.GetHeaderManipulation() != nil || options.GetAppendXForwardedHost() != nil || options.GetLbHash() != nil || options.GetUpgrades() != nil ||
		options.GetRatelimit() != nil || options.GetRatelimitBasic() != nil || options.GetWaf() != nil || options.GetJwtConfig() != nil || options.GetRbac() != nil ||
		options.GetDlp() != nil || options.GetStagedTransformations() != nil || options.GetEnvoyMetadata() != nil || options.GetMaxStreamDuration() != nil ||
		options.GetIdleTimeout() != nil || options.GetRegexRewrite() != nil || options.GetExtauth() != nil
}
