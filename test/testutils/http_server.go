package testutils

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func StartCancellableTracingServer(serverContext context.Context, address string, handler http.Handler) {
	tracingServer := &http.Server{
		Addr:    address,
		Handler: handler,
	}

	// Start a goroutine to handle requests
	go func() {
		defer GinkgoRecover()
		if err := tracingServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
		}
	}()

	// Start a goroutine to shutdown the server
	go func(serverCtx context.Context) {
		defer GinkgoRecover()

		<-serverCtx.Done()
		// tracingServer.Shutdown hangs with opentelemetry tests, probably
		// because the agent leaves the connection open. There's no need for a
		// graceful shutdown anyway, so just force it using Close() instead
		tracingServer.Close()
	}(serverContext)
}
