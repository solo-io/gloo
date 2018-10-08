package configproto

import (
	"github.com/solo-io/rate-limiter/pkg/config"
	server_pb_struct "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type RateLimitConfigGenerator interface {
	GenerateConfig(configs []*server_pb_struct.RateLimitConfig) (config.RateLimitConfig, error)
}
