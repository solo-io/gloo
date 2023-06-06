package matchers

import (
	"github.com/solo-io/gloo/test/gomega/transforms"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

// Percentile returns a matcher requiring the given slice of durations to be less than the given upperBound at the given percentile
func Percentile(percentile int, upperBound time.Duration) types.GomegaMatcher {
	return gomega.WithTransform(transforms.WithPercentile(percentile), gomega.BeNumerically("<", upperBound))
}

// Median returns a matcher requiring the given slice of durations have a median valude less than the given upperBound
func Median(upperBound time.Duration) types.GomegaMatcher {
	return gomega.WithTransform(transforms.WithMedian(), gomega.BeNumerically("<", upperBound))
}
