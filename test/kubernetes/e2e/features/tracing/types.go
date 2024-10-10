package tracing

import (
	"path/filepath"

	"github.com/solo-io/skv2/codegen/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	setupOtelcolManifest = filepath.Join(util.MustGetThisDir(), "testdata", "setup-otelcol.yaml")
	setupEchoServerManifest = filepath.Join(util.MustGetThisDir(), "testdata", "setup-echo-server.yaml")
	tracingConfigManifest = filepath.Join(util.MustGetThisDir(), "testdata", "tracing.yaml")

	otelcolService = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{ Name: "otel-collector", Namespace: "default" },
	}
	otelcolPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{ Name: "otel-collector", Namespace: "default" },
	}
	otelcolSelector = metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=opentelemetry",
	}
	echoServerPod = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{ Name: "echo-server", Namespace: "default" },
	}
	echoServerSelector = metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=echo-server",
	}
	proxyService = &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{ Name: "gloo-proxy-gw", Namespace: "default" },
	}
)
