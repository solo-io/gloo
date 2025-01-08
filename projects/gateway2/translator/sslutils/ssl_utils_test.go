package sslutils

import (
	"context"
	"testing"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestApplySslExtensionOptions(t *testing.T) {
	testCases := []struct {
		name string
		out  *ssl.SslConfig
		in   *gwv1.GatewayTLSConfig
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ApplySslExtensionOptions(ctx, tc.in, tc.out)
		})

	}
}
