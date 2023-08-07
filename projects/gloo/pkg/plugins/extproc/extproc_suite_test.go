package extproc_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExtProc(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "ExtProc Suite")
}
