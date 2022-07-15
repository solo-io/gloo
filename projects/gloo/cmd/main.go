package main

import (
	"context"
	"net/http"

	"github.com/solo-io/gloo/projects/gloo/pkg/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils/statuses"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/stats"
)

func main() {

	ctx := context.Background()
	statusHandler, writeChan := statuses.NewStatusHandler(ctx)
	stats.ConditionallyStartStatsServer(
		func(mux *http.ServeMux, profiles map[string]string) {
			profiles["/statuses"] = "Data about Gloo resource statuses"
			mux.Handle("/statuses", statusHandler)
		})

	if err := setup.Main(ctx, writeChan); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}
