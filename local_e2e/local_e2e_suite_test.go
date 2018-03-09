package local_e2e

import (
	"testing"

	"github.com/solo-io/gloo-testing/helpers/local"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	envoyFactory *localhelpers.EnvoyFactory
	glooFactory  *localhelpers.GlooFactory
)

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

var (
	envoyInstance *localhelpers.EnvoyInstance
	glooInstance  *localhelpers.GlooInstance
)

var _ = BeforeEach(func() {
	var err error
	envoyInstance, err = envoyFactory.NewEnvoyInstance()
	Expect(err).NotTo(HaveOccurred())
	glooInstance, err = glooFactory.NewGlooInstance()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterEach(func() {
	if envoyInstance != nil {
		envoyInstance.Clean()
	}
	if glooInstance != nil {
		glooInstance.Clean()
	}
})
