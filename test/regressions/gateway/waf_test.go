package gateway_test

import (
	"context"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/waf"

	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v2"

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

var _ = Describe("waf tests", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc
		cfg    *rest.Config

		cache                kube.SharedCache
		gatewayClient        v2.GatewayClient
		virtualServiceClient v1.VirtualServiceClient
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

	})

	AfterEach(func() {
		cancel()
		err := virtualServiceClient.Delete(testHelper.InstallNamespace, "vs", clients.DeleteOpts{})
		Expect(err).NotTo(HaveOccurred())
	})

	waitForGateway := func() {
		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		EventuallyWithOffset(2, func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))
	}

	checkConnection := func(status string) {
		waitForGateway()

		gatewayPort := int(80)
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              testMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			Verbose:           true,
		}, status, 1, time.Minute*5)
	}

	Context("Core Rule Set", func() {

		It("will return 200 on a standard request and no custom rules", func() {

			wafVhost := &waf.Settings{
				CoreRuleSet: &waf.CoreRuleSet{},
			}

			virtualHostPlugins := &gloov1.VirtualHostPlugins{
				Waf: wafVhost,
			}

			writeVirtualService(ctx, virtualServiceClient, virtualHostPlugins, nil, nil)
			checkConnection("200")
		})

		It("will reject an http 1.0 request", func() {

			wafVhost := &waf.Settings{
				CoreRuleSet: &waf.CoreRuleSet{
					CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
						CustomSettingsString: ruleStr,
					},
				},
			}

			virtualHostPlugins := &gloov1.VirtualHostPlugins{
				Waf: wafVhost,
			}

			writeVirtualService(ctx, virtualServiceClient, virtualHostPlugins, nil, nil)
			checkConnection("403")
		})

	})
})

const (
	ruleStr = `
	# default rules section
	SecRuleEngine On
	SecRequestBodyAccess On

	# CRS section
	# Will block by default
	SecDefaultAction "phase:1,log,auditlog,deny,status:403"
	SecDefaultAction "phase:2,log,auditlog,deny,status:403"

	# only allow http2 connections
	SecAction \
     "id:900230,\
      phase:1,\
      nolog,\
      pass,\
      t:none,\
      setvar:'tx.allowed_http_versions=HTTP/2 HTTP/2.0'"

    SecAction \
     "id:900990,\
      phase:1,\
      nolog,\
      pass,\
      t:none,\
      setvar:tx.crs_setup_version=310"
`
)
