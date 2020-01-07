package plugins

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/utils/modutils"
)

func TestPlugins(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExtAuth Plugin Loader Suite")
}

var (
	pluginFileDir string
)

var _ = BeforeSuite(func() {
	modPackageFile, err := modutils.GetCurrentModPackageFile()
	Expect(err).NotTo(HaveOccurred())
	repoPath := filepath.Dir(modPackageFile)
	pluginFileDir = filepath.Join(repoPath, "test/extauth/plugins")
})
