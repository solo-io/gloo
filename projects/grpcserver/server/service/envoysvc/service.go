package envoysvc

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails"
	"go.uber.org/zap"
)

type envoyGrpcService struct {
	ctx                context.Context
	envoyDetailsClient envoydetails.Client
	podNamespace       string
}

func NewEnvoyGrpcService(
	ctx context.Context,
	client envoydetails.Client,
	podNamespace string) v1.EnvoyApiServer {

	return &envoyGrpcService{
		ctx:                ctx,
		envoyDetailsClient: client,
		podNamespace:       podNamespace,
	}
}

func (s *envoyGrpcService) ListEnvoyDetails(ctx context.Context, request *v1.ListEnvoyDetailsRequest) (*v1.ListEnvoyDetailsResponse, error) {
	detailsList, err := s.envoyDetailsClient.List(s.ctx, s.podNamespace)
	if err != nil {
		wrapped := FailedToListEnvoyDetailsError(err)
		contextutils.LoggerFrom(s.ctx).Errorw(wrapped.Error(), zap.Error(err))
		return nil, wrapped
	}

	return &v1.ListEnvoyDetailsResponse{EnvoyDetails: detailsList}, nil
}
