package runner

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	GlooAddress     string `envconfig:"GLOO_ADDRESS" default:"control-plane:8080"`
	SigningKey      string `envconfig:"SIGNING_KEY" default:""`
	TlsEnabled      bool   `envconfig:"TLS_ENABLED" default:"false"`
	Cert            []byte `envconfig:"CERT" default:""`
	Key             []byte `envconfig:"KEY" default:""`
	CertPath        string `envconfig:"CERT_PATH" default:"/etc/envoy/ssl/tls.crt"`
	KeyPath         string `envconfig:"KEY_PATH" default:"/etc/envoy/ssl/tls.key"`
	DebugPort       int    `envconfig:"DEBUG_PORT" default:"9091"`
	ServerPort      int    `envconfig:"SERVER_PORT" default:"8083"`
	ServiceName     string `envconfig:"SERVICE_NAME" default:"ext-auth"`
	ServerUDSAddr   string `envconfig:"UDS_ADDR" default:""`
	UserIdHeader    string `envconfig:"USER_ID_HEADER" default:""`
	PluginDirectory string `envconfig:"PLUGIN_DIRECTORY" default:"/auth-plugins/"`
}

func NewSettings() Settings {
	var s Settings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	return s
}
