package runner

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/solo-io/ext-auth-service/pkg/server"
)

type Settings struct {
	// This port serves simple health check responses
	GlooAddress         string        `envconfig:"GLOO_ADDRESS" default:"control-plane:8080"`
	TlsEnabled          bool          `envconfig:"TLS_ENABLED" default:"false"`
	Cert                []byte        `envconfig:"CERT" default:""`
	Key                 []byte        `envconfig:"KEY" default:""`
	CertPath            string        `envconfig:"CERT_PATH" default:"/etc/envoy/ssl/tls.crt"`
	KeyPath             string        `envconfig:"KEY_PATH" default:"/etc/envoy/ssl/tls.key"`
	ServerUDSAddr       string        `envconfig:"UDS_ADDR" default:""`
	HeadersToRedact     string        `envconfig:"HEADERS_TO_REDACT" default:"authorization"`
	ShutdownGracePeriod time.Duration `envconfig:"SHUTDOWN_GRACE_PERIOD" default:"15s"`

	ExtAuthSettings server.Settings
}

func NewSettings() Settings {
	var s Settings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	s.ExtAuthSettings = server.NewSettings()

	return s
}
