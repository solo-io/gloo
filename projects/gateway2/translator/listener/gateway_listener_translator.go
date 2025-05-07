package listener

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"google.golang.org/protobuf/types/known/wrapperspb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"github.com/solo-io/gloo/projects/gateway2/ports"
	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/backendref"
	route "github.com/solo-io/gloo/projects/gateway2/translator/httproute"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/routeutils"
	"github.com/solo-io/gloo/projects/gateway2/translator/sslutils"
	"github.com/solo-io/gloo/projects/gateway2/translator/types"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

// TranslateListeners translates the set of gloo listeners required to produce a full output proxy (either form one Gateway or multiple merged Gateways)
func TranslateListeners(
	ctx context.Context,
	queries query.GatewayQueries,
	pluginRegistry registry.PluginRegistry,
	consolidatedGateway *types.ConsolidatedGateway,
	routesForGw *query.RoutesForGwResult,
	reporter reports.Reporter,
) []*v1.Listener {
	validatedConsolidatedListeners := validateConsolidatedGateway(consolidatedGateway, reporter)
	mergedListeners := mergeConsolidatedListeners(queries, validatedConsolidatedListeners, consolidatedGateway, routesForGw, reporter)
	translatedListeners := mergedListeners.translateListeners(ctx, pluginRegistry, queries, reporter)
	return translatedListeners
}

func mergeConsolidatedListeners(
	queries query.GatewayQueries,
	consolidatedListeners []types.ConsolidatedListener,
	consolidatedGateway *types.ConsolidatedGateway,
	routesForGw *query.RoutesForGwResult,
	reporter reports.Reporter,
) *MergedListeners {
	ml := &MergedListeners{
		parentGw:         *consolidatedGateway.Gateway,
		GatewayNamespace: consolidatedGateway.Gateway.Namespace,
		Queries:          queries,
	}

	for _, cl := range consolidatedListeners {
		listener := cl.Listener
		result := routesForGw.GetListenerResult(cl.GetParent(), string(listener.Name))
		listenerReporter := cl.GetParentReporter(reporter).Listener(listener)

		if result == nil || result.Error != nil {
			// TODO report
			// TODO, if Error is not nil, this is a user-config error on selectors
			// continue
		}
		var routes []*query.RouteInfo
		if result != nil {
			routes = result.Routes
		}
		ml.AppendListener(*listener, routes, listenerReporter, cl.ListenerSet)
	}
	return ml
}

type MergedListeners struct {
	GatewayNamespace string
	parentGw         gwv1.Gateway
	Listeners        []*MergedListener
	Queries          query.GatewayQueries
}

func (ml *MergedListeners) AppendListener(
	listener gwv1.Listener,
	routes []*query.RouteInfo,
	reporter reports.ListenerReporter,
	listenerSet *gwxv1a1.XListenerSet,
) error {
	switch listener.Protocol {
	case gwv1.HTTPProtocolType:
		ml.appendHttpListener(listener, routes, reporter, listenerSet)
	case gwv1.HTTPSProtocolType:
		ml.appendHttpsListener(listener, routes, reporter, listenerSet)
	// TODO default handling
	case gwv1.TCPProtocolType:
		ml.AppendTcpListener(listener, routes, reporter, listenerSet)
	case gwv1.TLSProtocolType:
		ml.AppendTlsListener(listener, routes, reporter, listenerSet)
	default:
		return eris.Errorf("unsupported protocol: %v", listener.Protocol)
	}

	return nil
}

