package virtualhost_options

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/test/gomega/matchers"
	e2edefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	setupManifests = []string{
		filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml"),
		e2edefaults.CurlPodManifest,
	}

	manifestVhoRemoveXBar              = filepath.Join(util.MustGetThisDir(), "testdata", "vho-remove-x-bar.yaml")
	manifestVhoSectionAddXFoo          = filepath.Join(util.MustGetThisDir(), "testdata", "vho-section-add-x-foo.yaml")
	manifestVhoRemoveXBaz              = filepath.Join(util.MustGetThisDir(), "testdata", "vho-remove-x-baz.yaml")
	manifestVhoWebhookReject           = filepath.Join(util.MustGetThisDir(), "testdata", "vho-webhook-reject.yaml")
	manifestVhoMergeRemoveXBaz         = filepath.Join(util.MustGetThisDir(), "testdata", "vho-merge-remove-x-baz.yaml")
	manifestVhoMultipleTargetRefs      = filepath.Join(util.MustGetThisDir(), "testdata", "vho-multiple-target-refs.yaml")
	manifestVhoMultipleGatewayWarnings = filepath.Join(util.MustGetThisDir(), "testdata", "vho-multiple-gateway-warnings.yaml")

	// When we apply the setup file, we expect resources to be created with this metadata
	// When we apply the setup file, we expect resources to be created with this metadata
	glooProxyObjectMeta1 = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw-1",
		Namespace: "default",
	}
	proxyService1     = &corev1.Service{ObjectMeta: glooProxyObjectMeta1}
	proxyService1Fqdn = kubeutils.ServiceFQDN(proxyService1.ObjectMeta)
	proxyDeployment1  = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gloo-proxy-gw-1",
			Namespace: "default",
		},
	}
	glooProxyObjectMeta2 = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw-2",
		Namespace: "default",
	}
	proxyService2     = &corev1.Service{ObjectMeta: glooProxyObjectMeta2}
	proxyService2Fqdn = kubeutils.ServiceFQDN(proxyService2.ObjectMeta)
	proxyDeployment2  = &appsv1.Deployment{
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
		Name:      "add-x-foo-header",
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

	// Expects default response with x-foo header
	expectedResponseWithXFoo = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]interface{}{
			"x-foo": gomega.Equal("foo"),
		},
		Body: gstruct.Ignore(),
	}

	expectedResponseWithoutXFoo = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaderKeys([]string{"x-foo"})),
		),
		Body: gstruct.Ignore(),
	}

	// Port numbers and mappings match the ports in the setup.yaml file
	gw1port1 = 8080
	gw1port2 = 8081
	// port 8082 is used by envoy's readiness probe
	gw2port1 = 8083
	gw2port2 = 8084

	// The keys in this map are the FQDNs of the gateway services
	// The values are the ports on which the gateway services are listening
	gatewayListenerPorts = map[string][]int{
		proxyService1Fqdn: {gw1port1, gw1port2},
		proxyService2Fqdn: {gw2port1, gw2port2},
	}
)
