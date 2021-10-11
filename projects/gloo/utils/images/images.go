package images

import (
	"fmt"
	"strings"

	dockerref "github.com/docker/distribution/reference"

	v1 "k8s.io/api/apps/v1"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
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
		repo, tag, _, err := ParseImageName(container.Image)
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
		repo, tag, _, err := ParseImageName(container.Image)
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

// copied from: https://github.com/kubernetes/kubernetes/blob/master/pkg/util/parsers/parsers.go#L29
// this was the only place we relies on the k8s.io/kubernetes dependency so this was the easiest
// way to maintain functionality but remove dependencies.
// this likely belongs in our k8s-utils repo, but I kept it here until there is a broader usage of it
const (
	DefaultImageTag = "latest"
)

// ParseImageName parses a docker image string into three parts: repo, tag and digest.
// If both tag and digest are empty, a default image tag will be returned.
func ParseImageName(image string) (string, string, string, error) {
	named, err := dockerref.ParseNormalizedNamed(image)
	if err != nil {
		return "", "", "", fmt.Errorf("couldn't parse image name: %v", err)
	}

	repoToPull := named.Name()
	var tag, digest string

	tagged, ok := named.(dockerref.Tagged)
	if ok {
		tag = tagged.Tag()
	}

	digested, ok := named.(dockerref.Digested)
	if ok {
		digest = digested.Digest().String()
	}
	// If no tag was specified, use the default "latest".
	if len(tag) == 0 && len(digest) == 0 {
		tag = DefaultImageTag
	}
	return repoToPull, tag, digest, nil
}
