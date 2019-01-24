package runner

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	GlooAddress  string `envconfig:"GLOO_ADDRESS" default:"control-plane:8080"`
	SigningKey   string `envconfig:"SIGNING_KEY_ADDRESS" default:""`
	DebugPort    int    `envconfig:"DEBUG_PORT" default:9091`
	ServerPort   int    `envconfig:"SERVER_PORT" default:8080`
	UserIdHeader string `envconfig:"USER_ID_HEADER" default:""`
}

func NewSettings() Settings {
	var s Settings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	return s
}
