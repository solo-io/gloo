package matchers

import (
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"sort"
	"time"
)

func withPercentile(percentile int) func(durations []time.Duration) time.Duration {
	return func(durations []time.Duration) time.Duration {
		sort.Slice(durations, func(i, j int) bool {
			return durations[i] < durations[j]
		})
		return durations[int(float64(len(durations))*(float64(percentile)/float64(100)))]
	}
}

func Percentile(percentile int, upperBound time.Duration) types.GomegaMatcher {
	return gomega.WithTransform(withPercentile(percentile), gomega.BeNumerically("<", upperBound))
}

func Median(upperBound time.Duration) types.GomegaMatcher {
	return gomega.WithTransform(func(durations []time.Duration) time.Duration {
		sort.Slice(durations, func(i, j int) bool {
			return durations[i] < durations[j]
		})
		var median time.Duration
		if l := len(durations); l%2 == 1 {
			median = durations[l/2]
		} else {
			median = (durations[l/2] + durations[l/2-1]) / 2
		}
		return median
	}, gomega.BeNumerically("<", upperBound))
}
