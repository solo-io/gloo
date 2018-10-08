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

func StartStatsServer(addhandlers ...func(mux *http.ServeMux, profiles map[string]string)) {
	StartStatsServerWithPort("9091", addhandlers...)
}
func StartStatsServerWithPort(port string, addhandlers ...func(mux *http.ServeMux, profiles map[string]string)) {
	logconfig := zap.NewProductionConfig()

	logger, logerr := logconfig.Build()
	contextutils.SetFallbackLogger(logger.Sugar())

	go RunGoroutineStat()

	go func() {
		mux := new(http.ServeMux)

		if logerr == nil {
			mux.Handle("/logging", logconfig.Level)
		}

		addhandlers = append(addhandlers, addPprof, addStats)

		for _, addhandler := range addhandlers {
			addhandler(mux, profileDescriptions)
		}

		// add the index
		mux.HandleFunc("/", Index)
		http.ListenAndServe("localhost:"+port, mux)
	}()
}

func addPprof(mux *http.ServeMux, profiles map[string]string) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	profiles["/debug/pprof/"] = `PProf related things:<br/>
	<a href="/debug/pprof/goroutine?debug=2">full goroutine stack dump</a>
	`
}

func addStats(mux *http.ServeMux, profiles map[string]string) {
	exporter, err := prometheus.NewExporter(prometheus.Options{})
	if err == nil {
		view.RegisterExporter(exporter)
		mux.Handle("/metrics", exporter)

		profiles["/metrics"] = "Prometheus format metrics"
	}

	zpages.Handle(mux, "/zpages")
	profiles["/zpages"] = `Tracing. See <a href="/zpages/tracez">list of spans</a>`
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

	"/logging": `View \ change the log level of the program. <br/>
	
log level:
<select id="loglevelselector">
<option value="debug">debug</option>
<option value="info">info</option>
<option value="warn">warn</option>
<option value="error">error</option>
</select>
<button onclick="setlevel(document.getElementById('loglevelselector').value)">click</button>

<script>	
function setlevel(l) {
	var xhr = new XMLHttpRequest();
	xhr.open('PUT', '/logging', true);
	xhr.setRequestHeader("Content-Type", "application/json");

	xhr.onreadystatechange = function() {
		if (this.readyState == 4 && this.status == 200) {
			var resp = JSON.parse(this.responseText);
			alert("log level set to:" + resp["level"]);
		}
	};

	xhr.send('{"level":"' + l + '"}');
}
</script>
	`,
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
