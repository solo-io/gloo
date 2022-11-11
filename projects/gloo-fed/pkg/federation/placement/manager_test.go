package placement_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/core/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation/placement"
)

var _ = Describe("Manager", func() {

	var (
		eastManager placement.Manager
		westManager placement.Manager
	)

	BeforeEach(func() {
		eastManager = placement.NewManager("east-namespace", "east-pod")
		westManager = placement.NewManager("west-namespace", "west-pod")
	})

	Describe("GetBuilder", func() {

		It("returns new builder each invocation", func() {
			builderOne := eastManager.GetBuilder()
			placementOne := builderOne.Build(1)

			builderTwo := eastManager.GetBuilder()
			placementTwo := builderTwo.Build(2)

			Expect(placementOne.ObservedGeneration).To(Equal(int64(1)))
			Expect(placementTwo.ObservedGeneration).To(Equal(int64(2)))
		})

	})

	Describe("Get/SetPlacementStatus", func() {

		It("can get and set placement status", func() {
			mockResource := &mockResourceWithPlacement{}

			eastManager.SetPlacementStatus(mockResource, &v1.PlacementStatus{
				ObservedGeneration: 1,
			})
			Expect(mockResource.PlacementStatuses["east-namespace"].ObservedGeneration).To(Equal(int64(1)))
			Expect(eastManager.GetPlacementStatus(mockResource).ObservedGeneration).To(Equal(int64(1)))
		})

		It("can get and set placement status using two different managers", func() {
			mockResource := &mockResourceWithPlacement{}

			eastManager.SetPlacementStatus(mockResource, &v1.PlacementStatus{
				ObservedGeneration: 1,
			})
			westManager.SetPlacementStatus(mockResource, &v1.PlacementStatus{
				ObservedGeneration: 2,
			})

			expectedEastGeneration := int64(1)
			Expect(mockResource.PlacementStatuses["east-namespace"].ObservedGeneration).To(Equal(expectedEastGeneration))
			Expect(eastManager.GetPlacementStatus(mockResource).ObservedGeneration).To(Equal(expectedEastGeneration))

			expectedWestGeneration := int64(2)
			Expect(mockResource.PlacementStatuses["west-namespace"].ObservedGeneration).To(Equal(expectedWestGeneration))
			Expect(westManager.GetPlacementStatus(mockResource).ObservedGeneration).To(Equal(expectedWestGeneration))
		})

	})

})

var _ placement.ResourceWithPlacement = new(mockResourceWithPlacement)

type mockResourceWithPlacement struct {
	PlacementStatus   *v1.PlacementStatus
	PlacementStatuses map[string]*v1.PlacementStatus
}

func (m *mockResourceWithPlacement) GetPlacementStatus() *v1.PlacementStatus {
	return m.PlacementStatus
}

func (m *mockResourceWithPlacement) SetPlacementStatus(status *v1.PlacementStatus) {
	m.PlacementStatus = status
}

func (m *mockResourceWithPlacement) GetNamespacedPlacementStatuses() map[string]*v1.PlacementStatus {
	return m.PlacementStatuses
}

func (m *mockResourceWithPlacement) SetNamespacedPlacementStatuses(statuses map[string]*v1.PlacementStatus) {
	m.PlacementStatuses = statuses
}
