package e2e_test

import (
	"net/http"

	"github.com/solo-io/solo-projects/test/e2e"
	"github.com/solo-io/solo-projects/test/services/tap_server"

	"github.com/solo-io/gloo/test/testutils"

	"github.com/solo-io/gloo/test/gomega/matchers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tap_service "github.com/solo-io/tap-extension-examples/pkg/tap_service"
)

var _ = Describe("Tap filter E2E Test", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	sendTestRequest := func(request *http.Request) *tap_service.TapRequest {
		Eventually(func(g Gomega) {
			g.Expect(testutils.DefaultHttpClient.Do(request)).Should(matchers.HaveOkResponse())
		}, "30s", ".5s").Should(Succeed(), "Simple GET request should return 200")
		var tapData *tap_service.TapRequest
		Eventually(func(g Gomega) {
			// try to get the logs from the tap server
			tapData = testContext.TapServerInstance().Logs()
			g.Expect(tapData).NotTo(BeNil())
		}, "5s", ".5s").Should(Succeed())
		return tapData
	}

	Context("without dlp", func() {
		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				TapServer: &tap_server.InstanceConfig{
					EnableDataScrubbing: false,
				},
			})
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

		Context("tap functionality", func() {

			// Test whether the tap server scrubs sensitive data from the trace
			// information. Note that this is happening in the tap server itself,
			// and not in the envoy tap filter, so do not take away the impression
			// that we provide this functionality to users!!!
			It("does not scrub sensitive data from headers", func() {
				ssnHeader := "x-some-ssn"
				ssnInput := "123456789"
				ssnExpected := ssnInput

				creditCardHeader := "x-some-credit-card"
				creditCardInput := "5151-5151-5151-5151"
				creditCardExpected := creditCardInput

				request := testContext.GetHttpRequestBuilder().
					WithHeader("X-some-ssn", ssnInput).
					WithHeader("X-some-credit-card", creditCardInput).
					Build()
				tapData := sendTestRequest(request)
				traceData := tapData.TraceData.GetHttpBufferedTrace()

				response := &http.Response{
					Header: make(http.Header),
				}
				for _, header := range traceData.GetRequest().GetHeaders() {
					response.Header.Add(header.Key, header.Value)
				}
				Expect(response).To(matchers.ConsistOfHeaders(
					http.Header{
						ssnHeader:        []string{ssnExpected},
						creditCardHeader: []string{creditCardExpected},
					}), "should have headers with sensitive information removed")
			})

			It("does not scrub sensitive data from the body", func() {
				bodyInput := "sample SSN: 012-34-5678 sample credit card: 5151-5151-5151-5151"
				bodyExpected := bodyInput
				request := testContext.GetHttpRequestBuilder().WithBody(bodyInput).Build()
				tapData := sendTestRequest(request)
				traceData := tapData.TraceData.GetHttpBufferedTrace()
				Expect(traceData).NotTo(BeNil())
				Expect(string(traceData.GetRequest().GetBody().GetAsBytes())).To(Equal(bodyExpected))
				// since the upstream is an echo server, we expect the body to be
				// the same. this assertion isn't strictly necessary - it doesn't
				// test any specific behaviour - but it doesn't harm anything, so
				// i'll just leave it in
				Expect(string(traceData.GetResponse().GetBody().GetAsBytes())).To(Equal(bodyExpected))
			})
		})

	})

	Context("with dlp", func() {
		BeforeEach(func() {
			testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
				TapServer: &tap_server.InstanceConfig{
					EnableDataScrubbing: true,
				},
			})
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

		Context("tap functionality", func() {

			// Test whether the tap server scrubs sensitive data from the trace
			// information. Note that this is happening in the tap server itself,
			// and not in the envoy tap filter, so do not take away the impression
			// that we provide this functionality to users!!!
			It("scrubs sensitive data from headers", func() {
				ssnHeader := "x-some-ssn"
				ssnInput := "123456789"
				ssnExpected := "XXXXXXX89"

				creditCardHeader := "x-some-credit-card"
				creditCardInput := "5151-5151-5151-5151"
				creditCardExpected := "XXXX-XXXX-XXXX-5151"

				request := testContext.GetHttpRequestBuilder().
					WithHeader("X-some-ssn", ssnInput).
					WithHeader("X-some-credit-card", creditCardInput).
					Build()
				tapData := sendTestRequest(request)
				traceData := tapData.TraceData.GetHttpBufferedTrace()

				response := &http.Response{
					Header: make(http.Header),
				}
				for _, header := range traceData.GetRequest().GetHeaders() {
					response.Header.Add(header.Key, header.Value)
				}
				Expect(response).To(matchers.ConsistOfHeaders(
					http.Header{
						ssnHeader:        []string{ssnExpected},
						creditCardHeader: []string{creditCardExpected},
					}), "should have headers with sensitive information removed")
			})

			It("scrubs sensitive data from the body", func() {
				bodyInput := "sample SSN: 012-34-5678 sample credit card: 5151-5151-5151-5151"
				bodyExpected := "sample SSN: XXX-XX-XX78 sample credit card: XXXX-XXXX-XXXX-5151"
				request := testContext.GetHttpRequestBuilder().WithBody(bodyInput).Build()
				tapData := sendTestRequest(request)
				traceData := tapData.TraceData.GetHttpBufferedTrace()
				Expect(traceData).NotTo(BeNil())
				Expect(string(traceData.GetRequest().GetBody().GetAsBytes())).To(Equal(bodyExpected))
				// since the upstream is an echo server, we expect the body to be
				// the same. this assertion isn't strictly necessary - it doesn't
				// test any specific behaviour - but it doesn't harm anything, so
				// i'll just leave it in
				Expect(string(traceData.GetResponse().GetBody().GetAsBytes())).To(Equal(bodyExpected))
			})
		})

	})
})
