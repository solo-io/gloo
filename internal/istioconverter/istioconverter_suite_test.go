package istioconverter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIstioconverter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Istioconverter Suite")
}
