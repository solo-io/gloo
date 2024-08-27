package serviceentry

import (
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
)

const PluginName = "ServiceEntryDiscoveryPlugin"

type sePlugin struct {
	istio istioclient.Interface

	kube          kubernetes.Interface
	kubeCoreCache corecache.KubeCoreCache

	settings *v1.Settings
}

func NewPlugin(kube kubernetes.Interface, kubeCoreCache corecache.KubeCoreCache) plugins.Plugin {
	return &sePlugin{
		istio:         mustBuildIstioClient(),
		kube:          kube,
		kubeCoreCache: kubeCoreCache,
	}
}

func mustBuildIstioClient() istioclient.Interface {
	cfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		// TODO move this init somewhere we can handle the err
		panic(err)
	}
	client, err := istioclient.NewForConfig(cfg)
	if err != nil {
		// TODO move this init somewhere we can handle the err
		panic(err)
	}
	return client
}

func (s *sePlugin) Init(params plugins.InitParams) {
	s.settings = params.Settings
}

func (s *sePlugin) Name() string {
	return PluginName
}
