package stats

import (
	"net/http"

	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
)

func StartStatsServer() {
	go func() {

		mux := new(http.ServeMux)

		exporter, err := prometheus.NewExporter(prometheus.Options{})
		if err == nil {
			view.RegisterExporter(exporter)
			mux.Handle("/prom", exporter)
		}

		zpages.Handle(mux, "/debug")
		http.ListenAndServe("localhost:9091", mux)
	}()
}
