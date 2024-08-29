package utils

import (
	"context"
	"strings"
	"time"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	e2edefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/testutils/assertions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	GatewayProxyName = "gateway-proxy"
)

func generateCurlOpts(ctx context.Context, host string) []curl.Option {
	var curlOpts = []curl.Option{
		curl.WithHost(kubeutils.ServiceFQDN(metav1.ObjectMeta{Name: GatewayProxyName, Namespace: ctx.Value("namespace").(string)})),
		curl.WithPort(80),
		curl.Silent(),
	}

	path := ""
	parts := strings.SplitN(host, "/", 2)
	host = parts[0]
	if len(parts) > 1 {
		path = parts[1] + path
	}

	parts = strings.SplitN(host, "@", 2)
	if len(parts) > 1 {
		host = parts[1]
		auth := strings.Split(parts[0], ":")
		curlOpts = append(curlOpts, curl.WithBasicAuth(auth[0], auth[1]))
	}

	return append(curlOpts,
		curl.WithHostHeader(host),
		curl.WithPath(path),
	)
}

func generateCurlOptsWithHeaders(ctx context.Context, host string, headers map[string]string) []curl.Option {
	curlOpts := generateCurlOpts(ctx, host)
	for k, v := range headers {
		curlOpts = append(curlOpts, curl.WithHeader(k, v))
	}
	return curlOpts
}

func CurlConsistentlyRespondsWithStatus(ctx context.Context, provider *assertions.Provider, host string, status int) {
	curlOptsHeader := generateCurlOpts(ctx, host)

	provider.AssertEventuallyConsistentCurlResponse(
		ctx,
		e2edefaults.CurlPodExecOpt,
		curlOptsHeader,
		&matchers.HttpResponse{StatusCode: status},
		10*time.Second,
		100*time.Millisecond,
	)
}

func CurlEventuallyRespondsWithStatus(ctx context.Context, provider *assertions.Provider, host string, status int) {
	curlOptsHeader := generateCurlOpts(ctx, host)

	provider.AssertEventualCurlResponse(
		ctx,
		e2edefaults.CurlPodExecOpt,
		curlOptsHeader,
		&matchers.HttpResponse{StatusCode: status},
	)
}

func CurlRespondsWithStatus(ctx context.Context, provider *assertions.Provider, host string, status int) {
	curlOptsHeader := generateCurlOpts(ctx, host)

	provider.AssertCurlResponse(
		ctx,
		e2edefaults.CurlPodExecOpt,
		curlOptsHeader,
		&matchers.HttpResponse{StatusCode: status},
	)
}

func CurlWithHeadersConsistentlyRespondsWithStatus(ctx context.Context, provider *assertions.Provider, host string, headers map[string]string, status int) {
	curlOptsHeader := generateCurlOptsWithHeaders(ctx, host, headers)

	provider.AssertEventuallyConsistentCurlResponse(
		ctx,
		e2edefaults.CurlPodExecOpt,
		curlOptsHeader,
		&matchers.HttpResponse{StatusCode: status},
		10*time.Second,
		100*time.Millisecond,
	)
}

func CurlWithHeadersEventuallyRespondsWithStatus(ctx context.Context, provider *assertions.Provider, host string, headers map[string]string, status int) {
	curlOptsHeader := generateCurlOptsWithHeaders(ctx, host, headers)

	provider.AssertEventualCurlResponse(
		ctx,
		e2edefaults.CurlPodExecOpt,
		curlOptsHeader,
		&matchers.HttpResponse{StatusCode: status},
	)
}

func CurlWithHeadersRespondsWithStatus(ctx context.Context, provider *assertions.Provider, host string, headers map[string]string, status int) {
	curlOptsHeader := generateCurlOptsWithHeaders(ctx, host, headers)

	provider.AssertCurlResponse(
		ctx,
		e2edefaults.CurlPodExecOpt,
		curlOptsHeader,
		&matchers.HttpResponse{StatusCode: status},
	)
}
