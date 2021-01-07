package settings

import (
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	DebugPort      int    `envconfig:"DEBUG_PORT" default:"9091"`
	PodNamespace   string `envconfig:"POD_NAMESPACE" default:"gloo-fed"`
	PodName        string `envconfig:"POD_NAME" default:"gloo-fed"`
	WriteNamespace string `envconfig:"WRITE_NAMESPACE" default:"gloo-fed"`
	LicenseKey     string `envconfig:"LICENSE_KEY" default:""`
	ApiServerSettings
}

type ApiServerSettings struct {
	GrpcPort        int `envconfig:"APISERVER_GRPC_PORT" default:"10101"`
	HealthCheckPort int `envconfig:"APISERVER_HEALTH_CHECK_PORT" default:"8081"`
}

func New() *Settings {
	var s Settings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	return &s
}
