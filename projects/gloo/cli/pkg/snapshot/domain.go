package snapshot

import (
	api "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/types"
)

type Instance struct {
	parseErrors          []error
	yamlObjects          []*YAMLWrapper
	settings             map[types.NamespacedName]*SettingsWrapper
	routeTables          map[types.NamespacedName]*RouteTableWrapper
	routeOptions         map[types.NamespacedName]*RouteOptionWrapper
	listenerOptions      map[types.NamespacedName]*ListenerOptionWrapper
	httpListenerOptions  map[types.NamespacedName]*HTTPListenerOptionWrapper
	upstreams            map[types.NamespacedName]*UpstreamWrapper
	virtualServices      map[types.NamespacedName]*VirtualServiceWrapper
	glooGateways         map[types.NamespacedName]*GlooGatewayWrapper
	authConfigs          map[types.NamespacedName]*AuthConfigWrapper
	delegatesRouteTables []*api.RouteTable

	//Gateway API Configs
	httpRoutes         map[types.NamespacedName]*HTTPRouteWrapper
	directResponses    map[types.NamespacedName]*DirectResponseWrapper
	listenerSets       map[types.NamespacedName]*ListenerSetWrapper
	virtualHostOptions map[types.NamespacedName]*VirtualHostOptionWrapper
	gateways           map[types.NamespacedName]*GatewayWrapper
	gatewayParameters  map[types.NamespacedName]*GatewayParametersWrapper
}

func (i *Instance) YAMLObjects() []*YAMLWrapper {
	if i.yamlObjects == nil {
		i.yamlObjects = []*YAMLWrapper{}
	}
	return i.yamlObjects
}
func (i *Instance) Settings() map[types.NamespacedName]*SettingsWrapper {
	if i.settings == nil {
		i.settings = make(map[types.NamespacedName]*SettingsWrapper)
	}
	return i.settings
}
func (i *Instance) RouteTables() map[types.NamespacedName]*RouteTableWrapper {
	if i.routeTables == nil {
		i.routeTables = make(map[types.NamespacedName]*RouteTableWrapper)
	}
	return i.routeTables
}
func (i *Instance) RouteOptions() map[types.NamespacedName]*RouteOptionWrapper {
	if i.routeOptions == nil {
		i.routeOptions = make(map[types.NamespacedName]*RouteOptionWrapper)
	}
	return i.routeOptions
}
func (i *Instance) ListenerOptions() map[types.NamespacedName]*ListenerOptionWrapper {
	if i.listenerOptions == nil {
		i.listenerOptions = make(map[types.NamespacedName]*ListenerOptionWrapper)
	}
	return i.listenerOptions
}
func (i *Instance) HTTPListenerOptions() map[types.NamespacedName]*HTTPListenerOptionWrapper {
	if i.httpListenerOptions == nil {
		i.httpListenerOptions = make(map[types.NamespacedName]*HTTPListenerOptionWrapper)
	}
	return i.httpListenerOptions
}
func (i *Instance) Upstreams() map[types.NamespacedName]*UpstreamWrapper {
	if i.upstreams == nil {
		i.upstreams = make(map[types.NamespacedName]*UpstreamWrapper)
	}
	return i.upstreams
}
func (i *Instance) VirtualServices() map[types.NamespacedName]*VirtualServiceWrapper {
	if i.virtualServices == nil {
		i.virtualServices = make(map[types.NamespacedName]*VirtualServiceWrapper)
	}
	return i.virtualServices
}
func (i *Instance) GlooGateways() map[types.NamespacedName]*GlooGatewayWrapper {
	if i.glooGateways == nil {
		i.glooGateways = make(map[types.NamespacedName]*GlooGatewayWrapper)
	}
	return i.glooGateways
}

