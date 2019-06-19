package server

import (
	"context"
	"net"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/setup"

	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service"
	"google.golang.org/grpc"
)

type GlooGrpcService struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGlooGrpcService(
	listener net.Listener,
	clientset setup.ClientSet) *GlooGrpcService {

	server := &GlooGrpcService{
		server:   grpc.NewServer(),
		listener: listener,
	}
	v1.RegisterUpstreamApiServer(server.server, service.NewUpstreamGrpcService(clientset.UpstreamClient))
	v1.RegisterArtifactApiServer(server.server, service.NewArtifactGrpcService(clientset.ArtifactClient))
	v1.RegisterConfigApiServer(server.server, service.NewConfigGrpcService(clientset.SettingsClient))
	v1.RegisterSecretApiServer(server.server, service.NewSecretGrpcService(clientset.SecretClient))
	v1.RegisterVirtualServiceApiServer(server.server, service.NewVirtualServiceGrpcService(clientset.VirtualServiceClient))
	return server
}

func (s *GlooGrpcService) Run(ctx context.Context) error {
	contextutils.LoggerFrom(ctx).Infow("Starting gloo grpc server")
	return s.server.Serve(s.listener)
}

func (s *GlooGrpcService) Stop() {
	s.server.Stop()
}
