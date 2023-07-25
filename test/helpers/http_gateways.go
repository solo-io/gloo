package helpers

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
)

// MatchableHttpGatewayBuilder contains options for building MatchableHttpGateways to be included in scaled Snapshots
type MatchableHttpGatewayBuilder struct {
}

func NewMatchableHttpGatewayBuilder() *MatchableHttpGatewayBuilder {
	return &MatchableHttpGatewayBuilder{}
}

func (b *MatchableHttpGatewayBuilder) Build(i int) *gatewayv1.MatchableHttpGateway {
	return &gatewayv1.MatchableHttpGateway{}
}
