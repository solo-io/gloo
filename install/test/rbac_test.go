package test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/go-utils/manifesttestutils"
)

var _ = Describe("RBAC Test", func() {
	var (
		testManifest TestManifest
		manifestYaml string
	)

	prepareMakefile := func(helmFlags string) {
		makefileSerializer.Lock()
		defer makefileSerializer.Unlock()

		f, err := ioutil.TempFile("", "*.yaml")
		Expect(err).NotTo(HaveOccurred())
		err = f.Close()
		Expect(err).ToNot(HaveOccurred())
		manifestYaml = f.Name()

		MustMake(".", "-C", "../..", "install/glooe-gateway.yaml", "HELMFLAGS="+helmFlags, "OUTPUT_YAML="+manifestYaml)
		testManifest = NewTestManifest(manifestYaml)
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
	// TODO(mitchdraft) test the Gloo read only UI permissions
})
