package customtypes

import (
	"fmt"
	"io"
	"strconv"
	"time"
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

// MarshalGQL implements the graphql.Marshaler interface
func (d Duration) MarshalGQL(w io.Writer) {
	timeDuration := time.Duration(d)
	// time.Duration is an int64. 4 decimal places will serve our needs
	w.Write([]byte(strconv.FormatFloat(timeDuration.Seconds(), 'f', 4, 64)))
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
