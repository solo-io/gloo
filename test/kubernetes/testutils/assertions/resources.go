package assertions

import (
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e/helper"
)

// Checks GetNamespacedStatuses status for gloo installation namespace
func (p *Provider) EventuallyResourceExsits(getter helpers.ResourceGetter, timeout ...time.Duration) {
	ginkgo.GinkgoHelper()

	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	gomega.Eventually(func(g gomega.Gomega) {
		_, err := getter()
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource")
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}
