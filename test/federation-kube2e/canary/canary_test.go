package canary_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/testutils"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	v12 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"path/filepath"
	"time"
)

var _ = Describe("Canary Test", func() {

	It("can place federated upstream in remote clusters", func() {
		err := testutils.Kubectl("apply", "-f", filepath.Join(resourcesFolder, "federated-upstream.yaml"))
		Expect(err).NotTo(HaveOccurred())

		federatedClientset, err := v1.NewClientsetFromConfig(managementClusterConfig.RestConfig)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func(g Gomega) {
			federatedUpstream, err := federatedClientset.FederatedUpstreams().GetFederatedUpstream(ctx, client.ObjectKey{
				Namespace: federationNamespace,
				Name:      "federated-upstream",
			})
			g.Expect(err).NotTo(HaveOccurred())

			statuses := federatedUpstream.Status.GetNamespacedPlacementStatuses()
			g.Expect(statuses).NotTo(BeNil())
			g.Expect(statuses[releaseNamespace].GetState()).To(Equal(v12.PlacementStatus_PLACED))
			g.Expect(statuses[canaryNamespace].GetState()).To(Equal(v12.PlacementStatus_PLACED))
		}, time.Second*15, time.Second).Should(Succeed())
	})

})
