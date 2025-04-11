package http_listener_options

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	gw1port1 = 8080
	gw1port2 = 8081
	// port 8082 is used by envoy's readiness probe
	gw2port1 = 8083
	gw2port2 = 8084
	lsPort1  = 8085
	lsPort2  = 8086
)

var (
	setup = func(ti *e2e.TestInstallation) base.SimpleTestCase {
		return base.SimpleTestCase{
			Manifests: setupManifests(ti),
			Resources: []client.Object{proxy1Service, proxy1Deployment, proxy2Service, proxy2Deployment, nginxPod, defaults.CurlPod},
		}
	}

	testCases = map[string]*base.TestCase{
		"TestConfigureHttpListenerOptions": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{basicLisOptManifest},
			},
		},
		"TestConfigureNotAttachedHttpListenerOptions": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{notAttachedLisOptManifest},
			},
		},
		"TestConfigureHttpListenerOptionsWithSection": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{basicLisOptSectionManifest},
			},
		},
		"TestConfigureHttpListenerOptionsWithListenerSetsAndSection": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{basicLisOptManifest, basicLisOptSectionManifest, basicLisOptListenerSetSectionManifest, basicLisOptListenerSetManifest},
			},
		},
	}

	commonSetupManifests = []string{
		filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml"),
		defaults.CurlPodManifest,
	}
	gw1NoListenerSetManifest              = filepath.Join(util.MustGetThisDir(), "testdata", "gw1-no-listenerset.yaml")
	gw1ListenerSetManifest                = filepath.Join(util.MustGetThisDir(), "testdata", "gw1-listenerset.yaml")
	listenerSetManifest                   = filepath.Join(util.MustGetThisDir(), "testdata", "listener-set.yaml")
	gatewayManifest                       = filepath.Join(util.MustGetThisDir(), "testdata", "gateway.yaml")
	basicLisOptManifest                   = filepath.Join(util.MustGetThisDir(), "testdata", "basic-http-lis-opt.yaml")
	notAttachedLisOptManifest             = filepath.Join(util.MustGetThisDir(), "testdata", "not-attached-http-lis-opt.yaml")
	basicLisOptSectionManifest            = filepath.Join(util.MustGetThisDir(), "testdata", "basic-http-lis-opt-section.yaml")
	basicLisOptListenerSetManifest        = filepath.Join(util.MustGetThisDir(), "testdata", "basic-http-lis-opt-listener-set.yaml")
	basicLisOptListenerSetSectionManifest = filepath.Join(util.MustGetThisDir(), "testdata", "basic-http-lis-opt-listener-set-section.yaml")

	// When we apply the setup file, we expect resources to be created with this metadata
	glooProxy1ObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw-1",
		Namespace: "default",
	}
	proxy1Service     = &corev1.Service{ObjectMeta: glooProxy1ObjectMeta}
	proxy1Deployment  = &appsv1.Deployment{ObjectMeta: glooProxy1ObjectMeta}
	proxyService1Fqdn = kubeutils.ServiceFQDN(proxy1Service.ObjectMeta)

	glooProxy2ObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw-2",
		Namespace: "default",
	}
	proxy2Service     = &corev1.Service{ObjectMeta: glooProxy2ObjectMeta}
	proxy2Deployment  = &appsv1.Deployment{ObjectMeta: glooProxy2ObjectMeta}
	proxyService2Fqdn = kubeutils.ServiceFQDN(proxy2Service.ObjectMeta)

	nginxPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
		},
	}
	exampleSvc = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-svc",
			Namespace: "default",
		},
	}

	expectedResponseWithServer = func(serverName string) *matchers.HttpResponse {
		return &matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring("Welcome to nginx!"),
			Headers: map[string]interface{}{
				"server": serverName,
			},
		}
	}

	defaultExpectedResponseWithServer = expectedResponseWithServer("server-override")

	expectedResponseWithoutServer = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaders(http.Header{"server": {"should-not-attach"}})),
		),
		Body: gomega.ContainSubstring("Welcome to nginx!"),
	}
)

func setupManifests(ti *e2e.TestInstallation) []string {
	manifests := commonSetupManifests
	if listenerset.RequiredCrdExists(ti) {
		manifests = append(manifests, gw1ListenerSetManifest, listenerSetManifest)
	} else {
		manifests = append(manifests, gw1NoListenerSetManifest)
	}
	return manifests
}
