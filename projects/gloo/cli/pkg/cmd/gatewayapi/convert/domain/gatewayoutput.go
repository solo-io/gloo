package domain

type GatewayAPICache struct {
	YamlObjects         []*YAMLWrapper
	HTTPRoutes          map[string]*HTTPRouteWrapper
	RouteOptions        map[string]*RouteOptionWrapper
	VirtualHostOptions  map[string]*VirtualHostOptionWrapper
	ListenerOptions     map[string]*ListenerOptionWrapper
	HTTPListenerOptions map[string]*HTTPListenerOptionWrapper
	DirectResponses     map[string]*DirectResponseWrapper

	Upstreams    map[string]*UpstreamWrapper
	AuthConfigs  map[string]*AuthConfigWrapper
	Gateways     map[string]*GatewayWrapper
	ListenerSets map[string]*ListenerSetWrapper
	Settings     map[string]*SettingsWrapper
}

func (g *GatewayAPICache) AddSettings(s *SettingsWrapper) {
	if g.Settings == nil {
		g.Settings = make(map[string]*SettingsWrapper)
	}
	g.Settings[s.NameIndex()] = s
}
func (g *GatewayAPICache) GetGateway(name string, namespace string) *GatewayWrapper {
	if g.Gateways == nil {
		return nil
	}
	return g.Gateways[NameNamespaceIndex(name, namespace)]
}
func (g *GatewayAPICache) AddGateway(gw *GatewayWrapper) {
	if g.Gateways == nil {
		g.Gateways = make(map[string]*GatewayWrapper)
	}
	g.Gateways[gw.NameIndex()] = gw
}
func (g *GatewayAPICache) AddDirectResponse(d *DirectResponseWrapper) {
	if g.DirectResponses == nil {
		g.DirectResponses = make(map[string]*DirectResponseWrapper)
	}
	g.DirectResponses[d.NameIndex()] = d
}
func (g *GatewayAPICache) AddYAML(y *YAMLWrapper) {
	if g.YamlObjects == nil {
		g.YamlObjects = []*YAMLWrapper{}
	}
	g.YamlObjects = append(g.YamlObjects, y)
}

func (g *GatewayAPICache) AddHTTPRoute(route *HTTPRouteWrapper) {
	if g.HTTPRoutes == nil {
		g.HTTPRoutes = make(map[string]*HTTPRouteWrapper)
	}
	g.HTTPRoutes[route.NameIndex()] = route
}
func (g *GatewayAPICache) AddRouteOption(r *RouteOptionWrapper) {
	if g.RouteOptions == nil {
		g.RouteOptions = make(map[string]*RouteOptionWrapper)
	}
	g.RouteOptions[r.NameIndex()] = r
}
func (g *GatewayAPICache) AddVirtualHostOption(v *VirtualHostOptionWrapper) {
	if g.VirtualHostOptions == nil {
		g.VirtualHostOptions = make(map[string]*VirtualHostOptionWrapper)
	}
	g.VirtualHostOptions[v.NameIndex()] = v
}
func (g *GatewayAPICache) AddListenerOption(l *ListenerOptionWrapper) {
	if g.ListenerOptions == nil {
		g.ListenerOptions = make(map[string]*ListenerOptionWrapper)
	}
	g.ListenerOptions[l.NameIndex()] = l
}
func (g *GatewayAPICache) AddHTTPListenerOption(h *HTTPListenerOptionWrapper) {
	if g.HTTPListenerOptions == nil {
		g.HTTPListenerOptions = make(map[string]*HTTPListenerOptionWrapper)
	}
	g.HTTPListenerOptions[h.NameIndex()] = h
}
func (g *GatewayAPICache) AddUpstream(u *UpstreamWrapper) {
	if g.Upstreams == nil {
		g.Upstreams = make(map[string]*UpstreamWrapper)
	}
	g.Upstreams[u.NameIndex()] = u
}
func (g *GatewayAPICache) AddAuthConfig(a *AuthConfigWrapper) {
	if g.AuthConfigs == nil {
		g.AuthConfigs = make(map[string]*AuthConfigWrapper)
	}
	g.AuthConfigs[a.NameIndex()] = a
}
func (g *GatewayAPICache) AddListenerSet(l *ListenerSetWrapper) {
	if g.ListenerSets == nil {
		g.ListenerSets = make(map[string]*ListenerSetWrapper)
	}
	g.ListenerSets[l.NameIndex()] = l
}
