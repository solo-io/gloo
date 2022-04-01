package checker

import (
	"sort"

	"github.com/solo-io/skv2/contrib/pkg/sets"
	v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
)

type Summary struct {
	Total    int32
	Errors   []*ResourceReport
	Warnings []*ResourceReport
}

type ResourceReport struct {
	Ref     *v1.ObjectRef
	Message string
}

func SortLists(summary *Summary) {
	sort.SliceStable(summary.Errors, func(i, j int) bool {
		return sets.Key(summary.Errors[i].Ref) < sets.Key(summary.Errors[j].Ref)
	})

	sort.SliceStable(summary.Warnings, func(i, j int) bool {
		return sets.Key(summary.Warnings[i].Ref) < sets.Key(summary.Warnings[j].Ref)
	})
}
