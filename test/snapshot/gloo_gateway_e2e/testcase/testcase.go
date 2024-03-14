package testcase

import (
	"github.com/solo-io/skv2/contrib/pkg/output"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type (
	TranslatorAssertion func(outputs *output.Snapshot)

	TestInput struct {
		// Gloo Edge input in YAML format. Only TranslatorInput or TranslatorInputYaml should be used.
		TranslatorInputYamlFile string

		// Gloo Edge input resources. Only TranslatorInput or TranslatorInputYaml should be used.
		TranslatorInput []*client.Object
	}

	// Testcase is a simple test case for bug tests
	Testcase struct {
		When string
		It   string

		// Test input to apply before running assertion
		TestInput *TestInput

		// insert these additional objects into the snapshot
		// map of type to yaml string
		Insertions map[client.Object][]TestInput

		// assertions to run
		Assertions []TranslatorAssertion

		// Focus this test
		Focus bool
		// Skip this test
		Skip bool
	}
)

var Testcases []Testcase
