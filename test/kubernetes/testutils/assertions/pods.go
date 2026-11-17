package assertions

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kube2e/helper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// defaultPodMatchTimeout is how long EventuallyPodsRunning and EventuallyPodsMatches wait when
// the caller passes no timeout. A new pod's startup includes pulling its image, and test images
// that are not preloaded into the kind cluster can take longer than the 20s curl default to pull.
const defaultPodMatchTimeout = 2 * time.Minute

// EventuallyPodsRunning asserts that the pod(s) are running and ready.
func (p *Provider) EventuallyPodsRunning(
	ctx context.Context,
	podNamespace string,
	listOpt metav1.ListOptions,
	timeout ...time.Duration,
) {
	p.EventuallyPodsMatches(ctx, podNamespace, listOpt, matchers.PodMatches(matchers.ExpectedPod{Status: corev1.PodRunning, Ready: true}), timeout...)
}

// EventuallyPodsMatches asserts that the pod(s) in the given namespace matches the provided matcher.
// Without a timeout argument, it waits up to defaultPodMatchTimeout.
func (p *Provider) EventuallyPodsMatches(
	ctx context.Context,
	podNamespace string,
	listOpt metav1.ListOptions,
	matcher types.GomegaMatcher,
	timeout ...time.Duration,
) {
	// Use the pod default timeout unless the caller set one; polling keeps the shared default
	if len(timeout) == 0 {
		timeout = []time.Duration{defaultPodMatchTimeout}
	}
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)

	// Poll until every matching pod satisfies the matcher
	p.Gomega.Eventually(func(g gomega.Gomega) {
		pods, err := p.clusterContext.Clientset.CoreV1().Pods(podNamespace).List(ctx, listOpt)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to list pods")
		g.Expect(pods.Items).NotTo(gomega.BeEmpty(), "No pods found")
		for _, pod := range pods.Items {
			g.Expect(pod).To(matcher)
		}
	}).
		WithTimeout(currentTimeout).
		WithPolling(pollingInterval).
		Should(gomega.Succeed(), fmt.Sprintf("Failed to match pod in namespace %s", podNamespace))
}

// EventuallyPodsNotExist asserts that the pod(s) are no longer present
func (p *Provider) EventuallyPodsNotExist(
	ctx context.Context,
	podNamespace string,
	listOpt metav1.ListOptions,
	timeout ...time.Duration,
) {
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)

	p.Gomega.Eventually(func(g gomega.Gomega) {
		pods, err := p.clusterContext.Clientset.CoreV1().Pods(podNamespace).List(ctx, listOpt)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to list pods")
		g.Expect(pods.Items).To(gomega.BeEmpty(), "Pods found")
	}).
		WithTimeout(currentTimeout).
		WithPolling(pollingInterval).
		Should(gomega.Succeed(), fmt.Sprintf("Failed to match pod in namespace %s", podNamespace))
}
