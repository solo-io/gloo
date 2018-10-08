package customtypes

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"
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
	d = &ourDur
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

type Struct map[string]interface{}

func NewStruct(protoStruct *types.Struct) *Struct {
	if protoStruct == nil {
		return nil
	}
	m, err := protoutils.MarshalMap(protoStruct)
	if err != nil {
		panic(err)
	}
	ourStruct := Struct(m)
	return &ourStruct
}

// UnmarshalGQL implements the graphql.Marshaler interface
func (s *Struct) UnmarshalGQL(v interface{}) error {
	if v == nil {
		s = nil
		return nil
	}
	mapStruct, ok := v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("structs must be map[string]interface")
	}

	ourStruct := Struct(mapStruct)

	s = &ourStruct
	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (s Struct) MarshalGQL(w io.Writer) {
	json.NewEncoder(w).Encode(s)
}

func (s *Struct) GetStruct() *types.Struct {
	if s == nil {
		return nil
	}
	var protoStruct types.Struct
	if err := protoutils.UnmarshalMap(map[string]interface{}(*s), &protoStruct); err != nil {
		panic(err)
	}
	return &protoStruct
}
