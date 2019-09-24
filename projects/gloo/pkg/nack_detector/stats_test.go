package nackdetector_test

import (
	"context"

	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/projects/gloo/pkg/nack_detector"
)

var _ = Describe("Stats", func() {
	var (
		changeState func(newstate State)
		s           *StatGen
	)
	BeforeEach(func() {
		s = NewStatsGen(context.Background())
		var id EnvoyStatusId
		st := New
		changeState = func(newstate State) {
			s.Stat(id, st, newstate)
			st = newstate
		}
	})

	It("should work", func() {
		changeState(New)
		VerifyState(1, 0, 0, 0)
		changeState(InSync)
		VerifyState(1, 1, 0, 0)
		changeState(OutOfSync)
		VerifyState(1, 0, 1, 0)
		changeState(OutOfSyncNack)
		VerifyState(1, 0, 0, 1)
		changeState(Gone)
		VerifyState(0, 0, 0, 0)
	})

	It("should record resource tags", func() {
		var id EnvoyStatusId
		id.StreamId = DiscoveryServiceId{TypeUrl: "cds"}
		var id2 EnvoyStatusId
		id2.StreamId = DiscoveryServiceId{TypeUrl: "rds"}
		s.Stat(id, New, New)
		s.Stat(id2, New, New)

		d, err := view.RetrieveData(GlooeTotalEntities.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(d)).To(BeNumerically(">=", 2))
		seenCds := 0
		seenRds := 0
		resourceTag := func(r string) tag.Tag { return tag.Tag{Key: tag.MustNewKey("resource"), Value: r} }
		for _, datum := range d {
			if len(datum.Tags) == 0 {
				continue
			}
			if datum.Tags[0] == resourceTag("cds") {
				seenCds += 1
			}
			if datum.Tags[0] == resourceTag("rds") {
				seenRds += 1
			}
		}
		Expect(seenCds).To(Equal(1))
		Expect(seenRds).To(Equal(1))
	})
})

func VerifyState(total, insync, outofsync, nacks int) {

	ExpectWithOffset(1, GetData(GlooeTotalEntities)).To(BeNumerically("==", total))
	ExpectWithOffset(1, GetData(GlooeInSync)).To(BeNumerically("==", insync))
	ExpectWithOffset(1, GetData(GlooeNack)).To(BeNumerically("==", nacks))
	ExpectWithOffset(1, GetData(GlooeOutOfSync)).To(BeNumerically("==", outofsync))
}

func GetData(v *view.View) float64 {

	d, err := view.RetrieveData(v.Name)
	ExpectWithOffset(2, err).NotTo(HaveOccurred())
	if len(d) == 0 {
		return 0
	}
	return d[len(d)-1].Data.(*view.SumData).Value
}
