package assertions

import (
	"context"
	"io"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
)

func (p *Provider) AssertEnvoyAdminApi(
	ctx context.Context,
	envoyDeployment metav1.ObjectMeta,
	adminAssertions ...func(ctx context.Context, adminClient *admincli.Client),
) {
	// Before opening a port-forward, we assert that there is at least one Pod that is ready
	p.EventuallyRunningReplicas(ctx, envoyDeployment, BeNumerically(">=", 1))

	portForwarder, err := p.clusterContext.Cli.StartPortForward(ctx,
		portforward.WithDeployment(envoyDeployment.GetName(), envoyDeployment.GetNamespace()),
		portforward.WithPorts(admincli.DefaultAdminPort, admincli.DefaultAdminPort),
	)
	p.Require.NoError(err, "can open port-forward")
	defer func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}()

	adminClient := admincli.NewClient().
		WithReceiver(io.Discard). // adminAssertion can overwrite this
		WithCurlOptions(
			curl.WithRetries(3, 0, 10),
			curl.WithPort(admincli.DefaultAdminPort),
		)

	for _, adminAssertion := range adminAssertions {
		adminAssertion(ctx, adminClient)
	}
}
