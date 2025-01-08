package assertions

import (
	"context"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	corev1 "k8s.io/api/core/v1"
)

func (p *Provider) ExternalTrafficPolicy(ctx context.Context, service corev1.Service, externalTrafficPolicyMatcher types.GomegaMatcher) {
	p.Gomega.Eventually(func(innerG Gomega) {
		service, err := kubeutils.GetService(ctx, p.clusterContext.Clientset, service.Name, service.Namespace)
		innerG.Expect(err).NotTo(HaveOccurred(), "can get service")
		innerG.Expect(service.Spec.ExternalTrafficPolicy).To(externalTrafficPolicyMatcher, "externalTrafficPolicy to match")
	}).
		WithContext(ctx).
		// It may take some time for pods to initialize and pull images from remote registries.
		// Therefore, we set a longer timeout, to account for latency that may exist.
		WithTimeout(time.Second * 30).
		WithPolling(time.Millisecond * 200).
		Should(Succeed())
}
