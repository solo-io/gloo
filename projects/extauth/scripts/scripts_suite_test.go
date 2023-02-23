package main

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/anyvendor/pkg/modutils"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"
)

func TestScripts(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Plugin Verification Script Suite")
}

var (
	testAssetDir        string
	pluginDir           string
	validManifest       string
	wrongNameManifest   string
	wrongSymbolManifest string
	malformedManifest   string
)

var _ = BeforeSuite(func() {
	contextutils.SetLogLevel(zap.DebugLevel)
	modPackageFile, err := modutils.GetCurrentModPackageFile()
	Expect(err).NotTo(HaveOccurred())
	repoPath := filepath.Dir(modPackageFile)
	testAssetDir = filepath.Join(repoPath, "test/extauth")

	pluginDir = filepath.Join(testAssetDir, "plugins")
	validManifest = filepath.Join(testAssetDir, "manifests", "valid.yaml")
	wrongNameManifest = filepath.Join(testAssetDir, "manifests", "wrong_name.yaml")
	wrongSymbolManifest = filepath.Join(testAssetDir, "manifests", "wrong_symbol.yaml")
	malformedManifest = filepath.Join(testAssetDir, "manifests", "malformed.yaml")
})
