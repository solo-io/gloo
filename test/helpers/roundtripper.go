package helpers

import (
	"net/http"

	"github.com/solo-io/glue/pkg/log"
)

func AddLoggingToTransport(rt http.RoundTripper) http.RoundTripper {
	return &loggingRoundTripper{wrapped: rt}
}

type loggingRoundTripper struct {
	wrapped http.RoundTripper
}

func (rt *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	log.GreyPrintf("Logging Request: %v", req)
	res, err := rt.wrapped.RoundTrip(req)
	log.GreyPrintf("Logging response: %v", res)
	if err != nil {
		log.Debugf("Logging err: %v", err)
	}
	return res, err
}
