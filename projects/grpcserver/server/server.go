package server

import (
	"context"
	"net"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
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
	virtualServiceClient gatewayv1.VirtualServiceClient,
	upstreamClient gloov1.UpstreamClient,
	artifactClient gloov1.ArtifactClient,
	secretClient gloov1.SecretClient,
	settingsClient gloov1.SettingsClient) *GlooGrpcService {

	server := &GlooGrpcService{
		server:   grpc.NewServer(),
		listener: listener,
	}
	v1.RegisterUpstreamApiServer(server.server, service.NewUpstreamGrpcService(upstreamClient))
	v1.RegisterArtifactApiServer(server.server, service.NewArtifactGrpcService(artifactClient))
	v1.RegisterConfigApiServer(server.server, service.NewConfigGrpcService(settingsClient))
	v1.RegisterSecretApiServer(server.server, service.NewSecretGrpcService(secretClient))
	v1.RegisterVirtualServiceApiServer(server.server, service.NewVirtualServiceGrpcService(virtualServiceClient))
	return server
}

func (s *GlooGrpcService) Run(ctx context.Context) error {
	contextutils.LoggerFrom(ctx).Infow("Starting gloo grpc server")
	return s.server.Serve(s.listener)
}

func (s *GlooGrpcService) Stop() {
	s.server.Stop()
}
