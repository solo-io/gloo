package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

type resolverMapQueryResolver struct{ *ApiResolver }

func (r *resolverMapQueryResolver) List(ctx context.Context, obj *customtypes.ResolverMapQuery, selector *models.InputMapStringString) ([]*models.ResolverMap, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.ResolverMaps.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputResolverMaps(list)
}

func (r *resolverMapQueryResolver) Get(ctx context.Context, obj *customtypes.ResolverMapQuery, name string) (*models.ResolverMap, error) {
	resolverMap, err := r.ResolverMaps.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputResolverMap(resolverMap)
}

type resolverMapMutationResolver struct{ *ApiResolver }

func (r *resolverMapMutationResolver) SetResolver(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMapName, resourceVersion, typeName, fieldName string, resolver models.InputGlooResolver) (*models.ResolverMap, error) {
	v1Resolver, err := ConvertInputResolver(models.InputResolver{GlooResolver: &resolver})
	if err != nil {
		return nil, err
	}

	resolverMap, err := r.ResolverMaps.Read(obj.Namespace, resolverMapName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if resolverMap.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, resolverMap.Metadata.ResourceVersion)
	}

	typResolver, ok := resolverMap.Types[typeName]
	if !ok {
		return nil, errors.Errorf("no type %v in resolver map %v", typeName, resolverMapName)
	}
	typResolver.Fields[fieldName] = v1Resolver

	out, err := r.ResolverMaps.Write(resolverMap, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputResolverMap(out)
}

func (r *resolverMapMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	ups, err := NewConverter(r.ApiResolver, ctx).ConvertInputResolverMap(resolverMap)
	if err != nil {
		return nil, err
	}
	out, err := r.ResolverMaps.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputResolverMap(out)
}

func (r *resolverMapMutationResolver) Create(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	return r.write(false, ctx, obj, resolverMap)
}
func (r *resolverMapMutationResolver) Update(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	return r.write(true, ctx, obj, resolverMap)
}
func (r *resolverMapMutationResolver) Delete(ctx context.Context, obj *customtypes.ResolverMapMutation, name string) (*models.ResolverMap, error) {
	resolverMap, err := r.ResolverMaps.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.ResolverMaps.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputResolverMap(resolverMap)
}
