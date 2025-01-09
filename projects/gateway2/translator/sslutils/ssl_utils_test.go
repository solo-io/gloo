package sslutils

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/wrapperspb"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestApplySslExtensionOptions(t *testing.T) {
	testCases := []struct {
		name   string
		out    *ssl.SslConfig
		in     *gwv1.GatewayTLSConfig
		errors []string
	}{
		{
			name: "one_way_tls_true",
			out: &ssl.SslConfig{
				OneWayTls: wrapperspb.Bool(true),
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslOneWayTls: "true",
				},
			},
		},
		{
			name: "one_way_tls_true_incorrect_casing",
			out: &ssl.SslConfig{
				OneWayTls: wrapperspb.Bool(true),
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslOneWayTls: "True",
				},
			},
		},
		{
			name: "one_way_tls_false",
			out: &ssl.SslConfig{
				OneWayTls: wrapperspb.Bool(false),
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslOneWayTls: "false",
				},
			},
		},
		{
			name: "one_way_tls_false_incorrect_casing",
			out: &ssl.SslConfig{
				OneWayTls: wrapperspb.Bool(false),
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslOneWayTls: "False",
				},
			},
		},
		{
			name: "cipher_suites",
			out: &ssl.SslConfig{
				Parameters: &ssl.SslParameters{
					CipherSuites: []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
				},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslCipherSuites: "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
				},
			},
		},
		{
			name: "subject_alt_names",
			out: &ssl.SslConfig{
				VerifySubjectAltName: []string{"foo", "bar"},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslVerifySubjectAltName: "foo,bar",
				},
			},
		},
		{
			name: "tls_max_version",
			out: &ssl.SslConfig{
				Parameters: &ssl.SslParameters{
					MaximumProtocolVersion: ssl.SslParameters_TLSv1_2,
				},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslMaximumTlsVersion: "TLSv1_2",
				},
			},
		},
		{
			name: "tls_min_version",
			out: &ssl.SslConfig{
				Parameters: &ssl.SslParameters{
					MinimumProtocolVersion: ssl.SslParameters_TLSv1_3,
				},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslMinimumTlsVersion: "TLSv1_3",
				},
			},
		},
		{
			name: "invalid_tls_versions",
			out: &ssl.SslConfig{
				Parameters: &ssl.SslParameters{},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslMinimumTlsVersion: "TLSv1.3",
					GatewaySslMaximumTlsVersion: "TLSv1.2",
				},
			},
			errors: []string{
				"invalid maximum tls version: TLSv1.2",
				"invalid minimum tls version: TLSv1.3",
			},
		},
		{
			name: "maximium_tls_version_less_than_minimum",
			out: &ssl.SslConfig{
				Parameters: &ssl.SslParameters{},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslMinimumTlsVersion: "TLSv1_3",
					GatewaySslMaximumTlsVersion: "TLSv1_2",
				},
			},
			errors: []string{
				"maximum tls version TLSv1_2 is less than minimum tls version TLSv1_3",
			},
		},
		{
			name: "multiple_options",
			out: &ssl.SslConfig{
				VerifySubjectAltName: []string{"foo", "bar"},
				OneWayTls:            wrapperspb.Bool(true),
				Parameters: &ssl.SslParameters{
					MaximumProtocolVersion: ssl.SslParameters_TLSv1_3,
					MinimumProtocolVersion: ssl.SslParameters_TLSv1_2,
					CipherSuites:           []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
				},
			},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslMaximumTlsVersion:    "TLSv1_3",
					GatewaySslMinimumTlsVersion:    "TLSv1_2",
					GatewaySslVerifySubjectAltName: "foo,bar",
					GatewaySslOneWayTls:            "true",
					GatewaySslCipherSuites:         "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
				},
			},
		},
		{
			name: "misspelled_option",
			out:  &ssl.SslConfig{},
			in: &gwv1.GatewayTLSConfig{
				Options: map[gwv1.AnnotationKey]gwv1.AnnotationValue{
					GatewaySslMinimumTlsVersion + "s": "TLSv1_3",
				},
			},
			errors: []string{
				"unknown ssl option: gateway.gloo.solo.io/ssl/minimum-tls-versions",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := &zaptest.Buffer{}
			logger := zap.New(zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()),
				b,
				zapcore.DebugLevel,
			))
			ctx := contextutils.WithExistingLogger(context.Background(), logger.Sugar())
			out := &ssl.SslConfig{}
			ApplySslExtensionOptions(ctx, tc.in, out)
			assert.Empty(t, cmp.Diff(tc.out, out, protocmp.Transform()))
			if len(tc.errors) > 0 {
				for _, err := range tc.errors {
					assert.Contains(t, b.String(), err)
				}
			} else {
				assert.Empty(t, b.String())
			}
		})

	}
}
