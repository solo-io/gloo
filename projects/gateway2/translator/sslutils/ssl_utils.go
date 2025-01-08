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

	GatewaySslCipherSuites   = GatewaySslOptionsPrefix + "/cipher-suites"
	GatewaySslMinimumTlsVersion    = GatewaySslOptionsPrefix + "/minimum-tls-version"
	GatewaySslMaximumTlsVersion    = GatewaySslOptionsPrefix + "/maximum-tls-version"
	GatewaySslOneWayTls            = GatewaySslOptionsPrefix + "/one-way-tls"
	GatewaySslVerifySubjectAltName = GatewaySslOptionsPrefix + "/verify-subject-alt-name"
)

var (
	InvalidTlsSecretError = func(secret *corev1.Secret, err error) error {
		errorString := fmt.Sprintf("%v.%v is not a valid TLS secret", secret.Namespace, secret.Name)
		return eris.Wrapf(err, errorString)
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

// ApplySslExtensionOptions applies the GatewayTLSConfig options to the SslConfig
// This function will never exit early, even if an error is encountered.
// It will apply all options and return all errors encountered.
func ApplySslExtensionOptions(ctx context.Context, in *gwv1.GatewayTLSConfig, out *ssl.SslConfig) error {
	if len(in.Options) == 0 {
		return nil
	}

	var err error

	if oneWayTls, ok := in.Options[GatewaySslOneWayTls]; ok {
		if strings.ToLower(string(oneWayTls)) == "true" {
			out.OneWayTls = wrapperspb.Bool(true)
		}
	}

	if verifySubjectAltNameStr, ok := in.Options[GatewaySslVerifySubjectAltName]; ok {
		altNames := strings.Split(string(verifySubjectAltNameStr), ",")
		out.VerifySubjectAltName = altNames
	}

	out.Parameters = &ssl.SslParameters{}
	if cipherSuitesStr, ok := in.Options[GatewaySslCipherSuitesOption]; ok {
		cipherSuites := strings.Split(string(cipherSuitesStr), ",")
		out.Parameters.CipherSuites = cipherSuites
	}

	if minTlsVersion, ok := in.Options[GatewaySslMinimumTlsVersion]; ok {
		if parsed, ok := ssl.SslParameters_ProtocolVersion_value[string(minTlsVersion)]; ok {
			out.Parameters.MinimumProtocolVersion = ssl.SslParameters_ProtocolVersion(parsed)
		} else {
			contextutils.LoggerFrom(ctx).Debugf("invalid minimum tls version: %s", minTlsVersion)
			err = multierror.Append(err, eris.Errorf("invalid minimum tls version: %s", minTlsVersion))
		}
	}

	if maxTlsVersion, ok := in.Options[GatewaySslMaximumTlsVersion]; ok {
		if parsed, ok := ssl.SslParameters_ProtocolVersion_value[string(maxTlsVersion)]; ok {
			out.Parameters.MaximumProtocolVersion = ssl.SslParameters_ProtocolVersion(parsed)
		} else {
			contextutils.LoggerFrom(ctx).Debugf("invalid maximum tls version: %s", maxTlsVersion)
			err = multierror.Append(err, eris.Errorf("invalid maximum tls version: %s", maxTlsVersion))
		}
	}
	return err
}
