package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/go-utils/manifesttestutils"
)

var _ = Describe("RBAC Test", func() {
	var (
		testManifest TestManifest
		manifestYaml string
	)

	AfterEach(func() {
		if manifestYaml != "" {
			err := os.Remove(manifestYaml)
			Expect(err).ToNot(HaveOccurred())
		}
	})

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
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true," +
					"global.glooRbac.namespaced=true")
				permissions := GetGlooEServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true," +
					"global.glooRbac.namespaced=false")
				permissions := GetGlooEServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})

			It("creates no permissions when rbac is disabled", func() {
				prepareMakefile(fmt.Sprintf("--namespace %s --set global.glooRbac.create=false,"+
					"prometheus.rbac.create=false,grafana.testFramework.enabled=false", namespace))
				contents, err := ioutil.ReadFile(manifestYaml)
				Expect(err).NotTo(HaveOccurred(), "should be able to read manifest file")

				Expect(strings.ToLower(string(contents))).NotTo(ContainSubstring("rbac.authorization.k8s.io"), "should not have any reference to the rbac api group")
			})
		})

		Context("all cluster-scoped RBAC resources", func() {
			checkSuffix := func(suffix string) {
				rbacResources := testManifest.SelectResources(func(resource *unstructured.Unstructured) bool {
					return (resource.GetKind() == "ClusterRole" || resource.GetKind() == "ClusterRoleBinding") &&
						!strings.Contains(resource.GetName(), "glooe-grafana") &&
						!strings.Contains(resource.GetName(), "glooe-prometheus")
				})

				Expect(rbacResources.NumResources()).NotTo(BeZero())

				rbacResources.ExpectAll(func(resource *unstructured.Unstructured) {
					Expect(resource.GetName()).To(HaveSuffix("-" + suffix))
				})
			}

			It("is all named appropriately when a custom suffix is specified", func() {
				suffix := "test-suffix"
				prepareMakefile("--namespace " + namespace + " --set global.glooRbac.nameSuffix=" + suffix)
				checkSuffix(suffix)
			})

			It("is all named appropriately in a non-namespaced install", func() {
				prepareMakefile("--namespace " + namespace)
				checkSuffix(namespace)
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
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true," +
					"global.glooRbac.namespaced=true")
				permissions := GetGlooWithReadOnlyUiServiceAccountPermissions("gloo-system")
				testManifest.ExpectPermissions(permissions)
			})

			It("correctly assigns permissions for cluster-scoped gloo", func() {
				prepareMakefile("--namespace " + namespace + " --set namespace.create=true," +
					"global.glooRbac.namespaced=false")
				permissions := GetGlooWithReadOnlyUiServiceAccountPermissions("")
				testManifest.ExpectPermissions(permissions)
			})
		})
	})
})
