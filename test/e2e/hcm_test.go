package e2e_test

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"

	"github.com/solo-io/gloo/test/e2e"
)

var _ = Describe("HttpConnectionManagerSettings", func() {

	var (
		testContext *e2e.TestContext
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

	Context("UuidRequestIdConfig", func() {

		BeforeEach(func() {
			gw := gatewaydefaults.DefaultGateway(writeNamespace)
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				HttpConnectionManagerSettings: &hcm.HttpConnectionManagerSettings{
					UuidRequestIdConfig: &hcm.HttpConnectionManagerSettings_UuidRequestIdConfigSettings{
						PackTraceReason:              &wrappers.BoolValue{Value: true},
						UseRequestIdForTraceSampling: &wrappers.BoolValue{Value: true},
					},
				},
			}

			testContext.ResourcesToCreate().Gateways = gatewayv1.GatewayList{
				gw,
			}
		})

		It("does not lead to error in logs", func() {
			Consistently(func(g Gomega) {
				logs, err := testContext.EnvoyInstance().Logs()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(logs).NotTo(ContainSubstring("Didn't find a registered implementation"))
			}).Should(Succeed())
		})
	})
})
