package settings

import (
	"github.com/kelseyhightower/envconfig"
)

type ApiServerSettings struct {
	LicenseKey      string `envconfig:"LICENSE_KEY" default:""`
	GrpcPort        int    `envconfig:"APISERVER_GRPC_PORT" default:"10101"`
	HealthCheckPort int    `envconfig:"APISERVER_HEALTH_CHECK_PORT" default:"8081"`
}

func New() *ApiServerSettings {
	var s ApiServerSettings

	err := envconfig.Process("", &s)
	if err != nil {
		panic(err)
	}

	return &s
}
