package testcases

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/utils"
)

var TestGatewayIngress = func(
	ctx context.Context,
	runner snapshot.TestRunner,
	env *snapshot.TestEnv,
	inputs []client.Object,
	customSetupAssertions func(),
) {
	err := runner.Run(ctx, inputs)
	Expect(err).NotTo(HaveOccurred())

	// check custom assertions before sending requests
	customSetupAssertions()

	// This tests assumes that curl pod is installed in the curl namespace
	// ie. curl gateway-proxy.gloo-system:80/headers -H "host: httpbin.example.com"  -v
	curl := &utils.CurlFromPod{
		Url: fmt.Sprintf("http://%s.%s:%d/headers", env.GatewayName, env.GatewayNamespace, env.GatewayPort),
		Cluster: &utils.KubeContext{
			Context:           env.ClusterContext,
			ClusterName:       env.ClusterName,
			KubernetesClients: runner.ClientSet.KubeClients(),
		},
		App:       "curl",
		Namespace: "curl",
		Headers:   []string{"host: httpbin.example.com"},
	}

	Eventually(func(g Gomega) {
		output, err := curl.Execute(GinkgoWriter)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(output).To(ContainSubstring("200 OK"), "expected 200 OK in output")
	}, 60*time.Second, 1*time.Second).Should(Succeed())
}
