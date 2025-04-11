package convert

import (
	"fmt"
	"os"

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
)

func init() {
	prometheus.MustRegister(filesMetrics)
	prometheus.MustRegister(glooConfigMetric)
	prometheus.MustRegister(gatewayAPIConfigMetrics)

	glooConfigMetric.WithLabelValues("AuthConfig").Inc()
	glooConfigMetric.WithLabelValues("RouteTable").Inc()
	glooConfigMetric.WithLabelValues("Upstream").Inc()
	glooConfigMetric.WithLabelValues("VirtualService").Inc()
	glooConfigMetric.WithLabelValues("RouteOption").Inc()
	glooConfigMetric.WithLabelValues("VirtualHostOption").Inc()
	glooConfigMetric.WithLabelValues("ListenerOption").Inc()
	glooConfigMetric.WithLabelValues("HTTPListenerOption").Inc()
	glooConfigMetric.WithLabelValues("Unknown").Inc()
	glooConfigMetric.WithLabelValues("Gateway").Inc()
	glooConfigMetric.WithLabelValues("Settings").Inc()

	gatewayAPIConfigMetrics.WithLabelValues("Gateway").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("AuthConfig").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("HTTPRoute").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("Upstream").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("RouteOption").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("VirtualHostOption").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("ListenerSets").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("ListenerOption").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("HTTPListenerOption").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("DirectResponse").Inc()
	gatewayAPIConfigMetrics.WithLabelValues("Unknown").Inc()
}

func printMetrics(output *GatewayAPIOutput) {

	//we need to save the output to metrics
	gatewayAPIConfigMetrics.WithLabelValues("Gateway").Add(float64(len(output.gatewayAPICache.Gateways)))
	gatewayAPIConfigMetrics.WithLabelValues("AuthConfig").Add(float64(len(output.gatewayAPICache.AuthConfigs)))
	gatewayAPIConfigMetrics.WithLabelValues("HTTPRoute").Add(float64(len(output.gatewayAPICache.HTTPRoutes)))
	gatewayAPIConfigMetrics.WithLabelValues("Upstream").Add(float64(len(output.gatewayAPICache.Upstreams)))
	gatewayAPIConfigMetrics.WithLabelValues("RouteOption").Add(float64(len(output.gatewayAPICache.RouteOptions)))
	gatewayAPIConfigMetrics.WithLabelValues("VirtualHostOption").Add(float64(len(output.gatewayAPICache.VirtualHostOptions)))
	gatewayAPIConfigMetrics.WithLabelValues("ListenerSets").Add(float64(len(output.gatewayAPICache.ListenerSets)))
	gatewayAPIConfigMetrics.WithLabelValues("HTTPListenerOption").Add(float64(len(output.gatewayAPICache.HTTPListenerOptions)))
	gatewayAPIConfigMetrics.WithLabelValues("DirectResponse").Add(float64(len(output.gatewayAPICache.DirectResponses)))
	gatewayAPIConfigMetrics.WithLabelValues("ListenerOption").Add(float64(len(output.gatewayAPICache.ListenerOptions)))
	gatewayAPIConfigMetrics.WithLabelValues("Unknown").Add(float64(len(output.gatewayAPICache.YamlObjects)))

	metrics, _ := prometheus.DefaultGatherer.Gather()
	_, err := fmt.Fprintf(os.Stdout, "-------------------------------------\n")
	if err != nil {
		return
	}
	for _, m := range metrics {
		if m.GetName() == "gloo_config_count" {
			var count float64
			for _, t := range m.GetMetric() {
				if t.GetLabel()[0] != nil && t.GetCounter() != nil {
					_, _ = fmt.Fprintf(os.Stdout, "Gloo Config: Number of %s: %v\n", t.GetLabel()[0].GetValue(), t.GetCounter().GetValue()-1)
					if t.GetLabel()[0].GetValue() != "Unknown" {
						count += t.GetCounter().GetValue() - 1
					}
				}
			}
			_, _ = fmt.Fprintf(os.Stdout, "Total Gloo Config: %v\n", count)
		}
	}
	_, err = fmt.Fprintf(os.Stdout, "-------------------------------------\n")
	if err != nil {
		return
	}
	for _, m := range metrics {
		if m.GetName() == "gatewayapi_config_count" {
			var count float64
			for _, t := range m.GetMetric() {
				if t.GetLabel()[0] != nil && t.GetCounter() != nil {
					_, _ = fmt.Fprintf(os.Stdout, "Gateway API Config: Number of %s: %v\n", t.GetLabel()[0].GetValue(), t.GetCounter().GetValue()-1)
					if t.GetLabel()[0].GetValue() != "Unknown" {
						count += t.GetCounter().GetValue() - 1
					}
				}
			}
			_, _ = fmt.Fprintf(os.Stdout, "Total Gateway API Config: %v\n", count)
		}
	}
	_, err = fmt.Fprintf(os.Stdout, "-------------------------------------\n")
	if err != nil {
		return
	}
	for _, m := range metrics {
		if m.GetName() == "files_evaluated" {
			if m.GetMetric()[0] != nil && m.GetMetric()[0].GetCounter() != nil {
				_, _ = fmt.Fprintf(os.Stdout, "Files evaluated: %v\n", m.GetMetric()[0].GetCounter().GetValue())
			}
		}
	}
}