func (i *Instance) AuthConfigs() map[types.NamespacedName]*AuthConfigWrapper {
	if i.authConfigs == nil {
		i.authConfigs = make(map[types.NamespacedName]*AuthConfigWrapper)
	}
	return i.authConfigs
}
func (i *Instance) HTTPRoutes() map[types.NamespacedName]*HTTPRouteWrapper {
	if i.httpRoutes == nil {
		i.httpRoutes = make(map[types.NamespacedName]*HTTPRouteWrapper)
	}
	return i.httpRoutes
}
func (i *Instance) DirectResponses() map[types.NamespacedName]*DirectResponseWrapper {
	if i.directResponses == nil {
		i.directResponses = make(map[types.NamespacedName]*DirectResponseWrapper)
	}
	return i.directResponses
}
func (i *Instance) ListenerSets() map[types.NamespacedName]*ListenerSetWrapper {
	if i.listenerSets == nil {
		i.listenerSets = make(map[types.NamespacedName]*ListenerSetWrapper)
	}
	return i.listenerSets
}
func (i *Instance) VirtualHostOptions() map[types.NamespacedName]*VirtualHostOptionWrapper {
	if i.virtualHostOptions == nil {
		i.virtualHostOptions = make(map[types.NamespacedName]*VirtualHostOptionWrapper)
	}
	return i.virtualHostOptions
}
func (i *Instance) Gateways() map[types.NamespacedName]*GatewayWrapper {
	if i.gateways == nil {
		i.gateways = make(map[types.NamespacedName]*GatewayWrapper)
	}
	return i.gateways
}
func (i *Instance) GatewayParameters() map[types.NamespacedName]*GatewayParametersWrapper {
	if i.gatewayParameters == nil {
		i.gatewayParameters = make(map[types.NamespacedName]*GatewayParametersWrapper)
	}
	return i.gatewayParameters
}

func (i *Instance) ParseErrors() []error {
	if i.parseErrors == nil {
		i.parseErrors = make([]error, 0)
	}
	return i.parseErrors
}

func (i *Instance) AddGatewayParameters(w *GatewayParametersWrapper) {
	if i.gatewayParameters == nil {
		i.gatewayParameters = make(map[types.NamespacedName]*GatewayParametersWrapper)
	}
	i.gatewayParameters[w.Index()] = w
}
func (i *Instance) AddGateway(route *GatewayWrapper) {
	if i.gateways == nil {
		i.gateways = make(map[types.NamespacedName]*GatewayWrapper)
	}
	i.gateways[route.Index()] = route
}
func (i *Instance) AddHTTPRoute(route *HTTPRouteWrapper) {
	if i.httpRoutes == nil {
		i.httpRoutes = make(map[types.NamespacedName]*HTTPRouteWrapper)
	}
	i.httpRoutes[route.Index()] = route
}

func (i *Instance) AddListenerSet(l *ListenerSetWrapper) {
	if i.listenerSets == nil {
		i.listenerSets = make(map[types.NamespacedName]*ListenerSetWrapper)
	}
	i.listenerSets[l.Index()] = l
}

func (i *Instance) AddDirectResponse(d *DirectResponseWrapper) {
	if i.directResponses == nil {
		i.directResponses = make(map[types.NamespacedName]*DirectResponseWrapper)
	}
	i.directResponses[d.Index()] = d
}

func (i *Instance) AddSettings(d *SettingsWrapper) {
	if i.settings == nil {
		i.settings = make(map[types.NamespacedName]*SettingsWrapper)
	}
	i.settings[d.Index()] = d
}

func (i *Instance) AddAuthConfig(d *AuthConfigWrapper) {

	if i.authConfigs == nil {
		i.authConfigs = make(map[types.NamespacedName]*AuthConfigWrapper)
	}
	i.authConfigs[d.Index()] = d
}

func (i *Instance) AddUpstream(u *UpstreamWrapper) {
	if i.upstreams == nil {
		i.upstreams = make(map[types.NamespacedName]*UpstreamWrapper)
	}
	i.upstreams[u.Index()] = u
}

func (i *Instance) AddRouteTable(r *RouteTableWrapper) {
	if i.routeTables == nil {
		i.routeTables = make(map[types.NamespacedName]*RouteTableWrapper)
	}
	i.routeTables[r.Index()] = r
}

func (i *Instance) AddVirtualService(v *VirtualServiceWrapper) {
	if i.virtualServices == nil {
		i.virtualServices = make(map[types.NamespacedName]*VirtualServiceWrapper)
	}
	i.virtualServices[v.Index()] = v
}

func (i *Instance) AddRouteOption(r *RouteOptionWrapper) {
	if i.routeOptions == nil {
		i.routeOptions = make(map[types.NamespacedName]*RouteOptionWrapper)
	}
	i.routeOptions[r.Index()] = r
}

func (i *Instance) AddVirtualHostOption(v *VirtualHostOptionWrapper) {
	if i.virtualHostOptions == nil {
		i.virtualHostOptions = make(map[types.NamespacedName]*VirtualHostOptionWrapper)
	}
	i.virtualHostOptions[v.Index()] = v
}

