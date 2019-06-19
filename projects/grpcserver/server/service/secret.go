package service

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type secretGrpcService struct {
	secretClient gloov1.SecretClient
}

func NewSecretGrpcService(secretClient gloov1.SecretClient) v1.SecretApiServer {
	return &secretGrpcService{
		secretClient: secretClient,
	}
}

func (secretGrpcService) GetSecret(context.Context, *v1.GetSecretRequest) (*v1.GetSecretResponse, error) {
	panic("implement me")
}

func (secretGrpcService) ListSecrets(context.Context, *v1.ListSecretsRequest) (*v1.ListSecretsResponse, error) {
	panic("implement me")
}

func (secretGrpcService) CreateSecret(context.Context, *v1.CreateSecretRequest) (*v1.CreateSecretResponse, error) {
	panic("implement me")
}

func (secretGrpcService) UpdateSecret(context.Context, *v1.UpdateSecretRequest) (*v1.UpdateSecretResponse, error) {
	panic("implement me")
}

func (secretGrpcService) DeleteSecret(context.Context, *v1.DeleteSecretRequest) (*v1.DeleteSecretResponse, error) {
	panic("implement me")
}
