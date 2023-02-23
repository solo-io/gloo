package main

import (
	"context"
	"path/filepath"

	structpb "github.com/golang/protobuf/ptypes/struct"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"

	plugins "github.com/solo-io/ext-auth-service/pkg/config/plugin"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	"github.com/solo-io/solo-kit/test/matchers"
)

var _ = Describe("Plugin verification script", func() {

	Describe("parsing manifest file", func() {

		It("can parse a correct manifest file", func() {
			pluginCfg, err := parseManifestFile(validManifest)
			Expect(err).NotTo(HaveOccurred())
			Expect(pluginCfg).NotTo(BeNil())
			Expect(pluginCfg).To(matchers.MatchProto(
				&extauth.AuthPlugin{
					Name:               "IsHeaderPresent",
					PluginFileName:     "IsHeaderPresent.so",
					ExportedSymbolName: "IsHeaderPresent",
					Config: &structpb.Struct{
						Fields: map[string]*structpb.Value{},
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
