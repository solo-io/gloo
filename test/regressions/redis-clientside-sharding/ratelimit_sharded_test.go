package clientside_sharding_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/test/regressions"

	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"

	"github.com/solo-io/go-utils/kubeutils"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"

	"k8s.io/client-go/rest"
)

var _ = Describe("RateLimit tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		cache                 kube.SharedCache
		gatewayClient         v2.GatewayClient
		virtualServiceClient  v1.VirtualServiceClient
		settingsClient        gloov1.SettingsClient
		origSettings          *gloov1.Settings // used to capture & restore initial Settings so each test can modify them
		uniqueDescriptorValue string
	)

	const (
		response200 = "HTTP/1.1 200 OK"
		response401 = "HTTP/1.1 401 Unauthorized"
		response429 = "HTTP/1.1 429 Too Many Requests"
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())
		cache = kube.NewKubeCache(ctx)
		settingsClientFactory := &factory.KubeResourceClientFactory{
			Crd:         gloov1.SettingsCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		settingsClient, err = gloov1.NewSettingsClient(settingsClientFactory)
		Expect(err).NotTo(HaveOccurred())
		if uniqueDescriptorValue == "" {
			uniqueDescriptorValue = "sharding-test-limit"
		}
		uniqueDescriptorValue = uniqueDescriptorValue + "1"

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
		gatewayClient, err = v2.NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		var err error
		origSettings, err = settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read initial settings")
	})

	AfterEach(func() {
		regressions.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})

		currentSettings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")

		if origSettings.Metadata.ResourceVersion != currentSettings.Metadata.ResourceVersion {
			origSettings.Metadata.ResourceVersion = currentSettings.Metadata.ResourceVersion // so we can overwrite settings
			_, err = settingsClient.Write(origSettings, clients.WriteOpts{Ctx: ctx, OverwriteExisting: true})
			Expect(err).ToNot(HaveOccurred())
		}
		cancel()
	})

	waitForGateway := func() {
		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		EventuallyWithOffset(2, func() (*v2.Gateway, error) {
			g, err := gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
			status, ok := g.Status.SubresourceStatuses["*v1.Proxy.gloo-system.gateway-proxy"]
			if !ok {
				return nil, fmt.Errorf("gateway proxy not yet ready")
			}
			if g.Status.State != core.Status_Accepted {
				return nil, fmt.Errorf("gateway resource not yet accepted")
			}
			if status.State != core.Status_Accepted {
				return nil, fmt.Errorf("gateway proxy resource not yet accepted")
			}
			return g, err
		}, "15s", "0.5s").Should(Not(BeNil()))
	}

	requestOnLimitedPath := func() (string, error) {
		res, err := testHelper.Curl(helper.CurlOpts{
			Protocol:          "http",
			Path:              regressions.TestMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              80,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			Verbose:           true,
		})
		if err != nil {
			return "", err
		}
		return res, nil
	}

	Context("rate limit at specific rate, even with scaled redis", func() {

		var (
			ingressRateLimit = &ratelimit.IngressRateLimit{
				AnonymousLimits: &rlv1alpha1.RateLimit{
					RequestsPerUnit: 10,
					Unit:            rlv1alpha1.RateLimit_DAY,
				},
			}
			virtualHostPlugins = &gloov1.VirtualHostOptions{
				RatelimitBasic: ingressRateLimit,
			}
		)

		It("can rate limit to upstream", func() {
			regressions.WriteVirtualService(ctx, testHelper, virtualServiceClient, virtualHostPlugins, nil, nil)
			waitForGateway()

			// If we're about to roll over to the next day and
			// risk resetting the daily rate-limit mid-test,
			// just wait a minute before starting.
			if time.Now().Format("1504") == "2359" {
				fmt.Fprintln(GinkgoWriter, "the current time is 11.59pm, there is a risk of the daily rate limit resetting mid-test, so we will wait 1 minute before starting")
				time.Sleep(1 * time.Minute)
			}

			// Wait for the service to be happily responding
			Eventually(func() (string, error) {
				return requestOnLimitedPath()
			}, "15s", "0.5s").Should(ContainSubstring(response200))

			for i := 0; i < 9; i++ {
				// First 10 (next 9) requests should work (all should go to the same Redis)
				res, err := requestOnLimitedPath()
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(ContainSubstring(response200), "request #"+strconv.Itoa(i+1)+" should not be rate limited")
			}

			// After we've hit our limit, we should consistently be rate limited:
			Consistently(func() error {
				res, err := requestOnLimitedPath()
				if err != nil {
					return err
				}
				if !strings.Contains(res, response429) {
					return fmt.Errorf("response was not rate limited, expected %v to be found within %v", response429, res)
				}
				return nil
			}, "5s", ".1s").Should(BeNil())

		})
	})
})