func (ml *MergedListeners) appendHttpListener(
	listener gwv1.Listener,
	routesWithHosts []*query.RouteInfo,
	reporter reports.ListenerReporter,
	listenerSet *gwxv1a1.XListenerSet,
) {
	parent := httpFilterChainParent{
		gatewayListenerName: string(listener.Name),
		routesWithHosts:     routesWithHosts,
	}

	fc := &httpFilterChain{
		parents: []httpFilterChainParent{parent},
	}
	listenerName := generateListenerName(listener)
	finalPort := gwv1.PortNumber(ports.TranslatePort(uint16(listener.Port)))

	for _, lis := range ml.Listeners {
		if lis.port == finalPort {
			if lis.httpFilterChain != nil {
				lis.httpFilterChain.parents = append(lis.httpFilterChain.parents, parent)
			} else {
				lis.httpFilterChain = fc
			}
			return
		}
	}

	// create a new filter chain for the listener
	ml.Listeners = append(ml.Listeners, &MergedListener{
		name:             listenerName,
		gatewayNamespace: ml.GatewayNamespace,
		port:             finalPort,
		httpFilterChain:  fc,
		listenerReporter: reporter,
		listener:         listener,
		listenerSet:      listenerSet,
	})
}

func (ml *MergedListeners) appendHttpsListener(
	listener gwv1.Listener,
	routesWithHosts []*query.RouteInfo,
	reporter reports.ListenerReporter,
	listenerSet *gwxv1a1.XListenerSet,
) {
	// create a new filter chain for the listener
	// protocol:            listener.Protocol,
	mfc := httpsFilterChain{
		gatewayListenerName: string(listener.Name),
		sniDomain:           listener.Hostname,
		tls:                 listener.TLS,
		routesWithHosts:     routesWithHosts,
		queries:             ml.Queries,
	}

	// Perform the port transformation away from privileged ports only once to use
	// during both lookup and when appending the listener.
	finalPort := gwv1.PortNumber(ports.TranslatePort(uint16(listener.Port)))

	listenerName := generateListenerName(listener)
	for _, lis := range ml.Listeners {
		if lis.port == finalPort {
			lis.httpsFilterChains = append(lis.httpsFilterChains, mfc)
			return
		}
	}
	ml.Listeners = append(ml.Listeners, &MergedListener{
		name:              listenerName,
		gatewayNamespace:  ml.GatewayNamespace,
		port:              finalPort,
		httpsFilterChains: []httpsFilterChain{mfc},
		listenerReporter:  reporter,
		listener:          listener,
		listenerSet:       listenerSet,
	})
}

func (ml *MergedListeners) AppendTcpListener(
	listener gwv1.Listener,
	routeInfos []*query.RouteInfo,
	reporter reports.ListenerReporter,
	listenerSet *gwxv1a1.XListenerSet,
) {
	var validRouteInfos []*query.RouteInfo

	for _, routeInfo := range routeInfos {
		tRoute, ok := routeInfo.Object.(*gwv1a2.TCPRoute)
		if !ok {
			continue
		}

		if len(tRoute.Spec.ParentRefs) == 0 {
			contextutils.LoggerFrom(context.Background()).Warnf(
				"No parent references found for TCPRoute %s", tRoute.Name,
			)
			continue
		}

		validRouteInfos = append(validRouteInfos, routeInfo)
	}

	// If no valid routes are found, do not create a listener
	if len(validRouteInfos) == 0 {
		contextutils.LoggerFrom(context.Background()).Errorf(
			"No valid routes found for listener %s", listener.Name,
		)
		return
	}

	parent := tcpFilterChainParent{
		gatewayListenerName: string(listener.Name),
		routesWithHosts:     validRouteInfos,
	}

	fc := tcpFilterChain{
		parents: []tcpFilterChainParent{parent},
	}
	listenerName := generateListenerName(listener)
	finalPort := gwv1.PortNumber(ports.TranslatePort(uint16(listener.Port)))

	for _, lis := range ml.Listeners {
		if lis.port == finalPort {
			lis.TcpFilterChains = append(lis.TcpFilterChains, fc)
			return
		}
	}

	// create a new filter chain for the listener
	ml.Listeners = append(ml.Listeners, &MergedListener{
		name:             listenerName,
		gatewayNamespace: ml.GatewayNamespace,
		port:             finalPort,
		TcpFilterChains:  []tcpFilterChain{fc},
		listenerReporter: reporter,
		listener:         listener,
		listenerSet:      listenerSet,
	})
}

