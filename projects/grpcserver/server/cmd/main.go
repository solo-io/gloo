package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"

	"github.com/solo-io/gloo/pkg/utils/setuputils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/stats"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/setup"
	"github.com/solo-io/solo-projects/projects/grpcserver/server"

	"github.com/solo-io/solo-projects/pkg/version"

	"github.com/solo-io/go-utils/envutils"

	"go.uber.org/zap"

	"github.com/solo-io/go-utils/contextutils"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	ctx := getInitialContext()
	startStatsIfConfigured()
	grpcPort := envutils.MustGetGrpcPort(ctx)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to setup listener",
			zap.Any("listener", lis),
			zap.Error(err))
	}
	glooGrpcService := mustGetGlooGrpcService(ctx, lis)
	if err := glooGrpcService.Run(ctx); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed while running gloo grpc service", zap.Error(err))
	}
}

func mustGetClientset(ctx context.Context) *setup.ClientSet {
	settings := mustGetSettings(ctx)
	// TODO: figure out how auth should work
	clientset, err := setup.NewClientSet(ctx, settings, "")
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to generate clientset", zap.Error(err))
	}
	return clientset
}

func mustGetGlooGrpcService(ctx context.Context, listener net.Listener) *server.GlooGrpcService {
	clientset := mustGetClientset(ctx)
	return server.NewGlooGrpcService(listener, *clientset)
}

func mustGetSettingsClient(ctx context.Context) v1.SettingsClient {
	settingsClient, err := setuputils.KubeOrFileSettingsClient(ctx, "")
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not create settings client", zap.Error(err))
	}
	if err := settingsClient.Register(); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not register settings client", zap.Error(err))
	}
	return settingsClient
}

func writeDefaultSettings(settingsNamespace, name string, cli v1.SettingsClient) error {
	settings := &v1.Settings{
		ConfigSource: &v1.Settings_KubernetesConfigSource{
			KubernetesConfigSource: &v1.Settings_KubernetesCrds{},
		},
		ArtifactSource: &v1.Settings_KubernetesArtifactSource{
			KubernetesArtifactSource: &v1.Settings_KubernetesConfigmaps{},
		},
		SecretSource: &v1.Settings_KubernetesSecretSource{
			KubernetesSecretSource: &v1.Settings_KubernetesSecrets{},
		},
		BindAddr:           "0.0.0.0:9977",
		RefreshRate:        types.DurationProto(time.Minute),
		DevMode:            true,
		DiscoveryNamespace: settingsNamespace,
		Metadata:           core.Metadata{Namespace: settingsNamespace, Name: name},
	}
	if _, err := cli.Write(settings, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
		return errors.Wrapf(err, "failed to create default settings")
	}
	return nil
}

func mustGetSettings(ctx context.Context) *v1.Settings {
	settingsClient := mustGetSettingsClient(ctx)
	namespace := envutils.MustGetPodNamespace(ctx)
	name := defaults.SettingsName
	err := writeDefaultSettings(namespace, name, settingsClient)
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to write default settings", zap.Error(err))
	}
	settings, err := settingsClient.Read(namespace, name, clients.ReadOpts{})
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to read settings", zap.Error(err))
	}
	return settings
}

func startStatsIfConfigured() {
	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}
}

func getInitialContext() context.Context {
	loggingContext := []interface{}{"version", version.Version}
	ctx := contextutils.WithLogger(context.Background(), "gloo-grpcserver")
	ctx = contextutils.WithLoggerValues(ctx, loggingContext...)
	return ctx
}
