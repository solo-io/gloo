package builders

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRouteBuilder simplifies the process of generating HTTPRoutes in tests for client.Object type
type HTTPRouteBuilder struct {
	name      string
	namespace string
	labels    map[string]string

	commonRouteSpec gwv1.CommonRouteSpec
	hostnames       []gwv1.Hostname
	httpRouteRules  []gwv1.HTTPRouteRule
}

// BuilderFromHTTPRoute creates a new HTTPRouteBuilder from an existing HTTPRoute
func BuilderFromHTTPRoute(httpRoute *gwv1.HTTPRoute) *HTTPRouteBuilder {

	builder := &HTTPRouteBuilder{
		name:            httpRoute.GetName(),
		namespace:       httpRoute.GetNamespace(),
		labels:          httpRoute.GetLabels(),
		commonRouteSpec: httpRoute.Spec.CommonRouteSpec,
		hostnames:       httpRoute.Spec.Hostnames,
		httpRouteRules:  httpRoute.Spec.Rules,
	}
	return builder
}

// NewHTTPRouteBuilder creates an empty HTTPRouteBuilder
func NewHTTPRouteBuilder() *HTTPRouteBuilder {
	return &HTTPRouteBuilder{}
}

func (b *HTTPRouteBuilder) WithName(name string) *HTTPRouteBuilder {
	b.name = name
	return b
}

func (b *HTTPRouteBuilder) WithNamespace(namespace string) *HTTPRouteBuilder {
	b.namespace = namespace
	return b
}

func (b *HTTPRouteBuilder) WithLabel(key, value string) *HTTPRouteBuilder {
	b.labels[key] = value
	return b
}

func (b *HTTPRouteBuilder) WithHostnames(hostnames []string) *HTTPRouteBuilder {
	hosts := make([]gwv1.Hostname, len(hostnames))
	for i, hostname := range hostnames {
		hosts[i] = gwv1.Hostname(hostname)
	}
	b.hostnames = hosts
	return b
}

func (b *HTTPRouteBuilder) WithCommonRoute(commonRoute gwv1.CommonRouteSpec) *HTTPRouteBuilder {
	b.commonRouteSpec = commonRoute
	return b
}

func (b *HTTPRouteBuilder) WithHTTPRouteRule(rule gwv1.HTTPRouteRule) *HTTPRouteBuilder {
	b.httpRouteRules = append(b.httpRouteRules, rule)
	return b
}

func (b *HTTPRouteBuilder) WithHTTPRouteRules(rules []gwv1.HTTPRouteRule) *HTTPRouteBuilder {
	b.httpRouteRules = append(b.httpRouteRules, rules...)
	return b
}

func (b *HTTPRouteBuilder) Build() *gwv1.HTTPRoute {
	httpRoute := &gwv1.HTTPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
			Labels:    b.labels,
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: b.commonRouteSpec,
			Hostnames:       b.hostnames,
			Rules:           b.httpRouteRules,
		},
	}
	return httpRoute
}