// buildTcpHost builds a Gloo TcpHost from the provided parameters.
func buildTcpHost(
	routeInfo *query.RouteInfo,
	parentRefReporters []reports.ParentRefReporter,
	tcpRouteName string,
	defaultPort gwv1.PortNumber,
	backendRefs []gwv1.BackendRef,
) *v1.TcpHost {
	// If there are no backendRefs, return nil to skip this TcpHost
	if len(backendRefs) == 0 {
		return nil
	}

	// Use the TCPRoute name for the tcpHost name
	tcpHost := &v1.TcpHost{Name: tcpRouteName}

	var weightedDestinations []*v1.WeightedDestination

	for _, ref := range backendRefs {
		// Try to get the backend object
		obj, err := routeInfo.GetBackendForRef(ref.BackendObjectReference)
		if err != nil {
			// Process error and set the appropriate status conditions
			for _, parentRefReporter := range parentRefReporters {
				query.ProcessBackendError(err, parentRefReporter)
			}
			continue
		}

		// Process the backend object
		var destination *v1.Destination
		if backendref.RefIsService(ref.BackendObjectReference) {
			var port uint32
			if ref.Port != nil {
				port = uint32(*ref.Port)
			} else {
				port = uint32(defaultPort)
			}

			destination = &v1.Destination{
				DestinationType: &v1.Destination_Kube{
					Kube: &v1.KubernetesServiceDestination{
						Ref: &core.ResourceRef{
							Name:      obj.GetName(),
							Namespace: obj.GetNamespace(),
						},
						Port: port,
					},
				},
			}
		} else {
			// Unsupported kind
			err := query.ErrUnknownBackendKind
			for _, parentRefReporter := range parentRefReporters {
				query.ProcessBackendError(err, parentRefReporter)
			}
			continue
		}

		weightedDestinations = append(weightedDestinations, &v1.WeightedDestination{
			Destination: destination,
			Weight:      getWeight(ref),
		})
	}

	// Set the TcpHost destination type
	if len(weightedDestinations) == 0 {
		// No valid destinations, return nil
		return nil
	} else if len(weightedDestinations) == 1 {
		tcpHost.Destination = &v1.TcpHost_TcpAction{
			Destination: &v1.TcpHost_TcpAction_Single{
				Single: weightedDestinations[0].GetDestination(),
			},
		}
	} else {
		// Multiple destinations, set up a Multi destination
		tcpHost.Destination = &v1.TcpHost_TcpAction{
			Destination: &v1.TcpHost_TcpAction_Multi{
				Multi: &v1.MultiDestination{
					Destinations: weightedDestinations,
				},
			},
		}
	}

	return tcpHost
}

func getWeight(backendRef gwv1.BackendRef) *wrapperspb.UInt32Value {
	if backendRef.Weight != nil {
		return &wrapperspb.UInt32Value{Value: uint32(*backendRef.Weight)}
	}
	// Default weight is 1
	return &wrapperspb.UInt32Value{Value: 1}
}

