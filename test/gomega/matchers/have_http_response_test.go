package matchers_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/test/gomega/matchers"
)

var _ = Describe("HaveHttpStatusMultiMatcher", func() {

	DescribeTable("successful matches",
		func(expected []interface{}, responseStatusCode int, responseStatus string) {
			matcher := &matchers.HaveHttpStatusMultiMatcher{
				Expected: expected,
			}
			resp := &http.Response{
				StatusCode: responseStatusCode,
				Status:     responseStatus,
			}

			success, err := matcher.Match(resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(success).To(BeTrue())
		},
		Entry("single int - expects 200, gets 200", []interface{}{http.StatusOK}, http.StatusOK, "200 OK"),
		Entry("single string - expects '200 OK', gets 200", []interface{}{"200 OK"}, http.StatusOK, "200 OK"),
		Entry("multiple int - expects [200,201,202], gets 200", []interface{}{[]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}}, http.StatusOK, "200 OK"),
		Entry("multiple int - expects [200,201,202], gets 201", []interface{}{[]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}}, http.StatusCreated, "201 Created"),
		Entry("multiple int - expects [200,201,202], gets 202", []interface{}{[]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}}, http.StatusAccepted, "202 Accepted"),
		Entry("multiple string - expects ['200 OK','201 Created'], gets 200", []interface{}{[]string{"200 OK", "201 Created"}}, http.StatusOK, "200 OK"),
		Entry("multiple string - expects ['200 OK','201 Created'], gets 201", []interface{}{[]string{"200 OK", "201 Created"}}, http.StatusCreated, "201 Created"),
		Entry("mixed types - expects [200,'201 Created',404], gets 200", []interface{}{http.StatusOK, "201 Created", http.StatusNotFound}, http.StatusOK, "200 OK"),
		Entry("mixed types - expects [200,'201 Created',404], gets 201", []interface{}{http.StatusOK, "201 Created", http.StatusNotFound}, http.StatusCreated, "201 Created"),
		Entry("mixed types - expects [200,'201 Created',404], gets 404", []interface{}{http.StatusOK, "201 Created", http.StatusNotFound}, http.StatusNotFound, "404 Not Found"),
	)

	DescribeTable("failed matches",
		func(expected []interface{}, responseStatusCode int, responseStatus string) {
			matcher := &matchers.HaveHttpStatusMultiMatcher{
				Expected: expected,
			}
			resp := &http.Response{
				StatusCode: responseStatusCode,
				Status:     responseStatus,
			}

			success, err := matcher.Match(resp)
			Expect(err).ToNot(HaveOccurred())
			Expect(success).To(BeFalse())
		},
		Entry("single int - expects 200, gets 404", []interface{}{http.StatusOK}, http.StatusNotFound, "404 Not Found"),
		Entry("single string - expects '200 OK', gets 404", []interface{}{"200 OK"}, http.StatusNotFound, "404 Not Found"),
		Entry("multiple int - expects [200,201,202], gets 404", []interface{}{[]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}}, http.StatusNotFound, "404 Not Found"),
		Entry("multiple string - expects ['200 OK','201 Created'], gets 404", []interface{}{[]string{"200 OK", "201 Created"}}, http.StatusNotFound, "404 Not Found"),
		Entry("mixed types - expects [200,'201 Created'], gets 500", []interface{}{http.StatusOK, "201 Created"}, http.StatusInternalServerError, "500 Internal Server Error"),
		Entry("empty expected array, gets 200", []interface{}{}, http.StatusOK, "200 OK"),
		Entry("empty int slice, gets 200", []interface{}{[]int{}}, http.StatusOK, "200 OK"),
		Entry("empty string slice, gets 200", []interface{}{[]string{}}, http.StatusOK, "200 OK"),
	)

	DescribeTable("error cases",
		func(expected []interface{}, input interface{}, expectedErrorSubstring string) {
			matcher := &matchers.HaveHttpStatusMultiMatcher{
				Expected: expected,
			}

			success, err := matcher.Match(input)
			Expect(success).To(BeFalse())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(expectedErrorSubstring))
		},
		Entry("non-http.Response input - expects 200, gets string", []interface{}{http.StatusOK}, "not a response", "HaveHTTPStatus matcher expects *http.Response"),
		Entry("invalid expected type - float, response 200", []interface{}{3.14}, &http.Response{StatusCode: http.StatusOK}, "HaveHTTPStatus matcher expects int, string, or []int"),
		Entry("invalid expected type - struct, response 200", []interface{}{struct{}{}}, &http.Response{StatusCode: http.StatusOK}, "HaveHTTPStatus matcher expects int, string, or []int"),
	)

	Context("failure messages", func() {
		var matcher *matchers.HaveHttpStatusMultiMatcher
		var resp *http.Response

		BeforeEach(func() {
			matcher = &matchers.HaveHttpStatusMultiMatcher{
				Expected: []interface{}{http.StatusOK, http.StatusCreated},
			}
			resp = &http.Response{StatusCode: http.StatusNotFound}
		})

		It("should provide informative failure message", func() {
			message := matcher.FailureMessage(resp)
			Expect(message).To(ContainSubstring("Expected"))
			Expect(message).To(ContainSubstring("to have HTTP status"))
		})

		It("should provide informative negated failure message", func() {
			message := matcher.NegatedFailureMessage(resp)
			Expect(message).To(ContainSubstring("Expected"))
			Expect(message).To(ContainSubstring("not to have HTTP status"))
		})

		It("should include expected values in failure message", func() {
			message := matcher.FailureMessage(resp)
			Expect(message).To(ContainSubstring("200"))
			Expect(message).To(ContainSubstring("201"))
		})
	})

	Context("integration tests", func() {
		It("should work with httptest.ResponseRecorder - expects 200, gets 200", func() {
			matcher := &matchers.HaveHttpStatusMultiMatcher{
				Expected: []interface{}{http.StatusOK},
			}

			recorder := httptest.NewRecorder()
			recorder.WriteHeader(http.StatusOK)

			success, err := matcher.Match(recorder.Result())
			Expect(err).ToNot(HaveOccurred())
			Expect(success).To(BeTrue())
		})
	})
})
