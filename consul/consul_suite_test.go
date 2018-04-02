package consul_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo-testing/helpers"
	"github.com/solo-io/gloo-testing/helpers/local"
)

func TestConsul(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Consul Suite")
}

var (
	consulFactory *localhelpers.ConsulFactory
	err           error
)

var _ = BeforeSuite(func() {
	consulFactory, err = localhelpers.NewConsulFactory()
	helpers.Must(err)
})

var _ = AfterSuite(func() {
	consulFactory.Clean()
})

var (
	consulInstance *localhelpers.ConsulInstance
)

var _ = BeforeEach(func() {
	consulInstance, err = consulFactory.NewConsulInstance()
	helpers.Must(err)
	err = consulInstance.Run()
	helpers.Must(err)
})

var _ = AfterEach(func() {
	consulInstance.Clean()
})
