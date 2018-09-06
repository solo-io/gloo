package e2e_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/xds"

	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/services"
)

var (
	envoyFactory *services.EnvoyFactory
)

var _ = BeforeSuite(func() {
	var err error
	envoyFactory, err = services.NewEnvoyFactory()
	Expect(err).NotTo(HaveOccurred())
	xds.FallbackBindPort = 8080

})

var _ = AfterSuite(func() {
	envoyFactory.Clean()
})

func TestE2e(t *testing.T) {
	helpers.RegisterCommonFailHandlers()
	// RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}
