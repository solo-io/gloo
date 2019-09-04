package configproto

import (
	server_pb_struct "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/rate-limiter/pkg/config"
)

type RateLimitConfigGenerator interface {
	GenerateConfig(configs []*server_pb_struct.RateLimitConfig) (config.RateLimitConfig, error)
}
