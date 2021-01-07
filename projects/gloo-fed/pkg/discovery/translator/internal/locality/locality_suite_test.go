package locality_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLocality(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Locality Suite")
}
