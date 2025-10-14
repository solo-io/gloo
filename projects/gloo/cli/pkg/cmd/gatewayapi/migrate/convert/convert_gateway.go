package convert

import (
	"encoding/json"
	"fmt"
	"strings"

	kgateway "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	gloogateway "github.com/solo-io/gloo-gateway/api/v1alpha1"
	glooratelimit "github.com/solo-io/gloo-gateway/external/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	v4 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/csrf/v3"
	gloo_type_matcher "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/protocol_upgrade"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tracing"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

func (g *GatewayAPIOutput) convertVirtualServices(glooGateway *snapshot.GlooGatewayWrapper, useListenersets bool) error {

	gatewayVs, err := g.edgeCache.GlooGatewayVirtualServices(glooGateway)
	if err != nil {
		return err
	}
	if len(gatewayVs) == 0 {
		g.AddErrorFromWrapper(ERROR_TYPE_NO_REFERENCES, glooGateway, "gateway does not contain virtual services")
	}
	for _, vs := range gatewayVs {
		proxyNames := glooGateway.Spec.GetProxyNames()
		if len(proxyNames) == 0 {
			proxyNames = append(proxyNames, "gateway-proxy")
		}
		for _, proxyName := range proxyNames {
			listenerName := fmt.Sprintf("%s-%d-%s-%s", proxyName, glooGateway.Spec.GetBindPort(), vs.Name, vs.Namespace)

			var policyRef kgateway.LocalPolicyTargetReference
			var parentRef gwv1.ParentReference
			if useListenersets {
				policyRef, err = g.addListenerSetToCache(glooGateway, listenerName, vs, proxyName)
				if err != nil {
					return err
				}
				parentRef = gwv1.ParentReference{
					Name:      policyRef.Name,
					Namespace: ptr.To(gwv1.Namespace(glooGateway.GetNamespace())),
					Kind:      ptr.To(policyRef.Kind),
					Group:     ptr.To(policyRef.Group),
				}
				// we need to grab the gateway and update its listener and add dummy one
				gateway := g.GetGatewayAPICache().GetGateway(types.NamespacedName{
					Namespace: glooGateway.GetNamespace(),
					Name:      proxyName,
				})
				// this is to get around a bug in the spec that requires 1 listener
				gateway.Spec.Listeners = []gwv1.Listener{
					{
						Name:     "dummy",
						Port:     8123,
						Protocol: "HTTP",
						AllowedRoutes: &gwv1.AllowedRoutes{
							Namespaces: &gwv1.RouteNamespaces{
								From: ptr.To(gwv1.NamespacesFromSelector),
								Selector: &metav1.LabelSelector{
									MatchLabels: map[string]string{
										"dummy": "dummy",
									},
								},
							},
						},
					},
				}
				g.GetGatewayAPICache().AddGateway(gateway)

			} else {

				// we need to generate listeners for the Gateway here and set them
				gateway := g.GetGatewayAPICache().GetGateway(types.NamespacedName{
					Namespace: glooGateway.GetNamespace(),
					Name:      proxyName,
				})

				listeners, err := g.generateListenersForGateway(vs, glooGateway, listenerName)
				if err != nil {
					return err
				}
				gateway.Spec.Listeners = append(gateway.Spec.Listeners, listeners...)

				g.GetGatewayAPICache().AddGateway(gateway)
				// we need the policy ref to be the gateway
				policyRef = kgateway.LocalPolicyTargetReference{
					Group: gwv1.GroupName,
					Kind:  "Gateway",
					Name:  gwv1.ObjectName(gateway.GetName()),
				}
				parentRef = gwv1.ParentReference{
					Name:      gwv1.ObjectName(gateway.GetName()),
					Namespace: ptr.To(gwv1.Namespace(gateway.GetNamespace())),
					Kind:      ptr.To(gwv1.Kind("Gateway")),
					Group:     ptr.To(gwv1.Group(gwv1.GroupName)),
				}
			}

			if vs.Spec.GetVirtualHost().GetOptionsConfigRefs() != nil && len(vs.Spec.GetVirtualHost().GetOptionsConfigRefs().GetDelegateOptions()) > 0 {
				g.convertDelegateOptions(vs, vs.Spec.GetVirtualHost().GetOptionsConfigRefs().GetDelegateOptions(), policyRef)
			}

			// we need to get the virtualhostoptions and update their references
			if vs.Spec.GetVirtualHost().GetOptions() != nil {
				g.convertInlineOptions(vs, policyRef)
			}

			// convert the routing portion of the virtual service
			err := g.convertVirtualServiceHTTPRoutes(vs, glooGateway, listenerName, parentRef)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *GatewayAPIOutput) convertInlineOptions(vs *snapshot.VirtualServiceWrapper, policyRef kgateway.LocalPolicyTargetReference) {
	// create a separate virtualhost option and link it
	gtp := &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vs.Name,
			Namespace: vs.Namespace,
		},
	}
	// go through each option and add it to traffic policy
	spec := g.convertVHOOptionsToTrafficPolicySpec(vs.Spec.GetVirtualHost().GetOptions(), vs)

	// attach the xListenerSet to the GlooTrafficPolicy
	if spec.TrafficPolicySpec.TargetRefs == nil {
		spec.TrafficPolicySpec.TargetRefs = []kgateway.LocalPolicyTargetReferenceWithSectionName{}
	}
	// add the policy reference to the GTP (could be Gateway or ListenerSet)
	spec.TrafficPolicySpec.TargetRefs = append(spec.TrafficPolicySpec.TargetRefs, kgateway.LocalPolicyTargetReferenceWithSectionName{
		LocalPolicyTargetReference: policyRef,
	})

	gtp.Spec = spec

	g.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(gtp, vs.FileOrigin()))
}

// convertDelegateOptions - finds and converts the delegated option if it does not exist
func (g *GatewayAPIOutput) convertDelegateOptions(vs *snapshot.VirtualServiceWrapper, delegateOptions []*core.ResourceRef, policyRef kgateway.LocalPolicyTargetReference) {
	for _, delegateOption := range delegateOptions {
		// check to see if this already exists in gatewayAPI cache, if not move it over from edge cache
		gtp, exists := g.gatewayAPICache.GlooTrafficPolicies[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
		if !exists {
			vho, exists := g.edgeCache.VirtualHostOptions()[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
			if !exists {
				g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, vs, "references VirtualHostOption %s that does not exist", types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()})
				continue
			}
			gtp = g.convertVirtualHostOptionToGlooTrafficPolicy(vho)
			if gtp == nil {
				g.AddErrorFromWrapper(ERROR_TYPE_IGNORED, vs, "references VirtualHostOption %s - No options converted", types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()})
			}
		}
		if vs.Namespace != gtp.GlooTrafficPolicy.GetNamespace() {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "VirtualHostOption %s references a listener set in a different namespace %s which is not supported", types.NamespacedName{Name: vs.GetName(), Namespace: vs.GetNamespace()}, types.NamespacedName{Name: gtp.GlooTrafficPolicy.GetName(), Namespace: gtp.GlooTrafficPolicy.GetNamespace()})
		}
		// add the target ref to the listener
		gtp.GlooTrafficPolicy.Spec.TargetRefs = append(gtp.GlooTrafficPolicy.Spec.TargetRefs, kgateway.LocalPolicyTargetReferenceWithSectionName{
			LocalPolicyTargetReference: policyRef,
		})
		g.gatewayAPICache.AddGlooTrafficPolicy(gtp)
	}
}

