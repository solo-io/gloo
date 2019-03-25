package nackdetector_test

import (
	"context"

	"go.opencensus.io/stats/view"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/projects/gloo/pkg/nack_detector"
)

var _ = Describe("Stats", func() {
	var (
		changeState func(newstate State)
	)
	BeforeEach(func() {
		s := NewStatsGen(context.Background())
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
