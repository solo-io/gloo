package matchers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"

	"github.com/onsi/gomega/matchers"
	"github.com/onsi/gomega/types"
)

var _ types.GomegaMatcher = new(HaveHttpResponseMatcher)

// HaveOkResponse expects a http response with a 200 status code
func HaveOkResponse() types.GomegaMatcher {
	return HaveStatusCode(http.StatusOK)
}

// HaveStatusCode expects a http response with a particular status code
func HaveStatusCode(statusCode int) types.GomegaMatcher {
	return HaveHttpResponse(&HttpResponse{
		StatusCode: statusCode,
		Body:       gstruct.Ignore(),
	})
}

// HaveExactResponseBody expects a 200 response with a body that matches the provided string
func HaveExactResponseBody(body string) types.GomegaMatcher {
	return HaveHttpResponse(&HttpResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	})
}

// HavePartialResponseBody expects a 200 response with a body that contains the provided substring
func HavePartialResponseBody(substring string) types.GomegaMatcher {
	return HaveHttpResponse(&HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gomega.ContainSubstring(substring),
	})
}

// HaveOkResponseWithHeaders expects an 200 response with a set of headers that match the provided headers
func HaveOkResponseWithHeaders(headers map[string]interface{}) types.GomegaMatcher {
	return HaveHttpResponse(&HttpResponse{
		StatusCode: http.StatusOK,
		Body:       gomega.BeEmpty(),
		Headers:    headers,
	})
}

// HaveOKResponseWithJSONContains expects a 200 response with a body that contains the provided JSON
func HaveOKResponseWithJSONContains(jsonBody []byte) types.GomegaMatcher {
	return HaveHttpResponse(&HttpResponse{
		StatusCode: http.StatusOK,
		Body:       JSONContains(jsonBody),
	})
}

// HttpResponse defines the set of properties that we can validate from an http.Response
type HttpResponse struct {
	// StatusCode is the expected status code for an http.Response
	// Required
	StatusCode int
	// Body is the expected response body for an http.Response
	// Body can be of type: {string, bytes, GomegaMatcher}
	// Optional: If not provided, defaults to an empty string
	Body interface{}
	// Headers is the set of expected header values for an http.Response
	// Each header can be of type: {string, GomegaMatcher}
	// Optional: If not provided, does not perform header validation
	Headers map[string]interface{}
	// Protocol is the expected protocol of an http.Response
	// Optional: If not provided, does not perform additional validation
	Protocol string
	// Custom is a generic matcher that can be applied to validate any other properties of an http.Response
	// Optional: If not provided, does not perform additional validation
	Custom types.GomegaMatcher
}

func (r *HttpResponse) String() string {
	var bodyString string
	switch bodyMatcher := r.Body.(type) {
	case string:
		bodyString = bodyMatcher
	case []byte:
		bodyString = string(bodyMatcher)
	case types.GomegaMatcher:
		bodyString = fmt.Sprintf("%#v", bodyMatcher)
	}

	return fmt.Sprintf("HttpResponse{StatusCode: %d, Body: %s, Headers: %v, Protocol: %s, Custom: %v}",
		r.StatusCode, bodyString, r.Headers, r.Protocol, r.Custom)

}

// HaveHttpResponse returns a GomegaMatcher which validates that an http.Response contains
// particular expected properties (status, body..etc)
// If an expected body isn't specified, the body is not matched
func HaveHttpResponse(expected *HttpResponse) types.GomegaMatcher {
	expectedCustomMatcher := expected.Custom
	if expected.Custom == nil {
		// Default to an always accept matcher
		expectedCustomMatcher = gstruct.Ignore()
	}

	var partialResponseMatchers []types.GomegaMatcher
	partialResponseMatchers = append(partialResponseMatchers, &matchers.HaveHTTPStatusMatcher{
		Expected: []interface{}{
			expected.StatusCode,
		},
	})
	if expected.Body != nil {
		partialResponseMatchers = append(partialResponseMatchers, &matchers.HaveHTTPBodyMatcher{
			Expected: expected.Body,
		})
	}
	if expected.Protocol != "" {
		partialResponseMatchers = append(partialResponseMatchers, HaveProtocol(expected.Protocol))
	}
	for headerName, headerMatch := range expected.Headers {
		partialResponseMatchers = append(partialResponseMatchers, &matchers.HaveHTTPHeaderWithValueMatcher{
			Header: headerName,
			Value:  headerMatch,
		})
	}
	partialResponseMatchers = append(partialResponseMatchers, expectedCustomMatcher)

	return &HaveHttpResponseMatcher{
		Expected:        expected,
		responseMatcher: gomega.And(partialResponseMatchers...),
	}
}

type HaveHttpResponseMatcher struct {
	Expected *HttpResponse

	responseMatcher types.GomegaMatcher

	// An internal utility for tracking whether we have evaluated this matcher
	// There is a comment within the Match method, outlining why we introduced this
	evaluated bool
}

func (m *HaveHttpResponseMatcher) Match(actual interface{}) (success bool, err error) {
	if m.evaluated {
		// Matchers are intended to be short-lived, and we have seen inconsistent behaviors
		// when evaluating the same matcher multiple times.
		// For example, the underlying http body matcher caches the response body, so if you are wrapping this
		// matcher in an Eventually, you need to create a new matcher each iteration.
		// This error is intended to help prevent developers hitting this edge case
		return false, errors.New("using the same matcher twice can lead to inconsistent behaviors")
	}
	m.evaluated = true

	if ok, matchErr := m.responseMatcher.Match(actual); !ok {
		return false, matchErr
	}

	return true, nil
}

func (m *HaveHttpResponseMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("%s \n%s",
		m.responseMatcher.FailureMessage(actual),
		informativeComparison(m.Expected, actual))
}

func (m *HaveHttpResponseMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("%s \n%s",
		m.responseMatcher.NegatedFailureMessage(actual),
		informativeComparison(m.Expected, actual))
}

// informativeComparison returns a string which presents data to the user to help them understand why a failure occurred.
// The HaveHttpResponseMatcher uses an And matcher, which intentionally short-circuits and only
// logs the first failure that occurred.
// To help developers, we print more details in this function.
// NOTE: Printing the actual http.Response is challenging (since the body has already been read), so for now
// we do not print it.
func informativeComparison(expected, _ interface{}) string {
	expectedJson, _ := json.MarshalIndent(expected, "", "  ")

	return fmt.Sprintf("\nexpected: %s", expectedJson)
}
