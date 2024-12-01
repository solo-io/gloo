package upstreams_test

import (
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/constants"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams"
	upstreams_kube "github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Conversions", func() {

	AfterEach(func() {
		err := os.Unsetenv(constants.GlooGatewayEnableK8sGwControllerEnv)
		Expect(err).NotTo(HaveOccurred())
	})

	DescribeTable("converting upstream to cluster and back",
		func(kubeGatewayEnabled bool, us *gloov1.Upstream, expectedCluster string) {
			err := os.Setenv(constants.GlooGatewayEnableK8sGwControllerEnv,
				strconv.FormatBool(kubeGatewayEnabled))
			Expect(err).NotTo(HaveOccurred())

			// make sure the upstream converts to the expected cluster name
			clusterName := upstreams.UpstreamToClusterName(us)
			Expect(clusterName).To(Equal(expectedCluster))

			// convert it back and make sure the ref is the same as the original upstream ref
			usRef, err := upstreams.ClusterToUpstreamRef(clusterName)
			Expect(err).ToNot(HaveOccurred())
			Expect(usRef).To(Equal(us.GetMetadata().Ref()))
		},
		Entry("kubeGateway enabled, real kube upstream", true, createKubeUpstream(true, "name", "ns", "svcName", "svcNs", 123), "kube-upstream:name_ns_svcNs_svcName_123"),
		Entry("kubeGateway enabled, fake kube upstream", true, createKubeUpstream(false, "name", "ns", "svcName", "svcNs", 123), "kube-svc:name_ns_svcNs_svcName_123"),
		Entry("kubeGateway enabled, non-kube upstream", true, createStaticUpstream("name", "ns"), "name_ns"),
		Entry("kubeGateway disabled, real kube upstream", false, createKubeUpstream(true, "name", "ns", "svcName", "svcNs", 123), "name_ns"),
		Entry("kubeGateway disabled, fake kube upstream", false, createKubeUpstream(false, "name", "ns", "svcName", "svcNs", 123), "kube-svc:name_ns"),
		Entry("kubeGateway disabled, non-kube upstream", false, createStaticUpstream("name", "ns"), "name_ns"),
	)
})

func createKubeUpstream(isReal bool, name string, namespace string,
	svcName string, svcNamespace string, svcPort uint32) *gloov1.Upstream {
	if !isReal {
		name = upstreams_kube.FakeUpstreamNamePrefix + name
	}
	return &gloov1.Upstream{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		UpstreamType: &gloov1.Upstream_Kube{
			Kube: &kubernetes.UpstreamSpec{
				ServiceName:      svcName,
				ServiceNamespace: svcNamespace,
				ServicePort:      svcPort,
			},
		},
	}
}

func createStaticUpstream(name string, namespace string) *gloov1.Upstream {
	return &gloov1.Upstream{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		UpstreamType: &gloov1.Upstream_Static{
			Static: &static.UpstreamSpec{},
		},
	}
}