func (ml *MergedListeners) AppendTlsListener(
	listener gwv1.Listener,
	routeInfos []*query.RouteInfo,
	reporter reports.ListenerReporter,
	listenerSet *gwxv1a1.XListenerSet,
) {
	var validRouteInfos []*query.RouteInfo

	for _, routeInfo := range routeInfos {
		tRoute, ok := routeInfo.Object.(*gwv1a2.TLSRoute)
		if !ok {
			continue
		}

		if len(tRoute.Spec.ParentRefs) == 0 {
			contextutils.LoggerFrom(context.Background()).Warnf(
				"No parent references found for TLSRoute %s", tRoute.Name,
			)
			continue
		}

		validRouteInfos = append(validRouteInfos, routeInfo)
	}

	// If no valid routes are found, do not create a listener
	if len(validRouteInfos) == 0 {
		contextutils.LoggerFrom(context.Background()).Errorf(
			"No valid routes found for listener %s", listener.Name,
		)
		return
	}

	parent := tcpFilterChainParent{
		gatewayListenerName: string(listener.Name),
		routesWithHosts:     validRouteInfos,
	}

	fc := tcpFilterChain{
		parents:   []tcpFilterChainParent{parent},
		tls:       listener.TLS,
		sniDomain: listener.Hostname,
	}
	listenerName := generateListenerName(listener)
	finalPort := gwv1.PortNumber(ports.TranslatePort(uint16(listener.Port)))

	for _, lis := range ml.Listeners {
		if lis.port == finalPort {
			lis.TcpFilterChains = append(lis.TcpFilterChains, fc)
			return
		}
	}

	// create a new filter chain for the listener
	ml.Listeners = append(ml.Listeners, &MergedListener{
		name:             listenerName,
		gatewayNamespace: ml.GatewayNamespace,
		port:             finalPort,
		TcpFilterChains:  []tcpFilterChain{fc},
		listenerReporter: reporter,
		listener:         listener,
		listenerSet:      listenerSet,
	})
}

func (ml *MergedListeners) translateListeners(
	ctx context.Context,
	pluginRegistry registry.PluginRegistry,
	queries query.GatewayQueries,
	reporter reports.Reporter,
) []*v1.Listener {
	var listeners []*v1.Listener
	for _, mergedListener := range ml.Listeners {
		listener := mergedListener.TranslateListener(ctx, pluginRegistry, queries, reporter)

		// run listener plugins
		for _, listenerPlugin := range pluginRegistry.GetListenerPlugins() {
			err := listenerPlugin.ApplyListenerPlugin(ctx, &plugins.ListenerContext{
				Gateway:     &ml.parentGw,
				ListenerSet: mergedListener.listenerSet,
				GwListener:  &mergedListener.listener,
			}, listener)
			if err != nil {
				contextutils.LoggerFrom(ctx).Errorf("error in ListenerPlugin: %v", err)
			}
		}

		listeners = append(listeners, listener)
	}
	return listeners
}

type MergedListener struct {
	name              string
	gatewayNamespace  string
	port              gwv1.PortNumber
	httpFilterChain   *httpFilterChain
	httpsFilterChains []httpsFilterChain
	TcpFilterChains   []tcpFilterChain
	listenerReporter  reports.ListenerReporter
	listener          gwv1.Listener
	listenerSet       *gwxv1a1.XListenerSet

	// TODO(policy via http listener options)
}

