package stats

import (
	"sort"
	"text/template"

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

		mux.HandleFunc("/", Index)

		http.ListenAndServe("localhost:9091", mux)
	}()
}

func Index(w http.ResponseWriter, r *http.Request) {

	type profile struct {
		Name string
		Href string
		Desc string
	}
	var profiles []profile

	// Adding other profiles exposed from within this package
	for p, pd := range profileDescriptions {
		profiles = append(profiles, profile{
			Name: p,
			Href: p,
			Desc: pd,
		})
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	indexTmpl.Execute(w, profiles)
}

var profileDescriptions = map[string]string{
	"/debug/pprof/": `PProf related things:<br/>
	<a href="/debug/pprof/goroutine?debug=2">full goroutine stack dump</a>
	`,
	"/zpages":  `Tracing. See <a href="/zpages/tracez">list of spans</a>`,
	"/logging": "View \\ change the log level of the program",
	"/metrics": "Prometheus format metrics",
}

var indexTmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html><html>
<head>
<title>/debug/pprof/</title>
<style>
.profile-name{
	display:inline-block;
	width:6rem;
}
</style>
</head>
<body>
Things to do:
{{range .}}
<h2><a href={{.Href}}>{{.Name}}</a></h2>
<p>
{{.Desc}}
</p>
{{end}}
</body>
</html>
`))
