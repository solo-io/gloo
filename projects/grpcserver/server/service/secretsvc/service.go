package secretsvc

import (
	"context"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"

	"github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/svccodes"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc/scrub"
	"go.uber.org/zap"
)

type secretGrpcService struct {
	ctx            context.Context
	clientCache    client.ClientCache
	secretScrubber scrub.Scrubber
	licenseClient  license.Client
}

func NewSecretGrpcService(
	ctx context.Context,
	clientCache client.ClientCache,
	secretScrubber scrub.Scrubber,
	licenseClient license.Client,
) v1.SecretApiServer {
	return &secretGrpcService{
		ctx:            ctx,
		clientCache:    clientCache,
		secretScrubber: secretScrubber,
		licenseClient:  licenseClient,
	}
}

func (s *secretGrpcService) GetSecret(ctx context.Context, request *v1.GetSecretRequest) (*v1.GetSecretResponse, error) {
	secret, err := s.readSecret(request.GetRef())
	if err != nil {
		wrapped := FailedToReadSecretError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	s.secretScrubber.Secret(s.ctx, secret)
	return &v1.GetSecretResponse{Secret: secret}, nil
}

func (s *secretGrpcService) ListSecrets(ctx context.Context, request *v1.ListSecretsRequest) (*v1.ListSecretsResponse, error) {
	var secretList gloov1.SecretList
	for _, ns := range request.GetNamespaces() {
		secrets, err := s.clientCache.GetSecretClient().List(ns, clients.ListOpts{Ctx: s.ctx})
		if err != nil {
			wrapped := FailedToListSecretsError(err, ns)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		secretList = append(secretList, secrets...)
	}

	for _, secret := range secretList {
		s.secretScrubber.Secret(s.ctx, secret)
	}

	return &v1.ListSecretsResponse{Secrets: secretList}, nil
}

func (s *secretGrpcService) CreateSecret(ctx context.Context, request *v1.CreateSecretRequest) (*v1.CreateSecretResponse, error) {
	// TODO(mitchdraft) move this (and all other calls to CheckLicenseForGlooUiMutations) to solo-kit as a client hook
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	var (
		secret *gloov1.Secret
		ref    *core.ResourceRef
	)

	if request.GetSecret() == nil {
		secret = &gloov1.Secret{
			Metadata: core.Metadata{
				Namespace: request.GetRef().GetNamespace(),
				Name:      request.GetRef().GetName(),
			},
		}

		switch request.GetKind().(type) {
		case *v1.CreateSecretRequest_Aws:
			secret.Kind = &gloov1.Secret_Aws{Aws: request.GetAws()}
		case *v1.CreateSecretRequest_Azure:
			secret.Kind = &gloov1.Secret_Azure{Azure: request.GetAzure()}
		case *v1.CreateSecretRequest_Extension:
			secret.Kind = &gloov1.Secret_Extension{Extension: request.GetExtension()}
		case *v1.CreateSecretRequest_Tls:
			secret.Kind = &gloov1.Secret_Tls{Tls: request.GetTls()}
		}

		ref = request.GetRef()
	} else {
		secret = request.GetSecret()
		secretRef := request.GetSecret().GetMetadata().Ref()
		ref = &secretRef
	}

	written, err := s.writeSecret(*secret, false)
	if err != nil {
		wrapped := FailedToCreateSecretError(err, ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	s.secretScrubber.Secret(s.ctx, written)
	return &v1.CreateSecretResponse{Secret: written}, nil
}

func (s *secretGrpcService) UpdateSecret(ctx context.Context, request *v1.UpdateSecretRequest) (*v1.UpdateSecretResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	var (
		ref           *core.ResourceRef
		secretToWrite *gloov1.Secret
	)

	if request.GetSecret() == nil {
		ref = request.GetRef()

		read, err := s.readSecret(ref)
		if err != nil {
			wrapped := FailedToUpdateSecretError(err, ref)
			contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}

		switch request.GetKind().(type) {
		case *v1.UpdateSecretRequest_Aws:
			read.Kind = &gloov1.Secret_Aws{Aws: request.GetAws()}
		case *v1.UpdateSecretRequest_Azure:
			read.Kind = &gloov1.Secret_Azure{Azure: request.GetAzure()}
		case *v1.UpdateSecretRequest_Extension:
			read.Kind = &gloov1.Secret_Extension{Extension: request.GetExtension()}
		case *v1.UpdateSecretRequest_Tls:
			read.Kind = &gloov1.Secret_Tls{Tls: request.GetTls()}
		}
		secretToWrite = read
	} else {
		metadataRef := request.GetSecret().GetMetadata().Ref()
		ref = &metadataRef
		secretToWrite = request.GetSecret()
	}

	written, err := s.writeSecret(*secretToWrite, true)
	if err != nil {
		wrapped := FailedToUpdateSecretError(err, ref)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	s.secretScrubber.Secret(s.ctx, written)
	return &v1.UpdateSecretResponse{Secret: written}, nil
}

func (s *secretGrpcService) DeleteSecret(ctx context.Context, request *v1.DeleteSecretRequest) (*v1.DeleteSecretResponse, error) {
	if err := svccodes.CheckLicenseForGlooUiMutations(ctx, s.licenseClient); err != nil {
		return nil, err
	}
	err := s.clientCache.GetSecretClient().Delete(request.GetRef().GetNamespace(), request.GetRef().GetName(), clients.DeleteOpts{Ctx: s.ctx})
	if err != nil {
		wrapped := FailedToDeleteSecretError(err, request.GetRef())
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &v1.DeleteSecretResponse{}, nil
}

func (s *secretGrpcService) readSecret(ref *core.ResourceRef) (*gloov1.Secret, error) {
	return s.clientCache.GetSecretClient().Read(ref.GetNamespace(), ref.GetName(), clients.ReadOpts{Ctx: s.ctx})
}

func (s *secretGrpcService) writeSecret(secret gloov1.Secret, shouldOverWrite bool) (*gloov1.Secret, error) {
	return s.clientCache.GetSecretClient().Write(&secret, clients.WriteOpts{Ctx: s.ctx, OverwriteExisting: shouldOverWrite})
}