func (g *GatewayAPIOutput) addListenerSetToCache(glooGateway *snapshot.GlooGatewayWrapper, listenerName string, vs *snapshot.VirtualServiceWrapper, proxyName string) (kgateway.LocalPolicyTargetReference, error) {
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
				Name:      gwv1.ObjectName(proxyName),
			},
			Listeners: []apixv1a1.ListenerEntry{},
		},
	}
	// convert the listener portion of the virtual service
	listenerEntries, err := g.generateListenerEntries(vs, glooGateway, listenerName)
	if err != nil {
		return kgateway.LocalPolicyTargetReference{}, err
	}
	listenerSet.Spec.Listeners = listenerEntries
	g.gatewayAPICache.AddListenerSet(snapshot.NewListenerSetWrapper(listenerSet, vs.FileOrigin()))

	ref := kgateway.LocalPolicyTargetReference{
		Group: apixv1a1.GroupName,
		Kind:  "XListenerSet",
		Name:  gwv1.ObjectName(listenerName),
	}

	return ref, nil
}
func (g *GatewayAPIOutput) generateListenersForGateway(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string) ([]gwv1.Listener, error) {
	var listeners []gwv1.Listener
	for _, hostname := range vs.Spec.GetVirtualHost().GetDomains() {
		if strings.Contains(hostname, ":") {
			g.AddErrorFromWrapper(ERROR_TYPE_IGNORED, vs, "contains port in hostname %s, its being ignored for ListenerSet %s/%s", hostname, listenerName, vs.GetNamespace())
			continue
		}

		// listener entry does not support wildcard
		listenerEntryName := strings.ReplaceAll(fmt.Sprintf("%s-%s", vs.Name, hostname), "*", "star")
		entry := gwv1.Listener{
			Name:     gwv1.SectionName(listenerEntryName),
			Hostname: ptr.To(gwv1.Hostname(hostname)),
			Port:     apixv1a1.PortNumber(glooGateway.Spec.GetBindPort()),
			Protocol: gwv1.HTTPProtocolType,
		}
		if vs.Spec.GetSslConfig() != nil {
			tlsConfig := g.generateTLSConfiguration(vs)
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
		listeners = append(listeners, entry)
	}
	return listeners, nil
}

func (g *GatewayAPIOutput) generateListenerEntries(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string) ([]apixv1a1.ListenerEntry, error) {
	var listeners []apixv1a1.ListenerEntry
	// we only create the listener part, not the http matchers
	for _, hostname := range vs.Spec.GetVirtualHost().GetDomains() {
		if strings.Contains(hostname, ":") {
			g.AddErrorFromWrapper(ERROR_TYPE_IGNORED, vs, "contains port in hostname %s, its being ignored for ListenerSet %s/%s", hostname, listenerName, vs.GetNamespace())
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
			tlsConfig := g.generateTLSConfiguration(vs)
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
		listeners = append(listeners, entry)
	}
	return listeners, nil
}

func (g *GatewayAPIOutput) convertVirtualHostOptionToGlooTrafficPolicy(vho *snapshot.VirtualHostOptionWrapper) *snapshot.GlooTrafficPolicyWrapper {

	policy := &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vho.GetName(),
			Namespace: vho.GetNamespace(),
		},
	}

	policy.Spec = g.convertVHOOptionsToTrafficPolicySpec(vho.VirtualHostOption.Spec.GetOptions(), vho)

	wrapper := snapshot.NewGlooTrafficPolicyWrapper(policy, vho.FileOrigin())
	return wrapper
}

