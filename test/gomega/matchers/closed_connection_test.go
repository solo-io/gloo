package matchers_test

import (
	"errors"
	"fmt"
	"io"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/gomega/matchers"
)

var _ = Describe("BeClosedConnectionError", func() {

	// The errors below are shaped like the ones net/http returns from Client.Do when the
	// server drops the connection, since that is what the e2e tests assert against.
	urlErr := func(err error) error {
		return &url.Error{Op: "Get", URL: "http://localhost:9086/", Err: err}
	}

	DescribeTable("matches an error from a dropped connection",
		func(err error) {
			Expect(err).To(MatchError(matchers.BeClosedConnectionError()))
		},
		Entry("EOF", urlErr(io.EOF)),
		Entry("unexpected EOF", urlErr(io.ErrUnexpectedEOF)),
		Entry("server closed idle connection", urlErr(errors.New("http: server closed idle connection"))),
	)

	DescribeTable("does not match an unrelated error",
		func(err error) {
			Expect(err).NotTo(MatchError(matchers.BeClosedConnectionError()))
		},
		Entry("connection refused", urlErr(errors.New("connect: connection refused"))),
		Entry("timeout", urlErr(fmt.Errorf("context deadline exceeded"))),
	)

})
