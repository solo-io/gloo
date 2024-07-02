package gloo_test

import (
	"context"
	"net"
	"time"

	"github.com/avast/retry-go/v4"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
)

var _ = Describe("Utils", func() {

	Context("portforward", func() {
		var (
			pf                  portforward.PortForwarder
			ctx                 context.Context
			cancel              context.CancelFunc
			defaultRetryOptions = []retry.Option{
				retry.LastErrorOnly(true),
				retry.Delay(100 * time.Millisecond),
				retry.DelayType(retry.BackOffDelay),
				retry.Attempts(5),
			}
		)

		AfterEach(func() {
			cancel()
			pf.Close()
			pf.WaitForStop()
		})

		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())
			// We are using kube-dns because we know we have a kubectl and a kind cluster on our CI
			// env. If that ever changes we need to re-assess this test strategy.
			// We are using port 30053 because 53 is privileged.
			pf = portforward.NewPortForwarder(portforward.WithService("kube-dns", "kube-system"), portforward.WithPorts(30053, 53))
		})

		It("Creates a usable portforward with Start", func() {
			err := pf.Start(ctx, defaultRetryOptions...)
			Expect(err).ToNot(HaveOccurred())
			_, err = net.Dial("tcp4", pf.Address())
			Expect(err).ToNot(HaveOccurred())

		})

	})
})
