package failover_scheme_handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRpcFailoverScheme(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FailoverScheme Suite")
}
