package sslutils

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/go-utils/contextutils"
	"google.golang.org/protobuf/types/known/wrapperspb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/cert"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// Gateway API has an extension point for implementation specific tls settings, they can be found [here](https://gateway-api.sigs.k8s.io/guides/tls/#extensions)
const (
	GatewaySslOptionsPrefix = wellknown.GatewayAnnotationPrefix + "/ssl"

	GatewaySslCipherSuites         = GatewaySslOptionsPrefix + "/cipher-suites"
	GatewaySslEcdhCurves           = GatewaySslOptionsPrefix + "/ecdh-curves"
	GatewaySslMinimumTlsVersion    = GatewaySslOptionsPrefix + "/minimum-tls-version"
	GatewaySslMaximumTlsVersion    = GatewaySslOptionsPrefix + "/maximum-tls-version"
	GatewaySslOneWayTls            = GatewaySslOptionsPrefix + "/one-way-tls"
	GatewaySslVerifySubjectAltName = GatewaySslOptionsPrefix + "/verify-subject-alt-name"
)

var (
	InvalidTlsSecretError = func(secret *corev1.Secret, err error) error {
		errorString := fmt.Sprintf("%v.%v is not a valid TLS secret", secret.Namespace, secret.Name)
		return eris.Wrap(err, errorString)
	}

	NoCertificateFoundError = eris.New("no certificate information found")
)

// ValidateTlsSecret and return a cleaned cert
func ValidateTlsSecret(sslSecret *corev1.Secret) (cleanedCertChain string, err error) {

	certChain := string(sslSecret.Data[corev1.TLSCertKey])
	privateKey := string(sslSecret.Data[corev1.TLSPrivateKeyKey])
	rootCa := string(sslSecret.Data[corev1.ServiceAccountRootCAKey])

	cleanedCertChain, err = cleanedSslKeyPair(certChain, privateKey, rootCa)
	if err != nil {
		err = InvalidTlsSecretError(sslSecret, err)
	}
	return cleanedCertChain, err
}

func cleanedSslKeyPair(certChain, privateKey, rootCa string) (cleanedChain string, err error) {

	// in the case where we _only_ provide a rootCa, we do not want to validate tls.key+tls.cert
	if (certChain == "") && (privateKey == "") && (rootCa != "") {
		return certChain, nil
	}

	// validate that the cert and key are a valid pair
	_, err = tls.X509KeyPair([]byte(certChain), []byte(privateKey))
	if err != nil {
		return "", err
	}

	// validate that the parsed piece is valid
	// this is still faster than a call out to openssl despite this second parsing pass of the cert
	// pem parsing in go is permissive while envoy is not
	// this might not be needed once we have larger envoy validation
	candidateCert, err := cert.ParseCertsPEM([]byte(certChain))
	if err != nil {
		// return err rather than sanitize. This is to maintain UX with older versions and to keep in line with gateway2 pkg.
		return "", err
	}
	cleanedChainBytes, err := cert.EncodeCertificates(candidateCert...)
	cleanedChain = string(cleanedChainBytes)

	return cleanedChain, err
}

type SslExtensionOptionFunc = func(ctx context.Context, in string, out *ssl.SslConfig) error

func ApplyCipherSuites(ctx context.Context, in string, out *ssl.SslConfig) error {
	if out.GetParameters() == nil {
		out.Parameters = &ssl.SslParameters{}
	}
	cipherSuites := strings.Split(in, ",")
	out.GetParameters().CipherSuites = cipherSuites
	return nil
}

func ApplyEcdhCurves(ctx context.Context, in string, out *ssl.SslConfig) error {
	if out.GetParameters() == nil {
		out.Parameters = &ssl.SslParameters{}
	}
	ecdhCurves := strings.Split(in, ",")
	out.GetParameters().EcdhCurves = ecdhCurves
	return nil
}

