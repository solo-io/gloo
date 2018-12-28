package server

import (
	"context"
)

type LicensingClient interface {
	Validate(key string) (bool, error)
}

type Settings struct {
	HealthPort int    `envconfig:"HTTP_PORT" default:"9001"` // This port serves simple health check responses
	GrpcPort   int    `envconfig:"GRPC_PORT" default:"9000"` // This port serves Licensing server requests
	LogLevel   string `envconfig:"LOG_LEVEL" default:"WARN"`
}

func Setup(ctx context.Context, settings Settings, debugMode bool, licensingClient LicensingClient) error {
	server, err := NewServer(licensingClient, settings, ctx)
	if err != nil {
		return err
	}
	err = server.Start()
	if err != nil {
		return err
	}
	return nil
}
