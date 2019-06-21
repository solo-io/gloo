package artifactsvc

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type artifactGrpcService struct {
	artifactClient gloov1.ArtifactClient
}

func NewArtifactGrpcService(artifactClient gloov1.ArtifactClient) v1.ArtifactApiServer {
	return &artifactGrpcService{
		artifactClient: artifactClient,
	}
}

func (s *artifactGrpcService) GetArtifact(context.Context, *v1.GetArtifactRequest) (*v1.GetArtifactResponse, error) {
	panic("implement me")
}

func (s *artifactGrpcService) ListArtifacts(context.Context, *v1.ListArtifactsRequest) (*v1.ListArtifactsResponse, error) {
	panic("implement me")
}

func (s *artifactGrpcService) CreateArtifact(context.Context, *v1.CreateArtifactRequest) (*v1.CreateArtifactResponse, error) {
	panic("implement me")
}

func (s *artifactGrpcService) UpdateArtifact(context.Context, *v1.UpdateArtifactRequest) (*v1.UpdateArtifactResponse, error) {
	panic("implement me")
}

func (s *artifactGrpcService) DeleteArtifact(context.Context, *v1.DeleteArtifactRequest) (*v1.DeleteArtifactResponse, error) {
	panic("implement me")
}
