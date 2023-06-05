package matchers

import (
	"sort"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func withPercentile(percentile int) func(durations []time.Duration) time.Duration {
	return func(durations []time.Duration) time.Duration {
		sort.Slice(durations, func(i, j int) bool {
			return durations[i] < durations[j]
		})
		return durations[int(float64(len(durations))*(float64(percentile)/float64(100)))]
	}
}

// Percentile returns a matcher requiring the given slice of durations to be less than the given upperBound at the given percentile
func Percentile(percentile int, upperBound time.Duration) types.GomegaMatcher {
	return gomega.WithTransform(withPercentile(percentile), gomega.BeNumerically("<", upperBound))
}

// Median returns a matcher requiring the given slice of durations have a median valude less than the given upperBound
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
