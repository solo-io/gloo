package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

func (r namespaceResolver) Upstreams(ctx context.Context, obj *customtypes.Namespace) ([]*models.Upstream, error) {
	list, err := r.UpstreamClient.List(obj.Name, clients.ListOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputUpstreams(list)
}

func (r namespaceResolver) Upstream(ctx context.Context, obj *customtypes.Namespace, name string) (*models.Upstream, error) {
	upstream, err := r.UpstreamClient.Read(obj.Name, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputUpstream(upstream)
}

type upstreamMutationResolver struct{ *ApiResolver }

func (r *upstreamMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	ups, err := NewConverter(r.ApiResolver, ctx).ConvertInputUpstream(upstream)
	if err != nil {
		return nil, err
	}
	out, err := r.UpstreamClient.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputUpstream(out)
}

func (r *upstreamMutationResolver) Create(ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	return r.write(false, ctx, obj, upstream)
}
func (r *upstreamMutationResolver) Update(ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	return r.write(true, ctx, obj, upstream)
}
func (r *upstreamMutationResolver) Delete(ctx context.Context, obj *customtypes.UpstreamMutation, guid string) (*models.Upstream, error) {
	_, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return &models.Upstream{}, err
	}
	upstream, err := r.UpstreamClient.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.UpstreamClient.Delete(namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputUpstream(upstream)
}
