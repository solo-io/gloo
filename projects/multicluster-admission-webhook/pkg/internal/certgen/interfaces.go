package certgen

import (
	"context"
)

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_interfaces.go

/*
Interface for ensuring the CA Certificate for a kubernetes webhook.

It uses the svcName, svcNamespace to construct the DNS name used as the host for the CA cert.
The The cert info is then saved to the secret specified by secretName, secretNamespace.

This function can be run on startup, and then on a timer to ensure the certs are rotated properly.

Currently the 2 implementations available are:
  - Self Signed: creates a self-signed CA and saves it.
  - k8s Signed: Creates a CSR which should be signed by the apiserver, and then saves the result.
*/
type WebhookCAManager interface {
	EnsureCaCerts(
		ctx context.Context,
		secretName, secretNamespace, svcName, svcNamespace string,
	) error
}
