package runner

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	DebugPort      int    `envconfig:"DEBUG_PORT" default:"9091"`
	PodNamespace   string `envconfig:"POD_NAMESPACE" default:"gloo-system"`
	PodName        string `envconfig:"POD_NAME" default:"gloo-fed"`
	WriteNamespace string `envconfig:"WRITE_NAMESPACE" default:"gloo-system"`
	LicenseKey     string `envconfig:"GLOO_LICENSE_KEY" default:""`
}

// NewSettings returns the Gloo Federation Settings based on the configured environment variables
func NewSettings() *Settings {
	var s Settings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	return &s
}
