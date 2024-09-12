package kube

import (
	k8sleaderelection "k8s.io/client-go/tools/leaderelection"
	k8smetrics "k8s.io/component-base/metrics"
)

var _ k8sleaderelection.MetricsProvider = new(prometheusMetricsProvider)

var leaderGauge = k8smetrics.NewGaugeVec(&k8smetrics.GaugeOpts{
	Name: "leader_election_leader_status",
	Help: "Gauge of if the reporting system is owner of the relevant lease, 0 indicates candidate, 1 indicates leader. 'name' is the string used to identify the lease. Please make sure to group by name.",
}, []string{"name"})


var slowpathCounter = k8smetrics.NewCounterVec(&k8smetrics.CounterOpts{
	Name: "leader_election_slowpath_exercised",
	Help: "Counter of occurrences where the slow path for Kubernetes leader election was used.",
}, []string{"name"})

func init() {
	k8sleaderelection.SetProvider(prometheusMetricsProvider{})
}

type prometheusMetricsProvider struct{}

func (prometheusMetricsProvider) NewLeaderMetric() k8sleaderelection.LeaderMetric {
	return &switchAdapter{gauge: leaderGauge, slowpath: slowpathCounter}
}

type switchAdapter struct {
	gauge    *k8smetrics.GaugeVec
	slowpath *k8smetrics.CounterVec
}

func (s *switchAdapter) SlowpathExercised(name string) {
	s.slowpath.WithLabelValues(name).Inc()
}

func (s *switchAdapter) On(name string) {
	s.gauge.WithLabelValues(name).Set(1.0)
}

func (s *switchAdapter) Off(name string) {
	s.gauge.WithLabelValues(name).Set(0.0)
}
