package stats

import (
	"net/http"
	"net/http/pprof"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/zpages"
	"go.uber.org/zap"
)

func StartStatsServer() {
	logconfig := zap.NewProductionConfig()

	logger, logerr := logconfig.Build()
	contextutils.SetFallbackLogger(logger.Sugar())

	go func() {

		mux := new(http.ServeMux)
		if logerr == nil {
			mux.Handle("/logging", logconfig.Level)
		}

		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		exporter, err := prometheus.NewExporter(prometheus.Options{})
		if err == nil {
			view.RegisterExporter(exporter)
			mux.Handle("/metrics", exporter)
		}

		zpages.Handle(mux, "/zpages")
		http.ListenAndServe("localhost:9091", mux)
	}()
}
