package consul_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/services"
)

func TestConsul(t *testing.T) {
	if os.Getenv("RUN_CONSUL_TESTS") != "1" {
		log.Printf("This test downloads and runs consul and is disabled by default. To enable, set RUN_CONSUL_TESTS=1 in your env.")
		return
	}
	RegisterFailHandler(Fail)
	log.DefaultOut = GinkgoWriter
	RunSpecs(t, "Consul Suite")
}

var (
	consulFactory  *services.ConsulFactory
	consulInstance *services.ConsulInstance
	err            error
)

var _ = BeforeSuite(func() {
	consulFactory, err = services.NewConsulFactory()
	Expect(err).NotTo(HaveOccurred())
	consulInstance, err = consulFactory.NewConsulInstance()
	Expect(err).NotTo(HaveOccurred())
	err = consulInstance.Run()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	consulInstance.Clean()
	consulFactory.Clean()
})
