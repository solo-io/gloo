package listener

import (
	"slices"

	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/types"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

const NormalizedHTTPSTLSType = "HTTPS/TLS"
const DefaultHostname = "*"
const AttachedListenerSetsConditionType = "AttachedListenerSets"

type portProtocol struct {
	hostnames map[gwv1.Hostname]int
	protocol  map[gwv1.ProtocolType]bool
	// needed for getting reporter? doesn't seem great
	listeners []types.ConsolidatedListener
}

type protocol = string
type groupName = string
type routeKind = string

func getSupportedProtocolsRoutes() map[protocol]map[groupName][]routeKind {
	// we currently only support HTTPRoute on HTTP and HTTPS protocols
	supportedProtocolToKinds := map[protocol]map[groupName][]routeKind{
		string(gwv1.HTTPProtocolType): {
			gwv1.GroupName: []string{
				wellknown.HTTPRouteKind,
			},
		},
		string(gwv1.HTTPSProtocolType): {
			gwv1.GroupName: []string{
				wellknown.HTTPRouteKind,
			},
		},
		string(gwv1.TCPProtocolType): {
			gwv1.GroupName: []string{
				wellknown.TCPRouteKind,
			},
		},
		string(gwv1.TLSProtocolType): {
			gwv1.GroupName: []string{
				wellknown.TLSRouteKind,
			},
		},
	}
	return supportedProtocolToKinds
}

func buildDefaultRouteKindsForProtocol(supportedRouteKindsForProtocol map[groupName][]routeKind) []gwv1.RouteGroupKind {
	rgks := []gwv1.RouteGroupKind{}
	for group, kinds := range supportedRouteKindsForProtocol {
		for _, kind := range kinds {
			rgks = append(rgks, gwv1.RouteGroupKind{
				Group: (*gwv1.Group)(&group),
				Kind:  gwv1.Kind(kind),
			})
		}
	}
	return rgks
}

func validateSupportedRoutes(consolidatedListeners []types.ConsolidatedListener, reporter reports.Reporter) []types.ConsolidatedListener {
	supportedProtocolToKinds := getSupportedProtocolsRoutes()
	var validListeners []types.ConsolidatedListener

	for _, cl := range consolidatedListeners {
		parentReporter := cl.GetParentReporter(reporter)

		listener := cl.Listener
		supportedRouteKindsForProtocol, ok := supportedProtocolToKinds[string(listener.Protocol)]
		if !ok {
			// todo: log?
			parentReporter.Listener(listener).SetCondition(reports.ListenerCondition{
				Type:   gwv1.ListenerConditionAccepted,
				Status: metav1.ConditionFalse,
				Reason: gwv1.ListenerReasonUnsupportedProtocol, //TODO: add message
			})
			continue
		}

		if listener.AllowedRoutes == nil || len(listener.AllowedRoutes.Kinds) == 0 {
			// default to whatever route kinds we support on this protocol
			// TODO(Law): confirm this matches spec
			rgks := buildDefaultRouteKindsForProtocol(supportedRouteKindsForProtocol)
			parentReporter.Listener(listener).SetSupportedKinds(rgks)
			validListeners = append(validListeners, cl)
			continue
		}

		foundSupportedRouteKinds := []gwv1.RouteGroupKind{}
		foundInvalidRouteKinds := []gwv1.RouteGroupKind{}
		for _, rgk := range listener.AllowedRoutes.Kinds {
			if rgk.Group == nil {
				// default to Gateway API group if not set
				rgk.Group = getGroupName()
			}
			supportedRouteKinds, ok := supportedRouteKindsForProtocol[string(*rgk.Group)]
			if !ok || !slices.Contains(supportedRouteKinds, string(rgk.Kind)) {
				foundInvalidRouteKinds = append(foundInvalidRouteKinds, rgk)
				continue
			}
			foundSupportedRouteKinds = append(foundSupportedRouteKinds, rgk)
		}

		parentReporter.Listener(listener).SetSupportedKinds(foundSupportedRouteKinds)
		if len(foundInvalidRouteKinds) > 0 {
			parentReporter.Listener(listener).SetCondition(reports.ListenerCondition{
				Type:   gwv1.ListenerConditionResolvedRefs,
				Status: metav1.ConditionFalse,
				Reason: gwv1.ListenerReasonInvalidRouteKinds,
			})
		} else {
			validListeners = append(validListeners, cl)
		}
	}

	return validListeners
}

func validateConsolidatedGateway(consolidatedGateway *types.ConsolidatedGateway, reporter reports.Reporter) []types.ConsolidatedListener {
	rejectDeniedListenerSets(consolidatedGateway, reporter)
	validatedListeners := validateListeners(consolidatedGateway, reporter)
	return validatedListeners
}

func rejectDeniedListenerSets(consolidatedGateway *types.ConsolidatedGateway, reporter reports.Reporter) {
	for _, ls := range consolidatedGateway.DeniedListenerSets {
		rejectListenerSet(ls, reporter.ListenerSet(ls))
	}
}

func rejectListenerSet(ls *gwxv1a1.XListenerSet, reporter reports.GatewayReporter) {
	reporter.SetCondition(reports.GatewayCondition{
		Type:   gwv1.GatewayConditionAccepted,
		Status: metav1.ConditionFalse,
		Reason: gwv1.GatewayConditionReason(gwxv1a1.ListenerSetReasonNotAllowed),
	})
	reporter.SetCondition(reports.GatewayCondition{
		Type:   gwv1.GatewayConditionProgrammed,
		Status: metav1.ConditionFalse,
		Reason: gwv1.GatewayConditionReason(gwxv1a1.ListenerSetReasonNotAllowed),
	})
}

func validateListeners(consolidatedGateway *types.ConsolidatedGateway, reporter reports.Reporter) []types.ConsolidatedListener {
	if len(consolidatedGateway.GetConsolidatedListeners()) == 0 {
		// gwReporter.Err("gateway must contain at least 1 listener")
	}

	validListeners := validateSupportedRoutes(consolidatedGateway.GetConsolidatedListeners(), reporter)

	portListeners := map[gwv1.PortNumber]*portProtocol{}
	for _, cl := range validListeners {
		listener := cl.Listener
		protocol := listener.Protocol
		if protocol == gwv1.HTTPSProtocolType || protocol == gwv1.TLSProtocolType {
			protocol = NormalizedHTTPSTLSType
		}

		if existingListener, ok := portListeners[listener.Port]; ok {
			existingListener.protocol[protocol] = true
			existingListener.listeners = append(existingListener.listeners, cl)
			//TODO(Law): handle validation that hostname empty for udp/tcp
			if listener.Hostname != nil {
				existingListener.hostnames[*listener.Hostname]++
			} else {
				existingListener.hostnames[DefaultHostname]++
			}
		} else {
			var hostname gwv1.Hostname
			if listener.Hostname == nil {
				hostname = DefaultHostname
			} else {
				hostname = *listener.Hostname
			}
			pp := portProtocol{
				hostnames: map[gwv1.Hostname]int{
					hostname: 1,
				},
				protocol: map[gwv1.ProtocolType]bool{
					protocol: true,
				},
				listeners: []types.ConsolidatedListener{cl},
			}
			portListeners[listener.Port] = &pp
		}
	}

	// reset valid listeners
	validListeners = []types.ConsolidatedListener{}
	for _, pp := range portListeners {

		protocolConflict := false
		if len(pp.protocol) > 1 {
			protocolConflict = true
		}

		for _, cl := range pp.listeners {
			listener := cl.Listener
			parentReporter := cl.GetParentReporter(reporter)
			if protocolConflict {
				parentReporter.Listener(listener).SetCondition(reports.ListenerCondition{
					Type:    gwv1.ListenerConditionConflicted,
					Status:  metav1.ConditionTrue,
					Reason:  gwv1.ListenerReasonProtocolConflict,
					Message: "Found conflicting protocols on listeners, a single port can only contain listeners with compatible protocols",
				})

				// continue as protocolConflict will take precedence over hostname conflicts
				continue
			}

			var hostname gwv1.Hostname
			if listener.Hostname == nil {
				hostname = DefaultHostname
			} else {
				hostname = *listener.Hostname
			}
			if count := pp.hostnames[hostname]; count > 1 {
				parentReporter.Listener(listener).SetCondition(reports.ListenerCondition{
					Type:    gwv1.ListenerConditionConflicted,
					Status:  metav1.ConditionTrue,
					Reason:  gwv1.ListenerReasonHostnameConflict,
					Message: "Found conflicting hostnames on listeners, all listeners on a single port must have unique hostnames",
				})
			} else {
				// TODO should check this is exactly 1?
				validListeners = append(validListeners, cl)
			}
		}
	}

	// Add the final conditions on the Gateway
	if consolidatedGateway.AllowedListenerSets == nil {
		reporter.Gateway(consolidatedGateway.Gateway).SetCondition(reports.GatewayCondition{
			Type:   AttachedListenerSetsConditionType,
			Status: metav1.ConditionUnknown,
			Reason: gwv1.GatewayReasonNoResources,
		})
	}

	if len(validListeners) == 0 {
		reporter.Gateway(consolidatedGateway.Gateway).SetCondition(reports.GatewayCondition{
			Type:   gwv1.GatewayConditionAccepted,
			Status: metav1.ConditionFalse,
			Reason: gwv1.GatewayReasonListenersNotValid,
		})
		reporter.Gateway(consolidatedGateway.Gateway).SetCondition(reports.GatewayCondition{
			Type:   gwv1.GatewayConditionProgrammed,
			Status: metav1.ConditionFalse,
			Reason: gwv1.GatewayReasonInvalid,
		})
		return validListeners
	}

	if validListeners[len(validListeners)-1].ListenerSet != nil {
		reporter.Gateway(consolidatedGateway.Gateway).SetCondition(reports.GatewayCondition{
			Type:   AttachedListenerSetsConditionType,
			Status: metav1.ConditionTrue,
			Reason: gwv1.GatewayReasonAccepted,
		})
	} else {
		reporter.Gateway(consolidatedGateway.Gateway).SetCondition(reports.GatewayCondition{
			Type:   AttachedListenerSetsConditionType,
			Status: metav1.ConditionFalse,
			Reason: gwv1.GatewayReasonNoResources,
		})
	}

	return validListeners
}

func getGroupName() *gwv1.Group {
	g := gwv1.Group(gwv1.GroupName)
	return &g
}
