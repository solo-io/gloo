package setup

import (
	"context"
	"os"

	"github.com/solo-io/go-utils/kubeutils"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/pkg/version"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

func MustSettings(ctx context.Context, podNamespace string) *gloov1.Settings {
	return mustGetSettings(ctx, podNamespace)
}

func NewOAuthEndpoint() v1.OAuthEndpoint {
	return v1.OAuthEndpoint{Url: os.Getenv("OAUTH_SERVER"), ClientName: os.Getenv("OAUTH_CLIENT")}
}

type BuildVersion string

func GetBuildVersion() BuildVersion {
	return BuildVersion(version.Version)
}

func NewKubeConfig() (*rest.Config, error) {
	// When running in-cluster, this configuration will hold a token associated with the pod service account
	return kubeutils.GetConfig("", "")
}

// have to type this as something else because wire is already injecting a string
type Token *string

func GetToken(cfg *rest.Config) (Token, error) {
	var token string

	// TODO: temporary solution to bypass authentication.
	if token == "" {
		// Use the token associated with the pod service account
		token = cfg.BearerToken
	}

	return &token, nil
}

func GetK8sCoreInterface(cfg *rest.Config) (corev1.CoreV1Interface, error) {
	kubeClientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	return kubeClientset.CoreV1(), nil
}

func NewNamespacesGetter(coreInterface corev1.CoreV1Interface) corev1.NamespacesGetter {
	return coreInterface
}

func NewPodsGetter(coreInterface corev1.CoreV1Interface) corev1.PodsGetter {
	return coreInterface
}