func (ml *MergedListener) TranslateListener(
	ctx context.Context,
	pluginRegistry registry.PluginRegistry,
	queries query.GatewayQueries,
	reporter reports.Reporter,
) *v1.Listener {
	var (
		httpFilterChains    []*v1.AggregateListener_HttpFilterChain
		mergedVhosts        = map[string]*v1.VirtualHost{}
		matchedTcpListeners []*v1.MatchedTcpListener
	)

	// Translate HTTP filter chains
	if ml.httpFilterChain != nil {
		httpFilterChain, vhostsForFilterchain := ml.httpFilterChain.translateHttpFilterChain(
			ctx,
			ml.name,
			ml.listener,
			pluginRegistry,
			reporter,
		)
		httpFilterChains = append(httpFilterChains, httpFilterChain)
		for vhostRef, vhost := range vhostsForFilterchain {
			if _, ok := mergedVhosts[vhostRef]; ok {
				// Handle potential error if duplicate vhosts are found
				contextutils.LoggerFrom(ctx).Errorf(
					"Duplicate virtual host found: %s", vhostRef,
				)
				continue
			}
			mergedVhosts[vhostRef] = vhost
		}
	}

	// Translate HTTPS filter chains
	for _, mfc := range ml.httpsFilterChains {
		httpsFilterChain, vhostsForFilterchain := mfc.translateHttpsFilterChain(
			ctx,
			pluginRegistry,
			mfc.gatewayListenerName,
			ml.gatewayNamespace,
			ml.listener,
			queries,
			reporter,
			ml.listenerReporter,
		)
		if httpsFilterChain == nil {
			// Log and skip invalid HTTPS filter chains
			contextutils.LoggerFrom(ctx).Errorf("Failed to translate HTTPS filter chain for listener: %s", ml.name)
			continue
		}

		httpFilterChains = append(httpFilterChains, httpsFilterChain)

		for vhostRef, vhost := range vhostsForFilterchain {
			if _, ok := mergedVhosts[vhostRef]; ok {
				contextutils.LoggerFrom(ctx).Errorf("Duplicate virtual host found: %s", vhostRef)
				continue
			}
			mergedVhosts[vhostRef] = vhost
		}
	}

	// Translate TCP listeners (if any exist)
	for _, tfc := range ml.TcpFilterChains {
		if tcpListener := tfc.translateTcpFilterChain(
			ctx,
			ml.listener,
			reporter,
			ml.listenerReporter,
			queries,
			ml.gatewayNamespace,
		); tcpListener != nil {
			matchedTcpListeners = append(matchedTcpListeners, &v1.MatchedTcpListener{
				TcpListener: tcpListener,
			})
		}
	}

	// Create and return the listener with all filter chains and TCP listeners
	return &v1.Listener{
		Name:        ml.name,
		BindAddress: "::",
		BindPort:    uint32(ml.port),
		ListenerType: &v1.Listener_AggregateListener{
			AggregateListener: &v1.AggregateListener{
				HttpResources: &v1.AggregateListener_HttpResources{
					VirtualHosts: mergedVhosts,
					HttpOptions:  nil, // HttpListenerOptions will be added by policy plugins
				},
				HttpFilterChains: httpFilterChains,
				TcpListeners:     matchedTcpListeners,
			},
		},
		// Used for tracing to create service_name
		OpaqueMetadata: &v1.Listener_MetadataStatic{
			MetadataStatic: &v1.SourceMetadata{
				Sources: []*v1.SourceMetadata_SourceRef{
					{
						ResourceRef: &core.ResourceRef{
							Name:      ml.name,
							Namespace: ml.gatewayNamespace,
						},
						ResourceKind: wellknown.GatewayGroup + "/" + wellknown.GatewayKind,
					},
				},
			},
		},
		Options:      nil, // Listener options will be added by policy plugins
		RouteOptions: nil,
	}
}

// tcpFilterChain each one represents a Gateway listener that has been merged into a single Gloo Listener
// (with distinct filter chains). In the case where no Gateway listener merging takes place, every listener
// will use a Gloo AggregatedListener with one TCP filter chain.
type tcpFilterChain struct {
	parents   []tcpFilterChainParent
	tls       *gwv1.GatewayTLSConfig
	sniDomain *gwv1.Hostname
}

type tcpFilterChainParent struct {
	gatewayListenerName string
	routesWithHosts     []*query.RouteInfo
}

