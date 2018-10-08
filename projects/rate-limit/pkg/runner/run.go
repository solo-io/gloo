package runner

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	pb "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v2"
	"github.com/solo-io/rate-limiter/pkg/redis"
	"github.com/solo-io/rate-limiter/pkg/server"
	"github.com/solo-io/rate-limiter/pkg/service"
	"github.com/solo-io/rate-limiter/pkg/settings"
	configproto "github.com/solo-io/solo-kit/projects/rate-limit/pkg/config"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"

	"github.com/solo-io/solo-kit/pkg/utils/stats"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"

	"google.golang.org/grpc"
)

func Run() {
	s := settings.NewSettings()
	clientSettings := NewSettings()
	var perSecondPool redis.Pool
	if s.RedisPerSecond {
		perSecondPool = redis.NewPoolImpl(s.RedisPerSecondSocketType, s.RedisPerSecondUrl, s.RedisPerSecondPoolSize)
	}
	redisPool := redis.NewPoolImpl(s.RedisSocketType, s.RedisUrl, s.RedisPoolSize)

	service := NewService(s, redisPool, perSecondPool)

	debugPort := fmt.Sprintf("%d", s.DebugPort)
	// TODO(yuval-k): we need to start the stats server before calling contextutils
	// need to think of a better way to express this dependency, or preferably, fix it.
	stats.StartStatsServerWithPort(debugPort, addConfigDumpHandler(service))

	ctx := context.Background()
	ctx = contextutils.WithLogger(ctx, "ratelimit")

	StartRateLimit(ctx, s, clientSettings, service)
}

func NewService(s settings.Settings, redisPool, perSecondPool redis.Pool) ratelimit.RateLimitServiceServer {
	return ratelimit.NewService(
		redis.NewRateLimitCacheImpl(
			redisPool,
			perSecondPool,
			redis.NewTimeSourceImpl(),
			rand.New(redis.NewLockedSource(time.Now().Unix())),
			s.ExpirationJitterMaxSeconds),
	)
}

func StartRateLimit(ctx context.Context, s settings.Settings, clientSettings Settings, service ratelimit.RateLimitServiceServer) {
	srv := server.NewServer("ratelimit", s)
	StartRateLimitWithGrpcServer(ctx, clientSettings, service, srv.GrpcServer())
	srv.Start()
}

func StartRateLimitWithGrpcServer(ctx context.Context, clientSettings Settings, service ratelimit.RateLimitServiceServer, grpcServer *grpc.Server) {
	err := startClient(ctx, clientSettings, service)
	if err != nil {
		panic(err)
	}

	pb.RegisterRateLimitServiceServer(grpcServer, service)
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

		// We are using non secure grpc to gloo with the asumption that it will be
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

func addConfigDumpHandler(service ratelimit.RateLimitServiceServer) func(mux *http.ServeMux, profiles map[string]string) {
	return func(mux *http.ServeMux, profiles map[string]string) {

		mux.HandleFunc(
			"/rlconfig",
			func(writer http.ResponseWriter, request *http.Request) {
				io.WriteString(writer, service.GetCurrentConfig().Dump())
			})

		profiles["/rlconfig"] = "print out the currently loaded configuration for debugging"
	}
}
