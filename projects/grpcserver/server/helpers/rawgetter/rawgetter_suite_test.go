package rawgetter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRawGetter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RawGetter Suite")
}