func (tc *tcpFilterChain) translateTcpFilterChain(
	ctx context.Context,
	listener gwv1.Listener,
	reporter reports.Reporter,
	listenerReporter reports.ListenerReporter,
	queries query.GatewayQueries,
	gatewayNamespace string,
) *v1.TcpListener {
	var tcpHosts []*v1.TcpHost
	for _, parent := range tc.parents {
		for _, r := range parent.routesWithHosts {
			// TODO(puertomontt): deduplicate logic for TLSRoute and TCPRoute
			switch r.Object.(type) {
			case *gwv1a2.TLSRoute:
				tRoute := r.Object.(*gwv1a2.TLSRoute)
				// Collect ParentRefReporters for the TLSRoute
				parentRefReporters := make([]reports.ParentRefReporter, 0, len(tRoute.Spec.ParentRefs))
				for _, parentRef := range tRoute.Spec.ParentRefs {
					parentRefReporter := reporter.Route(tRoute).ParentRef(&parentRef)
					parentRefReporter.SetCondition(reports.RouteCondition{
						Type:   gwv1.RouteConditionAccepted,
						Status: metav1.ConditionTrue,
						Reason: gwv1.RouteReasonAccepted,
					})
					parentRefReporters = append(parentRefReporters, parentRefReporter)
				}

				for i, rule := range tRoute.Spec.Rules {
					// Ensure unique names by appending the rule index to the TLSRoute name
					tcpHostName := fmt.Sprintf("%s-rule-%d", tRoute.Name, i)
					tcpHost := buildTcpHost(r, parentRefReporters, tcpHostName, listener.Port, rule.BackendRefs)
					if tcpHost != nil {
						sslConfig, err := translateSslConfig(
							ctx,
							gatewayNamespace,
							tc.sniDomain,
							tc.tls,
							queries,
						)
						if err != nil {
							reason := gwv1.ListenerReasonRefNotPermitted
							if !errors.Is(err, query.ErrMissingReferenceGrant) {
								reason = gwv1.ListenerReasonInvalidCertificateRef
							}
							listenerReporter.SetCondition(reports.ListenerCondition{
								Type:   gwv1.ListenerConditionResolvedRefs,
								Status: metav1.ConditionFalse,
								Reason: reason,
							})
							// listener with no ssl is invalid. We return nil so set programmed to false
							listenerReporter.SetCondition(reports.ListenerCondition{
								Type:   gwv1.ListenerConditionProgrammed,
								Status: metav1.ConditionFalse,
								Reason: gwv1.ListenerReasonInvalid,
							})
							return nil
						}
						tcpHost.SslConfig = sslConfig
						tcpHosts = append(tcpHosts, tcpHost)
					}
				}
			case *gwv1a2.TCPRoute:
				tRoute := r.Object.(*gwv1a2.TCPRoute)
				// Collect ParentRefReporters for the TCPRoute
				parentRefReporters := make([]reports.ParentRefReporter, 0, len(tRoute.Spec.ParentRefs))
				for _, parentRef := range tRoute.Spec.ParentRefs {
					parentRefReporter := reporter.Route(tRoute).ParentRef(&parentRef)
					parentRefReporter.SetCondition(reports.RouteCondition{
						Type:   gwv1.RouteConditionAccepted,
						Status: metav1.ConditionTrue,
						Reason: gwv1.RouteReasonAccepted,
					})
					parentRefReporters = append(parentRefReporters, parentRefReporter)
				}

				for i, rule := range tRoute.Spec.Rules {
					// Ensure unique names by appending the rule index to the TCPRoute name
					tcpHostName := fmt.Sprintf("%s-rule-%d", tRoute.Name, i)
					tcpHost := buildTcpHost(r, parentRefReporters, tcpHostName, listener.Port, rule.BackendRefs)
					if tcpHost != nil {
						tcpHosts = append(tcpHosts, tcpHost)
					}
				}
			}
		}
	}

	// Avoid creating a TcpListener if there are no TcpHosts
	if len(tcpHosts) == 0 {
		return nil
	}

	return &v1.TcpListener{
		TcpHosts: tcpHosts,
	}
}

// httpFilterChain each one represents a GW Listener that has been merged into a single Gloo Listener (with distinct filter chains).
// In the case where no GW Listener merging takes place, every listener will use a Gloo AggregatedListeener with 1 HTTP filter chain.
type httpFilterChain struct {
	parents []httpFilterChainParent
}

type httpFilterChainParent struct {
	gatewayListenerName string
	routesWithHosts     []*query.RouteInfo
}

