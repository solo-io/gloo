package rawgetter

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	v1 "github.com/solo-io/solo-projects/projects/grpcserver/api/v1"
	"go.uber.org/zap"
)

var (
	FailedToGetKubeYaml = func(resource, namespace, name string) string {
		return fmt.Sprintf("Failed to get kubernetes yaml for %v %v.%v", resource, namespace, name)
	}
)

type kubeYamlGetter struct{}

var _ RawGetter = kubeYamlGetter{}

func NewKubeYamlRawGetter() RawGetter {
	return kubeYamlGetter{}
}

func (kubeYamlGetter) GetRaw(ctx context.Context, in resources.InputResource, resourceCrd crd.Crd) *v1.Raw {
	var contentRenderError string
	content, err := yaml.Marshal(resourceCrd.KubeResource(in))
	if err != nil {
		contentRenderError = FailedToGetKubeYaml(resourceCrd.KindName, in.GetMetadata().Namespace, in.GetMetadata().Name)
		contextutils.LoggerFrom(ctx).Warnw(contentRenderError, zap.Error(err), zap.Any("resource", in))
	}
	return &v1.Raw{
		FileName:           in.GetMetadata().Name + ".yaml",
		Content:            string(content),
		ContentRenderError: contentRenderError,
	}
}
