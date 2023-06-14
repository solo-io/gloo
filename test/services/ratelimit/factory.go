package ratelimit

import (
	"sync/atomic"

	"github.com/solo-io/rate-limiter/pkg/cache/aerospike"
	"github.com/solo-io/rate-limiter/pkg/cache/dynamodb"
	"github.com/solo-io/rate-limiter/pkg/cache/redis"

	"github.com/solo-io/gloo/test/ginkgo/parallel"
	ratelimitserver "github.com/solo-io/rate-limiter/pkg/server"
)

const (
	basePort = uint32(18081)
)

type Factory struct {
	basePort uint32
}

func NewFactory() *Factory {
	return &Factory{
		basePort: basePort,
	}
}

func (f Factory) NewInstance(address string) *Instance {
	serverSettings := ratelimitserver.NewSettings()
	serverSettings.HealthFailTimeout = 2 // seconds
	serverSettings.RateLimitPort = int(advancePort(&f.basePort))
	serverSettings.ReadyPort = int(advancePort(&f.basePort))
	serverSettings.LogSettings = ratelimitserver.LogSettings{
		LoggerName: "rate-limit-service",
		DebugMode:  "true",
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
