package builders

import (
	errors "github.com/rotisserie/eris"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	soloapisv1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualServiceBuilder simplifies the process of generating VirtualServices in tests for client.Object type
type VirtualServiceBuilder struct {
	name      string
	namespace string
	labels    map[string]string

	domains            []string
	virtualHostOptions *soloapisv1.VirtualHostOptions
	routesByName       map[string]*v1.Route
	sslConfig          *soloapisv1.SslConfig
}

// BuilderFromVirtualService creates a new VirtualServiceBuilder from an existing VirtualService
func BuilderFromVirtualService(vs *v1.VirtualService) *VirtualServiceBuilder {
	builder := &VirtualServiceBuilder{
		name:               vs.GetName(),
		namespace:          vs.GetNamespace(),
		labels:             vs.GetLabels(),
		domains:            vs.Spec.GetVirtualHost().GetDomains(),
		virtualHostOptions: vs.Spec.GetVirtualHost().GetOptions(),
		sslConfig:          vs.Spec.GetSslConfig(),
		routesByName:       make(map[string]*v1.Route, len(vs.Spec.GetVirtualHost().GetRoutes())),
	}
	for _, r := range vs.Spec.GetVirtualHost().GetRoutes() {
		builder.WithRoute(r.GetName(), r)
	}
	return builder
}

// NewVirtualServiceBuilder creates an empty VirtualServiceBuilder
func NewVirtualServiceBuilder() *VirtualServiceBuilder {
	return &VirtualServiceBuilder{
		labels:       make(map[string]string, 2),
		routesByName: make(map[string]*v1.Route, 10),
	}
}

func (b *VirtualServiceBuilder) WithSslConfig(sslConfig *soloapisv1.SslConfig) *VirtualServiceBuilder {
	b.sslConfig = sslConfig
	return b
}

func (b *VirtualServiceBuilder) WithName(name string) *VirtualServiceBuilder {
	b.name = name
	return b
}

func (b *VirtualServiceBuilder) WithNamespace(namespace string) *VirtualServiceBuilder {
	b.namespace = namespace
	return b
}

func (b *VirtualServiceBuilder) WithLabel(key, value string) *VirtualServiceBuilder {
	b.labels[key] = value
	return b
}

func (b *VirtualServiceBuilder) WithDomain(domain string) *VirtualServiceBuilder {
	return b.WithDomains([]string{domain})
}

func (b *VirtualServiceBuilder) WithDomains(domains []string) *VirtualServiceBuilder {
	b.domains = domains
	return b
}

func (b *VirtualServiceBuilder) WithVirtualHostOptions(virtualHostOptions *soloapisv1.VirtualHostOptions) *VirtualServiceBuilder {
	b.virtualHostOptions = virtualHostOptions
	return b
}

func (b *VirtualServiceBuilder) WithRouteOptions(routeName string, routeOptions *soloapisv1.RouteOptions) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Options = routeOptions
	})
}

func (b *VirtualServiceBuilder) WithRoute(routeName string, route *v1.Route) *VirtualServiceBuilder {
	b.routesByName[routeName] = route
	return b
}

func (b *VirtualServiceBuilder) getOrDefaultRoute(routeName string) *v1.Route {
	route, ok := b.routesByName[routeName]
	if !ok {
		return &v1.Route{
			Name: routeName,
		}
	}
	return route
}

func (b *VirtualServiceBuilder) WithRouteActionToUpstream(routeName string, upstream *gloov1.Upstream) *VirtualServiceBuilder {
	return b.WithRouteActionToUpstreamRef(routeName, upstream.GetMetadata().Ref())
}

func (b *VirtualServiceBuilder) WithRouteActionToUpstreamRef(routeName string, upstreamRef *core.ResourceRef) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Action = &v1.Route_RouteAction{
			RouteAction: &soloapisv1.RouteAction{
				Destination: &soloapisv1.RouteAction_Single{
					Single: &soloapisv1.Destination{
						DestinationType: &soloapisv1.Destination_Upstream{
							Upstream: upstreamRef,
						},
					},
				},
			},
		}
	})
}

