package faultinjection

import (
	"reflect"

	envoytype "github.com/envoyproxy/go-control-plane/envoy/type"

	"testing"
	"time"

	"github.com/gogo/protobuf/types"
)

func TestToGoDuration(t *testing.T) {
	assertEqualDurations(&types.Duration{Seconds: 3}, 3*time.Second, t)
	assertEqualDurations(&types.Duration{Nanos: 3}, 3*time.Nanosecond, t)
	assertEqualDurations(&types.Duration{Seconds: 3, Nanos: 3}, 3*time.Second+3, t)
}

func assertEqualDurations(protoDuration *types.Duration, expectedGoDuration time.Duration, t *testing.T) {
	if toGoDuration(protoDuration) != expectedGoDuration {
		t.Errorf("Duration %v should have been %v.", protoDuration, expectedGoDuration)
	}
}

func TestToEnvoyPercentage(t *testing.T) {
	assertEqualPercent(1, 1000000, t)
	assertEqualPercent(50.0005, 50000500, t)
	// assertEqualPercent(50.000005, 50000005, t) cannot test for this level of precision
}

func assertEqualPercent(actual float32, expectedNumerator uint32, t *testing.T) {
	expectedPercentage := envoytype.FractionalPercent{
		Numerator:   expectedNumerator,
		Denominator: envoytype.FractionalPercent_MILLION,
	}

	actualPercentage := toEnvoyPercentage(actual)
	if !reflect.DeepEqual(expectedPercentage, *actualPercentage) {
		t.Errorf("Expected %v but got %v.", expectedPercentage, actualPercentage)
	}
}