func (httpFilterChain *httpFilterChain) translateHttpFilterChain(
	ctx context.Context,
	parentName string,
	listener gwv1.Listener,
	pluginRegistry registry.PluginRegistry,
	reporter reports.Reporter,
) (*v1.AggregateListener_HttpFilterChain, map[string]*v1.VirtualHost) {
	routesByHost := map[string]routeutils.SortableRoutes{}
	for _, parent := range httpFilterChain.parents {
		buildRoutesPerHost(
			ctx,
			routesByHost,
			parent.routesWithHosts,
			listener,
			pluginRegistry,
			reporter,
		)
	}

	var (
		virtualHostRefs []string
		virtualHosts    = map[string]*v1.VirtualHost{}
	)
	for host, vhostRoutes := range routesByHost {
		sort.Stable(vhostRoutes)
		vhostName := makeVhostName(ctx, parentName, host)
		virtualHosts[vhostName] = &v1.VirtualHost{
			Name:    vhostName,
			Domains: []string{host},
			Routes:  vhostRoutes.ToRoutes(),
			Options: nil,
		}

		virtualHostRefs = append(virtualHostRefs, vhostName)
	}
	sort.Strings(virtualHostRefs)

	return &v1.AggregateListener_HttpFilterChain{
		Matcher:         &v1.Matcher{}, // http filter chain matcher is not used
		VirtualHostRefs: virtualHostRefs,
	}, virtualHosts
}

type httpsFilterChain struct {
	gatewayListenerName string
	sniDomain           *gwv1.Hostname
	tls                 *gwv1.GatewayTLSConfig
	routesWithHosts     []*query.RouteInfo
	queries             query.GatewayQueries
}

func (httpsFilterChain *httpsFilterChain) translateHttpsFilterChain(
	ctx context.Context,
	pluginRegistry registry.PluginRegistry,
	parentName string,
	gatewayNamespace string,
	listener gwv1.Listener,
	queries query.GatewayQueries,
	reporter reports.Reporter,
	listenerReporter reports.ListenerReporter,
) (*v1.AggregateListener_HttpFilterChain, map[string]*v1.VirtualHost) {
	// process routes first, so any route related errors are reported on the httproute.
	routesByHost := map[string]routeutils.SortableRoutes{}
	buildRoutesPerHost(
		ctx,
		routesByHost,
		httpsFilterChain.routesWithHosts,
		listener,
		pluginRegistry,
		reporter,
	)

	var (
		virtualHostRefs []string
		virtualHosts    = map[string]*v1.VirtualHost{}
	)
	for host, vhostRoutes := range routesByHost {
		sort.Stable(vhostRoutes)
		vhostName := makeVhostName(ctx, parentName, host)
		virtualHosts[vhostName] = &v1.VirtualHost{
			Name:    vhostName,
			Domains: []string{host},
			Routes:  vhostRoutes.ToRoutes(),
			Options: nil,
		}

		virtualHostRefs = append(virtualHostRefs, vhostName)
	}
	sort.Strings(virtualHostRefs)

	sslConfig, err := translateSslConfig(
		ctx,
		gatewayNamespace,
		httpsFilterChain.sniDomain,
		httpsFilterChain.tls,
		queries,
	)
	if err != nil {
		reason := gwv1.ListenerReasonRefNotPermitted
		if !errors.Is(err, query.ErrMissingReferenceGrant) {
			reason = gwv1.ListenerReasonInvalidCertificateRef
		}
		listenerReporter.SetCondition(reports.ListenerCondition{
			Type:   gwv1.ListenerConditionResolvedRefs,
			Status: metav1.ConditionFalse,
			Reason: reason,
		})
		// listener with no ssl is invalid. We return nil so set programmed to false
		listenerReporter.SetCondition(reports.ListenerCondition{
			Type:   gwv1.ListenerConditionProgrammed,
			Status: metav1.ConditionFalse,
			Reason: gwv1.ListenerReasonInvalid,
		})
		return nil, nil
	}
	matcher := &v1.Matcher{SslConfig: sslConfig, SourcePrefixRanges: nil, PassthroughCipherSuites: nil}

	return &v1.AggregateListener_HttpFilterChain{
		Matcher:         matcher,
		VirtualHostRefs: virtualHostRefs,
	}, virtualHosts
}

