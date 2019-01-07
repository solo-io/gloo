package main

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"

	"github.com/solo-io/solo-projects/projects/licensingserver/pkg/clients"
	"github.com/solo-io/solo-projects/projects/licensingserver/pkg/server"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {

	var s server.Settings
	err := envconfig.Process("", &s)
	if err != nil {
		return err
	}

	stats.StartStatsServer()
	debugMode := s.Debug == 1
	ctx := contextutils.WithLogger(context.Background(), "apiserver")

	var lc server.LicensingClient
	switch s.Client {
	case "test":
		lc = &testClient{}
	default:
		var kla clients.KeygenAuthConfig
		err = envconfig.Process(clients.KEYGEN_ENV_EXPANSION, &kla)
		lc, err = clients.NewKeygenLicensingClient(&kla)
		if err != nil {
			return err
		}
	}
	return server.Setup(ctx, s, debugMode, lc)
}

type testClient struct{}

func (*testClient) Validate(key string) (bool, error) {
	return true, nil
}
