package pkgmetrics // import "go.pedge.io/pkg/metrics"

import (
	"time"

	"github.com/mihasya/go-metrics-librato"
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/stathat"
)

// Env defines a struct for environment variables that can be parsed with go.pedge.io/env.
type Env struct {
	// The email address for the Librato account to send stats to.
	// Must be set with LibratoAPIToken.
	// If not set and StathatUserKey not set, no metrics.Registry for stats will be created.
	LibratoEmailAddress string `env:"LIBRATO_EMAIL_ADDRESS"`
	// The API Token for the Librato account to send stats to.
	// Must be set with LibratoEmailAddress.
	// If not set and StathatUserKey not set, no metrics.Registry for stats will be created.
	LibratoAPIToken string `env:"LIBRATO_API_TOKEN"`
	// The StatHat user key to send stats to.
	// If not set and LibratoEmailAddress and LibratoAPIToken not set, no metrics.Registry for stats will be created.
	StathatUserKey string `env:"STATHAT_USER_KEY"`
}

// SetupMetrics sets up metrics.
//
// The returned metrics.Registry can be nil.
func SetupMetrics(appName string, env Env) (metrics.Registry, error) {
	if env.StathatUserKey == "" && env.LibratoEmailAddress == "" && env.LibratoAPIToken == "" {
		return nil, nil
	}
	registry := metrics.NewPrefixedRegistry(appName)
	if env.StathatUserKey != "" {
		go stathat.Stathat(
			registry,
			time.Hour,
			env.StathatUserKey,
		)
	}
	if env.LibratoEmailAddress != "" && env.LibratoAPIToken != "" {
		go librato.Librato(
			registry,
			5*time.Minute,
			env.LibratoEmailAddress,
			env.LibratoAPIToken,
			appName,
			[]float64{0.95},
			time.Millisecond,
		)
	}
	return registry, nil
}
