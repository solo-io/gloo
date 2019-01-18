package customtypes

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/solo-io/solo-kit/pkg/utils/log"
)

type Duration time.Duration

// UnmarshalGQL implements the graphql.Marshaler interface
func (d *Duration) UnmarshalGQL(v interface{}) error {
	durationStr, ok := v.(string)
	if !ok {
		return fmt.Errorf("durations must be strings")
	}

	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}
	ourDur := Duration(dur)
	*d = ourDur

	return nil
}

// MarshalGQL implements the graphql.Marshaller interface
func (d Duration) MarshalGQL(w io.Writer) {
	timeDuration := time.Duration(d)
	// time.Duration is an int64. 4 decimal places will serve our needs
	_, err := w.Write([]byte(strconv.FormatFloat(timeDuration.Seconds(), 'f', 4, 64)))
	if err != nil {
		log.Sprintf("failed to marshal Duration with value %v", d)
	}
}

func NewDuration(duration time.Duration) *Duration {
	ourDuration := Duration(duration)
	return &ourDuration
}

func (d *Duration) GetDuration() time.Duration {
	if d == nil {
		return 0
	}
	return time.Duration(*d)
}

// The GQL `Int` scalar type is converted to the Go `int` type during code generation.
// This custom type is meant to help working with unsigned integers.
type UnsignedInt uint32

func (uint *UnsignedInt) UnmarshalGQL(v interface{}) error {
	stringValue, ok := v.(string)
	if !ok {
		return fmt.Errorf("unsigned ints must be strings")
	}

	uintVal, err := strconv.ParseUint(stringValue, 0, 32)
	if err != nil {
		return fmt.Errorf("failed to convert %v to uint%v", stringValue, 32)
	}
	*uint = UnsignedInt(uintVal)
	return nil
}

func (uint UnsignedInt) MarshalGQL(w io.Writer) {
	_, err := w.Write([]byte(strconv.FormatUint(uint64(uint), 10)))
	if err != nil {
		log.Printf("failed to marshal UnsignedInt with value %v", uint)
	}
}
