package setup

import (
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/setuputils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
	"time"
)

func Main(devMode bool, settingsDir string) error {
	// TODO(ilackarms) provide a way to create bootstrap settings
	if devMode {
		settingsClient, err := setuputils.KubeOrFileSettingsClient(settingsDir)
		if err != nil {
			return err
		}
		if err := settingsClient.Register(); err != nil {
			return err
		}
		if err := writeSettings(settingsClient); err != nil && !errors.IsExist(err) {
			return err
		}
	}

	return setuputils.Main("gloo", syncer.NewSetupSyncer(memory.NewInMemoryResourceCache(), kube.NewKubeCache()), settingsDir)
}

// TODO(ilackarms): remove this or move it to a test package, only use settings watch for prodution gloo
func writeSettings(cli v1.SettingsClient) error {
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
			Namespace: defaults.GlooSystem,
			Name:      "gloo",
		},
	}
	_, err := cli.Write(settings, clients.WriteOpts{})
	return err
}
