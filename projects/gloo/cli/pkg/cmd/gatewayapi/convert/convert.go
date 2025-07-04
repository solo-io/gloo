package convert

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"

	gloogwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaykube "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	v3 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
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
			o.gatewayAPICache.AddUpstream(upstream)
		}
	}

	for _, settings := range o.edgeCache.Settings() {
		o.gatewayAPICache.AddSettings(settings)
	}

	// copy over any existing options
	return nil
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
			vho, exists := o.gatewayAPICache.VirtualHostOptions[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
			if !exists {
				vho, exists = o.edgeCache.VirtualHostOptions()[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
				if !exists {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, vs, "references VirtualHostOption %s that does not exist", types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()})
					continue
				}
			}
			// add the target ref to the listener
			vho.VirtualHostOption.Spec.TargetRefs = append(vho.VirtualHostOption.Spec.GetTargetRefs(), &v3.PolicyTargetReferenceWithSectionName{
				Group:     apixv1a1.GroupName,
				Kind:      "XListenerSet",
				Name:      listenerSet.Name,
				Namespace: wrapperspb.String(listenerSet.Namespace),
			})
			o.gatewayAPICache.AddVirtualHostOption(snapshot.NewVirtualHostOptionWrapper(vho.VirtualHostOption, vs.FileOrigin()))
		}
	}

	// we need to get the virtualhostoptions and update their references
	if vs.Spec.GetVirtualHost().GetOptions() != nil {
		// create a separate virtualhost option and link it
		vho := &gatewaykube.VirtualHostOption{
			TypeMeta: metav1.TypeMeta{
				Kind:       "VirtualHostOption",
				APIVersion: gatewaykube.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      listenerSet.Name,
				Namespace: listenerSet.Namespace,
			},
			Spec: gloogwv1.VirtualHostOption{
				Options: vs.Spec.GetVirtualHost().GetOptions(),
				TargetRefs: []*v3.PolicyTargetReferenceWithSectionName{
					{
						Group:     apixv1a1.GroupName,
						Kind:      "XListenerSet",
						Name:      listenerSet.Name,
						Namespace: wrapperspb.String(listenerSet.Namespace),
					},
				},
			},
		}
		o.gatewayAPICache.AddVirtualHostOption(snapshot.NewVirtualHostOptionWrapper(vho, vs.FileOrigin()))
	}
	o.gatewayAPICache.AddListenerSet(snapshot.NewListenerSetWrapper(listenerSet, vs.FileOrigin()))

	return nil
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
			o.convertHTTPListenerOptions(glooGateway, proxyName)
		}
		if glooGateway.Spec.GetOptions() != nil && glooGateway.Spec.GetOptions() != nil {
			o.convertListenerOptions(glooGateway, proxyName)
		}
	}
}

func (o *GatewayAPIOutput) convertListenerOptions(glooGateway *snapshot.GlooGatewayWrapper, proxyName string) {
	options := glooGateway.Spec.GetOptions()
	listenerOption := &gatewaykube.ListenerOption{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ListenerOption",
			APIVersion: gatewaykube.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      glooGateway.GetName(),
			Namespace: glooGateway.GetNamespace(),
			Labels:    glooGateway.Gateway.Labels,
		},
		Spec: gloogwv1.ListenerOption{
			Options: options,
		},
	}

	listenerOption.Spec.TargetRefs = append(listenerOption.Spec.GetTargetRefs(), &v3.PolicyTargetReferenceWithSectionName{
		Group:     gwv1.GroupVersion.Group,
		Kind:      "Gateway",
		Name:      proxyName,
		Namespace: wrapperspb.String(glooGateway.GetNamespace()),
	})

	o.gatewayAPICache.AddListenerOption(snapshot.NewListenerOptionWrapper(listenerOption, glooGateway.FileOrigin()))
}

