package gloo_mtls_test

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	v12 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	osskube2e "github.com/solo-io/gloo/test/kube2e"
)

// This file is largely copied from test/kube2e/gateway/gateway_test.go (May 2020)

var _ = Describe("Installing gloo in gloo mtls mode", func() {

	var (
		testContext *kube2e.TestContext
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	It("can route request to upstream", func() {
		testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
			return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, &v12.RouteOptions{
				PrefixRewrite: &wrappers.StringValue{
					Value: "/",
				},
			}).Build()
		})

		curlOpts := testContext.DefaultCurlOptsBuilder().WithConnectionTimeout(10).Build()
		testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, osskube2e.GetSimpleTestRunnerHttpResponse(), 1, time.Minute*5)
	})

	It("can recover from the ext-auth sidecar container being deleted", func() {
		podName, err := services.GetGatewayPodName()
		Expect(err).ToNot(HaveOccurred())
		// getting out of memory issue on CI containers, try to attempt multiple times
		numOfAttempts := 5
		// if we successfully kill the container, then the container will be in a bad state and is not ready
		result := eventuallyKillExtAuthSideCarUnGracefully(testContext.Ctx(), testContext.InstallNamespace(), podName, numOfAttempts)
		Expect(result).To(Equal(true))
		eventuallyPodWillBeNonReady(testContext.Ctx(), podName)
		// after some time kubernetes will repost the container, and this should resolve with a ready status
		eventuallyPodWillBeReady(testContext.Ctx(), podName)
	})

})

func killExtAuthSideCarContainer(namespace, podName string) bool {
	// it takes 2 requests to kill the container un-gracefully
	// it takes some time for the extauth container to be ready
	Eventually(func(g Gomega) {
		_, err := services.KubectlOut("exec", "-n", namespace, podName, "-c", "extauth", "--", "kill", "1")
		g.Expect(err).ToNot(HaveOccurred())
	}, "5s", "1s")
	_, err := services.KubectlOut("exec", "-n", namespace, podName, "-c", "extauth", "--", "kill", "1")
	return err == nil
}

func eventuallyKillExtAuthSideCarUnGracefully(ctx context.Context, namespace, podName string, attempts int) bool {
	for i := 0; i < attempts; i++ {
		result := killExtAuthSideCarContainer(namespace, podName)
		if result {
			return true
		} else {
			eventuallyPodWillBeNonReady(ctx, podName)
			// not sure if the container will fix itself... as it should
			eventuallyPodWillBeReady(ctx, podName)
		}
	}
	return false
}

func eventuallyPodWillBeReady(ctx context.Context, podName string) {
	eventuallyPodWillReachReadyState(ctx, namespace, podName, true)
}

func eventuallyPodWillBeNonReady(ctx context.Context, podName string) {
	eventuallyPodWillReachReadyState(ctx, namespace, podName, false)
}

func eventuallyPodWillReachReadyState(ctx context.Context, namespace, podName string, readyState bool) {
	Eventually(func() bool {
		ready, err := services.PodIsReady(ctx, namespace, podName)
		if err != nil {
			return false
		}
		return ready
	}, "15s", time.Millisecond*50).Should(Equal(readyState))
}
