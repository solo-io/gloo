package rawgetter

import (
	"github.com/ghodss/yaml"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
)

type kubeYamlRawGetter struct{}

var _ RawGetter = kubeYamlRawGetter{}

func NewKubeYamlRawGetter() RawGetter {
	return kubeYamlRawGetter{}
}

func (kubeYamlRawGetter) GetRaw(in resources.InputResource, resourceCrd crd.Crd) (*v1.Raw, error) {
	content, err := yaml.Marshal(resourceCrd.KubeResource(in))
	if err != nil {
		return nil, err
	}
	return &v1.Raw{
		FileName: in.GetMetadata().Name + ".yaml",
		Content:  string(content),
	}, nil
}
