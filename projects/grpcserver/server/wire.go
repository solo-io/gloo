//+build wireinject

package server

import (
	"context"
	"net"

	"github.com/google/wire"
	"github.com/solo-io/go-utils/envutils"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/artifactsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/configsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/gatewaysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/proxysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc"
	us_mutation "github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/mutation"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/converter"
	vs_mutation "github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/mutation"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/virtualservicesvc/selection"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"
)

func InitializeServer(ctx context.Context, listener net.Listener) (*GlooGrpcService, error) {
	wire.Build(
		// If more strings are provided after PodNamespace, they will have to be typed. see configsvc.BuildVersion
		envutils.MustGetPodNamespace,
		setup.GetBuildVersion,
		setup.MustSettings,
		setup.NewOAuthEndpoint,

		// Resource clients.
		setup.NewClientSet,
		setup.NewNamespacesGetter,
		setup.NewPodsGetter,
		setup.NewSettingsClient,
		setup.NewVirtualServiceClient,
		setup.NewUpstreamClient,
		setup.NewArtifactClient,
		setup.NewSecretClient,
		setup.NewGatewayClient,
		setup.NewProxyClient,

		// Derived and simple clients.
		license.NewClient,
		kube.NewNamespaceClient,
		settings.NewSettingsValuesClient,
		vs_mutation.NewMutator,
		vs_mutation.NewMutationFactory,
		rawgetter.NewKubeYamlRawGetter,
		converter.NewVirtualServiceDetailsConverter,
		selection.NewVirtualServiceSelector,
		us_mutation.NewMutator,
		us_mutation.NewFactory,
		envoydetails.NewClient,
		envoydetails.NewHttpGetter,

		// Services
		upstreamsvc.NewUpstreamGrpcService,
		artifactsvc.NewArtifactGrpcService,
		configsvc.NewConfigGrpcService,
		secretsvc.NewSecretGrpcService,
		virtualservicesvc.NewVirtualServiceGrpcService,
		gatewaysvc.NewGatewayGrpcService,
		proxysvc.NewProxyGrpcService,
		envoysvc.NewEnvoyGrpcService,
		NewGlooGrpcService,
	)
	return &GlooGrpcService{}, nil
}
