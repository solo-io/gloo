package setup

import (
	gatewayV1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooV1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	sqoopV1 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"
	"k8s.io/client-go/kubernetes"
)

type ApiServerOpts struct {
	UpstreamsRCF       factory.ResourceClientFactory
	VirtualServicesRCF factory.ResourceClientFactory
	SecretsRCF         factory.ResourceClientFactory
	ArtifactsRCF       factory.ResourceClientFactory
	SchemasRCF         factory.ResourceClientFactory
	ResolverMapsRCF    factory.ResourceClientFactory
	SettingsClient     glooV1.SettingsClient
}

func InitOpts() (ApiServerOpts, error) {

	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return ApiServerOpts{}, err
	}
	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return ApiServerOpts{}, err
	}

	cache := kube.NewKubeCache()

	// TODO(ilackarms): pass in settings configuration from an environment variable or CLI flag, rather than hard-coding to k8s
	settingsClient, err := glooV1.NewSettingsClient(&factory.KubeResourceClientFactory{
		Crd:         glooV1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: cache,
	})
	if err := settingsClient.Register(); err != nil {
		return ApiServerOpts{}, err
	}

	return ApiServerOpts{
		UpstreamsRCF: &factory.KubeResourceClientFactory{
			Crd:         glooV1.UpstreamCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		VirtualServicesRCF: &factory.KubeResourceClientFactory{
			Crd:         gatewayV1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		SchemasRCF: &factory.KubeResourceClientFactory{
			Crd:         sqoopV1.SchemaCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		ResolverMapsRCF: &factory.KubeResourceClientFactory{
			Crd:         sqoopV1.ResolverMapCrd,
			Cfg:         cfg,
			SharedCache: cache,
		},
		SecretsRCF: &factory.KubeSecretClientFactory{
			Clientset: clientSet,
		},
		ArtifactsRCF: &factory.KubeConfigMapClientFactory{
			Clientset: clientSet,
		},
		SettingsClient: settingsClient,
	}, nil
}
