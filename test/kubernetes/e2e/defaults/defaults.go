package defaults

import (
	"context"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/skv2/codegen/util"
)

var (
	CurlPodExecOpt = kubectl.PodExecOptions{
		Name:      "curl",
		Namespace: "curl",
		Container: "curl",
	}

	CurlPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "curl",
			Namespace: "curl",
		},
	}

	CurlPodManifest = filepath.Join(util.MustGetThisDir(), "testdata", "curl_pod.yaml")
)

func SetupCurlPod(ctx context.Context, t *e2e.TestInstallation) error {
	err := t.Actions.Kubectl().ApplyFile(ctx, CurlPodManifest)
	if err != nil {
		return err
	}
	t.Assertions.EventuallyPodsRunning(ctx, CurlPod.ObjectMeta.GetNamespace(), metav1.ListOptions{})
	return nil
}

func TeardownCurlPod(ctx context.Context, t *e2e.TestInstallation) error {
	return t.Actions.Kubectl().DeleteFile(ctx, CurlPodManifest)
}
