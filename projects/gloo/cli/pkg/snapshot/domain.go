package snapshot

import (
	"fmt"
	"strings"

	api "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type Instance struct {
	parseErrors          []error
	yamlObjects          []*YAMLWrapper
	settings             map[string]*SettingsWrapper
	routeTables          map[string]*RouteTableWrapper
	routeOptions         map[string]*RouteOptionWrapper
	listenerOptions      map[string]*ListenerOptionWrapper
	httpListenerOptions  map[string]*HTTPListenerOptionWrapper
	upstreams            map[string]*UpstreamWrapper
	virtualServices      map[string]*VirtualServiceWrapper
	glooGateways         map[string]*GlooGatewayWrapper
	authConfigs          map[string]*AuthConfigWrapper
	delegatesRouteTables []*api.RouteTable

	//Gateway API Configs
	httpRoutes         map[string]*HTTPRouteWrapper
	directResponses    map[string]*DirectResponseWrapper
	listenerSets       map[string]*ListenerSetWrapper
	virtualHostOptions map[string]*VirtualHostOptionWrapper
	gateways           map[string]*GatewayWrapper
	gatewayParameters  map[string]*GatewayParametersWrapper
}

func (i *Instance) YAMLObjects() []*YAMLWrapper {
	if i.yamlObjects == nil {
		i.yamlObjects = []*YAMLWrapper{}
	}
	return i.yamlObjects
}
func (i *Instance) Settings() map[string]*SettingsWrapper {
	if i.settings == nil {
		i.settings = make(map[string]*SettingsWrapper)
	}
	return i.settings
}
func (i *Instance) RouteTables() map[string]*RouteTableWrapper {
	if i.routeTables == nil {
		i.routeTables = make(map[string]*RouteTableWrapper)
	}
	return i.routeTables
}
func (i *Instance) RouteOptions() map[string]*RouteOptionWrapper {
	if i.routeOptions == nil {
		i.routeOptions = make(map[string]*RouteOptionWrapper)
	}
	return i.routeOptions
}
func (i *Instance) ListenerOptions() map[string]*ListenerOptionWrapper {
	if i.listenerOptions == nil {
		i.listenerOptions = make(map[string]*ListenerOptionWrapper)
	}
	return i.listenerOptions
}
func (i *Instance) HTTPListenerOptions() map[string]*HTTPListenerOptionWrapper {
	if i.httpListenerOptions == nil {
		i.httpListenerOptions = make(map[string]*HTTPListenerOptionWrapper)
	}
	return i.httpListenerOptions
}
func (i *Instance) Upstreams() map[string]*UpstreamWrapper {
	if i.upstreams == nil {
		i.upstreams = make(map[string]*UpstreamWrapper)
	}
	return i.upstreams
}
func (i *Instance) VirtualServices() map[string]*VirtualServiceWrapper {
	if i.virtualServices == nil {
		i.virtualServices = make(map[string]*VirtualServiceWrapper)
	}
	return i.virtualServices
}
func (i *Instance) GlooGateways() map[string]*GlooGatewayWrapper {
	if i.glooGateways == nil {
		i.glooGateways = make(map[string]*GlooGatewayWrapper)
	}
	return i.glooGateways
}

func (i *Instance) AuthConfigs() map[string]*AuthConfigWrapper {
	if i.authConfigs == nil {
		i.authConfigs = make(map[string]*AuthConfigWrapper)
	}
	return i.authConfigs
}
func (i *Instance) HTTPRoutes() map[string]*HTTPRouteWrapper {
	if i.httpRoutes == nil {
		i.httpRoutes = make(map[string]*HTTPRouteWrapper)
	}
	return i.httpRoutes
}
func (i *Instance) DirectResponses() map[string]*DirectResponseWrapper {
	if i.directResponses == nil {
		i.directResponses = make(map[string]*DirectResponseWrapper)
	}
	return i.directResponses
}
func (i *Instance) ListenerSets() map[string]*ListenerSetWrapper {
	if i.listenerSets == nil {
		i.listenerSets = make(map[string]*ListenerSetWrapper)
	}
	return i.listenerSets
}
func (i *Instance) VirtualHostOptions() map[string]*VirtualHostOptionWrapper {
	if i.virtualHostOptions == nil {
		i.virtualHostOptions = make(map[string]*VirtualHostOptionWrapper)
	}
	return i.virtualHostOptions
}
func (i *Instance) Gateways() map[string]*GatewayWrapper {
	if i.gateways == nil {
		i.gateways = make(map[string]*GatewayWrapper)
	}
	return i.gateways
}
func (i *Instance) GatewayParameters() map[string]*GatewayParametersWrapper {
	if i.gatewayParameters == nil {
		i.gatewayParameters = make(map[string]*GatewayParametersWrapper)
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
		i.gatewayParameters = make(map[string]*GatewayParametersWrapper)
	}
	i.gatewayParameters[w.Index()] = w
}
func (i *Instance) AddGateway(route *GatewayWrapper) {
	if i.gateways == nil {
		i.gateways = make(map[string]*GatewayWrapper)
	}
	i.gateways[route.Index()] = route
}
func (i *Instance) AddHTTPRoute(route *HTTPRouteWrapper) {
	if i.httpRoutes == nil {
		i.httpRoutes = make(map[string]*HTTPRouteWrapper)
	}
	i.httpRoutes[route.Index()] = route
}

