package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

func (r namespaceResolver) Secrets(ctx context.Context, obj *customtypes.Namespace) ([]*models.Secret, error) {
	list, err := r.SecretClient.List(obj.Name, clients.ListOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSecrets(list), nil
}

func (r namespaceResolver) Secret(ctx context.Context, obj *customtypes.Namespace, name string) (*models.Secret, error) {
	secret, err := r.SecretClient.Read(obj.Name, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSecret(secret), nil
}

type secretMutationResolver struct{ *ApiResolver }

func (r *secretMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.SecretMutation, secret models.InputSecret) (*models.Secret, error) {
	ups, err := NewConverter(r.ApiResolver, ctx).ConvertInputSecret(secret)
	if err != nil {
		return nil, err
	}
	out, err := r.SecretClient.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSecret(out), nil
}

func (r *secretMutationResolver) Create(ctx context.Context, obj *customtypes.SecretMutation, secret models.InputSecret) (*models.Secret, error) {
	return r.write(false, ctx, obj, secret)
}
func (r *secretMutationResolver) Update(ctx context.Context, obj *customtypes.SecretMutation, secret models.InputSecret) (*models.Secret, error) {
	return r.write(true, ctx, obj, secret)
}
func (r *secretMutationResolver) Delete(ctx context.Context, obj *customtypes.SecretMutation, guid string) (*models.Secret, error) {
	_, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return &models.Secret{}, err
	}
	secret, err := r.SecretClient.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.SecretClient.Delete(namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSecret(secret), nil
}
