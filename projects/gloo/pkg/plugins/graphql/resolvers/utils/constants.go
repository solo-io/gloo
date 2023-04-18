package utils

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// DefaultTimeout is the timeout to apply to the resolver in case there is no
// user-defined timeout on the resolver definition or the upstream ConnectionConfig
var DefaultTimeout = durationpb.New(time.Second)
