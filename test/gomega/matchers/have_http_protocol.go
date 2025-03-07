package matchers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

var _ types.GomegaMatcher = new(HaveHTTPProtocolMatcher)

// HaveHTTPProtocolMatcher matches the request with the expected protocol
// This has been inspired by gomega.HaveHTTPStatusMatcher
type HaveHTTPProtocolMatcher struct {
	Expected interface{}
}

func (matcher *HaveHTTPProtocolMatcher) Match(actual interface{}) (success bool, err error) {
	var resp *http.Response
	switch a := actual.(type) {
	case *http.Response:
		resp = a
	case *httptest.ResponseRecorder:
		resp = a.Result()
	default:
		return false, fmt.Errorf("HaveHTTPProtocol matcher expects *http.Response or *httptest.ResponseRecorder. Got:\n%s", format.Object(actual, 1))
	}

	switch e := matcher.Expected.(type) {
	case string:
		if resp.Proto == e {
			return true, nil
		}
	default:
		return false, fmt.Errorf("HaveHTTPProtocol matcher must be passed a string type. Got:\n%s", format.Object(matcher.Expected, 1))
	}

	return false, nil
}

func (matcher *HaveHTTPProtocolMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\n%s\n%s", formatHttpResponse(actual), "to have HTTP protocol", matcher.Expected)
}

func (matcher *HaveHTTPProtocolMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\n%s\n%s", formatHttpResponse(actual), "not to have HTTP protocol", matcher.Expected)
}

func formatHttpResponse(input interface{}) string {
	var resp *http.Response
	switch r := input.(type) {
	case *http.Response:
		resp = r
	case *httptest.ResponseRecorder:
		resp = r.Result()
	default:
		return "cannot format invalid HTTP response"
	}

	body := "<nil>"
	if resp.Body != nil {
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			data = []byte("<error reading body>")
		}
		body = format.Object(string(data), 0)
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s<%s>: {\n", format.Indent, reflect.TypeOf(input)))
	s.WriteString(fmt.Sprintf("%s%sProtocol:   %s\n", format.Indent, format.Indent, resp.Proto))
	s.WriteString(fmt.Sprintf("%s%sStatus:     %s\n", format.Indent, format.Indent, format.Object(resp.Status, 0)))
	s.WriteString(fmt.Sprintf("%s%sStatusCode: %s\n", format.Indent, format.Indent, format.Object(resp.StatusCode, 0)))
	s.WriteString(fmt.Sprintf("%s%sBody:       %s\n", format.Indent, format.Indent, body))
	s.WriteString(fmt.Sprintf("%s}", format.Indent))

	return s.String()
}
