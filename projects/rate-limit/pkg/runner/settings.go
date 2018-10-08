package runner

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	GlooAddress string `envconfig:"GLOO_ADDRESS" default:"control-plane:8080"` // This port serves simple health check responses
}

func NewSettings() Settings {
	var s Settings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	return s
}
