package runner

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2alpha"
	extauthconfig "github.com/solo-io/ext-auth-service/pkg/config"
	extauth "github.com/solo-io/ext-auth-service/pkg/service"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	configproto "github.com/solo-io/solo-projects/projects/extauth/pkg/config"
	xdsproto "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	plugin "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"github.com/solo-io/solo-kit/pkg/utils/stats"

	"google.golang.org/grpc"
)

func Run() {
	clientSettings := NewSettings()
	ctx := context.Background()

	if clientSettings.DebugPort != 0 {

		debugPort := fmt.Sprintf("%d", clientSettings.DebugPort)
		// TODO(yuval-k): we need to start the stats server before calling contextutils
		// need to think of a better way to express this dependency, or preferably, fix it.
		stats.StartStatsServerWithPort(debugPort)
	}

	err := RunWithSettings(ctx, clientSettings)

	if err != nil {
		if ctx.Err() == nil {
			// not a context error - panic
			panic(err)
		}
	}
}
func RunWithSettings(ctx context.Context, clientSettings Settings) error {

	service := extauth.NewServer()
	service.VhostContextExtension = plugin.ContextExtensionVhost

	ctx = contextutils.WithLogger(ctx, "extauth")

	err := StartExtAuth(ctx, clientSettings, service)
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

func StartExtAuth(ctx context.Context, clientSettings Settings, service *extauth.Server) error {
	srv := grpc.NewServer(grpc.UnaryInterceptor(nil))

	pb.RegisterAuthorizationServer(srv, service)
	healthpb.RegisterHealthServer(srv, &healthChecker{})
	reflection.Register(srv)

	logger := contextutils.LoggerFrom(ctx)

	err := StartExtAuthWithGrpcServer(ctx, clientSettings, service)
	if err != nil {
		logger.Error("Failed to starting ext auth server: %v", err)
		return err
	}

	addr := fmt.Sprintf(":%d", clientSettings.ServerPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error("Failed to listen for gRPC: %v", err)
		return err
	}
	go func() {
		<-ctx.Done()
		srv.Stop()
		lis.Close()
	}()

	return srv.Serve(lis)
}

func StartExtAuthWithGrpcServer(ctx context.Context, clientSettings Settings, service extauthconfig.ConfigMutator) error {
	var nodeinfo core.Node
	var err error
	nodeinfo.Id, err = os.Hostname()
	// TODO(yuval-k): unhardcode this
	if err != nil {
		nodeinfo.Id = "extauth-unknown"
	}
	nodeinfo.Cluster = "extauth"
	role := "extauth"
	nodeinfo.Metadata = &types.Struct{
		Fields: map[string]*types.Value{
			"role": &types.Value{
				Kind: &types.Value_StringValue{
					StringValue: role,
				},
			},
		},
	}

	go clientLoop(ctx, clientSettings, nodeinfo, service)
	return nil
}

func clientLoop(ctx context.Context, clientSettings Settings, nodeinfo core.Node, service extauthconfig.ConfigMutator) {
	generator := configproto.NewConfigGenerator(ctx, []byte(clientSettings.SigningKey), clientSettings.UserIdHeader)

	contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{}).Backoff(ctx, func(ctx context.Context) error {
		client := xdsproto.NewExtAuthConfigClient(&nodeinfo, func(version string, resources []*xdsproto.ExtAuthConfig) error {
			config, err := generator.GenerateConfig(resources)
			if err != nil {
				return err
			}
			service.UpdateConfig(config)
			return nil
		})

		// We are using non secure grpc to gloo with the asumption that it will be
		// secured by envoy. if this assumption is not correct this needs to change.
		conn, err := grpc.DialContext(ctx, clientSettings.GlooAddress, grpc.WithInsecure())
		if err != nil {
			return err
		}
		// TODO(yuval-k): a stat that indicates we are connected, with the reverse one deferred.
		// TODO(yuval-k): write a warning log
		return client.Start(ctx, conn)
	})
}

type healthChecker struct{}

func (h *healthChecker) Check(context.Context, *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
