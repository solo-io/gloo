package multicluster_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
)

var _ = Describe("ClusterSet", func() {
	It("works", func() {
		clusterSet := multicluster.NewClusterSet()

		clusterSet.AddCluster(nil, "foo", nil)
		clusterSet.AddCluster(nil, "bar", nil)

		Expect(clusterSet.Exists("foo")).To(BeTrue())
		Expect(clusterSet.Exists("bar")).To(BeTrue())
		Expect(clusterSet.ListClusters()).To(HaveLen(2))
		Expect(clusterSet.ListClusters()).To(Equal([]string{"bar", "foo"}), "the clusters should be sorted alphabetically")

		clusterSet.RemoveCluster("foo")
		Expect(clusterSet.Exists("foo")).To(BeFalse())
		Expect(clusterSet.Exists("bar")).To(BeTrue())
		Expect(clusterSet.ListClusters()).To(HaveLen(1))
		Expect(clusterSet.ListClusters()).NotTo(ContainElement("foo"))
		Expect(clusterSet.ListClusters()).To(ContainElement("bar"))
	})
})
