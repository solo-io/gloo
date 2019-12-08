package runner

import (
	"context"
	"os"

	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	_struct "github.com/golang/protobuf/ptypes/struct"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/rate-limiter/cmd/service/runner"
	ratelimit "github.com/solo-io/rate-limiter/pkg/service"
	configproto "github.com/solo-io/solo-projects/projects/rate-limit/pkg/config"
	"google.golang.org/grpc"
)

func Run() {
	clientSettings := NewSettings()
	RunWithClientSettingsBlocking(clientSettings)
}

func RunWithClientSettingsBlocking(clientSettings Settings) {
	ctx, cancel := context.WithCancel(context.Background())
	RunWithClientSettings(ctx, cancel, clientSettings)
	<-ctx.Done()
}

func RunWithClientSettings(ctx context.Context, cancel context.CancelFunc, clientSettings Settings) ratelimit.RateLimitServiceServer {
	svc := runner.Run(cancel, ctx)
	err := startClient(ctx, clientSettings, svc)
	if err != nil {
		panic(err)
	}
	return svc
}

func startClient(ctx context.Context, s Settings, service ratelimit.RateLimitServerConfigMutator) error {
	var nodeinfo core.Node
	var err error
	nodeinfo.Id, err = os.Hostname()
	// TODO(yuval-k): unhardcode this
	if err != nil {
		nodeinfo.Id = "ratelimit-unknown"
	}
	nodeinfo.Cluster = "ratelimit"
	role := "ratelimit"
	nodeinfo.Metadata = &_struct.Struct{
		Fields: map[string]*_struct.Value{
			"role": {
				Kind: &_struct.Value_StringValue{
					StringValue: role,
				},
			},
		},
	}

	go clientLoop(ctx, s.GlooAddress, nodeinfo, service)
	return nil
}

func clientLoop(ctx context.Context, dialString string, nodeinfo core.Node, service ratelimit.RateLimitServerConfigMutator) {
	generator := configproto.NewConfigGenerator(contextutils.LoggerFrom(ctx))

	contextutils.NewExponentioalBackoff(contextutils.ExponentioalBackoff{}).Backoff(ctx, func(ctx context.Context) error {
		client := v1.NewRateLimitConfigClient(&nodeinfo, func(version string, resources []*v1.RateLimitConfig) error {
			config, err := generator.GenerateConfig(resources)
			if err != nil {
				return err
			}
			service.SetCurrentConfig(config)
			return nil
		})

		// We are using non secure grpc to gloo with the assumption that it will be
		// secured by envoy. if this assumption is not correct this needs to change.
		conn, err := grpc.DialContext(ctx, dialString, grpc.WithInsecure())
		if err != nil {
			return err
		}
		// TODO(yuval-k): a stat that indicates we are connected, with the reverse one deferred.
		// TODO(yuval-k): write a warning log
		return client.Start(ctx, conn)
	})
}
