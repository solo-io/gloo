package queries_test

import (
	"testing"

	. "github.com/onsi/ginkgo"

	skhelpers "github.com/solo-io/solo-kit/test/helpers"
)

func TestGateway(t *testing.T) {
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	RunSpecs(t, "Gateway Suite")
}
