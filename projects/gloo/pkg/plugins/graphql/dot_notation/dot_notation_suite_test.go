package dot_notation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAdvancedHttp(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "DotNotation Suite")
}
