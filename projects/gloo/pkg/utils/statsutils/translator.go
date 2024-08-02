package statsutils

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	translationRecorder  *prometheus.HistogramVec
	statusSyncerRecorder *prometheus.HistogramVec
)

func init() {
	translationRecorder = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "translation_time_sec",
			Namespace: "gloo_edge",
			Help:      "how long the translator takes in seconds",
			Buckets:   []float64{0.25, 0.5, 0.75, 1, 2.5, 5, 10, 20, 30, 60},
		},
		[]string{"translator_name"},
	)

	statusSyncerRecorder = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:      "status_syncer_time_sec",
			Namespace: "gloo_edge",
			Help:      "how long the status syncer takes in seconds",
			Buckets:   []float64{0.25, 0.5, 0.75, 1, 2.5, 5, 10, 20, 30, 60},
		},
		[]string{"status_syncer_name"},
	)

	metrics.Registry.MustRegister(translationRecorder)
	metrics.Registry.MustRegister(statusSyncerRecorder)
}

type StopWatchFactory interface {
	// NewTranslatorStopWatch returns a stop watch for the given translator
	NewTranslatorStopWatch(translatorName string) StopWatch

	// NewStatusSyncerStopWatch returns a stop watch for the given status syncer
	NewStatusSyncerStopWatch(statusSyncerName string) StopWatch
}

func NewStatusSyncerStopWatch(statusSyncerName string) StopWatch {
	return &stopwatch{
		hist: statusSyncerRecorder.WithLabelValues(statusSyncerName),
	}
}

func NewTranslatorStopWatch(translatorName string) StopWatch {
	return &stopwatch{
		hist: translationRecorder.WithLabelValues(translatorName),
	}
}
