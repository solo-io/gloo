package virtualhost_options

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
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

	// ports used by the listener set
	lsPort1 = 8085
	lsPort2 = 8086
)

var (
	setup = func(ti *e2e.TestInstallation) base.SimpleTestCase {
		return base.SimpleTestCase{
			Manifests: setupManifests(ti),
			Resources: []client.Object{proxy1Service, proxy1Deployment, proxy2Service, proxy2Deployment, nginxPod, defaults.CurlPod},
		}
	}

	testCases = map[string]*base.TestCase{
		"TestConfirmSetup": {},
		"TestConfigureVirtualHostOptions": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoRemoveXBar},
			},
		},
		"TestConfigureVirtualHostOptionsMultipleTargetRefs": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoMultipleTargetRefs},
			},
		},
		"TestConfigureVirtualHostOptionsListenerSetTargetRef": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoListenerSetTargetRef},
			},
		},
		"TestConfigureVirtualHostOptionsListenerSetSectionedTargetRef": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoListenerSetSectionedTargetRef},
			},
		},
		"TestConfigureVirtualHostOptionsWithConflictingVHO": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoSectionAddXFoo, manifestVhoGwAddXFoo, manifestVhoListenerSetTargetRef, manifestVhoListenerSetSectionedTargetRef},
			},
		},
		"TestConfigureInvalidVirtualHostOptions": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoRemoveXBar, manifestVhoWebhookReject},
			},
		},
		"TestConfigureVirtualHostOptionsWithSectionNameManualSetup": {
			// "Manual setup" so let the test handle it.
		},
		"TestMultipleVirtualHostOptionsSetup": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoRemoveXBar, manifestVhoRemoveXBaz},
			},
		},
		"TestConfigureVirtualHostOptionsWarningMultipleGatewaysSetup": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoMultipleGatewayWarnings},
			},
		},
		"TestDeletingConflictingVirtualHostOptions": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoRemoveXBaz},
			},
		},
		"TestOptionsMerge": {
			SimpleTestCase: base.SimpleTestCase{
				Manifests: []string{manifestVhoRemoveXBar, manifestVhoMergeRemoveXBaz},
			},
		},
	}
	commonSetupManifests = []string{
		filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml"),
		e2edefaults.CurlPodManifest,
	}

	manifestGw1NoListenerSet                 = filepath.Join(util.MustGetThisDir(), "testdata", "gw1-no-listenerset.yaml")
	manifestGw1ListenerSet                   = filepath.Join(util.MustGetThisDir(), "testdata", "gw1-listenerset.yaml")
	manifestListenerSetup                    = filepath.Join(util.MustGetThisDir(), "testdata", "listenerset.yaml")
	manifestVhoRemoveXBar                    = filepath.Join(util.MustGetThisDir(), "testdata", "vho-remove-x-bar.yaml")
	manifestVhoSectionAddXFoo                = filepath.Join(util.MustGetThisDir(), "testdata", "vho-section-add-x-foo.yaml")
	manifestVhoGwAddXFoo                     = filepath.Join(util.MustGetThisDir(), "testdata", "vho-gw-add-x-foo.yaml")
	manifestVhoRemoveXBaz                    = filepath.Join(util.MustGetThisDir(), "testdata", "vho-remove-x-baz.yaml")
	manifestVhoWebhookReject                 = filepath.Join(util.MustGetThisDir(), "testdata", "vho-webhook-reject.yaml")
	manifestVhoMergeRemoveXBaz               = filepath.Join(util.MustGetThisDir(), "testdata", "vho-merge-remove-x-baz.yaml")
	manifestVhoMultipleTargetRefs            = filepath.Join(util.MustGetThisDir(), "testdata", "vho-multiple-target-refs.yaml")
	manifestVhoListenerSetTargetRef          = filepath.Join(util.MustGetThisDir(), "testdata", "vho-listener-set-target-ref.yaml")
	manifestVhoListenerSetSectionedTargetRef = filepath.Join(util.MustGetThisDir(), "testdata", "vho-listener-set-sectioned-target-ref.yaml")
	manifestVhoMultipleGatewayWarnings       = filepath.Join(util.MustGetThisDir(), "testdata", "vho-multiple-gateway-warnings.yaml")

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

	// VHO to add a x-foo header
	vhoRemoveXBar = metav1.ObjectMeta{
		Name:      "remove-x-bar-header",
		Namespace: "default",
	}
	// VHO to remove a x-baz header
	vhoRemoveXBaz = metav1.ObjectMeta{
		Name:      "remove-x-baz-header",
		Namespace: "default",
	}
	// VHO to remove a x-baz header
	vhoMergeRemoveXBaz = metav1.ObjectMeta{
		Name:      "remove-x-baz-merge",
		Namespace: "default",
	}
	// VHO to add a x-foo header in a section
	vhoSectionAddXFoo = metav1.ObjectMeta{
		Name:      "add-x-foo-header-section",
		Namespace: "default",
	}
	// VHO to add a x-foo header to a gateway
	vhoGwAddXFoo = metav1.ObjectMeta{
		Name:      "add-x-foo-header-gw",
		Namespace: "default",
	}
	// VHO that should be rejected by the validating webhook
	vhoWebhookReject = metav1.ObjectMeta{
		Name:      "bad-retries",
		Namespace: "default",
	}
	// VHO to add a x-foo header with multiple target refs
	vhoMultipleTargetRefs = metav1.ObjectMeta{
		Name:      "add-x-foo-header-multiple-target-refs",
		Namespace: "default",
	}
	// VHO to add a x-foo header with multiple target refs
	vhoListenerSetTargetRef = metav1.ObjectMeta{
		Name:      "add-x-foo-header-listener-set-target-ref",
		Namespace: "default",
	}
	// VHO to add a x-foo header with multiple target refs
	vhoListenerSetSectionedTargetRef = metav1.ObjectMeta{
		Name:      "add-x-foo-header-listener-set-sectioned-target-ref",
		Namespace: "default",
	}

	// Expects a 200 response with x-bar and x-baz headers
	defaultResponseGw1 = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaderKeys([]string{"x-foo"})),
			matchers.ContainHeaderKeys([]string{"x-bar", "x-baz"}),
		),
		Body: gstruct.Ignore(),
	}

	defaultResponseGw2 = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaderKeys([]string{"x-foo"})),
			matchers.ContainHeaderKeys([]string{"x-bar-2", "x-baz-2"}),
		),
		Body: gstruct.Ignore(),
	}

	// Expects default response with no x-bar header
	expectedResponseWithoutXBar = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaderKeys([]string{"x-bar"})),
			matchers.ContainHeaderKeys([]string{"x-baz"}),
		),
		Body: gstruct.Ignore(),
	}

	// Expects default response with no x-baz header
	expectedResponseWithoutXBaz = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			matchers.ContainHeaderKeys([]string{"x-bar"}),
			gomega.Not(matchers.ContainHeaderKeys([]string{"x-baz"})),
		),
	}

	// Expects default response with x-foo header
	expectedResponseWithXFooBarBaz = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]interface{}{
			"x-foo": gomega.Equal("foo"),
		},
		// Make sure the x-bar isn't being removed as a function of the unwanted VHO
		Custom: gomega.And(
			matchers.ContainHeaderKeys([]string{"x-foo", "x-bar", "x-baz"}),
		),
		Body: gstruct.Ignore(),
	}

	expectedResponseWithoutXFoo = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaderKeys([]string{"x-foo"})),
		),
		Body: gstruct.Ignore(),
	}

	expectedResponseWithXFoo = func(val string) *matchers.HttpResponse {
		return &matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]interface{}{
				"x-foo": gomega.Equal(val),
			},
			Body: gstruct.Ignore(),
		}
	}
)

func setupManifests(ti *e2e.TestInstallation) []string {
	manifests := commonSetupManifests
	if listenerset.RequiredCrdExists(ti) {
		manifests = append(manifests, manifestGw1ListenerSet, manifestListenerSetup)
	} else {
		manifests = append(manifests, manifestGw1NoListenerSet)
	}
	return manifests
}
