package summarize

import (
	"sort"

	"github.com/solo-io/skv2/contrib/pkg/sets"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
)

// SortLists sorts the lists inside a GlooInstanceSpec_Check_Summary to ensure idempotence
func SortLists(summary *types.GlooInstanceSpec_Check_Summary) {
	sort.SliceStable(summary.Errors, func(i, j int) bool {
		return sets.Key(summary.Errors[i].GetRef()) < sets.Key(summary.Errors[j].GetRef())
	})

	sort.SliceStable(summary.Warnings, func(i, j int) bool {
		return sets.Key(summary.Warnings[i].GetRef()) < sets.Key(summary.Warnings[j].GetRef())
	})
}