func (g *GatewayAPIOutput) convertVHOOptionsToTrafficPolicySpec(vho *gloov1.VirtualHostOptions, wrapper snapshot.Wrapper) gloogateway.GlooTrafficPolicySpec {

	spec := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs:      nil, // existing
			TargetSelectors: nil, // existing
			AI:              nil, // existing
			Transformation:  nil, // existing
			ExtProc:         nil, // existing
			ExtAuth:         nil, // existing
			RateLimit:       nil, // existing
			Cors:            nil, // existing
			Csrf:            nil, // existing
			HeaderModifiers: nil,
			AutoHostRewrite: nil, // existing
			Buffer:          nil, // existing
			Timeouts:        nil, // existing
			Retry:           nil, // existing
			RBAC:            nil,
		},
		GlooRateLimit:      nil, // existing
		GlooExtAuth:        nil, // existing
		GlooTransformation: nil, // existing
		GlooJWT:            nil, // existing
		GlooRBAC:           nil, // existing
	}
	if vho != nil {
		if vho.GetExtauth() != nil {
			// we need to copy over the auth config ref if it exists
			ref := vho.GetExtauth().GetConfigRef()
			ac, exists := g.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
			if !exists {
				g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
			} else {
				ac.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   "extauth.solo.io",
					Version: "v1",
					Kind:    "AuthConfig",
				})

				g.gatewayAPICache.AddAuthConfig(ac)

				spec.GlooExtAuth = &gloogateway.GlooExtAuth{
					ExtensionRef: &kgateway.NamespacedObjectReference{
						Name:      "ext-authz",
						Namespace: ptr.To(gwv1.Namespace(wrapper.GetNamespace())),
					},
					AuthConfigRef: &gloogateway.AuthConfigRef{
						Name:      gwv1.ObjectName(vho.GetExtauth().GetConfigRef().GetName()),
						Namespace: ptr.To(gwv1.Namespace(vho.GetExtauth().GetConfigRef().GetNamespace())),
					},
				}
			}
		}
		if vho.GetExtProc() != nil {
			// TODO(nick) the extproc on the VHO allows a user to disable the global
			// one but idk if there is an equivalent in gateway api?
			//extProc := &kgateway.ExtProcPolicy{
			//	ExtensionRef:   &corev1.LocalObjectReference{Name: "global-ext-proc"},
			//	ProcessingMode: &kgateway.ProcessingMode{
			//		RequestHeaderMode:   nil,
			//		ResponseHeaderMode:  nil,
			//		RequestBodyMode:     nil,
			//		ResponseBodyMode:    nil,
			//		RequestTrailerMode:  nil,
			//		ResponseTrailerMode: nil,
			//	},
			//}
			//
			//if vho.GetExtProc().GetDisabled() != nil {
			//
			//}
			//if vho.GetExtProc().GetOverride() != nil {}
			//
			//
			//spec.TrafficPolicySpec.ExtProc = extProc
		}
		if vho.GetWaf() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF is not supported")

			//waf := &gloogateway.Waf{
			//	Disabled:      ptr.To(vho.GetWaf().GetDisabled()),
			//	CustomMessage: ptr.To(vho.GetWaf().GetCustomInterventionMessage()),
			//	Rules:         []gloogateway.WafRule{},
			//}
			//for _, r := range vho.GetWaf().GetRuleSets() {
			//	waf.Rules = append(waf.Rules, gloogateway.WafRule{
			//		RuleStr: ptr.To(r.GetRuleStr()),
			//	})
			//	if r.GetFiles() != nil && len(r.GetFiles()) > 0 {
			//		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF files is not supported")
			//	}
			//}
			//spec.Waf = waf
		}
		if vho.GetRatelimitBasic() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimitBasic is not supported")
		}
		if vho.GetRatelimitEarly() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimitEarly is not supported, defaulting to regular rate limiting")

			// We need to create a glooratelimit.RateLimitConfig and reference it to the GTP
			rlcName := fmt.Sprintf("ratelimit-%s", RandStringRunes(4))
			rlc := &glooratelimit.RateLimitConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "RateLimitConfig",
					APIVersion: glooratelimit.RateLimitConfigGVK.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      rlcName,
					Namespace: wrapper.GetNamespace(),
				},
				Spec: glooratelimit.RateLimitConfigSpec{
					//ConfigType:
				},
			}

			//TODO(nick) where do we need to set descriptors?
			raw := &glooratelimit.RateLimitConfigSpec_Raw_{
				Raw: &glooratelimit.RateLimitConfigSpec_Raw{
					Descriptors:    nil,
					RateLimits:     []*glooratelimit.RateLimitActions{},
					SetDescriptors: nil,
				},
			}
			for _, rl := range vho.GetRatelimitEarly().GetRateLimits() {
				rateLimit := &glooratelimit.RateLimitActions{
					Actions:    []*glooratelimit.Action{},
					SetActions: []*glooratelimit.Action{},
				}
				for _, action := range rl.GetActions() {
					rateLimitAction := g.convertRateLimitAction(action)
					rateLimit.Actions = append(rateLimit.Actions, rateLimitAction)
				}
				for _, action := range rl.GetSetActions() {
					rateLimitAction := g.convertRateLimitAction(action)
					rateLimit.SetActions = append(rateLimit.SetActions, rateLimitAction)
				}
				if rl.GetLimit() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimit action limit is not supported")
				}
			}
			rlc.Spec = glooratelimit.RateLimitConfigSpec{
				ConfigType: raw,
			}

			g.gatewayAPICache.AddRateLimitConfigs(snapshot.NewRateLimitConfigPolicyWrapper(rlc, wrapper.FileOrigin()))

			spec.GlooRateLimit = &gloogateway.GlooRateLimit{
				Global: &gloogateway.GlobalRateLimit{
					// Need to find the Gateway Extension for Global Rate Limit Server
					ExtensionRef: &kgateway.NamespacedObjectReference{
						Name:      "rate-limit",
						Namespace: ptr.To(gwv1.Namespace(wrapper.GetNamespace())),
					},
					// RateLimitConfig for the policy, not sure how it works for rate limit basic
					RateLimitConfigRef: gloogateway.RateLimitConfigRef{
						Name:      gwv1.ObjectName(rlc.GetName()),
						Namespace: ptr.To(gwv1.Namespace(rlc.GetNamespace())),
					},
				},
			}
		}
		if vho.GetRatelimitRegular() != nil {
			// We need to create a glooratelimit.RateLimitConfig and reference it to the GTP
			rlcName := fmt.Sprintf("ratelimit-%s", RandStringRunes(4))
			rlc := &glooratelimit.RateLimitConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       "RateLimitConfig",
					APIVersion: glooratelimit.RateLimitConfigGVK.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      rlcName,
					Namespace: wrapper.GetNamespace(),
				},
				Spec: glooratelimit.RateLimitConfigSpec{
					//ConfigType:
				},
			}

			//TODO(nick) where do we need to set descriptors?
			raw := &glooratelimit.RateLimitConfigSpec_Raw_{
				Raw: &glooratelimit.RateLimitConfigSpec_Raw{
					Descriptors:    nil,
					RateLimits:     []*glooratelimit.RateLimitActions{},
					SetDescriptors: nil,
				},
			}
			for _, rl := range vho.GetRatelimitRegular().GetRateLimits() {
				rateLimit := &glooratelimit.RateLimitActions{
					Actions:    []*glooratelimit.Action{},
					SetActions: []*glooratelimit.Action{},
				}
				for _, action := range rl.GetActions() {
					rateLimitAction := g.convertRateLimitAction(action)
					rateLimit.Actions = append(rateLimit.Actions, rateLimitAction)
				}
				for _, action := range rl.GetSetActions() {
					rateLimitAction := g.convertRateLimitAction(action)
					rateLimit.SetActions = append(rateLimit.SetActions, rateLimitAction)
				}
				if rl.GetLimit() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimit action limit is not supported")
				}
			}
			rlc.Spec = glooratelimit.RateLimitConfigSpec{
				ConfigType: raw,
			}

			g.gatewayAPICache.AddRateLimitConfigs(snapshot.NewRateLimitConfigPolicyWrapper(rlc, wrapper.FileOrigin()))

			spec.GlooRateLimit = &gloogateway.GlooRateLimit{
				Global: &gloogateway.GlobalRateLimit{
					// Need to find the Gateway Extension for Global Rate Limit Server
					ExtensionRef: &kgateway.NamespacedObjectReference{
						Name:      "rate-limit",
						Namespace: ptr.To(gwv1.Namespace(wrapper.GetNamespace())),
					},
					// RateLimitConfig for the policy, not sure how it works for rate limit basic
					// TODO(nick) grab the global rate limit config ref
					RateLimitConfigRef: gloogateway.RateLimitConfigRef{},
				},
			}
		}
		if vho.GetHeaderManipulation() != nil {
			// this is natively supported on the HTTPRoute
		}
		if vho.GetCors() != nil {
			policy := g.convertCORS(vho.GetCors(), wrapper)
			spec.Cors = policy
		}
		if vho.GetTransformations() != nil {
			// TODO(nick) should we try to translate this or require the end user to migrate to staged?
		}
		if vho.GetStagedTransformations() != nil {
			transformation := g.convertStagedTransformation(vho.GetStagedTransformations(), wrapper)
			spec.GlooTransformation = transformation
		}
		if vho.GetJwt() != nil {
			// TODO(nick) should we try to translate this or require the end user to migrate to staged?
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "jwt is deprecated in edge and not supported")
		}
		if vho.GetJwtStaged() != nil {
			spec.GlooJWT = &gloogateway.StagedJWT{
				AfterExtAuth:  nil, // existing
				BeforeExtAuth: nil, // existing
			}
			if vho.GetJwtStaged().GetBeforeExtAuth() != nil {
				jwte := g.convertJWTStagedExtAuth(vho.GetJwtStaged().GetBeforeExtAuth(), wrapper)
				spec.GlooJWT.BeforeExtAuth = jwte
			}
			if vho.GetJwtStaged().GetAfterExtAuth() != nil {
				jwte := g.convertJWTStagedExtAuth(vho.GetJwtStaged().GetAfterExtAuth(), wrapper)
				spec.GlooJWT.BeforeExtAuth = jwte
			}
		}
		if vho.GetRbac() != nil {
			rbe := g.convertRBAC(vho.GetRbac())
			spec.GlooRBAC = rbe
		}
		if vho.GetBufferPerRoute() != nil && vho.GetBufferPerRoute().GetBuffer() != nil && vho.GetBufferPerRoute().GetBuffer().GetMaxRequestBytes() != nil {
			spec.Buffer = &kgateway.Buffer{
				MaxRequestSize: resource.NewQuantity(int64(vho.GetBufferPerRoute().GetBuffer().GetMaxRequestBytes().GetValue()), resource.BinarySI),
			}
			if vho.GetBufferPerRoute().GetDisabled() {
				spec.Buffer.Disable = &kgateway.PolicyDisable{}
			}
		}
		if vho.GetIncludeRequestAttemptCount() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "includeRequestAttemptCount is not supported")
		}
		if vho.GetIncludeAttemptCountInResponse() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "includeRequestAttemptCountInResponse is not supported")
		}
		if vho.GetCorsPolicyMergeSettings() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "corsPolicyMergeSettings is not supported")
		}
		if vho.GetDlp() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported")
		}
		if vho.GetCsrf() != nil {
			csrf := g.convertCSRF(vho.GetCsrf())
			spec.TrafficPolicySpec.Csrf = csrf
		}
		if vho.GetExtensions() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported")
		}
	}
	return spec
}

func (g *GatewayAPIOutput) generateTLSConfiguration(vs *snapshot.VirtualServiceWrapper) *gwv1.GatewayTLSConfig {
	tlsConfig := &gwv1.GatewayTLSConfig{
		Mode: ptr.To(gwv1.TLSModeTerminate),
		//FrontendValidation: nil, // TODO(nick) do we need to set this?
		//Options:            nil, // TODO(nick) do we need to set this?
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
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has SSLFiles but its not supported in Gateway API")
	}
	if vs.Spec.GetSslConfig().GetSds() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has SDS Certificates but its not supported in Gateway API")
	}
	if len(vs.Spec.GetSslConfig().GetVerifySubjectAltName()) > 0 {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has VerifySubjectAltName but its not supported in Gateway API")
	}
	if len(vs.Spec.GetSslConfig().GetAlpnProtocols()) > 0 {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has AlpnProtocols but its not supported in Gateway API")
	}
	if vs.Spec.GetSslConfig().GetOcspStaplePolicy() > 0 {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has OcspStaplePolicy %d but its not supported in Gateway API", vs.Spec.GetSslConfig().GetOcspStaplePolicy())
	}

	return tlsConfig
}

