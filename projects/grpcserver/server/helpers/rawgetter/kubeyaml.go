package rawgetter

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/go-utils/protoutils"
	kubecrd "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"

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

// init `emptyInputResource` from the resource provided in the yaml
func (kubeYamlGetter) InitResourceFromYamlString(ctx context.Context,
	yamlString string,
	refToValidate *core.ResourceRef,
	emptyInputResource resources.InputResource,
) error {
	resourceFromYaml := &kubecrd.Resource{}
	unmarshalError := yaml.Unmarshal([]byte(yamlString), resourceFromYaml)

	if unmarshalError != nil {
		return UnmarshalError(unmarshalError, refToValidate)
	}

	// we enforce that:
	// 1. there is a resource version provided, and
	// 2. the name and namespace have not been edited
	if resourceFromYaml.ResourceVersion == "" {
		return NoResourceVersionError(refToValidate)
	}
	if resourceFromYaml.Namespace != refToValidate.Namespace ||
		resourceFromYaml.Name != refToValidate.Name {
		return EditedRefError(refToValidate)
	}

	emptyInputResource.SetMetadata(kubeutils.FromKubeMeta(resourceFromYaml.ObjectMeta))

	if withStatus, ok := emptyInputResource.(resources.InputResource); ok {
		resources.UpdateStatus(withStatus, func(status *core.Status) {
			*status = resourceFromYaml.Status
		})
	}

	if err := protoutils.UnmarshalMap(*resourceFromYaml.Spec, emptyInputResource.(proto.Message)); err != nil {
		return FailedToReadCrdSpec(err, refToValidate)
	}

	return nil
}
