package discovery

import (
	"github.com/solo-io/solo-kit/pkg/bootstrap"
)

func KubeUpstreamDiscovery(bstrp bootstrap.Config) {
	kube := bstrp.KubeClient()

}
