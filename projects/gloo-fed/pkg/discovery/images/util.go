package images

import (
	"strings"

	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/util/parsers"
)

const (
	OpenSourceGatewayProxyImageName     = "gloo-envoy-wrapper"
	EnterpriseGatewayProxyImageName     = "gloo-ee-envoy-wrapper"
	OpenSourceGatewayProxyWasmImageName = "gloo-envoy-wasm-wrapper"

	OpenSourceControlPlaneImageName = "gloo"
	EnterpriseControlPlaneImageName = "gloo-ee"
)

func GetGatewayProxyImage(podTemplate *core_v1.PodTemplateSpec) (image string, wasmEnabled, isProxy bool) {
	for _, container := range podTemplate.Spec.Containers {
		repo, _, _, err := parsers.ParseImageName(container.Image)
		if err != nil {
			continue
		}

		if strings.HasSuffix(repo, OpenSourceGatewayProxyImageName) ||
			strings.HasSuffix(repo, OpenSourceGatewayProxyWasmImageName) ||
			strings.HasSuffix(repo, EnterpriseGatewayProxyImageName) {
			return container.Image, strings.HasSuffix(repo, OpenSourceGatewayProxyWasmImageName), true
		}
	}
	return "", false, false
}

func GetControlPlaneImage(deployment *v1.Deployment) (image string, isEnterprise, isControlPlane bool) {
	for _, container := range deployment.Spec.Template.Spec.Containers {
		repo, _, _, err := parsers.ParseImageName(container.Image)
		if err != nil {
			continue
		}

		if strings.HasSuffix(repo, OpenSourceControlPlaneImageName) ||
			strings.HasSuffix(repo, EnterpriseControlPlaneImageName) {
			return container.Image, strings.HasSuffix(repo, EnterpriseControlPlaneImageName), true
		}
	}
	return "", false, false
}

func GetVersion(image string) string {
	_, tag, _, err := parsers.ParseImageName(image)
	if err != nil {
		return ""
	}
	return tag
}
