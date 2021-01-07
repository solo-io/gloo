package placement_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	fed_core_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation/placement"
)

var _ = Describe("StatusBuilder", func() {
	var (
		builder placement.StatusBuilder

		cluster1, cluster2 = "cluster-one", "cluster-two"
		ns1, ns2           = "ns-one", "ns-two"
	)

	BeforeEach(func() {
		builder = placement.NewFactory("pod").GetBuilder()
	})

	Describe("Build", func() {
		goodPlace := fed_core_v1.PlacementStatus_Namespace{
			State: fed_core_v1.PlacementStatus_PLACED,
		}

		badPlace := fed_core_v1.PlacementStatus_Namespace{
			State: fed_core_v1.PlacementStatus_FAILED,
		}

		It("works when all statuses are PLACED", func() {
			builder.AddDestinations([]string{cluster1}, []string{ns1, ns2}, goodPlace)
			builder.AddDestination(cluster2, ns2, goodPlace)

			actual := builder.Build(3)
			Expect(actual).To(Equal(&fed_core_v1.PlacementStatus{
				Clusters: map[string]*fed_core_v1.PlacementStatus_Cluster{
					cluster1: {
						Namespaces: map[string]*fed_core_v1.PlacementStatus_Namespace{
							ns1: &goodPlace,
							ns2: &goodPlace,
						},
					},
					cluster2: {
						Namespaces: map[string]*fed_core_v1.PlacementStatus_Namespace{
							ns2: &goodPlace,
						},
					},
				},
				State:              fed_core_v1.PlacementStatus_PLACED,
				ObservedGeneration: 3,
				WrittenBy:          "pod",
			}))
		})

		It("works when a status has state PLACEMENT_FAILED", func() {
			builder.AddDestination(cluster1, ns1, goodPlace)
			builder.AddDestination(cluster1, ns2, badPlace)
			builder.AddDestination(cluster2, ns2, goodPlace)

			actual := builder.Build(4)
			Expect(actual).To(Equal(&fed_core_v1.PlacementStatus{
				Clusters: map[string]*fed_core_v1.PlacementStatus_Cluster{
					cluster1: {
						Namespaces: map[string]*fed_core_v1.PlacementStatus_Namespace{
							ns1: &goodPlace,
							ns2: &badPlace,
						},
					},
					cluster2: {
						Namespaces: map[string]*fed_core_v1.PlacementStatus_Namespace{
							ns2: &goodPlace,
						},
					},
				},
				State:              fed_core_v1.PlacementStatus_FAILED,
				Message:            placement.FailedToPlaceResource,
				ObservedGeneration: 4,
				WrittenBy:          "pod",
			}))
		})
	})

	Describe("Eject", func() {
		It("works when a resource has been marked INVALID", func() {
			actual := builder.UpdateUnprocessed(&fed_core_v1.PlacementStatus{}, "foo", fed_core_v1.PlacementStatus_INVALID).
				Eject(100)

			Expect(actual).To(Equal(&fed_core_v1.PlacementStatus{
				State:              fed_core_v1.PlacementStatus_INVALID,
				Message:            "foo",
				ObservedGeneration: 100,
				WrittenBy:          "pod",
			}))
		})
	})
})