func (g *GatewayAPIOutput) generateGatewaysFromProxyNames(glooGateway *snapshot.GlooGatewayWrapper) {

	proxyNames := glooGateway.Gateway.Spec.GetProxyNames()

	if len(proxyNames) == 0 {
		proxyNames = append(proxyNames, "gateway-proxy")
	}

	for _, proxyName := range proxyNames {
		// check to see if we already created the Gateway, if we did then just move on
		existingGw := g.gatewayAPICache.GetGateway(types.NamespacedName{Name: proxyName, Namespace: glooGateway.Gateway.Namespace})
		if existingGw == nil {
			// create a new gateway
			existingGw = snapshot.NewGatewayWrapper(&gwv1.Gateway{
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
					Listeners:        []gwv1.Listener{},
					GatewayClassName: "gloo-gateway-v2",
				},
			}, glooGateway.FileOrigin())
		}

		// special case for per connection buffer limits to apply to the gateway as an annotation - https://github.com/kgateway-dev/kgateway/pull/11505
		if glooGateway.Spec.GetOptions() != nil && glooGateway.Spec.GetOptions().GetPerConnectionBufferLimitBytes() != nil && glooGateway.Spec.GetOptions().GetPerConnectionBufferLimitBytes().GetValue() != 0 {
			if existingGw.Annotations == nil {
				existingGw.Annotations = make(map[string]string)
			}
			existingGw.Annotations[perConnectionBufferLimit] = glooGateway.Spec.GetOptions().GetPerConnectionBufferLimitBytes().String()
		}

		g.gatewayAPICache.AddGateway(existingGw)

		if glooGateway.Spec.GetHttpGateway() != nil && glooGateway.Spec.GetHttpGateway().GetOptions() != nil {
			g.convertHTTPListenerOptions(glooGateway.Spec.GetHttpGateway().GetOptions(), glooGateway, proxyName)
		}
		if glooGateway.Spec.GetOptions() != nil && glooGateway.Spec.GetOptions() != nil {
			g.convertListenerOptions(glooGateway, proxyName)
		}
	}
}

func (g *GatewayAPIOutput) convertListenerOptions(glooGateway *snapshot.GlooGatewayWrapper, proxyName string) {
	options := glooGateway.Spec.GetOptions()
	if options == nil {
		return
	}
	listenerPolicy := &kgateway.HTTPListenerPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPListenerPolicy",
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
			TargetSelectors:            nil, // existing
			AccessLog:                  nil, // existing
			Tracing:                    nil, // existing
			UpgradeConfig:              nil, // existing
			UseRemoteAddress:           nil, // existing
			XffNumTrustedHops:          nil, // existing
			ServerHeaderTransformation: nil, // existing
			StreamIdleTimeout:          nil, // existing
			HealthCheck:                nil, // existing
			PreserveHttp1HeaderCase:    nil, // existing
		},
	}
	if options.GetExtensions() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option extensions are not supported for HTTPTrafficPolicy")
	}
	if options.GetSocketOptions() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option socket options are not supported for HTTPTrafficPolicy")
	}
	if options.GetAccessLoggingService() != nil {
		g.convertListenerOptionAccessLogging(glooGateway, options, listenerPolicy)
	}
	if options.GetListenerAccessLoggingService() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option listenerAccessLoggingService is not supported for HTTPTrafficPolicy")
	}
	if options.GetConnectionBalanceConfig() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option connectionBalanceConfig is not supported for HTTPTrafficPolicy")
	}
	//if options.GetPerConnectionBufferLimitBytes() != nil {
	// This is now set as an annotation on Gateway
	//	o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option perConnectionBufferLimitBytes is not supported for HTTPTrafficPolicy")
	//}
	if options.GetProxyProtocol() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option proxyProtocol is not supported for HTTPTrafficPolicy")
	}
	if options.GetTcpStats() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option tcpStats is not supported for HTTPTrafficPolicy")
	}

	g.gatewayAPICache.AddHTTPListenerPolicy(snapshot.NewHTTPListenerPolicyWrapper(listenerPolicy, glooGateway.FileOrigin()))
}

func (g *GatewayAPIOutput) convertListenerOptionAccessLogging(glooGateway *snapshot.GlooGatewayWrapper, options *gloov1.ListenerOptions, listenerPolicy *kgateway.HTTPListenerPolicy) {
	accessLoggingService := options.GetAccessLoggingService()

	for _, edgeAccessLog := range accessLoggingService.GetAccessLog() {
		if listenerPolicy.Spec.AccessLog == nil {
			listenerPolicy.Spec.AccessLog = []kgateway.AccessLog{}
		}
		accessLog := kgateway.AccessLog{
			FileSink:      nil, // existing
			GrpcService:   nil, // existing
			Filter:        nil, // existing
			OpenTelemetry: nil, //existing
		}
		if edgeAccessLog.GetFileSink() != nil {
			fileSink := &kgateway.FileSink{
				Path: edgeAccessLog.GetFileSink().GetPath(),
			}
			if jsonFormat := edgeAccessLog.GetFileSink().GetJsonFormat(); jsonFormat != nil {
				jsonBytes, err := json.Marshal(jsonFormat.AsMap())
				if err != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_IGNORED, glooGateway, "unable to marshal json format for accessLoggingService %v", err)
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
			accessLog.GrpcService = &kgateway.AccessLogGrpcService{
				CommonAccessLogGrpcService: kgateway.CommonAccessLogGrpcService{
					//CommonGrpcService: nil,// TODO(nick) what do we need to set here?
					LogName: edgeAccessLog.GetGrpcService().GetLogName(),
				},
				AdditionalRequestHeadersToLog:   edgeAccessLog.GetGrpcService().GetAdditionalRequestHeadersToLog(),
				AdditionalResponseHeadersToLog:  edgeAccessLog.GetGrpcService().GetAdditionalResponseHeadersToLog(),
				AdditionalResponseTrailersToLog: edgeAccessLog.GetGrpcService().GetAdditionalResponseTrailersToLog(),
			}

			// backend Ref
			switch edgeAccessLog.GetGrpcService().GetServiceRef().(type) {
			case *als.GrpcService_StaticClusterName:
				accessLog.GrpcService.BackendRef = &gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      gwv1.ObjectName(edgeAccessLog.GetGrpcService().GetStaticClusterName()),
						Namespace: nil,
						Port:      nil,
					},
				}
				g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, glooGateway, "", edgeAccessLog.GetGrpcService().GetStaticClusterName())
			}
			if edgeAccessLog.GetGrpcService().GetFilterStateObjectsToLog() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "accessLoggingService grpcService has filterStateObjectsToLog but its not supported")
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
					accessLog.Filter.AndFilter = append(accessLog.Filter.AndFilter, *g.convertAccessLogFitler(filter, glooGateway))
				}
			} else if edgeAccessLog.GetFilter().GetAndFilter() != nil {
				if accessLog.Filter.AndFilter == nil {
					accessLog.Filter.AndFilter = []kgateway.FilterType{}
				}
				for _, filter := range edgeAccessLog.GetFilter().GetAndFilter().GetFilters() {
					accessLog.Filter.AndFilter = append(accessLog.Filter.AndFilter, *g.convertAccessLogFitler(filter, glooGateway))
				}
			} else {
				// just and inline filter
				accessLog.Filter.FilterType = g.convertAccessLogFitler(edgeAccessLog.GetFilter(), glooGateway)
			}
		}
		listenerPolicy.Spec.AccessLog = append(listenerPolicy.Spec.AccessLog, accessLog)
	}
}

