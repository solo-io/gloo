package assertions

import (
	"context"
	"io"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/gloo/pkg/utils/cmdutils"
	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/gomega/transforms"
	"github.com/solo-io/go-utils/threadsafe"
)

// EventuallyEnvoyReachable checks that the deployment is ready, opens a port-forward, and
// asserts the Envoy instance is eventually reachable on the port-forward. It returns a func
// to close the port-forward, which the caller is responsible for invoking, usually as a defer
// immediately after this function returns.
func (p *Provider) EventuallyEnvoyReachable(
	ctx context.Context,
	envoyDeployment metav1.ObjectMeta,
	localPort, remotePort int,
	checkPath string,
	checkCode int,
) func() {
	// Before opening a port-forward, we assert that there is at least one Pod that is ready
	p.EventuallyRunningReplicas(ctx, envoyDeployment, BeNumerically(">=", 1))

	pf, err := p.clusterContext.Cli.StartPortForward(ctx,
		portforward.WithDeployment(envoyDeployment.GetName(), envoyDeployment.GetNamespace()),
		portforward.WithPorts(localPort, remotePort),
	)
	p.Require.NoError(err, "can open port-forward")
	pfClose := func() {
		pf.Close()
		pf.WaitForStop()
	}

	// In case the Eventually fails, which manifests as a panic, we close the port-forward and re-assert the error
	defer func() {
		err := recover()
		if err != nil {
			pfClose()
		}
		Expect(err).ToNot(HaveOccurred())
	}()

	// We are NOT using the AssertEventualCurlResponse here because we are testing connectivity from
	// the port forward, i.e. outside the cluster.
	Eventually(func(g Gomega) {
		expectedResponseMatcher := WithTransform(transforms.WithCurlHttpResponse, matchers.HaveHttpResponse(&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gstruct.Ignore(),
		}))

		var buf threadsafe.Buffer
		curlCmd := cmdutils.Command(ctx, "curl", curl.BuildArgs(curl.WithHostPort(pf.Address()), curl.WithPath(checkPath), curl.VerboseOutput(), curl.WithHeadersOnly())...)
		err := curlCmd.WithStdout(&buf).WithStderr(&buf).Run()
		g.Expect(err).ToNot(HaveOccurred(), "executing curl command")
		g.Expect(buf.String()).To(expectedResponseMatcher)
	}).WithTimeout(time.Second * 3).WithPolling(time.Millisecond * 100).Should(Succeed())

	return pfClose
}

func (p *Provider) AssertEnvoyAdminApi(
	ctx context.Context,
	envoyDeployment metav1.ObjectMeta,
	adminAssertions ...func(ctx context.Context, adminClient *admincli.Client),
) {
	pfClose := p.EventuallyEnvoyReachable(ctx, envoyDeployment, admincli.DefaultAdminPort, admincli.DefaultAdminPort, "/ready", http.StatusOK)
	defer pfClose()

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
