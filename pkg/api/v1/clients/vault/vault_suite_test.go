package vault_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/test/services"
)

func TestVault(t *testing.T) {
	RegisterFailHandler(Fail)
	log.DefaultOut = GinkgoWriter
	RunSpecs(t, "Vault Suite")
}

var (
	vaultFactory  *services.VaultFactory
	vaultInstance *services.VaultInstance
	err           error
)

var _ = BeforeSuite(func() {
	vaultFactory, err = services.NewVaultFactory()
	Expect(err).NotTo(HaveOccurred())
	vaultInstance, err = vaultFactory.NewVaultInstance()
	Expect(err).NotTo(HaveOccurred())
	err = vaultInstance.Run()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	vaultInstance.Clean()
	vaultFactory.Clean()
})
