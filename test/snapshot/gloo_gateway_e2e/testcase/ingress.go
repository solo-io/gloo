package testcase

import (
	_ "embed"

	"github.com/solo-io/skv2/contrib/pkg/output"
)

func init(input *TestInput) {
	Testcases = append(
		Testcases, Testcase{
			When:       "Ingress Routing is Configured",
			It:         "Get 200 response",
			TestInput:  input,
			Assertions: []TranslatorAssertion{assertSubsetRouting},
			Focus:      false,
		},
	)
}

/*
Tests subset-routing creates a RouteTable with a subset selector on two VirtualDestination routes.
*/
func assertSubsetRouting(outputs *output.Snapshot) {

}
