package pkg

import (
	"time"
)

type Fault interface {
}

type DelayFault interface {
	Fault
	Duration() time.Duration
}

type AbortFault interface {
	Fault
	Status() int
}

type Service interface {
	Name() string

	AddFault(f Fault)
	GetFault(experimentId ExperimentID) Fault
	RemoveAllFaults() Fault
}

type Mesh interface {
	GetServices() []Service
}

type MetricsProvider interface {
	GetMetric(name string) float64
}

type FaultSelector interface {
	// between 0.0 and 1.0
	RequestRatio() float64
	// TODO:
	// other things come to mind like specific header, or specific service instance.
}

type ExperimentID string

type Experiment interface {
	ID() ExperimentID
	// automatically stop if this condition is met
	MetricCondition() MetricCondition
	// type of fault to apply
	Fault() Fault
	// What service to target
	Service() Service
	// Select specific properties (i.e.: fail only a subset, fail only a certain ratio of the requests)
	// perhaps we can com up with a better name
	FaultSelector() FaultSelector
	// how long the expirement
	Duration() time.Duration
}

type MetricRelation int

const (
	Larger MetricRelation = iota
	Smaller
)

type MetricCondition interface {
	MetricName() string
	MetricRelation() MetricRelation
	MetricValue() float64
}

type ExperimentRunner interface {
	Run(e Experiment) ExperimentHandle
}

type ExperimentHandle interface {
	Stop()
}
