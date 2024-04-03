package example

import (
	"github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Assertion is a function that asserts a particular behavior.
type Assertion func(g gomega.Gomega)

type TestScenario interface {
	RegisterResources([]v1.Object)

	BeforeScenario()
	RunScenario()
	AfterScenario()
}

type SpecRunner struct {
	// The thing responsible for "running" a Spec (a test)

	//

	// assumes that a cluster has been created

	// Maintains configuration about how to install gloo edge
}

// ClusterContext contains the metadata about a Kubernetes cluster
// It also includes useful utilities for interacting with that cluster
type ClusterContext struct {
	// metadata about a kubernetes cluster

	// The name of the Kubernetes cluster
	Name string

	KubeContext string
}

type InfrastructureProvider interface {
}

// 1. Cluster is created

// 2.