func ApplyMinimumTlsVersion(ctx context.Context, in string, out *ssl.SslConfig) error {
	if out.GetParameters() == nil {
		out.Parameters = &ssl.SslParameters{}
	}
	if parsed, ok := ssl.SslParameters_ProtocolVersion_value[in]; ok {
		out.GetParameters().MinimumProtocolVersion = ssl.SslParameters_ProtocolVersion(parsed)
		if out.GetParameters().GetMaximumProtocolVersion() != ssl.SslParameters_TLS_AUTO && out.GetParameters().GetMaximumProtocolVersion() < out.GetParameters().GetMinimumProtocolVersion() {
			err := eris.Errorf("maximum tls version %s is less than minimum tls version %s", out.GetParameters().GetMaximumProtocolVersion().String(), in)
			out.GetParameters().MaximumProtocolVersion = ssl.SslParameters_TLS_AUTO
			out.GetParameters().MinimumProtocolVersion = ssl.SslParameters_TLS_AUTO
			return err
		}
	} else {
		return eris.Errorf("invalid minimum tls version: %s", in)
	}
	return nil
}

func ApplyMaximumTlsVersion(ctx context.Context, in string, out *ssl.SslConfig) error {
	if out.GetParameters() == nil {
		out.Parameters = &ssl.SslParameters{}
	}
	if parsed, ok := ssl.SslParameters_ProtocolVersion_value[in]; ok {
		out.GetParameters().MaximumProtocolVersion = ssl.SslParameters_ProtocolVersion(parsed)
		if out.GetParameters().GetMaximumProtocolVersion() != ssl.SslParameters_TLS_AUTO && out.GetParameters().GetMaximumProtocolVersion() < out.GetParameters().GetMinimumProtocolVersion() {
			err := eris.Errorf("maximum tls version %s is less than minimum tls version %s", in, out.GetParameters().GetMinimumProtocolVersion().String())
			out.GetParameters().MaximumProtocolVersion = ssl.SslParameters_TLS_AUTO
			out.GetParameters().MinimumProtocolVersion = ssl.SslParameters_TLS_AUTO
			return err
		}
	} else {
		return eris.Errorf("invalid maximum tls version: %s", in)
	}
	return nil
}

func ApplyOneWayTls(ctx context.Context, in string, out *ssl.SslConfig) error {
	if strings.ToLower(in) == "true" {
		out.OneWayTls = wrapperspb.Bool(true)
	} else if strings.ToLower(in) == "false" {
		out.OneWayTls = wrapperspb.Bool(false)
	} else {
		return eris.Errorf("invalid value for one-way-tls: %s", in)
	}
	return nil
}

func ApplyVerifySubjectAltName(ctx context.Context, in string, out *ssl.SslConfig) error {
	altNames := strings.Split(in, ",")
	out.VerifySubjectAltName = altNames
	return nil
}

var SslExtensionOptionFuncs = map[string]SslExtensionOptionFunc{
	GatewaySslCipherSuites:         ApplyCipherSuites,
	GatewaySslEcdhCurves:           ApplyEcdhCurves,
	GatewaySslMinimumTlsVersion:    ApplyMinimumTlsVersion,
	GatewaySslMaximumTlsVersion:    ApplyMaximumTlsVersion,
	GatewaySslOneWayTls:            ApplyOneWayTls,
	GatewaySslVerifySubjectAltName: ApplyVerifySubjectAltName,
}

// ApplySslExtensionOptions applies the GatewayTLSConfig options to the SslConfig
// This function will never exit early, even if an error is encountered.
// It will apply all options and log all errors encountered.
func ApplySslExtensionOptions(ctx context.Context, in *gwv1.GatewayTLSConfig, out *ssl.SslConfig) {
	var wrapped error
	for key, _ := range in.Frontend.PerPort {
		if extensionFunc, ok := SslExtensionOptionFuncs[string(key)]; ok {
			if err := extensionFunc(ctx, "(nick) deleted this to be able to compile code", out); err != nil {
				wrapped = multierror.Append(wrapped, err)
			}
		} else {
			wrapped = multierror.Append(wrapped, eris.Errorf("unknown ssl option: %s", key))
		}
	}

	if wrapped != nil {
		contextutils.LoggerFrom(ctx).Warnf("error applying ssl extension options: %v", wrapped)
	}
}
