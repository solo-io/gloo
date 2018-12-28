package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"golang.org/x/sync/errgroup"

	"net"

	pb "github.com/solo-io/solo-projects/projects/licensingserver/pkg/api/v1"
	"google.golang.org/grpc"
)

const (
	service_name          = "licensing"
	health_check_endpoint = "/healthcheck"
)

type LicenseValidationServer struct {
	licensingClient LicensingClient
	health          *healthChecker
	grpcPort        int
	grpcServer      *grpc.Server
	ctx             context.Context
}

func (lvs *LicenseValidationServer) ValidateKey(ctx context.Context, in *pb.ValidateKeyRequest) (*pb.ValidateKeyReponse, error) {
	v := &pb.ValidateKeyReponse{Valid: false}
	valid, err := lvs.licensingClient.Validate(in.Key.Key)
	if err != nil {
		return nil, err
	}
	v.Valid = valid
	return v, nil
}

func (lvs *LicenseValidationServer) Start() error {
	var g errgroup.Group
	g.Go(func() error {
		return lvs.startGrpcServer()
	})
	g.Go(func() error {
		return lvs.health.start()
	})

	return g.Wait()
}

func (lvs *LicenseValidationServer) Close() {
	lvs.grpcServer.GracefulStop()
}

func NewServer(client LicensingClient, settings Settings, ctx context.Context) (*LicenseValidationServer, error) {
	if client == nil {
		return nil, fmt.Errorf("licensing client cannot be nil ")
	}

	if ctx == nil {
		ctx = contextutils.WithLogger(context.Background(), "licensing-server")
	}

	server := &LicenseValidationServer{
		grpcPort:        settings.GrpcPort,
		licensingClient: client,
		ctx:             ctx,
	}
	// setup ports
	server.setupGrpcServer()
	// setup healthcheck path
	server.health = newHealthChecker(server.grpcServer, settings.HealthPort, ctx)
	return server, nil
}

func (lvs *LicenseValidationServer) setupGrpcServer() {
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLicenseValidationServer(grpcServer, lvs)

	lvs.grpcServer = grpc.NewServer(opts...)

	pb.RegisterLicenseValidationServer(lvs.grpcServer, lvs)
}

func (lvs *LicenseValidationServer) startGrpcServer() error {
	addr := fmt.Sprintf(":%d", lvs.grpcPort)
	contextutils.LoggerFrom(lvs.ctx).Infof("Listening for gRPC on '%s'", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		contextutils.LoggerFrom(lvs.ctx).Fatalf("Failed to listen for gRPC: %v", err)
	}
	lvs.handleGracefulShutdown(lvs.ctx)
	return lvs.grpcServer.Serve(lis)
}

func (lvs *LicenseValidationServer) handleGracefulShutdown(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		select {
		case sig := <-sigs:
			contextutils.LoggerFrom(lvs.ctx).Warnf("Licensing server received %v, shutting down gracefully", sig)
		case <-ctx.Done():
		}

		lvs.grpcServer.GracefulStop()
		lvs.health.close()
	}()
}