func (o *GatewayAPIOutput) convertHTTPListenerOptions(glooGateway *snapshot.GlooGatewayWrapper, proxyName string) {
	options := glooGateway.Spec.GetHttpGateway().GetOptions()
	listenerOption := &gatewaykube.HttpListenerOption{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HttpListenerOption",
			APIVersion: gatewaykube.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      glooGateway.GetName(),
			Namespace: glooGateway.GetNamespace(),
			Labels:    glooGateway.Gateway.Labels,
		},
		Spec: gloogwv1.HttpListenerOption{
			Options: options,
		},
	}

	listenerOption.Spec.TargetRefs = append(listenerOption.Spec.GetTargetRefs(), &v3.PolicyTargetReferenceWithSectionName{
		Group:     gwv1.GroupVersion.Group,
		Kind:      "Gateway",
		Name:      proxyName,
		Namespace: wrapperspb.String(glooGateway.GetNamespace()),
	})

	o.gatewayAPICache.AddHTTPListenerOption(snapshot.NewHTTPListenerOptionWrapper(listenerOption, glooGateway.FileOrigin()))
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
) (*gatewaykube.RouteOption, *gwv1.HTTPRouteFilter) {

	var ro *gatewaykube.RouteOption
	var filter *gwv1.HTTPRouteFilter
	associationID := RandStringRunes(RandomSuffix)
	if routeName == "" {
		routeName = "route-association"
	}
	associationName := fmt.Sprintf("%s-%s", routeName, associationID)

	// converts options to RouteOptions but we need to this for everything except prefixrewrite
	if isRouteOptionsSet(options) {
		ro = &gatewaykube.RouteOption{
			TypeMeta: metav1.TypeMeta{
				Kind:       "RouteOption",
				APIVersion: gatewaykube.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      associationName,
				Namespace: wrapper.GetNamespace(),
			},
			Spec: gloogwv1.RouteOption{
				Options: options,
			},
		}

		// Because we move rewrites to a filter we need to remove it from RouteOptions
		if options.GetPrefixRewrite() != nil {
			ro.Spec.GetOptions().PrefixRewrite = nil
		}

		filter = &gwv1.HTTPRouteFilter{
			Type: gwv1.HTTPRouteFilterExtensionRef,
			ExtensionRef: &gwv1.LocalObjectReference{
				Group: glookube.GroupName,
				Kind:  "RouteOption",
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
	}
	return ro, filter
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

		ro, filter := o.convertRouteOptions(options, r.GetName(), wrapper)
		if filter != nil {
			rr.Filters = append(rr.Filters, *filter)
		}
		if ro != nil {
			o.gatewayAPICache.AddRouteOption(snapshot.NewRouteOptionWrapper(ro, wrapper.FileOrigin()))
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
			//TODO(nick) what if route name is nil?
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
			rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
				Type: gwv1.HTTPRouteFilterExtensionRef,
				ExtensionRef: &gwv1.LocalObjectReference{
					Group: glookube.GroupName,
					Kind:  "RouteOption",
					Name:  gwv1.ObjectName(delegateOptions.GetName()),
				},
			})
			// grab that route option and add it to the cache
			ro, exists := o.edgeCache.RouteOptions()[types.NamespacedName{Name: delegateOptions.GetName(), Namespace: delegateOptions.GetNamespace()}]
			if !exists {
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find RouteOption %s/%s for delegated route option reference", delegateOptions.GetNamespace(), delegateOptions.GetName())
			}
			o.gatewayAPICache.AddRouteOption(ro)

			if ro.Spec.GetOptions() != nil && ro.Spec.GetOptions().GetExtauth() != nil && ro.Spec.GetOptions().GetExtauth().GetConfigRef() != nil {
				// we need to copy over the auth config ref if it exists
				ref := ro.Spec.GetOptions().GetExtauth().GetConfigRef()
				ac, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
				if !exists {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, ro, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
				}
				o.gatewayAPICache.AddAuthConfig(ac)
			}
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
func convertDirectResponse(action *gloov1.DirectResponseAction) *v1alpha1.DirectResponse {
	if action == nil {
		return nil
	}
	dr := &v1alpha1.DirectResponse{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DirectResponse",
			APIVersion: v1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1alpha1.DirectResponseSpec{
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
					Kind:      (*gwv1.Kind)(ptr.To("Upstream")),
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
					Kind:      (*gwv1.Kind)(ptr.To("Upstream")),
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
