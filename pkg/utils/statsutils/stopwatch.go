package statsutils

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/solo-io/go-utils/contextutils"
)

type StopWatch interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

type stopwatch struct {
	hist        prometheus.Observer
	inProgress  prometheus.Gauge
	currentTime *time.Time
}

func NewStopWatch(hist prometheus.Observer) StopWatch {
	return &stopwatch{hist: hist}
}

// Start begins the counter for the recorded value.
func (s *stopwatch) Start(ctx context.Context) {
	if s.currentTime != nil {
		contextutils.LoggerFrom(ctx).DPanic("Start() should not be called twice without calling Stop()")
		return
	}
	now := time.Now()
	s.currentTime = &now
	if s.inProgress != nil {
		s.inProgress.Inc()
	}
}

// Stop ends the counter for the recorded value. This method should be called directly after
// Start(), in a defer statement.
func (s *stopwatch) Stop(ctx context.Context) {
	if s.currentTime == nil {
		contextutils.LoggerFrom(ctx).DPanic("Stop() should only be called after Start()")
		return
	}
	s.hist.Observe(time.Since(*s.currentTime).Seconds())
	s.currentTime = nil
	if s.inProgress != nil {
		s.inProgress.Dec()
	}
}
