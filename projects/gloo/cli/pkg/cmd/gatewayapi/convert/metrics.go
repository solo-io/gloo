package convert

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	filesMetrics = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "files_evaluated",
		},
	)
	glooConfigMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gloo_config_count",
		}, []string{"type"},
	)
	gatewayAPIConfigMetrics = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gatewayapi_config_count",
		}, []string{"type"},
	)
	unknownObjects = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "unknown_objects",
		}, []string{"kind"},
	)
)

func init() {
	prometheus.MustRegister(filesMetrics)
	prometheus.MustRegister(glooConfigMetric)
	prometheus.MustRegister(gatewayAPIConfigMetrics)
	prometheus.MustRegister(unknownObjects)
	glooConfigMetric.WithLabelValues("none").Inc()

	glooConfigMetric.WithLabelValues("AuthConfig").Inc()
	glooConfigMetric.WithLabelValues("RouteTable").Inc()
	glooConfigMetric.WithLabelValues("Upstream").Inc()
	glooConfigMetric.WithLabelValues("VirtualService").Inc()
	glooConfigMetric.WithLabelValues("RouteOption").Inc()
	glooConfigMetric.WithLabelValues("VirtualHostOption").Inc()

	gatewayAPIConfigMetrics.WithLabelValues("AuthConfig").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("HTTPRoute").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("Upstream").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("RouteOption").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("VirtualHostOption").Inc()
}

func printMetrics() {
	metrics, _ := prometheus.DefaultGatherer.Gather()

	for _, m := range metrics {
		if *m.Name == "gloo_config_count" {
			var count float64
			for _, t := range m.Metric {
				log.Printf("Gloo Config: Number of %s: %v", *t.Label[0].Value, *t.Counter.Value-1)
				count += *t.Counter.Value - 1
			}
			log.Printf("Total Gloo Config: %v", count)
		}
		if *m.Name == "gatewayapi_config_count" {
			var count float64
			for _, t := range m.Metric {
				log.Printf("Gateway API Config: Number of %s: %v", *t.Label[0].Value, *t.Counter.Value-1)
				count += *t.Counter.Value - 1
			}
			log.Printf("Total Gateway API Config: %v", count)
		}
		if *m.Name == "files_evaluated" {
			log.Printf("Files evaluated: %v", *m.Metric[0].Counter.Value)
		}
	}
}
