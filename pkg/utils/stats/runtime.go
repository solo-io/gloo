package stats

import (
	"context"
	"runtime"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
)

var (
	MNumGoRoutines = stats.Int64("runtime/goroutines", "The number of goroutines", "1")

	GoroutinesNumView = &view.View{
		Name:        "runtime/goroutines",
		Measure:     MNumGoRoutines,
		Description: "The number of goroutines",
		Aggregation: view.Sum(),
	}
)

func init() {
	view.Register(GoroutinesNumView)
}

func RunGoroutineStat() {
	numgoroutines := int64(0)
	for {
		time.Sleep(time.Second)
		newnumgoroutines := int64(runtime.NumGoroutine())
		diff := newnumgoroutines - numgoroutines
		numgoroutines = newnumgoroutines
		if diff != 0 {
			stats.Record(context.TODO(), MNumGoRoutines.M(diff))
		}
	}
}
