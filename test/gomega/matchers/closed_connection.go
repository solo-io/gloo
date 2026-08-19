package matchers

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// serverClosedIdleMessage is the message of net/http's errServerClosedIdle. It is matched as a
// literal because net/http does not export that error.
const serverClosedIdleMessage = "http: server closed idle connection"

// BeClosedConnectionError produces a matcher that will match the message of an error
// returned by net/http when the server closed the connection without sending a response.
//
// Which message the transport surfaces is a race, not a behavioral difference: when Envoy
// drops a connection immediately (as the connection limit and L4 local rate limit filters
// do), the read loop may observe the close either after the request was written, yielding
// "EOF", or while the connection is still considered idle, yielding "http: server closed
// idle connection". Both mean the same thing for these tests, so matching only one of them
// makes the assertion depend on how fast Envoy closes relative to the request write. Note
// this stays strict about the connection actually being dropped: a request that was not
// limited returns a response with a nil error and still fails the assertion.
//
// Intended for use with gomega.MatchError, which applies a sub-matcher to err.Error():
//
//	Expect(err).To(MatchError(matchers.BeClosedConnectionError()))
func BeClosedConnectionError() types.GomegaMatcher {
	return gomega.Or(
		gomega.ContainSubstring("EOF"),
		gomega.ContainSubstring(serverClosedIdleMessage),
	)
}
