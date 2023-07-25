package helpers

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
)

// MatchableTcpGatewayBuilder contains options for building MatchableTcpGateways to be included in scaled Snapshots
type MatchableTcpGatewayBuilder struct {
}

func NewMatchableTcpGatewayBuilder() *MatchableTcpGatewayBuilder {
	return &MatchableTcpGatewayBuilder{}
}

func (b *MatchableTcpGatewayBuilder) Build(i int) *gatewayv1.MatchableTcpGateway {
	return &gatewayv1.MatchableTcpGateway{}
}
