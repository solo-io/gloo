package matchers

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/gloo/test/gomega/transforms"
)

const (
	HTTP1Protocol = "HTTP/1.1"
	HTTP2Protocol = "HTTP/2"
)

// HaveProtocol expects an http response with the given protocol
func HaveProtocol(protocol string) types.GomegaMatcher {
	if protocol == "" {
		// If no keys are defined, we create a matcher that always succeeds
		// If we do not this we will create an And matcher for 0 objects, which leads to a panic
		return gstruct.Ignore()
	}
	//nolint:bodyclose // The caller of this matcher constructor should be responsible for ensuring the body close
	return gomega.WithTransform(transforms.WithProtocol(), gomega.Equal(protocol))
}

// HaveHTTP1Protocol expects an http response with the HTTP/1.1 protocol
func HaveHTTP1Protocol() types.GomegaMatcher {
	return HaveProtocol(HTTP1Protocol)
}

// HaveHTTP2Protocol expects an http response with the HTTP/2 protocol
func HaveHTTP2Protocol() types.GomegaMatcher {
	return HaveProtocol(HTTP2Protocol)
}