func buildRoutesPerHost(
	ctx context.Context,
	routesByHost map[string]routeutils.SortableRoutes,
	routes []*query.RouteInfo,
	gwListener gwv1.Listener,
	pluginRegistry registry.PluginRegistry,
	reporter reports.Reporter,
) {
	for _, routeWithHosts := range routes {
		parentRefReporter := reporter.Route(routeWithHosts.Object).ParentRef(&routeWithHosts.ParentRef)
		routes := route.TranslateGatewayHTTPRouteRules(
			ctx,
			pluginRegistry,
			gwListener,
			routeWithHosts,
			parentRefReporter,
			reporter,
		)

		if len(routes) == 0 {
			// TODO report
			continue
		}

		hostnames := routeWithHosts.Hostnames()
		if len(hostnames) == 0 {
			hostnames = []string{"*"}
		}

		for _, host := range hostnames {
			routesByHost[host] = append(routesByHost[host], routeutils.ToSortable(routeWithHosts.Object, routes)...)
		}
	}
}

func translateSslConfig(
	ctx context.Context,
	parentNamespace string,
	sniDomain *gwv1.Hostname,
	tls *gwv1.GatewayTLSConfig,
	queries query.GatewayQueries,
) (*ssl.SslConfig, error) {
	if tls == nil {
		return nil, nil
	}

	// TODO support passthrough mode
	if tls.Mode == nil {
		return nil, nil
	}

	var secretRef *core.ResourceRef
	if *tls.Mode == gwv1.TLSModeTerminate {
		for _, certRef := range tls.CertificateRefs {
			// validate via query
			secret, err := queries.GetSecretForRef(ctx, query.FromGkNs{
				Gk: metav1.GroupKind{
					Group: gwv1.GroupName,
					Kind:  wellknown.GatewayKind,
				},
				Ns: parentNamespace,
			}, certRef)
			if err != nil {
				return nil, err
			}
			// The resulting sslconfig will still have to go through a real translation where we run through this again.
			// This means that while its nice to still fail early here we dont need to scrub the actual contents of the secret.
			if _, err := sslutils.ValidateTlsSecret(secret.(*corev1.Secret)); err != nil {
				return nil, err
			}

			// TODO verify secret ref / grant using query
			secretNamespace := parentNamespace
			if certRef.Namespace != nil {
				secretNamespace = string(*certRef.Namespace)
			}
			secretRef = &core.ResourceRef{
				Name:      string(certRef.Name),
				Namespace: secretNamespace,
			}
			break // TODO support multiple certs
		}
		if secretRef == nil {
			return nil, nil
		}
	}

	var sniDomains []string
	if sniDomain != nil {
		sniDomains = []string{string(*sniDomain)}
	}
	cfg := &ssl.SslConfig{
		SslSecrets:                    nil,
		SniDomains:                    sniDomains,
		VerifySubjectAltName:          nil,
		Parameters:                    nil,
		AlpnProtocols:                 nil,
		OneWayTls:                     nil,
		DisableTlsSessionResumption:   nil,
		TransportSocketConnectTimeout: nil,
		OcspStaplePolicy:              0,
	}
	if secretRef != nil {
		cfg.SslSecrets = &ssl.SslConfig_SecretRef{SecretRef: secretRef}
	}

	// Apply known SSL Extension options
	sslutils.ApplySslExtensionOptions(ctx, tls, cfg)

	return cfg, nil
}

// makeVhostName computes the name of a virtual host based on the parent name and domain.
func makeVhostName(
	ctx context.Context,
	parentName string,
	domain string,
) string {
	return utils.SanitizeForEnvoy(ctx, parentName+"~"+domain, "vHost")
}

func generateListenerName(listener gwv1.Listener) string {
	// Add a ~ to make sure the name won't collide with user provided names in other listeners
	return fmt.Sprintf("listener~%d", listener.Port)
}
