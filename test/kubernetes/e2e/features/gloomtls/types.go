package gloomtls

import (
	"net/http"
	"path/filepath"

	. "github.com/onsi/gomega"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	expectedHealthyResponse = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       ContainSubstring("Welcome to nginx!"),
	}

	edgeRoutingResources = filepath.Join(util.MustGetThisDir(), "testdata", "edge_resources.yaml")

	edgeGatewaySetupSuite = base.SimpleTestCase{
		Manifests: []string{defaults.CurlPodManifest, defaults.NginxPodManifest},
		Resources: []client.Object{defaults.CurlPod, defaults.NginxPod},
	}

	// K8s Gateway tests
	routeWithServiceManifest = filepath.Join(util.MustGetThisDir(), "testdata", "route-with-service.yaml")
	serviceManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "service-for-route.yaml")

	glooProxyObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	proxyDeployment = &appsv1.Deployment{ObjectMeta: glooProxyObjectMeta}
	proxyService    = &corev1.Service{ObjectMeta: glooProxyObjectMeta}

	k8sGatewayTestCases = map[string]*base.TestCase{
		"TestRouteSecureRequestToUpstream": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{defaults.CurlPodManifest, serviceManifest, routeWithServiceManifest},
				Resources: []client.Object{proxyDeployment, proxyService, defaults.CurlPod},
			},
		},
	}
)
