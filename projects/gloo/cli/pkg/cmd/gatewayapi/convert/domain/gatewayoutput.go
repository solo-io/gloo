package domain

//
//type GatewayAPICache struct {
//	YamlObjects         []*snapshot.YAMLWrapper
//	HTTPRoutes          map[string]*snapshot.HTTPRouteWrapper
//	RouteOptions        map[string]*snapshot.RouteOptionWrapper
//	VirtualHostOptions  map[string]*snapshot.VirtualHostOptionWrapper
//	ListenerOptions     map[string]*snapshot.ListenerOptionWrapper
//	HTTPListenerOptions map[string]*snapshot.HTTPListenerOptionWrapper
//	DirectResponses     map[string]*snapshot.DirectResponseWrapper
//
//	Upstreams    map[string]*snapshot.UpstreamWrapper
//	AuthConfigs  map[string]*snapshot.AuthConfigWrapper
//	Gateways     map[string]*snapshot.GatewayWrapper
//	ListenerSets map[string]*snapshot.ListenerSetWrapper
//	Settings     map[string]*snapshot.SettingsWrapper
//}
//
//func (g *GatewayAPICache) AddSettings(s *snapshot.SettingsWrapper) {
//	if g.Settings == nil {
//		g.Settings = make(map[string]*snapshot.SettingsWrapper)
//	}
//	g.Settings[s.Index()] = s
//}
//func (g *GatewayAPICache) GetGateway(name string, namespace string) *snapshot.GatewayWrapper {
//	if g.Gateways == nil {
//		return nil
//	}
//	return g.Gateways[snapshot.NameNamespaceIndex(name, namespace)]
//}
//func (g *GatewayAPICache) AddGateway(gw *snapshot.GatewayWrapper) {
//	if g.Gateways == nil {
//		g.Gateways = make(map[string]*snapshot.GatewayWrapper)
//	}
//	g.Gateways[gw.Index()] = gw
//}
//func (g *GatewayAPICache) AddDirectResponse(d *snapshot.DirectResponseWrapper) {
//	if g.DirectResponses == nil {
//		g.DirectResponses = make(map[string]*snapshot.DirectResponseWrapper)
//	}
//	g.DirectResponses[d.Index()] = d
//}
//func (g *GatewayAPICache) AddYAML(y *snapshot.YAMLWrapper) {
//	if g.YamlObjects == nil {
//		g.YamlObjects = []*snapshot.YAMLWrapper{}
//	}
//	g.YamlObjects = append(g.YamlObjects, y)
//}
//
//func (g *GatewayAPICache) AddHTTPRoute(route *snapshot.HTTPRouteWrapper) {
//	if g.HTTPRoutes == nil {
//		g.HTTPRoutes = make(map[string]*snapshot.HTTPRouteWrapper)
//	}
//	g.HTTPRoutes[route.Index()] = route
//}
//func (g *GatewayAPICache) AddRouteOption(r *snapshot.RouteOptionWrapper) {
//	if g.RouteOptions == nil {
//		g.RouteOptions = make(map[string]*snapshot.RouteOptionWrapper)
//	}
//	g.RouteOptions[r.Index()] = r
//}
//func (g *GatewayAPICache) AddVirtualHostOption(v *snapshot.VirtualHostOptionWrapper) {
//	if g.VirtualHostOptions == nil {
//		g.VirtualHostOptions = make(map[string]*snapshot.VirtualHostOptionWrapper)
//	}
//	g.VirtualHostOptions[v.Index()] = v
//}
//func (g *GatewayAPICache) AddListenerOption(l *snapshot.ListenerOptionWrapper) {
//	if g.ListenerOptions == nil {
//		g.ListenerOptions = make(map[string]*snapshot.ListenerOptionWrapper)
//	}
//	g.ListenerOptions[l.Index()] = l
//}
//func (g *GatewayAPICache) AddHTTPListenerOption(h *snapshot.HTTPListenerOptionWrapper) {
//	if g.HTTPListenerOptions == nil {
//		g.HTTPListenerOptions = make(map[string]*snapshot.HTTPListenerOptionWrapper)
//	}
//	g.HTTPListenerOptions[h.Index()] = h
//}
//func (g *GatewayAPICache) AddUpstream(u *snapshot.UpstreamWrapper) {
//	if g.Upstreams == nil {
//		g.Upstreams = make(map[string]*snapshot.UpstreamWrapper)
//	}
//	g.Upstreams[u.Index()] = u
//}
//func (g *GatewayAPICache) AddAuthConfig(a *snapshot.AuthConfigWrapper) {
//	if g.AuthConfigs == nil {
//		g.AuthConfigs = make(map[string]*snapshot.AuthConfigWrapper)
//	}
//	g.AuthConfigs[a.Index()] = a
//}
//func (g *GatewayAPICache) AddListenerSet(l *snapshot.ListenerSetWrapper) {
//	if g.ListenerSets == nil {
//		g.ListenerSets = make(map[string]*snapshot.ListenerSetWrapper)
//	}
//	g.ListenerSets[l.Index()] = l
//}
