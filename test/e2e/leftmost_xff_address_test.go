package e2e_test

import (
	"fmt"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-projects/test/e2e"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

const (
	downstreamRemoteAddressFormat = "[%DOWNSTREAM_REMOTE_ADDRESS_WITHOUT_PORT%]"
)

var _ = Describe("Leftmost x-forwarded-for address ", func() {

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

	When("Leftmost XFF Address = true", func() {

		BeforeEach(func() {
			gw := defaults.DefaultGateway(e2e.WriteNamespace)
			gw.Options = &gloov1.ListenerOptions{
				AccessLoggingService: &als.AccessLoggingService{
					AccessLog: []*als.AccessLog{
						{
							OutputDestination: &als.AccessLog_FileSink{
								FileSink: &als.FileSink{
									Path: "/dev/stdout",
									OutputFormat: &als.FileSink_StringFormat{
										StringFormat: downstreamRemoteAddressFormat,
									},
								},
							},
						},
					},
				},
			}
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				LeftmostXffAddress: &wrappers.BoolValue{
					Value: true,
				},
			}

			testContext.ResourcesToCreate().Gateways = v1.GatewayList{
				gw,
			}
		})

		It("uses leftmost xff address as downstream remote address", func() {
			address := "192.168.2.1"

			requestBuilder := testContext.GetHttpRequestBuilder().
				WithHeader("x-forwarded-for", fmt.Sprintf("%s,192.123.3.1,192.123.3.2", address))

			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(requestBuilder.Build())).To(matchers.HaveOkResponse(),
					"request should succeed")
				g.Expect(testContext.EnvoyInstance().Logs()).To(ContainSubstring(fmt.Sprintf("[%s]", address)),
					"downstream remote address should be set to leftmost x-forwarded-for address")
			})
		})

	})

	When("Leftmost XFF Address = false", func() {

		BeforeEach(func() {
			gw := defaults.DefaultGateway(e2e.WriteNamespace)
			gw.Options = &gloov1.ListenerOptions{
				AccessLoggingService: &als.AccessLoggingService{
					AccessLog: []*als.AccessLog{
						{
							OutputDestination: &als.AccessLog_FileSink{
								FileSink: &als.FileSink{
									Path: "/dev/stdout",
									OutputFormat: &als.FileSink_StringFormat{
										StringFormat: downstreamRemoteAddressFormat,
									},
								},
							},
						},
					},
				},
			}
			gw.GetHttpGateway().Options = &gloov1.HttpListenerOptions{
				LeftmostXffAddress: &wrappers.BoolValue{
					Value: false,
				},
			}

			testContext.ResourcesToCreate().Gateways = v1.GatewayList{
				gw,
			}
		})

		It("uses rightmost xff address as downstream remote address", func() {
			address := "192.168.2.1"

			requestBuilder := testContext.GetHttpRequestBuilder().
				WithHeader("x-forwarded-for", fmt.Sprintf("192.123.3.1,192.123.3.2,%s", address))

			Eventually(func(g Gomega) {
				g.Expect(testutils.DefaultHttpClient.Do(requestBuilder.Build())).To(matchers.HaveOkResponse(),
					"request should succeed")
				g.Expect(testContext.EnvoyInstance().Logs()).To(ContainSubstring(fmt.Sprintf("[%s]", address)),
					"downstream remote address should be set to rightmost x-forwarded-for address")
			})

		})

	})

})
