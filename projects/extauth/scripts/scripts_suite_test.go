package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins"

	"github.com/gogo/protobuf/types"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestScripts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin Verification Script Suite")
}

var _ = Describe("Plugin verification script", func() {

	var (
		testAssetDir = os.ExpandEnv("$GOPATH/src/github.com/solo-io/solo-projects/test/extauth")

		pluginDir           = filepath.Join(testAssetDir, "plugins")
		validManifest       = filepath.Join(testAssetDir, "manifests", "valid.yaml")
		wrongNameManifest   = filepath.Join(testAssetDir, "manifests", "wrong_name.yaml")
		wrongSymbolManifest = filepath.Join(testAssetDir, "manifests", "wrong_symbol.yaml")
		malformedManifest   = filepath.Join(testAssetDir, "manifests", "malformed.yaml")
	)

	BeforeSuite(func() {
		contextutils.SetLogLevel(zap.DebugLevel)
	})

	Describe("parsing manifest file", func() {

		It("can parse a correct manifest file", func() {
			pluginCfg, err := parseManifestFile(validManifest)
			Expect(err).NotTo(HaveOccurred())
			Expect(pluginCfg).NotTo(BeNil())
			Expect(pluginCfg).To(Equal(
				&extauth.AuthPlugin{
					Name:               "IsHeaderPresent",
					PluginFileName:     "IsHeaderPresent.so",
					ExportedSymbolName: "IsHeaderPresent",
					Config: &types.Struct{
						Fields: map[string]*types.Value{},
					},
				},
			))
		})

		It("fails when parsing a malformed manifest file", func() {
			pluginCfg, err := parseManifestFile(malformedManifest)
			Expect(err).To(HaveOccurred())
			Expect(pluginCfg).To(BeNil())
		})
	})

	Describe("verify whether plugin can be loaded", func() {

		It("returns without error if plugin could be loaded", func() {
			err := verifyPlugin(context.Background(), pluginDir, validManifest)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("failures", func() {

			It("returns an error if plugin name is incorrect", func() {
				err := verifyPlugin(context.Background(), pluginDir, wrongNameManifest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(plugins.PluginFileOpenError(errors.New("")).Error()))
			})

			It("returns an error if symbol name is incorrect", func() {
				err := verifyPlugin(context.Background(), pluginDir, wrongSymbolManifest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(plugins.InvalidExportedSymbolError(errors.New("")).Error()))
			})

			It("returns an error if plugin dir is incorrect", func() {
				err := verifyPlugin(context.Background(), filepath.Join(pluginDir, "wrong"), validManifest)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(plugins.PluginFileOpenError(errors.New("")).Error()))
			})
		})
	})
})
