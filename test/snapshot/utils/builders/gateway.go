package builders

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// KubernetesGatewayBuilder contains options for building k8s Gateways
type KubernetesGatewayBuilder struct {
	name      string
	namespace string
	labels    map[string]string

	gatewayClassName gwv1.ObjectName
	listeners        []gwv1.Listener
}

func NewKubernetesGatewayBuilder() *KubernetesGatewayBuilder {
	return &KubernetesGatewayBuilder{}
}

func BuilderFromGateway(gw *gwv1.Gateway) *KubernetesGatewayBuilder {
	builder := &KubernetesGatewayBuilder{
		name:             gw.GetName(),
		namespace:        gw.GetNamespace(),
		labels:           gw.GetLabels(),
		gatewayClassName: gw.Spec.GatewayClassName,
		listeners:        gw.Spec.Listeners,
	}
	return builder
}

func (b *KubernetesGatewayBuilder) WithName(name string) *KubernetesGatewayBuilder {
	b.name = name
	return b
}

func (b *KubernetesGatewayBuilder) WithNamespace(namespace string) *KubernetesGatewayBuilder {
	b.namespace = namespace
	return b
}

func (b *KubernetesGatewayBuilder) WithLabel(key, value string) *KubernetesGatewayBuilder {
	b.labels[key] = value
	return b
}

func (b *KubernetesGatewayBuilder) WithGatewayClassName(gatewayClassName gwv1.ObjectName) *KubernetesGatewayBuilder {
	b.gatewayClassName = gatewayClassName
	return b
}

func (b *KubernetesGatewayBuilder) WithListeners(listener []gwv1.Listener) *KubernetesGatewayBuilder {
	b.listeners = append(b.listeners, listener...)
	return b
}

func (b *KubernetesGatewayBuilder) Build() *gwv1.Gateway {
	gw := &gwv1.Gateway{
		ObjectMeta: metav1.ObjectMeta{
			Name:      b.name,
			Namespace: b.namespace,
			Labels:    b.labels,
		},
		Spec: gwv1.GatewaySpec{
			GatewayClassName: b.gatewayClassName,
			Listeners:        b.listeners,
		},
	}
	return gw
}
