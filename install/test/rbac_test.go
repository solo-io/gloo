package test

import (
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/go-utils/manifesttestutils"
)

var _ = Describe("RBAC Test", func() {
	var (
		testManifest TestManifest
		manifestYaml string
	)

	Context("GlooE", func() {
		prepareTestManifest := func(customHelmArgs ...string) {
			makefileSerializer.Lock()
			defer makefileSerializer.Unlock()

			f, err := ioutil.TempFile("", "*.yaml")
			Expect(err).NotTo(HaveOccurred())

			Expect(WriteGlooETestManifest(f, customHelmArgs...)).NotTo(HaveOccurred())
			Expect(f.Close()).NotTo(HaveOccurred())

			manifestYaml = f.Name()
			testManifest = NewTestManifest(manifestYaml)
		}
		prepareMakefile := func(customHelmArgs string) {
			prepareTestManifest(strings.Split(customHelmArgs, " ")...)
		}
		Context("implementation-agnostic permissions", func() {
			It("correctly assigns permissions for single-namespace gloo", func() {
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true --set global.glooRbac.namespaced=true")
				permissions := GetGlooEServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true --set global.glooRbac.namespaced=false")
				permissions := GetGlooEServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})
		})
	})

	Context("Gloo OS with read-only UI", func() {
		prepareTestManifest := func(customHelmArgs ...string) {
			makefileSerializer.Lock()
			defer makefileSerializer.Unlock()

			f, err := ioutil.TempFile("", "*.yaml")
			Expect(err).NotTo(HaveOccurred())

			Expect(WriteGlooOsWithRoUiTestManifest(f, customHelmArgs...)).NotTo(HaveOccurred())
			Expect(f.Close()).NotTo(HaveOccurred())

			manifestYaml = f.Name()
			testManifest = NewTestManifest(manifestYaml)
		}
		prepareMakefile := func(customHelmArgs string) {
			prepareTestManifest(strings.Split(customHelmArgs, " ")...)
		}
		Context("implementation-agnostic permissions", func() {
			It("correctly assigns permissions for single-namespace gloo", func() {
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true --set global.glooRbac.namespaced=true")
				permissions := GetGlooWithReadOnlyUiServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true --set global.glooRbac.namespaced=false")
				permissions := GetGlooWithReadOnlyUiServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})
		})
	})
})
