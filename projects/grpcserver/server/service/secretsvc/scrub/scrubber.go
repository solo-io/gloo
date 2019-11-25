package scrub

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

//go:generate mockgen -destination mocks/mock_scrubber.go -package mocks github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc/scrub Scrubber

type Scrubber interface {
	Secret(ctx context.Context, secret *gloov1.Secret)
}

type scrubber struct{}

var _ Scrubber = scrubber{}

func NewScrubber() Scrubber {
	return scrubber{}
}

// Modify the given secret to have the correct type but no value.
func (s scrubber) Secret(ctx context.Context, secret *gloov1.Secret) {
	switch secret.GetKind().(type) {
	case *gloov1.Secret_Aws:
		secret.Kind = &gloov1.Secret_Aws{Aws: &gloov1.AwsSecret{}}
	case *gloov1.Secret_Azure:
		secret.Kind = &gloov1.Secret_Azure{Azure: &gloov1.AzureSecret{}}
	case *gloov1.Secret_Tls:
		secret.Kind = &gloov1.Secret_Tls{Tls: &gloov1.TlsSecret{}}
	default:
		contextutils.LoggerFrom(ctx).Warnw(
			"Scrubbing secret with unhandled type. Details will be removed but type will be undefined.",
			zap.String("namespace", secret.GetMetadata().Namespace),
			zap.String("names", secret.GetMetadata().Name))
		secret.Kind = nil
	}
}
