package testcases

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/snapshot"
	"github.com/solo-io/gloo/test/snapshot/utils"
)

var TestGatewaySubset = func(
	ctx context.Context,
	runner snapshot.TestRunner,
	env *snapshot.TestEnv,
	customSetupAssertions func(),
) {
	err := runner.Run(ctx)
	Expect(err).NotTo(HaveOccurred())

	// check custom assertions before sending requests
	customSetupAssertions()

	curlV1 := &utils.CurlFromPod{
		Url: fmt.Sprintf("http://%s.%s:%d/get", env.GatewayName, env.GatewayNamespace, env.GatewayPort),
		Cluster: &utils.KubeContext{
			Context:           env.ClusterContext,
			ClusterName:       env.ClusterName,
			KubernetesClients: runner.ClientSet.KubeClients(),
		},
		App:       "curl",
		Namespace: "curl",
		Headers:   []string{"host: httpbin.example.com"},
	}

	// This tests assumes that curl pod is installed in the curl namespace
	curlV2 := &utils.CurlFromPod{
		Url: fmt.Sprintf("http://%s.%s:%d/get", env.GatewayName, env.GatewayNamespace, env.GatewayPort),
		Cluster: &utils.KubeContext{
			Context:           env.ClusterContext,
			ClusterName:       env.ClusterName,
			KubernetesClients: runner.ClientSet.KubeClients(),
		},
		App:       "curl",
		Namespace: "curl",
		Headers:   []string{"host: httpbin.example.com", "version: v2"},
	}

	// First check we get a healthy response
	Eventually(func(g Gomega) {
		output, err := curlV1.Execute(GinkgoWriter)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(output).To(ContainSubstring("200 OK"), "expected 200 OK in output")
	}, 60*time.Second, 1*time.Second).Should(Succeed())

	// Check only subset v1 is hit when no header is provided
	Consistently(func(g Gomega) {
		output, err := curlV1.Execute(GinkgoWriter)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(output).To(ContainSubstring("200 OK"), "expected 200 OK in output")
		g.Expect(output).To(ContainSubstring("httpbin-v1"), "expected httpbin-v1")
		g.Expect(output).To(Not(ContainSubstring("httpbin-v2")), "should not hit httpbin-v2")
	}, 5*time.Second, 1*time.Second).Should(Succeed())

	// Check only subset v2 is hit when header is set
	Consistently(func(g Gomega) {
		output, err := curlV2.Execute(GinkgoWriter)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(output).To(ContainSubstring("200 OK"), "expected 200 OK in output")
		g.Expect(output).To(ContainSubstring("httpbin-v2"), "expected httpbin-v2")
		g.Expect(output).To(Not(ContainSubstring("httpbin-v1")), "should not hit httpbin-v1")
	}, 5*time.Second, 1*time.Second).Should(Succeed())
}
