package assertions

import (
	"context"
	"io"
	"time"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/testutils/cmd/envoyadmin"
)

type EnvoyAdminAssertion func(ctx context.Context, adminClient *envoyadmin.Client, clientRefresh func(ctx context.Context, adminClient *envoyadmin.Client, expectedReplicas int) *envoyadmin.Client)

func (p *Provider) AssertEnvoyAdminApi(
	ctx context.Context,
	envoyDeployment metav1.ObjectMeta,
	adminAssertions ...EnvoyAdminAssertion,
) {
	// Before opening a port-forward, we assert that there is at least one Pod that is ready
	p.EventuallyRunningReplicas(ctx, envoyDeployment, BeNumerically(">=", 1))

	pods, err := kubeutils.GetReadyPodsForDeployment(ctx, p.clusterContext.Clientset, envoyDeployment)
	p.Require.NoError(err)
	p.Require.NotEmpty(pods)

	adminClient := envoyadmin.NewClient().
		WithReceiver(io.Discard). // adminAssertion can overwrite this
		WithCurlOptions(
			curl.WithRetries(3, 0, 10),
			curl.WithPort(admincli.DefaultAdminPort),
		).
		WithEnvoyMeta(metav1.ObjectMeta{Name: pods[0], Namespace: envoyDeployment.Namespace})

	refreshFunc := func(ctx context.Context, adminClient *envoyadmin.Client, expectedReplicas int) *envoyadmin.Client {
		return p.getReadyEnvoyForClient(ctx, adminClient, envoyDeployment, expectedReplicas)
	}
	for _, adminAssertion := range adminAssertions {
		// If there is a pod that is transitioning from Ready to Terminating it, we could be targeting it with
		// our ephemeral exec. There isn't a solid way to guarantee this isn't going to occur is to wrap it in
		// Eventually. The assertion we call should panic
		adminAssertion(ctx, adminClient, refreshFunc)
	}
}

func (p *Provider) getReadyEnvoyForClient(ctx context.Context, adminClient *envoyadmin.Client, envoyDeployment metav1.ObjectMeta, expectedReplicas int) *envoyadmin.Client {

	var pods []string
	var err error
	Eventually(func(g Gomega) {
		pods, err = kubeutils.GetReadyPodsForDeployment(ctx, p.clusterContext.Clientset, envoyDeployment)
		p.Require.NoError(err)
		p.Require.NotEmpty(pods)
		// We have to make sure that the pod count is what we expect otherwise we could be selecting a pod which
		// is Terminating BUT HAS NO STATUS CONDITION THAT WE CAN CHECK??!!
		g.Expect(pods).Should(HaveLen(expectedReplicas))
	}).WithTimeout(time.Second * 30).WithPolling(time.Second).Should(Succeed())

	adminClient = adminClient.
		WithEnvoyMeta(metav1.ObjectMeta{Name: pods[0], Namespace: envoyDeployment.Namespace})

	Eventually(func() bool {
		return adminClient.GetReady(ctx)
	}).
		WithTimeout(10 * time.Second).
		WithPolling(100 * time.Millisecond).
		Should(BeTrue())

	return adminClient

}
