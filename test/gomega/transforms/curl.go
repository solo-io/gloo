package transforms

import (
	"net/http"
	"strings"
)

// WithCurlHttpResponse is a Gomega Transform that converts the string return by an exec.Curl
// and transforms it into an http.Response. This is useful to be used in tandem with matchers.HaveHttpResponse
// NOTE: This is not feature complete, as we do not convert the entire response.
// For now, we handle HTTP/1.1 response headers and status
func WithCurlHttpResponse(curlResponse string) *http.Response {
	headers := make(http.Header)
	// response headers start with "< "
	for _, header := range strings.Split(curlResponse, "< ") {
		headerParts := strings.Split(header, ": ")
		if len(headerParts) == 2 {
			// strip "\r\n" from the end of the value
			headers.Add(strings.ToLower(headerParts[0]), strings.TrimSuffix(headerParts[1], "\r\n"))
		}
	}

	// A more robust solution would be to search for a regex and extract the statusCode from it
	// For now, we only use 200 and 404 assertions, so this works
	statusCode := 0
	if strings.Contains(curlResponse, "HTTP/1.1 200") {
		statusCode = 200
	}
	if strings.Contains(curlResponse, "HTTP/1.1 404") {
		statusCode = 404
	}
	return &http.Response{
		StatusCode: statusCode,
		Header:     headers,
	}
}
