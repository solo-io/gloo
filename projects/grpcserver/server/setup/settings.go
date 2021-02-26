package setup

import (
	"context"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	kube2 "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"go.uber.org/zap"
)

func mustGetSettings(ctx context.Context, podNamespace string) *gloov1.Settings {
	settingsClient := mustGetSettingsClient(ctx, podNamespace)
	name := defaults.SettingsName
	settings, err := settingsClient.Read(podNamespace, name, clients.ReadOpts{})
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Failed to read settings", zap.Error(err))
	}
	return settings
}

func kubeSettingsClient(ctx context.Context, podNamespace string) (gloov1.SettingsClient, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, err
	}

	return gloov1.NewSettingsClient(ctx, &factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: kube2.NewKubeCache(ctx),
		// Restrict this client to the pod namespace in case we're running single-namespace Gloo.
		NamespaceWhitelist: []string{podNamespace},
	})
}

func mustGetSettingsClient(ctx context.Context, podNamespace string) gloov1.SettingsClient {
	settingsClient, err := kubeSettingsClient(ctx, podNamespace)
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not create settings client", zap.Error(err))
	}
	if err := settingsClient.Register(); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not register settings client", zap.Error(err))
	}
	return settingsClient
}
