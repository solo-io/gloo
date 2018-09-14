package stats

import (
	"net/http"

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

		exporter, err := prometheus.NewExporter(prometheus.Options{})
		if err == nil {
			view.RegisterExporter(exporter)
			mux.Handle("/metrics", exporter)
		}

		zpages.Handle(mux, "/debug")
		http.ListenAndServe("localhost:9091", mux)
	}()
}
