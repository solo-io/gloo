package resources

import (
	"context"
	"path/filepath"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
	//e2edefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/skv2/codegen/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resource struct {
	manifest  string
	validator func(ctx context.Context, t *e2e.TestInstallation)
}

// List of manifests
type ResourceBuilder []Resource

type ResourceBuilderOption func(*ResourceBuilder)

func NewResourceBuilder(opts ...ResourceBuilderOption) ResourceBuilder {
	builder := ResourceBuilder{}
	return builder
}

func (b ResourceBuilder) WithCurl() ResourceBuilder {
	return ResourceBuilder(append(b, Resource{
		manifest: CurlPodManifest,
		validator: func(ctx context.Context, t *e2e.TestInstallation) {
			t.Assertions.EventuallyPodsRunning(ctx, CurlPod.ObjectMeta.GetNamespace(), metav1.ListOptions{})
		},
	}))
}

func (b ResourceBuilder) Install(ctx context.Context, t *e2e.TestInstallation) error {
	for _, resource := range b {
		// apply manifest
		err := t.Actions.Kubectl().ApplyFile(ctx, resource.manifest)
		if err != nil {
			return err
		}
		// validate
		resource.validator(ctx, t)
	}
	return nil
}

func (b ResourceBuilder) Uninstall(ctx context.Context, t *e2e.TestInstallation) error {
	for _, resource := range b {
		// apply manifest
		err := t.Actions.Kubectl().DeleteFile(ctx, resource.manifest)
		if err != nil {
			return err
		}

	}
	return nil
}

var (
	CurlPodManifest = filepath.Join(util.MustGetThisDir(), "../../testdata", "curl_pod.yaml")

	CurlPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "curl",
			Namespace: "curl",
		},
	}
)
