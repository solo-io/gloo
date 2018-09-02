package v1

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
)

/*
	Add these two methods to any upstream spec that supports
*/
func (us *UpstreamSpec_Kube) GetServiceSpec() *plugins.ServiceSpec {
	return us.Kube.ServiceSpec
}

func (us *UpstreamSpec_Kube) SetServiceSpec(spec *plugins.ServiceSpec) {
	us.Kube.ServiceSpec = spec
}
