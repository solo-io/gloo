package duration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/solo-io/gloo/pkg/utils/protoutils/duration"
)

var _ = Describe("MillisToDuration", func() {

	DescribeTable("converts milliseconds to durationpb.Duration correctly",
		func(inputMillis uint32, expectedSeconds int64, expectedNanos int32) {
			result := duration.MillisToDuration(inputMillis)

			Expect(result).NotTo(BeNil())
			Expect(result).To(BeAssignableToTypeOf(&durationpb.Duration{}))
			Expect(result.Seconds).To(Equal(expectedSeconds))
			Expect(result.Nanos).To(Equal(expectedNanos))
		},
		Entry("zero milliseconds", uint32(0), int64(0), int32(0)),
		Entry("single millisecond", uint32(1), int64(0), int32(1_000_000)),
		Entry("1 second", uint32(1000), int64(1), int32(0)),
		Entry("100ms", uint32(100), int64(0), int32(100_000_000)),
		Entry("1.5 seconds", uint32(1500), int64(1), int32(500_000_000)),
		Entry("12.345 seconds", uint32(12345), int64(12), int32(345_000_000)),
		Entry("60.123 seconds", uint32(60123), int64(60), int32(123_000_000)),
		Entry("maximum seconds with no nanos", uint32(4294967000), int64(4294967), int32(0)),
		Entry("maximum uint32 value", uint32(4294967295), int64(4294967), int32(295_000_000)),
	)
})
