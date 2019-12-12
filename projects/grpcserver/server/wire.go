//+build wireinject

package server

import (
	"context"
	"net"

	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/client"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/search"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamsvc/truncate"

	"github.com/google/wire"
	"github.com/solo-io/go-utils/envutils"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/rawgetter"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/helpers/status"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/kube"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/internal/settings"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/artifactsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/configsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/envoysvc/envoydetails"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/gatewaysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/proxysvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/routetablesvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/secretsvc/scrub"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/upstreamgroupsvc"
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

		setup.NewKubeConfig,
		setup.GetToken,
		setup.GetK8sCoreInterface,
		setup.NewNamespacesGetter,
		setup.NewPodsGetter,

		// Resource clients.
		client.NewClientCache,
		client.NewClientUpdater,

		// Derived and simple clients.
		license.NewClient,
		kube.NewNamespaceClient,
		settings.NewSettingsValuesClient,
		vs_mutation.NewMutator,
		vs_mutation.NewMutationFactory,
		rawgetter.NewKubeYamlRawGetter,
		search.NewUpstreamSearcher,
		converter.NewVirtualServiceDetailsConverter,
		selection.NewVirtualServiceSelector,
		us_mutation.NewMutator,
		us_mutation.NewFactory,
		envoydetails.NewClient,
		envoydetails.NewHttpGetter,
		envoydetails.NewProxyStatusGetter,
		scrub.NewScrubber,
		status.NewInputResourceStatusGetter,
		truncate.NewUpstreamTruncator,
		wire.Bind(new(truncate.UpstreamTruncator), truncate.Truncator{}),

		// Services
		upstreamsvc.NewUpstreamGrpcService,
		upstreamgroupsvc.NewUpstreamGroupGrpcService,
		artifactsvc.NewArtifactGrpcService,
		configsvc.NewConfigGrpcService,
		secretsvc.NewSecretGrpcService,
		virtualservicesvc.NewVirtualServiceGrpcService,
		routetablesvc.NewRouteTableGrpcService,
		gatewaysvc.NewGatewayGrpcService,
		proxysvc.NewProxyGrpcService,
		envoysvc.NewEnvoyGrpcService,
		NewGlooGrpcService,
	)
	return &GlooGrpcService{}, nil
}
