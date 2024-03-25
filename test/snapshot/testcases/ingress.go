package testcases

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/skv2/codegen/util"

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
	dir := util.MustGetThisDir()
	runner.ResultsByGateway = map[types.NamespacedName]snapshot.ExpectedTestResult{
		{
			Namespace: "default",
			Name:      "example-gateway-http",
		}: {
			Proxy: dir + "/outputs/http-routing-proxy.yaml",
		},
	}

	err := runner.Run(ctx, inputs)
	Expect(err).NotTo(HaveOccurred())

	// check custom assertions before sending requests
	customSetupAssertions()

	// This tests assumes that curl pod is installed in the curl namespace
	// curl gateway-proxy.gloo-system:80/headers -H "host: httpbin.example.com"  -v
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
		g.Expect(output).To(ContainSubstring("200 OK"))
	}).Should(Succeed())
}
