package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
)

type schemaQueryResolver struct{ *ApiResolver }

func (r *schemaQueryResolver) List(ctx context.Context, obj *customtypes.SchemaQuery, selector *models.InputMapStringString) ([]*models.Schema, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.Schemas.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSchemas(list), nil
}

func (r *schemaQueryResolver) Get(ctx context.Context, obj *customtypes.SchemaQuery, name string) (*models.Schema, error) {
	schema, err := r.Schemas.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
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
	out, err := r.Schemas.Write(ups, clients.WriteOpts{
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
func (r *schemaMutationResolver) Delete(ctx context.Context, obj *customtypes.SchemaMutation, name string) (*models.Schema, error) {
	schema, err := r.Schemas.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.Schemas.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSchema(schema), nil
}
