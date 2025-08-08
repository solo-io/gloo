package duration

import "google.golang.org/protobuf/types/known/durationpb"

// Convert milliseconds to durationpb.Duration
func MillisToDuration(millis uint32) *durationpb.Duration {
	nanos := millis % 1000 * 1_000_000
	seconds := millis / 1000
	return &durationpb.Duration{
		Seconds: int64(seconds),
		Nanos:   int32(nanos),
	}
}
