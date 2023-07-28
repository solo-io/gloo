package gateway_test

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/solo-projects/test/kube2e"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/waf"

	. "github.com/onsi/ginkgo/v2"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var _ = Describe("waf tests", func() {

	var (
		testContext *kube2e.TestContext
	)

	const (
		response200 = "HTTP/1.1 200 OK"
		response403 = "HTTP/1.1 403 Forbidden"
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

	checkConnection := func(status string) {
		curlOpts := testContext.DefaultCurlOptsBuilder().WithConnectionTimeout(10).WithVerbose(true).Build()
		testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, status, 1, time.Minute*5)
	}

	Context("Core Rule Set", func() {

		It("will return 200 on a standard request and no custom rules", func() {
			wafVhost := &waf.Settings{
				CoreRuleSet: &waf.CoreRuleSet{},
			}
			virtualHostPlugins := &gloov1.VirtualHostOptions{
				Waf: wafVhost,
			}

			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, &gloov1.RouteOptions{
					PrefixRewrite: &wrappers.StringValue{
						Value: "/",
					},
				}).WithVirtualHostOptions(virtualHostPlugins).Build()
			})
			checkConnection(response200)
		})

		It("will reject an http 1.0 request", func() {
			wafVhost := &waf.Settings{
				CoreRuleSet: &waf.CoreRuleSet{
					CustomSettingsType: &waf.CoreRuleSet_CustomSettingsString{
						CustomSettingsString: ruleStr,
					},
				},
			}

			virtualHostPlugins := &gloov1.VirtualHostOptions{
				Waf: wafVhost,
			}

			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				return helpers.BuilderFromVirtualService(service).WithRouteOptions(kube2e.DefaultRouteName, &gloov1.RouteOptions{
					PrefixRewrite: &wrappers.StringValue{
						Value: "/",
					},
				}).WithVirtualHostOptions(virtualHostPlugins).Build()
			})
			checkConnection(response403)
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