func (g *GatewayAPIOutput) convertAccessLogFitler(filter *als.AccessLogFilter, wrapper snapshot.Wrapper) *kgateway.FilterType {

	filterType := &kgateway.FilterType{
		StatusCodeFilter:     nil,   // existing
		DurationFilter:       nil,   // existing
		NotHealthCheckFilter: false, // existing
		TraceableFilter:      false, // existing
		HeaderFilter:         nil,   // existing
		ResponseFlagFilter:   nil,   // existing
		GrpcStatusFilter:     nil,   // existing
		CELFilter:            nil,   // existing
	}

	if filter.GetDurationFilter() != nil {
		filterType.DurationFilter = &kgateway.DurationFilter{
			Op:    kgateway.Op(filter.GetDurationFilter().GetComparison().GetOp()),
			Value: filter.GetDurationFilter().GetComparison().GetValue().GetDefaultValue(),
		}
	}
	if filter.GetHeaderFilter() != nil && filter.GetHeaderFilter().GetHeader() != nil {
		headerMatch := gwv1.HTTPHeaderMatch{
			Name: gwv1.HTTPHeaderName(filter.GetHeaderFilter().GetHeader().GetName()),
		}

		if filter.GetHeaderFilter().GetHeader().GetExactMatch() != "" {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService grpcService has filterStateObjectsToLog but its not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetInvertMatch() == true {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header invert match is not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetPresentMatch() == true {
			// TODO(nick): is this supported in Gateway API?
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header range match match is not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetPrefixMatch() != "" {
			//	HeaderMatchExact             HeaderMatchType = "Exact"
			//	HeaderMatchRegularExpression HeaderMatchType = "RegularExpression"
			//TODO(nick): can someone verify this is the equivalent?
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = filter.GetHeaderFilter().GetHeader().GetPrefixMatch() + ".*"
		}

		if filter.GetHeaderFilter().GetHeader().GetRangeMatch() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header range match match is not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetSafeRegexMatch() != nil {
			// Edge only supported Googles Regex (RE2) which might not be compatible with Gateway API regex
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = filter.GetHeaderFilter().GetHeader().GetSafeRegexMatch().GetRegex()
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
			Exclude:  filter.GetGrpcStatusFilter().GetExclude(),
		}
		for _, status := range filter.GetGrpcStatusFilter().GetStatuses() {
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
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService runtimeFilter is not supported")
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
// TODO(nick) - need to figure out which fields go to which policy. For example: httpConnectionManagerSettings.streamIdleTimeout: 3600s
func (g *GatewayAPIOutput) convertHTTPListenerOptions(options *gloov1.HttpListenerOptions, wrapper snapshot.Wrapper, proxyName string) {
	if options == nil {
		return
	}

	trafficPolicy := &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
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
			TargetSelectors: nil, // existing
			AI:              nil, // existing
			Transformation:  nil, // existing
			ExtProc:         nil, // existing
			ExtAuth:         nil, // existing
			RateLimit:       nil, // existing
			Cors:            nil, // existing
			Csrf:            nil, // existing
			HeaderModifiers: nil,
			AutoHostRewrite: nil,
			Buffer:          nil, // existing
			Timeouts:        nil, // existing
			Retry:           nil, // existing
			RBAC:            nil,
		},
		GlooRateLimit:      nil, // existing
		GlooExtAuth:        nil, // existing
		GlooTransformation: nil, // existing
		GlooJWT:            nil, // existing
		GlooRBAC:           nil, // existing
	}

	// go through each option in Gateway Options and convert to listener policy

	// inline extAuth settings
	if options.GetExtauth() != nil {
		// These are global extAuthSettings that are also on the Settings Object.
		// If this exists we need to generate a GatewayExtensionObject for this
		gatewayExtensions := g.generateGatewayExtensionForExtAuth(options.GetExtauth(), wrapper.GetName(), wrapper)
		g.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	// inline extProc settings
	if options.GetExtProc() != nil {
		gatewayExtensions := g.generateGatewayExtensionForExtProc(options.GetExtProc(), wrapper.GetName(), wrapper)
		g.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	// inline rate limit settings
	if options.GetRatelimitServer() != nil {
		gatewayExtensions := g.generateGatewayExtensionForRateLimit(options.GetRatelimitServer(), wrapper.GetName(), wrapper)
		g.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	if options.GetHttpLocalRatelimit() != nil {
		if options.GetHttpLocalRatelimit().GetEnableXRatelimitHeaders() != nil && options.GetHttpLocalRatelimit().GetEnableXRatelimitHeaders().GetValue() {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpLocalRateLimit enableXRateLimitHeaders is not supported")
		}
		if options.GetHttpLocalRatelimit().GetLocalRateLimitPerDownstreamConnection() != nil && options.GetHttpLocalRatelimit().GetLocalRateLimitPerDownstreamConnection().GetValue() {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpLocalRateLimit localRateLimitPerDownstreamConnection is not supported")
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
				rl.Local.TokenBucket.FillInterval = metav1.Duration{Duration: options.GetHttpLocalRatelimit().GetDefaultLimit().GetFillInterval().AsDuration()}
			}
			tps.TrafficPolicySpec.RateLimit = rl
		}
	}

	if options.GetWaf() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF is not supported")
		//waf := &gloogateway.Waf{
		//	Disabled:      ptr.To(options.GetWaf().GetDisabled()),
		//	CustomMessage: ptr.To(options.GetWaf().GetCustomInterventionMessage()),
		//	Rules:         []gloogateway.WafRule{},
		//}
		//for _, r := range options.GetWaf().GetRuleSets() {
		//	waf.Rules = append(waf.Rules, gloogateway.WafRule{
		//		RuleStr: ptr.To(r.GetRuleStr()),
		//	})
		//	if r.GetFiles() != nil && len(r.GetFiles()) > 0 {
		//		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF files is not supported")
		//	}
		//}
		//tps.Waf = waf
	}
	if options.GetDisableExtProc() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "disableExtProc is not supported")
	}
	if options.GetNetworkLocalRatelimit() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "networkLocalRateLimit is not supported")
	}
	if options.GetDlp() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported")
	}
	if options.GetCsrf() != nil {
		csrf := g.convertCSRF(options.GetCsrf())
		tps.TrafficPolicySpec.Csrf = csrf
	}
	if options.GetBuffer() != nil {
		tps.TrafficPolicySpec.Buffer = &kgateway.Buffer{
			MaxRequestSize: resource.NewQuantity(int64(options.GetBuffer().GetMaxRequestBytes().GetValue()), resource.BinarySI),
		}
	}
	if options.GetCaching() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "caching is not supported")
	}
	if options.GetConnectionLimit() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "connectionlimit is not supported")
	}
	if options.GetDynamicForwardProxy() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dynamicForwardProxy (DFP) is not supported")
	}
	if options.GetExtensions() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported")
	}
	if options.GetGrpcJsonTranscoder() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "grpcToJson is not supported")
	}
	if options.GetGrpcWeb() != nil {
		//TODO(nick) : GRPCWeb is enabled by default in edge. we need to verify the same.
		//o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "grpcWeb is not supported")
	}
	if options.GetGzip() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gzip is not supported")
	}
	if options.GetHeaderValidationSettings() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "header validation is not supported")
	}

	if options.GetProxyLatency() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "proxy latency is not supported")
	}
	if options.GetRouter() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "router (envoy filter maps) is not supported")
	}
	if options.GetSanitizeClusterHeader() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "sanitize cluster header is not supported")
	}
	if options.GetStatefulSession() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "statefulSession is not supported")
	}
	if options.GetTap() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "Tap filter is not supported")
	}
	if options.GetWasm() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WASM is not supported")
	}
	trafficPolicy.Spec = tps

	g.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(trafficPolicy, wrapper.FileOrigin()))

	if options.GetHealthCheck() != nil || options.GetHttpConnectionManagerSettings() != nil {
		// we need to create an HTTPListenerPolicy for these
		httpListenerPolicy := &kgateway.HTTPListenerPolicy{
			TypeMeta: metav1.TypeMeta{
				Kind:       "HTTPListenerPolicy",
				APIVersion: kgateway.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      wrapper.GetName(),
				Namespace: wrapper.GetNamespace(),
				Labels:    wrapper.GetLabels(),
			},
		}

		// TODO(nick) this is going to be an issue where multiple Gloo Gateways applies to the same Gateway (due to multiple ports)
		hlp := kgateway.HTTPListenerPolicySpec{
			TargetRefs: []kgateway.LocalPolicyTargetReference{
				{
					Group: gwv1.Group(gwv1.GroupVersion.Group),
					Kind:  "Gateway",
					Name:  gwv1.ObjectName(proxyName),
				},
			},
			TargetSelectors:            nil, // existing
			AccessLog:                  nil, // existing
			Tracing:                    nil, // existing
			UpgradeConfig:              nil, // existing
			UseRemoteAddress:           nil, // existing
			XffNumTrustedHops:          nil, // existing
			ServerHeaderTransformation: nil, // existing
			StreamIdleTimeout:          nil, // existing
			HealthCheck:                nil, // existing
		}
		// now in httplistenersettings
		if options.GetHealthCheck() != nil {
			hc := &kgateway.EnvoyHealthCheck{
				Path: options.GetHealthCheck().GetPath(),
			}
			hlp.HealthCheck = hc
		}
		// now in httplistenersettings
		if options.GetHttpConnectionManagerSettings() != nil {
			if options.GetHttpConnectionManagerSettings().GetTracing() != nil {
				t := g.convertHTTPListenerOptionsTracing(options.GetHttpConnectionManagerSettings().GetTracing(), wrapper)
				if t != nil {
					hlp.Tracing = t
				}
			}
			if options.GetHttpConnectionManagerSettings().GetSkipXffAppend() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.skipXffAppend is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetVia() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.via is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetXffNumTrustedHops() != nil {
				hlp.XffNumTrustedHops = ptr.To(options.GetHttpConnectionManagerSettings().GetXffNumTrustedHops().GetValue())
			}
			if options.GetHttpConnectionManagerSettings().GetUseRemoteAddress() != nil && options.GetHttpConnectionManagerSettings().GetUseRemoteAddress().GetValue() {
				hlp.UseRemoteAddress = ptr.To(options.GetHttpConnectionManagerSettings().GetUseRemoteAddress().GetValue())
			}
			if options.GetHttpConnectionManagerSettings().GetGenerateRequestId() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.generateRequestId is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetProxy_100Continue() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.proxy_100Continue is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetStreamIdleTimeout() != nil {
				hlp.StreamIdleTimeout = ptr.To(metav1.Duration{Duration: options.GetHttpConnectionManagerSettings().GetStreamIdleTimeout().AsDuration()})
			}
			if options.GetHttpConnectionManagerSettings().GetIdleTimeout() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.idleTimeout is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetMaxRequestHeadersKb() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.maxRequestHeadersKb is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetRequestTimeout() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.requestTimeout is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetRequestHeadersTimeout() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.requestHeadersTimeout is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetDrainTimeout() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.drainTimeout is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetDelayedCloseTimeout() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.delayCloseTimeout is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetServerName() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.serverName is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetStripAnyHostPort() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.stripAnyHostPort is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetAcceptHttp_10() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.acceptHttp_10 is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetDefaultHostForHttp_10() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.defaultHostForHttp_10 is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetAllowChunkedLength() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.allowChunkedLength is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetEnableTrailers() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.enableTrailers is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetProperCaseHeaderKeyFormat() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.properCaseHeaderKeyFormat is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetPreserveCaseHeaderKeyFormat() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.preserveCaseHeaderKeyFormat is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetSetCurrentClientCertDetails() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.setCurrentClientCertDetails is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetPreserveExternalRequestId() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.preserveExternalRequestId is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetUpgrades() != nil {
				var upgrades []string
				for _, upgrade := range options.GetHttpConnectionManagerSettings().GetUpgrades() {
					switch upgrade.GetUpgradeType().(type) {
					case *protocol_upgrade.ProtocolUpgradeConfig_Websocket:
						upgrades = append(upgrades, "websocket")
					case *protocol_upgrade.ProtocolUpgradeConfig_Connect:
						upgrades = append(upgrades, "CONNECT")
					}
				}
				hlp.UpgradeConfig = &kgateway.UpgradeConfig{EnabledUpgrades: upgrades}
			}
			if options.GetHttpConnectionManagerSettings().GetMaxConnectionDuration() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.maxConnectionDuration is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetMaxStreamDuration() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.maxStreamDuration is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetMaxHeadersCount() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.maxHeadersCount is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetHeadersWithUnderscoresAction() > 0 {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.headersWithUnderscoresAction is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetMaxRequestsPerConnection() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.maxRequestsPerConnection is not supported")
			}

			switch options.GetHttpConnectionManagerSettings().GetServerHeaderTransformation() {
			case hcm.HttpConnectionManagerSettings_OVERWRITE:
				hlp.ServerHeaderTransformation = ptr.To(kgateway.OverwriteServerHeaderTransformation)
			case hcm.HttpConnectionManagerSettings_APPEND_IF_ABSENT:
				hlp.ServerHeaderTransformation = ptr.To(kgateway.AppendIfAbsentServerHeaderTransformation)
			case hcm.HttpConnectionManagerSettings_PASS_THROUGH:
				hlp.ServerHeaderTransformation = ptr.To(kgateway.PassThroughServerHeaderTransformation)
			}

			if options.GetHttpConnectionManagerSettings().GetPathWithEscapedSlashesAction() > 0 {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.pathWithEscapedSlashesAction is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetCodecType() > 0 {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.codecType is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetMergeSlashes() != nil && options.GetHttpConnectionManagerSettings().GetMergeSlashes().GetValue() {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.mergeSlashes is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetNormalizePath() != nil && options.GetHttpConnectionManagerSettings().GetNormalizePath().GetValue() {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.normalizePath is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetUuidRequestIdConfig() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.uuidRequestIdConfig is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetHttp2ProtocolOptions() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.http2ProtocolOptions is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetInternalAddressConfig() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.internalAddressConfig is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetAppendXForwardedPort() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.appendXForwardedPort is not supported")
			}
			if options.GetHttpConnectionManagerSettings().GetEarlyHeaderManipulation() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManager.earlyHeaderManipulation is not supported")
			}
		}
		httpListenerPolicy.Spec = hlp
		g.gatewayAPICache.AddHTTPListenerPolicy(snapshot.NewHTTPListenerPolicyWrapper(httpListenerPolicy, wrapper.FileOrigin()))
	}
}

func (g *GatewayAPIOutput) convertHTTPListenerOptionsTracing(tracing *tracing.ListenerTracingSettings, wrapper snapshot.Wrapper) *kgateway.Tracing {
	kgatewayTracing := &kgateway.Tracing{
		Provider:          kgateway.TracingProvider{},   // existing
		ClientSampling:    nil,                          // existing
		RandomSampling:    nil,                          // existing
		OverallSampling:   nil,                          // existing
		Verbose:           nil,                          // existing
		MaxPathTagLength:  nil,                          // existing
		Attributes:        []kgateway.CustomAttribute{}, // existing
		SpawnUpstreamSpan: nil,                          // existing
	}
	if tracing.GetTracePercentages() != nil {
		if tracing.GetTracePercentages().GetClientSamplePercentage() != nil {
			kgatewayTracing.ClientSampling = ptr.To(uint32(tracing.GetTracePercentages().GetClientSamplePercentage().GetValue()))
		}
		if tracing.GetTracePercentages().GetOverallSamplePercentage() != nil {
			kgatewayTracing.OverallSampling = ptr.To(uint32(tracing.GetTracePercentages().GetOverallSamplePercentage().GetValue()))
		}
		if tracing.GetTracePercentages().GetRandomSamplePercentage() != nil {
			kgatewayTracing.RandomSampling = ptr.To(uint32(tracing.GetTracePercentages().GetRandomSamplePercentage().GetValue()))
		}
	}
	if tracing.GetDatadogConfig() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "Datadog kgatewayTracing is not supported")
	}
	if tracing.GetLiteralsForTags() != nil {
		for _, tag := range tracing.GetLiteralsForTags() {
			attribute := kgateway.CustomAttribute{
				Name: tag.GetTag().GetValue(),
				Literal: &kgateway.CustomAttributeLiteral{
					Value: tag.GetValue().GetValue(),
				},
			}
			kgatewayTracing.Attributes = append(kgatewayTracing.Attributes, attribute)
		}
	}
	if tracing.GetEnvironmentVariablesForTags() != nil {
		for _, tag := range tracing.GetEnvironmentVariablesForTags() {
			attribute := kgateway.CustomAttribute{
				Name: tag.GetTag().GetValue(),
				Environment: &kgateway.CustomAttributeEnvironment{
					Name:         tag.GetName().GetValue(),
					DefaultValue: nil,
				},
			}
			if tag.GetDefaultValue() != nil {
				attribute.Environment.DefaultValue = ptr.To(tag.GetDefaultValue().GetValue())
			}
			kgatewayTracing.Attributes = append(kgatewayTracing.Attributes, attribute)
		}
	}
	if tracing.GetMetadataForTags() != nil {
		for _, tag := range tracing.GetMetadataForTags() {
			attribute := kgateway.CustomAttribute{
				Name: tag.GetTag(),
				Metadata: &kgateway.CustomAttributeMetadata{
					Kind: kgateway.MetadataKind(tag.GetKind().String()),
					MetadataKey: kgateway.MetadataKey{
						Key: tag.GetValue().GetKey(),
					},
					DefaultValue: nil,
				},
			}
			if tag.GetDefaultValue() != "" {
				attribute.Environment.DefaultValue = ptr.To(tag.GetDefaultValue())
			}
			if tag.GetValue().GetNamespace() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "metadataForTags value.namespace is not supported")
			}
			if tag.GetValue().GetNestedFieldDelimiter() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "metadataForTags value.nestedFieldDelimiter is not supported")
			}

			kgatewayTracing.Attributes = append(kgatewayTracing.Attributes, attribute)
		}
	}
	if tracing.GetOpenCensusConfig() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "OpenCensus kgatewayTracing is not supported")
	}
	if tracing.GetOpenTelemetryConfig() != nil {
		otlp := &kgateway.OpenTelemetryTracingConfig{
			GrpcService:       kgateway.CommonGrpcService{}, // existing
			ServiceName:       nil,                          // existing
			ResourceDetectors: nil,                          // existing
			Sampler:           nil,                          // existing
		}
		if tracing.GetOpenTelemetryConfig().GetServiceName() != "" {
			otlp.ServiceName = ptr.To(tracing.GetOpenTelemetryConfig().GetServiceName())
		}
		if tracing.GetOpenTelemetryConfig().GetGrpcService() != nil {
			gs := tracing.GetOpenTelemetryConfig().GetGrpcService()
			grpcService := kgateway.CommonGrpcService{
				BackendRef:              nil, // existing
				Authority:               nil, // existing
				MaxReceiveMessageLength: nil, // existing
				SkipEnvoyHeaders:        nil, // existing
				Timeout:                 nil, // existing
				InitialMetadata:         nil, // existing
				RetryPolicy:             nil, // existing
			}
			if gs.GetAuthority() != "" {
				grpcService.Authority = ptr.To(gs.GetAuthority())
			}

			otlp.GrpcService = grpcService
		}
		kgatewayTracing.Provider.OpenTelemetry = otlp
	}
	if tracing.GetRequestHeadersForTags() != nil {
		for _, tag := range tracing.GetRequestHeadersForTags() {
			attribute := kgateway.CustomAttribute{
				Name: tag.GetValue(),
				RequestHeader: &kgateway.CustomAttributeHeader{
					Name:         tag.GetValue(),
					DefaultValue: nil,
				},
			}
			kgatewayTracing.Attributes = append(kgatewayTracing.Attributes, attribute)
		}
	}
	if tracing.GetSpawnUpstreamSpan() {
		kgatewayTracing.SpawnUpstreamSpan = ptr.To(tracing.GetSpawnUpstreamSpan())
	}
	if tracing.GetVerbose() != nil && tracing.GetVerbose().GetValue() {
		kgatewayTracing.Verbose = ptr.To(tracing.GetVerbose().GetValue())
	}
	if tracing.GetZipkinConfig() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "Zipkin kgatewayTracing is not supported")
	}
	return kgatewayTracing
}

