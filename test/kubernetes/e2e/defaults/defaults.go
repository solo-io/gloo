package defaults

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/skv2/codegen/util"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"

	_ "embed"
)

//go:embed testdata/curl_pod.yaml
var CurlPodYaml []byte

//go:embed testdata/http_echo.yaml
var HttpEchoPodYaml []byte

//go:embed testdata/tcp_echo.yaml
var TcpEchoPodYaml []byte

//go:embed testdata/nginx_pod.yaml
var NginxPodYaml []byte

//go:embed testdata/httpbin.yaml
var HttpbinYaml []byte

var (
	CurlPodExecOpt = kubectl.PodExecOptions{
		Name:      "curl",
		Namespace: "curl",
		Container: "curl",
	}

	CurlPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "curl",
			Namespace: "curl",
		},
	}

	CurlPodManifest = filepath.Join(util.MustGetThisDir(), "testdata", "curl_pod.yaml")

	CurlPodLabelSelector = "app.kubernetes.io/name=curl"

	HttpEchoPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http-echo",
			Namespace: "http-echo",
		},
	}

	HttpEchoPodManifest = filepath.Join(util.MustGetThisDir(), "testdata", "http_echo.yaml")

	TcpEchoPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "tcp-echo",
			Namespace: "tcp-echo",
		},
	}

	TcpEchoPodManifest = filepath.Join(util.MustGetThisDir(), "testdata", "tcp_echo.yaml")

	NginxPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "nginx",
		},
	}

	NginxSvc = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "nginx",
		},
	}

	NginxPodManifest = filepath.Join(util.MustGetThisDir(), "testdata", "nginx_pod.yaml")
)
