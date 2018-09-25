package setup

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
)

// TODO(ilackarms): remove this or move it to a test package, only use settings watch for prodution gloo
func writeSettings(settingsDir string) error {
	cli, err := v1.NewSettingsClient(&factory.FileResourceClientFactory{RootDir: settingsDir})
	if err != nil {
		return err
	}
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
		BindAddr:    "0.0.0.0:9977",
		RefreshRate: types.DurationProto(time.Minute),
		DevMode:     true,
		Metadata: core.Metadata{
			Namespace: "settings",
			Name:      "gloo",
		},
	}
	_, err = cli.Write(settings, clients.WriteOpts{})
	return err
}

func Main(settingsDir string) error {
	if err := writeSettings(settingsDir); err != nil && !errors.IsExist(err) {
		return err
	}
	settingsClient, err := v1.NewSettingsClient(&factory.FileResourceClientFactory{
		RootDir: settingsDir,
	})
	if err != nil {
		return err
	}
	cache := v1.NewSetupEmitter(settingsClient)
	ctx := contextutils.WithLogger(context.Background(), "gloo")
	eventLoop := v1.NewSetupEventLoop(cache, syncer.NewSetupSyncer())
	errs, err := eventLoop.Run([]string{"settings"}, clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: time.Second,
	})
	if err != nil {
		return err
	}
	for err := range errs {
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}
	return nil
}
