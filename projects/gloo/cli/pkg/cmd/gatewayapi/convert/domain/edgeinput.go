package domain

//
//type GlooEdgeCache struct {
//	YamlObjects         []*YAMLWrapper
//	Settings            map[string]*SettingsWrapper
//	RouteTables         map[string]*RouteTableWrapper
//	RouteOptions        map[string]*RouteOptionWrapper
//	ListenerOptions     map[string]*ListenerOptionWrapper
//	HTTPListenerOptions map[string]*HTTPListenerOptionWrapper
//	Upstreams           map[string]*UpstreamWrapper
//	VirtualServices     map[string]*VirtualServiceWrapper
//	GlooGateways        map[string]*GlooGatewayWrapper
//	AuthConfigs         map[string]*AuthConfigWrapper
//
//	//Gateway API Configs
//	HTTPRoutes         map[string]*HTTPRouteWrapper
//	DirectResponses    map[string]*DirectResponseWrapper
//	ListenerSets       map[string]*ListenerSetWrapper
//	VirtualHostOptions map[string]*VirtualHostOptionWrapper
//	Gateways           map[string]*GatewayWrapper
//	GatewayParameters  map[string]*GatewayParametersWrapper
//}
//
//func (g *GlooEdgeCache) AddGatewayParameters(w *GatewayParametersWrapper) {
//	if g.GatewayParameters == nil {
//		g.GatewayParameters = make(map[string]*GatewayParametersWrapper)
//	}
//	g.GatewayParameters[w.NameIndex()] = w
//}
//func (g *GlooEdgeCache) AddGateway(route *GatewayWrapper) {
//	if g.Gateways == nil {
//		g.Gateways = make(map[string]*GatewayWrapper)
//	}
//	g.Gateways[route.NameIndex()] = route
//}
//func (g *GlooEdgeCache) AddHTTPRoute(route *HTTPRouteWrapper) {
//	if g.HTTPRoutes == nil {
//		g.HTTPRoutes = make(map[string]*HTTPRouteWrapper)
//	}
//	g.HTTPRoutes[route.NameIndex()] = route
//}
//
//func (g *GlooEdgeCache) AddListenerSet(l *ListenerSetWrapper) {
//	if g.ListenerSets == nil {
//		g.ListenerSets = make(map[string]*ListenerSetWrapper)
//	}
//	g.ListenerSets[l.NameIndex()] = l
//}
//
//func (g *GlooEdgeCache) AddDirectResponse(d *DirectResponseWrapper) {
//	if g.DirectResponses == nil {
//		g.DirectResponses = make(map[string]*DirectResponseWrapper)
//	}
//	g.DirectResponses[d.NameIndex()] = d
//}
//
//func (g *GlooEdgeCache) VirtualServicesByNameNamespace(name, namespace string) []*VirtualServiceWrapper {
//	var response []*VirtualServiceWrapper
//	if name == "" && namespace == "" {
//		return response
//	}
//	if name == "" {
//		// namespace only search
//		for vsName, virtualService := range g.VirtualServices {
//			split := strings.Split(vsName, "/")
//			if split[0] == namespace {
//				response = append(response, virtualService)
//			}
//		}
//	} else if namespace == "" {
//		for vsName, virtualService := range g.VirtualServices {
//			split := strings.Split(vsName, "/")
//			if split[1] == name {
//				response = append(response, virtualService)
//			}
//		}
//	} else {
//		vs := g.VirtualServices[NameNamespaceIndex(namespace, name)]
//		if vs != nil {
//			response = append(response, vs)
//		}
//	}
//	return response
//}
//func (g *GlooEdgeCache) GlooGatewayContainsVirtualService(gateway *GlooGatewayWrapper) ([]*VirtualServiceWrapper, error) {
//	var response []*VirtualServiceWrapper
//
//	for _, virtualService := range g.VirtualServices {
//
//		var matches bool
//		var err error
//		// this logic reqires the VS to have a namespace and metadata on the inner object which we dont have. So we do the namespace lookup ourselves
//		vs := &virtualService.Spec
//		vs.SetMetadata(&core.Metadata{
//			Name:        virtualService.Name,
//			Namespace:   virtualService.Namespace,
//			Labels:      virtualService.Labels,
//			Annotations: virtualService.Annotations,
//		})
//		matches, err = translator.HttpGatewayContainsVirtualService(gateway.Gateway.Spec.GetHttpGateway(), vs, gateway.Gateway.Spec.Ssl)
//		if err != nil {
//			return nil, err
//		}
//		//}
//
//		if matches {
//			response = append(response, virtualService)
//		}
//	}
//	return response, nil
//}
//
//func (g *GlooEdgeCache) AddYamlObject(w *YAMLWrapper) {
//	if g.YamlObjects == nil {
//		g.YamlObjects = []*YAMLWrapper{}
//	}
//	g.YamlObjects = append(g.YamlObjects, w)
//}
//func (g *GlooEdgeCache) AddSettings(s *SettingsWrapper) {
//	if g.Settings == nil {
//		g.Settings = map[string]*SettingsWrapper{}
//	}
//	g.Settings[s.NameIndex()] = s
//}
//func (g *GlooEdgeCache) AddRouteTable(r *RouteTableWrapper) {
//	if g.RouteTables == nil {
//		g.RouteTables = map[string]*RouteTableWrapper{}
//	}
//	g.RouteTables[r.NameIndex()] = r
//}
//func (g *GlooEdgeCache) AddRouteOption(r *RouteOptionWrapper) {
//	if g.RouteOptions == nil {
//		g.RouteOptions = map[string]*RouteOptionWrapper{}
//	}
//	g.RouteOptions[r.NameIndex()] = r
//}
//func (g *GlooEdgeCache) AddListenerOption(l *ListenerOptionWrapper) {
//	if g.ListenerOptions == nil {
//		g.ListenerOptions = map[string]*ListenerOptionWrapper{}
//	}
//	g.ListenerOptions[l.NameIndex()] = l
//}
//func (g *GlooEdgeCache) AddHTTPListenerOption(h *HTTPListenerOptionWrapper) {
//	if g.HTTPListenerOptions == nil {
//		g.HTTPListenerOptions = map[string]*HTTPListenerOptionWrapper{}
//	}
//	g.HTTPListenerOptions[h.NameIndex()] = h
//}
//func (g *GlooEdgeCache) AddVirtualHostOption(v *VirtualHostOptionWrapper) {
//	if g.VirtualHostOptions == nil {
//		g.VirtualHostOptions = map[string]*VirtualHostOptionWrapper{}
//	}
//	g.VirtualHostOptions[v.NameIndex()] = v
//}
//func (g *GlooEdgeCache) AddUpstream(u *UpstreamWrapper) {
//	if g.Upstreams == nil {
//		g.Upstreams = map[string]*UpstreamWrapper{}
//	}
//	g.Upstreams[u.NameIndex()] = u
//}
//
//func (g *GlooEdgeCache) AddVirtualService(v *VirtualServiceWrapper) {
//	if g.VirtualServices == nil {
//		g.VirtualServices = map[string]*VirtualServiceWrapper{}
//	}
//	g.VirtualServices[v.NameIndex()] = v
//}
//func (g *GlooEdgeCache) AddGlooGateway(w *GlooGatewayWrapper) {
//	if g.GlooGateways == nil {
//		g.GlooGateways = map[string]*GlooGatewayWrapper{}
//	}
//	g.GlooGateways[w.NameIndex()] = w
//}
//func (g *GlooEdgeCache) AddAuthConfig(a *AuthConfigWrapper) {
//	if g.AuthConfigs == nil {
//		g.AuthConfigs = map[string]*AuthConfigWrapper{}
//	}
//	g.AuthConfigs[a.NameIndex()] = a
//}
//
//func NameNamespaceIndex(name string, namespace string) string {
//	return fmt.Sprintf("%s/%s", namespace, name)
//}
//
//func (g *GlooEdgeCache) GetUpstream(name string, namespace string) *UpstreamWrapper {
//	return g.Upstreams[NameNamespaceIndex(name, namespace)]
//}
