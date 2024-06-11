package kube

import (
	"context"
	"fmt"
	"io"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/testutils/hack/go-utils/testutils/kubectl"
)

var (
	DefaultCurlImage = Image{
		Registry:   "curlimages",
		Repository: "curl",
		Tag:        "",
		Digest:     "aa45e9d93122a3cfdf8d7de272e2798ea63733eeee6d06bd2ee4f2f8c4027d7c",
		PullPolicy: "IfNotPresent",
	}
)

func CurlWithEphemeralPod(ctx context.Context, logger io.Writer, kubeContext, fromns, frompod string, args ...string) string {
	execParams := EphemeralPodParams{
		Logger:        logger,
		KubeContext:   kubeContext,
		Image:         DefaultCurlImage,
		FromContainer: "curl",
		FromNamespace: fromns,
		FromPod:       frompod,
		ExecCmdPath:   "curl",
		Args:          args,
	}
	out, _ := ExecFromEphemeralPod(ctx, execParams)
	return out
}

// labelSelector is a string map e.g. gloo=gateway-proxy
func FindPodNameByLabel(cfg *rest.Config, ctx context.Context, ns, labelSelector string) string {
	clientset, err := kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
	pl, err := clientset.CoreV1().Pods(ns).List(ctx, v1.ListOptions{LabelSelector: labelSelector})
	Expect(err).NotTo(HaveOccurred())
	Expect(pl.Items).NotTo(BeEmpty())
	return pl.Items[0].GetName()
}

func WaitForRollout(ctx context.Context, logger io.Writer, kubeContext, ns, deployment string) {
	args := []string{"-n", ns, "rollout", "status", "deployment", deployment}
	mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
}

func Curl(ctx context.Context, logger io.Writer, kubeContext, ns, fromDeployment, fromContainer, url string) string {
	args := []string{
		"-n", ns,
		"exec", fmt.Sprintf("deployment/%s", fromDeployment),
		"-c", fromContainer,
		"--", "curl", url,
	}
	return mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
}

func CreateNamespace(ctx context.Context, logger io.Writer, kubeContext, ns string) {
	args := []string{"create", "namespace", ns}
	out := mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
	fmt.Fprintln(logger, out)
}

func DeleteNamespace(ctx context.Context, logger io.Writer, kubeContext, ns string) {
	args := []string{"delete", "namespace", ns}
	out := mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
	fmt.Fprintln(logger, out)
}

func LabelNamespace(ctx context.Context, logger io.Writer, kubeContext, ns, label string) {
	args := []string{"label", "namespace", ns, label}
	out := mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
	fmt.Fprintln(logger, out)
}

func SetDeploymentEnvVars(
	ctx context.Context,
	kubeContext string,
	logger io.Writer,
	ns string,
	deploymentName string,
	containerName string,
	envVars map[string]string,
) {
	var envVarStrings []string
	for name, value := range envVars {
		envVarStrings = append(envVarStrings, fmt.Sprintf("%s=%s", name, value))
	}
	args := append([]string{"set", "env", "-n", ns, fmt.Sprintf("deployment/%s", deploymentName), "-c", containerName}, envVarStrings...)
	out := mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
	fmt.Fprintln(logger, out)
}

func DisableContainer(
	ctx context.Context,
	logger io.Writer,
	kubeContext string,
	ns string,
	deploymentName string,
	containerName string,
) {
	args := []string{
		"-n", ns,
		"patch", "deployment", deploymentName,
		"--patch",
		fmt.Sprintf(`{"spec": {"template": {"spec": {"containers": [{"name": "%s","command": ["sleep", "20h"]}]}}}}`,
			containerName),
	}
	out := mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
	fmt.Fprintln(logger, out)
}

func EnableContainer(
	ctx context.Context,
	logger io.Writer,
	kubeContext string,
	ns string,
	deploymentName string,
) {
	args := []string{
		"-n", ns,
		"patch", "deployment", deploymentName,
		"--type", "json",
		"-p", `[{"op": "remove", "path": "/spec/template/spec/containers/0/command"}]`,
	}
	out := mustExecute(ctx, KubectlParams{
		KubectlCmdParams: kubectl.NewParams(args...),
		KubeContext:      kubeContext,
		Logger:           logger,
	})
	fmt.Fprintln(logger, out)
}
