package certgen

import (
	"context"
	"crypto/x509"
	"fmt"

	"github.com/rotisserie/eris"
	k8s_ar_v1 "github.com/solo-io/external-apis/pkg/api/k8s/admissionregistration.k8s.io/v1"
	k8s_core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/k8s-utils/certutils"
	"github.com/solo-io/solo-projects/projects/multicluster-admission-webhook/pkg/config"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/cert"
	"knative.dev/pkg/network"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	CaSecretError = func(err error, name, namespace string) error {
		return eris.Wrapf(err, "unable to ensure CA secret %s.%s", name, namespace)
	}
)

func NewSelfSignedWebhookCAManager(
	v1Clientset k8s_core_v1.Clientset,
	admissionV1Clientset k8s_ar_v1.Clientset,
	cfg *config.Config,
) WebhookCAManager {
	return &selfSignedWebhookCAManager{
		secretClient:        v1Clientset.Secrets(),
		webhookConfigClient: admissionV1Clientset.ValidatingWebhookConfigurations(),
		vwcName:             cfg.GetString(config.ValidatingWebhookConfigurationName),
	}
}

type selfSignedWebhookCAManager struct {
	secretClient        k8s_core_v1.SecretClient
	webhookConfigClient k8s_ar_v1.ValidatingWebhookConfigurationClient
	vwcName             string
}

func (w *selfSignedWebhookCAManager) EnsureCaCerts(
	ctx context.Context,
	secretName, secretNamespace, svcName, svcNamespace string,
) error {
	caBundle, err := w.ensureCaSecret(ctx, secretName, secretNamespace, svcName, svcNamespace)
	if err != nil {
		return CaSecretError(err, secretName, secretNamespace)
	}
	return w.ensureValidatingWebhookConfiguration(ctx, caBundle, svcName, svcNamespace)

}

func (w *selfSignedWebhookCAManager) ensureCaSecret(
	ctx context.Context,
	secretName, secretNamespace, svcName, svcNamespace string,
) ([]byte, error) {
	secret, err := w.secretClient.GetSecret(ctx, client.ObjectKey{
		Name:      secretName,
		Namespace: secretNamespace,
	})
	if err != nil {
		if errors.IsNotFound(err) {
			cert, err := genCerts(svcName, svcNamespace)
			if err != nil {
				return nil, err
			}
			caSecret := makeTlsSecret(
				secretName,
				secretNamespace,
				cert.ServerCertificate,
				cert.ServerCertKey,
				cert.CaCertificate,
			)
			if err = w.secretClient.CreateSecret(ctx, caSecret); err != nil {
				return nil, err
			}
			return cert.CaCertificate, nil
		}
		return nil, err
	}
	return secret.Data[core_v1.ServiceAccountRootCAKey], nil
}

func (w *selfSignedWebhookCAManager) ensureValidatingWebhookConfiguration(
	ctx context.Context,
	caBundle []byte,
	svcName, svcNamespace string,
) error {
	vwc, err := w.webhookConfigClient.GetValidatingWebhookConfiguration(ctx, client.ObjectKey{
		Name: w.vwcName,
	})
	if err != nil {
		return err
	}
	for i, wh := range vwc.Webhooks {
		if wh.ClientConfig.Service == nil {
			continue
		}

		// if we find a webhook cfg that targets our service, update it
		if svcName == wh.ClientConfig.Service.Name && svcNamespace == wh.ClientConfig.Service.Namespace {
			wh.ClientConfig.CABundle = caBundle
			vwc.Webhooks[i] = wh
		}
	}
	return w.webhookConfigClient.UpdateValidatingWebhookConfiguration(ctx, vwc)
}

func makeTlsSecret(
	secretName, secretNamespace string,
	cert, key, caCert []byte,
) *core_v1.Secret {
	return &core_v1.Secret{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      secretName,
			Namespace: secretNamespace,
		},
		Type: core_v1.SecretTypeTLS,
		Data: map[string][]byte{
			core_v1.TLSCertKey:              cert,
			core_v1.TLSPrivateKeyKey:        key,
			core_v1.ServiceAccountRootCAKey: caCert,
		},
	}
}

// Copy pasted from Gloo
// https://github.com/solo-io/gloo/blob/2900b083257fa23aaaa2cf8dacbebdf20651dbb1/jobs/pkg/certgen/gen_certs.go#L12
// TODO remove this and depend on Gloo function once it's bumped to k8s.io/* v0.18.6
func genCerts(svcName, svcNamespace string) (*certutils.Certificates, error) {
	return certutils.GenerateSelfSignedCertificate(cert.Config{
		CommonName:   fmt.Sprintf("%s.%s.svc", svcName, svcNamespace),
		Organization: []string{"solo.io"},
		AltNames: cert.AltNames{
			DNSNames: []string{
				svcName,
				fmt.Sprintf("%s.%s", svcName, svcNamespace),
				fmt.Sprintf("%s.%s.svc", svcName, svcNamespace),
				fmt.Sprintf("%s.%s.svc.%s", svcName, svcNamespace, network.GetClusterDomainName()),
			},
		},
		Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	})
}
