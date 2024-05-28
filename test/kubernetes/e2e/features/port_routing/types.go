package port_routing

import (
	"net/http"
	"path/filepath"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/skv2/codegen/util"

	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
)

type testManifest struct {
	manifestFile string
	extraArgs    []string
}

var (
	setupManifest = filepath.Join(util.MustGetThisDir(), "testdata/setup.yaml")

	k8sHTTPRouteInvalidPortAndValidTargetportManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-backendref-routing", "invalid-port-and-valid-targetport.yaml")
	k8sHTTPRouteInvalidPortAndInvalidTargetportManifest = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-backendref-routing", "invalid-port-and-invalid-targetport.yaml")
	k8sHTTPRouteMatchPodPortWithoutTargetportManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-backendref-routing", "match-pod-port-without-targetport.yaml")
	k8sHTTPRouteMatchPortandTargetportManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-backendref-routing", "match-port-and-targetport.yaml")
	k8sHTTPRouteInvalidPortWithoutTargetportManifest    = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-backendref-routing", "invalid-port-without-targetport.yaml")

	k8sGatewayBackendrefRoutingManifest = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-backendref-routing", "gw.yaml")
	k8sGatewayUpstreamRoutingManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "k8s-gateway-upstream-routing", "routing.yaml")
	glooGatewayRoutingManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "gloo-gateway-routing", "routing.yaml")

	upstreamInvalidPortAndValidTargetportManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "upstreams", "invalid-port-and-valid-targetport.yaml")
	upstreamInvalidPortAndInvalidTargetportManifest = filepath.Join(util.MustGetThisDir(), "testdata", "upstreams", "invalid-port-and-invalid-targetport.yaml")
	upstreamMatchPodPortWithoutTargetportManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "upstreams", "match-pod-port-without-targetport.yaml")
	upstreamMatchPortandTargetportManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "upstreams", "match-port-and-targetport.yaml")
	upstreamInvalidPortWithoutTargetportManifest    = filepath.Join(util.MustGetThisDir(), "testdata", "upstreams", "invalid-port-without-targetport.yaml")

	svcInvalidPortAndValidTargetportManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "svc", "invalid-port-and-valid-targetport.yaml")
	svcInvalidPortAndInvalidTargetportManifest = filepath.Join(util.MustGetThisDir(), "testdata", "svc", "invalid-port-and-invalid-targetport.yaml")
	svcMatchPodPortWithoutTargetportManifest   = filepath.Join(util.MustGetThisDir(), "testdata", "svc", "match-pod-port-without-targetport.yaml")
	svcMatchPortandTargetportManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "svc", "match-port-and-targetport.yaml")
	svcInvalidPortWithoutTargetportManifest    = filepath.Join(util.MustGetThisDir(), "testdata", "svc", "invalid-port-without-targetport.yaml")

	// When we apply the setup.yaml file, we expect resources to be created with this metadata
	glooProxyObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	proxyDeployment = &appsv1.Deployment{ObjectMeta: glooProxyObjectMeta}
	proxyService    = &corev1.Service{ObjectMeta: glooProxyObjectMeta}

	// curlPod is the Pod that will be used to execute curl requests, and is defined in the port routing setup.yaml manifest files
	curlPodExecOpt = kubectl.PodExecOptions{
		Name:      "curl",
		Namespace: "curl",
		Container: "curl",
	}

	exampleSvc = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-svc",
			Namespace: "default",
		},
	}

	nginxPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
		},
	}

	expectedHealthyResponse = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       ContainSubstring("Welcome to nginx!"),
	}

	expectedServiceUnavailableResponse = &testmatchers.HttpResponse{
		StatusCode: http.StatusServiceUnavailable,
		Body:       gstruct.Ignore(), // ignore the body
	}
)
