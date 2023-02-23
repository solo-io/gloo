package failover_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFailover(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Failover Suite")
}
