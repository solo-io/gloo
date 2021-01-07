package images

import (
	"strings"

	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/util/parsers"
)

const (
	OpenSourceGatewayProxyImageName = "gloo-envoy-wrapper"
	EnterpriseGatewayProxyImageName = "gloo-ee-envoy-wrapper"

	OpenSourceControlPlaneImageName = "gloo"
	EnterpriseControlPlaneImageName = "gloo-ee"
)

func GetGatewayProxyImage(podTemplate *core_v1.PodTemplateSpec) (image string, wasmEnabled, isProxy bool) {
	return getImageInfo(podTemplate.Spec.Containers,
		OpenSourceGatewayProxyImageName, EnterpriseGatewayProxyImageName)
}

func GetControlPlaneImage(deployment *v1.Deployment) (image string, isEnterprise, isControlPlane bool) {
	return getImageInfo(deployment.Spec.Template.Spec.Containers,
		OpenSourceControlPlaneImageName, EnterpriseControlPlaneImageName)
}

func getImageInfo(containers []core_v1.Container, ossImageName, enterpriseImageName string) (image string, enterprise, ok bool) {
	for _, container := range containers {
		repo, _, _, err := parsers.ParseImageName(container.Image)
		if err != nil {
			continue
		}

		if strings.HasSuffix(repo, ossImageName) ||
			strings.HasSuffix(repo, enterpriseImageName) {
			return container.Image, strings.HasSuffix(repo, enterpriseImageName), true
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
