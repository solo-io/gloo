package gateway_test

import (
	"context"
	"fmt"
	"time"

	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/dlp"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"k8s.io/client-go/rest"
)

var _ = Describe("dlp tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		cache                kube.SharedCache
		gatewayClient        v2.GatewayClient
		virtualServiceClient v1.VirtualServiceClient

		httpEcho helper.TestRunner
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		cache = kube.NewKubeCache(ctx)
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

		httpEcho, err = helper.NewEchoHttp(testHelper.InstallNamespace)
		Expect(err).NotTo(HaveOccurred())

		err = httpEcho.Deploy(2 * time.Minute)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		cancel()
		deleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
		err := httpEcho.Terminate()
		Expect(err).NotTo(HaveOccurred())
	})

	waitForGateway := func() {
		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		EventuallyWithOffset(2, func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))
	}

	checkConnection := func(body string) {
		waitForGateway()

		gatewayPort := int(80)
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              testMatcherPrefix,
			Method:            "GET",
			Headers:           map[string]string{"hello": "world"},
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			Verbose:           true,
		}, body, 1, time.Minute*5)
	}

	Context("data loss prevention", func() {

		It("will mask regex matches", func() {

			dlpVhost := &dlp.Config{
				Actions: []*dlp.Action{
					{
						ActionType: dlp.Action_CUSTOM,
						CustomAction: &dlp.CustomAction{
							Name:     "test",
							Regex:    []string{"hello", "world"},
							MaskChar: "Y",
							Percent: &envoy_type.Percent{
								Value: 60,
							},
						},
					},
				},
			}

			virtualHostPlugins := &gloov1.VirtualHostPlugins{
				Dlp: dlpVhost,
			}

			httpEchoRef := &core.ResourceRef{
				Namespace: testHelper.InstallNamespace,
				Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, helper.HttpEchoName, helper.HttpEchoPort),
			}
			writeCustomVirtualService(ctx, virtualServiceClient, virtualHostPlugins, nil, nil, httpEchoRef)
			checkConnection(`"YYYlo":"YYYld"`)
		})
	})
})
