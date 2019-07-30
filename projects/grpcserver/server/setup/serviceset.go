package setup

import (
	"context"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/utils/setuputils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/envutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/pkg/version"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/config"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/setup"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube"
	settings_values "github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/artifactsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/configsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc"
	us_converter "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/converter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc"
	vs_converter "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation"
	"go.uber.org/zap"
)

type ServiceSet struct {
	UpstreamService       v1.UpstreamApiServer
	ArtifactService       v1.ArtifactApiServer
	ConfigService         v1.ConfigApiServer
	SecretService         v1.SecretApiServer
	VirtualServiceService v1.VirtualServiceApiServer
}

func MustGetServiceSet(ctx context.Context) ServiceSet {
	clientset := mustGetClientset(ctx)

	// Create simple and derived clients
	licenseClient := license.NewClient(ctx)
	namespaceClient := kube.NewNamespaceClient(clientset.CoreV1Interface)
	settingsValues := settings_values.NewSettingsValuesClient(ctx, clientset.SettingsClient)
	upstreamInputConverter := us_converter.NewUpstreamInputConverter()
	mutator := mutation.NewMutator(ctx, clientset.VirtualServiceClient)
	mutationFactory := mutation.NewMutationFactory()
	virtualServiceDetailsConverter := vs_converter.NewVirtualServiceDetailsConverter()

	// Read env
	oAuthUrl, oAuthClient := config.GetOAuthEndpointValues()
	oAuthEndpoint := v1.OAuthEndpoint{Url: oAuthUrl, ClientName: oAuthClient}

	upstreamService := upstreamsvc.NewUpstreamGrpcService(ctx, clientset.UpstreamClient, upstreamInputConverter, settingsValues)
	artifactService := artifactsvc.NewArtifactGrpcService(ctx, clientset.ArtifactClient)
	configService := configsvc.NewConfigGrpcService(ctx, clientset.SettingsClient, licenseClient, namespaceClient, oAuthEndpoint, version.Version)
	secretService := secretsvc.NewSecretGrpcService(ctx, clientset.SecretClient)
	virtualServiceService := virtualservicesvc.NewVirtualServiceGrpcService(ctx, clientset.VirtualServiceClient, settingsValues, mutator, mutationFactory, virtualServiceDetailsConverter)

	return ServiceSet{
		UpstreamService:       upstreamService,
		ArtifactService:       artifactService,
		ConfigService:         configService,
		SecretService:         secretService,
		VirtualServiceService: virtualServiceService,
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

func mustGetSettingsClient(ctx context.Context) gloov1.SettingsClient {
	settingsClient, err := setuputils.KubeOrFileSettingsClient(ctx, "")
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
