package images

import (
	"strings"

	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/kubernetes/pkg/util/parsers"
)

const (
	OpenSourceGatewayProxyImageName     = "gloo-envoy-wrapper"
	EnterpriseGatewayProxyImageName     = "gloo-ee-envoy-wrapper"
	OpenSourceGatewayProxyWasmImageName = "gloo-envoy-wasm-wrapper"

	OpenSourceControlPlaneImageName = "gloo"
	EnterpriseControlPlaneImageName = "gloo-ee"

	ControlPlaneDeploymentName = "gloo"
)

var (
	GatewayProxyPodLabels = map[string]string{
		"gloo": "gateway-proxy",
	}

	ControlPlaneDeploymentLabels = map[string]string{
		"app":  "gloo",
		"gloo": "gloo",
	}
)

func GetGatewayProxyImage(podTemplate *core_v1.PodTemplateSpec) (tag string, wasmEnabled, isProxy bool) {
	// make sure it has the gateway-proxy labels
	if !labels.SelectorFromSet(GatewayProxyPodLabels).Matches(labels.Set(podTemplate.GetLabels())) {
		return "", false, false
	}

	for _, container := range podTemplate.Spec.Containers {
		repo, tag, _, err := parsers.ParseImageName(container.Image)
		if err != nil {
			continue
		}

		if strings.Contains(repo, OpenSourceGatewayProxyImageName) ||
			strings.Contains(repo, OpenSourceGatewayProxyWasmImageName) ||
			strings.Contains(repo, EnterpriseGatewayProxyImageName) {
			return tag, strings.Contains(repo, OpenSourceGatewayProxyWasmImageName), true
		}
	}
	return "", false, false
}

func GetControlPlaneImage(deployment *v1.Deployment) (tag string, isEnterprise, isControlPlane bool) {
	// make sure it has the gloo name and labels
	if deployment.GetName() != ControlPlaneDeploymentName || !labels.SelectorFromSet(ControlPlaneDeploymentLabels).Matches(labels.Set(deployment.GetLabels())) {
		return "", false, false
	}

	for _, container := range deployment.Spec.Template.Spec.Containers {
		repo, tag, _, err := parsers.ParseImageName(container.Image)
		if err != nil {
			continue
		}

		if strings.Contains(repo, OpenSourceControlPlaneImageName) ||
			strings.Contains(repo, EnterpriseControlPlaneImageName) {
			return tag, strings.Contains(repo, EnterpriseControlPlaneImageName), true
		}
	}
	return "", false, false
}
