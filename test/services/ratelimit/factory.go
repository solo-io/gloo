package ratelimit

import (
	"sync/atomic"

	"github.com/solo-io/gloo/test/services"

	"github.com/solo-io/rate-limiter/pkg/cache/aerospike"
	"github.com/solo-io/rate-limiter/pkg/cache/dynamodb"
	"github.com/solo-io/rate-limiter/pkg/cache/redis"

	"github.com/solo-io/gloo/test/ginkgo/parallel"
	ratelimitserver "github.com/solo-io/rate-limiter/pkg/server"
)

const (
	// basePort is the starting port for the rate limit server
	// This was the previous static port used in tests, but it is not a special value
	basePort = uint32(18081)

	ServiceName = "rate-limit-service"
)

type Factory struct {
	basePort uint32
}

func NewFactory() *Factory {
	return &Factory{
		basePort: basePort,
	}
}

func (f *Factory) NewInstance(address string) *Instance {
	serverSettings := ratelimitserver.NewSettings()
	// The number of seconds that the server will remain alive, but actively failing health checks
	// After this time elapses, the server will exit
	// This is useful for testing that the server will exit when it fails health checks
	// This is useful in production to ensure that the server will handle in flight requests before exiting
	serverSettings.HealthFailTimeout = 0 // seconds
	serverSettings.RateLimitPort = int(advancePort(&f.basePort))
	serverSettings.ReadyPort = int(advancePort(&f.basePort))
	serverSettings.LogSettings = ratelimitserver.LogSettings{
		LoggerName: ServiceName,
		LogLevel:   services.GetLogLevel(ServiceName).String(),
	}
	serverSettings.RedisSettings = redis.Settings{}
	serverSettings.DynamoDbSettings = dynamodb.Settings{}
	serverSettings.AerospikeSettings = aerospike.Settings{}

	return &Instance{
		address:        address,
		serverSettings: serverSettings,
	}
}

func advancePort(p *uint32) uint32 {
	return atomic.AddUint32(p, 1) + uint32(parallel.GetPortOffset())
}
