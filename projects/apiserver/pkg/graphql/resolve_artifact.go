package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
)

type artifactQueryResolver struct{ *ApiResolver }

func (r *artifactQueryResolver) List(ctx context.Context, obj *customtypes.ArtifactQuery, selector *models.InputMapStringString) ([]*models.Artifact, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.Artifacts.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputArtifacts(list), nil
}

func (r *artifactQueryResolver) Get(ctx context.Context, obj *customtypes.ArtifactQuery, name string) (*models.Artifact, error) {
	artifact, err := r.Artifacts.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputArtifact(artifact), nil
}

type artifactMutationResolver struct{ *ApiResolver }

func (r *artifactMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.ArtifactMutation, artifact models.InputArtifact) (*models.Artifact, error) {
	ups, err := NewConverter(r.ApiResolver, ctx).ConvertInputArtifact(artifact)
	if err != nil {
		return nil, err
	}
	out, err := r.Artifacts.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputArtifact(out), nil
}

func (r *artifactMutationResolver) Create(ctx context.Context, obj *customtypes.ArtifactMutation, artifact models.InputArtifact) (*models.Artifact, error) {
	return r.write(false, ctx, obj, artifact)
}
func (r *artifactMutationResolver) Update(ctx context.Context, obj *customtypes.ArtifactMutation, artifact models.InputArtifact) (*models.Artifact, error) {
	return r.write(true, ctx, obj, artifact)
}
func (r *artifactMutationResolver) Delete(ctx context.Context, obj *customtypes.ArtifactMutation, name string) (*models.Artifact, error) {
	artifact, err := r.Artifacts.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.Artifacts.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputArtifact(artifact), nil
}
