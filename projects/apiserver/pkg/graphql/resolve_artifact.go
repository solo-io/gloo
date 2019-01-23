package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

func (r namespaceResolver) Artifacts(ctx context.Context, obj *customtypes.Namespace) ([]*models.Artifact, error) {
	list, err := r.ArtifactClient.List(obj.Name, clients.ListOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputArtifacts(list), nil
}

func (r namespaceResolver) Artifact(ctx context.Context, obj *customtypes.Namespace, name string) (*models.Artifact, error) {
	artifact, err := r.ArtifactClient.Read(obj.Name, name, clients.ReadOpts{Ctx: ctx})
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
	out, err := r.ArtifactClient.Write(ups, clients.WriteOpts{
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
func (r *artifactMutationResolver) Delete(ctx context.Context, obj *customtypes.ArtifactMutation, guid string) (*models.Artifact, error) {
	_, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return &models.Artifact{}, err
	}
	artifact, err := r.ArtifactClient.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.ArtifactClient.Delete(namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputArtifact(artifact), nil
}
