package server

import (
	"context"
	"net"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"

	"github.com/solo-io/go-utils/contextutils"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"google.golang.org/grpc"
)

type GlooGrpcService struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGlooGrpcService(listener net.Listener, serviceSet setup.ServiceSet) *GlooGrpcService {
	server := &GlooGrpcService{
		server:   grpc.NewServer(),
		listener: listener,
	}

	v1.RegisterUpstreamApiServer(server.server, serviceSet.UpstreamService)
	v1.RegisterArtifactApiServer(server.server, serviceSet.ArtifactService)
	v1.RegisterConfigApiServer(server.server, serviceSet.ConfigService)
	v1.RegisterSecretApiServer(server.server, serviceSet.SecretService)
	v1.RegisterVirtualServiceApiServer(server.server, serviceSet.VirtualServiceService)
	return server
}

func (s *GlooGrpcService) Run(ctx context.Context) error {
	contextutils.LoggerFrom(ctx).Infow("Starting gloo grpc server")
	return s.server.Serve(s.listener)
}

func (s *GlooGrpcService) Stop() {
	s.server.Stop()
}
