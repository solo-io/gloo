package gloo_mtls_test

import (
	"context"
	"time"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"k8s.io/client-go/rest"
)

// This file is largely copied from test/kube2e/gateway/gateway_test.go (May 2020)

var _ = Describe("Installing gloo in gloo mtls mode", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		gatewayClient        v2.GatewayClient
		virtualServiceClient v1.VirtualServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		cache := kube.NewKubeCache(ctx)
		gatewayClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v2.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}

		gatewayClient, err = v2.NewGatewayClient(ctx, gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(ctx, virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		kube2e.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
		cancel()
	})

	It("can route request to upstream", func() {

		kube2e.WriteVirtualService(ctx, testHelper, virtualServiceClient, nil, nil, nil)

		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		Eventually(func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))

		gatewayPort := 80
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              kube2e.TestMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		}, kube2e.GetSimpleTestRunnerHttpResponse(), 1, time.Minute*5)
	})

	It("can recover from the ext-auth sidecar container being deleted", func() {
		podName, err := services.GetGatewayPodName()
		Expect(err).ToNot(HaveOccurred())
		// it takes 2 requests to kill the container un-gracefully
		_, err = services.KubectlOut("exec", "-n", testHelper.InstallNamespace, podName, "-c", "extauth", "--", "kill", "1")
		Expect(err).ToNot(HaveOccurred())
		_, err = services.KubectlOut("exec", "-n", testHelper.InstallNamespace, podName, "-c", "extauth", "--", "kill", "1")
		Expect(err).ToNot(HaveOccurred())

		Eventually(func() bool {
			ready, err := services.PodIsReady(ctx, testHelper.InstallNamespace, podName)
			if err != nil {
				return true
			}
			return ready
		}, "2s", "0.1s").Should(Equal(false))

		// after some time kubernetes will repost the container, and this should resolve with a ready status
		Eventually(func() bool {
			ready, err := services.PodIsReady(ctx, testHelper.InstallNamespace, podName)
			if err != nil {
				return false
			}
			return ready
		}, "15s", "1s").Should(Equal(true))
	})

})
