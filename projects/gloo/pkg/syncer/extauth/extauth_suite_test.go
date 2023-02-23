package extauth_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
)

const (
	// All resources will live in a single namespace to simplify the tests
	writeNamespace = defaults.GlooSystem
)

func TestExtAuth(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "ExtAuth Suite")
}
