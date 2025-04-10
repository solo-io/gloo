package listenerset

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega/gstruct"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	// manifests
	setupManifest                           = filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml")
	validListenerSetManifest                = filepath.Join(util.MustGetThisDir(), "testdata", "valid-listenerset.yaml")
	invalidListenerSetNotAllowedManifest    = filepath.Join(util.MustGetThisDir(), "testdata", "invalid-listenerset-not-allowed.yaml")
	invalidListenerSetNonExistingGWManifest = filepath.Join(util.MustGetThisDir(), "testdata", "invalid-listenerset-non-existing-gw.yaml")

	// objects
	gatewayObjectMeta = metav1.ObjectMeta{
		Name:      "gw",
		Namespace: "default",
	}
	glooProxyObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	proxyDeployment = &appsv1.Deployment{ObjectMeta: glooProxyObjectMeta}
	proxyService    = &corev1.Service{ObjectMeta: glooProxyObjectMeta}

	exampleSvc = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-svc",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": "nginx",
			},
			Ports: []corev1.ServicePort{
				{
					Port:       8080,
					Protocol:   corev1.ProtocolTCP,
					TargetPort: intstr.FromString("http-web-svc"),
				},
			},
		},
	}
	nginxPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
		},
	}

	// TestValidListenerSet
	validListenerSet = types.NamespacedName{
		Name:      "valid-ls",
		Namespace: "default",
	}

	// TestInvalidListenerSetNotAllowed
	invalidListenerSetNotAllowed = types.NamespacedName{
		Name:      "invalid-ls-not-allowed",
		Namespace: "curl",
	}

	// TestInvalidListenerSetNonExistingGW
	invalidListenerSetNonExistingGW = types.NamespacedName{
		Name:      "invalid-ls-non-existing-gw",
		Namespace: "default",
	}

	expectOK = &testmatchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gstruct.Ignore(),
	}

	expectNotFound = &testmatchers.HttpResponse{
		StatusCode: http.StatusNotFound,
		Body:       gstruct.Ignore(),
	}

	curlExitErrorCode = 28

	setup = base.SimpleTestCase{
		Manifests: []string{testdefaults.CurlPodManifest, setupManifest},
		Resources: []client.Object{testdefaults.CurlPod, exampleSvc, nginxPod, proxyDeployment, proxyService},
	}

	// test cases
	testCases = map[string]*base.TestCase{
		"TestValidListenerSet": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{validListenerSetManifest},
			},
		},
		"TestInvalidListenerSetNotAllowed": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{invalidListenerSetNotAllowedManifest},
			},
		},
		"TestInvalidListenerSetNonExistingGW": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{invalidListenerSetNonExistingGWManifest},
			},
		},
	}
)
