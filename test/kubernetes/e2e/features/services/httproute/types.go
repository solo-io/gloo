package httproute

import (
	"net/http"
	"path/filepath"

	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/skv2/codegen/util"

	"github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/projects/gateway2/crds"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	routeWithServiceManifest = filepath.Join(util.MustGetThisDir(), "testdata", "route-with-service.yaml")
	serviceManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "service-for-route.yaml")
	http2ServiceManifest     = filepath.Join(util.MustGetThisDir(), "testdata", "http2-service-for-route.yaml")
	tcpRouteCrdManifest      = filepath.Join(crds.AbsPathToCrd("tcproute-crd.yaml"))
	// Proxy resource to be translated
	glooProxyObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	proxyDeployment = &appsv1.Deployment{ObjectMeta: glooProxyObjectMeta}
	proxyService    = &corev1.Service{ObjectMeta: glooProxyObjectMeta}

	expectedSvcResp = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gstruct.Ignore(),
	}

	expectedHTTP2SvcResp = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
		Protocol:   "HTTP/2",
	}

	// test cases
	testCases = map[string]*base.TestCase{
		"TestConfigureHTTPRouteBackingDestinationsWithService": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{routeWithServiceManifest, serviceManifest},
				Resources: []client.Object{proxyService, proxyDeployment},
			},
		},
		"TestConfigureHTTPRouteBackingDestinationsWithServiceAndWithoutTCPRoute": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{tcpRouteCrdManifest, routeWithServiceManifest, serviceManifest},
				Resources: []client.Object{proxyService, proxyDeployment},
			},
		},
		"TestHTTP2AppProtocol": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{routeWithServiceManifest, http2ServiceManifest},
				Resources: []client.Object{proxyService, proxyDeployment},
			},
		},
	}
)
