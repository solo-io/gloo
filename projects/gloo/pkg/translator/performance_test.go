package translator_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	"github.com/solo-io/gloo/test/ginkgo/labels"
	"sort"
	"time"
)

var _ = FDescribe("Translation - Benchmarking Tests", Serial, Label(labels.Performance), func() {

	DescribeTable("Translate",
		func(desc, s string, max90Dur time.Duration) {
			experiment := gmeasure.NewExperiment("print strings")

			n := 20
			AddReportEntry(experiment.Name, experiment)

			statName := fmt.Sprintf("printing %s", desc)
			experiment.Sample(func(idx int) {
				experiment.MeasureDuration(statName, func() { print(s) })
			}, gmeasure.SamplingConfig{N: n, Duration: time.Minute})

			durations := experiment.Get(statName).Durations
			sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
			ninetyPct := experiment.Get(statName).Durations[int(float64(n)*.9)]
			Expect(ninetyPct).To(BeNumerically("<", max90Dur))
		},
		Entry("foo", "foo", "foo", time.Millisecond),
		Entry("bar", "bar", "bar", 10*time.Nanosecond),
		Entry("100", "100", longString(100), time.Millisecond),
		Entry("1000", "1000", longString(1000), time.Millisecond),
		Entry("10000", "10000", longString(10000), time.Millisecond),
		Entry("100000", "100000", longString(100000), 10*time.Millisecond),
	)
})

func longString(size int) string {
	s := ""
	for i := 0; i < size; i++ {
		s += "s"
	}

	return s
}
