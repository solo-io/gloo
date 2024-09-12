package serviceentry

import (
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	kubeclient "istio.io/istio/pkg/kube"

	// side-effect to setup CRDWatcher builder
	_ "istio.io/istio/pkg/kube/kclient"
)

const PluginName = "ServiceEntryDiscoveryPlugin"

// TODO this is both a upstreams.ClientPlugin and  DiscoveryPlugin. split?
type sePlugin struct {
	client kubeclient.Client

	settings *v1.Settings
}

func NewPlugin() plugins.Plugin {
	return &sePlugin{
		// TODO build client using shared rest config
		client: mustBuildIstioClient(),
	}
}

// TODO share rest cfg setup with the rest of the app
// TODO move this init somewhere we can handle the err
func mustBuildIstioClient() kubeclient.Client {
	restCfg, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		panic(err)
	}
	client, err := kubeclient.NewClient(kubeclient.NewClientConfigForRestConfig(restCfg), "")
	if err != nil {
		// TODO move this init somewhere we can handle the err
		panic(err)
	}

	kubeclient.EnableCrdWatcher(client)

	return client
}

func (s *sePlugin) Init(params plugins.InitParams) {
	s.settings = params.Settings
}

func (s *sePlugin) Name() string {
	return PluginName
}
