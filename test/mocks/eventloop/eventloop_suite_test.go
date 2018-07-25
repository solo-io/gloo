package eventloop_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEventloop(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Eventloop Suite")
}
