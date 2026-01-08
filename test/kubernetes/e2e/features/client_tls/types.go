package client_tls

import (
	_ "embed"
	"net/http"

	"github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/test/gomega/matchers"
)

//go:embed testdata/nginx-upstreams.yaml
var NginxUpstreamsYaml []byte

//go:embed testdata/nginx-annotated-services.yaml
var NginxAnnotatedServicesYaml []byte

//go:embed testdata/vs-targeting-upstream.yaml
var VSTargetingUpstreamYaml []byte

//go:embed testdata/vs-targeting-kube.yaml
var VSTargetingKubeYaml []byte

//go:embed testdata/vs-oneway-downstream-tls.yaml
var VSOnewayDownstreamTlsYaml []byte

var (
	vSTargetingUpstreamObject = func(ns string) *kubev1.VirtualService {
		return &kubev1.VirtualService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vs-targeting-upstream",
				Namespace: ns,
			},
		}
	}

	vSTargetingKubeObject = func(ns string) *kubev1.VirtualService {
		return &kubev1.VirtualService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vs-targeting-kube",
				Namespace: ns,
			},
		}
	}

	vSOnewayDownstreamTlsObject = func(ns string) *kubev1.VirtualService {
		return &kubev1.VirtualService{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "vs-oneway-downstream-tls",
				Namespace: ns,
			},
		}
	}

	expectedHealthyResponse = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gomega.ContainSubstring("Welcome to nginx!"),
	}

	expectedCertVerifyFailedResponse = &matchers.HttpResponse{
		StatusCode: http.StatusServiceUnavailable,
		Body:       gomega.ContainSubstring("CERTIFICATE_VERIFY_FAILED"),
	}

	nginxTlsSecret = func(ns string) *corev1.Secret {
		return &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nginx-tls",
				Namespace: ns,
			},
		}
	}

	coreSecretGVK = schema.GroupVersionKind{
		Version: "v1",
		Group:   "",
		Kind:    "Secret",
	}
)
