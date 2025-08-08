package client_tls

import (
	_ "embed"
	"net/http"

	"github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	expectedHealthyResponse = &matchers.HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gomega.ContainSubstring("Welcome to nginx!"),
	}

	expectedCertVerifyFailedResponse = &matchers.HttpResponse{
		StatusCode: http.StatusServiceUnavailable,
		Body:       gomega.ContainSubstring("CERTIFICATE_VERIFY_FAILED"),
	}
)