func (i *Instance) AddGlooGateway(g *GlooGatewayWrapper) {
	if i.glooGateways == nil {
		i.glooGateways = make(map[types.NamespacedName]*GlooGatewayWrapper)
	}
	i.glooGateways[g.Index()] = g
}

func (i *Instance) AddHTTPListenerOption(h *HTTPListenerOptionWrapper) {
	if i.httpListenerOptions == nil {
		i.httpListenerOptions = make(map[types.NamespacedName]*HTTPListenerOptionWrapper)
	}
	i.httpListenerOptions[h.Index()] = h
}

func (i *Instance) AddListenerOption(l *ListenerOptionWrapper) {
	if i.listenerOptions == nil {
		i.listenerOptions = make(map[types.NamespacedName]*ListenerOptionWrapper)
	}
	i.listenerOptions[l.Index()] = l
}

func (i *Instance) AddYamlObject(y *YAMLWrapper) {
	if i.yamlObjects == nil {
		i.yamlObjects = []*YAMLWrapper{}
	}
	i.yamlObjects = append(i.yamlObjects, y)
}

// GlooGatewayVirtualServices finds all VirtualServiceWrapper that match the gateways selectors
func (i *Instance) GlooGatewayVirtualServices(gateway *GlooGatewayWrapper) ([]*VirtualServiceWrapper, error) {
	var response []*VirtualServiceWrapper

	for _, virtualService := range i.VirtualServices() {
		// this logic requires the VS to have a namespace and metadata on the inner object which we dont have. So we do the namespace lookup ourselves
		vs := &virtualService.Spec
		vs.SetMetadata(&core.Metadata{
			Name:        virtualService.Name,
			Namespace:   virtualService.Namespace,
			Labels:      virtualService.Labels,
			Annotations: virtualService.Annotations,
		})
		matches, err := translator.HttpGatewayContainsVirtualService(gateway.Gateway.Spec.GetHttpGateway(), vs, gateway.Gateway.Spec.GetSsl())
		if err != nil {
			return nil, err
		}

		if matches {
			response = append(response, virtualService)
		}
	}
	return response, nil
}

// DelegatedRouteTables finds all RouteTables that match the virtual service delegate selector
func (i *Instance) DelegatedRouteTables(namespace string, delegateAction *api.DelegateAction) ([]*RouteTableWrapper, error) {
	var response []*RouteTableWrapper

	//get all routeTables
	rtlist := []*api.RouteTable{}
	if len(i.delegatesRouteTables) != len(i.routeTables) {
		//TODO this could be expensive and may need to move this to a special property
		for _, rtw := range i.RouteTables() {
			rtSpec := &rtw.Spec
			// set all the metadata on the route table object.
			rtSpec.Metadata = &core.Metadata{
				Name:        rtw.Name,
				Namespace:   rtw.Namespace,
				Labels:      rtw.Labels,
				Annotations: rtw.Annotations,
			}

			rtlist = append(rtlist, &rtw.Spec)
		}
		i.delegatesRouteTables = rtlist
	}

	rts := translator.NewRouteTableSelector(i.delegatesRouteTables)

	selectedRouteTables, err := rts.SelectRouteTables(delegateAction, namespace)
	if err != nil {
		return nil, err
	}
	for _, rt := range selectedRouteTables {
		routeTable := i.routeTables[types.NamespacedName{Name: rt.GetMetadata().GetName(), Namespace: rt.GetMetadata().GetNamespace()}]
		response = append(response, routeTable)
	}

	return response, nil
}

func (i *Instance) VirtualServicesByNameNamespace(name, namespace string) []*VirtualServiceWrapper {
	var response []*VirtualServiceWrapper
	if name == "" && namespace == "" {
		return response
	}
	if name == "" {
		// namespace only search
		for vsName, virtualService := range i.virtualServices {
			if vsName.Namespace == namespace {
				response = append(response, virtualService)
			}
		}
	} else if namespace == "" {
		for vsName, virtualService := range i.virtualServices {
			if vsName.Name == name {
				response = append(response, virtualService)
			}
		}
	} else {
		vs := i.virtualServices[types.NamespacedName{Namespace: namespace, Name: name}]
		if vs != nil {
			response = append(response, vs)
		}
	}
	return response
}

func (i *Instance) GetUpstream(namespacedName types.NamespacedName) *UpstreamWrapper {
	return i.Upstreams()[namespacedName]
}
