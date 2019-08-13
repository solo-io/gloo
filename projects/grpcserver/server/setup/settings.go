package setup

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/envutils"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	kube2 "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"go.uber.org/zap"
)

func mustGetSettings(ctx context.Context) *gloov1.Settings {
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

func kubeSettingsClient(ctx context.Context) (gloov1.SettingsClient, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}
	return gloov1.NewSettingsClient(&factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: kube2.NewKubeCache(ctx),
	})
}

func mustGetSettingsClient(ctx context.Context) gloov1.SettingsClient {
	settingsClient, err := kubeSettingsClient(ctx)
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not create settings client", zap.Error(err))
	}
	if err := settingsClient.Register(); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not register settings client", zap.Error(err))
	}
	return settingsClient
}

func writeDefaultSettings(settingsNamespace, name string, cli gloov1.SettingsClient) error {
	settings := &gloov1.Settings{
		ConfigSource: &gloov1.Settings_KubernetesConfigSource{
			KubernetesConfigSource: &gloov1.Settings_KubernetesCrds{},
		},
		ArtifactSource: &gloov1.Settings_KubernetesArtifactSource{
			KubernetesArtifactSource: &gloov1.Settings_KubernetesConfigmaps{},
		},
		SecretSource: &gloov1.Settings_KubernetesSecretSource{
			KubernetesSecretSource: &gloov1.Settings_KubernetesSecrets{},
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
