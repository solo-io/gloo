package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

func (r namespaceResolver) Schemas(ctx context.Context, obj *customtypes.Namespace) ([]*models.Schema, error) {
	list, err := r.SchemaClient.List(obj.Name, clients.ListOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSchemas(list), nil
}

func (r namespaceResolver) Schema(ctx context.Context, obj *customtypes.Namespace, name string) (*models.Schema, error) {
	schema, err := r.SchemaClient.Read(obj.Name, name, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSchema(schema), nil
}

type schemaMutationResolver struct{ *ApiResolver }

func (r *schemaMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.SchemaMutation, schema models.InputSchema) (*models.Schema, error) {
	ups, err := NewConverter(r.ApiResolver, ctx).ConvertInputSchema(schema)
	if err != nil {
		return nil, err
	}
	out, err := r.SchemaClient.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSchema(out), nil
}

func (r *schemaMutationResolver) Create(ctx context.Context, obj *customtypes.SchemaMutation, schema models.InputSchema) (*models.Schema, error) {
	return r.write(false, ctx, obj, schema)
}
func (r *schemaMutationResolver) Update(ctx context.Context, obj *customtypes.SchemaMutation, schema models.InputSchema) (*models.Schema, error) {
	return r.write(true, ctx, obj, schema)
}
func (r *schemaMutationResolver) Delete(ctx context.Context, obj *customtypes.SchemaMutation, guid string) (*models.Schema, error) {
	_, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return &models.Schema{}, err
	}
	schema, err := r.SchemaClient.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.SchemaClient.Delete(namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSchema(schema), nil
}