func (g *GatewayAPIOutput) convertCSRF(policy *v4.CsrfPolicy) *kgateway.CSRFPolicy {
	csrf := &kgateway.CSRFPolicy{
		PercentageEnabled:  nil,
		PercentageShadowed: nil,
		AdditionalOrigins:  nil,
	}
	if policy.GetFilterEnabled() != nil {
		filterEnabled := policy.GetFilterEnabled()

		// Convert FractionalPercent to numerical percentage
		var percentage float64
		switch filterEnabled.GetDefaultValue().GetDenominator() {
		case v3.FractionalPercent_HUNDRED:
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator())
		case v3.FractionalPercent_TEN_THOUSAND:
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator()) / 100.0
		case v3.FractionalPercent_MILLION:
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator()) / 10000.0
		default:
			// Default to HUNDRED if denominator is not set
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator())
		}
		csrf.PercentageEnabled = ptr.To(uint32(percentage))
	}
	if policy.GetAdditionalOrigins() != nil {
		// Convert the additional origins from Gloo Edge format to kgateway format
		var additionalOrigins []kgateway.StringMatcher
		for _, origin := range policy.GetAdditionalOrigins() {
			switch typed := origin.GetMatchPattern().(type) {
			case *gloo_type_matcher.StringMatcher_Exact:
				additionalOrigins = append(additionalOrigins, kgateway.StringMatcher{Exact: ptr.To(typed.Exact)})
			case *gloo_type_matcher.StringMatcher_Prefix:
				additionalOrigins = append(additionalOrigins, kgateway.StringMatcher{Prefix: ptr.To(typed.Prefix)})
			case *gloo_type_matcher.StringMatcher_Suffix:
				additionalOrigins = append(additionalOrigins, kgateway.StringMatcher{Suffix: ptr.To(typed.Suffix)})
			case *gloo_type_matcher.StringMatcher_SafeRegex:
				additionalOrigins = append(additionalOrigins, kgateway.StringMatcher{SafeRegex: ptr.To(typed.SafeRegex.GetRegex())})
			}
		}
		csrf.AdditionalOrigins = additionalOrigins
	}
	if policy.GetShadowEnabled() != nil {
		shadowEnabled := policy.GetShadowEnabled()

		// Convert FractionalPercent to numerical percentage
		var percentage float64
		switch shadowEnabled.GetDefaultValue().GetDenominator() {
		case v3.FractionalPercent_HUNDRED:
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator())
		case v3.FractionalPercent_TEN_THOUSAND:
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator()) / 100.0
		case v3.FractionalPercent_MILLION:
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator()) / 10000.0
		default:
			// Default to HUNDRED if denominator is not set
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator())
		}

		csrf.PercentageShadowed = ptr.To(uint32(percentage))
	}
	return csrf
}

