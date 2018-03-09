package local_e2e

import (
	"testing"

	"github.com/solo-io/gloo-testing/helpers/local"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var envoyFactory *localhelpers.EnvoyFactory
var glooFactory *localhelpers.GlooFactory

func TestLocalE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LocalE2e Suite")
}

var _ = BeforeSuite(func() {
	var err error
	envoyFactory, err = localhelpers.NewEnvoyFactory()
	Expect(err).NotTo(HaveOccurred())
	glooFactory, err = localhelpers.NewGlooFactory()
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	envoyFactory.Clean()
	glooFactory.Clean()
})