func (i *Instance) AddListenerSet(l *ListenerSetWrapper) {
	if i.listenerSets == nil {
		i.listenerSets = make(map[string]*ListenerSetWrapper)
	}
	i.listenerSets[l.Index()] = l
}

func (i *Instance) AddDirectResponse(d *DirectResponseWrapper) {
	if i.directResponses == nil {
		i.directResponses = make(map[string]*DirectResponseWrapper)
	}
	i.directResponses[d.Index()] = d
}

func (i *Instance) AddSettings(d *SettingsWrapper) {
	if i.settings == nil {
		i.settings = make(map[string]*SettingsWrapper)
	}
	i.settings[d.Index()] = d
}

func (i *Instance) AddAuthConfig(d *AuthConfigWrapper) {

	if i.authConfigs == nil {
		i.authConfigs = make(map[string]*AuthConfigWrapper)
	}
	i.authConfigs[d.Index()] = d
}

func (i *Instance) AddUpstream(u *UpstreamWrapper) {
	if i.upstreams == nil {
		i.upstreams = make(map[string]*UpstreamWrapper)
	}
	i.upstreams[u.Index()] = u
}

func (i *Instance) AddRouteTable(r *RouteTableWrapper) {
	if i.routeTables == nil {
		i.routeTables = make(map[string]*RouteTableWrapper)
	}
	i.routeTables[r.Index()] = r
}

func (i *Instance) AddVirtualService(v *VirtualServiceWrapper) {
	if i.virtualServices == nil {
		i.virtualServices = make(map[string]*VirtualServiceWrapper)
	}
	i.virtualServices[v.Index()] = v
}

func (i *Instance) AddRouteOption(r *RouteOptionWrapper) {
	if i.routeOptions == nil {
		i.routeOptions = make(map[string]*RouteOptionWrapper)
	}
	i.routeOptions[r.Index()] = r
}

func (i *Instance) AddVirtualHostOption(v *VirtualHostOptionWrapper) {
	if i.virtualHostOptions == nil {
		i.virtualHostOptions = make(map[string]*VirtualHostOptionWrapper)
	}
	i.virtualHostOptions[v.Index()] = v
}

func (i *Instance) AddGlooGateway(g *GlooGatewayWrapper) {
	if i.glooGateways == nil {
		i.glooGateways = make(map[string]*GlooGatewayWrapper)
	}
	i.glooGateways[g.Index()] = g
}

func (i *Instance) AddHTTPListenerOption(h *HTTPListenerOptionWrapper) {
	if i.httpListenerOptions == nil {
		i.httpListenerOptions = make(map[string]*HTTPListenerOptionWrapper)
	}
	i.httpListenerOptions[h.Index()] = h
}

func (i *Instance) AddListenerOption(l *ListenerOptionWrapper) {
	if i.listenerOptions == nil {
		i.listenerOptions = make(map[string]*ListenerOptionWrapper)
	}
	i.listenerOptions[l.Index()] = l
}

func (i *Instance) AddYamlObject(y *YAMLWrapper) {
	if i.yamlObjects == nil {
		i.yamlObjects = []*YAMLWrapper{}
	}
	i.yamlObjects = append(i.yamlObjects, y)
}

// NameNamespaceIndex is the unique identifier that identifies a specific one
func NameNamespaceIndex(name string, namespace string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
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
		matches, err := translator.HttpGatewayContainsVirtualService(gateway.Gateway.Spec.GetHttpGateway(), vs, gateway.Gateway.Spec.Ssl)
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
		response = append(response, i.routeTables[NameNamespaceIndex(rt.GetMetadata().GetName(), rt.GetMetadata().GetNamespace())])
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
			split := strings.Split(vsName, "/")
			if split[0] == namespace {
				response = append(response, virtualService)
			}
		}
	} else if namespace == "" {
		for vsName, virtualService := range i.virtualServices {
			split := strings.Split(vsName, "/")
			if split[1] == name {
				response = append(response, virtualService)
			}
		}
	} else {
		vs := i.virtualServices[NameNamespaceIndex(namespace, name)]
		if vs != nil {
			response = append(response, vs)
		}
	}
	return response
}

func (i *Instance) GetUpstream(name string, namespace string) *UpstreamWrapper {
	return i.Upstreams()[NameNamespaceIndex(name, namespace)]
}
