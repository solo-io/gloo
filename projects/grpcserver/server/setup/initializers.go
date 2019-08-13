package setup

import (
	"context"
	"os"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewayv2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/pkg/version"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/service/configsvc"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func MustSettings(ctx context.Context) *gloov1.Settings {
	return mustGetSettings(ctx)
}

func NewOAuthEndpoint() v1.OAuthEndpoint {
	return v1.OAuthEndpoint{Url: os.Getenv("OAUTH_SERVER"), ClientName: os.Getenv("OAUTH_CLIENT")}
}

func NewCoreV1Interface(set *ClientSet) corev1.CoreV1Interface {
	return set.CoreV1Interface
}

func NewSettingsClient(set *ClientSet) gloov1.SettingsClient {
	return set.SettingsClient
}

func NewVirtualServiceClient(set *ClientSet) gatewayv1.VirtualServiceClient {
	return set.VirtualServiceClient
}

func NewUpstreamClient(set *ClientSet) gloov1.UpstreamClient {
	return set.UpstreamClient
}

func NewSecretClient(set *ClientSet) gloov1.SecretClient {
	return set.SecretClient
}

func NewArtifactClient(set *ClientSet) gloov1.ArtifactClient {
	return set.ArtifactClient
}

func NewGatewayClient(set *ClientSet) gatewayv2.GatewayClient {
	return set.GatewayClient
}

func NewProxyClient(set *ClientSet) gloov1.ProxyClient {
	return set.ProxyClient
}

func GetBuildVersion() configsvc.BuildVersion {
	return configsvc.BuildVersion(version.Version)
}
