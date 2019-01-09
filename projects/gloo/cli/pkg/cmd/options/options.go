package options

import (
	"sort"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"
)

var RateLimit_TimeUnits = func() []string {
	var vals []string
	for _, name := range ratelimit.RateLimit_Unit_name {
		vals = append(vals, name)
	}
	sort.Strings(vals)
	return vals
}()

type RateLimit struct {
	Enable              bool
	TimeUnit            string
	RequestsPerTimeUnit uint32
}