func (b *VirtualServiceBuilder) WithRouteDelegateActionRef(routeName string, delegateRef *core.ResourceRef) *VirtualServiceBuilder {
	return b.WithRouteDelegateAction(routeName,
		&v1.DelegateAction{
			DelegationType: &v1.DelegateAction_Ref{
				Ref: delegateRef,
			},
		})
}

func (b *VirtualServiceBuilder) WithRouteDelegateActionSelector(routeName string, delegateSelector *v1.RouteTableSelector) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Action = &v1.Route_DelegateAction{
			DelegateAction: &v1.DelegateAction{
				DelegationType: &v1.DelegateAction_Selector{
					Selector: delegateSelector,
				},
			},
		}
	})
}

func (b *VirtualServiceBuilder) WithRouteDelegateAction(routeName string, delegateAction *v1.DelegateAction) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Action = &v1.Route_DelegateAction{
			DelegateAction: delegateAction,
		}
	})
}

func (b *VirtualServiceBuilder) WithRouteAction(routeName string, routeAction *soloapisv1.RouteAction) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Action = &v1.Route_RouteAction{
			RouteAction: routeAction,
		}
	})
}

func (b *VirtualServiceBuilder) WithRouteActionToSingleDestination(routeName string, destination *soloapisv1.Destination) *VirtualServiceBuilder {
	return b.WithRouteAction(routeName, &soloapisv1.RouteAction{
		Destination: &soloapisv1.RouteAction_Single{
			Single: destination,
		},
	})
}

func (b *VirtualServiceBuilder) WithRouteActionToMultiDestination(routeName string, destination *soloapisv1.MultiDestination) *VirtualServiceBuilder {
	return b.WithRouteAction(routeName, &soloapisv1.RouteAction{
		Destination: &soloapisv1.RouteAction_Multi{
			Multi: destination,
		},
	})
}

func (b *VirtualServiceBuilder) WithRouteDirectResponseAction(routeName string, action *soloapisv1.DirectResponseAction) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Action = &v1.Route_DirectResponseAction{
			DirectResponseAction: action,
		}
	})
}

func (b *VirtualServiceBuilder) WithRoutePrefixMatcher(routeName string, prefixMatch string) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Matchers = []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: prefixMatch,
			},
		}}
	})
}

func (b *VirtualServiceBuilder) WithRouteMatcher(routeName string, matcher *matchers.Matcher) *VirtualServiceBuilder {
	return b.WithRouteMutation(routeName, func(route *v1.Route) {
		route.Matchers = []*matchers.Matcher{matcher}
	})
}

func (b *VirtualServiceBuilder) WithRouteMutation(routeName string, mutation func(route *v1.Route)) *VirtualServiceBuilder {
	route := b.getOrDefaultRoute(routeName)
	mutation(route)
	return b.WithRoute(routeName, route)
}

func (b *VirtualServiceBuilder) errorIfInvalid() error {
	if len(b.domains) == 0 {
		// Unset domains will behave like a wildcard "*", which contributes to test flakes
		return errors.New("attempting to not set a VirtualService domain")

	}
	for _, domain := range b.domains {
		if domain == "*" {
			// Wildcard domains contribute to test flakes
			return errors.New("attempting to set * as a VirtualService domain")
		}
	}
	return nil
}

func (b *VirtualServiceBuilder) Build() (*v1.VirtualService, error) {
	if err := b.errorIfInvalid(); err != nil {
		// We error loudly here
		// These types of errors are intended to prevent developers from creating resources
		// which are semantically correct, but lead to test flakes/confusion
		return nil, err
	}

	var routes []*v1.Route
	for _, route := range b.routesByName {
		routes = append(routes, route)
	}

	vs := &v1.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
			Labels:    b.labels,
		},
		Spec: v1.VirtualServiceSpec{
			VirtualHost: &v1.VirtualHost{
				Domains: b.domains,
				Routes:  routes,
				Options: b.virtualHostOptions,
			},
			SslConfig: b.sslConfig,
		},
	}
	return vs, nil
}
