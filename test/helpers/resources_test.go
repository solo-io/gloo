package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"
)

var _ = Describe("Resource Creation", func() {

	It("can make distinict resources as expected", func() {
		Expect(helpers.Upstream(1)).To(Not(BeEquivalentTo(helpers.Upstream(2))))
		Expect(helpers.Endpoint(1)).To(Not(BeEquivalentTo(helpers.Endpoint(2))))
	})

	It("creates snapshots deterministically by default", func() {
		scfg := helpers.ScaleConfig{
			Upstreams: 1,
			Endpoints: 10,
		}
		Expect(helpers.ScaledSnapshot(scfg)).To(Equal(helpers.ScaledSnapshot(scfg)))
	})

	It("can mutate snapshot upstreams", func() {
		scfg := helpers.ScaleConfig{
			Upstreams: 1,
			Endpoints: 10,
		}
		snapOriginal := helpers.ScaledSnapshot(scfg)
		snapToMutate := snapOriginal.Clone()
		mutatedSnap := helpers.MutateSnapUpstreams(&snapToMutate,
			func(up *v1.Upstream) {
				up.Metadata.Name = "mutated"
			},
		)

		Expect(mutatedSnap).To(Not(Equal(snapOriginal)))

	})

})


