package portforward_test

import (
	"context"
	"net"
	"time"

	"github.com/avast/retry-go/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
)

var _ = Describe("Portforwarder test", func() {
	var (
		pf                  PortForwarder
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
	})

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		pf = NewPortForwarder(WithService("kube-dns", "kube-system"), WithPorts(30053, 53))
	})

	It("Creates a flaky portforward with Start", func() {
		Eventually(func(g Gomega) error {
			err := pf.Start(ctx, defaultRetryOptions...)
			if err != nil {
				return err
			}

			_, err = net.Dial("tcp4", pf.Address())
			if err != nil {
				pf.Close()
				return err
			}
			pf.Close()
			return nil
		}, time.Minute, time.Millisecond*100).ShouldNot(Succeed())
	})

	It("Creates a usable portforward with StartAndWaitForConn", func() {
		err := pf.StartAndWaitForConn(ctx, defaultRetryOptions...)
		Expect(err).ToNot(HaveOccurred())
		_, err = net.Dial("tcp4", pf.Address())
		Expect(err).ToNot(HaveOccurred())
		pf.Close()

	})
})
