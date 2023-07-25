package helpers

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
)

// GatewayBuilder contains options for building Gateways to be included in scaled Snapshots
type GatewayBuilder struct {
}

func NewGatewayBuilder() *GatewayBuilder {
	return &GatewayBuilder{}
}

func (b *GatewayBuilder) Build(i int) *gatewayv1.Gateway {
	return &gatewayv1.Gateway{}
}
