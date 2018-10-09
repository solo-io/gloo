package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
)

type secretQueryResolver struct{ *ApiResolver }

func (r *secretQueryResolver) List(ctx context.Context, obj *customtypes.SecretQuery, selector *models.InputMapStringString) ([]*models.Secret, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.Secrets.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSecrets(list), nil
}

func (r *secretQueryResolver) Get(ctx context.Context, obj *customtypes.SecretQuery, name string) (*models.Secret, error) {
	secret, err := r.Secrets.Read(obj.Namespace, name, clients.ReadOpts{
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
	out, err := r.Secrets.Write(ups, clients.WriteOpts{
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
func (r *secretMutationResolver) Delete(ctx context.Context, obj *customtypes.SecretMutation, name string) (*models.Secret, error) {
	secret, err := r.Secrets.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.Secrets.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSecret(secret), nil
}
