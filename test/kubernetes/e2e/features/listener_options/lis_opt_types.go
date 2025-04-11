package listener_options

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	e2edefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Port numbers and mappings match the ports in the setup.yaml file
	gw1port1 = 8080
	gw1port2 = 8081
	// port 8082 is used by envoy's readiness probe
	gw2port1 = 8083
	gw2port2 = 8084
	ls1port1 = 8085
	ls1port2 = 8086
)

var (
	setup = func(ti *e2e.TestInstallation) base.SimpleTestCase {
		return base.SimpleTestCase{
			Manifests: setupManifests(ti),
			Resources: []client.Object{proxy1Service, proxy1Deployment, proxy2Service, proxy2Deployment, nginxPod, defaults.CurlPod},
		}
	}

	testCases = map[string]*base.TestCase{
		"TestConfigureListenerOptions": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{basicLisOptManifest},
			},
		},
		"TestConfigureListenerOptionsWithSectionedTargetRefs": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{basicLisOptManifest, lisOptWithSectionedTargetRefsManifest},
			},
		},
	}

	commonSetupManifests = []string{
		filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml"),
		e2edefaults.CurlPodManifest,
	}
	gw1NoListenerSetManifest              = filepath.Join(util.MustGetThisDir(), "testdata", "gw1-no-listenerset.yaml")
	gw1ListenerSetManifest                = filepath.Join(util.MustGetThisDir(), "testdata", "gw1-listenerset.yaml")
	listenerSetManifest                   = filepath.Join(util.MustGetThisDir(), "testdata", "listener-set.yaml")
	basicLisOptManifest                   = filepath.Join(util.MustGetThisDir(), "testdata", "basic-lisopt.yaml")
	lisOptWithSectionedTargetRefsManifest = filepath.Join(util.MustGetThisDir(), "testdata", "listopt-with-sectioned-target-refs.yaml")
	lisOptWithListenerSetRefsManifest     = filepath.Join(util.MustGetThisDir(), "testdata", "listopt-with-listenerset-refs.yaml")

	// When we apply the setup file, we expect resources to be created with this metadata
	glooProxy1ObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw-1",
		Namespace: "default",
	}
	proxy1Service     = &corev1.Service{ObjectMeta: glooProxy1ObjectMeta}
	proxy1ServiceFqdn = kubeutils.ServiceFQDN(proxy1Service.ObjectMeta)
	proxy1Deployment  = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gloo-proxy-gw-1",
			Namespace: "default",
		},
	}
	glooProxy2ObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw-2",
		Namespace: "default",
	}
	proxy2Service     = &corev1.Service{ObjectMeta: glooProxy2ObjectMeta}
	proxy2ServiceFqdn = kubeutils.ServiceFQDN(proxy2Service.ObjectMeta)
	proxy2Deployment  = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gloo-proxy-gw-2",
			Namespace: "default",
		},
	}
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

	expectedHealthyResponse = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gomega.ContainSubstring("Welcome to nginx!"),
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
