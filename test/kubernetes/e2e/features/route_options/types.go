package route_options

import (
	"net/http"
	"path/filepath"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/skv2/codegen/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	setupManifest                        = filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml")
	basicRtoManifest                     = filepath.Join(util.MustGetThisDir(), "testdata", "basic-rto.yaml")
	basicRtoTargetRefManifest            = filepath.Join(util.MustGetThisDir(), "testdata", "basic-rto-targetref.yaml")
	basicRtoMultipleTargetRefManifest    = filepath.Join(util.MustGetThisDir(), "testdata", "basic-rto-multiple-targetref.yaml")
	extraRtoManifest                     = filepath.Join(util.MustGetThisDir(), "testdata", "extra-rto.yaml")
	extraRtoTargetRefManifest            = filepath.Join(util.MustGetThisDir(), "testdata", "extra-rto-targetref.yaml")
	badRtoManifest                       = filepath.Join(util.MustGetThisDir(), "testdata", "bad-rto.yaml")
	badRtoTargetRefManifest              = filepath.Join(util.MustGetThisDir(), "testdata", "bad-rto-targetref.yaml")
	httproute1Manifest                   = filepath.Join(util.MustGetThisDir(), "testdata", "httproute1.yaml")
	httproute1ExtensionManifest          = filepath.Join(util.MustGetThisDir(), "testdata", "httproute1-extension.yaml")
	httproute1BadExtensionManifest       = filepath.Join(util.MustGetThisDir(), "testdata", "httproute1-bad-extension.yaml")
	httproute1MultipleExtensionsManifest = filepath.Join(util.MustGetThisDir(), "testdata", "httproute1-multiple-extensions.yaml")
	httproute2Manifest                   = filepath.Join(util.MustGetThisDir(), "testdata", "httproute2.yaml")
	mergeManifest                        = filepath.Join(util.MustGetThisDir(), "testdata", "merge.yaml")

	// When we apply the fault injection manifest files, we expect resources to be created with this metadata
	proxyService = &corev1.Service{
		ObjectMeta: objectMetaInDefault("gloo-proxy-gw"),
	}
	proxyDeployment = &appsv1.Deployment{
		ObjectMeta: objectMetaInDefault("gloo-proxy-gw"),
	}
	nginxPod = &corev1.Pod{
		ObjectMeta: objectMetaInDefault("nginx"),
	}
	exampleSvc = &corev1.Service{
		ObjectMeta: objectMetaInDefault("example-svc"),
	}

	// RouteOption resources to be created
	basicRtoMeta          = objectMetaInDefault("basic-rto")
	basicRtoTargetRefMeta = objectMetaInDefault("basic-rto-targetref")
	extraRtoMeta          = objectMetaInDefault("extra-rto")
	extraRtoTargetRefMeta = objectMetaInDefault("extra-rto-targetref")
	badRtoMeta            = objectMetaInDefault("bad-rto")
	badRtoTargetRefMeta   = objectMetaInDefault("bad-rto-targetref")
	extref1RtoMeta        = objectMetaInDefault("extref1")
	extref2RtoMeta        = objectMetaInDefault("extref2")
	target1RtoMeta        = objectMetaInDefault("target-1")
	target2RtoMeta        = objectMetaInDefault("target-2")

	// Expected response matchers for various route options applied
	expectedResponseWithBasicHeader          = expectedResponseWithFooHeader("basic-rto")
	expectedResponseWithBasicTargetRefHeader = expectedResponseWithFooHeader("basic-rto-targetref")
	expectedResponseWithExtraHeader          = expectedResponseWithFooHeader("extra-rto")
	expectedResponseWithExtraTargetRefHeader = expectedResponseWithFooHeader("extra-rto-targetref")

	expectedResponseWithoutFooHeader = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Custom: gomega.And(
			gomega.Not(matchers.ContainHeaderKeys([]string{"foo"})),
		),
		Body: gstruct.Ignore(),
	}
)

func objectMetaInDefault(name string) metav1.ObjectMeta {
	return objectMeta(name, "default")
}

func objectMeta(name, namespace string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
	}
}

func expectedResponseWithFooHeader(value string) *matchers.HttpResponse {
	return &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]interface{}{
			"foo": gomega.Equal(value),
		},
		Body: gstruct.Ignore(),
	}
}