func (g *GatewayAPIOutput) generateGatewayExtensionForExtProc(extProc *extproc.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {

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
			Type:    kgateway.GatewayExtensionTypeExtProc,
			ExtProc: &kgateway.ExtProcProvider{},
		},
	}

	//TODO(nick): Implement ExtProc - https://github.com/kgateway-dev/kgateway/issues/11424
	if extProc.GetStatPrefix() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc statPrefix is not supported")
	}
	if extProc.GetFailureModeAllow() != nil && extProc.GetFailureModeAllow().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc failureModeAllow is not supported")
	}
	if extProc.GetAllowModeOverride() != nil && extProc.GetAllowModeOverride().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc allowModeOverride is not supported")
	}
	if extProc.GetAsyncMode() != nil && extProc.GetAsyncMode().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc asyncMode is not supported")
	}
	if extProc.GetDisableClearRouteCache() != nil && extProc.GetDisableClearRouteCache().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc disableClearRouteCache is not supported")
	}
	if extProc.GetFilterMetadata() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc filterMetadata is not supported")
	}
	if extProc.GetFilterStage() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc filterStage is not supported")
	}
	if extProc.GetForwardRules() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc forwardRules is not supported")
	}
	if extProc.GetMaxMessageTimeout() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc maxMessageTimeout is not supported")
	}
	if extProc.GetMessageTimeout() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc messageTimeout is not supported")
	}
	if extProc.GetMetadataContextNamespaces() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc metadataContextNamespaces is not supported")
	}
	if extProc.GetMutationRules() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc mutationRules is not supported")
	}
	if extProc.GetRequestAttributes() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc requestAttributes is not supported")
	}
	if extProc.GetResponseAttributes() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc responseAttributes is not supported")
	}
	if extProc.GetTypedMetadataContextNamespaces() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc typedMetadataContextNamespaces is not supported")
	}
	if extProc.GetGrpcService() != nil {

		backend := g.GetGatewayAPICache().GetBackend(types.NamespacedName{Namespace: extProc.GetGrpcService().GetExtProcServerRef().GetNamespace(), Name: extProc.GetGrpcService().GetExtProcServerRef().GetName()})
		if backend == nil {
			g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "extProc grpcService backend not found")
		}
		grpcService := &kgateway.ExtGrpcService{
			BackendRef: &gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(backend.GetName()),
					Namespace: ptr.To(gwv1.Namespace(backend.GetNamespace())),
					// using the default port here
					Port: ptr.To(gwv1.PortNumber(4444)),
				},
			},
		}
		gatewayExtension.Spec.ExtProc.GrpcService = grpcService
	}
	return gatewayExtension
}
func (g *GatewayAPIOutput) generateGatewayExtensionForRateLimit(rateLimitSettings *ratelimit.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {

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
			Type:      kgateway.GatewayExtensionTypeRateLimit,
			RateLimit: &kgateway.RateLimitProvider{},
		},
	}

	//TODO(nick): Implement RateLimitSettings - https://github.com/kgateway-dev/kgateway/issues/11424
	if rateLimitSettings.GetRateLimitBeforeAuth() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "RateLimitSettings rateLimitBeforeAuth is not supported")
	}
	if rateLimitSettings.GetEnableXRatelimitHeaders() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "RateLimitSettings enableXRatelimitHeaders is not supported")
	}
	if rateLimitSettings.GetDenyOnFail() == false {
		gatewayExtension.Spec.RateLimit.FailOpen = !rateLimitSettings.GetDenyOnFail()
	}
	if rateLimitSettings.GetRequestTimeout() != nil {
		gatewayExtension.Spec.RateLimit.Timeout = metav1.Duration{Duration: rateLimitSettings.GetRequestTimeout().AsDuration()}
	}
	if rateLimitSettings.GetRatelimitServerRef() != nil {
		backend := g.GetGatewayAPICache().GetBackend(types.NamespacedName{Namespace: rateLimitSettings.GetRatelimitServerRef().GetNamespace(), Name: rateLimitSettings.GetRatelimitServerRef().GetName()})
		var grpcService *kgateway.ExtGrpcService
		if backend != nil {
			// backend exists so we use that reference
			grpcService = &kgateway.ExtGrpcService{
				BackendRef: &gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      gwv1.ObjectName(backend.GetName()),
						Namespace: ptr.To(gwv1.Namespace(backend.GetNamespace())),
						// using the default port here
						Port: ptr.To(gwv1.PortNumber(18081)),
					},
				},
			}
		} else {
			//TODO(nick): just assuming its a kube service
			grpcService = &kgateway.ExtGrpcService{
				BackendRef: &gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      gwv1.ObjectName(rateLimitSettings.GetRatelimitServerRef().GetName()),
						Namespace: ptr.To(gwv1.Namespace(rateLimitSettings.GetRatelimitServerRef().GetNamespace())),
						// using the default port here
						Port: ptr.To(gwv1.PortNumber(18081)),
					},
				},
			}
		}

		if rateLimitSettings.GetGrpcService() != nil {
			grpcService.Authority = ptr.To(rateLimitSettings.GetGrpcService().GetAuthority())
		}
		gatewayExtension.Spec.RateLimit.GrpcService = grpcService
	}

	return gatewayExtension
}
func (g *GatewayAPIOutput) generateGatewayExtensionForExtAuth(extauth *v1.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {
	if extauth == nil {
		return nil
	}
	var grpcService *kgateway.ExtGrpcService
	if extauth.GetExtauthzServerRef() != nil {
		backend := g.GetGatewayAPICache().GetBackend(types.NamespacedName{Namespace: extauth.GetExtauthzServerRef().GetNamespace(), Name: extauth.GetExtauthzServerRef().GetName()})
		if backend != nil {
			// this has a backend definition
			grpcService = &kgateway.ExtGrpcService{
				BackendRef: &gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      gwv1.ObjectName(backend.GetName()),
						Namespace: ptr.To(gwv1.Namespace(backend.GetNamespace())),
						// using the default port here
						Port: ptr.To(gwv1.PortNumber(8083)),
					},
				},
			}
		} else {
			// TODO(nick): if the backend wasnt found im just assuming it was a kube service which is a pretty safe assumption
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
		}

		if extauth.GetGrpcService() != nil {
			grpcService.Authority = ptr.To(extauth.GetGrpcService().GetAuthority())
		}
	}
	//TODO(nick): Implement ExtAuthSettings - https://github.com/kgateway-dev/kgateway/issues/11424
	if extauth.GetClearRouteCache() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings clearRouteCache is not supported")
	}
	if extauth.GetFailureModeAllow() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings clearRouteCache is not supported")
	}
	if extauth.GetHttpService() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings httpService is not supported")
	}
	if extauth.GetRequestBody() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings requestBody is not supported")
	}
	if extauth.GetRequestTimeout() != nil {
		// Moved to GatewayExtension.spec.extAuth.grpcService.requestTimeout
		g.AddErrorFromWrapper(ERROR_TYPE_UPDATE_OBJECT, wrapper, "extAuth settings requestTimeout needs to be set on the GatewayExtension.spec.extAuth.grpcService.requestTimeout")
	}
	if extauth.GetStatPrefix() != "" {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings statPrefix is not supported")
	}
	if extauth.GetStatusOnError() != 0 {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings statusOnError is not supported")
	}
	if extauth.GetTransportApiVersion() != 0 {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings transportApiVersion is not supported")
	}
	if extauth.GetUserIdHeader() != "" {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings userIdHeader is not supported")
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
